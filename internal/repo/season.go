package repo

import (
	"context"
	"github.com/liluoliluoli/gnboot/internal/repo/gen"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/samber/lo"
)

type SeasonRepo struct {
	Data *Data
}

func NewSeasonRepo(data *Data) *SeasonRepo {
	return &SeasonRepo{
		Data: data,
	}
}

func (r *SeasonRepo) do(ctx context.Context, tx *gen.Query) gen.ISeasonDo {
	if tx == nil {
		return gen.Use(r.Data.DB(ctx)).Season.WithContext(ctx)
	} else {
		return tx.Season.WithContext(ctx)
	}
}

func (r *SeasonRepo) Get(ctx context.Context, id int64) (*sdomain.Season, error) {
	find, err := r.do(ctx, nil).Where(gen.Season.ID.Eq(id)).First()
	if err != nil {
		return nil, handleQueryError(err)
	}
	return (&sdomain.Season{}).ConvertFromRepo(find), nil
}

func (r *SeasonRepo) QueryBySeriesId(ctx context.Context, seriesId int64) ([]*sdomain.Season, error) {
	finds, err := r.do(ctx, nil).Order(gen.Season.Season).Find()
	if err != nil {
		return nil, handleQueryError(err)
	}
	return lo.Map(finds, func(item *model.Season, index int) *sdomain.Season {
		return (&sdomain.Season{}).ConvertFromRepo(item)
	}), nil
}
