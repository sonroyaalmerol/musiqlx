package artist

import "github.com/gorilla/mux"

func Router(r *mux.Router) {
	r.HandleFunc("/api/v1/artist/{id}", getArtist).Methods("GET")
	r.HandleFunc("/api/v1/artist/{id}", updateArtist).Methods("PUT")
	r.HandleFunc("/api/v1/artist/{id}", deleteArtist).Methods("DELETE")
	r.HandleFunc("/api/v1/artist", getArtists).Methods("GET")
	r.HandleFunc("/api/v1/artist", createArtist).Methods("POST")
	r.HandleFunc("/api/v1/artist/editor", updateArtistEditor).Methods("PUT")
	r.HandleFunc("/api/v1/artist/editor", deleteArtistEditor).Methods("DELETE")
	r.HandleFunc("/api/v1/artist/lookup", artistLookup).Methods("GET")
}
