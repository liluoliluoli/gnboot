package biz

import (
	"context"
	"gnboot/internal/utils/cache_util"
	"strings"

	"github.com/go-cinch/common/constant"
	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/page"
	"github.com/go-cinch/common/utils"
	"github.com/pkg/errors"
	"gnboot/internal/conf"
)

type CreateGenre struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
}

type Genre struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
}

type FindGenre struct {
	Page page.Page `json:"page"`
	Name *string   `json:"name"`
}

type FindGenreCache struct {
	Page page.Page `json:"page"`
	List []Genre   `json:"list"`
}

type UpdateGenre struct {
	ID   uint64  `json:"id,string"`
	Name *string `json:"name,omitempty"`
}

type GenreRepo interface {
	Create(ctx context.Context, item *CreateGenre) error
	Get(ctx context.Context, id uint64) (*Genre, error)
	Find(ctx context.Context, condition *FindGenre) []Genre
	Update(ctx context.Context, item *UpdateGenre) error
	Delete(ctx context.Context, ids ...uint64) error
}

type GenreUseCase struct {
	c     *conf.Bootstrap
	repo  GenreRepo
	tx    Transaction
	cache Cache
}

func NewGenreUseCase(c *conf.Bootstrap, repo GenreRepo, tx Transaction, cache Cache) *GenreUseCase {
	return &GenreUseCase{
		c:    c,
		repo: repo,
		tx:   tx,
		cache: cache.WithPrefix(strings.Join([]string{
			c.Name, "genre",
		}, "_")),
	}
}

func (uc *GenreUseCase) Create(ctx context.Context, item *CreateGenre) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Create(ctx, item)
		})
	})
}

func (uc *GenreUseCase) Get(ctx context.Context, id uint64) (rp *Genre, err error) {
	rp = &Genre{}
	str, err := uc.cache.Get(ctx, cache_util.GetCacheActionName(id), func(action string, ctx context.Context) (string, error) {
		return uc.get(ctx, action, id)
	})
	if err != nil {
		return
	}
	utils.Json2Struct(&rp, str)
	if rp.ID == constant.UI0 {
		err = ErrRecordNotFound(ctx)
		return
	}
	return
}

func (uc *GenreUseCase) get(ctx context.Context, action string, id uint64) (res string, err error) {
	// read data from db and write to cache
	rp := &Genre{}
	item, err := uc.repo.Get(ctx, id)
	notFound := errors.Is(err, ErrRecordNotFound(ctx))
	if err != nil && !notFound {
		return
	}
	copierx.Copy(&rp, item)
	res = utils.Struct2Json(rp)
	uc.cache.Set(ctx, action, res, notFound)
	return
}

func (uc *GenreUseCase) Find(ctx context.Context, condition *FindGenre) (rp []Genre, err error) {
	// use md5 string as cache replay json str, key is short
	str, err := uc.cache.Get(ctx, cache_util.GetCacheActionName(condition), func(action string, ctx context.Context) (string, error) {
		return uc.find(ctx, action, condition)
	})
	if err != nil {
		return
	}
	var cache FindGenreCache
	utils.Json2Struct(&cache, str)
	condition.Page = cache.Page
	rp = cache.List
	return
}

func (uc *GenreUseCase) find(ctx context.Context, action string, condition *FindGenre) (res string, err error) {
	// read data from db and write to cache
	list := uc.repo.Find(ctx, condition)
	var cache FindGenreCache
	cache.List = list
	cache.Page = condition.Page
	res = utils.Struct2Json(cache)
	uc.cache.Set(ctx, action, res, len(list) == 0)
	return
}

func (uc *GenreUseCase) Update(ctx context.Context, item *UpdateGenre) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Update(ctx, item)
			return
		})
	})
}

func (uc *GenreUseCase) Delete(ctx context.Context, ids ...uint64) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Delete(ctx, ids...)
			return
		})
	})
}
