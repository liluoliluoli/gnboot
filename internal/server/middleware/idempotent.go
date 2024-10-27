package middleware

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
)

const (
	WhitelistIdempotentCategory uint32 = 2
)

func Idempotent() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (rp interface{}, err error) {
			return handler(ctx, req)
		}
	}
}
