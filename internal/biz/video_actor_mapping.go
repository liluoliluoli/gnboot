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

type CreateVideoActorMapping struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
}

type VideoActorMapping struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
}

type FindVideoActorMapping struct {
	Page page.Page `json:"page"`
	Name *string   `json:"name"`
}

type FindVideoActorMappingCache struct {
	Page page.Page           `json:"page"`
	List []VideoActorMapping `json:"list"`
}

type UpdateVideoActorMapping struct {
	ID   uint64  `json:"id,string"`
	Name *string `json:"name,omitempty"`
}

type VideoActorMappingRepo interface {
	Create(ctx context.Context, item *CreateVideoActorMapping) error
	Get(ctx context.Context, id uint64) (*VideoActorMapping, error)
	Find(ctx context.Context, condition *FindVideoActorMapping) []VideoActorMapping
	Update(ctx context.Context, item *UpdateVideoActorMapping) error
	Delete(ctx context.Context, ids ...uint64) error
}

type VideoActorMappingUseCase struct {
	c     *conf.Bootstrap
	repo  VideoActorMappingRepo
	tx    Transaction
	cache Cache
}

func NewVideoActorMappingUseCase(c *conf.Bootstrap, repo VideoActorMappingRepo, tx Transaction, cache Cache) *VideoActorMappingUseCase {
	return &VideoActorMappingUseCase{
		c:    c,
		repo: repo,
		tx:   tx,
		cache: cache.WithPrefix(strings.Join([]string{
			c.Name, "video_actor_mapping",
		}, "_")),
	}
}

func (uc *VideoActorMappingUseCase) Create(ctx context.Context, item *CreateVideoActorMapping) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Create(ctx, item)
		})
	})
}

func (uc *VideoActorMappingUseCase) Get(ctx context.Context, id uint64) (rp *VideoActorMapping, err error) {
	rp = &VideoActorMapping{}
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

func (uc *VideoActorMappingUseCase) get(ctx context.Context, action string, id uint64) (res string, err error) {
	// read data from db and write to cache
	rp := &VideoActorMapping{}
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

func (uc *VideoActorMappingUseCase) Find(ctx context.Context, condition *FindVideoActorMapping) (rp []VideoActorMapping, err error) {
	// use md5 string as cache replay json str, key is short
	action := strings.Join([]string{"find", utils.StructMd5(condition)}, "_")
	str, err := uc.cache.Get(ctx, action, func(ctx context.Context) (string, error) {
		return uc.find(ctx, action, condition)
	})
	if err != nil {
		return
	}
	var cache FindVideoActorMappingCache
	utils.Json2Struct(&cache, str)
	condition.Page = cache.Page
	rp = cache.List
	return
}

func (uc *VideoActorMappingUseCase) find(ctx context.Context, action string, condition *FindVideoActorMapping) (res string, err error) {
	// read data from db and write to cache
	list := uc.repo.Find(ctx, condition)
	var cache FindVideoActorMappingCache
	cache.List = list
	cache.Page = condition.Page
	res = utils.Struct2Json(cache)
	uc.cache.Set(ctx, action, res, len(list) == 0)
	return
}

func (uc *VideoActorMappingUseCase) Update(ctx context.Context, item *UpdateVideoActorMapping) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Update(ctx, item)
			return
		})
	})
}

func (uc *VideoActorMappingUseCase) Delete(ctx context.Context, ids ...uint64) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Delete(ctx, ids...)
			return
		})
	})
}
