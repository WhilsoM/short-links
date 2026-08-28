package handler

import (
	"context"
	"log/slog"
	"net/http"
	custommiddlewares "short-links/internal/custom-middlewares"
	"short-links/internal/links/dto"
	"short-links/internal/utils"
)

type LinkService interface {
	CreateLink(ctx context.Context, userID int, link string) (dto.Link, error)
}

type LinkHandler struct {
	svc LinkService
}

func NewLinkHandler(svc LinkService) *LinkHandler {
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
		ShortLink: "http://localhost:8080/" + link.Code,
		CreatedAt: link.CreatedAt,
	}

	utils.WriteJSONResponse(w, http.StatusCreated, res)

}
