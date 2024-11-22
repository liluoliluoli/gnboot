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

type CreateSeries struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
}

type Series struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
}

type FindSeries struct {
	Page page.Page `json:"page"`
	Name *string   `json:"name"`
}

type FindSeriesCache struct {
	Page page.Page `json:"page"`
	List []Series  `json:"list"`
}

type UpdateSeries struct {
	ID   uint64  `json:"id,string"`
	Name *string `json:"name,omitempty"`
}

type SeriesRepo interface {
	Create(ctx context.Context, item *CreateSeries) error
	Get(ctx context.Context, id uint64) (*Series, error)
	Find(ctx context.Context, condition *FindSeries) []Series
	Update(ctx context.Context, item *UpdateSeries) error
	Delete(ctx context.Context, ids ...uint64) error
}

type SeriesUseCase struct {
	c     *conf.Bootstrap
	repo  SeriesRepo
	tx    Transaction
	cache Cache
}

func NewSeriesUseCase(c *conf.Bootstrap, repo SeriesRepo, tx Transaction, cache Cache) *SeriesUseCase {
	return &SeriesUseCase{
		c:    c,
		repo: repo,
		tx:   tx,
		cache: cache.WithPrefix(strings.Join([]string{
			c.Name, "series",
		}, "_")),
	}
}

func (uc *SeriesUseCase) Create(ctx context.Context, item *CreateSeries) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Create(ctx, item)
		})
	})
}

func (uc *SeriesUseCase) Get(ctx context.Context, id uint64) (rp *Series, err error) {
	rp = &Series{}
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

func (uc *SeriesUseCase) get(ctx context.Context, action string, id uint64) (res string, err error) {
	// read data from db and write to cache
	rp := &Series{}
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

func (uc *SeriesUseCase) Find(ctx context.Context, condition *FindSeries) (rp []Series, err error) {
	// use md5 string as cache replay json str, key is short
	str, err := uc.cache.Get(ctx, cache_util.GetCacheActionName(condition), func(action string, ctx context.Context) (string, error) {
		return uc.find(ctx, action, condition)
	})
	if err != nil {
		return
	}
	var cache FindSeriesCache
	utils.Json2Struct(&cache, str)
	condition.Page = cache.Page
	rp = cache.List
	return
}

func (uc *SeriesUseCase) find(ctx context.Context, action string, condition *FindSeries) (res string, err error) {
	// read data from db and write to cache
	list := uc.repo.Find(ctx, condition)
	var cache FindSeriesCache
	cache.List = list
	cache.Page = condition.Page
	res = utils.Struct2Json(cache)
	uc.cache.Set(ctx, action, res, len(list) == 0)
	return
}

func (uc *SeriesUseCase) Update(ctx context.Context, item *UpdateSeries) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Update(ctx, item)
			return
		})
	})
}

func (uc *SeriesUseCase) Delete(ctx context.Context, ids ...uint64) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Delete(ctx, ids...)
			return
		})
	})
}
