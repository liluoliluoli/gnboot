package sdomain

import (
	actordto "github.com/liluoliluoli/gnboot/api/actor"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/samber/lo"
)

type VideoActorMapping struct {
	ID         int64  `json:"id"`
	VideoId    int64  `json:"videoId"`
	ActorId    int64  `json:"actorId"`
	Character  string `json:"character"`
	IsDirector bool   `json:"isDirector"`
	Actor      *Actor `json:"actor"`
}

func (d *VideoActorMapping) ConvertFromRepo(m *model.VideoActorMapping) *VideoActorMapping {
	return &VideoActorMapping{
		ID:         m.ID,
		VideoId:    m.VideoID,
		ActorId:    m.ActorID,
		Character:  lo.FromPtr(m.Character),
		IsDirector: m.IsDirector,
	}
}

func (d *VideoActorMapping) ConvertToDto() *actordto.Actor {
	return &actordto.Actor{
		Id:         int32(d.ID),
		Name:       d.Actor.Name,
		Thumbnail:  d.Actor.Thumbnail,
		Region:     d.Actor.Region,
		IsDirector: d.IsDirector,
		Character:  d.Character,
	}
}
