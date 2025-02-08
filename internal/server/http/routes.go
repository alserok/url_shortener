package http

import (
	"github.com/alserok/url_shortener/internal/server/http/middleware"
	"github.com/alserok/url_shortener/internal/utils"
	"net/http"
	"time"
)

func (s *server) setupRoutes(mux *http.ServeMux, h handler) {
	middleware.WithLogger(s.log)(s.s.Handler)
	middleware.WithRateLimiter(utils.NewLimiter(100, time.Second))(s.s.Handler)
	middleware.WithRecovery(s.s.Handler)

	mux.HandleFunc("GET /get", middleware.WithErrorHandler(h.GetURL))
	mux.HandleFunc("POST /save", middleware.WithErrorHandler(h.ShortenAndSaveURL))
}
