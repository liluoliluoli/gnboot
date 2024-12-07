package service

import (
	"context"
	"github.com/liluoliluoli/gnboot/internal/common/utils/cache_util"
	"github.com/liluoliluoli/gnboot/internal/conf"
	"github.com/liluoliluoli/gnboot/internal/repo"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
)

type SeasonService struct {
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
	seasonRepo               *repo.SeasonRepo
	episodeRepo              *repo.EpisodeRepo
	cache                    sdomain.Cache[*sdomain.Season]
}

func NewSeasonService(c *conf.Bootstrap,
	movieRepo *repo.MovieRepo,
	genreRepo *repo.GenreRepo, videoGenreMappingRepo *repo.VideoGenreMappingRepo,
	actorRepo *repo.ActorRepo, videoActorMappingRepo *repo.VideoActorMappingRepo,
	studioRepo *repo.StudioRepo, videoStudioMappingRepo *repo.VideoStudioMappingRepo,
	keywordRepo *repo.KeywordRepo, videoKeywordMappingRepo *repo.VideoKeywordMappingRepo,
	videoSubtitleMappingRepo *repo.VideoSubtitleMappingRepo,
	seasonRepo *repo.SeasonRepo, episodeRepo *repo.EpisodeRepo) *SeasonService {
	return &SeasonService{
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
		seasonRepo:               seasonRepo,
		episodeRepo:              episodeRepo,
		cache:                    repo.NewCache[*sdomain.Season](c, movieRepo.Data.Cache()),
	}
}

func (s *SeasonService) Get(ctx context.Context, id int64) (*sdomain.Season, error) {
	return s.cache.Get(ctx, cache_util.GetCacheActionName(id), func(action string, ctx context.Context) (*sdomain.Season, error) {
		return s.get(ctx, id)
	})
}

func (s *SeasonService) get(ctx context.Context, id int64) (*sdomain.Season, error) {
	season, err := s.seasonRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if season == nil {
		return nil, nil
	}
	episodes, err := s.episodeRepo.QueryBySeasonId(ctx, season.ID)
	if err != nil {
		return nil, err
	}
	season.Episodes = episodes
	return season, nil
}
