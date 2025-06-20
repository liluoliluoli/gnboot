package adaptor

import (
	"context"
	"github.com/liluoliluoli/gnboot/api/episode"
	"github.com/liluoliluoli/gnboot/internal/common/constant"
	"github.com/liluoliluoli/gnboot/internal/common/gerror"
	"github.com/liluoliluoli/gnboot/internal/common/utils/json_util"
	"github.com/liluoliluoli/gnboot/internal/service"
	"google.golang.org/protobuf/types/known/emptypb"
)

type EpisodeProvider struct {
	episode.UnimplementedEpisodeRemoteServiceServer
	episode *service.EpisodeService
	user    *service.UserService
}

func NewEpisodeProvider(episode *service.EpisodeService, user *service.UserService) *EpisodeProvider {
	return &EpisodeProvider{
		episode: episode,
		user:    user,
	}
}

func (s *EpisodeProvider) GetEpisode(ctx context.Context, req *episode.GetEpisodeRequest) (*episode.Episode, error) {
	user, err := s.user.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	if user.PackageType == constant.None {
		return nil, gerror.ErrAccountPackageExpire(ctx)
	}
	res, err := s.episode.Get(ctx, int64(req.Id), true)
	if err != nil {
		return nil, err
	}
	return res.ConvertToDto(), nil
}

func (s *EpisodeProvider) UpdateBoxIps(ctx context.Context, req *episode.UpdateBoxIpsRequest) (*emptypb.Empty, error) {
	marshalString, err := json_util.MarshalString(req.BoxIps)
	if err != nil {
		return nil, err
	}
	err = s.episode.UpdateBoxIps(ctx, marshalString)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
