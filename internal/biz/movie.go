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

type CreateMovie struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
}

type Movie struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
}

type FindMovie struct {
	Page  page.Page `json:"page"`
	Title *string   `json:"title"`
}

type FindMovieCache struct {
	Page page.Page `json:"page"`
	List []Movie   `json:"list"`
}

type UpdateMovie struct {
	ID    uint64  `json:"id,string"`
	Title *string `json:"title,omitempty"`
}

type MovieRepo interface {
	Create(ctx context.Context, item *CreateMovie) error
	Get(ctx context.Context, id uint64) (*Movie, error)
	Find(ctx context.Context, condition *FindMovie) []Movie
	Update(ctx context.Context, item *UpdateMovie) error
	Delete(ctx context.Context, ids ...uint64) error
}

type MovieUseCase struct {
	c     *conf.Bootstrap
	repo  MovieRepo
	tx    Transaction
	cache Cache
}

func NewMovieUseCase(c *conf.Bootstrap, repo MovieRepo, tx Transaction, cache Cache) *MovieUseCase {
	return &MovieUseCase{
		c:    c,
		repo: repo,
		tx:   tx,
		cache: cache.WithPrefix(strings.Join([]string{
			c.Name, "movie",
		}, "_")),
	}
}

func (uc *MovieUseCase) Create(ctx context.Context, item *CreateMovie) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Create(ctx, item)
		})
	})
}

func (uc *MovieUseCase) Get(ctx context.Context, id uint64) (rp *Movie, err error) {
	rp = &Movie{}
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

func (uc *MovieUseCase) get(ctx context.Context, action string, id uint64) (res string, err error) {
	// read data from db and write to cache
	rp := &Movie{}
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

func (uc *MovieUseCase) Find(ctx context.Context, condition *FindMovie) (rp []Movie, err error) {
	// use md5 string as cache replay json str, key is short
	action := strings.Join([]string{"find", utils.StructMd5(condition)}, "_")
	str, err := uc.cache.Get(ctx, action, func(ctx context.Context) (string, error) {
		return uc.find(ctx, action, condition)
	})
	if err != nil {
		return
	}
	var cache FindMovieCache
	utils.Json2Struct(&cache, str)
	condition.Page = cache.Page
	rp = cache.List
	return
}

func (uc *MovieUseCase) find(ctx context.Context, action string, condition *FindMovie) (res string, err error) {
	// read data from db and write to cache
	list := uc.repo.Find(ctx, condition)
	var cache FindMovieCache
	cache.List = list
	cache.Page = condition.Page
	res = utils.Struct2Json(cache)
	uc.cache.Set(ctx, action, res, len(list) == 0)
	return
}

func (uc *MovieUseCase) Update(ctx context.Context, item *UpdateMovie) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Update(ctx, item)
			return
		})
	})
}

func (uc *MovieUseCase) Delete(ctx context.Context, ids ...uint64) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Delete(ctx, ids...)
			return
		})
	})
}
