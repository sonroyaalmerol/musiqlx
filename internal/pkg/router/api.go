package router

import (
	"github.com/gorilla/mux"
	"github.com/sonroyaalmerol/musiqlx/internal/app/album"
	"github.com/sonroyaalmerol/musiqlx/internal/pkg/middlewares"
)

func RegisterRoutes(r *mux.Router) {
	// Apply the middleware for API key validation
	r.Use(middlewares.ApiKeyMiddleware)

	album.Router(r)

	// ApiInfo endpoints
	r.HandleFunc("/api", getApiInfo).Methods("GET")

	// StaticResource endpoints
	r.HandleFunc("/login", staticLogin).Methods("GET")
	r.HandleFunc("/content/{path}", getContent).Methods("GET")
	r.HandleFunc("/", getIndex).Methods("GET")
	r.HandleFunc("/{path}", getPath).Methods("GET")

	// AutoTagging endpoints
	r.HandleFunc("/api/v1/autotagging/{id}", getAutoTagging).Methods("GET")
	r.HandleFunc("/api/v1/autotagging/{id}", updateAutoTagging).Methods("PUT")
	r.HandleFunc("/api/v1/autotagging/{id}", deleteAutoTagging).Methods("DELETE")
	r.HandleFunc("/api/v1/autotagging", createAutoTagging).Methods("POST")
	r.HandleFunc("/api/v1/autotagging", listAutoTagging).Methods("GET")
	r.HandleFunc("/api/v1/autotagging/schema", getAutoTaggingSchema).Methods("GET")

	// Backup endpoints
	r.HandleFunc("/api/v1/system/backup", getBackups).Methods("GET")
	r.HandleFunc("/api/v1/system/backup/{id}", deleteBackup).Methods("DELETE")
	r.HandleFunc("/api/v1/system/backup/restore/{id}", restoreBackup).Methods("POST")
	r.HandleFunc("/api/v1/system/backup/restore/upload", uploadRestoreBackup).Methods("POST")

	// Blocklist endpoints
	r.HandleFunc("/api/v1/blocklist", getBlocklist).Methods("GET")
	r.HandleFunc("/api/v1/blocklist/{id}", deleteBlocklist).Methods("DELETE")
	r.HandleFunc("/api/v1/blocklist/bulk", bulkDeleteBlocklist).Methods("DELETE")

	// Calendar endpoints
	r.HandleFunc("/api/v1/calendar", getCalendar).Methods("GET")
	r.HandleFunc("/api/v1/calendar/{id}", getCalendarByID).Methods("GET")

	// CalendarFeed endpoints
	r.HandleFunc("/feed/v1/calendar/lidarr.ics", getCalendarFeed).Methods("GET")

	// Command endpoints
	r.HandleFunc("/api/v1/command/{id}", getCommand).Methods("GET")
	r.HandleFunc("/api/v1/command/{id}", deleteCommand).Methods("DELETE")
	r.HandleFunc("/api/v1/command", createCommand).Methods("POST")
	r.HandleFunc("/api/v1/command", getCommands).Methods("GET")

	// CustomFilter endpoints
	r.HandleFunc("/api/v1/customfilter/{id}", getCustomFilter).Methods("GET")
	r.HandleFunc("/api/v1/customfilter/{id}", updateCustomFilter).Methods("PUT")
	r.HandleFunc("/api/v1/customfilter/{id}", deleteCustomFilter).Methods("DELETE")
	r.HandleFunc("/api/v1/customfilter", getCustomFilters).Methods("GET")
	r.HandleFunc("/api/v1/customfilter", createCustomFilter).Methods("POST")

	// CustomFormat endpoints
	r.HandleFunc("/api/v1/customformat/{id}", getCustomFormat).Methods("GET")
	r.HandleFunc("/api/v1/customformat/{id}", updateCustomFormat).Methods("PUT")
	r.HandleFunc("/api/v1/customformat/{id}", deleteCustomFormat).Methods("DELETE")
	r.HandleFunc("/api/v1/customformat", getCustomFormats).Methods("GET")
	r.HandleFunc("/api/v1/customformat", createCustomFormat).Methods("POST")
	r.HandleFunc("/api/v1/customformat/schema", getCustomFormatSchema).Methods("GET")

	// Cutoff endpoints
	r.HandleFunc("/api/v1/wanted/cutoff", getCutoffs).Methods("GET")
	r.HandleFunc("/api/v1/wanted/cutoff/{id}", getCutoff).Methods("GET")

	// DelayProfile endpoints
	r.HandleFunc("/api/v1/delayprofile", createDelayProfile).Methods("POST")
	r.HandleFunc("/api/v1/delayprofile", getDelayProfiles).Methods("GET")
	r.HandleFunc("/api/v1/delayprofile/{id}", deleteDelayProfile).Methods("DELETE")
	r.HandleFunc("/api/v1/delayprofile/{id}", updateDelayProfile).Methods("PUT")
	r.HandleFunc("/api/v1/delayprofile/{id}", getDelayProfile).Methods("GET")
	r.HandleFunc("/api/v1/delayprofile/reorder/{id}", reorderDelayProfile).Methods("PUT")

	// DiskSpace endpoints
	r.HandleFunc("/api/v1/diskspace", getDiskSpace).Methods("GET")

	// DownloadClient endpoints
	r.HandleFunc("/api/v1/downloadclient/{id}", getDownloadClient).Methods("GET")
	r.HandleFunc("/api/v1/downloadclient/{id}", updateDownloadClient).Methods("PUT")
	r.HandleFunc("/api/v1/downloadclient/{id}", deleteDownloadClient).Methods("DELETE")
	r.HandleFunc("/api/v1/downloadclient", getDownloadClients).Methods("GET")
	r.HandleFunc("/api/v1/downloadclient", createDownloadClient).Methods("POST")
	r.HandleFunc("/api/v1/downloadclient/bulk", bulkUpdateDownloadClient).Methods("PUT")
	r.HandleFunc("/api/v1/downloadclient/bulk", bulkDeleteDownloadClient).Methods("DELETE")
	r.HandleFunc("/api/v1/downloadclient/schema", getDownloadClientSchema).Methods("GET")
	r.HandleFunc("/api/v1/downloadclient/test", testDownloadClient).Methods("POST")
	r.HandleFunc("/api/v1/downloadclient/testall", testAllDownloadClients).Methods("POST")
	r.HandleFunc("/api/v1/downloadclient/action/{name}", actionDownloadClient).Methods("POST")

	// DownloadClientConfig endpoints
	r.HandleFunc("/api/v1/config/downloadclient/{id}", getDownloadClientConfig).Methods("GET")
	r.HandleFunc("/api/v1/config/downloadclient/{id}", updateDownloadClientConfig).Methods("PUT")
	r.HandleFunc("/api/v1/config/downloadclient", getDownloadClientConfigs).Methods("GET")

	// FileSystem endpoints
	r.HandleFunc("/api/v1/filesystem", getFileSystems).Methods("GET")
	r.HandleFunc("/api/v1/filesystem/type", getFileSystemTypes).Methods("GET")
	r.HandleFunc("/api/v1/filesystem/mediafiles", getFileSystemMediaFiles).Methods("GET")

	// Health endpoints
	r.HandleFunc("/api/v1/health", getHealth).Methods("GET")

	// History endpoints
	r.HandleFunc("/api/v1/history", getHistory).Methods("GET")
	r.HandleFunc("/api/v1/history/since", getHistorySince).Methods("GET")
	r.HandleFunc("/api/v1/history/artist", getArtistHistory).Methods("GET")
	r.HandleFunc("/api/v1/history/failed/{id}", markHistoryFailed).Methods("POST")

	// HostConfig endpoints
	r.HandleFunc("/api/v1/config/host/{id}", getHostConfig).Methods("GET")
	r.HandleFunc("/api/v1/config/host/{id}", updateHostConfig).Methods("PUT")
	r.HandleFunc("/api/v1/config/host", getHostConfigs).Methods("GET")
}
