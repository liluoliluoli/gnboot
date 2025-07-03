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
	"github.com/liluoliluoli/gnboot/internal/integration/dto/alidto"
	"github.com/liluoliluoli/gnboot/internal/integration/dto/xiaoyadto"
	"github.com/liluoliluoli/gnboot/internal/repo"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"runtime/debug"
	"time"
)

type EpisodeService struct {
	c                          *conf.Bootstrap
	videoRepo                  *repo.VideoRepo
	actorRepo                  *repo.ActorRepo
	videoActorMappingRepo      *repo.VideoActorMappingRepo
	episodeSubtitleMappingRepo *repo.EpisodeSubtitleMappingRepo
	episodeRepo                *repo.EpisodeRepo
	userRepo                   *repo.UserRepo
	client                     redis.UniversalClient
	cache                      sdomain.Cache[*sdomain.Episode]
	configRepo                 *repo.ConfigRepo
	episodeAudioMappingRepo    *repo.EpisodeAudioMappingRepo
}

func NewEpisodeService(c *conf.Bootstrap,
	videoRepo *repo.VideoRepo,
	actorRepo *repo.ActorRepo, videoActorMappingRepo *repo.VideoActorMappingRepo,
	episodeSubtitleMappingRepo *repo.EpisodeSubtitleMappingRepo, episodeRepo *repo.EpisodeRepo,
	userRepo *repo.UserRepo, client redis.UniversalClient, configRepo *repo.ConfigRepo, episodeAudioMappingRepo *repo.EpisodeAudioMappingRepo) *EpisodeService {
	return &EpisodeService{
		c:                          c,
		videoRepo:                  videoRepo,
		actorRepo:                  actorRepo,
		videoActorMappingRepo:      videoActorMappingRepo,
		episodeSubtitleMappingRepo: episodeSubtitleMappingRepo,
		episodeRepo:                episodeRepo,
		userRepo:                   userRepo,
		client:                     client,
		cache:                      repo.NewCache[*sdomain.Episode](c, videoRepo.Data.Cache()),
		configRepo:                 configRepo,
		episodeAudioMappingRepo:    episodeAudioMappingRepo,
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
	episode, err := s.episodeRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	var currentWatchs = 0
	if episode.VideoId != constant.TrailVideoId { //不是试看目录才记录观看数
		currentWatchs, err = s.client.HGet(ctx, fmt.Sprintf(constant.RK_UserWatchCountPrefix, userName), time_util.FormatYYYYMMDD(time.Now())).Int()
		if gerror.HandleRedisNotFoundError(err) != nil {
			return nil, err
		}
		if currentWatchs > constant.MaxWatchCountByDay {
			return nil, gerror.ErrExceedWatchCount(ctx, fmt.Sprintf("%d", constant.MaxWatchCountByDay))
		}
		currentWatchs++
	}
	packageType := s.client.Get(ctx, fmt.Sprintf(constant.RK_UserPackagePrefix, userName)).Val()
	packageType = lo.Ternary(packageType == "", constant.None, packageType)

	subtitleMappings, err := s.episodeSubtitleMappingRepo.FindByEpisodeId(ctx, episode.ID)
	if err != nil {
		return nil, err
	}
	episode.Subtitles = subtitleMappings
	audioMappings, err := s.episodeAudioMappingRepo.FindByEpisodeId(ctx, episode.ID)
	if err != nil {
		return nil, err
	}
	episode.Audios = audioMappings

	var newRatio = ""
	if containPlayUrl {
		if episode.Url == "" || episode.ExpiredTime == nil || episode.ExpiredTime.Before(time.Now()) {
			url, duration, ratio, subtitles, err := s.getPlayUrl(ctx, episode.XiaoYaPath+"/"+episode.EpisodeTitle, packageType, episode.VideoId, episode.ID)
			if err != nil {
				return nil, err
			}
			episode.Url = url
			episode.Duration = duration
			expiredTime, _ := httpclient_util.ExtractAliM3u8UrlExpireTime(url)
			if expiredTime != nil && expiredTime.Unix() > time.Now().Unix()+duration+constant.AliyunM3u8EarlyExpireSeconds {
				episode.ExpiredTime = lo.ToPtr(expiredTime.Add(-time.Duration(duration+constant.AliyunM3u8EarlyExpireSeconds) * time.Second))
			}
			episode.Ratio = ratio
			newRatio = ratio
			err = s.episodeRepo.Updates(ctx, nil, episode)
			episode.Subtitles = append(episode.Subtitles, subtitles...)
			if err != nil {
				return nil, err
			}
		}
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("AddWatchCount panic recovered: %v\n%s", r, debug.Stack())
			}
		}()
		if episode.VideoId != constant.TrailVideoId {
			s.client.HSet(ctx, fmt.Sprintf(constant.RK_UserWatchCountPrefix, userName), time_util.FormatYYYYMMDD(time.Now()), currentWatchs)
		}
		err := s.videoRepo.AddWatchCount(ctx, nil, episode.VideoId, newRatio)
		if err != nil {
			log.Errorf("增加总观看次数失败: %v", err)
		}
		log.Infof("增加总观看次数成功：%d", episode.VideoId)
	}()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("TransferStoreNextEpisodeToAliyun panic recovered: %v\n%s", r, debug.Stack())
			}
		}()
		err := s.TransferStoreNextEpisodeToAliyun(ctx, episode.VideoId, episode.ID)
		if err != nil {
			log.Errorf("转存阿里云失败: %v", err)
		}
		log.Infof("转存阿里云成功：%d的下一集", episode.ID)
	}()
	return episode, nil
}

func (s *EpisodeService) getPlayUrl(ctx context.Context, xiaoyaPath string, packageType string, videoId int64, episodeId int64) (string, int64, string, []*sdomain.EpisodeSubtitleMapping, error) {
	//if len(s.c.Dynamic.BoxServerIps) == 0 {
	//	return "", 0, nil
	//}
	url, err := s.transferStoreToAliyun(ctx, xiaoyaPath)
	if err != nil {
		return "", 0, "", nil, err
	}
	driveId, fileId, _ := httpclient_util.ExtractAliVideoUrlDriveAndFileId(url)
	clientIp, err := context_util.GetGenericContext[string](ctx, constant.CTX_ClientIp)
	if err != nil {
		clientIp = "127.0.0.1"
	}
	useAliOpenapiStr, err := s.configRepo.GetConfigBySubKey(ctx, constant.Key_BoxIpMapping, constant.SubKey_UseAliOpenapi)
	useAliOpenAPi := lo.Ternary(useAliOpenapiStr == "true", true, false)
	if useAliOpenAPi {
		aliOpenapiDomain, err := s.configRepo.GetConfigBySubKey(ctx, constant.Key_BoxIpMapping, constant.SubKey_AliOpenapiDomain)
		if err != nil {
			return "", 0, "", nil, err
		}
		aliOpenapiToken, err := s.configRepo.GetConfigBySubKey(ctx, constant.Key_BoxIpMapping, constant.SubKey_AliOpenapiToken)
		if err != nil {
			return "", 0, "", nil, err
		}
		headerMap := make(map[string]string)
		headerMap["Authorization"] = aliOpenapiToken
		m3u8Result, err := httpclient_util.DoPost[alidto.GetVideoPreviewPlayInfoReq, alidto.VideoPreviewPlayInfoResp](ctx, aliOpenapiDomain+constant.AliyunM3u8Path, &alidto.GetVideoPreviewPlayInfoReq{
			DriveId:         driveId,
			FileId:          fileId,
			Category:        "live_transcoding",
			UrlExpireSec:    14400,
			GetSubtitleInfo: true,
		}, headerMap)
		if err != nil {
			return "", 0, "", nil, err
		}
		if m3u8Result == nil {
			return "", 0, "", nil, nil
		}
		var m3u8Url = ""
		var radio = ""
		for _, item := range m3u8Result.VideoPreviewPlayInfo.LiveTranscodingTaskList {
			if item.TemplateId == "LD" && item.Status == "finished" && item.Url != "" {
				m3u8Url = item.Url
				radio = item.TemplateId
			}
			if item.TemplateId == "SD" && item.Status == "finished" && item.Url != "" {
				m3u8Url = item.Url
				radio = item.TemplateId
			}
			if packageType == constant.Year || videoId == constant.TrailVideoId {
				if item.TemplateId == "HD" && item.Status == "finished" && item.Url != "" {
					m3u8Url = item.Url
					radio = item.TemplateId
				}
				if item.TemplateId == "QHD" && item.Status == "finished" && item.Url != "" {
					m3u8Url = item.Url
					radio = item.TemplateId
				}
			} else {
				if item.TemplateId == "HD" {
					radio = item.TemplateId
				}
				if item.TemplateId == "QHD" {
					radio = item.TemplateId
				}
			}
		}
		subtitles := make([]*sdomain.EpisodeSubtitleMapping, 0)
		for i, item := range m3u8Result.VideoPreviewPlayInfo.LiveTranscodingSubtitleTaskList {
			if item.Status == "finished" {
				subtitles = append(subtitles, &sdomain.EpisodeSubtitleMapping{
					ID:        int64(-i),
					EpisodeId: episodeId,
					Url:       item.Url,
					Language:  item.Language,
				})
			}
		}
		return m3u8Url, 0, radio, subtitles, nil
	} else {
		boxIps, err := s.configRepo.GetConfigBySubKey(ctx, constant.Key_BoxIpMapping, constant.SubKey_XiaoYaBoxIp)
		if err != nil {
			return "", 0, "", nil, err
		}
		boxIp := array_util.GetHashElement(boxIps, clientIp)
		m3u8Result, err := httpclient_util.DoPost[xiaoyadto.M3u8Req, xiaoyadto.XiaoyaResult[xiaoyadto.M3u8Resp]](ctx, boxIp+constant.XiaoYaM3u8Path, &xiaoyadto.M3u8Req{
			Path:     xiaoyaPath,
			Password: "",
			Method:   "video_preview",
		}, nil)
		if err != nil {
			return "", 0, "", nil, err
		}
		if m3u8Result == nil || m3u8Result.Code != 200 || m3u8Result.Data == nil || m3u8Result.Data.VideoPreviewPlayInfo == nil || len(m3u8Result.Data.VideoPreviewPlayInfo.LiveTranscodingTaskList) == 0 {
			return "", 0, "", nil, gerror.ErrInternal(ctx, "获取播放地址失败")
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
			if packageType == constant.Year {
				if item.TemplateId == "HD" && item.Status == "finished" && item.Url != "" {
					m3u8Url = item.Url
					radio = item.TemplateId
				}
				if item.TemplateId == "QHD" && item.Status == "finished" && item.Url != "" {
					m3u8Url = item.Url
					radio = item.TemplateId
				}
			} else {
				if item.TemplateId == "HD" {
					radio = item.TemplateId
				}
				if item.TemplateId == "QHD" {
					radio = item.TemplateId
				}
			}
		}
		return m3u8Url, int64(m3u8Result.Data.VideoPreviewPlayInfo.Meta.Duration), radio, nil, nil
	}
}

func (s *EpisodeService) TransferStoreNextEpisodeToAliyun(ctx context.Context, videoId int64, currentEpisodeId int64) error {
	episode, err := s.episodeRepo.Next(ctx, videoId, currentEpisodeId)
	if err != nil {
		return err
	}
	if episode == nil {
		return nil
	}
	url, err := s.transferStoreToAliyun(ctx, episode.XiaoYaPath+"/"+episode.EpisodeTitle)
	if err != nil {
		return err
	}
	driveId, fileId, _ := httpclient_util.ExtractAliVideoUrlDriveAndFileId(url)
	episode.AliDriveId = driveId
	episode.AliFileId = fileId
	err = s.episodeRepo.Updates(ctx, nil, episode)
	if err != nil {
		return err
	}
	return nil
}

func (s *EpisodeService) transferStoreToAliyun(ctx context.Context, xiaoyaPath string) (string, error) {
	clientIp, err := context_util.GetGenericContext[string](ctx, constant.CTX_ClientIp)
	if err != nil {
		clientIp = "127.0.0.1"
	}
	boxIps, err := s.configRepo.GetConfigBySubKey(ctx, constant.Key_BoxIpMapping, constant.SubKey_XiaoYaBoxIp)
	if err != nil {
		return "", err
	}
	boxIp := array_util.GetHashElement(boxIps, clientIp)
	headerMap := make(map[string]string)
	headerMap["Authorization"] = constant.XiaoYaToken
	transferStoreResult, err := httpclient_util.DoPost[xiaoyadto.TransferStoreReq, xiaoyadto.XiaoyaResult[xiaoyadto.TransferStoreResp]](ctx, boxIp+constant.XiaoYaTransferStorePath, &xiaoyadto.TransferStoreReq{
		Path:     xiaoyaPath,
		Password: "",
	}, headerMap)
	if err != nil {
		return "", err
	}
	if transferStoreResult == nil || transferStoreResult.Code != 200 || transferStoreResult.Data == nil || transferStoreResult.Data.RawUrl == "" {
		return "", gerror.ErrInternal(ctx, "获取播放地址转存失败")
	}
	return transferStoreResult.Data.RawUrl, nil
}

func (s *EpisodeService) UpdateConfigs(ctx context.Context, configs string) error {
	cmd := s.client.Set(ctx, constant.RK_Configs, configs, 0)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	s.configRepo.InitConfig(ctx)
	return nil
}
