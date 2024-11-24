package service

import (
	"context"
	"gnboot/internal/conf"
	"gnboot/internal/repo"
	"gnboot/internal/repo/gen"
	"gnboot/internal/service/sdomain"
	"gnboot/internal/utils/cache_util"
)

type MovieUseCase struct {
	c     *conf.Bootstrap
	repo  *repo.MovieRepo
	cache sdomain.Cache[*sdomain.Movie]
}

func NewMovieUseCase(c *conf.Bootstrap, rp *repo.MovieRepo) *MovieUseCase {
	return &MovieUseCase{
		c:     c,
		repo:  rp,
		cache: repo.NewCache[*sdomain.Movie](c, rp.Data.Cache()),
	}
}

func (uc *MovieUseCase) Create(ctx context.Context, item *sdomain.CreateMovie) error {
	err := gen.Use(uc.repo.Data.DB(ctx)).Transaction(func(tx *gen.Query) error {
		err := uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Create(ctx, tx, item)
		})
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (uc *MovieUseCase) Get(ctx context.Context, id int64) (*sdomain.Movie, error) {
	return uc.cache.Get(ctx, cache_util.GetCacheActionName(id), func(action string, ctx context.Context) (*sdomain.Movie, error) {
		return uc.get(ctx, id)
	})
}

func (uc *MovieUseCase) get(ctx context.Context, id int64) (*sdomain.Movie, error) {
	item, err := uc.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (uc *MovieUseCase) Page(ctx context.Context, condition *sdomain.FindMovie) (*sdomain.PageResult[*sdomain.Movie], error) {
	rp, err := uc.cache.GetPage(ctx, cache_util.GetCacheActionName(condition), func(action string, ctx context.Context) (*sdomain.PageResult[*sdomain.Movie], error) {
		return uc.page(ctx, condition)
	})
	if err != nil {
		return nil, err
	}
	return rp, nil
}

func (uc *MovieUseCase) page(ctx context.Context, condition *sdomain.FindMovie) (*sdomain.PageResult[*sdomain.Movie], error) {
	pageResult, err := uc.repo.Page(ctx, condition)
	if err != nil {
		return nil, err
	}
	return pageResult, nil
}

func (uc *MovieUseCase) Update(ctx context.Context, item *sdomain.UpdateMovie) error {
	err := gen.Use(uc.repo.Data.DB(ctx)).Transaction(func(tx *gen.Query) error {
		err := uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Update(ctx, tx, item)
		})
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (uc *MovieUseCase) Delete(ctx context.Context, ids ...int64) error {
	err := gen.Use(uc.repo.Data.DB(ctx)).Transaction(func(tx *gen.Query) error {
		err := uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Delete(ctx, tx, ids...)
		})
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
