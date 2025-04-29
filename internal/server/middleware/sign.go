package middleware

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	kratosHttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/liluoliluoli/gnboot/internal/common/constant"
	"github.com/liluoliluoli/gnboot/internal/common/gerror"
	"github.com/liluoliluoli/gnboot/internal/common/utils/json_util"
	"github.com/liluoliluoli/gnboot/internal/common/utils/security_util"
)

func Sign(encryptKey string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (rp interface{}, err error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				err = gerror.ErrIdempotentMissingToken(ctx)
				return
			}
			var path string
			switch tr.Kind() {
			case transport.KindHTTP:
				if ht, ok3 := tr.(kratosHttp.Transporter); ok3 {
					path = ht.Request().URL.Path
				}
			}
			if path == "/api/version/getLastVersion" {
				return handler(ctx, req)
			}

			signature := tr.RequestHeader().Get("Signature")
			timestamp := tr.RequestHeader().Get("Timestamp")
			if signature == "" || timestamp == "" {
				err = gerror.ErrMissingSignature(ctx)
				return
			}
			if signature != security_util.SignByHMACSha256(path+timestamp, constant.SYS_PWD) {
				err = gerror.ErrSignature(ctx)
				return
			}
			return handler(ctx, req)
		}
	}
}

// decryptRequest 解密HTTP请求
func decryptRequest(body interface{}, key string) (interface{}, error) {
	bodyStr, err := json_util.MarshalString(body)
	if err != nil {
		return body, err
	}
	if len(bodyStr) > 0 {
		//todo 解密
		newBody, err := json_util.Unmarshal[interface{}](bodyStr)
		if err != nil {
			return body, err
		}
		return newBody, nil
	}
	return bodyStr, nil
}

// encryptResponse 加密HTTP响应
func encryptResponse(reply interface{}, key string) (interface{}, error) {
	replyStr, err := json_util.MarshalString(reply)
	if err != nil {
		return reply, err
	}
	if len(replyStr) > 0 {
		//todo 加密
		return replyStr, nil
	}
	return replyStr, nil
}
