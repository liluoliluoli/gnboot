package sdomain

import (
	"encoding/json"
	actordto "github.com/liluoliluoli/gnboot/api/actor"
	episodedto "github.com/liluoliluoli/gnboot/api/episode"
	videodto "github.com/liluoliluoli/gnboot/api/video"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/samber/lo"
	"strings"
	"time"
)

type Video struct {
	ID                  int64
	Title               string
	VideoType           string
	VoteRate            float32
	WatchCount          int64
	Region              string
	TotalEpisode        int32
	Description         string
	PublishDay          string
	Thumbnail           string
	Ratio               string
	Genres              []string
	Actors              []*VideoActorMapping
	Episodes            []*Episode
	LastPlayedTime      *time.Time
	LastPlayedEpisodeId int64
	LastPlayedPosition  int64
	IsFavorite          bool
	CreateTime          time.Time
	UpdateTime          time.Time
	IsValid             bool
	Ext                 string
	JellyfinId          string
	JellyfinCreateTime  time.Time
	JellyfinRootPathId  string
}

func (d *Video) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

func (d *Video) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

func (d *Video) ConvertFromRepo(m *model.Video) *Video {
	return &Video{
		ID:                 m.ID,
		Title:              m.Title,
		VideoType:          m.VideoType,
		VoteRate:           lo.FromPtr(m.VoteRate),
		WatchCount:         m.WatchCount,
		Region:             lo.FromPtr(m.Region),
		TotalEpisode:       lo.FromPtr(m.TotalEpisode),
		Description:        lo.FromPtr(m.Description),
		Ext:                lo.FromPtr(m.Ext),
		PublishDay:         lo.FromPtr(m.PublishDay),
		Thumbnail:          lo.FromPtr(m.Thumbnail),
		Ratio:              lo.FromPtr(m.Ratio),
		Genres:             strings.Split(lo.FromPtr(m.Genres), ","),
		CreateTime:         m.CreateTime,
		UpdateTime:         m.UpdateTime,
		IsValid:            m.IsValid,
		JellyfinId:         m.JellyfinID,
		JellyfinCreateTime: m.JellyfinCreateTime,
		JellyfinRootPathId: lo.FromPtr(m.JellyfinRootPathID),
	}
}

func (d *Video) ConvertToRepo() *model.Video {
	return &model.Video{
		Title:              d.Title,
		VideoType:          d.VideoType,
		VoteRate:           lo.ToPtr(d.VoteRate),
		WatchCount:         d.WatchCount,
		Region:             lo.ToPtr(d.Region),
		TotalEpisode:       lo.ToPtr(d.TotalEpisode),
		Description:        lo.ToPtr(d.Description),
		Ext:                lo.ToPtr(d.Ext),
		PublishDay:         lo.ToPtr(d.PublishDay),
		Thumbnail:          lo.ToPtr(d.Thumbnail),
		Genres:             lo.ToPtr(strings.Join(d.Genres, ",")),
		CreateTime:         d.CreateTime,
		UpdateTime:         d.CreateTime,
		IsValid:            d.IsValid,
		JellyfinID:         d.JellyfinId,
		JellyfinCreateTime: d.JellyfinCreateTime,
		Ratio:              lo.ToPtr(d.Ratio),
		JellyfinRootPathID: lo.ToPtr(d.JellyfinRootPathId),
	}
}

func (d *Video) ConvertFromDto(req *videodto.CreateVideoRequest) *Video {
	return &Video{}
}

func (d *Video) ConvertToDto() *videodto.Video {
	actors := lo.Map(d.Actors, func(item *VideoActorMapping, index int) *actordto.Actor {
		return item.ConvertToDto()
	})
	actors = lo.Filter(actors, func(item *actordto.Actor, index int) bool {
		return item != nil
	})
	return &videodto.Video{
		Id:           int32(d.ID),
		Title:        d.Title,
		VideoType:    d.VideoType,
		VoteRate:     d.VoteRate,
		WatchCount:   int32(d.WatchCount),
		Region:       d.Region,
		TotalEpisode: d.TotalEpisode,
		Description:  d.Description,
		Ext:          d.Ext,
		PublishDay:   d.PublishDay,
		Thumbnail:    d.Thumbnail,
		Ratio:        d.Ratio,
		Genres:       d.Genres,
		LastPlayedTime: lo.TernaryF(d.LastPlayedTime != nil, func() *int32 {
			return lo.ToPtr(int32(d.LastPlayedTime.Unix()))
		}, func() *int32 {
			return nil
		}),
		LastPlayedEpisodeId: lo.Ternary(d.LastPlayedEpisodeId == 0, nil, lo.ToPtr(int32(d.LastPlayedEpisodeId))),
		LastPlayedPosition:  lo.Ternary(d.LastPlayedPosition == 0, nil, lo.ToPtr(int32(d.LastPlayedPosition))),
		IsFavorite:          d.IsFavorite,
		Actors: lo.Filter(actors, func(item *actordto.Actor, index int) bool {
			return !item.IsDirector
		}),
		Directors: lo.Filter(actors, func(item *actordto.Actor, index int) bool {
			return item.IsDirector
		}),
		Episodes: lo.Map(d.Episodes, func(item *Episode, index int) *episodedto.Episode {
			return item.ConvertToDto()
		}),
	}
}

type VideoSearch struct {
	Page      *Page   `json:"page"`
	Search    string  `json:"search"`
	Ids       []int64 `json:"ids"`
	Type      string  `json:"type"`
	Sort      string  `json:"sort"`
	Genre     string  `json:"genre"`
	Region    string  `json:"region"`
	Year      string  `json:"year"`
	IsHistory bool    `json:"isHistory"`
}

type UpdateVideo struct {
	ID    int64  `json:"id,string"`
	Title string `json:"title,omitempty"`
}

func (d *UpdateVideo) ConvertToRepo() *model.Video {
	return &model.Video{}
}

func (d *UpdateVideo) ConvertFromDto(dto *videodto.UpdateVideoRequest) *UpdateVideo {
	return &UpdateVideo{}
}
