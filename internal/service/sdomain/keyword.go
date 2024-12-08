package sdomain

import (
	"encoding/json"
	dto "github.com/liluoliluoli/gnboot/api/keyword"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
)

type Keyword struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (d *Keyword) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

func (d *Keyword) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

func (d *Keyword) ConvertFromRepo(m *model.Keyword) *Keyword {
	return &Keyword{
		ID:   m.ID,
		Name: m.Name,
	}
}

func (d *Keyword) ConvertToDto() *dto.KeywordResp {
	return &dto.KeywordResp{
		Id:   int32(d.ID),
		Name: d.Name,
	}
}
