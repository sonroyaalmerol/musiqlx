package authentication

import "github.com/gorilla/mux"

func Router(r *mux.Router) {
	// Authentication endpoints
	r.HandleFunc("/login", login).Methods("POST")
	r.HandleFunc("/logout", logout).Methods("GET")
}
