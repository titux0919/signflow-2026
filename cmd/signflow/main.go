package main

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"

	"signflow-2026/internal/config"
	"signflow-2026/internal/db"
	"signflow-2026/internal/handlers"
	"signflow-2026/static"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	queries := db.New(pool)

	h := &handlers.Handlers{
		Queries: queries,
	}

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
