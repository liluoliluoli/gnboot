package sdomain

import (
	"encoding/json"
	videodto "github.com/liluoliluoli/gnboot/api/video"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/samber/lo"
	"time"
)

type CreateVideo struct {
	Title        string
	VideoType    string
	VoteRate     float32
	VoteCount    int64
	Region       string
	TotalEpisode int32
	Description  string
	PublishMonth string
	Thumbnail    string
	Genres       string
	CreateTime   time.Time
	UpdateTime   time.Time
	IsValid      bool
	Ext          string
}

func (d *CreateVideo) ConvertToRepo() *model.Video {
	return &model.Video{
		Title:        d.Title,
		VideoType:    d.VideoType,
		VoteRate:     lo.ToPtr(d.VoteRate),
		VoteCount:    lo.ToPtr(d.VoteCount),
		Region:       lo.ToPtr(d.Region),
		TotalEpisode: lo.ToPtr(d.TotalEpisode),
		Description:  lo.ToPtr(d.Description),
		Ext:          lo.ToPtr(d.Ext),
		PublishMonth: lo.ToPtr(d.PublishMonth),
		Thumbnail:    lo.ToPtr(d.Thumbnail),
		Genres:       lo.ToPtr(d.Genres),
		CreateTime:   d.CreateTime,
		UpdateTime:   d.CreateTime,
		IsValid:      d.IsValid,
	}
}

func (d *CreateVideo) ConvertFromDto(req *videodto.CreateVideoRequest) *CreateVideo {
	return &CreateVideo{}
}

type Video struct {
	ID                  int64
	Title               string
	VideoType           string
	VoteRate            float32
	VoteCount           int64
	Region              string
	TotalEpisode        int32
	Description         string
	PublishMonth        string
	Thumbnail           string
	Genres              string
	Actors              []*Actor
	Episodes            []*Episode
	LastPlayedTime      *time.Time
	LastPlayedEpisodeId int64
	LastPlayedPosition  int64
	CreateTime          time.Time
	UpdateTime          time.Time
	IsValid             bool
	Ext                 string
}

func (d *Video) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

func (d *Video) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

func (d *Video) ConvertFromRepo(m *model.Video) *Video {
	return &Video{
		ID:           m.ID,
		Title:        d.Title,
		VideoType:    d.VideoType,
		VoteRate:     lo.FromPtr(m.VoteRate),
		VoteCount:    lo.FromPtr(m.VoteCount),
		Region:       lo.FromPtr(m.Region),
		TotalEpisode: lo.FromPtr(m.TotalEpisode),
		Description:  lo.FromPtr(m.Description),
		Ext:          lo.FromPtr(m.Ext),
		PublishMonth: lo.FromPtr(m.PublishMonth),
		Thumbnail:    lo.FromPtr(m.Thumbnail),
		Genres:       lo.FromPtr(m.Genres),
		CreateTime:   m.CreateTime,
		UpdateTime:   m.UpdateTime,
		IsValid:      m.IsValid,
	}
}

func (d *Video) ConvertToDto() *videodto.Video {
	return &videodto.Video{
		Id:  int32(d.ID),
		Ext: d.Ext,
	}
}

type VideoSearch struct {
	Page   *Page  `json:"page"`
	Search string `json:"search"`
	Id     int64  `json:"id"`
	Type   string `json:"type"`
	Sort   string `json:"sort"`
	Genre  string `json:"genre"`
	Region string `json:"region"`
	Year   string `json:"year"`
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
