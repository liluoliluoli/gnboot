package middleware

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	"google.golang.org/grpc"
)

func GrpcUnaryDisableTimeoutPropagation() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if _, ok := ctx.Deadline(); ok {
			ctx = context.WithoutCancel(ctx)
		}

		resp, err = handler(ctx, req)
		return
	}
}

func HttpDisableTimeoutPropagation() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (rp interface{}, err error) {
			if _, ok := ctx.Deadline(); ok {
				ctx = context.WithoutCancel(ctx)
			}
			return handler(ctx, req)
		}
	}
}
