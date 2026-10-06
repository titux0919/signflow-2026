package handlers

import (
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (h *Handlers) Router(staticFS fs.FS) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Handle("/static/*",
		http.StripPrefix("/static/",
			http.FileServer(http.FS(staticFS))))

	r.Get("/", h.Home)
	r.Get("/healthz", h.Health)
	r.Get("/about", h.About)

	return r
}
