package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	migrations "signflow-2026/db"
	"signflow-2026/internal/config"
	"signflow-2026/internal/db"
	"signflow-2026/internal/handlers"
	"signflow-2026/internal/session"
	"signflow-2026/static"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	if err := run(log); err != nil {
		log.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run(log *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	log.Info("starting signflow",
		"env", cfg.Env,
		"port", cfg.Port,
		"base_url", cfg.BaseURL,
	)

	ctx := context.Background()

	if err := runMigrations(cfg.DatabaseURL, log); err != nil {
		return err
	}

	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	queries := db.New(pool)

	// Create the session manager and pass all dependencies to handlers.
	sessions := session.NewManager(queries, cfg.IsProd())
	h := handlers.New(cfg, queries, sessions, log)

	staticFS, err := fs.Sub(static.FS, "assets")
	if err != nil {
		return err
	}

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           h.Router(staticFS),
		ReadHeaderTimeout: 10 * time.Second,
	}

	shutdownCtx, stop := signal.NotifyContext(
		ctx,
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		log.Info("listening", "addr", srv.Addr)

		if err := srv.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			log.Error("server error", "err", err)
			stop()
		}
	}()

	<-shutdownCtx.Done()
	log.Info("shutting down")

	timeout, cancel := context.WithTimeout(
		context.Background(),
		15*time.Second,
	)
	defer cancel()

	return srv.Shutdown(timeout)
}

func runMigrations(databaseURL string, log *slog.Logger) error {
	sqlDB, err := goose.OpenDBWithDriver("pgx", databaseURL)
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	goose.SetBaseFS(migrations.FS)
	goose.SetLogger(gooseLogger{log})

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	return goose.Up(sqlDB, "migrations")
}

type gooseLogger struct {
	log *slog.Logger
}

func (g gooseLogger) Fatalf(format string, v ...interface{}) {
	g.log.Error("goose", "msg", fmt.Sprintf(format, v...))
}

func (g gooseLogger) Printf(format string, v ...interface{}) {
	g.log.Info("goose", "msg", fmt.Sprintf(format, v...))
}
