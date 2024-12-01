package sdomain

import (
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
		Id:     d.ID,
		Season: d.Season,
	}
}
