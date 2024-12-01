package sdomain

import (
	seriesdto "github.com/liluoliluoli/gnboot/api/series"
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
	}
}

func (d *Series) ConvertToDto() *seriesdto.SeriesResp {
	return &seriesdto.SeriesResp{
		Id: d.ID,
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
