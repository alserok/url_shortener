package service

import (
	"context"
	"fmt"
	"github.com/alserok/url_shortener/internal/db"
	"github.com/alserok/url_shortener/internal/utils"
	"github.com/alserok/url_shortener/pkg/logger"
)

type Service interface {
	ShortenAndSaveURL(ctx context.Context, url string) (string, error)
	GetURL(ctx context.Context, shortened string) (string, error)
}

func New(repo db.Repository) *service {
	return &service{
		repo: repo,
	}
}

type service struct {
	repo db.Repository
}

const (
	shortenedLength = 10
)

func (s *service) ShortenAndSaveURL(ctx context.Context, url string) (string, error) {
	log := logger.ExtractLogger(ctx)

	log.Debug("started ShortenAndSaveURL service", logger.WithArg("url", url))

	shortened, err := utils.ShortenURL(ctx, url, shortenedLength)
	if err != nil {
		return "", fmt.Errorf("failed to shorten url: %w", err)
	}

	log.Debug("shortened url in ShortenAndSaveURL service", logger.WithArg("url", url))

	if err = s.repo.SaveURL(ctx, url, shortened); err != nil {
		return "", fmt.Errorf("failed to save url: %w", err)
	}

	log.Debug("successfully finished ShortenAndSaveURL service", logger.WithArg("url", url))

	return shortened, nil
}

func (s *service) GetURL(ctx context.Context, shortened string) (string, error) {
	log := logger.ExtractLogger(ctx)

	log.Debug("started GetURL service", logger.WithArg("shortened_url", shortened))

	url, err := s.repo.GetURL(ctx, shortened)
	if err != nil {
		return "", fmt.Errorf("failed to get url: %w", err)
	}

	log.Debug("successfully finished GetURL service", logger.WithArg("shortened_url", shortened))

	return url, nil
}
