package service

import (
	"context"
	"errors"
	"short-links/internal/links/dto"
	"short-links/internal/utils"
	"strings"
)

type linkRepo interface {
	CreateLink(ctx context.Context, userID int, link, code string) (dto.Link, error)
	GetLinks(ctx context.Context, userID int) ([]dto.Link, error)
}

type LinkService struct {
	repo linkRepo
}

func NewLinkService(repo linkRepo) *LinkService {
	return &LinkService{
		repo,
	}
}

func (l *LinkService) CreateLink(ctx context.Context, userID int, link string) (dto.Link, error) {
	if link == "" {
		return dto.Link{}, errors.New("link cannot be empty")
	}

	if !strings.HasPrefix(link, "http") {
		return dto.Link{}, errors.New("link must be only http or https")
	}

	code, err := utils.GenerateCode(7)
	if err != nil {
		return dto.Link{}, err
	}

	createdLink, err := l.repo.CreateLink(ctx, userID, link, code)
	if err != nil {
		return dto.Link{}, err
	}

	return createdLink, nil
}

func (l *LinkService) GetLinks(ctx context.Context, userID int) ([]dto.Link, error) {
	links, err := l.repo.GetLinks(ctx, userID)
	if err != nil {
		return nil, err
	}

	return links, nil
}
