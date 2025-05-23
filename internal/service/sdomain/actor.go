package sdomain

import (
	dto "github.com/liluoliluoli/gnboot/api/actor"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/samber/lo"
)

type Actor struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Thumbnail string `json:"thumbnail"`
	Region    string `json:"region"`
}

func (d *Actor) ConvertFromRepo(m *model.Actor) *Actor {
	return &Actor{
		ID:        m.ID,
		Name:      m.Name,
		Thumbnail: lo.FromPtr(m.Thumbnail),
		Region:    lo.FromPtr(m.Region),
	}
}

func (d *Actor) ConvertToDto() *dto.Actor {
	return &dto.Actor{
		Id:        int32(d.ID),
		Name:      d.Name,
		Thumbnail: d.Thumbnail,
		Region:    d.Region,
	}
}
