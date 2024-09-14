//nolint:lll,gocyclo,forbidigo,nilerr,errcheck
package main

import (
	"context"
	"errors"
	"expvar"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"net/url"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"

	// avatar encode/decode
	_ "image/gif"
	_ "image/png"

	"github.com/google/shlex"
	"github.com/gorilla/securecookie"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/sentriz/gormstore"
	"golang.org/x/sync/errgroup"

	"github.com/sonroyaalmerol/musiqlx"
	"github.com/sonroyaalmerol/musiqlx/db"
	"github.com/sonroyaalmerol/musiqlx/handlerutil"
	"github.com/sonroyaalmerol/musiqlx/infocache/albuminfocache"
	"github.com/sonroyaalmerol/musiqlx/infocache/artistinfocache"
	"github.com/sonroyaalmerol/musiqlx/jukebox"
	"github.com/sonroyaalmerol/musiqlx/lastfm"
	"github.com/sonroyaalmerol/musiqlx/listenbrainz"
	"github.com/sonroyaalmerol/musiqlx/playlist"
	"github.com/sonroyaalmerol/musiqlx/scrobble"
	"github.com/sonroyaalmerol/musiqlx/server/ctrlsubsonic"
	"github.com/sonroyaalmerol/musiqlx/transcode"
	"go.senan.xyz/flagconf"
)

func main() {
	confListenAddr := flag.String("listen-addr", "0.0.0.0:4747", "listen address (optional)")

	confTLSCert := flag.String("tls-cert", "", "path to TLS certificate (optional)")
	confTLSKey := flag.String("tls-key", "", "path to TLS private key (optional)")

	confCachePath := flag.String("cache-path", "", "path to cache")

	var confMusicPaths pathAliases
	flag.Var(&confMusicPaths, "music-path", "path to music")

	confPlaylistsPath := flag.String("playlists-path", "", "path to your list of new or existing m3u playlists that musiqlx can manage")

	confDBPath := flag.String("db-path", "musiqlx.db", "path to database (optional)")

	confJukeboxEnabled := flag.Bool("jukebox-enabled", false, "whether the subsonic jukebox api should be enabled (optional)")
	confJukeboxMPVExtraArgs := flag.String("jukebox-mpv-extra-args", "", "extra command line arguments to pass to the jukebox mpv daemon (optional)")

	confProxyPrefix := flag.String("proxy-prefix", "", "url path prefix to use if behind proxy. eg '/musiqlx' (optional)")
	confHTTPLog := flag.Bool("http-log", true, "http request logging (optional)")

	confShowVersion := flag.Bool("version", false, "show musiqlx version")
	confConfigPath := flag.String("config-path", "", "path to config (optional)")

	confExcludePattern := flag.String("exclude-pattern", "", "regex pattern to exclude files from scan (optional)")

	confPprof := flag.Bool("pprof", false, "enable the /debug/pprof endpoint (optional)")
	confExpvar := flag.Bool("expvar", false, "enable the /debug/vars endpoint (optional)")

	flag.Parse()
	flagconf.ParseEnv()
	flagconf.ParseConfig(*confConfigPath)

	if *confShowVersion {
		fmt.Printf("v%s\n", musiqlx.Version)
		os.Exit(0)
	}

	if _, err := regexp.Compile(*confExcludePattern); err != nil {
		log.Fatalf("invalid exclude pattern: %v\n", err)
	}

	if len(confMusicPaths) == 0 {
		log.Fatalf("please provide a music directory")
	}

	var err error
	for i, confMusicPath := range confMusicPaths {
		if confMusicPaths[i].path, err = validatePath(confMusicPath.path); err != nil {
			log.Fatalf("checking music dir %q: %v", confMusicPath.path, err)
		}
	}

	if *confCachePath, err = validatePath(*confCachePath); err != nil {
		log.Fatalf("checking cache directory: %v", err)
	}
	if *confPlaylistsPath, err = validatePath(*confPlaylistsPath); err != nil {
		log.Fatalf("checking playlist directory: %v", err)
	}

	cacheDirAudio := path.Join(*confCachePath, "audio")
	cacheDirCovers := path.Join(*confCachePath, "covers")
	if err := os.MkdirAll(cacheDirAudio, os.ModePerm); err != nil {
		log.Fatalf("couldn't create audio cache path: %v\n", err)
	}
	if err := os.MkdirAll(cacheDirCovers, os.ModePerm); err != nil {
		log.Fatalf("couldn't create covers cache path: %v\n", err)
	}

	dbc, err := db.New(*confDBPath, db.DefaultOptions())
	if err != nil {
		log.Fatalf("error opening database: %v\n", err)
	}
	defer dbc.Close()

	err = dbc.Migrate(db.MigrationContext{
		Production:        true,
		DBPath:            *confDBPath,
		OriginalMusicPath: confMusicPaths[0].path,
		PlaylistsPath:     *confPlaylistsPath,
	})
	if err != nil {
		log.Panicf("error migrating database: %v\n", err)
	}

	var musicPaths []ctrlsubsonic.MusicPath
	for _, pa := range confMusicPaths {
		musicPaths = append(musicPaths, ctrlsubsonic.MusicPath{Alias: pa.alias, Path: pa.path})
	}

	proxyPrefixExpr := regexp.MustCompile(`^\/*(.*?)\/*$`)
	*confProxyPrefix = proxyPrefixExpr.ReplaceAllString(*confProxyPrefix, `/$1`)

	log.Printf("starting musiqlx v%s\n", musiqlx.Version)
	log.Printf("provided config\n")
	flag.VisitAll(func(f *flag.Flag) {
		value := strings.ReplaceAll(f.Value.String(), "\n", "")
		log.Printf("    %-25s %s\n", f.Name, value)
	})

	transcoder := transcode.NewCachingTranscoder(
		transcode.NewFFmpegTranscoder(),
		cacheDirAudio,
	)

	lastfmClientKeySecretFunc := func() (string, string, error) {
		apiKey, _ := dbc.GetSetting(db.LastFMAPIKey)
		secret, _ := dbc.GetSetting(db.LastFMSecret)
		if apiKey == "" || secret == "" {
			return "", "", fmt.Errorf("not configured")
		}
		return apiKey, secret, nil
	}

	listenbrainzClient := listenbrainz.NewClient()
	lastfmClient := lastfm.NewClient(lastfmClientKeySecretFunc)

	playlistStore, err := playlist.NewStore(*confPlaylistsPath)
	if err != nil {
		log.Panicf("error creating playlists store: %v", err)
	}

	var jukebx *jukebox.Jukebox
	if *confJukeboxEnabled {
		jukebx = jukebox.New()
	}

	sessKey, err := dbc.GetSetting("session_key")
	if err != nil {
		log.Panicf("error getting session key: %v\n", err)
	}
	if sessKey == "" {
		sessKey = string(securecookie.GenerateRandomKey(32))
		if err := dbc.SetSetting("session_key", sessKey); err != nil {
			log.Panicf("error setting session key: %v\n", err)
		}
	}
	sessDB := gormstore.New(dbc.DB, []byte(sessKey))
	sessDB.SessionOpts.HttpOnly = true
	sessDB.SessionOpts.SameSite = http.SameSiteLaxMode

	artistInfoCache := artistinfocache.New(dbc, lastfmClient)
	albumInfoCache := albuminfocache.New(dbc, lastfmClient)

	scrobblers := []scrobble.Scrobbler{lastfmClient, listenbrainzClient}

	resolveProxyPath := func(in string) string {
		url, _ := url.Parse(in)
		url.Path = path.Join(*confProxyPrefix, url.Path)
		return url.String()
	}

	ctrlSubsonic, err := ctrlsubsonic.New(dbc, musicPaths, cacheDirAudio, cacheDirCovers, jukebx, playlistStore, scrobblers, transcoder, lastfmClient, artistInfoCache, albumInfoCache, resolveProxyPath)
	if err != nil {
		log.Panicf("error creating subsonic controller: %v\n", err)
	}

	chain := handlerutil.Chain()
	if *confHTTPLog {
		chain = handlerutil.Chain(handlerutil.Log)
	}
	chain = handlerutil.Chain(
		chain,
		handlerutil.BasicCORS,
	)
	trim := handlerutil.TrimPathSuffix(".view") // /x.view and /x should match the same

	mux := http.NewServeMux()
	mux.Handle("/rest/", http.StripPrefix("/rest", chain(trim(ctrlSubsonic))))
	mux.Handle("/ping", chain(handlerutil.Message("ok")))
	mux.Handle("/", chain(http.RedirectHandler(resolveProxyPath("/admin/home"), http.StatusSeeOther)))

	if *confExpvar {
		mux.Handle("/debug/vars", expvar.Handler())
		expvar.Publish("stats", expvar.Func(func() any {
			stats, _ := dbc.Stats()
			return stats
		}))
	}

	var (
		readTimeout  = 5 * time.Second
		writeTimeout = 5 * time.Second
		idleTimeout  = 5 * time.Second
	)

	if *confPprof {
		// overwrite global WriteTimeout. in future we should set this only for these handlers
		// https://github.com/golang/go/issues/62358
		writeTimeout = 0

		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	errgrp, ctx := errgroup.WithContext(ctx)

	errgrp.Go(func() error {
		defer logJob("http")()

		server := &http.Server{
			Addr:        *confListenAddr,
			ReadTimeout: readTimeout, WriteTimeout: writeTimeout, IdleTimeout: idleTimeout,
			Handler: mux,
		}
		errgrp.Go(func() error {
			<-ctx.Done()
			return server.Shutdown(context.Background())
		})
		if *confTLSCert != "" && *confTLSKey != "" {
			return server.ListenAndServeTLS(*confTLSCert, *confTLSKey)
		}
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	errgrp.Go(func() error {
		if jukebx == nil {
			return nil
		}

		defer logJob("jukebox")()

		extraArgs, _ := shlex.Split(*confJukeboxMPVExtraArgs)
		jukeboxTempDir := filepath.Join(*confCachePath, "musiqlx-jukebox")
		if err := os.RemoveAll(jukeboxTempDir); err != nil {
			return fmt.Errorf("remove jubebox tmp dir: %w", err)
		}
		if err := os.MkdirAll(jukeboxTempDir, os.ModePerm); err != nil {
			return fmt.Errorf("create tmp sock file: %w", err)
		}
		sockPath := filepath.Join(jukeboxTempDir, "sock")
		if err := jukebx.Start(ctx, sockPath, extraArgs); err != nil {
			return fmt.Errorf("start jukebox: %w", err)
		}
		return nil
	})

	errgrp.Go(func() error {
		defer logJob("session clean")()

		ctxTick(ctx, 10*time.Minute, func() {
			sessDB.Cleanup()
		})
		return nil
	})

	errgrp.Go(func() error {
		if _, _, err := lastfmClientKeySecretFunc(); err != nil {
			return nil
		}

		defer logJob("refresh artist info")()

		ctxTick(ctx, 8*time.Second, func() {
			if err := artistInfoCache.Refresh(); err != nil {
				log.Printf("error in artist info cache: %v", err)
			}
		})
		return nil
	})

	if err := errgrp.Wait(); err != nil {
		log.Panic(err)
	}

	fmt.Println("shutdown complete")
}

const pathAliasSep = "->"

type (
	pathAliases []pathAlias
	pathAlias   struct{ alias, path string }
)

func (pa pathAliases) String() string {
	var strs []string
	for _, p := range pa {
		if p.alias != "" {
			strs = append(strs, fmt.Sprintf("%s %s %s", p.alias, pathAliasSep, p.path))
			continue
		}
		strs = append(strs, p.path)
	}
	return strings.Join(strs, ", ")
}

func (pa *pathAliases) Set(value string) error {
	if name, path, ok := strings.Cut(value, pathAliasSep); ok {
		*pa = append(*pa, pathAlias{alias: name, path: path})
		return nil
	}
	*pa = append(*pa, pathAlias{path: value})
	return nil
}

func validatePath(p string) (string, error) {
	if p == "" {
		return "", errors.New("path can't be empty")
	}
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return "", errors.New("path does not exist, please provide one")
	}
	p, err := filepath.Abs(p)
	if err != nil {
		return "", fmt.Errorf("make absolute: %w", err)
	}
	return p, nil
}

func logJob(jobName string) func() {
	log.Printf("starting job %q", jobName)
	return func() { log.Printf("stopped job %q", jobName) }
}

func ctxTick(ctx context.Context, interval time.Duration, f func()) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			f()
		}
	}
}
