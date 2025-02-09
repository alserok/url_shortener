package middleware

import (
	"github.com/alserok/url_shortener/internal/utils"
	"net/http"
)

func WithRateLimiter(limiter utils.Limiter) func(fn http.Handler) http.HandlerFunc {
	return func(fn http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow(r.Context()) {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}

			fn.ServeHTTP(w, r)
		}
	}
}
