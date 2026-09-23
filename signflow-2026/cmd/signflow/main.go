package main

import (
	"fmt"
	"net/http"

	"signflow-2026/internal/config"
	"signflow-2026/internal/handlers"
)

func main() {
	cfg, err := config.Load()

	if err != nil {
		panic(err)
	}

	h := &handlers.Handlers{}
	router := handlers.NewRouter(h)

	fmt.Println("Starting SignFlow")
	fmt.Println("Port:", cfg.Port)
	fmt.Println("Environment:", cfg.Env)

	http.ListenAndServe(":"+cfg.Port, router)
}
