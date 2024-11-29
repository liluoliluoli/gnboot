package service

import (
	"context"
	"github.com/liluoliluoli/gnboot/internal/conf"
	"github.com/liluoliluoli/gnboot/internal/repo"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/liluoliluoli/gnboot/internal/utils/cache_util"
)

type StudioService struct {
	c          *conf.Bootstrap
	studioRepo *repo.StudioRepo
	cache      sdomain.Cache[*sdomain.Studio]
}

func NewStudioService(c *conf.Bootstrap,
	studioRepo *repo.StudioRepo,
) *StudioService {
	return &StudioService{
		c:          c,
		studioRepo: studioRepo,
		cache:      repo.NewCache[*sdomain.Studio](c, studioRepo.Data.Cache()),
	}
}

func (s *StudioService) Get(ctx context.Context, id int64) (*sdomain.Studio, error) {
	return s.cache.Get(ctx, cache_util.GetCacheActionName(id), func(action string, ctx context.Context) (*sdomain.Studio, error) {
		return s.get(ctx, id)
	})
}

func (s *StudioService) get(ctx context.Context, id int64) (*sdomain.Studio, error) {
	item, err := s.studioRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (s *StudioService) FindAll(ctx context.Context) ([]*sdomain.Studio, error) {
	rp, err := s.cache.List(ctx, cache_util.GetCacheActionName(""), func(action string, ctx context.Context) ([]*sdomain.Studio, error) {
		return s.findAll(ctx)
	})
	if err != nil {
		return nil, err
	}
	return rp, nil
}

func (s *StudioService) findAll(ctx context.Context) ([]*sdomain.Studio, error) {
	finds, err := s.studioRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return finds, nil
}
