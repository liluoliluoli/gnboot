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

type CreateVideoGenreMapping struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
}

type VideoGenreMapping struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
}

type FindVideoGenreMapping struct {
	Page page.Page `json:"page"`
	Name *string   `json:"name"`
}

type FindVideoGenreMappingCache struct {
	Page page.Page           `json:"page"`
	List []VideoGenreMapping `json:"list"`
}

type UpdateVideoGenreMapping struct {
	ID   uint64  `json:"id,string"`
	Name *string `json:"name,omitempty"`
}

type VideoGenreMappingRepo interface {
	Create(ctx context.Context, item *CreateVideoGenreMapping) error
	Get(ctx context.Context, id uint64) (*VideoGenreMapping, error)
	Find(ctx context.Context, condition *FindVideoGenreMapping) []VideoGenreMapping
	Update(ctx context.Context, item *UpdateVideoGenreMapping) error
	Delete(ctx context.Context, ids ...uint64) error
}

type VideoGenreMappingUseCase struct {
	c     *conf.Bootstrap
	repo  VideoGenreMappingRepo
	tx    Transaction
	cache Cache
}

func NewVideoGenreMappingUseCase(c *conf.Bootstrap, repo VideoGenreMappingRepo, tx Transaction, cache Cache) *VideoGenreMappingUseCase {
	return &VideoGenreMappingUseCase{
		c:    c,
		repo: repo,
		tx:   tx,
		cache: cache.WithPrefix(strings.Join([]string{
			c.Name, "video_genre_mapping",
		}, "_")),
	}
}

func (uc *VideoGenreMappingUseCase) Create(ctx context.Context, item *CreateVideoGenreMapping) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Create(ctx, item)
		})
	})
}

func (uc *VideoGenreMappingUseCase) Get(ctx context.Context, id uint64) (rp *VideoGenreMapping, err error) {
	rp = &VideoGenreMapping{}
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

func (uc *VideoGenreMappingUseCase) get(ctx context.Context, action string, id uint64) (res string, err error) {
	// read data from db and write to cache
	rp := &VideoGenreMapping{}
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

func (uc *VideoGenreMappingUseCase) Find(ctx context.Context, condition *FindVideoGenreMapping) (rp []VideoGenreMapping, err error) {
	// use md5 string as cache replay json str, key is short
	action := strings.Join([]string{"find", utils.StructMd5(condition)}, "_")
	str, err := uc.cache.Get(ctx, action, func(ctx context.Context) (string, error) {
		return uc.find(ctx, action, condition)
	})
	if err != nil {
		return
	}
	var cache FindVideoGenreMappingCache
	utils.Json2Struct(&cache, str)
	condition.Page = cache.Page
	rp = cache.List
	return
}

func (uc *VideoGenreMappingUseCase) find(ctx context.Context, action string, condition *FindVideoGenreMapping) (res string, err error) {
	// read data from db and write to cache
	list := uc.repo.Find(ctx, condition)
	var cache FindVideoGenreMappingCache
	cache.List = list
	cache.Page = condition.Page
	res = utils.Struct2Json(cache)
	uc.cache.Set(ctx, action, res, len(list) == 0)
	return
}

func (uc *VideoGenreMappingUseCase) Update(ctx context.Context, item *UpdateVideoGenreMapping) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Update(ctx, item)
			return
		})
	})
}

func (uc *VideoGenreMappingUseCase) Delete(ctx context.Context, ids ...uint64) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Delete(ctx, ids...)
			return
		})
	})
}
