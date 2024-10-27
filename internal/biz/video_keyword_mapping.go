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

type CreateVideoKeywordMapping struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
}

type VideoKeywordMapping struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
}

type FindVideoKeywordMapping struct {
	Page page.Page `json:"page"`
	Name *string   `json:"name"`
}

type FindVideoKeywordMappingCache struct {
	Page page.Page             `json:"page"`
	List []VideoKeywordMapping `json:"list"`
}

type UpdateVideoKeywordMapping struct {
	ID   uint64  `json:"id,string"`
	Name *string `json:"name,omitempty"`
}

type VideoKeywordMappingRepo interface {
	Create(ctx context.Context, item *CreateVideoKeywordMapping) error
	Get(ctx context.Context, id uint64) (*VideoKeywordMapping, error)
	Find(ctx context.Context, condition *FindVideoKeywordMapping) []VideoKeywordMapping
	Update(ctx context.Context, item *UpdateVideoKeywordMapping) error
	Delete(ctx context.Context, ids ...uint64) error
}

type VideoKeywordMappingUseCase struct {
	c     *conf.Bootstrap
	repo  VideoKeywordMappingRepo
	tx    Transaction
	cache Cache
}

func NewVideoKeywordMappingUseCase(c *conf.Bootstrap, repo VideoKeywordMappingRepo, tx Transaction, cache Cache) *VideoKeywordMappingUseCase {
	return &VideoKeywordMappingUseCase{
		c:    c,
		repo: repo,
		tx:   tx,
		cache: cache.WithPrefix(strings.Join([]string{
			c.Name, "video_keyword_mapping",
		}, "_")),
	}
}

func (uc *VideoKeywordMappingUseCase) Create(ctx context.Context, item *CreateVideoKeywordMapping) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Create(ctx, item)
		})
	})
}

func (uc *VideoKeywordMappingUseCase) Get(ctx context.Context, id uint64) (rp *VideoKeywordMapping, err error) {
	rp = &VideoKeywordMapping{}
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

func (uc *VideoKeywordMappingUseCase) get(ctx context.Context, action string, id uint64) (res string, err error) {
	// read data from db and write to cache
	rp := &VideoKeywordMapping{}
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

func (uc *VideoKeywordMappingUseCase) Find(ctx context.Context, condition *FindVideoKeywordMapping) (rp []VideoKeywordMapping, err error) {
	// use md5 string as cache replay json str, key is short
	action := strings.Join([]string{"find", utils.StructMd5(condition)}, "_")
	str, err := uc.cache.Get(ctx, action, func(ctx context.Context) (string, error) {
		return uc.find(ctx, action, condition)
	})
	if err != nil {
		return
	}
	var cache FindVideoKeywordMappingCache
	utils.Json2Struct(&cache, str)
	condition.Page = cache.Page
	rp = cache.List
	return
}

func (uc *VideoKeywordMappingUseCase) find(ctx context.Context, action string, condition *FindVideoKeywordMapping) (res string, err error) {
	// read data from db and write to cache
	list := uc.repo.Find(ctx, condition)
	var cache FindVideoKeywordMappingCache
	cache.List = list
	cache.Page = condition.Page
	res = utils.Struct2Json(cache)
	uc.cache.Set(ctx, action, res, len(list) == 0)
	return
}

func (uc *VideoKeywordMappingUseCase) Update(ctx context.Context, item *UpdateVideoKeywordMapping) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Update(ctx, item)
			return
		})
	})
}

func (uc *VideoKeywordMappingUseCase) Delete(ctx context.Context, ids ...uint64) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Delete(ctx, ids...)
			return
		})
	})
}
