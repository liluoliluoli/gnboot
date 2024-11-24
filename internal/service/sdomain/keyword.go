package sdomain

import (
	dto "gnboot/api/keyword"
	"gnboot/internal/repo/model"
)

type Keyword struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (d *Keyword) ConvertFromRepo(m *model.Keyword) *Keyword {
	return &Keyword{
		ID:   m.ID,
		Name: m.Name,
	}
}

func (d *Keyword) ConvertToDto() *dto.KeywordResp {
	return &dto.KeywordResp{
		Id:   d.ID,
		Name: d.Name,
	}
}
