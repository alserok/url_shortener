package middleware

import (
	"github.com/alserok/url_shortener/pkg/logger"
	"net/http"
)

func WithRecovery(fn http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.ExtractLogger(r.Context()).Error("panic recovery", logger.WithArg("error", err))
			}
		}()

		fn.ServeHTTP(w, r)
	}
}
