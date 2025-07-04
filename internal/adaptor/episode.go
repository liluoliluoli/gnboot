package adaptor

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/liluoliluoli/gnboot/api/episode"
	"github.com/liluoliluoli/gnboot/internal/common/constant"
	"github.com/liluoliluoli/gnboot/internal/common/gerror"
	"github.com/liluoliluoli/gnboot/internal/common/utils/json_util"
	"github.com/liluoliluoli/gnboot/internal/service"
	"github.com/liluoliluoli/gnboot/internal/task/video"
	"google.golang.org/protobuf/types/known/emptypb"
	"runtime/debug"
)

type EpisodeProvider struct {
	episode.UnimplementedEpisodeRemoteServiceServer
	episode       *service.EpisodeService
	user          *service.UserService
	jfVideoTask   *video.JfVideoTask
	embyVideoTask *video.EmbyVideoTask
}

func NewEpisodeProvider(episode *service.EpisodeService, user *service.UserService, jfVideoTask *video.JfVideoTask, embyVideoTask *video.EmbyVideoTask) *EpisodeProvider {
	return &EpisodeProvider{
		episode:       episode,
		user:          user,
		jfVideoTask:   jfVideoTask,
		embyVideoTask: embyVideoTask,
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

func (s *EpisodeProvider) TestLatestSyncTask(ctx context.Context, req *episode.TestLatestSyncTaskRequest) (*emptypb.Empty, error) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("TestLatestSyncTask panic recovered: %v\n%s", r, debug.Stack())
			}
		}()

		err := s.embyVideoTask.LatestSync(ctx, req.ScanPathIds)
		if err != nil {
			log.Errorf("TestLatestSyncTask fail:%v", err)
		}
	}()
	return &emptypb.Empty{}, nil
}

func (s *EpisodeProvider) TestFullSyncTask(ctx context.Context, req *episode.TestFullSyncTaskRequest) (*emptypb.Empty, error) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("TestFullSyncTask panic recovered: %v\n%s", r, debug.Stack())
			}
		}()

		err := s.embyVideoTask.FullSync(ctx, req.ScanPathIds)
		if err != nil {
			log.Errorf("TestFullSyncTask fail:%v", err)
		}
	}()
	return &emptypb.Empty{}, nil
}
