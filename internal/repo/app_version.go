package repo

import (
	"context"
	"github.com/liluoliluoli/gnboot/internal/repo/gen"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
)

type AppVersionRepo struct {
	Data *Data
}

func NewAppVersionRepo(data *Data) *AppVersionRepo {
	return &AppVersionRepo{
		Data: data,
	}
}

func (r *AppVersionRepo) do(ctx context.Context, tx *gen.Query) gen.IAppVersionDo {
	if tx == nil {
		return gen.Use(r.Data.DB(ctx)).AppVersion.WithContext(ctx)
	} else {
		return tx.AppVersion.WithContext(ctx)
	}
}

func (r *AppVersionRepo) GetLastVersion(ctx context.Context) (*sdomain.AppVersion, error) {
	find, err := r.do(ctx, nil).Order(gen.AppVersion.VersionCode.Desc()).First()
	if err != nil {
		return nil, handleQueryError(ctx, err)
	}
	return (&sdomain.AppVersion{}).ConvertFromRepo(find), nil
}
