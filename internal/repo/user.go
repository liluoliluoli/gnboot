package repo

import (
	"context"
	"github.com/liluoliluoli/gnboot/internal/repo/gen"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"gorm.io/gorm"
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

func (r *UserRepo) GetByUserName(ctx context.Context, userName string) (*sdomain.User, error) {
	find, err := r.do(ctx, nil).Where(gen.User.UserName.Eq(userName)).First()
	if err != nil {
		return nil, handleQueryError(ctx, err)
	}
	return (&sdomain.User{}).ConvertFromRepo(find), nil
}

func (r *UserRepo) Create(ctx context.Context, tx *gen.Query, userName string, password string) error {
	return r.do(ctx, tx).Create(&model.User{
		UserName: userName,
		Password: password,
	})
}

func (r *UserRepo) Update(ctx context.Context, tx *gen.Query, user *sdomain.User) error {
	updates, err := r.do(ctx, tx).Updates(user.ConvertToRepo())
	if err != nil {
		return err
	}
	if updates.RowsAffected != 1 {
		return gorm.ErrDuplicatedKey
	}
	return nil
}
