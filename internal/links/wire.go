package links

import (
	"short-links/internal/links/handler"
	"short-links/internal/links/repository"
	"short-links/internal/links/service"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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

func NewBuildHandler(db *pgxpool.Pool) *handler.LinkHandler {
	repo := repository.NewLinkRepo(db)
	svc := service.NewLinkService(repo)

	return handler.NewLinkHandler(svc)
}
