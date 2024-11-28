package repo

import (
	"context"
	"github.com/go-cinch/common/page"
	"github.com/liluoliluoli/gnboot/internal/repo/gen"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/samber/lo"
)

type SeriesRepo struct {
	Data *Data
}

func NewSeriesRepo(data *Data) *SeriesRepo {
	return &SeriesRepo{
		Data: data,
	}
}

func (r *SeriesRepo) do(ctx context.Context, tx *gen.Query) gen.ISeriesDo {
	if tx == nil {
		return gen.Use(r.Data.DB(ctx)).Series.WithContext(ctx)
	} else {
		return tx.Series.WithContext(ctx)
	}
}

func (r *SeriesRepo) Get(ctx context.Context, id int64) (*sdomain.Series, error) {
	find, err := r.do(ctx, nil).Where(gen.Series.ID.Eq(id)).First()
	if err != nil {
		return nil, handleQueryError(err)
	}
	return (&sdomain.Series{}).ConvertFromRepo(find), nil
}

func (r *SeriesRepo) Page(ctx context.Context, condition *sdomain.SearchSeries) (*sdomain.PageResult[*sdomain.Series], error) {
	do := r.do(ctx, nil)
	if condition.Search != "" {
		do = do.Where(gen.Series.OriginalTitle.Like("%" + condition.Search + "%"))
	}
	if len(condition.FilterIds) != 0 {
		do = do.Where(gen.Series.ID.In(condition.FilterIds...))
	}
	list, total, err := do.Order(gen.Series.UpdateTime.Desc()).FindByPage(int((condition.Page.Num-1)*condition.Page.Size), int(condition.Page.Size))
	if err != nil {
		return nil, handleQueryError(err)
	}
	return &sdomain.PageResult[*sdomain.Series]{
		Page: &page.Page{
			Num:   condition.Page.Num,
			Size:  condition.Page.Size,
			Total: total,
		},
		List: lo.Map(list, func(item *model.Series, index int) *sdomain.Series {
			return (&sdomain.Series{}).ConvertFromRepo(item)
		}),
	}, nil
}
