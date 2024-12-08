package adaptor

import (
	"context"
	"github.com/liluoliluoli/gnboot/api/studio"
	"github.com/liluoliluoli/gnboot/internal/service"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/samber/lo"
)

type StudioProvider struct {
	studio.UnimplementedStudioRemoteServiceServer
	studio *service.StudioService
}

func NewStudioProvider(studio *service.StudioService) *StudioProvider {
	return &StudioProvider{studio: studio}
}

func (s *StudioProvider) FindStudio(ctx context.Context, req *studio.FindStudioRequest) (*studio.FindStudioResp, error) {
	res, err := s.studio.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return &studio.FindStudioResp{
		Studios: lo.Map(res, func(item *sdomain.Studio, index int) *studio.StudioResp {
			return &studio.StudioResp{
				Name:    item.Name,
				Id:      int32(item.ID),
				Country: item.Country,
				Logo:    item.Logo,
			}
		}),
	}, nil

}
