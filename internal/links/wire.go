package links

import (
	"short-links/internal/links/handler"
	"short-links/internal/links/repository"
	"short-links/internal/links/service"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Init(r chi.Router, db *pgxpool.Pool) {
	repo := repository.NewLinkRepo(db)
	svc := service.NewLinkService(repo)
	h := handler.NewLinkHandler(svc)

	r.Route("/links", func(r chi.Router) {
		r.Post("/", h.CreateLink)
	})
}
