package gerror

import (
	"context"
	"strings"

	"github.com/go-cinch/common/constant"
	"github.com/go-cinch/common/middleware/i18n"
	"github.com/liluoliluoli/gnboot/api/reason"
)

var (
	ErrIdempotentMissingToken = func(ctx context.Context) error {
		return reason.ErrorIllegalParameter(i18n.FromContext(ctx).T(constant.IdempotentMissingToken))
	}
	ErrAuthMissingToken = func(ctx context.Context) error {
		return reason.ErrorIllegalParameter(i18n.FromContext(ctx).T(constant.JwtMissingToken))
	}
	ErrAuthInvalidToken = func(ctx context.Context) error {
		return reason.ErrorIllegalParameter(i18n.FromContext(ctx).T(constant.JwtTokenInvalid))
	}
	ErrMissingSignature = func(ctx context.Context) error {
		return i18n.NewError(ctx, "sign:signature:missing", reason.ErrorForbidden)
	}
	ErrSignature = func(ctx context.Context) error {
		return i18n.NewError(ctx, "sign:signature:error", reason.ErrorForbidden)
	}
	ErrTooManyRequests = func(ctx context.Context) error {
		return reason.ErrorTooManyRequests(i18n.FromContext(ctx).T(constant.TooManyRequests))
	}
	ErrExceedWatchCount = func(ctx context.Context, args ...string) error {
		return i18n.NewError(ctx, "exceed.watch.count", reason.ErrorForbidden, args...)
	}
	ErrAccountPackageExpire = func(ctx context.Context, args ...string) error {
		return i18n.NewError(ctx, "account.package.expire", reason.ErrorForbidden, args...)
	}
	ErrDataNotChange = func(ctx context.Context, args ...string) error {
		return i18n.NewError(ctx, constant.DataNotChange, reason.ErrorIllegalParameter, args...)
	}
	ErrDuplicateField = func(ctx context.Context, args ...string) error {
		return i18n.NewError(ctx, constant.DuplicateField, reason.ErrorIllegalParameter, args...)
	}
	ErrRecordNotFound = func(ctx context.Context, args ...string) error {
		return i18n.NewError(ctx, constant.RecordNotFound, reason.ErrorNotFound, args...)
	}
	ErrInternal = func(ctx context.Context, args ...string) error {
		return i18n.NewError(ctx, constant.InternalError, reason.ErrorInternal, args...)
	}
	ErrIllegalParameter = func(ctx context.Context, args ...string) error {
		return i18n.NewError(ctx, constant.IllegalParameter, reason.ErrorIllegalParameter, args...)
	}
	ErrDataConvert = func(ctx context.Context, args ...string) error {
		return i18n.NewError(ctx, "data.convert", reason.ErrorIllegalParameter, args...)
	}
	ErrJwtSign = func(ctx context.Context, args ...string) error {
		return i18n.NewError(ctx, constant.JwtTokenParseFail, reason.ErrorIllegalParameter, args...)
	}
)

func HandleRedisNotFoundError(err error) error {
	if err != nil && err.Error() == "redis: nil" {
		return nil
	}
	return err
}

func HandleRedisStringNotFound(str string) string {
	if strings.Contains(str, "redis: nil") {
		return ""
	}
	return str
}
