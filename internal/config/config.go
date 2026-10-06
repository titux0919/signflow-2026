package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Port        string
	Env         string
	BaseURL     string
}

func (c Config) IsProd() bool {
	return c.Env == "prod"
}

func Load() (Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		DatabaseURL: getenv(
			"DATABASE_URL",
			"postgres://postgres@localhost:5432/signflow?sslmode=disable",
		),
		Port:    getenv("PORT", "8080"),
		Env:     getenv("APP_ENV", "dev"),
		BaseURL: os.Getenv("BASE_URL"),
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = "http://localhost:" + cfg.Port
	}

	if cfg.IsProd() && os.Getenv("BASE_URL") == "" {
		return Config{}, fmt.Errorf("BASE_URL is required in production")
	}

	return cfg, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
