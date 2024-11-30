package repo

import (
	"context"
	"github.com/liluoliluoli/gnboot/internal/repo/gen"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type ActorRepo struct {
	Data *Data
}

func NewActorRepo(data *Data) *ActorRepo {
	return &ActorRepo{
		Data: data,
	}
}

func (r *ActorRepo) do(ctx context.Context, tx *gen.Query) gen.IActorDo {
	if tx == nil {
		return gen.Use(r.Data.DB(ctx)).Actor.WithContext(ctx)
	} else {
		return tx.Actor.WithContext(ctx)
	}
}

func (r *ActorRepo) Get(ctx context.Context, id int64) (*sdomain.Actor, error) {
	find, err := r.do(ctx, nil).Where(gen.Actor.ID.Eq(id)).First()
	if err != nil {
		return nil, handleQueryError(err)
	}
	return (&sdomain.Actor{}).ConvertFromRepo(find), nil
}

func (r *ActorRepo) Page(ctx context.Context, condition *sdomain.SearchMovie) (*sdomain.PageResult[*sdomain.Actor], error) {
	do := r.do(ctx, nil)
	if condition.Search != "" {
		do = do.Where(gen.Movie.OriginalTitle.Like("%" + condition.Search + "%"))
	}
	list, total, err := do.Order(gen.Movie.UpdateTime.Desc()).FindByPage(int((condition.Page.CurrentPage-1)*condition.Page.PageSize), int(condition.Page.PageSize))
	if err != nil {
		return nil, handleQueryError(err)
	}
	return &sdomain.PageResult[*sdomain.Actor]{
		Page: &sdomain.Page{
			CurrentPage: condition.Page.CurrentPage,
			PageSize:    condition.Page.PageSize,
			Count:       total,
		},
		List: lo.Map(list, func(item *model.Actor, index int) *sdomain.Actor {
			return (&sdomain.Actor{}).ConvertFromRepo(item)
		}),
	}, nil
}

func (r *ActorRepo) FindByIds(ctx context.Context, ids []int64) ([]*sdomain.Actor, error) {
	finds, err := r.do(ctx, nil).Where(gen.Actor.ID.In(ids...)).Find()
	if err != nil {
		return nil, handleQueryError(err)
	}
	return lo.Map(finds, func(item *model.Actor, index int) *sdomain.Actor {
		return (&sdomain.Actor{}).ConvertFromRepo(item)
	}), nil
}

func (r *ActorRepo) FindAll(ctx context.Context) ([]*sdomain.Actor, error) {
	finds, err := r.do(ctx, nil).Find()
	if err != nil {
		return nil, handleQueryError(err)
	}
	return lo.Map(finds, func(item *model.Actor, index int) *sdomain.Actor {
		return (&sdomain.Actor{}).ConvertFromRepo(item)
	}), nil
}

func (r *ActorRepo) Update(ctx context.Context, tx *gen.Query, movie *sdomain.UpdateMovie) error {
	updates, err := r.do(ctx, tx).Updates(movie.ConvertToRepo())
	if err != nil {
		return err
	}
	if updates.RowsAffected != 1 {
		return gorm.ErrDuplicatedKey
	}
	return nil
}

func (r *ActorRepo) Delete(ctx context.Context, tx *gen.Query, ids ...int64) error {
	_, err := r.do(ctx, tx).Where(gen.Movie.ID.In(ids...)).Delete()
	if err != nil {
		return err
	}
	return nil
}
