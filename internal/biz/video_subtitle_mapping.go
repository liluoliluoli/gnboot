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

type CreateVideoSubtitleMapping struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
}

type VideoSubtitleMapping struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
}

type FindVideoSubtitleMapping struct {
	Page page.Page `json:"page"`
	Name *string   `json:"name"`
}

type FindVideoSubtitleMappingCache struct {
	Page page.Page              `json:"page"`
	List []VideoSubtitleMapping `json:"list"`
}

type UpdateVideoSubtitleMapping struct {
	ID   uint64  `json:"id,string"`
	Name *string `json:"name,omitempty"`
}

type VideoSubtitleMappingRepo interface {
	Create(ctx context.Context, item *CreateVideoSubtitleMapping) error
	Get(ctx context.Context, id uint64) (*VideoSubtitleMapping, error)
	Find(ctx context.Context, condition *FindVideoSubtitleMapping) []VideoSubtitleMapping
	Update(ctx context.Context, item *UpdateVideoSubtitleMapping) error
	Delete(ctx context.Context, ids ...uint64) error
}

type VideoSubtitleMappingUseCase struct {
	c     *conf.Bootstrap
	repo  VideoSubtitleMappingRepo
	tx    Transaction
	cache Cache
}

func NewVideoSubtitleMappingUseCase(c *conf.Bootstrap, repo VideoSubtitleMappingRepo, tx Transaction, cache Cache) *VideoSubtitleMappingUseCase {
	return &VideoSubtitleMappingUseCase{
		c:    c,
		repo: repo,
		tx:   tx,
		cache: cache.WithPrefix(strings.Join([]string{
			c.Name, "video_subtitle_mapping",
		}, "_")),
	}
}

func (uc *VideoSubtitleMappingUseCase) Create(ctx context.Context, item *CreateVideoSubtitleMapping) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Create(ctx, item)
		})
	})
}

func (uc *VideoSubtitleMappingUseCase) Get(ctx context.Context, id uint64) (rp *VideoSubtitleMapping, err error) {
	rp = &VideoSubtitleMapping{}
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

func (uc *VideoSubtitleMappingUseCase) get(ctx context.Context, action string, id uint64) (res string, err error) {
	// read data from db and write to cache
	rp := &VideoSubtitleMapping{}
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

func (uc *VideoSubtitleMappingUseCase) Find(ctx context.Context, condition *FindVideoSubtitleMapping) (rp []VideoSubtitleMapping, err error) {
	// use md5 string as cache replay json str, key is short
	action := strings.Join([]string{"find", utils.StructMd5(condition)}, "_")
	str, err := uc.cache.Get(ctx, action, func(ctx context.Context) (string, error) {
		return uc.find(ctx, action, condition)
	})
	if err != nil {
		return
	}
	var cache FindVideoSubtitleMappingCache
	utils.Json2Struct(&cache, str)
	condition.Page = cache.Page
	rp = cache.List
	return
}

func (uc *VideoSubtitleMappingUseCase) find(ctx context.Context, action string, condition *FindVideoSubtitleMapping) (res string, err error) {
	// read data from db and write to cache
	list := uc.repo.Find(ctx, condition)
	var cache FindVideoSubtitleMappingCache
	cache.List = list
	cache.Page = condition.Page
	res = utils.Struct2Json(cache)
	uc.cache.Set(ctx, action, res, len(list) == 0)
	return
}

func (uc *VideoSubtitleMappingUseCase) Update(ctx context.Context, item *UpdateVideoSubtitleMapping) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Update(ctx, item)
			return
		})
	})
}

func (uc *VideoSubtitleMappingUseCase) Delete(ctx context.Context, ids ...uint64) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Delete(ctx, ids...)
			return
		})
	})
}
