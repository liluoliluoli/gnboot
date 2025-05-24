package adaptor

import (
	"context"
	"github.com/liluoliluoli/gnboot/api/episode"
	"github.com/liluoliluoli/gnboot/internal/common/constant"
	"github.com/liluoliluoli/gnboot/internal/common/gerror"
	"github.com/liluoliluoli/gnboot/internal/service"
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
	res, err := s.episode.Get(ctx, int64(req.Id))
	if err != nil {
		return nil, err
	}
	////从盒子获取播放地址
	//url, err := httpclient_util.DoPost[string, string](ctx, "", nil)
	//if err != nil {
	//	return nil, err
	//}
	//res.Url = lo.FromPtr(url)
	return res.ConvertToDto(), nil
}
