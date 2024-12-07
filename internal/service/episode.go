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
	movieRepo                *repo.MovieRepo
	genreRepo                *repo.GenreRepo
	videoGenreMappingRepo    *repo.VideoGenreMappingRepo
	actorRepo                *repo.ActorRepo
	videoActorMappingRepo    *repo.VideoActorMappingRepo
	studioRepo               *repo.StudioRepo
	videoStudioMappingRepo   *repo.VideoStudioMappingRepo
	keywordRepo              *repo.KeywordRepo
	videoKeywordMappingRepo  *repo.VideoKeywordMappingRepo
	videoSubtitleMappingRepo *repo.VideoSubtitleMappingRepo
	episodeRepo              *repo.EpisodeRepo
	cache                    sdomain.Cache[*sdomain.Episode]
}

func NewEpisodeService(c *conf.Bootstrap,
	movieRepo *repo.MovieRepo,
	genreRepo *repo.GenreRepo, videoGenreMappingRepo *repo.VideoGenreMappingRepo,
	actorRepo *repo.ActorRepo, videoActorMappingRepo *repo.VideoActorMappingRepo,
	studioRepo *repo.StudioRepo, videoStudioMappingRepo *repo.VideoStudioMappingRepo,
	keywordRepo *repo.KeywordRepo, videoKeywordMappingRepo *repo.VideoKeywordMappingRepo,
	videoSubtitleMappingRepo *repo.VideoSubtitleMappingRepo, episodeRepo *repo.EpisodeRepo) *EpisodeService {
	return &EpisodeService{
		c:                        c,
		movieRepo:                movieRepo,
		genreRepo:                genreRepo,
		videoGenreMappingRepo:    videoGenreMappingRepo,
		actorRepo:                actorRepo,
		videoActorMappingRepo:    videoActorMappingRepo,
		studioRepo:               studioRepo,
		videoStudioMappingRepo:   videoStudioMappingRepo,
		keywordRepo:              keywordRepo,
		videoKeywordMappingRepo:  videoKeywordMappingRepo,
		videoSubtitleMappingRepo: videoSubtitleMappingRepo,
		episodeRepo:              episodeRepo,
		cache:                    repo.NewCache[*sdomain.Episode](c, movieRepo.Data.Cache()),
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
	return episode, nil
}
