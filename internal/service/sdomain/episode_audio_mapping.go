package sdomain

import (
	"encoding/json"
	dto "github.com/liluoliluoli/gnboot/api/episode"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
)

type EpisodeAudioMapping struct {
	ID        int64  `json:"id"`
	EpisodeId int64  `json:"episodeId"`
	Url       string `json:"url"`
	Title     string `json:"title"`
	Language  string `json:"language"`
	MimeType  string `json:"mimeType"`
}

func (d *EpisodeAudioMapping) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

func (d *EpisodeAudioMapping) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

func (d *EpisodeAudioMapping) ConvertFromRepo(m *model.EpisodeAudioMapping) *EpisodeAudioMapping {
	return &EpisodeAudioMapping{
		ID:        m.ID,
		EpisodeId: m.EpisodeID,
		Url:       m.URL,
		Title:     m.Title,
		Language:  m.Language,
		MimeType:  m.MimeType,
	}
}

func (d *EpisodeAudioMapping) ConvertToDto() *dto.Audio {
	return &dto.Audio{
		Id:       int32(d.ID),
		Url:      d.Url,
		Title:    d.Title,
		Language: d.Language,
		MimeType: d.MimeType,
	}
}
