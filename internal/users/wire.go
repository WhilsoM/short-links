package users

import (
	"short-links/internal/users/handler"
	"short-links/internal/users/repository"
	"short-links/internal/users/service"
	"short-links/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Init(r chi.Router, db *pgxpool.Pool, jwtmanager utils.JWTManagerInterface) {
	repo := repository.NewUsersRepo(db)
	svc := service.NewUsersService(repo, jwtmanager)
	h := handler.NewUsersHandler(svc)

	r.Route("/users", func(r chi.Router) {
		r.Post("/register", h.RegisterUser)
		r.Post("/login", h.LoginUser)
	})
}
