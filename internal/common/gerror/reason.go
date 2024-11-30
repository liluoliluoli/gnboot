package gerror

import (
	"context"

	"github.com/go-cinch/common/constant"
	"github.com/go-cinch/common/middleware/i18n"
	"github.com/liluoliluoli/gnboot/api/reason"
)

var (
	ErrIdempotentMissingToken = func(ctx context.Context) error {
		return reason.ErrorIllegalParameter(i18n.FromContext(ctx).T(constant.IdempotentMissingToken))
	}
	ErrTooManyRequests = func(ctx context.Context) error {
		return reason.ErrorTooManyRequests(i18n.FromContext(ctx).T(constant.TooManyRequests))
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
)

func HandleRedisNotFoundError(err error) error {
	if err != nil && err.Error() == "redis:nil" {
		return nil
	}
	return err
}
