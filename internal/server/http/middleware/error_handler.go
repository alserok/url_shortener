package middleware

import (
	"encoding/json"
	"github.com/alserok/url_shortener/internal/utils"
	"net/http"
)

func WithErrorHandler(fn func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			msg, code := utils.FromError(r.Context(), err)

			switch code {
			case utils.BadRequestErr:
				w.WriteHeader(http.StatusBadRequest)
			case utils.NotFoundErr:
				w.WriteHeader(http.StatusNotFound)
			case utils.InternalErr:
				w.WriteHeader(http.StatusInternalServerError)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}

			_ = json.NewEncoder(w).Encode(map[string]any{"error": msg})
		}
	}
}
