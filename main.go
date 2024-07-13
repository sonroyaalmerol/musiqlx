package main

import (
	"log"
	"net/http"

	"github.com/sonroyaalmerol/musiqlx/internal/pkg/router"
)

func main() {
	api := router.ApiRouter()

	log.Fatal(http.ListenAndServe(":8080", api))
}
