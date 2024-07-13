package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sonroyaalmerol/musiqlx/internal/pkg/router"
	"github.com/sonroyaalmerol/musiqlx/web"
)

func main() {
	r := mux.NewRouter()

	router.RegisterRoutes(r)
	web.RegisterRoutes(r)

	log.Fatal(http.ListenAndServe(":8686", r))
}
