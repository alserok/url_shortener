package middleware

import (
	"context"
	"github.com/alserok/url_shortener/internal/utils"
	"github.com/alserok/url_shortener/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func WithErrorHandler() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		res, err := handler(ctx, req)
		if err != nil {
			msg, code := utils.FromError(ctx, err)

			switch code {
			case utils.BadRequestErr:
				return nil, status.Error(codes.InvalidArgument, msg)
			case utils.NotFoundErr:
				return nil, status.Error(codes.NotFound, msg)
			case utils.InternalErr:
				logger.ExtractLogger(ctx).Error(err.Error())
				return nil, status.Error(codes.Internal, msg)
			default:
				logger.ExtractLogger(ctx).Error(err.Error())
				return nil, status.Error(codes.Internal, msg)
			}
		}

		return res, nil
	}
}
