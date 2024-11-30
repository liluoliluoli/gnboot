package repo

import (
	"context"
	"github.com/liluoliluoli/gnboot/internal/repo/gen"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
)

type UserRepo struct {
	Data *Data
}

func NewUserRepo(data *Data) *UserRepo {
	return &UserRepo{
		Data: data,
	}
}

func (r *UserRepo) do(ctx context.Context, tx *gen.Query) gen.IUserDo {
	if tx == nil {
		return gen.Use(r.Data.DB(ctx)).User.WithContext(ctx)
	} else {
		return tx.User.WithContext(ctx)
	}
}

func (r *UserRepo) Get(ctx context.Context, id int64) (*sdomain.User, error) {
	find, err := r.do(ctx, nil).Where(gen.User.ID.Eq(id)).First()
	if err != nil {
		return nil, handleQueryError(ctx, err)
	}
	return (&sdomain.User{}).ConvertFromRepo(find), nil
}
