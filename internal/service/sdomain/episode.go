package sdomain

import (
	episodedto "github.com/liluoliluoli/gnboot/api/episode"
	subtitledto "github.com/liluoliluoli/gnboot/api/subtitle"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/samber/lo"
)

type Episode struct {
	ID           int64                     `json:"id"`
	VideoId      int64                     `json:"videoId"`
	Episode      int32                     `json:"episode"`
	EpisodeTitle string                    `json:"episodeTitle"`
	Url          string                    `json:"url"`
	Platform     string                    `json:"platform"`
	Duration     int64                     `json:"duration"`
	Ext          string                    `json:"ext"`
	Subtitles    []*EpisodeSubtitleMapping `json:"subtitles"`
}

func (d *Episode) ConvertFromRepo(m *model.Episode) *Episode {
	return &Episode{
		ID:           m.ID,
		VideoId:      m.VideoID,
		Episode:      m.Episode,
		EpisodeTitle: m.EpisodeTitle,
		Url:          lo.FromPtr(m.URL),
		Platform:     lo.FromPtr(m.Platform),
		Duration:     lo.FromPtr(m.Duration),
		Ext:          lo.FromPtr(m.Ext),
	}
}

func (d *Episode) ConvertToDto() *episodedto.Episode {
	return &episodedto.Episode{
		Id:           int32(d.ID),
		VideoId:      int32(d.VideoId),
		Episode:      d.Episode,
		EpisodeTitle: d.EpisodeTitle,
		Url:          d.Url,
		Platform:     d.Platform,
		Ext:          d.Ext,
		Duration:     int32(d.Duration),
		Subtitles: lo.Map(d.Subtitles, func(item *EpisodeSubtitleMapping, index int) *subtitledto.Subtitle {
			return item.ConvertToDto()
		}),
	}
}
