package http

import (
	"encoding/json"
	"fmt"
	"github.com/alserok/url_shortener/internal/service"
	"github.com/alserok/url_shortener/internal/service/models"
	"github.com/alserok/url_shortener/internal/utils"
	"net/http"
)

type handler struct {
	srvc service.Service
}

func (h *handler) ShortenAndSaveURL(w http.ResponseWriter, r *http.Request) error {
	var url models.URL
	if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
		return utils.NewError(err.Error(), utils.BadRequestErr)
	}

	shortened, err := h.srvc.ShortenAndSaveURL(r.Context(), url.OriginURL)
	if err != nil {
		return fmt.Errorf("failed to shorten and save url: %w", err)
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(map[string]any{"shortened": shortened}); err != nil {
		return utils.NewError(err.Error(), utils.InternalErr)
	}

	w.WriteHeader(http.StatusCreated)

	return nil
}

func (h *handler) GetURL(w http.ResponseWriter, r *http.Request) error {
	shortened := r.URL.Query().Get("shortened")

	url, err := h.srvc.GetURL(r.Context(), shortened)
	if err != nil {
		return fmt.Errorf("failed to shorten and save url: %w", err)
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(map[string]any{"url": url}); err != nil {
		return utils.NewError(err.Error(), utils.InternalErr)
	}

	return nil
}
