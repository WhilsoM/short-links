package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"short-links/internal/clients/postgresql"
	"short-links/internal/config"
	"short-links/internal/users"
	"short-links/internal/utils"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.MustConfigLoad()

	ctx := context.Background()

	jwtmanager := utils.NewJWTManager(cfg.JWTSecret)

	dbpool, err := postgresql.NewPostgresClient(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Info("failed to run or ping postgresql", "error", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		users.Init(r, dbpool, jwtmanager)
	})

	if err := http.ListenAndServe(cfg.Port, r); err != nil {
		slog.Info("server shutdowned", "error", err, "port", cfg.Port)
		os.Exit(1)
	}
}
