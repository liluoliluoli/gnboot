package middleware

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	kratosHttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/liluoliluoli/gnboot/internal/common/gerror"
)

func Idempotent() middleware.Middleware {
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
			// check idempotent token
			token := tr.RequestHeader().Get("x-idempotent")
			//if token == "" {
			//	err = gerror.ErrIdempotentMissingToken(ctx)
			//	return
			//}
			idempotentKey := method + path + token
			//todo 查询redis是否有key
			fmt.Sprintf(idempotentKey)
			return handler(ctx, req)
		}
	}
}
