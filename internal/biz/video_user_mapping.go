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

type CreateVideoUserMapping struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
}

type VideoUserMapping struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
}

type FindVideoUserMapping struct {
	Page page.Page `json:"page"`
	Name *string   `json:"name"`
}

type FindVideoUserMappingCache struct {
	Page page.Page          `json:"page"`
	List []VideoUserMapping `json:"list"`
}

type UpdateVideoUserMapping struct {
	ID   uint64  `json:"id,string"`
	Name *string `json:"name,omitempty"`
}

type VideoUserMappingRepo interface {
	Create(ctx context.Context, item *CreateVideoUserMapping) error
	Get(ctx context.Context, id uint64) (*VideoUserMapping, error)
	Find(ctx context.Context, condition *FindVideoUserMapping) []VideoUserMapping
	Update(ctx context.Context, item *UpdateVideoUserMapping) error
	Delete(ctx context.Context, ids ...uint64) error
}

type VideoUserMappingUseCase struct {
	c     *conf.Bootstrap
	repo  VideoUserMappingRepo
	tx    Transaction
	cache Cache
}

func NewVideoUserMappingUseCase(c *conf.Bootstrap, repo VideoUserMappingRepo, tx Transaction, cache Cache) *VideoUserMappingUseCase {
	return &VideoUserMappingUseCase{
		c:    c,
		repo: repo,
		tx:   tx,
		cache: cache.WithPrefix(strings.Join([]string{
			c.Name, "video_user_mapping",
		}, "_")),
	}
}

func (uc *VideoUserMappingUseCase) Create(ctx context.Context, item *CreateVideoUserMapping) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Create(ctx, item)
		})
	})
}

func (uc *VideoUserMappingUseCase) Get(ctx context.Context, id uint64) (rp *VideoUserMapping, err error) {
	rp = &VideoUserMapping{}
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

func (uc *VideoUserMappingUseCase) get(ctx context.Context, action string, id uint64) (res string, err error) {
	// read data from db and write to cache
	rp := &VideoUserMapping{}
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

func (uc *VideoUserMappingUseCase) Find(ctx context.Context, condition *FindVideoUserMapping) (rp []VideoUserMapping, err error) {
	// use md5 string as cache replay json str, key is short
	str, err := uc.cache.Get(ctx, cache_util.GetCacheActionName(condition), func(action string, ctx context.Context) (string, error) {
		return uc.find(ctx, action, condition)
	})
	if err != nil {
		return
	}
	var cache FindVideoUserMappingCache
	utils.Json2Struct(&cache, str)
	condition.Page = cache.Page
	rp = cache.List
	return
}

func (uc *VideoUserMappingUseCase) find(ctx context.Context, action string, condition *FindVideoUserMapping) (res string, err error) {
	// read data from db and write to cache
	list := uc.repo.Find(ctx, condition)
	var cache FindVideoUserMappingCache
	cache.List = list
	cache.Page = condition.Page
	res = utils.Struct2Json(cache)
	uc.cache.Set(ctx, action, res, len(list) == 0)
	return
}

func (uc *VideoUserMappingUseCase) Update(ctx context.Context, item *UpdateVideoUserMapping) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Update(ctx, item)
			return
		})
	})
}

func (uc *VideoUserMappingUseCase) Delete(ctx context.Context, ids ...uint64) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Delete(ctx, ids...)
			return
		})
	})
}
