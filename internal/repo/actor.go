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
		return nil, handleQueryError(ctx, err)
	}
	return (&sdomain.Actor{}).ConvertFromRepo(find), nil
}

func (r *ActorRepo) FindByIds(ctx context.Context, ids []int64) ([]*sdomain.Actor, error) {
	finds, err := r.do(ctx, nil).Where(gen.Actor.ID.In(ids...)).Find()
	if err != nil {
		return nil, handleQueryError(ctx, err)
	}
	return lo.Map(finds, func(item *model.Actor, index int) *sdomain.Actor {
		return (&sdomain.Actor{}).ConvertFromRepo(item)
	}), nil
}

func (r *ActorRepo) FindAll(ctx context.Context) ([]*sdomain.Actor, error) {
	finds, err := r.do(ctx, nil).Find()
	if err != nil {
		return nil, handleQueryError(ctx, err)
	}
	return lo.Map(finds, func(item *model.Actor, index int) *sdomain.Actor {
		return (&sdomain.Actor{}).ConvertFromRepo(item)
	}), nil
}

func (r *ActorRepo) Update(ctx context.Context, tx *gen.Query, movie *sdomain.UpdateVideo) error {
	updates, err := r.do(ctx, tx).Updates(movie.ConvertToRepo())
	if err != nil {
		return err
	}
	if updates.RowsAffected != 1 {
		return gorm.ErrDuplicatedKey
	}
	return nil
}
