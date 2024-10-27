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

type CreateGnboot struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
}

type Gnboot struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
}

type FindGnboot struct {
	Page page.Page `json:"page"`
	Name *string   `json:"name"`
}

type FindGnbootCache struct {
	Page page.Page `json:"page"`
	List []Gnboot  `json:"list"`
}

type UpdateGnboot struct {
	ID   uint64  `json:"id,string"`
	Name *string `json:"name,omitempty"`
}

type GnbootRepo interface {
	Create(ctx context.Context, item *CreateGnboot) error
	Get(ctx context.Context, id uint64) (*Gnboot, error)
	Find(ctx context.Context, condition *FindGnboot) []Gnboot
	Update(ctx context.Context, item *UpdateGnboot) error
	Delete(ctx context.Context, ids ...uint64) error
}

type GnbootUseCase struct {
	c     *conf.Bootstrap
	repo  GnbootRepo
	tx    Transaction
	cache Cache
}

func NewGnbootUseCase(c *conf.Bootstrap, repo GnbootRepo, tx Transaction, cache Cache) *GnbootUseCase {
	return &GnbootUseCase{
		c:    c,
		repo: repo,
		tx:   tx,
		cache: cache.WithPrefix(strings.Join([]string{
			c.Name, "gnboot",
		}, "_")),
	}
}

func (uc *GnbootUseCase) Create(ctx context.Context, item *CreateGnboot) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Create(ctx, item)
		})
	})
}

func (uc *GnbootUseCase) Get(ctx context.Context, id uint64) (rp *Gnboot, err error) {
	rp = &Gnboot{}
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

func (uc *GnbootUseCase) get(ctx context.Context, action string, id uint64) (res string, err error) {
	// read data from db and write to cache
	rp := &Gnboot{}
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

func (uc *GnbootUseCase) Find(ctx context.Context, condition *FindGnboot) (rp []Gnboot, err error) {
	// use md5 string as cache replay json str, key is short
	action := strings.Join([]string{"find", utils.StructMd5(condition)}, "_")
	str, err := uc.cache.Get(ctx, action, func(ctx context.Context) (string, error) {
		return uc.find(ctx, action, condition)
	})
	if err != nil {
		return
	}
	var cache FindGnbootCache
	utils.Json2Struct(&cache, str)
	condition.Page = cache.Page
	rp = cache.List
	return
}

func (uc *GnbootUseCase) find(ctx context.Context, action string, condition *FindGnboot) (res string, err error) {
	// read data from db and write to cache
	list := uc.repo.Find(ctx, condition)
	var cache FindGnbootCache
	cache.List = list
	cache.Page = condition.Page
	res = utils.Struct2Json(cache)
	uc.cache.Set(ctx, action, res, len(list) == 0)
	return
}

func (uc *GnbootUseCase) Update(ctx context.Context, item *UpdateGnboot) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Update(ctx, item)
			return
		})
	})
}

func (uc *GnbootUseCase) Delete(ctx context.Context, ids ...uint64) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Delete(ctx, ids...)
			return
		})
	})
}
