package sdomain

import (
	episodedto "github.com/liluoliluoli/gnboot/api/episode"
	seasondto "github.com/liluoliluoli/gnboot/api/season"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Season struct {
	ID           int64      `json:"id"`
	SeriesId     int64      `json:"seriesId"`
	Season       int32      `json:"season"`
	SeriesTitle  string     `json:"seriesTitle"`
	SkipIntro    int32      `json:"skipIntro"`
	SkipEnding   int32      `json:"skipEnding"`
	EpisodeCount int32      `json:"episodeCount"`
	Episodes     []*Episode `json:"episodes"`
	Title        string     `json:"title"`
	Poster       string     `json:"poster"`
	Logo         string     `json:"logo"`
	Favorite     bool       `json:"favorite"`
	AirDate      *time.Time `json:"airDate"`
	Overview     string     `json:"overview"`
}

func (d *Season) ConvertFromRepo(m *model.Season) *Season {
	return &Season{
		ID:           m.ID,
		SeriesId:     m.SeriesID,
		Season:       m.Season,
		SkipIntro:    lo.FromPtr(m.SkipIntro),
		SkipEnding:   lo.FromPtr(m.SkipEnding),
		EpisodeCount: m.EpisodeCount,
		Title:        lo.FromPtr(m.Title),
		Poster:       lo.FromPtr(m.Poster),
		Logo:         lo.FromPtr(m.Logo),
		AirDate:      m.AirDate,
		Overview:     lo.FromPtr(m.Overview),
	}
}

func (d *Season) ConvertToDto() *seasondto.SeasonResp {
	return &seasondto.SeasonResp{
		Id:           int32(d.ID),
		Season:       d.Season,
		SkipIntro:    d.SkipIntro,
		SkipEnding:   d.SkipEnding,
		EpisodeCount: d.EpisodeCount,
		Episodes: lo.Map(d.Episodes, func(item *Episode, index int) *episodedto.EpisodeResp {
			return item.ConvertToDto()
		}),
		Title:    d.Title,
		Poster:   d.Poster,
		Logo:     d.Logo,
		AirDate:  lo.Ternary(d.AirDate != nil, timestamppb.New(lo.FromPtr(d.AirDate)), nil),
		Overview: d.Overview,
		SeriesId: int32(d.SeriesId),
	}
}
