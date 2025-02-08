package grpc

import (
	"context"
	"fmt"
	"github.com/alserok/url_shortener/internal/cache"
	"github.com/alserok/url_shortener/internal/service"
	"github.com/alserok/url_shortener/pkg/logger"
	"github.com/alserok/url_shortener/pkg/proto"
)

type handler struct {
	proto.UnimplementedURLShortenerServer

	srvc service.Service

	cache cache.Cache
}

func (h *handler) ShortenAndSaveURL(ctx context.Context, url *proto.URL) (*proto.ShortenedURL, error) {
	shortened, err := h.srvc.ShortenAndSaveURL(ctx, url.OriginUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to shorten and save url: %w", err)
	}

	return &proto.ShortenedURL{ShortenedUrl: shortened}, nil
}

func (h *handler) GetURL(ctx context.Context, url *proto.ShortenedURL) (*proto.URL, error) {
	if cachedURL, err := h.cache.Get(ctx, url.ShortenedUrl); err == nil {
		return &proto.URL{OriginUrl: cachedURL}, nil
	}

	originURL, err := h.srvc.GetURL(ctx, url.ShortenedUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get url: %w", err)
	}

	if err = h.cache.Set(ctx, url.ShortenedUrl, originURL); err != nil {
		logger.ExtractLogger(ctx).Warn("failed to insert in cache", logger.WithArg("error", err.Error()))
	}

	return &proto.URL{OriginUrl: originURL}, nil
}
