package config

import "os"

type Config struct {
	Port    string
	Env     string
	BaseURL string
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}

func Load() (Config, error) {
	cfg := Config{
		Port:    getenv("PORT", "8080"),
		Env:     getenv("APP_ENV", "dev"),
		BaseURL: os.Getenv("BASE_URL"),
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = "http://localhost:" + cfg.Port
	}

	return cfg, nil
}
