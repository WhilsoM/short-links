package handler

import (
	"context"
	"log/slog"
	"net/http"
	custommiddlewares "short-links/internal/custom-middlewares"
	"short-links/internal/links/dto"
	"short-links/internal/utils"

	"github.com/go-chi/chi/v5"
)

type linkService interface {
	CreateLink(ctx context.Context, userID int, link string) (dto.Link, error)
	GetLinks(ctx context.Context, userID int) ([]dto.Link, error)
	GetLinkByCode(ctx context.Context, code string) (string, error)
}

type LinkHandler struct {
	svc linkService
}

func NewLinkHandler(svc linkService) *LinkHandler {
	return &LinkHandler{
		svc,
	}
}

func (l *LinkHandler) CreateLink(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateLinkRequest

	if err := utils.DecodeJSON(r.Body, &req); err != nil {
		utils.WriteJSONErrorResponse(w, http.StatusBadRequest, "invalid json")
		return
	}

	userID, err := custommiddlewares.GetUserIDFromContext(r.Context())
	if err != nil {
		slog.Info("failed to get user id from context", "error", err)
		utils.WriteJSONErrorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	link, err := l.svc.CreateLink(r.Context(), userID, req.Link)
	if err != nil {
		slog.Info("failed to create link", "error", err, "user_id", userID, "link", req.Link)
		utils.WriteJSONErrorResponse(w, http.StatusInternalServerError, "failed to create link")
		return
	}

	res := dto.CreateLinkResponse{
		ID:        link.ID,
		ShortLink: "http://localhost:8080/r/" + link.Code,
		CreatedAt: link.CreatedAt,
	}

	utils.WriteJSONResponse(w, http.StatusCreated, res)
}

func (l *LinkHandler) GetLinks(w http.ResponseWriter, r *http.Request) {
	userID, err := custommiddlewares.GetUserIDFromContext(r.Context())
	if err != nil {
		slog.Info("failed to get user id from context", "error", err)
		utils.WriteJSONErrorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	links, err := l.svc.GetLinks(r.Context(), userID)
	if err != nil {
		slog.Info("failed to create link", "error", err, "user_id", userID)
		utils.WriteJSONErrorResponse(w, http.StatusInternalServerError, "failed to create link")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, links)
}

func (l *LinkHandler) RedirectLink(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	link, err := l.svc.GetLinkByCode(r.Context(), code)
	if err != nil {
		slog.Info("failed to get link by code", "code", code, "error", err)
		utils.WriteJSONErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	http.Redirect(w, r, link, http.StatusFound)
}
