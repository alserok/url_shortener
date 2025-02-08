package grpc

import (
	"github.com/alserok/url_shortener/internal/server/grpc/middleware"
	"github.com/alserok/url_shortener/internal/service"
	"github.com/alserok/url_shortener/internal/utils"
	"github.com/alserok/url_shortener/pkg/logger"
	"github.com/alserok/url_shortener/pkg/proto"
	"google.golang.org/grpc"
	"net"
	"time"
)

func New(srvc service.Service, log logger.Logger) *server {
	return &server{
		s: grpc.NewServer(
			middleware.WithChain(
				middleware.WithLogger(log),
				middleware.WithRecovery(),
				middleware.WithErrorHandler(),
				middleware.WithRateLimiter(utils.NewLimiter(100, time.Second)),
			),
		),
		handler: handler{
			srvc: srvc,
		},
	}
}

type server struct {
	s       *grpc.Server
	handler handler
}

func (s *server) MustServe(port string) {
	proto.RegisterURLShortenerServer(s.s, &s.handler)

	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic("failed to listen: " + err.Error())
	}

	if err = s.s.Serve(l); err != nil {
		panic("failed to serve: " + err.Error())
	}
}
