package context_util

import (
	"context"
	"github.com/liluoliluoli/gnboot/internal/common/constant"
	"github.com/liluoliluoli/gnboot/internal/common/gerror"
)

func GetGenericContext[R any](ctx context.Context, key string) (R, error) {
	var tempR R
	clu := ctx.Value(key)
	if clu == nil {
		return tempR, nil
	}
	r, ok := clu.(R)
	if !ok {
		return r, gerror.NewBizError(gerror.DataConvertError)
	}
	return r, nil
}

func SetContext[R any](ctx context.Context, key string, value R) context.Context {
	return context.WithValue(ctx, key, value)
}

func GetOperatorUid(ctx context.Context) string {
	uid, err := GetGenericContext[string](ctx, constant.GN_OPERATOR_CONTEXT)
	if err != nil {
		return "-1"
	}
	if uid == "" {
		return "0" //xxljob等未携带操作人信息的
	}
	return uid
}

func SetOperatorUidContext(ctx context.Context, uid string) context.Context {
	return SetContext(ctx, constant.GN_OPERATOR_CONTEXT, uid)
}
