package http

import (
	"context"
	"errors"
	"github.com/alserok/url_shortener/internal/cache"
	"github.com/alserok/url_shortener/internal/service"
	"github.com/alserok/url_shortener/pkg/logger"
	"net/http"
	"time"
)

func New(srvc service.Service, cache cache.Cache, log logger.Logger) *server {
	return &server{
		s: &http.Server{
			WriteTimeout:      time.Second * 3,
			ReadHeaderTimeout: time.Second * 1,
		},
		log: log,
		handler: handler{
			srvc:  srvc,
			cache: cache,
		},
	}
}

type server struct {
	s       *http.Server
	handler handler
	log     logger.Logger
}

func (s *server) Shutdown(ctx context.Context) {
	_ = s.s.Shutdown(ctx)
}

func (s *server) MustServe(port string) {
	mux := http.NewServeMux()
	s.setupRoutes(mux, s.handler)

	s.s.Handler = mux
	s.s.Addr = ":" + port

	if err := s.s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic("failed to listen and serve: " + err.Error())
	}
}
