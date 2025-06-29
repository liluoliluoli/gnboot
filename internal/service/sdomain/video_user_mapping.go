package sdomain

import (
	"encoding/json"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/samber/lo"
	"time"
)

type VideoUserMapping struct {
	ID                  int64      `json:"id"`
	VideoId             int64      `json:"videoId"`
	LastPlayedPosition  int64      `json:"lastPlayedPosition"`
	LastPlayedTime      *time.Time `json:"lastPlayedTime"`
	LastPlayedEpisodeId int64      `json:"lastPlayedEpisodeId"`
	IsFavorite          bool       `json:"isFavorite"`
	UserId              int64      `json:"userId"`
	CreateTime          time.Time  `json:"createTime"`
	UpdateTime          time.Time  `json:"updateTime"`
}

func (d *VideoUserMapping) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

func (d *VideoUserMapping) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

func (d *VideoUserMapping) ConvertFromRepo(m *model.VideoUserMapping) *VideoUserMapping {
	return &VideoUserMapping{
		ID:                  m.ID,
		VideoId:             m.VideoID,
		LastPlayedPosition:  lo.FromPtr(m.LastPlayedPosition),
		LastPlayedTime:      m.LastPlayedTime,
		LastPlayedEpisodeId: lo.FromPtr(m.LastPlayedEpisodeID),
		IsFavorite:          m.IsFavorite,
		UserId:              m.UserID,
		CreateTime:          m.CreateTime,
		UpdateTime:          m.UpdateTime,
	}
}
