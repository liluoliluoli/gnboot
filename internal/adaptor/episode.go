package adaptor

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/liluoliluoli/gnboot/api/episode"
	"github.com/liluoliluoli/gnboot/internal/common/constant"
	"github.com/liluoliluoli/gnboot/internal/common/gerror"
	"github.com/liluoliluoli/gnboot/internal/common/utils/json_util"
	"github.com/liluoliluoli/gnboot/internal/service"
	"github.com/liluoliluoli/gnboot/internal/task/xiaoya/video"
	"google.golang.org/protobuf/types/known/emptypb"
)

type EpisodeProvider struct {
	episode.UnimplementedEpisodeRemoteServiceServer
	episode     *service.EpisodeService
	user        *service.UserService
	jfVideoTask *video.JfVideoTask
}

func NewEpisodeProvider(episode *service.EpisodeService, user *service.UserService, jfVideoTask *video.JfVideoTask) *EpisodeProvider {
	return &EpisodeProvider{
		episode:     episode,
		user:        user,
		jfVideoTask: jfVideoTask,
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

func (s *EpisodeProvider) UpdateConfigs(ctx context.Context, req *episode.UpdateConfigRequest) (*emptypb.Empty, error) {
	marshalString, err := json_util.MarshalString(req)
	if err != nil {
		return nil, err
	}
	err = s.episode.UpdateConfigs(ctx, marshalString)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *EpisodeProvider) TestSyncTask(ctx context.Context, req *episode.TestSyncTaskRequest) (*emptypb.Empty, error) {
	go func() {
		err := s.jfVideoTask.DoProcess(ctx)
		if err != nil {
			log.Errorf("TestSyncTask fail:%v", err)
		}
	}()
	return &emptypb.Empty{}, nil
}
