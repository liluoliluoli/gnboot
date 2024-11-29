package service

import (
	"context"
	"github.com/liluoliluoli/gnboot/internal/conf"
	"github.com/liluoliluoli/gnboot/internal/repo"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/liluoliluoli/gnboot/internal/utils/cache_util"
)

type ActorService struct {
	c         *conf.Bootstrap
	actorRepo *repo.ActorRepo
	cache     sdomain.Cache[*sdomain.Actor]
}

func NewActorService(c *conf.Bootstrap,
	actorRepo *repo.ActorRepo,
) *ActorService {
	return &ActorService{
		c:         c,
		actorRepo: actorRepo,
		cache:     repo.NewCache[*sdomain.Actor](c, actorRepo.Data.Cache()),
	}
}

func (s *ActorService) Get(ctx context.Context, id int64) (*sdomain.Actor, error) {
	return s.cache.Get(ctx, cache_util.GetCacheActionName(id), func(action string, ctx context.Context) (*sdomain.Actor, error) {
		return s.get(ctx, id)
	})
}

func (s *ActorService) get(ctx context.Context, id int64) (*sdomain.Actor, error) {
	item, err := s.actorRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (s *ActorService) FindAll(ctx context.Context) ([]*sdomain.Actor, error) {
	rp, err := s.cache.List(ctx, cache_util.GetCacheActionName(""), func(action string, ctx context.Context) ([]*sdomain.Actor, error) {
		return s.findAll(ctx)
	})
	if err != nil {
		return nil, err
	}
	return rp, nil
}

func (s *ActorService) findAll(ctx context.Context) ([]*sdomain.Actor, error) {
	finds, err := s.actorRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return finds, nil
}
