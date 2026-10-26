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

	// Load the user from the session for every request.
	r.Use(h.LoadUser)

	// Static files.
	r.Handle("/static/*",
		http.StripPrefix("/static/",
			http.FileServer(http.FS(staticFS))))

	// Public routes.
	r.Get("/", h.Home)
	r.Get("/healthz", h.Health)
	r.Get("/about", h.About)

	// Registration.
	r.Get("/register", h.RegisterForm)
	r.Post("/register", h.Register)

	// Login.
	r.Get("/login", h.LoginForm)
	r.Post("/login", h.Login)

	// Logout.
	r.Post("/logout", h.Logout)

	// Protected routes.
	r.Group(func(r chi.Router) {
		r.Use(h.RequireAuth)

		r.Get("/dashboard", h.Dashboard)
		r.Get("/account", h.Account)
	})

	return r
}
