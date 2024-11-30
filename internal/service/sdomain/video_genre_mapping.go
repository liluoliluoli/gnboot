package sdomain

import (
	"github.com/liluoliluoli/gnboot/internal/repo/model"
)

type VideoGenreMapping struct {
	ID        int64  `json:"id"`
	VideoId   int64  `json:"videoId"`
	VideoType string `json:"type"`
	GenreId   int64  `json:"genreId"`
}

func (d *VideoGenreMapping) ConvertFromRepo(m *model.VideoGenreMapping) *VideoGenreMapping {
	return &VideoGenreMapping{
		ID:        m.ID,
		VideoId:   m.VideoID,
		VideoType: m.VideoType,
		GenreId:   m.GenreID,
	}
}
