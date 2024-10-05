package main

import (
	"log"
	"net/http"

	"github.com/sonroyaalmerol/musiqlx/config"
	v1 "github.com/sonroyaalmerol/musiqlx/internal/api/v1"
	"github.com/sonroyaalmerol/musiqlx/pkg/db"
)

func main() {
	config.LoadConfig() // Load configuration
	db.ConnectDB()      // Initialize the database

	mux := v1.NewRouter() // Set up the router and handlers

	log.Println("Starting OpenSubsonic API server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
