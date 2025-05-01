package middleware

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	kratosHttp "github.com/go-kratos/kratos/v2/transport/http"
)

func Cors() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				tr.ReplyHeader().Set("Access-Control-Allow-Origin", "*")
				tr.ReplyHeader().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
				tr.ReplyHeader().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
				var method string
				switch tr.Kind() {
				case transport.KindHTTP:
					if ht, ok3 := tr.(kratosHttp.Transporter); ok3 {
						method = ht.Request().Method
						if method == "OPTIONS" {
							return nil, nil
						}
					}
				}
			}
			return handler(ctx, req)
		}
	}
}
