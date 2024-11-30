package service

import (
	"context"
	"github.com/liluoliluoli/gnboot/internal/common/utils/cache_util"
	"github.com/liluoliluoli/gnboot/internal/conf"
	"github.com/liluoliluoli/gnboot/internal/repo"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
)

type UserService struct {
	c        *conf.Bootstrap
	userRepo *repo.UserRepo
	cache    sdomain.Cache[*sdomain.User]
}

func NewUserService(c *conf.Bootstrap,
	userRepo *repo.UserRepo,
) *UserService {
	return &UserService{
		c:        c,
		userRepo: userRepo,
		cache:    repo.NewCache[*sdomain.User](c, userRepo.Data.Cache()),
	}
}

func (s *UserService) Get(ctx context.Context, id int64) (*sdomain.User, error) {
	return s.cache.Get(ctx, cache_util.GetCacheActionName(id), func(action string, ctx context.Context) (*sdomain.User, error) {
		return s.get(ctx, id)
	})
}

func (s *UserService) get(ctx context.Context, id int64) (*sdomain.User, error) {
	item, err := s.userRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return item, nil
}
