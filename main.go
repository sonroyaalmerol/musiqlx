package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sonroyaalmerol/musiqlx/pages"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		page := pages.Hello("World")
		page.Render(r.Context(), w)
	})

	http.ListenAndServe(":3000", r)
}
