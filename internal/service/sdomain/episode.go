package sdomain

import (
	episodedto "github.com/liluoliluoli/gnboot/api/episode"
	subtitledto "github.com/liluoliluoli/gnboot/api/subtitle"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/samber/lo"
	"time"
)

type Episode struct {
	ID           int64                     `json:"id"`
	VideoId      int64                     `json:"videoId"`
	Episode      int32                     `json:"episode"`
	EpisodeTitle string                    `json:"episodeTitle"`
	Url          string                    `json:"url"`
	Platform     string                    `json:"platform"`
	Duration     int64                     `json:"duration"`
	Size         string                    `json:"size"`
	Ext          string                    `json:"ext"`
	Subtitles    []*EpisodeSubtitleMapping `json:"subtitles"`
	XiaoYaPath   string                    `json:"xiaoYaPath"`
	ExpiredTime  *time.Time                `json:"expiredTime"`
	CreateTime   time.Time                 `json:"createTime"`
	UpdateTime   time.Time                 `json:"updateTime"`
	Ratio        string                    `json:"ratio"`
	JellyfinId   string                    `json:"jellyfinId"`
}

func (d *Episode) ConvertFromRepo(m *model.Episode) *Episode {
	if m == nil {
		return nil
	}
	return &Episode{
		ID:           m.ID,
		VideoId:      m.VideoID,
		Episode:      m.Episode,
		EpisodeTitle: m.EpisodeTitle,
		Url:          lo.FromPtr(m.URL),
		Platform:     lo.FromPtr(m.Platform),
		Duration:     lo.FromPtr(m.Duration),
		Ext:          lo.FromPtr(m.Ext),
		XiaoYaPath:   lo.FromPtr(m.XiaoyaPath),
		ExpiredTime:  m.ExpiredTime,
		Size:         m.Size,
		CreateTime:   m.CreateTime,
		UpdateTime:   m.UpdateTime,
		Ratio:        lo.FromPtr(m.Ratio),
		JellyfinId:   lo.FromPtr(m.JellyfinID),
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
		Ratio: d.Ratio,
	}
}

func (d *Episode) ConvertToRepo() *model.Episode {
	return &model.Episode{
		ID:           d.ID,
		VideoID:      d.VideoId,
		Episode:      d.Episode,
		EpisodeTitle: d.EpisodeTitle,
		URL:          lo.ToPtr(d.Url),
		Platform:     lo.ToPtr(d.Platform),
		Ext:          lo.ToPtr(d.Ext),
		Duration:     lo.ToPtr(d.Duration),
		Size:         d.Size,
		ExpiredTime:  d.ExpiredTime,
		XiaoyaPath:   lo.ToPtr(d.XiaoYaPath),
		CreateTime:   d.CreateTime,
		UpdateTime:   time.Now(),
		Ratio:        lo.ToPtr(d.Ratio),
		JellyfinID:   lo.ToPtr(d.JellyfinId),
	}
}
