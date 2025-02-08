package server

import (
	"github.com/alserok/url_shortener/internal/server/grpc"
	"github.com/alserok/url_shortener/internal/server/http"
	"github.com/alserok/url_shortener/internal/service"
	"github.com/alserok/url_shortener/pkg/logger"
)

type Server interface {
	MustServe(port string)
}

const (
	GRPC = iota
	HTTP
)

func New(t uint, srvc service.Service, log logger.Logger) Server {
	switch t {
	case HTTP:
		return http.New(srvc, log)
	case GRPC:
		return grpc.New(srvc, log)
	default:
		panic("invalid server type")
	}
}
