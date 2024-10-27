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

type CreateSeason struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
}

type Season struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
}

type FindSeason struct {
	Page page.Page `json:"page"`
	Name *string   `json:"name"`
}

type FindSeasonCache struct {
	Page page.Page `json:"page"`
	List []Season  `json:"list"`
}

type UpdateSeason struct {
	ID   uint64  `json:"id,string"`
	Name *string `json:"name,omitempty"`
}

type SeasonRepo interface {
	Create(ctx context.Context, item *CreateSeason) error
	Get(ctx context.Context, id uint64) (*Season, error)
	Find(ctx context.Context, condition *FindSeason) []Season
	Update(ctx context.Context, item *UpdateSeason) error
	Delete(ctx context.Context, ids ...uint64) error
}

type SeasonUseCase struct {
	c     *conf.Bootstrap
	repo  SeasonRepo
	tx    Transaction
	cache Cache
}

func NewSeasonUseCase(c *conf.Bootstrap, repo SeasonRepo, tx Transaction, cache Cache) *SeasonUseCase {
	return &SeasonUseCase{
		c:    c,
		repo: repo,
		tx:   tx,
		cache: cache.WithPrefix(strings.Join([]string{
			c.Name, "season",
		}, "_")),
	}
}

func (uc *SeasonUseCase) Create(ctx context.Context, item *CreateSeason) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Create(ctx, item)
		})
	})
}

func (uc *SeasonUseCase) Get(ctx context.Context, id uint64) (rp *Season, err error) {
	rp = &Season{}
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

func (uc *SeasonUseCase) get(ctx context.Context, action string, id uint64) (res string, err error) {
	// read data from db and write to cache
	rp := &Season{}
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

func (uc *SeasonUseCase) Find(ctx context.Context, condition *FindSeason) (rp []Season, err error) {
	// use md5 string as cache replay json str, key is short
	action := strings.Join([]string{"find", utils.StructMd5(condition)}, "_")
	str, err := uc.cache.Get(ctx, action, func(ctx context.Context) (string, error) {
		return uc.find(ctx, action, condition)
	})
	if err != nil {
		return
	}
	var cache FindSeasonCache
	utils.Json2Struct(&cache, str)
	condition.Page = cache.Page
	rp = cache.List
	return
}

func (uc *SeasonUseCase) find(ctx context.Context, action string, condition *FindSeason) (res string, err error) {
	// read data from db and write to cache
	list := uc.repo.Find(ctx, condition)
	var cache FindSeasonCache
	cache.List = list
	cache.Page = condition.Page
	res = utils.Struct2Json(cache)
	uc.cache.Set(ctx, action, res, len(list) == 0)
	return
}

func (uc *SeasonUseCase) Update(ctx context.Context, item *UpdateSeason) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Update(ctx, item)
			return
		})
	})
}

func (uc *SeasonUseCase) Delete(ctx context.Context, ids ...uint64) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Delete(ctx, ids...)
			return
		})
	})
}
