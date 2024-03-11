package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sonroyaalmerol/musiqlx/pages"
)

var count = 0

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		page := pages.Index(count)
		page.Render(r.Context(), w)
	})

	r.Post("/add", func(w http.ResponseWriter, r *http.Request) {
		count++
		page := pages.Index(count)
		page.Render(r.Context(), w)
	})

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.ListenAndServe(":3000", r)
}
