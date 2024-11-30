package service

import (
	"context"
	"github.com/liluoliluoli/gnboot/internal/common/utils/cache_util"
	"github.com/liluoliluoli/gnboot/internal/conf"
	"github.com/liluoliluoli/gnboot/internal/repo"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
)

type KeywordService struct {
	c           *conf.Bootstrap
	keywordRepo *repo.KeywordRepo
	cache       sdomain.Cache[*sdomain.Keyword]
}

func NewKeywordService(c *conf.Bootstrap,
	keywordRepo *repo.KeywordRepo,
) *KeywordService {
	return &KeywordService{
		c:           c,
		keywordRepo: keywordRepo,
		cache:       repo.NewCache[*sdomain.Keyword](c, keywordRepo.Data.Cache()),
	}
}

func (s *KeywordService) Get(ctx context.Context, id int64) (*sdomain.Keyword, error) {
	return s.cache.Get(ctx, cache_util.GetCacheActionName(id), func(action string, ctx context.Context) (*sdomain.Keyword, error) {
		return s.get(ctx, id)
	})
}

func (s *KeywordService) get(ctx context.Context, id int64) (*sdomain.Keyword, error) {
	item, err := s.keywordRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (s *KeywordService) FindAll(ctx context.Context) ([]*sdomain.Keyword, error) {
	rp, err := s.cache.List(ctx, cache_util.GetCacheActionName(""), func(action string, ctx context.Context) ([]*sdomain.Keyword, error) {
		return s.findAll(ctx)
	})
	if err != nil {
		return nil, err
	}
	return rp, nil
}

func (s *KeywordService) findAll(ctx context.Context) ([]*sdomain.Keyword, error) {
	finds, err := s.keywordRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return finds, nil
}
