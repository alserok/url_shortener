package http

import (
	"github.com/alserok/url_shortener/internal/server/http/middleware"
	"github.com/alserok/url_shortener/internal/utils"
	"net/http"
	"time"
)

func (s *server) setupRoutes(mux *http.ServeMux, h handler) {
	mux.HandleFunc("GET /get/{shortenedURL}", s.wrapHandlerWithMiddlewares(h.GetURL))
	mux.HandleFunc("POST /save", s.wrapHandlerWithMiddlewares(h.ShortenAndSaveURL))
}

func (s *server) wrapHandlerWithMiddlewares(handlerFn func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	fn := middleware.WithLogger(s.log)(
		middleware.WithRecovery(
			middleware.WithRateLimiter(utils.NewLimiter(1000, time.Second))(
				middleware.WithErrorHandler(handlerFn),
			),
		),
	)

	return fn
}
