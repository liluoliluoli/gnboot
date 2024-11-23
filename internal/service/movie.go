package service

import (
	"context"
	"gnboot/internal/utils/cache_util"
	"strings"

	"github.com/go-cinch/common/page"
	"github.com/go-cinch/common/utils"
	"gnboot/internal/conf"
)

type CreateMovie struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
}

type Movie struct {
	ID            uint64  `json:"id,string"`
	OriginalTitle string  `json:"originalTitle"` // 标题
	Status        string  `json:"status"`        // 状态，Returning Series, Ended, Released, Unknown
	VoteAverage   float32 `json:"voteAverage"`   // 平均评分
	VoteCount     int64   `json:"voteCount"`     // 评分数
	Country       string  `json:"country"`       // 国家
	Trailer       string  `json:"trailer"`       // 预告片地址
	URL           string  `json:"url"`           // 影片地址
	Downloaded    bool    `json:"downloaded"`    // 是否可以下载
	FileSize      int64   `json:"fileSize"`      // 文件大小
	Filename      string  `json:"filename"`      // 文件名
	Ext           string  `json:"ext"`           //扩展参数
	//Genres             []*Genre               `json:"genres"`             //流派
	//Studios            []*Studio              `json:"studios"`            //出品方
	Keywords           []string `json:"keywords"`           //关键词
	LastPlayedPosition int64    `json:"lastPlayedPosition"` //上次播放位置
	LastPlayedTime     string   `json:"lastPlayedTime"`     //YYYY-MM-DD HH:MM:SS
	//Subtitles          []VideoSubtitleMapping `json:"subtitles"`          //字幕
	//Actors             []*Actor               `json:"actors"`             //演员
}

func (*Movie) ConvertFromRepo() {

}

type FindMovie struct {
	Page   page.Page `json:"page"`
	Search *string   `json:"search"`
	Sort   *Sort     `json:"sort"`
}

type Sort struct {
	Filter    *string `json:"filter"`
	Type      *string `json:"type"`
	Direction *string `json:"direction"`
}

type UpdateMovie struct {
	ID    uint64  `json:"id,string"`
	Title *string `json:"title,omitempty"`
}

type MovieRepo interface {
	Create(ctx context.Context, item *CreateMovie) error
	Get(ctx context.Context, id uint64) (*Movie, error)
	Find(ctx context.Context, condition *FindMovie) []*Movie
	Update(ctx context.Context, item *UpdateMovie) error
	Delete(ctx context.Context, ids ...uint64) error
}

type MovieUseCase struct {
	c     *conf.Bootstrap
	repo  MovieRepo
	tx    Transaction
	cache Cache[*Movie]
}

func NewMovieUseCase(c *conf.Bootstrap, repo MovieRepo, tx Transaction, cache Cache[*Movie]) *MovieUseCase {
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
	rp, err = uc.cache.Get(ctx, cache_util.GetCacheActionName(id), func(action string, ctx context.Context) (*Movie, error) {
		return uc.get(ctx, id)
	})
	return
}

func (uc *MovieUseCase) get(ctx context.Context, id uint64) (*Movie, error) {
	item, err := uc.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (uc *MovieUseCase) Find(ctx context.Context, condition *FindMovie) ([]*Movie, error) {
	rp, err := uc.cache.GetPage(ctx, cache_util.GetCacheActionName(condition), func(action string, ctx context.Context) ([]*Movie, error) {
		return uc.find(ctx, action, condition)
	})
	if err != nil {
		return nil, err
	}
	var cache FindMovieCache
	utils.Json2Struct(&cache, str)
	condition.Page = cache.Page
	rp = cache.List
	return
}

func (uc *MovieUseCase) find(ctx context.Context, action string, condition *FindMovie) (res []*Movie, err error) {
	// read repo from db and write to cache
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
