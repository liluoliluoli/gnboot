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

type CreateKeyword struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
}

type Keyword struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
}

type FindKeyword struct {
	Page page.Page `json:"page"`
	Name *string   `json:"name"`
}

type FindKeywordCache struct {
	Page page.Page `json:"page"`
	List []Keyword `json:"list"`
}

type UpdateKeyword struct {
	ID   uint64  `json:"id,string"`
	Name *string `json:"name,omitempty"`
}

type KeywordRepo interface {
	Create(ctx context.Context, item *CreateKeyword) error
	Get(ctx context.Context, id uint64) (*Keyword, error)
	Find(ctx context.Context, condition *FindKeyword) []Keyword
	Update(ctx context.Context, item *UpdateKeyword) error
	Delete(ctx context.Context, ids ...uint64) error
}

type KeywordUseCase struct {
	c     *conf.Bootstrap
	repo  KeywordRepo
	tx    Transaction
	cache Cache
}

func NewKeywordUseCase(c *conf.Bootstrap, repo KeywordRepo, tx Transaction, cache Cache) *KeywordUseCase {
	return &KeywordUseCase{
		c:    c,
		repo: repo,
		tx:   tx,
		cache: cache.WithPrefix(strings.Join([]string{
			c.Name, "keyword",
		}, "_")),
	}
}

func (uc *KeywordUseCase) Create(ctx context.Context, item *CreateKeyword) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Create(ctx, item)
		})
	})
}

func (uc *KeywordUseCase) Get(ctx context.Context, id uint64) (rp *Keyword, err error) {
	rp = &Keyword{}
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

func (uc *KeywordUseCase) get(ctx context.Context, action string, id uint64) (res string, err error) {
	// read data from db and write to cache
	rp := &Keyword{}
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

func (uc *KeywordUseCase) Find(ctx context.Context, condition *FindKeyword) (rp []Keyword, err error) {
	// use md5 string as cache replay json str, key is short
	action := strings.Join([]string{"find", utils.StructMd5(condition)}, "_")
	str, err := uc.cache.Get(ctx, action, func(ctx context.Context) (string, error) {
		return uc.find(ctx, action, condition)
	})
	if err != nil {
		return
	}
	var cache FindKeywordCache
	utils.Json2Struct(&cache, str)
	condition.Page = cache.Page
	rp = cache.List
	return
}

func (uc *KeywordUseCase) find(ctx context.Context, action string, condition *FindKeyword) (res string, err error) {
	// read data from db and write to cache
	list := uc.repo.Find(ctx, condition)
	var cache FindKeywordCache
	cache.List = list
	cache.Page = condition.Page
	res = utils.Struct2Json(cache)
	uc.cache.Set(ctx, action, res, len(list) == 0)
	return
}

func (uc *KeywordUseCase) Update(ctx context.Context, item *UpdateKeyword) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Update(ctx, item)
			return
		})
	})
}

func (uc *KeywordUseCase) Delete(ctx context.Context, ids ...uint64) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Delete(ctx, ids...)
			return
		})
	})
}