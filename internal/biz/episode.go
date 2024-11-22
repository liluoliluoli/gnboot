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

type CreateEpisode struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
}

type Episode struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
}

type FindEpisode struct {
	Page page.Page `json:"page"`
	Name *string   `json:"name"`
}

type FindEpisodeCache struct {
	Page page.Page `json:"page"`
	List []Episode `json:"list"`
}

type UpdateEpisode struct {
	ID   uint64  `json:"id,string"`
	Name *string `json:"name,omitempty"`
}

type EpisodeRepo interface {
	Create(ctx context.Context, item *CreateEpisode) error
	Get(ctx context.Context, id uint64) (*Episode, error)
	Find(ctx context.Context, condition *FindEpisode) []Episode
	Update(ctx context.Context, item *UpdateEpisode) error
	Delete(ctx context.Context, ids ...uint64) error
}

type EpisodeUseCase struct {
	c     *conf.Bootstrap
	repo  EpisodeRepo
	tx    Transaction
	cache Cache
}

func NewEpisodeUseCase(c *conf.Bootstrap, repo EpisodeRepo, tx Transaction, cache Cache) *EpisodeUseCase {
	return &EpisodeUseCase{
		c:    c,
		repo: repo,
		tx:   tx,
		cache: cache.WithPrefix(strings.Join([]string{
			c.Name, "episode",
		}, "_")),
	}
}

func (uc *EpisodeUseCase) Create(ctx context.Context, item *CreateEpisode) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Create(ctx, item)
		})
	})
}

func (uc *EpisodeUseCase) Get(ctx context.Context, id uint64) (rp *Episode, err error) {
	rp = &Episode{}
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

func (uc *EpisodeUseCase) get(ctx context.Context, action string, id uint64) (res string, err error) {
	// read data from db and write to cache
	rp := &Episode{}
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

func (uc *EpisodeUseCase) Find(ctx context.Context, condition *FindEpisode) (rp []Episode, err error) {
	// use md5 string as cache replay json str, key is short
	str, err := uc.cache.Get(ctx, cache_util.GetCacheActionName(condition), func(action string, ctx context.Context) (string, error) {
		return uc.find(ctx, action, condition)
	})
	if err != nil {
		return
	}
	var cache FindEpisodeCache
	utils.Json2Struct(&cache, str)
	condition.Page = cache.Page
	rp = cache.List
	return
}

func (uc *EpisodeUseCase) find(ctx context.Context, action string, condition *FindEpisode) (res string, err error) {
	// read data from db and write to cache
	list := uc.repo.Find(ctx, condition)
	var cache FindEpisodeCache
	cache.List = list
	cache.Page = condition.Page
	res = utils.Struct2Json(cache)
	uc.cache.Set(ctx, action, res, len(list) == 0)
	return
}

func (uc *EpisodeUseCase) Update(ctx context.Context, item *UpdateEpisode) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Update(ctx, item)
			return
		})
	})
}

func (uc *EpisodeUseCase) Delete(ctx context.Context, ids ...uint64) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Delete(ctx, ids...)
			return
		})
	})
}
