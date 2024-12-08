package service

import (
	"context"
	"github.com/liluoliluoli/gnboot/internal/common/utils/cache_util"
	"github.com/liluoliluoli/gnboot/internal/conf"
	"github.com/liluoliluoli/gnboot/internal/repo"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/samber/lo"
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
	seriesRepo               *repo.SeriesRepo
	cache                    sdomain.Cache[*sdomain.Season]
}

func NewSeasonService(c *conf.Bootstrap,
	movieRepo *repo.MovieRepo,
	genreRepo *repo.GenreRepo, videoGenreMappingRepo *repo.VideoGenreMappingRepo,
	actorRepo *repo.ActorRepo, videoActorMappingRepo *repo.VideoActorMappingRepo,
	studioRepo *repo.StudioRepo, videoStudioMappingRepo *repo.VideoStudioMappingRepo,
	keywordRepo *repo.KeywordRepo, videoKeywordMappingRepo *repo.VideoKeywordMappingRepo,
	videoSubtitleMappingRepo *repo.VideoSubtitleMappingRepo,
	seasonRepo *repo.SeasonRepo, episodeRepo *repo.EpisodeRepo, seriesRepo *repo.SeriesRepo) *SeasonService {
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
		seriesRepo:               seriesRepo,
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

func (s *SeasonService) FindBySeriesId(ctx context.Context, seriesId int64, actors []*sdomain.Actor) ([]*sdomain.Season, error) {
	return s.cache.List(ctx, cache_util.GetCacheActionName(seriesId), func(action string, ctx context.Context) ([]*sdomain.Season, error) {
		return s.findBySeriesId(ctx, seriesId, actors)
	})
}

func (s *SeasonService) findBySeriesId(ctx context.Context, seriesId int64, actors []*sdomain.Actor) ([]*sdomain.Season, error) {
	series, err := s.seriesRepo.Get(ctx, seriesId)
	if err != nil {
		return nil, err
	}
	seasons, err := s.seasonRepo.QueryBySeriesId(ctx, seriesId)
	if err != nil {
		return nil, err
	}
	for _, season := range seasons {
		episodes, err := s.episodeRepo.QueryBySeasonId(ctx, season.ID)
		if err != nil {
			return nil, err
		}
		season.Episodes = lo.Map(episodes, func(item *sdomain.Episode, index int) *sdomain.Episode {
			item.SeriesTitle = series.Title
			item.SeasonTitle = season.Title
			item.SeasonId = season.ID
			item.Season = season.Season
			item.Actors = actors
			return item
		})
	}
	return seasons, nil
}
