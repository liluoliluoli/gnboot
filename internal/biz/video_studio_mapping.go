package biz

import (
	"context"
	"strconv"
	"strings"

	"github.com/go-cinch/common/constant"
	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/page"
	"github.com/go-cinch/common/utils"
	"github.com/pkg/errors"
	"gnboot/internal/conf"
)

type CreateVideoStudioMapping struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
}

type VideoStudioMapping struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
}

type FindVideoStudioMapping struct {
	Page page.Page `json:"page"`
	Name *string   `json:"name"`
}

type FindVideoStudioMappingCache struct {
	Page page.Page            `json:"page"`
	List []VideoStudioMapping `json:"list"`
}

type UpdateVideoStudioMapping struct {
	ID   uint64  `json:"id,string"`
	Name *string `json:"name,omitempty"`
}

type VideoStudioMappingRepo interface {
	Create(ctx context.Context, item *CreateVideoStudioMapping) error
	Get(ctx context.Context, id uint64) (*VideoStudioMapping, error)
	Find(ctx context.Context, condition *FindVideoStudioMapping) []VideoStudioMapping
	Update(ctx context.Context, item *UpdateVideoStudioMapping) error
	Delete(ctx context.Context, ids ...uint64) error
}

type VideoStudioMappingUseCase struct {
	c     *conf.Bootstrap
	repo  VideoStudioMappingRepo
	tx    Transaction
	cache Cache
}

func NewVideoStudioMappingUseCase(c *conf.Bootstrap, repo VideoStudioMappingRepo, tx Transaction, cache Cache) *VideoStudioMappingUseCase {
	return &VideoStudioMappingUseCase{
		c:    c,
		repo: repo,
		tx:   tx,
		cache: cache.WithPrefix(strings.Join([]string{
			c.Name, "video_studio_mapping",
		}, "_")),
	}
}

func (uc *VideoStudioMappingUseCase) Create(ctx context.Context, item *CreateVideoStudioMapping) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Create(ctx, item)
		})
	})
}

func (uc *VideoStudioMappingUseCase) Get(ctx context.Context, id uint64) (rp *VideoStudioMapping, err error) {
	rp = &VideoStudioMapping{}
	action := strings.Join([]string{"get", strconv.FormatUint(id, 10)}, "_")
	str, err := uc.cache.Get(ctx, action, func(ctx context.Context) (string, error) {
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

func (uc *VideoStudioMappingUseCase) get(ctx context.Context, action string, id uint64) (res string, err error) {
	// read data from db and write to cache
	rp := &VideoStudioMapping{}
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

func (uc *VideoStudioMappingUseCase) Find(ctx context.Context, condition *FindVideoStudioMapping) (rp []VideoStudioMapping, err error) {
	// use md5 string as cache replay json str, key is short
	action := strings.Join([]string{"find", utils.StructMd5(condition)}, "_")
	str, err := uc.cache.Get(ctx, action, func(ctx context.Context) (string, error) {
		return uc.find(ctx, action, condition)
	})
	if err != nil {
		return
	}
	var cache FindVideoStudioMappingCache
	utils.Json2Struct(&cache, str)
	condition.Page = cache.Page
	rp = cache.List
	return
}

func (uc *VideoStudioMappingUseCase) find(ctx context.Context, action string, condition *FindVideoStudioMapping) (res string, err error) {
	// read data from db and write to cache
	list := uc.repo.Find(ctx, condition)
	var cache FindVideoStudioMappingCache
	cache.List = list
	cache.Page = condition.Page
	res = utils.Struct2Json(cache)
	uc.cache.Set(ctx, action, res, len(list) == 0)
	return
}

func (uc *VideoStudioMappingUseCase) Update(ctx context.Context, item *UpdateVideoStudioMapping) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Update(ctx, item)
			return
		})
	})
}

func (uc *VideoStudioMappingUseCase) Delete(ctx context.Context, ids ...uint64) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Delete(ctx, ids...)
			return
		})
	})
}