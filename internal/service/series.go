package service

import (
	"context"
	"github.com/liluoliluoli/gnboot/internal/common/constant"
	"github.com/liluoliluoli/gnboot/internal/common/utils/cache_util"
	"github.com/liluoliluoli/gnboot/internal/conf"
	"github.com/liluoliluoli/gnboot/internal/repo"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/samber/lo"
)

type SeriesService struct {
	c                        *conf.Bootstrap
	seriesRepo               *repo.SeriesRepo
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
	videoUserMappingRepo     *repo.VideoUserMappingRepo
	cache                    sdomain.Cache[*sdomain.Series]
}

func NewSeriesService(c *conf.Bootstrap,
	seriesRepo *repo.SeriesRepo,
	genreRepo *repo.GenreRepo, videoGenreMappingRepo *repo.VideoGenreMappingRepo,
	actorRepo *repo.ActorRepo, videoActorMappingRepo *repo.VideoActorMappingRepo,
	studioRepo *repo.StudioRepo, videoStudioMappingRepo *repo.VideoStudioMappingRepo,
	keywordRepo *repo.KeywordRepo, videoKeywordMappingRepo *repo.VideoKeywordMappingRepo,
	videoSubtitleMappingRepo *repo.VideoSubtitleMappingRepo,
	seasonRepo *repo.SeasonRepo,
	episodeRepo *repo.EpisodeRepo,
	videoUserMappingRepo *repo.VideoUserMappingRepo) *SeriesService {
	return &SeriesService{
		c:                        c,
		seriesRepo:               seriesRepo,
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
		videoUserMappingRepo:     videoUserMappingRepo,
		cache:                    repo.NewCache[*sdomain.Series](c, seriesRepo.Data.Cache()),
	}
}

func (s *SeriesService) Get(ctx context.Context, id int64, userId int64) (*sdomain.Series, error) {
	return s.cache.Get(ctx, cache_util.GetCacheActionName(id), func(action string, ctx context.Context) (*sdomain.Series, error) {
		return s.get(ctx, id, userId)
	})
}

func (s *SeriesService) get(ctx context.Context, id int64, userId int64) (*sdomain.Series, error) {
	item, err := s.seriesRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	genresMap, err := s.buildSeriesGenresMap(ctx, []*sdomain.Series{item})
	if err != nil {
		return nil, err
	}
	keywordsMap, err := s.buildSeriesKeywordsMap(ctx, []*sdomain.Series{item})
	if err != nil {
		return nil, err
	}
	studiosMap, err := s.buildSeriesStudiosMap(ctx, []*sdomain.Series{item})
	if err != nil {
		return nil, err
	}
	actorsMap, err := s.buildSeriesActorsMap(ctx, []*sdomain.Series{item})
	if err != nil {
		return nil, err
	}
	seasonsMap, err := s.buildSeriesSeasonsMap(ctx, []*sdomain.Series{item})
	if err != nil {
		return nil, err
	}
	nextToPlayEpisodeMap, err := s.buildSeriesNextToPlayEpisodeMap(ctx, []*sdomain.Series{item}, userId)
	if err != nil {
		return nil, err
	}
	item.Genres = genresMap[item.ID]
	item.Keywords = keywordsMap[item.ID]
	item.Studios = studiosMap[item.ID]
	item.Actors = actorsMap[item.ID]
	item.Seasons = seasonsMap[item.ID]
	item.NextToPlay = nextToPlayEpisodeMap[item.ID]
	return item, nil
}

func (s *SeriesService) Page(ctx context.Context, condition *sdomain.SearchSeries, userId int64) (*sdomain.PageResult[*sdomain.Series], error) {
	rp, err := s.cache.Page(ctx, cache_util.GetCacheActionName(condition), func(action string, ctx context.Context) (*sdomain.PageResult[*sdomain.Series], error) {
		return s.page(ctx, condition, userId)
	})
	if err != nil {
		return nil, err
	}
	return rp, nil
}

func (s *SeriesService) page(ctx context.Context, condition *sdomain.SearchSeries, userId int64) (*sdomain.PageResult[*sdomain.Series], error) {
	filterIds := make([]int64, 0)
	if condition.Id != 0 && condition.Type == "" {
		if condition.Type == constant.FilterType_genre {
			genreMappings, err := s.videoGenreMappingRepo.FindByGenreIdAndVideoType(ctx, condition.Id, constant.VideoType_series)
			if err != nil {
				return nil, err
			}
			filterIds = append(filterIds, lo.Map(genreMappings, func(item *sdomain.VideoGenreMapping, index int) int64 {
				return item.VideoId
			})...)
		}
		if condition.Type == constant.FilterType_studio {
			studioMappings, err := s.videoStudioMappingRepo.FindByStudioIdAndVideoType(ctx, condition.Id, constant.VideoType_series)
			if err != nil {
				return nil, err
			}
			filterIds = append(filterIds, lo.Map(studioMappings, func(item *sdomain.VideoStudioMapping, index int) int64 {
				return item.VideoId
			})...)
		}
		if condition.Type == constant.FilterType_keyword {
			keywordMappings, err := s.videoKeywordMappingRepo.FindByKeywordIdAndVideoType(ctx, condition.Id, constant.VideoType_series)
			if err != nil {
				return nil, err
			}
			filterIds = append(filterIds, lo.Map(keywordMappings, func(item *sdomain.VideoKeywordMapping, index int) int64 {
				return item.VideoId
			})...)
		}
		if condition.Type == constant.FilterType_actor {
			actorMappings, err := s.videoActorMappingRepo.FindByActorIdAndVideoType(ctx, condition.Id, constant.VideoType_series)
			if err != nil {
				return nil, err
			}
			filterIds = append(filterIds, lo.Map(actorMappings, func(item *sdomain.VideoActorMapping, index int) int64 {
				return item.VideoId
			})...)
		}
	} else if condition.FilterByNextPlay {
		userSeriesMappings, err := s.videoUserMappingRepo.FindByUserIdAndVideoIdAndType(ctx, userId, nil, constant.VideoType_episode)
		if err != nil {
			return nil, err
		}
		episodes, err := s.episodeRepo.QueryByIds(ctx, lo.Map(userSeriesMappings, func(item *sdomain.VideoUserMapping, index int) int64 {
			return item.VideoId
		}))
		if err != nil {
			return nil, err
		}
		seasons, err := s.seasonRepo.QueryByIds(ctx, lo.Map(episodes, func(item *sdomain.Episode, index int) int64 {
			return item.SeasonId
		}))
		if err != nil {
			return nil, err
		}
		filterIds = append(filterIds, lo.Map(seasons, func(item *sdomain.Season, index int) int64 {
			return item.SeriesId
		})...)
	}
	condition.FilterIds = filterIds
	pageResult, err := s.seriesRepo.Page(ctx, condition)
	if err != nil {
		return nil, err
	}
	if pageResult != nil && len(pageResult.List) != 0 {
		genresMap, err := s.buildSeriesGenresMap(ctx, pageResult.List)
		if err != nil {
			return nil, err
		}
		keywordsMap, err := s.buildSeriesKeywordsMap(ctx, pageResult.List)
		if err != nil {
			return nil, err
		}
		studiosMap, err := s.buildSeriesStudiosMap(ctx, pageResult.List)
		if err != nil {
			return nil, err
		}
		actorsMap, err := s.buildSeriesActorsMap(ctx, pageResult.List)
		if err != nil {
			return nil, err
		}
		seasonsMap, err := s.buildSeriesSeasonsMap(ctx, pageResult.List)
		if err != nil {
			return nil, err
		}
		nextToPlayEpisodeMap, err := s.buildSeriesNextToPlayEpisodeMap(ctx, pageResult.List, userId)
		if err != nil {
			return nil, err
		}
		for _, item := range pageResult.List {
			item.Genres = genresMap[item.ID]
			item.Keywords = keywordsMap[item.ID]
			item.Studios = studiosMap[item.ID]
			item.Actors = actorsMap[item.ID]
			item.Seasons = seasonsMap[item.ID]
			item.NextToPlay = nextToPlayEpisodeMap[item.ID]
		}
	}
	return pageResult, nil
}

// 补齐genre
func (s *SeriesService) buildSeriesGenresMap(ctx context.Context, series []*sdomain.Series) (map[int64][]*sdomain.Genre, error) {
	seriesIds := lo.Map(series, func(item *sdomain.Series, index int) int64 {
		return item.ID
	})
	genreMappings, err := s.videoGenreMappingRepo.FindByVideoIdAndType(ctx, seriesIds, constant.VideoType_series)
	if err != nil {
		return nil, err
	}
	genres, err := s.genreRepo.FindByIds(ctx, lo.Map(genreMappings, func(item *sdomain.VideoGenreMapping, index int) int64 {
		return item.GenreId
	}))
	if err != nil {
		return nil, err
	}
	genresMap := lo.SliceToMap(genres, func(item *sdomain.Genre) (int64, *sdomain.Genre) {
		return item.ID, item
	})
	rsMap := make(map[int64][]*sdomain.Genre)
	for _, genreMapping := range genreMappings {
		if _, ok := rsMap[genreMapping.VideoId]; !ok {
			rsMap[genreMapping.VideoId] = make([]*sdomain.Genre, 0)
		}
		rsMap[genreMapping.VideoId] = append(rsMap[genreMapping.VideoId], genresMap[genreMapping.GenreId])
	}
	return rsMap, nil
}

// 补齐actor
func (s *SeriesService) buildSeriesActorsMap(ctx context.Context, series []*sdomain.Series) (map[int64][]*sdomain.Actor, error) {
	seriesIds := lo.Map(series, func(item *sdomain.Series, index int) int64 {
		return item.ID
	})
	actorMappings, err := s.videoActorMappingRepo.FindByVideoIdAndType(ctx, seriesIds, constant.VideoType_series)
	if err != nil {
		return nil, err
	}
	actors, err := s.actorRepo.FindByIds(ctx, lo.Map(actorMappings, func(item *sdomain.VideoActorMapping, index int) int64 {
		return item.ActorId
	}))
	if err != nil {
		return nil, err
	}
	actorsMap := lo.SliceToMap(actors, func(item *sdomain.Actor) (int64, *sdomain.Actor) {
		return item.ID, item
	})
	rsMap := make(map[int64][]*sdomain.Actor)
	for _, actorMapping := range actorMappings {
		if _, ok := rsMap[actorMapping.VideoId]; !ok {
			rsMap[actorMapping.VideoId] = make([]*sdomain.Actor, 0)
		}
		actorsMap[actorMapping.ActorId].Character = actorMapping.Character
		rsMap[actorMapping.VideoId] = append(rsMap[actorMapping.VideoId], actorsMap[actorMapping.ActorId])
	}
	return rsMap, nil
}

// 补齐keyword
func (s *SeriesService) buildSeriesKeywordsMap(ctx context.Context, series []*sdomain.Series) (map[int64][]*sdomain.Keyword, error) {
	seriesIds := lo.Map(series, func(item *sdomain.Series, index int) int64 {
		return item.ID
	})
	keywordMappings, err := s.videoKeywordMappingRepo.FindByVideoIdAndType(ctx, seriesIds, constant.VideoType_series)
	if err != nil {
		return nil, err
	}
	keywords, err := s.keywordRepo.FindByIds(ctx, lo.Map(keywordMappings, func(item *sdomain.VideoKeywordMapping, index int) int64 {
		return item.KeywordId
	}))
	if err != nil {
		return nil, err
	}
	keywordsMap := lo.SliceToMap(keywords, func(item *sdomain.Keyword) (int64, *sdomain.Keyword) {
		return item.ID, item
	})
	rsMap := make(map[int64][]*sdomain.Keyword)
	for _, keywordMapping := range keywordMappings {
		if _, ok := rsMap[keywordMapping.VideoId]; !ok {
			rsMap[keywordMapping.VideoId] = make([]*sdomain.Keyword, 0)
		}
		rsMap[keywordMapping.VideoId] = append(rsMap[keywordMapping.VideoId], keywordsMap[keywordMapping.KeywordId])
	}
	return rsMap, nil
}

// 补齐studio
func (s *SeriesService) buildSeriesStudiosMap(ctx context.Context, series []*sdomain.Series) (map[int64][]*sdomain.Studio, error) {
	seriesIds := lo.Map(series, func(item *sdomain.Series, index int) int64 {
		return item.ID
	})
	studioMappings, err := s.videoStudioMappingRepo.FindByVideoIdAndType(ctx, seriesIds, constant.VideoType_series)
	if err != nil {
		return nil, err
	}
	studios, err := s.studioRepo.FindByIds(ctx, lo.Map(studioMappings, func(item *sdomain.VideoStudioMapping, index int) int64 {
		return item.StudioId
	}))
	if err != nil {
		return nil, err
	}
	studiosMap := lo.SliceToMap(studios, func(item *sdomain.Studio) (int64, *sdomain.Studio) {
		return item.ID, item
	})
	rsMap := make(map[int64][]*sdomain.Studio)
	for _, studioMapping := range studioMappings {
		if _, ok := rsMap[studioMapping.VideoId]; !ok {
			rsMap[studioMapping.VideoId] = make([]*sdomain.Studio, 0)
		}
		rsMap[studioMapping.VideoId] = append(rsMap[studioMapping.VideoId], studiosMap[studioMapping.StudioId])
	}
	return rsMap, nil
}

// 补齐subtitle
func (s *SeriesService) buildSeriesSubtitlesMap(ctx context.Context, series []*sdomain.Series) (map[int64][]*sdomain.VideoSubtitleMapping, error) {
	seriesIds := lo.Map(series, func(item *sdomain.Series, index int) int64 {
		return item.ID
	})
	subtitleMappings, err := s.videoSubtitleMappingRepo.FindByVideoIdAndType(ctx, seriesIds, constant.VideoType_series)
	if err != nil {
		return nil, err
	}
	return lo.GroupBy(subtitleMappings, func(item *sdomain.VideoSubtitleMapping) int64 {
		return item.VideoId
	}), nil
}

// 补齐season
func (s *SeriesService) buildSeriesSeasonsMap(ctx context.Context, series []*sdomain.Series) (map[int64][]*sdomain.Season, error) {
	seriesIds := lo.Map(series, func(item *sdomain.Series, index int) int64 {
		return item.ID
	})
	rsMap := make(map[int64][]*sdomain.Season)
	for _, seriesId := range seriesIds {
		seasons, err := s.seasonRepo.QueryBySeriesId(ctx, seriesId)
		if err != nil {
			return nil, err
		}
		for _, season := range seasons {
			episodes, err := s.episodeRepo.QueryBySeasonId(ctx, season.ID)
			if err != nil {
				return nil, err
			}
			season.Episodes = episodes
		}
		rsMap[seriesId] = seasons
	}
	return rsMap, nil
}

// 补齐next to play episode，写入时要保证一个series只会有一条episode记录
func (s *SeriesService) buildSeriesNextToPlayEpisodeMap(ctx context.Context, series []*sdomain.Series, userId int64) (map[int64]*sdomain.Episode, error) {
	seriesIds := lo.Map(series, func(item *sdomain.Series, index int) int64 {
		return item.ID
	})
	rsMap := make(map[int64]*sdomain.Episode)
	userEpisodeMappings, err := s.videoUserMappingRepo.FindByUserIdAndVideoIdAndType(ctx, userId, nil, constant.VideoType_episode)
	if err != nil {
		return nil, err
	}
	userEpisodeMappingsMap := lo.SliceToMap(userEpisodeMappings, func(item *sdomain.VideoUserMapping) (int64, *sdomain.VideoUserMapping) {
		return item.VideoId, item
	})
	episodes, err := s.episodeRepo.QueryByIds(ctx, lo.Map(userEpisodeMappings, func(item *sdomain.VideoUserMapping, index int) int64 {
		return item.VideoId
	}))
	if err != nil {
		return nil, err
	}
	seasonEpisodeMap := lo.SliceToMap(episodes, func(item *sdomain.Episode) (int64, *sdomain.Episode) {
		item.LastPlayedTime = userEpisodeMappingsMap[item.ID].LastPlayedTime
		item.LastPlayedPosition = userEpisodeMappingsMap[item.ID].LastPlayedPosition
		return item.SeasonId, item
	})
	seasons, err := s.seasonRepo.QueryByIds(ctx, lo.Map(episodes, func(item *sdomain.Episode, index int) int64 {
		return item.SeasonId
	}))
	seriesSeasonMap := lo.SliceToMap(seasons, func(item *sdomain.Season) (int64, *sdomain.Season) {
		return item.SeriesId, item
	})
	if err != nil {
		return nil, err
	}

	for _, seriesId := range seriesIds {
		if _, ok := seriesSeasonMap[seriesId]; ok {
			if _, ok = seasonEpisodeMap[seriesSeasonMap[seriesId].ID]; ok {
				rsMap[seriesId] = seasonEpisodeMap[seriesSeasonMap[seriesId].ID]
			}
		}
	}
	return rsMap, nil
}
