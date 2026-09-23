package handlers

import (
	"net/http"
	"signflow-2026/internal/config"
)

type Handlers struct {
	cfg config.Config
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte(`<!doctype html>
<title>SignFlow</title>
<h1>✍️ SignFlow</h1>
<p>The skeleton is up and serving.</p>`))
}

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func (h *Handlers) About(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<!doctype html>
<title>About SignFlow</title>
<h1>About SignFlow</h1>
<p>SignFlow is a document signing application.</p>`))
}
