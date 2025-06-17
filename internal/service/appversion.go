package service

import (
	"context"
	"github.com/liluoliluoli/gnboot/internal/common/utils/cache_util"
	"github.com/liluoliluoli/gnboot/internal/conf"
	"github.com/liluoliluoli/gnboot/internal/repo"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/redis/go-redis/v9"
)

type AppVersionService struct {
	c              *conf.Bootstrap
	appVersionRepo *repo.AppVersionRepo
	cache          sdomain.Cache[*sdomain.AppVersion]
	client         redis.UniversalClient
}

func NewAppVersionService(c *conf.Bootstrap,
	appVersionRepo *repo.AppVersionRepo,
	client redis.UniversalClient,
) *AppVersionService {
	return &AppVersionService{
		c:              c,
		appVersionRepo: appVersionRepo,
		client:         client,
		cache:          repo.NewCache[*sdomain.AppVersion](c, appVersionRepo.Data.Cache()),
	}
}

func (s *AppVersionService) GetLastVersion(ctx context.Context) (*sdomain.AppVersion, error) {
	return s.cache.Get(ctx, cache_util.GetCacheActionName(), func(action string, ctx context.Context) (*sdomain.AppVersion, error) {
		return s.get(ctx)
	})
}

func (s *AppVersionService) get(ctx context.Context) (*sdomain.AppVersion, error) {
	lastVersion, err := s.appVersionRepo.GetLastVersion(ctx)
	if err != nil {
		return nil, err
	}
	return lastVersion, nil
}
