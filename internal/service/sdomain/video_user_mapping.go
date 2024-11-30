package sdomain

import (
	"encoding/json"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/samber/lo"
	"time"
)

type VideoUserMapping struct {
	ID                 int64     `json:"id"`
	VideoId            int64     `json:"videoId"`
	VideoType          string    `json:"type"`
	LastPlayedPosition int32     `json:"lastPlayedPosition"`
	LastPlayedTime     time.Time `json:"lastPlayedTime"`
	Favorited          bool      `json:"favorited"`
	UserId             int64     `json:"userId"`
}

func (d *VideoUserMapping) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

func (d *VideoUserMapping) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

func (d *VideoUserMapping) ConvertFromRepo(m *model.VideoUserMapping) *VideoUserMapping {
	return &VideoUserMapping{
		ID:                 m.ID,
		VideoId:            m.VideoID,
		VideoType:          m.VideoType,
		LastPlayedPosition: lo.FromPtr(m.LastPlayedPosition),
		LastPlayedTime:     lo.FromPtr(m.LastPlayedTime),
		Favorited:          lo.FromPtr(m.Favorited),
		UserId:             m.UserID,
	}
}
