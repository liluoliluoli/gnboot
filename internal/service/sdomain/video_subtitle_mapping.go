package sdomain

import (
	"encoding/json"
	dto "github.com/liluoliluoli/gnboot/api/subtitle"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
)

type VideoSubtitleMapping struct {
	ID        int64  `json:"id"`
	VideoId   int64  `json:"videoId"`
	VideoType string `json:"type"`
	Url       string `json:"url"`
	Title     string `json:"title"`
	Language  string `json:"language"`
	MimeType  string `json:"mimeType"`
}

func (d *VideoSubtitleMapping) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

func (d *VideoSubtitleMapping) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

func (d *VideoSubtitleMapping) ConvertFromRepo(m *model.VideoSubtitleMapping) *VideoSubtitleMapping {
	return &VideoSubtitleMapping{
		ID:        m.ID,
		VideoId:   m.VideoID,
		VideoType: m.VideoType,
		Url:       m.URL,
		Title:     m.Title,
		Language:  m.Language,
		MimeType:  m.MimeType,
	}
}

func (d *VideoSubtitleMapping) ConvertToDto() *dto.SubtitleResp {
	return &dto.SubtitleResp{
		Id:       int32(d.ID),
		Url:      d.Url,
		Title:    d.Title,
		Language: d.Language,
		MimeType: d.MimeType,
	}
}
