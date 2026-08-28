package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"short-links/internal/links/dto"
	"short-links/internal/utils"
	"strings"

	"github.com/segmentio/kafka-go"
)

type linkRepo interface {
	CreateLink(ctx context.Context, userID int, link, code string) (dto.Link, error)
	GetLinks(ctx context.Context, userID int) ([]dto.Link, error)
	GetLinkByCode(ctx context.Context, code string) (string, error)
}

type linkCache interface {
	GetLinkByCode(ctx context.Context, code string) (string, error)
	SetLink(ctx context.Context, code, original_link string) error
}

type LinkService struct {
	repo   linkRepo
	cache  linkCache
	broker *kafka.Writer
}

func NewLinkService(repo linkRepo, cache linkCache, broker *kafka.Writer) *LinkService {
	return &LinkService{
		repo,
		cache,
		broker,
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

	if err := l.cache.SetLink(ctx, createdLink.Code, link); err != nil {
		slog.Info("failed to set link in cache", "error", err, "code", createdLink.Code)
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

func (l *LinkService) GetLinkByCode(ctx context.Context, code string) (string, error) {
	if code == "" {
		return "", errors.New("code cannot be empty")
	}

	linkFromCache, err := l.cache.GetLinkByCode(ctx, code)
	if err != nil {
		slog.Info("failed to take original link from cache", "error", err, "code", code)
	}

	if err == nil {
		err = l.broker.WriteMessages(ctx, kafka.Message{
			Value: fmt.Appendf([]byte(""), "%s", linkFromCache),
		})
		if err != nil {
			slog.Info("failed to send a message to broker messages", "error", err)
		}
		slog.Info("successfuly send a message to broker messages")

		slog.Info("from cache hit!", "link", linkFromCache)
		return linkFromCache, nil
	}
	slog.Info("cache doesn't exist")

	originalUrl, err := l.repo.GetLinkByCode(ctx, code)
	if err != nil {
		return "", err
	}

	if err := l.cache.SetLink(ctx, code, originalUrl); err != nil {
		slog.Info("failed to set link in cache", "error", err, "code", code)
	}

	err = l.broker.WriteMessages(ctx, kafka.Message{
		Value: fmt.Appendf([]byte(""), "%s", linkFromCache),
	})
	if err != nil {
		slog.Info("failed to send a message to broker messages", "error", err)
	}
	slog.Info("successfuly send a message to broker messages")

	return originalUrl, nil
}
