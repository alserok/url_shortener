package grpc

import (
	"context"
	"fmt"
	"github.com/alserok/url_shortener/internal/cache"
	"github.com/alserok/url_shortener/internal/service"
	"github.com/alserok/url_shortener/internal/utils"
	"github.com/alserok/url_shortener/pkg/logger"
	"github.com/alserok/url_shortener/pkg/proto"
)

type handler struct {
	proto.UnimplementedURLShortenerServer

	srvc service.Service

	cache cache.Cache
}

func (h *handler) ShortenAndSaveURL(ctx context.Context, url *proto.URL) (*proto.ShortenedURL, error) {
	if url.OriginUrl == "" {
		return nil, utils.NewError("invalid url", utils.BadRequestErr)
	}

	log := logger.ExtractLogger(ctx)

	log.Debug("started ShortenAndSaveURL handler", logger.WithArg("url", url.OriginUrl))

	shortened, err := h.srvc.ShortenAndSaveURL(ctx, url.OriginUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to shorten and save url: %w", err)
	}

	log.Debug("successfully finished ShortenAndSaveURL handler", logger.WithArg("url", url.OriginUrl))

	return &proto.ShortenedURL{ShortenedUrl: shortened}, nil
}

func (h *handler) GetURL(ctx context.Context, url *proto.ShortenedURL) (*proto.URL, error) {
	if url.ShortenedUrl == "" {
		return nil, utils.NewError("invalid url", utils.BadRequestErr)
	}

	log := logger.ExtractLogger(ctx)

	log.Debug("started GetURL handler", logger.WithArg("shortened_url", url.ShortenedUrl))

	if cachedURL, err := h.cache.Get(ctx, url.ShortenedUrl); err == nil {
		log.Debug("returned cached GetURL handler response", logger.WithArg("shortened_url", url.ShortenedUrl))
		return &proto.URL{OriginUrl: cachedURL}, nil
	}

	originURL, err := h.srvc.GetURL(ctx, url.ShortenedUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get url: %w", err)
	}

	if err = h.cache.Set(ctx, url.ShortenedUrl, originURL); err != nil {
		log.Warn("failed to insert in cache", logger.WithArg("error", err.Error()))
	}

	log.Debug("successfully finished GetURL handler", logger.WithArg("shortened_url", url.ShortenedUrl))

	return &proto.URL{OriginUrl: originURL}, nil
}
