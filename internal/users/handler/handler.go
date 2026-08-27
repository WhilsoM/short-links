package handler

import (
	"context"
	"net/http"
	"short-links/internal/users/dto"
	"short-links/internal/utils"
)

type UsersService interface {
	RegisterUser(ctx context.Context, user dto.User) (string, error)
	LoginUser(ctx context.Context, user dto.User) (string, error)
}

type UsersHandler struct {
	svc UsersService
}

func NewUsersHandler(svc UsersService) *UsersHandler {
	return &UsersHandler{
		svc,
	}
}

func (u *UsersHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterUserRequest
	if err := utils.DecodeJSON(r.Body, &req); err != nil {
		utils.WriteJSONErrorResponse(w, http.StatusBadRequest, "invalid payload")
		return
	}

	user := dto.User{
		Username: req.Username,
		Password: req.Password,
	}

	tokens, err := u.svc.RegisterUser(r.Context(), user)
	if err != nil {
		utils.WriteJSONErrorResponse(w, http.StatusBadRequest, "error: "+err.Error())
		return
	}

	res := dto.RegisterUserResponse{
		AccessToken: tokens,
	}

	utils.WriteJSONResponse(w, http.StatusCreated, res)
}

func (u *UsersHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginUserRequest
	if err := utils.DecodeJSON(r.Body, &req); err != nil {
		utils.WriteJSONErrorResponse(w, http.StatusBadRequest, "invalid payload")
		return
	}

	user := dto.User{
		Username: req.Username,
		Password: req.Password,
	}

	tokens, err := u.svc.LoginUser(r.Context(), user)
	if err != nil {
		utils.WriteJSONErrorResponse(w, http.StatusBadRequest, "error: "+err.Error())
		return
	}

	res := dto.LoginUserResponse{
		AccessToken: tokens,
	}

	utils.WriteJSONResponse(w, http.StatusCreated, res)
}
