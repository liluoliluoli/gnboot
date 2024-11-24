package sdomain

import (
	"github.com/go-cinch/common/page"
	dto "gnboot/api/genre"
	"gnboot/internal/repo/model"
)

type Genre struct {
	ID   int64  `json:"id"`
	Name string `json:"name"` // 标题
}

func (d *Genre) ConvertFromRepo(m *model.Genre) *Genre {
	return &Genre{
		ID:   m.ID,
		Name: m.Name,
	}
}

func (d *Genre) ConvertToDto() *dto.GenreResp {
	return &dto.GenreResp{
		Id:   d.ID,
		Name: d.Name,
	}
}

type FindGenre struct {
	Page   page.Page `json:"page"`
	Search *string   `json:"search"`
	Sort   *Sort     `json:"sort"`
}
