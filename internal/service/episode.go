package service

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/liluoliluoli/gnboot/internal/common/constant"
	"github.com/liluoliluoli/gnboot/internal/common/gerror"
	"github.com/liluoliluoli/gnboot/internal/common/utils/array_util"
	"github.com/liluoliluoli/gnboot/internal/common/utils/cache_util"
	"github.com/liluoliluoli/gnboot/internal/common/utils/context_util"
	"github.com/liluoliluoli/gnboot/internal/common/utils/httpclient_util"
	"github.com/liluoliluoli/gnboot/internal/common/utils/time_util"
	"github.com/liluoliluoli/gnboot/internal/conf"
	"github.com/liluoliluoli/gnboot/internal/integration/dto/xiaoyadto"
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
	configRepo               *repo.ConfigRepo
}

func NewEpisodeService(c *conf.Bootstrap,
	videoRepo *repo.VideoRepo,
	actorRepo *repo.ActorRepo, videoActorMappingRepo *repo.VideoActorMappingRepo,
	videoSubtitleMappingRepo *repo.EpisodeSubtitleMappingRepo, episodeRepo *repo.EpisodeRepo,
	userRepo *repo.UserRepo, client redis.UniversalClient, configRepo *repo.ConfigRepo) *EpisodeService {
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
		configRepo:               configRepo,
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

	var radio = ""
	if containPlayUrl {
		if episode.Url == "" || episode.ExpiredTime == nil || episode.Duration == 0 || episode.ExpiredTime.Sub(time.Now()) < constant.AliyunM3u8EarlyExpireMinutes*time.Minute {
			url, duration, vRadio, err := s.getPlayUrl(ctx, episode.XiaoYaPath+"/"+episode.EpisodeTitle)
			if err != nil {
				return nil, err
			}
			episode.Url = url
			episode.Duration = duration
			episode.ExpiredTime = lo.ToPtr(time.Now().Add(constant.AliyunM3u8ReallyExpireMinutes * time.Minute))
			radio = vRadio
			err = s.episodeRepo.Updates(ctx, nil, episode)
			if err != nil {
				return nil, err
			}
		}
	}
	go func() {
		err := s.videoRepo.AddWatchCount(ctx, nil, episode.VideoId, radio)
		if err != nil {
			log.Errorf("增加总观看次数失败: %v", err)
		}
	}()
	go func() {
		err := s.TransferStoreNextEpisodeToAliyun(ctx, episode.VideoId, episode.ID)
		if err != nil {
			log.Errorf("转存阿里云失败: %v", err)
		}
	}()
	return episode, nil
}

func (s *EpisodeService) getPlayUrl(ctx context.Context, xiaoyaPath string) (string, int64, string, error) {
	//if len(s.c.Dynamic.BoxServerIps) == 0 {
	//	return "", 0, nil
	//}
	err := s.transferStoreToAliyun(ctx, xiaoyaPath)
	if err != nil {
		return "", 0, "", err
	}
	clientIp, err := context_util.GetGenericContext[string](ctx, constant.CTX_ClientIp)
	if err != nil {
		clientIp = "127.0.0.1"
	}
	boxIps, err := s.configRepo.GetConfigBySubKey(ctx, constant.Key_BoxIpMapping, constant.SubKey_XiaoYaBoxIp)
	if err != nil {
		return "", 0, "", err
	}
	boxIp := array_util.GetHashElement(boxIps, clientIp)
	m3u8Result, err := httpclient_util.DoPost[xiaoyadto.M3u8Req, xiaoyadto.XiaoyaResult[xiaoyadto.M3u8Resp]](ctx, boxIp+constant.XiaoYaM3u8Path, &xiaoyadto.M3u8Req{
		Path:     xiaoyaPath,
		Password: "",
		Method:   "video_preview",
	}, nil)
	if err != nil {
		return "", 0, "", err
	}
	if m3u8Result == nil || m3u8Result.Code != 200 || m3u8Result.Data == nil || m3u8Result.Data.VideoPreviewPlayInfo == nil || len(m3u8Result.Data.VideoPreviewPlayInfo.LiveTranscodingTaskList) == 0 {
		return "", 0, "", gerror.ErrInternal(ctx, "获取播放地址失败")
	}

	var m3u8Url = ""
	var radio = ""
	for _, item := range m3u8Result.Data.VideoPreviewPlayInfo.LiveTranscodingTaskList {
		if item.TemplateId == "LD" && item.Status == "finished" && item.Url != "" {
			m3u8Url = item.Url
			radio = item.TemplateId
		}
		if item.TemplateId == "SD" && item.Status == "finished" && item.Url != "" {
			m3u8Url = item.Url
			radio = item.TemplateId
		}
		if item.TemplateId == "HD" && item.Status == "finished" && item.Url != "" {
			m3u8Url = item.Url
			radio = item.TemplateId
		}
		if item.TemplateId == "QHD" && item.Status == "finished" && item.Url != "" {
			m3u8Url = item.Url
			radio = item.TemplateId
		}
	}
	return m3u8Url, int64(m3u8Result.Data.VideoPreviewPlayInfo.Meta.Duration), radio, nil
}

func (s *EpisodeService) TransferStoreNextEpisodeToAliyun(ctx context.Context, videoId int64, currentEpisodeId int64) error {
	episode, err := s.episodeRepo.Next(ctx, videoId, currentEpisodeId)
	if err != nil {
		return err
	}
	if episode == nil {
		return nil
	}
	err = s.transferStoreToAliyun(ctx, episode.XiaoYaPath+"/"+episode.EpisodeTitle)
	if err != nil {
		return err
	}
	return nil
}

func (s *EpisodeService) transferStoreToAliyun(ctx context.Context, xiaoyaPath string) error {
	clientIp, err := context_util.GetGenericContext[string](ctx, constant.CTX_ClientIp)
	if err != nil {
		clientIp = "127.0.0.1"
	}
	boxIps, err := s.configRepo.GetConfigBySubKey(ctx, constant.Key_BoxIpMapping, constant.SubKey_XiaoYaBoxIp)
	if err != nil {
		return err
	}
	boxIp := array_util.GetHashElement(boxIps, clientIp)
	headerMap := make(map[string]string)
	headerMap["Authorization"] = constant.XiaoYaToken
	transferStoreResult, err := httpclient_util.DoPost[xiaoyadto.TransferStoreReq, xiaoyadto.XiaoyaResult[xiaoyadto.TransferStoreResp]](ctx, boxIp+constant.XiaoYaTransferStorePath, &xiaoyadto.TransferStoreReq{
		Path:     xiaoyaPath,
		Password: "",
	}, headerMap)
	if err != nil {
		return err
	}
	if transferStoreResult == nil || transferStoreResult.Code != 200 {
		return gerror.ErrInternal(ctx, "获取播放地址转存失败")
	}
	return nil
}

func (s *EpisodeService) UpdateConfigs(ctx context.Context, configs string) error {
	cmd := s.client.Set(ctx, constant.RK_Configs, configs, 0)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	s.configRepo.InitConfig(ctx)
	return nil
}
