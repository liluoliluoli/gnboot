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

type CreateActor struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
}

type Actor struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
}

type FindActor struct {
	Page page.Page `json:"page"`
	Name *string   `json:"name"`
}

type FindActorCache struct {
	Page page.Page `json:"page"`
	List []Actor   `json:"list"`
}

type UpdateActor struct {
	ID   uint64  `json:"id,string"`
	Name *string `json:"name,omitempty"`
}

type ActorRepo interface {
	Create(ctx context.Context, item *CreateActor) error
	Get(ctx context.Context, id uint64) (*Actor, error)
	Find(ctx context.Context, condition *FindActor) []Actor
	Update(ctx context.Context, item *UpdateActor) error
	Delete(ctx context.Context, ids ...uint64) error
}

type ActorUseCase struct {
	c     *conf.Bootstrap
	repo  ActorRepo
	tx    Transaction
	cache Cache
}

func NewActorUseCase(c *conf.Bootstrap, repo ActorRepo, tx Transaction, cache Cache) *ActorUseCase {
	return &ActorUseCase{
		c:    c,
		repo: repo,
		tx:   tx,
		cache: cache.WithPrefix(strings.Join([]string{
			c.Name, "actor",
		}, "_")),
	}
}

func (uc *ActorUseCase) Create(ctx context.Context, item *CreateActor) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Create(ctx, item)
		})
	})
}

func (uc *ActorUseCase) Get(ctx context.Context, id uint64) (rp *Actor, err error) {
	rp = &Actor{}
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

func (uc *ActorUseCase) get(ctx context.Context, action string, id uint64) (res string, err error) {
	// read data from db and write to cache
	rp := &Actor{}
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

func (uc *ActorUseCase) Find(ctx context.Context, condition *FindActor) (rp []Actor, err error) {
	// use md5 string as cache replay json str, key is short
	action := strings.Join([]string{"find", utils.StructMd5(condition)}, "_")
	str, err := uc.cache.Get(ctx, action, func(ctx context.Context) (string, error) {
		return uc.find(ctx, action, condition)
	})
	if err != nil {
		return
	}
	var cache FindActorCache
	utils.Json2Struct(&cache, str)
	condition.Page = cache.Page
	rp = cache.List
	return
}

func (uc *ActorUseCase) find(ctx context.Context, action string, condition *FindActor) (res string, err error) {
	// read data from db and write to cache
	list := uc.repo.Find(ctx, condition)
	var cache FindActorCache
	cache.List = list
	cache.Page = condition.Page
	res = utils.Struct2Json(cache)
	uc.cache.Set(ctx, action, res, len(list) == 0)
	return
}

func (uc *ActorUseCase) Update(ctx context.Context, item *UpdateActor) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Update(ctx, item)
			return
		})
	})
}

func (uc *ActorUseCase) Delete(ctx context.Context, ids ...uint64) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Delete(ctx, ids...)
			return
		})
	})
}
