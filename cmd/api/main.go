package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"short-links/internal/clients/postgresql"
	"short-links/internal/config"
	custommiddlewares "short-links/internal/custom-middlewares"
	"short-links/internal/links"
	"short-links/internal/users"
	"short-links/internal/utils"

	"github.com/go-chi/chi/middleware"
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

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		users.Init(r, dbpool, jwtmanager)

		r.Group(func(r chi.Router) {
			r.Use(custommiddlewares.AuthMiddleware(jwtmanager))

			links.Init(r, dbpool)
		})
	})

	slog.Info("server started", "port", cfg.Port)

	if err := http.ListenAndServe(cfg.Port, r); err != nil {
		slog.Info("server shutdowned", "error", err, "port", cfg.Port)
		os.Exit(1)
	}
}
