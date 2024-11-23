package sdomain

import (
	"github.com/go-cinch/common/page"
	moviedto "gnboot/api/movie"
	"gnboot/internal/repo/model"
)

type CreateMovie struct {
	ID   int64  `json:"id,string"`
	Name string `json:"name"`
}

func (*CreateMovie) ConvertToRepo() *model.Movie {
	return &model.Movie{}
}

func (*CreateMovie) ConvertFromDto(req *moviedto.CreateMovieRequest) *CreateMovie {
	return &CreateMovie{}
}

type Movie struct {
	ID            int64   `json:"id,string"`
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

func (*Movie) ConvertFromRepo(movie *model.Movie) *Movie {
	return &Movie{}
}

func (*Movie) ConvertFromDto() *moviedto.GetMovieResp {
	return &moviedto.GetMovieResp{}
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
	ID    int64   `json:"id,string"`
	Title *string `json:"title,omitempty"`
}

func (*UpdateMovie) ConvertToRepo() *model.Movie {
	return &model.Movie{}
}

func (*UpdateMovie) ConvertFromDto(dto *moviedto.UpdateMovieRequest) *UpdateMovie {
	return &UpdateMovie{}
}
