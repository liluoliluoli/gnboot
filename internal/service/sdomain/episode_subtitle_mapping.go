package sdomain

import (
	"encoding/json"
	dto "github.com/liluoliluoli/gnboot/api/subtitle"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
)

type EpisodeSubtitleMapping struct {
	ID        int64  `json:"id"`
	EpisodeId int64  `json:"episodeId"`
	Url       string `json:"url"`
	Title     string `json:"title"`
	Language  string `json:"language"`
	MimeType  string `json:"mimeType"`
}

func (d *EpisodeSubtitleMapping) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

func (d *EpisodeSubtitleMapping) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

func (d *EpisodeSubtitleMapping) ConvertFromRepo(m *model.EpisodeSubtitleMapping) *EpisodeSubtitleMapping {
	return &EpisodeSubtitleMapping{
		ID:        m.ID,
		EpisodeId: m.EpisodeID,
		Url:       m.URL,
		Title:     m.Title,
		Language:  m.Language,
		MimeType:  m.MimeType,
	}
}

func (d *EpisodeSubtitleMapping) ConvertToDto() *dto.Subtitle {
	return &dto.Subtitle{
		Id:       int32(d.ID),
		Url:      d.Url,
		Title:    d.Title,
		Language: d.Language,
		MimeType: d.MimeType,
	}
}
