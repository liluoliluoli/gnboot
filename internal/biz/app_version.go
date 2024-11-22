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

type CreateAppVersion struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
}

type AppVersion struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
}

type FindAppVersion struct {
	Page page.Page `json:"page"`
	Name *string   `json:"name"`
}

type FindAppVersionCache struct {
	Page page.Page    `json:"page"`
	List []AppVersion `json:"list"`
}

type UpdateAppVersion struct {
	ID   uint64  `json:"id,string"`
	Name *string `json:"name,omitempty"`
}

type AppVersionRepo interface {
	Create(ctx context.Context, item *CreateAppVersion) error
	Get(ctx context.Context, id uint64) (*AppVersion, error)
	Find(ctx context.Context, condition *FindAppVersion) []AppVersion
	Update(ctx context.Context, item *UpdateAppVersion) error
	Delete(ctx context.Context, ids ...uint64) error
}

type AppVersionUseCase struct {
	c     *conf.Bootstrap
	repo  AppVersionRepo
	tx    Transaction
	cache Cache
}

func NewAppVersionUseCase(c *conf.Bootstrap, repo AppVersionRepo, tx Transaction, cache Cache) *AppVersionUseCase {
	return &AppVersionUseCase{
		c:    c,
		repo: repo,
		tx:   tx,
		cache: cache.WithPrefix(strings.Join([]string{
			c.Name, "app_version",
		}, "_")),
	}
}

func (uc *AppVersionUseCase) Create(ctx context.Context, item *CreateAppVersion) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Create(ctx, item)
		})
	})
}

func (uc *AppVersionUseCase) Get(ctx context.Context, id uint64) (rp *AppVersion, err error) {
	rp = &AppVersion{}
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

func (uc *AppVersionUseCase) get(ctx context.Context, action string, id uint64) (res string, err error) {
	// read data from db and write to cache
	rp := &AppVersion{}
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

func (uc *AppVersionUseCase) Find(ctx context.Context, condition *FindAppVersion) (rp []AppVersion, err error) {
	// use md5 string as cache replay json str, key is short
	str, err := uc.cache.Get(ctx, cache_util.GetCacheActionName(condition), func(action string, ctx context.Context) (string, error) {
		return uc.find(ctx, action, condition)
	})
	if err != nil {
		return
	}
	var cache FindAppVersionCache
	utils.Json2Struct(&cache, str)
	condition.Page = cache.Page
	rp = cache.List
	return
}

func (uc *AppVersionUseCase) find(ctx context.Context, action string, condition *FindAppVersion) (res string, err error) {
	// read data from db and write to cache
	list := uc.repo.Find(ctx, condition)
	var cache FindAppVersionCache
	cache.List = list
	cache.Page = condition.Page
	res = utils.Struct2Json(cache)
	uc.cache.Set(ctx, action, res, len(list) == 0)
	return
}

func (uc *AppVersionUseCase) Update(ctx context.Context, item *UpdateAppVersion) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Update(ctx, item)
			return
		})
	})
}

func (uc *AppVersionUseCase) Delete(ctx context.Context, ids ...uint64) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Delete(ctx, ids...)
			return
		})
	})
}
