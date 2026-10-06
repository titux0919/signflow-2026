package handlers

import (
	"net/http"

	"signflow-2026/internal/config"
	"signflow-2026/internal/web"
)

type Handlers struct {
	cfg config.Config
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	render(w, r, http.StatusOK, web.Home())
}

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func (h *Handlers) About(w http.ResponseWriter, r *http.Request) {
	render(w, r, http.StatusOK, web.About())
}
