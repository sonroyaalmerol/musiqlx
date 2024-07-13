package album

import "github.com/gorilla/mux"

func Router(r *mux.Router) {
	r.HandleFunc("/api/v1/album", getAlbums).Methods("GET")
	r.HandleFunc("/api/v1/album", createAlbum).Methods("POST")
	r.HandleFunc("/api/v1/album/{id}", updateAlbum).Methods("PUT")
	r.HandleFunc("/api/v1/album/{id}", deleteAlbum).Methods("DELETE")
	r.HandleFunc("/api/v1/album/{id}", getAlbum).Methods("GET")
	r.HandleFunc("/api/v1/album/monitor", monitorAlbum).Methods("PUT")
	r.HandleFunc("/api/v1/album/lookup", albumLookup).Methods("GET")
	r.HandleFunc("/api/v1/albumstudio", createAlbumStudio).Methods("POST")
}
