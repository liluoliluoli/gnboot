package middleware

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	kratosHttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/liluoliluoli/gnboot/internal/common/constant"
	"github.com/liluoliluoli/gnboot/internal/common/gerror"
	"github.com/liluoliluoli/gnboot/internal/common/utils/context_util"
	"github.com/redis/go-redis/v9"
)

func Auth(client redis.UniversalClient) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (rp interface{}, err error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				err = gerror.ErrIdempotentMissingToken(ctx)
				return
			}
			var method, path string
			switch tr.Kind() {
			case transport.KindHTTP:
				if ht, ok3 := tr.(kratosHttp.Transporter); ok3 {
					method = ht.Request().Method
					path = ht.Request().URL.Path
				}
			}
			if method == "POST" && (path == "/user/create" || path == "/user/login") {
				return handler(ctx, req)
			}

			token := tr.RequestHeader().Get("authorization")
			if token == "" {
				err = gerror.ErrAuthMissingToken(ctx)
				return
			}
			result, err := client.Get(ctx, token).Result()
			if err != nil {
				return nil, gerror.ErrAuthInvalidToken(ctx)
			}
			if result == "" {
				return nil, gerror.ErrAuthInvalidToken(ctx)
			}
			ctx = context_util.SetContext[string](ctx, constant.CTX_UserName, result)
			ctx = context_util.SetContext[string](ctx, constant.CTX_SessionToken, token)
			return handler(ctx, req)
		}
	}
}
