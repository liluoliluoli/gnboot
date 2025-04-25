package service

import (
	"context"
	"github.com/liluoliluoli/gnboot/internal/common/utils/cache_util"
	"github.com/liluoliluoli/gnboot/internal/conf"
	"github.com/liluoliluoli/gnboot/internal/repo"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
)

type EpisodeService struct {
	c                        *conf.Bootstrap
	videoRepo                *repo.VideoRepo
	actorRepo                *repo.ActorRepo
	videoActorMappingRepo    *repo.VideoActorMappingRepo
	videoSubtitleMappingRepo *repo.EpisodeSubtitleMappingRepo
	episodeRepo              *repo.EpisodeRepo
	cache                    sdomain.Cache[*sdomain.Episode]
}

func NewEpisodeService(c *conf.Bootstrap,
	videoRepo *repo.VideoRepo,
	actorRepo *repo.ActorRepo, videoActorMappingRepo *repo.VideoActorMappingRepo,
	videoSubtitleMappingRepo *repo.EpisodeSubtitleMappingRepo, episodeRepo *repo.EpisodeRepo) *EpisodeService {
	return &EpisodeService{
		c:                        c,
		videoRepo:                videoRepo,
		actorRepo:                actorRepo,
		videoActorMappingRepo:    videoActorMappingRepo,
		videoSubtitleMappingRepo: videoSubtitleMappingRepo,
		episodeRepo:              episodeRepo,
		cache:                    repo.NewCache[*sdomain.Episode](c, videoRepo.Data.Cache()),
	}
}

func (s *EpisodeService) Get(ctx context.Context, id int64) (*sdomain.Episode, error) {
	return s.cache.Get(ctx, cache_util.GetCacheActionName(id), func(action string, ctx context.Context) (*sdomain.Episode, error) {
		return s.get(ctx, id)
	})
}

func (s *EpisodeService) get(ctx context.Context, id int64) (*sdomain.Episode, error) {
	episode, err := s.episodeRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	subtitleMappings, err := s.videoSubtitleMappingRepo.FindByEpisodeId(ctx, episode.ID)
	if err != nil {
		return nil, err
	}
	episode.Subtitles = subtitleMappings
	return episode, nil
}
