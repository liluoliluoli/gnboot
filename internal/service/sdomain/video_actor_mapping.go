package sdomain

import (
	"github.com/samber/lo"
	"gnboot/internal/repo/model"
)

type VideoActorMapping struct {
	ID        int64  `json:"id"`
	VideId    int64  `json:"videId"`
	VideoType string `json:"type"`
	ActorId   int64  `json:"actorId"`
	Character string `json:"character"`
}

func (d *VideoActorMapping) ConvertFromRepo(m *model.VideoActorMapping) *VideoActorMapping {
	return &VideoActorMapping{
		ID:        m.ID,
		VideId:    m.VideoID,
		VideoType: m.VideoType,
		ActorId:   m.ActorID,
		Character: lo.FromPtr(m.Character),
	}
}