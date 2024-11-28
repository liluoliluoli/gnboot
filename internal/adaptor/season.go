package adaptor

import (
	"context"
	"github.com/liluoliluoli/gnboot/api/season"
	"github.com/liluoliluoli/gnboot/internal/service"
)

type SeasonProvider struct {
	season.UnimplementedSeasonRemoteServiceServer
	season *service.SeasonService
}

func NewSeasonProvider(season *service.SeasonService) *SeasonProvider {
	return &SeasonProvider{season: season}
}

func (s *SeasonProvider) GetSeason(ctx context.Context, req *season.GetSeasonRequest) (*season.SeasonResp, error) {
	res, err := s.season.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return res.ConvertToDto(), nil
}
