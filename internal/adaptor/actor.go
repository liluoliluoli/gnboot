package adaptor

import (
	"context"
	"github.com/liluoliluoli/gnboot/api/actor"
	"github.com/liluoliluoli/gnboot/internal/service"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/samber/lo"
)

type ActorProvider struct {
	actor.UnimplementedActorRemoteServiceServer
	actor *service.ActorService
}

func NewActorProvider(actor *service.ActorService) *ActorProvider {
	return &ActorProvider{actor: actor}
}

func (s *ActorProvider) FindActor(ctx context.Context, req *actor.FindActorRequest) (*actor.FindActorResp, error) {
	res, err := s.actor.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return &actor.FindActorResp{
		Actors: lo.Map(res, func(item *sdomain.Actor, index int) *actor.ActorResp {
			return &actor.ActorResp{
				Name: item.Name,
				Id:   int32(item.ID),
			}
		}),
	}, nil

}
