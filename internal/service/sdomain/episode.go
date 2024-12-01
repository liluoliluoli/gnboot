package sdomain

import (
	episodedto "github.com/liluoliluoli/gnboot/api/episode"
	subtitledto "github.com/liluoliluoli/gnboot/api/subtitle"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/samber/lo"
)

type Episode struct {
	ID         int64                   `json:"id"`
	SeasonId   int64                   `json:"seasonId"`
	Episode    int32                   `json:"episode"`
	Url        string                  `json:"url"`
	Downloaded bool                    `json:"downloaded"`
	Ext        string                  `json:"ext"`
	FileSize   int32                   `json:"fileSize"`
	Subtitles  []*VideoSubtitleMapping `json:"subtitles"`
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
	}
}

func (d *Episode) ConvertToDto() *episodedto.EpisodeResp {
	return &episodedto.EpisodeResp{
		Id:       d.ID,
		Episode:  d.Episode,
		Url:      d.Url,
		Download: d.Downloaded,
		Ext:      d.Ext,
		FileSize: d.FileSize,
		Subtitles: lo.Map(d.Subtitles, func(item *VideoSubtitleMapping, index int) *subtitledto.SubtitleResp {
			return item.ConvertToDto()
		}),
	}
}
