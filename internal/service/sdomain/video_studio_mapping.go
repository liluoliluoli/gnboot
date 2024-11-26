package sdomain

import (
	"github.com/liluoliluoli/gnboot/internal/repo/model"
)

type VideoStudioMapping struct {
	ID        int64  `json:"id"`
	VideId    int64  `json:"videId"`
	VideoType string `json:"type"`
	StudioId  int64  `json:"studioId"`
}

func (d *VideoStudioMapping) ConvertFromRepo(m *model.VideoStudioMapping) *VideoStudioMapping {
	return &VideoStudioMapping{
		ID:        m.ID,
		VideId:    m.VideoID,
		VideoType: m.VideoType,
		StudioId:  m.StudioID,
	}
}
