package server

import (
	"context"
	"github.com/alserok/url_shortener/internal/cache"
	"github.com/alserok/url_shortener/internal/server/grpc"
	"github.com/alserok/url_shortener/internal/server/http"
	"github.com/alserok/url_shortener/internal/service"
	"github.com/alserok/url_shortener/pkg/logger"
)

type Server interface {
	MustServe(port string)
	Shutdown(ctx context.Context)
}

const (
	GRPC = iota
	HTTP
)

func New(t uint, srvc service.Service, cache cache.Cache, log logger.Logger) Server {
	switch t {
	case HTTP:
		return http.New(srvc, cache, log)
	case GRPC:
		return grpc.New(srvc, cache, log)
	default:
		panic("invalid server type")
	}
}
