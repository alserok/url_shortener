package http

import (
	"encoding/json"
	"fmt"
	"github.com/alserok/url_shortener/internal/cache"
	"github.com/alserok/url_shortener/internal/service"
	"github.com/alserok/url_shortener/internal/service/models"
	"github.com/alserok/url_shortener/internal/utils"
	"github.com/alserok/url_shortener/pkg/logger"
	"net/http"
)

type handler struct {
	srvc service.Service

	cache cache.Cache
}

func (h *handler) ShortenAndSaveURL(w http.ResponseWriter, r *http.Request) error {
	var url models.URL
	if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
		return utils.NewError(err.Error(), utils.BadRequestErr)
	}
	if url.OriginURL == "" {
		return utils.NewError("invalid url", utils.BadRequestErr)
	}

	log := logger.ExtractLogger(r.Context())

	log.Debug("started ShortenAndSaveURL handler", logger.WithArg("url", url.OriginURL))

	shortened, err := h.srvc.ShortenAndSaveURL(r.Context(), url.OriginURL)
	if err != nil {
		return fmt.Errorf("failed to shorten and save url: %w", err)
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(map[string]any{"shortenedURL": shortened}); err != nil {
		return utils.NewError(err.Error(), utils.InternalErr)
	}

	log.Debug("successfully finished ShortenAndSaveURL handler", logger.WithArg("url", url.OriginURL))

	return nil
}

func (h *handler) GetURL(w http.ResponseWriter, r *http.Request) error {
	log := logger.ExtractLogger(r.Context())

	shortened := r.PathValue("shortenedURL")

	log.Debug("started GetURL handler", logger.WithArg("shortened_url", shortened))

	if cachedURL, err := h.cache.Get(r.Context(), shortened); err == nil {
		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(map[string]any{"url": cachedURL}); err != nil {
			return utils.NewError(err.Error(), utils.InternalErr)
		}

		log.Debug("returned cached GetURL handler response", logger.WithArg("shortened_url", shortened))

		return nil
	}

	url, err := h.srvc.GetURL(r.Context(), shortened)
	if err != nil {
		return fmt.Errorf("failed to shorten and save url: %w", err)
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(map[string]any{"originURL": url}); err != nil {
		return utils.NewError(err.Error(), utils.InternalErr)
	}

	if err = h.cache.Set(r.Context(), shortened, url); err != nil {
		logger.ExtractLogger(r.Context()).Warn("failed to insert in cache", logger.WithArg("error", err.Error()))
	}

	log.Debug("successfully finished GetURL handler", logger.WithArg("shortened_url", shortened))

	return nil
}
