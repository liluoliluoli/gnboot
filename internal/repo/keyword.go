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

type KeywordRepo struct {
	Data *Data
}

func NewKeywordRepo(data *Data) *KeywordRepo {
	return &KeywordRepo{
		Data: data,
	}
}

func (r *KeywordRepo) do(ctx context.Context, tx *gen.Query) gen.IKeywordDo {
	if tx == nil {
		return gen.Use(r.Data.DB(ctx)).Keyword.WithContext(ctx)
	} else {
		return tx.Keyword.WithContext(ctx)
	}
}

func (r *KeywordRepo) Get(ctx context.Context, id int64) (*sdomain.Keyword, error) {
	find, err := r.do(ctx, nil).Where(gen.Keyword.ID.Eq(id)).First()
	if err != nil {
		return nil, handleQueryError(err)
	}
	return (&sdomain.Keyword{}).ConvertFromRepo(find), nil
}

func (r *KeywordRepo) Page(ctx context.Context, condition *sdomain.FindMovie) (*sdomain.PageResult[*sdomain.Keyword], error) {
	do := r.do(ctx, nil)
	if condition.Search != nil {
		do = do.Where(gen.Movie.OriginalTitle.Like("%" + *condition.Search + "%"))
	}
	list, total, err := do.Order(gen.Movie.UpdateTime.Desc()).FindByPage(int((condition.Page.Num-1)*condition.Page.Size), int(condition.Page.Size))
	if err != nil {
		return nil, handleQueryError(err)
	}
	return &sdomain.PageResult[*sdomain.Keyword]{
		Page: &page.Page{
			Num:   condition.Page.Num,
			Size:  condition.Page.Size,
			Total: total,
		},
		List: lo.Map(list, func(item *model.Keyword, index int) *sdomain.Keyword {
			return (&sdomain.Keyword{}).ConvertFromRepo(item)
		}),
	}, nil
}

func (r *KeywordRepo) FindByIds(ctx context.Context, ids []int64) ([]*sdomain.Keyword, error) {
	finds, err := r.do(ctx, nil).Where(gen.Genre.ID.In(ids...)).Find()
	if err != nil {
		return nil, handleQueryError(err)
	}
	return lo.Map(finds, func(item *model.Keyword, index int) *sdomain.Keyword {
		return (&sdomain.Keyword{}).ConvertFromRepo(item)
	}), nil
}

func (r *KeywordRepo) Update(ctx context.Context, tx *gen.Query, movie *sdomain.UpdateMovie) error {
	updates, err := r.do(ctx, tx).Updates(movie.ConvertToRepo())
	if err != nil {
		return err
	}
	if updates.RowsAffected != 1 {
		return gorm.ErrDuplicatedKey
	}
	return nil
}

func (r *KeywordRepo) Delete(ctx context.Context, tx *gen.Query, ids ...int64) error {
	_, err := r.do(ctx, tx).Where(gen.Movie.ID.In(ids...)).Delete()
	if err != nil {
		return err
	}
	return nil
}
