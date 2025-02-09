package middleware

import (
	"encoding/json"
	"github.com/alserok/url_shortener/internal/utils"
	"github.com/alserok/url_shortener/pkg/logger"
	"net/http"
)

func WithErrorHandler(next func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := next(w, r); err != nil {
			msg, code := utils.FromError(r.Context(), err)

			switch code {
			case utils.BadRequestErr:
				w.WriteHeader(http.StatusBadRequest)
			case utils.NotFoundErr:
				w.WriteHeader(http.StatusNotFound)
			case utils.InternalErr:
				logger.ExtractLogger(r.Context()).Error(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
			default:
				logger.ExtractLogger(r.Context()).Error(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
			}

			_ = json.NewEncoder(w).Encode(map[string]any{"error": msg})
		}
	}
}
