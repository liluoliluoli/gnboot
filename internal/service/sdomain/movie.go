package sdomain

import (
	"encoding/json"
	"github.com/liluoliluoli/gnboot/api/actor"
	"github.com/liluoliluoli/gnboot/api/genre"
	"github.com/liluoliluoli/gnboot/api/keyword"
	moviedto "github.com/liluoliluoli/gnboot/api/movie"
	"github.com/liluoliluoli/gnboot/api/studio"
	"github.com/liluoliluoli/gnboot/api/subtitle"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/samber/lo"
	"time"
)

type CreateMovie struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (d *CreateMovie) ConvertToRepo() *model.Movie {
	return &model.Movie{}
}

func (d *CreateMovie) ConvertFromDto(req *moviedto.CreateMovieRequest) *CreateMovie {
	return &CreateMovie{}
}

type Movie struct {
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
	LastPlayedPosition int32                   `json:"lastPlayedPosition"` //上次播放位置
	LastPlayedTime     *time.Time              `json:"lastPlayedTime"`     //YYYY-MM-DD HH:MM:SS
	Subtitles          []*VideoSubtitleMapping `json:"subtitles"`          //字幕
	Actors             []*Actor                `json:"actors"`             //演员
}

func (d *Movie) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

func (d *Movie) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

func (d *Movie) ConvertFromRepo(movie *model.Movie) *Movie {
	return &Movie{
		ID:            movie.ID,
		OriginalTitle: movie.OriginalTitle,
		Status:        movie.Status,
		VoteAverage:   lo.FromPtr(movie.VoteAverage),
		VoteCount:     lo.FromPtr(movie.VoteCount),
		Country:       lo.FromPtr(movie.Country),
		Trailer:       lo.FromPtr(movie.Trailer),
		URL:           movie.URL,
		Downloaded:    movie.Downloaded,
		FileSize:      lo.FromPtr(movie.FileSize),
		Filename:      lo.FromPtr(movie.Filename),
		Ext:           lo.FromPtr(movie.Ext),
	}
}

func (d *Movie) ConvertToDto() *moviedto.MovieResp {
	return &moviedto.MovieResp{
		Id:            d.ID,
		OriginalTitle: d.OriginalTitle,
		Status:        d.Status,
		VoteAverage:   d.VoteAverage,
		VoteCount:     d.VoteCount,
		Country:       d.Country,
		Trailer:       d.Trailer,
		Url:           d.URL,
		Downloaded:    d.Downloaded,
		FileSize:      d.FileSize,
		Filename:      d.Filename,
		Ext:           d.Ext,
		Genres: lo.Map(d.Genres, func(item *Genre, index int) *genre.GenreResp {
			return item.ConvertToDto()
		}),
		Studios: lo.Map(d.Studios, func(item *Studio, index int) *studio.StudioResp {
			return item.ConvertToDto()
		}),
		Keywords: lo.Map(d.Keywords, func(item *Keyword, index int) *keyword.KeywordResp {
			return item.ConvertToDto()
		}),
		Actors: lo.Map(d.Actors, func(item *Actor, index int) *actor.ActorResp {
			return item.ConvertToDto()
		}),
		Subtitles: lo.Map(d.Subtitles, func(item *VideoSubtitleMapping, index int) *subtitle.SubtitleResp {
			return item.ConvertToDto()
		}),
	}
}

type SearchMovie struct {
	Page      *Page   `json:"page"`
	Search    string  `json:"search"`
	Id        int64   `json:"id"`
	Type      string  `json:"type"`
	FilterIds []int64 `json:"filterIds"`
	Sort      *Sort   `json:"sort"`
}

type UpdateMovie struct {
	ID    int64  `json:"id,string"`
	Title string `json:"title,omitempty"`
}

func (d *UpdateMovie) ConvertToRepo() *model.Movie {
	return &model.Movie{}
}

func (d *UpdateMovie) ConvertFromDto(dto *moviedto.UpdateMovieRequest) *UpdateMovie {
	return &UpdateMovie{}
}
