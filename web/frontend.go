package web

import (
	"embed"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
)

var (
	//go:embed dist/*
	dist embed.FS

	//go:embed dist/index.html
	indexHTML embed.FS

	distDirFS     = http.FS(dist)
	distIndexHTML = http.FS(indexHTML)
)

func RegisterRoutes(r *mux.Router) {
	if os.Getenv("ENV") == "dev" {
		log.Println("Running in dev mode")
		setupDevProxy(r)
		return
	}
	// Use the static assets from the dist directory
	r.PathPrefix("/").Handler(http.FileServer(distIndexHTML))
	r.PathPrefix("/").Handler(http.FileServer(distDirFS))
}

func setupDevProxy(r *mux.Router) {
	devURL, err := url.Parse("http://localhost:8686")
	if err != nil {
		log.Fatal(err)
	}

	// Set up a reverse proxy to the vite dev server on localhost:8686
	proxy := httputil.NewSingleHostReverseProxy(devURL)
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if len(req.URL.Path) >= 4 && req.URL.Path[:4] == "/api" {
			// Skip the proxy if the path starts with /api
			http.NotFound(w, req)
			return
		}
		proxy.ServeHTTP(w, req)
	})
}
