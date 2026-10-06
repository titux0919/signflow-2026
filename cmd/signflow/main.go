package main

import (
	"fmt"
	"io/fs"
	"net/http"

	"signflow-2026/internal/config"
	"signflow-2026/internal/handlers"
	"signflow-2026/static"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	h := &handlers.Handlers{}

	staticFS, err := fs.Sub(static.FS, "assets")
	if err != nil {
		panic(err)
	}

	router := h.Router(staticFS)

	fmt.Println("Starting SignFlow")
	fmt.Println("Port:", cfg.Port)
	fmt.Println("Environment:", cfg.Env)

	http.ListenAndServe(":"+cfg.Port, router)
}
