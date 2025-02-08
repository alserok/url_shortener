package middleware

import (
	"context"
	"github.com/alserok/url_shortener/pkg/logger"
	"google.golang.org/grpc"
)

func WithRecovery() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		defer func() {
			if err := recover(); err != nil {
				logger.ExtractLogger(ctx).Error("panic recovery", logger.WithArg("error", err))
			}
		}()

		return handler(ctx, req)
	}
}
