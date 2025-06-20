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
	"net"
	"net/http"
	"strings"
)

func Auth(client redis.UniversalClient) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (rp interface{}, err error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				err = gerror.ErrIdempotentMissingToken(ctx)
				return
			}
			var path string
			var clientIp string
			switch tr.Kind() {
			case transport.KindHTTP:
				if ht, ok3 := tr.(kratosHttp.Transporter); ok3 {
					path = ht.Request().URL.Path
					clientIp = GetClientIP(ht.Request())
				}
			}
			if path == "/api/user/create" || path == "/api/user/login" || path == "/api/version/getLastVersion" || strings.HasPrefix(path, "/api/test/") {
				return handler(ctx, req)
			}

			token := tr.RequestHeader().Get("Authorization")
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
			ctx = context_util.SetContext[string](ctx, constant.CTX_ClientIp, clientIp)
			return handler(ctx, req)
		}
	}
}

func GetClientIP(r *http.Request) string {
	// 1. 优先从 X-Forwarded-For 获取（多级代理时，第一个 IP 是客户端真实 IP）
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// 2. 尝试从 X-Real-IP 获取
	xRealIP := r.Header.Get("X-Real-IP")
	if xRealIP != "" {
		return xRealIP
	}

	// 3. 最后从 RemoteAddr 获取（直接连接时）
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr // 退化情况
	}
	return ip
}
