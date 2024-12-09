package sdomain

import (
	actordto "github.com/liluoliluoli/gnboot/api/actor"
	episodedto "github.com/liluoliluoli/gnboot/api/episode"
	subtitledto "github.com/liluoliluoli/gnboot/api/subtitle"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Episode struct {
	ID                 int64                   `json:"id"`
	SeasonId           int64                   `json:"seasonId"`
	Season             int32                   `json:"season"`
	SeriesTitle        string                  `json:"seriesTitle"`
	SeasonTitle        string                  `json:"seasonTitle"`
	Episode            int32                   `json:"episode"`
	Url                string                  `json:"url"`
	Downloaded         bool                    `json:"downloaded"`
	Ext                string                  `json:"ext"`
	FileSize           int32                   `json:"fileSize"`
	Subtitles          []*VideoSubtitleMapping `json:"subtitles"`
	Filename           string                  `json:"filename"`           // 文件名
	LastPlayedPosition int32                   `json:"lastPlayedPosition"` //上次播放位置
	LastPlayedTime     *time.Time              `json:"lastPlayedTime"`     //YYYY-MM-DD HH:MM:SS
	Title              string                  `json:"title"`              //标题
	Poster             string                  `json:"poster"`             //海报
	Logo               string                  `json:"logo"`               //logo
	AirDate            *time.Time              `json:"airDate"`            //发行日期
	Overview           string                  `json:"overview"`           //简介
	Favorite           bool                    `json:"favorite"`           //是否喜欢
	Actors             []*Actor                `json:"actors"`
}

func (d *Episode) ConvertFromRepo(m *model.Episode) *Episode {
	return &Episode{
		ID:         m.ID,
		SeasonId:   m.SeasonID,
		Episode:    m.Episode,
		Url:        m.URL,
		Downloaded: m.Downloaded,
		Ext:        lo.FromPtr(m.Ext),
		FileSize:   lo.FromPtr(m.FileSize),
		Filename:   m.Filename,
		Title:      lo.FromPtr(m.Title),
		Poster:     lo.FromPtr(m.Poster),
		Logo:       lo.FromPtr(m.Logo),
		AirDate:    m.AirDate,
		Overview:   lo.FromPtr(m.Overview),
	}
}

func (d *Episode) ConvertToDto() *episodedto.EpisodeResp {
	return &episodedto.EpisodeResp{
		Id:       int32(d.ID),
		Episode:  d.Episode,
		Url:      d.Url,
		Download: d.Downloaded,
		Ext:      d.Ext,
		FileSize: d.FileSize,
		Filename: d.Filename,
		Subtitles: lo.Map(d.Subtitles, func(item *VideoSubtitleMapping, index int) *subtitledto.SubtitleResp {
			return item.ConvertToDto()
		}),
		LastPlayedPosition: d.LastPlayedPosition,
		LastPlayedTime:     lo.Ternary(d.LastPlayedTime != nil, timestamppb.New(lo.FromPtr(d.LastPlayedTime)), nil),
		Title:              d.Title,
		Poster:             d.Poster,
		Logo:               d.Logo,
		AirDate:            lo.Ternary(d.AirDate != nil, timestamppb.New(lo.FromPtr(d.AirDate)), nil),
		Overview:           d.Overview,
		Favorite:           d.Favorite,
		Season:             d.Season,
		SeasonTitle:        d.SeasonTitle,
		SeasonId:           int32(d.SeasonId),
		SeriesTitle:        d.SeriesTitle,
		Actors: lo.Map(d.Actors, func(item *Actor, index int) *actordto.ActorResp {
			return item.ConvertToDto()
		}),
	}
}
