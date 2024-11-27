package adaptor

import (
	"context"
	"github.com/liluoliluoli/gnboot/api/episode"
	"github.com/liluoliluoli/gnboot/internal/service"
)

type EpisodeProvider struct {
	episode.UnimplementedEpisodeRemoteServiceServer
	episode *service.EpisodeService
}

func NewEpisodeProvider(episode *service.EpisodeService) *EpisodeProvider {
	return &EpisodeProvider{episode: episode}
}

func (s *EpisodeProvider) GetEpisode(ctx context.Context, req *episode.GetEpisodeRequest) (*episode.EpisodeResp, error) {
	res, err := s.episode.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return res.ConvertToDto(), nil
}
