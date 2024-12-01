package sdomain

import (
	"github.com/liluoliluoli/gnboot/api/genre"
	"github.com/liluoliluoli/gnboot/api/keyword"
	"github.com/liluoliluoli/gnboot/api/season"
	seriesdto "github.com/liluoliluoli/gnboot/api/series"
	"github.com/liluoliluoli/gnboot/api/studio"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/samber/lo"
)

type Series struct {
	ID          int64      `json:"id"`
	VoteAverage float32    `json:"voteAverage"` // 平均评分
	VoteCount   int32      `json:"voteCount"`   // 评分数
	Country     string     `json:"country"`     // 国家
	Trailer     string     `json:"trailer"`     // 预告片地址
	Status      string     `json:"status"`      // 状态，Returning Series, Ended, Released, Unknown
	SkipIntro   int32      `json:"skipIntro"`
	SkipEnding  int32      `json:"skipEnding"`
	Genres      []*Genre   `json:"genres"`     //流派
	Studios     []*Studio  `json:"studios"`    //出品方
	Keywords    []*Keyword `json:"keywords"`   //关键词
	Seasons     []*Season  `json:"seasons"`    //季
	NextToPlay  *Episode   `json:"nextToPlay"` //下一集
}

func (d *Series) ConvertFromRepo(m *model.Series) *Series {
	return &Series{
		ID:          m.ID,
		VoteAverage: lo.FromPtr(m.VoteAverage),
		VoteCount:   lo.FromPtr(m.VoteCount),
		Country:     lo.FromPtr(m.Country),
		Trailer:     lo.FromPtr(m.Trailer),
		Status:      m.Status,
		SkipIntro:   lo.FromPtr(m.SkipIntro),
		SkipEnding:  lo.FromPtr(m.SkipEnding),
	}
}

func (d *Series) ConvertToDto() *seriesdto.SeriesResp {
	return &seriesdto.SeriesResp{
		Id:          d.ID,
		VoteAverage: d.VoteAverage,
		VoteCount:   d.VoteCount,
		Country:     d.Country,
		Trailer:     d.Trailer,
		Status:      d.Status,
		SkipIntro:   d.SkipIntro,
		SkipEnding:  d.SkipEnding,
		Genres: lo.Map(d.Genres, func(item *Genre, index int) *genre.GenreResp {
			return item.ConvertToDto()
		}),
		Studios: lo.Map(d.Studios, func(item *Studio, index int) *studio.StudioResp {
			return item.ConvertToDto()
		}),
		Keywords: lo.Map(d.Keywords, func(item *Keyword, index int) *keyword.KeywordResp {
			return item.ConvertToDto()
		}),
		Seasons: lo.Map(d.Seasons, func(item *Season, index int) *season.SeasonResp {
			return item.ConvertToDto()
		}),
		NextToPlay: d.NextToPlay.ConvertToDto(),
	}
}

type SearchSeries struct {
	Page      *Page   `json:"page"`
	Search    string  `json:"search"`
	Id        int64   `json:"id"`
	Type      string  `json:"type"`
	FilterIds []int64 `json:"filterIds"`
	Sort      *Sort   `json:"sort"`
}
