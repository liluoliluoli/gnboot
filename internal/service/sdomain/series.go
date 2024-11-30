package sdomain

import (
	seriesdto "github.com/liluoliluoli/gnboot/api/series"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
)

type Series struct {
	ID                 int64                   `json:"id"`
	OriginalTitle      string                  `json:"originalTitle"`      // 标题
	Status             string                  `json:"status"`             // 状态，Returning Series, Ended, Released, Unknown
	VoteAverage        float32                 `json:"voteAverage"`        // 平均评分
	VoteCount          int64                   `json:"voteCount"`          // 评分数
	Country            string                  `json:"country"`            // 国家
	Trailer            string                  `json:"trailer"`            // 预告片地址
	URL                string                  `json:"url"`                // 影片地址
	Downloaded         bool                    `json:"downloaded"`         // 是否可以下载
	FileSize           int64                   `json:"fileSize"`           // 文件大小
	Filename           string                  `json:"filename"`           // 文件名
	Ext                string                  `json:"ext"`                //扩展参数
	Genres             []*Genre                `json:"genres"`             //流派
	Studios            []*Studio               `json:"studios"`            //出品方
	Keywords           []*Keyword              `json:"keywords"`           //关键词
	LastPlayedPosition int64                   `json:"lastPlayedPosition"` //上次播放位置
	LastPlayedTime     string                  `json:"lastPlayedTime"`     //YYYY-MM-DD HH:MM:SS
	Subtitles          []*VideoSubtitleMapping `json:"subtitles"`          //字幕
	Actors             []*Actor                `json:"actors"`             //演员
}

func (d *Series) ConvertFromRepo(movie *model.Series) *Series {
	return &Series{
		ID:            movie.ID,
		OriginalTitle: movie.OriginalTitle,
	}
}

func (d *Series) ConvertToDto() *seriesdto.SeriesResp {
	return &seriesdto.SeriesResp{
		Id:            d.ID,
		OriginalTitle: d.OriginalTitle,
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
