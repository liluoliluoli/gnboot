package sdomain

import (
	"encoding/json"
	dto "github.com/liluoliluoli/gnboot/api/genre"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
)

type Genre struct {
	ID   int64  `json:"id"`
	Name string `json:"name"` // 标题
}

func (d *Genre) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

func (d *Genre) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

func (d *Genre) ConvertFromRepo(m *model.Genre) *Genre {
	return &Genre{
		ID:   m.ID,
		Name: m.Name,
	}
}

func (d *Genre) ConvertToDto() *dto.GenreResp {
	return &dto.GenreResp{
		Id:   int32(d.ID),
		Name: d.Name,
	}
}

type FindGenre struct {
	Page   Page    `json:"page"`
	Search *string `json:"search"`
	Sort   *Sort   `json:"sort"`
}
