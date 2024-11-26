package sdomain

import (
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/samber/lo"
)

type VideoGenreMapping struct {
	ID        int64  `json:"id"`
	VideId    int64  `json:"videId"`
	VideoType string `json:"type"`
	GenreId   int64  `json:"genreId"`
}

func (d *VideoGenreMapping) ConvertFromRepo(m *model.VideoGenreMapping) *VideoGenreMapping {
	return &VideoGenreMapping{
		ID:        m.ID,
		VideId:    lo.FromPtr(m.VideoID),
		VideoType: lo.FromPtr(m.VideoType),
		GenreId:   m.GenreID,
	}
}
