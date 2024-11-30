package repo

import (
	"context"
	"github.com/liluoliluoli/gnboot/internal/repo/gen"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type MovieRepo struct {
	Data *Data
}

func NewMovieRepo(data *Data) *MovieRepo {
	return &MovieRepo{
		Data: data,
	}
}

func (r *MovieRepo) do(ctx context.Context, tx *gen.Query) gen.IMovieDo {
	if tx == nil {
		return gen.Use(r.Data.DB(ctx)).Movie.WithContext(ctx)
	} else {
		return tx.Movie.WithContext(ctx)
	}
}

func (r *MovieRepo) Get(ctx context.Context, id int64) (*sdomain.Movie, error) {
	find, err := r.do(ctx, nil).Where(gen.Movie.ID.Eq(id)).First()
	if err != nil {
		return nil, handleQueryError(err)
	}
	return (&sdomain.Movie{}).ConvertFromRepo(find), nil
}

func (r *MovieRepo) Page(ctx context.Context, condition *sdomain.SearchMovie) (*sdomain.PageResult[*sdomain.Movie], error) {
	do := r.do(ctx, nil)
	if condition.Search != "" {
		do = do.Where(gen.Movie.OriginalTitle.Like("%" + condition.Search + "%"))
	}
	if len(condition.FilterIds) != 0 {
		do = do.Where(gen.Movie.ID.In(condition.FilterIds...))
	}
	list, total, err := do.Order(gen.Movie.UpdateTime.Desc()).FindByPage(int((condition.Page.CurrentPage-1)*condition.Page.PageSize), int(condition.Page.PageSize))
	if err != nil {
		return nil, handleQueryError(err)
	}
	return &sdomain.PageResult[*sdomain.Movie]{
		Page: &sdomain.Page{
			CurrentPage: condition.Page.CurrentPage,
			PageSize:    condition.Page.PageSize,
			Count:       total,
		},
		List: lo.Map(list, func(item *model.Movie, index int) *sdomain.Movie {
			return (&sdomain.Movie{}).ConvertFromRepo(item)
		}),
	}, nil
}

func (r *MovieRepo) Create(ctx context.Context, tx *gen.Query, movie *sdomain.CreateMovie) error {
	err := r.do(ctx, tx).Save(movie.ConvertToRepo())
	if err != nil {
		return err
	}
	return nil
}

func (r *MovieRepo) Update(ctx context.Context, tx *gen.Query, movie *sdomain.UpdateMovie) error {
	updates, err := r.do(ctx, tx).Updates(movie.ConvertToRepo())
	if err != nil {
		return err
	}
	if updates.RowsAffected != 1 {
		return gorm.ErrDuplicatedKey
	}
	return nil
}

func (r *MovieRepo) Delete(ctx context.Context, tx *gen.Query, ids ...int64) error {
	_, err := r.do(ctx, tx).Where(gen.Movie.ID.In(ids...)).Delete()
	if err != nil {
		return err
	}
	return nil
}
