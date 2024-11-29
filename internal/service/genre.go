package service

import (
	"context"
	"github.com/liluoliluoli/gnboot/internal/conf"
	"github.com/liluoliluoli/gnboot/internal/repo"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/liluoliluoli/gnboot/internal/utils/cache_util"
)

type GenreService struct {
	c         *conf.Bootstrap
	genreRepo *repo.GenreRepo
	cache     sdomain.Cache[*sdomain.Genre]
}

func NewGenreService(c *conf.Bootstrap,
	genreRepo *repo.GenreRepo,
) *GenreService {
	return &GenreService{
		c:         c,
		genreRepo: genreRepo,
		cache:     repo.NewCache[*sdomain.Genre](c, genreRepo.Data.Cache()),
	}
}

func (s *GenreService) Get(ctx context.Context, id int64) (*sdomain.Genre, error) {
	return s.cache.Get(ctx, cache_util.GetCacheActionName(id), func(action string, ctx context.Context) (*sdomain.Genre, error) {
		return s.get(ctx, id)
	})
}

func (s *GenreService) get(ctx context.Context, id int64) (*sdomain.Genre, error) {
	item, err := s.genreRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (s *GenreService) FindAll(ctx context.Context) ([]*sdomain.Genre, error) {
	rp, err := s.cache.List(ctx, cache_util.GetCacheActionName(""), func(action string, ctx context.Context) ([]*sdomain.Genre, error) {
		return s.findAll(ctx)
	})
	if err != nil {
		return nil, err
	}
	return rp, nil
}

func (s *GenreService) findAll(ctx context.Context) ([]*sdomain.Genre, error) {
	finds, err := s.genreRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return finds, nil
}
