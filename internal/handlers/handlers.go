package handlers

import (
	"net/http"

	"signflow-2026/internal/config"
	"signflow-2026/internal/db"
	"signflow-2026/internal/web"
)

type Handlers struct {
	cfg     config.Config
	Queries *db.Queries
}

func New(cfg config.Config, q *db.Queries) *Handlers {
	return &Handlers{
		cfg:     cfg,
		Queries: q,
	}
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	count, err := h.Queries.CountUsers(r.Context())
	if err != nil {
		println("COUNT USERS ERROR:", err.Error())
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	render(w, r, http.StatusOK, web.Home(count))
}

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func (h *Handlers) About(w http.ResponseWriter, r *http.Request) {
	render(w, r, http.StatusOK, web.About())
}
