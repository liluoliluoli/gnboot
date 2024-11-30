package sdomain

import (
	"encoding/json"
	dto "github.com/liluoliluoli/gnboot/api/actor"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/samber/lo"
)

type Actor struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	OriginalName string `json:"originalName"`
	Adult        bool   `json:"adult"`
	Gender       int32  `json:"gender"`
	Profile      string `json:"profile"`
}

func (d *Actor) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

func (d *Actor) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

func (d *Actor) ConvertFromRepo(m *model.Actor) *Actor {
	return &Actor{
		ID:           m.ID,
		Name:         lo.FromPtr(m.Name),
		OriginalName: lo.FromPtr(m.OriginalAme),
		Adult:        m.Adult,
		Gender:       1, // todo
		Profile:      lo.FromPtr(m.Profile),
	}
}

func (d *Actor) ConvertToDto() *dto.ActorResp {
	return &dto.ActorResp{
		Id:           d.ID,
		Name:         d.Name,
		OriginalName: d.OriginalName,
		Adult:        d.Adult,
		Gender:       d.Gender,
		Profile:      d.Profile,
	}
}
