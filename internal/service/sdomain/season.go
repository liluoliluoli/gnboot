package sdomain

import (
	episodedto "github.com/liluoliluoli/gnboot/api/episode"
	seasondto "github.com/liluoliluoli/gnboot/api/season"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/samber/lo"
)

type Season struct {
	ID           int64      `json:"id"`
	SeriesId     int64      `json:"seriesId"`
	Season       int32      `json:"season"`
	SeasonTitle  string     `json:"seasonTitle"`
	SkipIntro    int32      `json:"skipIntro"`
	SkipEnding   int32      `json:"skipEnding"`
	EpisodeCount int32      `json:"episodeCount"`
	Episodes     []*Episode `json:"episodes"`
}

func (d *Season) ConvertFromRepo(m *model.Season) *Season {
	return &Season{
		ID:           m.ID,
		SeriesId:     m.SeriesID,
		Season:       m.Season,
		SkipIntro:    lo.FromPtr(m.SkipIntro),
		SkipEnding:   lo.FromPtr(m.SkipEnding),
		EpisodeCount: m.EpisodeCount,
	}
}

func (d *Season) ConvertToDto() *seasondto.SeasonResp {
	return &seasondto.SeasonResp{
		Id:           d.ID,
		Season:       d.Season,
		SkipIntro:    d.SkipIntro,
		SkipEnding:   d.SkipEnding,
		EpisodeCount: d.EpisodeCount,
		Episodes: lo.Map(d.Episodes, func(item *Episode, index int) *episodedto.EpisodeResp {
			return item.ConvertToDto()
		}),
	}
}
