package service

import (
	"context"
	"fmt"
	"github.com/liluoliluoli/gnboot/internal/common/constant"
	"github.com/liluoliluoli/gnboot/internal/common/gerror"
	"github.com/liluoliluoli/gnboot/internal/common/utils/cache_util"
	"github.com/liluoliluoli/gnboot/internal/common/utils/context_util"
	"github.com/liluoliluoli/gnboot/internal/common/utils/time_util"
	"github.com/liluoliluoli/gnboot/internal/conf"
	"github.com/liluoliluoli/gnboot/internal/repo"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/redis/go-redis/v9"
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

func (s *EpisodeService) Get(ctx context.Context, id int64) (*sdomain.Episode, error) {
	return s.cache.Get(ctx, cache_util.GetCacheActionName(id), func(action string, ctx context.Context) (*sdomain.Episode, error) {
		return s.get(ctx, id)
	})
}

func (s *EpisodeService) get(ctx context.Context, id int64) (*sdomain.Episode, error) {
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
	return episode, nil
}
