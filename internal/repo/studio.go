package repo

import (
	"context"
	"github.com/go-cinch/common/page"
	"github.com/liluoliluoli/gnboot/internal/repo/gen"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type StudioRepo struct {
	Data *Data
}

func NewStudioRepo(data *Data) *StudioRepo {
	return &StudioRepo{
		Data: data,
	}
}

func (r *StudioRepo) do(ctx context.Context, tx *gen.Query) gen.IStudioDo {
	if tx == nil {
		return gen.Use(r.Data.DB(ctx)).Studio.WithContext(ctx)
	} else {
		return tx.Studio.WithContext(ctx)
	}
}

func (r *StudioRepo) Get(ctx context.Context, id int64) (*sdomain.Studio, error) {
	find, err := r.do(ctx, nil).Where(gen.Studio.ID.Eq(id)).First()
	if err != nil {
		return nil, handleQueryError(err)
	}
	return (&sdomain.Studio{}).ConvertFromRepo(find), nil
}

func (r *StudioRepo) Page(ctx context.Context, condition *sdomain.SearchMovie) (*sdomain.PageResult[*sdomain.Studio], error) {
	do := r.do(ctx, nil)
	if condition.Search != "" {
		do = do.Where(gen.Movie.OriginalTitle.Like("%" + condition.Search + "%"))
	}
	list, total, err := do.Order(gen.Movie.UpdateTime.Desc()).FindByPage(int((condition.Page.Num-1)*condition.Page.Size), int(condition.Page.Size))
	if err != nil {
		return nil, handleQueryError(err)
	}
	return &sdomain.PageResult[*sdomain.Studio]{
		Page: &page.Page{
			Num:   condition.Page.Num,
			Size:  condition.Page.Size,
			Total: total,
		},
		List: lo.Map(list, func(item *model.Studio, index int) *sdomain.Studio {
			return (&sdomain.Studio{}).ConvertFromRepo(item)
		}),
	}, nil
}

func (r *StudioRepo) FindByIds(ctx context.Context, ids []int64) ([]*sdomain.Studio, error) {
	finds, err := r.do(ctx, nil).Where(gen.Studio.ID.In(ids...)).Find()
	if err != nil {
		return nil, handleQueryError(err)
	}
	return lo.Map(finds, func(item *model.Studio, index int) *sdomain.Studio {
		return (&sdomain.Studio{}).ConvertFromRepo(item)
	}), nil
}

func (r *StudioRepo) Update(ctx context.Context, tx *gen.Query, movie *sdomain.UpdateMovie) error {
	updates, err := r.do(ctx, tx).Updates(movie.ConvertToRepo())
	if err != nil {
		return err
	}
	if updates.RowsAffected != 1 {
		return gorm.ErrDuplicatedKey
	}
	return nil
}

func (r *StudioRepo) Delete(ctx context.Context, tx *gen.Query, ids ...int64) error {
	_, err := r.do(ctx, tx).Where(gen.Movie.ID.In(ids...)).Delete()
	if err != nil {
		return err
	}
	return nil
}
