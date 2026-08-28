package links

import (
	"short-links/internal/links/cache"
	"short-links/internal/links/handler"
	"short-links/internal/links/repository"
	"short-links/internal/links/service"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

func RegisterProtected(r chi.Router, h *handler.LinkHandler) {
	r.Route("/links", func(r chi.Router) {
		r.Post("/", h.CreateLink)
		r.Get("/", h.GetLinks)
	})
}

func RegisterPublic(r chi.Router, h *handler.LinkHandler) {
	r.Get("/r/{code}", h.RedirectLink)
}

func NewBuildHandler(db *pgxpool.Pool, redis *redis.Client, broker *kafka.Writer) *handler.LinkHandler {
	cache := cache.NewLinksCache(redis)
	repo := repository.NewLinkRepo(db)
	svc := service.NewLinkService(repo, cache, broker)

	return handler.NewLinkHandler(svc)
}
