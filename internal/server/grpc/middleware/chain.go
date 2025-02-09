package middleware

import "google.golang.org/grpc"

func WithChain(inters ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(inters...)
}
