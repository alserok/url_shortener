package middleware

import (
	"github.com/alserok/url_shortener/pkg/logger"
	"net/http"
)

func WithLogger(log logger.Logger) func(handlerFunc http.Handler) http.HandlerFunc {
	return func(fn http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(logger.WrapLogger(r.Context(), log))
			fn.ServeHTTP(w, r)
		}
	}
}
