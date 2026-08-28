package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"short-links/internal/clients/kafka"
	"short-links/internal/clients/postgresql"
	"short-links/internal/clients/redis"
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

	redisClient := redis.NewRedisClient(cfg.RedisAddr)
	kafkaClient := kafka.NewKafkaClient(cfg.KafkaAddr, cfg.KafkaAnalyticTopic)
	defer kafkaClient.Close()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	linkHandler := links.NewBuildHandler(dbpool, redisClient, kafkaClient)

	links.RegisterPublic(r, linkHandler)

	r.Route("/api", func(r chi.Router) {
		users.Init(r, dbpool, jwtmanager)

		r.Group(func(r chi.Router) {
			r.Use(custommiddlewares.AuthMiddleware(jwtmanager))

			links.RegisterProtected(r, linkHandler)
		})
	})

	slog.Info("server started", "port", cfg.Port)

	if err := http.ListenAndServe(cfg.Port, r); err != nil {
		slog.Info("server shutdowned", "error", err, "port", cfg.Port)
		os.Exit(1)
	}
}
