package sdomain

import (
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/samber/lo"
)

type VideoActorMapping struct {
	ID        int64  `json:"id"`
	VideoId   int64  `json:"videoId"`
	VideoType string `json:"type"`
	ActorId   int64  `json:"actorId"`
	Character string `json:"character"`
}

func (d *VideoActorMapping) ConvertFromRepo(m *model.VideoActorMapping) *VideoActorMapping {
	return &VideoActorMapping{
		ID:        m.ID,
		VideoId:   m.VideoID,
		ActorId:   m.ActorID,
		Character: lo.FromPtr(m.Character),
	}
}
