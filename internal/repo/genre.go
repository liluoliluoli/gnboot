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

type GenreRepo struct {
	Data *Data
}

func NewGenreRepo(data *Data) *GenreRepo {
	return &GenreRepo{
		Data: data,
	}
}

func (r *GenreRepo) do(ctx context.Context, tx *gen.Query) gen.IGenreDo {
	if tx == nil {
		return gen.Use(r.Data.DB(ctx)).Genre.WithContext(ctx)
	} else {
		return tx.Genre.WithContext(ctx)
	}
}

func (r *GenreRepo) Get(ctx context.Context, id int64) (*sdomain.Genre, error) {
	find, err := r.do(ctx, nil).Where(gen.Movie.ID.Eq(id)).First()
	if err != nil {
		return nil, handleQueryError(err)
	}
	return (&sdomain.Genre{}).ConvertFromRepo(find), nil
}

func (r *GenreRepo) Page(ctx context.Context, condition *sdomain.SearchMovie) (*sdomain.PageResult[*sdomain.Genre], error) {
	do := r.do(ctx, nil)
	if condition.Search != "" {
		do = do.Where(gen.Movie.OriginalTitle.Like("%" + condition.Search + "%"))
	}
	list, total, err := do.Order(gen.Movie.UpdateTime.Desc()).FindByPage(int((condition.Page.Num-1)*condition.Page.Size), int(condition.Page.Size))
	if err != nil {
		return nil, handleQueryError(err)
	}
	return &sdomain.PageResult[*sdomain.Genre]{
		Page: &page.Page{
			Num:   condition.Page.Num,
			Size:  condition.Page.Size,
			Total: total,
		},
		List: lo.Map(list, func(item *model.Genre, index int) *sdomain.Genre {
			return (&sdomain.Genre{}).ConvertFromRepo(item)
		}),
	}, nil
}

func (r *GenreRepo) FindByIds(ctx context.Context, ids []int64) ([]*sdomain.Genre, error) {
	finds, err := r.do(ctx, nil).Where(gen.Genre.ID.In(ids...)).Find()
	if err != nil {
		return nil, handleQueryError(err)
	}
	return lo.Map(finds, func(item *model.Genre, index int) *sdomain.Genre {
		return (&sdomain.Genre{}).ConvertFromRepo(item)
	}), nil
}

func (r *GenreRepo) Update(ctx context.Context, tx *gen.Query, movie *sdomain.UpdateMovie) error {
	updates, err := r.do(ctx, tx).Updates(movie.ConvertToRepo())
	if err != nil {
		return err
	}
	if updates.RowsAffected != 1 {
		return gorm.ErrDuplicatedKey
	}
	return nil
}

func (r *GenreRepo) Delete(ctx context.Context, tx *gen.Query, ids ...int64) error {
	_, err := r.do(ctx, tx).Where(gen.Movie.ID.In(ids...)).Delete()
	if err != nil {
		return err
	}
	return nil
}
