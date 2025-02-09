package middleware

import (
	"context"
	"github.com/alserok/url_shortener/pkg/logger"
	"google.golang.org/grpc"
)

func WithLogger(log logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		return handler(logger.WrapLogger(ctx, log), req)
	}
}
