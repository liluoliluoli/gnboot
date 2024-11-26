package sdomain

import (
	dto "github.com/liluoliluoli/gnboot/api/studio"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/samber/lo"
)

type Studio struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"` // 标题
	Country string `json:"country"`
	Logo    string `json:"logo"`
}

func (d *Studio) ConvertFromRepo(m *model.Studio) *Studio {
	return &Studio{
		ID:      m.ID,
		Name:    m.Name,
		Country: m.Country,
		Logo:    lo.FromPtr(m.Logo),
	}
}

func (d *Studio) ConvertToDto() *dto.StudioResp {
	return &dto.StudioResp{
		Id:      d.ID,
		Name:    d.Name,
		Country: d.Country,
		Logo:    d.Logo,
	}
}
