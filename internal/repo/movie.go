package repo

import (
	"context"
	"github.com/go-cinch/common/page"
	"github.com/samber/lo"
	"gnboot/internal/repo/gen"
	"gnboot/internal/repo/model"
	"gnboot/internal/service/sdomain"
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

func (r *MovieRepo) Page(ctx context.Context, condition *sdomain.FindMovie) (*sdomain.PageResult[*sdomain.Movie], error) {
	do := r.do(ctx, nil)
	if condition.Search != nil {
		do = do.Where(gen.Q.Movie.OriginalTitle.Like("%" + *condition.Search + "%"))
	}
	list, total, err := do.Order(gen.Q.Movie.UpdateTime.Desc()).FindByPage(0, 10)
	if err != nil {
		return nil, handleQueryError(err)
	}
	return &sdomain.PageResult[*sdomain.Movie]{
		Page: &page.Page{
			Num:   condition.Page.Num,
			Size:  condition.Page.Size,
			Total: total,
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
