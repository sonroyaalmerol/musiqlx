package main

import (
	"flag"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sonroyaalmerol/musiqlx/components"
	"github.com/sonroyaalmerol/musiqlx/pages"
	"github.com/sonroyaalmerol/musiqlx/utils"
)

var count = 0

func main() {
	// Parse command-line flags
	isDevMode := flag.Bool("dev", false, "Enable development mode")
	flag.Parse()

	if *isDevMode {
		utils.Dev()
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		page := pages.Index(count)
		page.Render(r.Context(), w)
	})

	r.Post("/add", func(w http.ResponseWriter, r *http.Request) {
		count++
		counter := components.Counter(count)
		counter.Render(r.Context(), w)
	})

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.ListenAndServe(":3000", r)
}
