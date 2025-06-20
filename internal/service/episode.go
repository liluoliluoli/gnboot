package service

import (
	"context"
	"fmt"
	errors "github.com/go-kratos/kratos/v2/errors"
	"github.com/liluoliluoli/gnboot/internal/common/constant"
	"github.com/liluoliluoli/gnboot/internal/common/gerror"
	"github.com/liluoliluoli/gnboot/internal/common/utils/array_util"
	"github.com/liluoliluoli/gnboot/internal/common/utils/cache_util"
	"github.com/liluoliluoli/gnboot/internal/common/utils/context_util"
	"github.com/liluoliluoli/gnboot/internal/common/utils/httpclient_util"
	"github.com/liluoliluoli/gnboot/internal/common/utils/json_util"
	"github.com/liluoliluoli/gnboot/internal/common/utils/time_util"
	"github.com/liluoliluoli/gnboot/internal/conf"
	"github.com/liluoliluoli/gnboot/internal/integration/xiaoyadto"
	"github.com/liluoliluoli/gnboot/internal/repo"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"time"
)

type EpisodeService struct {
	c                        *conf.Bootstrap
	videoRepo                *repo.VideoRepo
	actorRepo                *repo.ActorRepo
	videoActorMappingRepo    *repo.VideoActorMappingRepo
	videoSubtitleMappingRepo *repo.EpisodeSubtitleMappingRepo
	episodeRepo              *repo.EpisodeRepo
	userRepo                 *repo.UserRepo
	client                   redis.UniversalClient
	cache                    sdomain.Cache[*sdomain.Episode]
}

func NewEpisodeService(c *conf.Bootstrap,
	videoRepo *repo.VideoRepo,
	actorRepo *repo.ActorRepo, videoActorMappingRepo *repo.VideoActorMappingRepo,
	videoSubtitleMappingRepo *repo.EpisodeSubtitleMappingRepo, episodeRepo *repo.EpisodeRepo,
	userRepo *repo.UserRepo, client redis.UniversalClient) *EpisodeService {
	return &EpisodeService{
		c:                        c,
		videoRepo:                videoRepo,
		actorRepo:                actorRepo,
		videoActorMappingRepo:    videoActorMappingRepo,
		videoSubtitleMappingRepo: videoSubtitleMappingRepo,
		episodeRepo:              episodeRepo,
		userRepo:                 userRepo,
		client:                   client,
		cache:                    repo.NewCache[*sdomain.Episode](c, videoRepo.Data.Cache()),
	}
}

func (s *EpisodeService) Get(ctx context.Context, id int64, containPlayUrl bool) (*sdomain.Episode, error) {
	return s.cache.Get(ctx, cache_util.GetCacheActionName(id), func(action string, ctx context.Context) (*sdomain.Episode, error) {
		return s.get(ctx, id, containPlayUrl)
	})
}

func (s *EpisodeService) get(ctx context.Context, id int64, containPlayUrl bool) (*sdomain.Episode, error) {
	userName, err := context_util.GetGenericContext[string](ctx, constant.CTX_UserName)
	if err != nil {
		return nil, err
	}
	currentWatchs, err := s.client.HGet(ctx, fmt.Sprintf(constant.RK_UserWatchCountPrefix, userName), time_util.FormatYYYYMMDD(time.Now())).Int()
	if gerror.HandleRedisNotFoundError(err) != nil {
		return nil, err
	}
	if currentWatchs > constant.MaxWatchCountByDay {
		return nil, gerror.ErrExceedWatchCount(ctx, fmt.Sprintf("%d", constant.MaxWatchCountByDay))
	}
	episode, err := s.episodeRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	subtitleMappings, err := s.videoSubtitleMappingRepo.FindByEpisodeId(ctx, episode.ID)
	if err != nil {
		return nil, err
	}
	episode.Subtitles = subtitleMappings

	currentWatchs++
	s.client.HSet(ctx, fmt.Sprintf(constant.RK_UserWatchCountPrefix, userName), time_util.FormatYYYYMMDD(time.Now()), currentWatchs)

	if containPlayUrl {
		if episode.Url == "" || episode.ExpiredTime == nil || episode.Duration == 0 || episode.ExpiredTime.Sub(time.Now()) < constant.AliyunM3u8EarlyExpireMinutes*time.Minute {
			clientIp, err := context_util.GetGenericContext[string](ctx, constant.CTX_ClientIp)
			if err != nil {
				clientIp = "127.0.0.1"
			}
			url, duration, err := s.getPlayUrl(ctx, episode.XiaoYaPath+"/"+episode.EpisodeTitle, clientIp)
			if err != nil {
				return nil, err
			}
			episode.Url = url
			episode.Duration = duration
			episode.ExpiredTime = lo.ToPtr(time.Now().Add(constant.AliyunM3u8ReallyExpireMinutes * time.Minute))
			err = s.episodeRepo.Updates(ctx, nil, episode)
			if err != nil {
				return nil, err
			}
		}
	}
	return episode, nil
}

func (s *EpisodeService) getPlayUrl(ctx context.Context, xiaoyaPath string, clientIp string) (string, int64, error) {
	//if len(s.c.Dynamic.BoxServerIps) == 0 {
	//	return "", 0, nil
	//}
	boxIps, err := json_util.Unmarshal[[]map[string]string](gerror.HandleRedisStringNotFound(s.client.Get(ctx, constant.RK_BoxIps).Val()))
	if err != nil {
		return "", 0, err
	}
	boxIp, _ := array_util.GetHashElement(boxIps, clientIp)

	transferStoreResult, err := httpclient_util.DoPost[xiaoyadto.TransferStoreReq, xiaoyadto.XiaoyaResult[xiaoyadto.TransferStoreResp]](ctx, boxIp+constant.XiaoYaTransferStorePath, constant.XiaoYaToken, &xiaoyadto.TransferStoreReq{
		Path:     xiaoyaPath,
		Password: "",
	})
	_, ok := err.(*errors.Error)
	if ok {
		//请求登录接口获取token
		loginResult, err := httpclient_util.DoPost[xiaoyadto.LoginReq, xiaoyadto.XiaoyaResult[xiaoyadto.LoginResp]](ctx, boxIp+constant.XiaoYaLoginPath, constant.XiaoYaToken, &xiaoyadto.LoginReq{
			Username: constant.XiaoYaLoginName,
			Password: constant.XiaoYaLoginPassword,
			OtpCode:  "",
		})
		if err != nil {
			return "", 0, err
		}
		if loginResult != nil && loginResult.Code == 200 && loginResult.Data.Token != "" {
			constant.XiaoYaToken = loginResult.Data.Token
			return s.getPlayUrl(ctx, xiaoyaPath, clientIp)
		}
	}
	if err != nil {
		return "", 0, err
	}
	if transferStoreResult == nil || transferStoreResult.Code != 200 {
		return "", 0, gerror.ErrInternal(ctx, "获取播放地址转存失败")
	}

	m3u8Result, err := httpclient_util.DoPost[xiaoyadto.M3u8Req, xiaoyadto.XiaoyaResult[xiaoyadto.M3u8Resp]](ctx, boxIp+constant.XiaoYaM3u8Path, constant.XiaoYaToken, &xiaoyadto.M3u8Req{
		Path:     xiaoyaPath,
		Password: "",
		Method:   "video_preview",
	})
	if err != nil {
		return "", 0, err
	}
	if m3u8Result == nil || m3u8Result.Code != 200 || m3u8Result.Data == nil || m3u8Result.Data.VideoPreviewPlayInfo == nil || len(m3u8Result.Data.VideoPreviewPlayInfo.LiveTranscodingTaskList) == 0 {
		return "", 0, gerror.ErrInternal(ctx, "获取播放地址失败")
	}

	var m3u8Url = ""
	for _, item := range m3u8Result.Data.VideoPreviewPlayInfo.LiveTranscodingTaskList {
		if item.TemplateId == "LD" && item.Status == "finished" && item.Url != "" {
			m3u8Url = item.Url
		}
		if item.TemplateId == "SD" && item.Status == "finished" && item.Url != "" {
			m3u8Url = item.Url
		}
		if item.TemplateId == "HD" && item.Status == "finished" && item.Url != "" {
			m3u8Url = item.Url
		}
		if item.TemplateId == "QHD" && item.Status == "finished" && item.Url != "" {
			m3u8Url = item.Url
		}
	}
	return m3u8Url, int64(m3u8Result.Data.VideoPreviewPlayInfo.Meta.Duration), nil
}

func (s *EpisodeService) UpdateBoxIps(ctx context.Context, boxIps string) error {
	cmd := s.client.Set(ctx, constant.RK_BoxIps, boxIps, 0)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}
