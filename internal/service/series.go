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
	cache                    sdomain.Cache[*sdomain.Series]
}

func NewSeriesService(c *conf.Bootstrap,
	seriesRepo *repo.SeriesRepo,
	genreRepo *repo.GenreRepo, videoGenreMappingRepo *repo.VideoGenreMappingRepo,
	actorRepo *repo.ActorRepo, videoActorMappingRepo *repo.VideoActorMappingRepo,
	studioRepo *repo.StudioRepo, videoStudioMappingRepo *repo.VideoStudioMappingRepo,
	keywordRepo *repo.KeywordRepo, videoKeywordMappingRepo *repo.VideoKeywordMappingRepo,
	videoSubtitleMappingRepo *repo.VideoSubtitleMappingRepo) *SeriesService {
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
		cache:                    repo.NewCache[*sdomain.Series](c, seriesRepo.Data.Cache()),
	}
}

func (s *SeriesService) Get(ctx context.Context, id int64) (*sdomain.Series, error) {
	return s.cache.Get(ctx, cache_util.GetCacheActionName(id), func(action string, ctx context.Context) (*sdomain.Series, error) {
		return s.get(ctx, id)
	})
}

func (s *SeriesService) get(ctx context.Context, id int64) (*sdomain.Series, error) {
	item, err := s.seriesRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	genresMap, err := s.buildSeriesGenresMap(ctx, []*sdomain.Series{item})
	if err != nil {
		return nil, err
	}
	actorsMap, err := s.buildSeriesActorsMap(ctx, []*sdomain.Series{item})
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
	subtitlesMap, err := s.buildSeriesSubtitlesMap(ctx, []*sdomain.Series{item})
	if err != nil {
		return nil, err
	}
	item.Actors = actorsMap[item.ID]
	item.Genres = genresMap[item.ID]
	item.Keywords = keywordsMap[item.ID]
	item.Studios = studiosMap[item.ID]
	item.Subtitles = subtitlesMap[item.ID]
	return item, nil
}

func (s *SeriesService) Page(ctx context.Context, condition *sdomain.SearchSeries) (*sdomain.PageResult[*sdomain.Series], error) {
	rp, err := s.cache.Page(ctx, cache_util.GetCacheActionName(condition), func(action string, ctx context.Context) (*sdomain.PageResult[*sdomain.Series], error) {
		return s.page(ctx, condition)
	})
	if err != nil {
		return nil, err
	}
	return rp, nil
}

func (s *SeriesService) page(ctx context.Context, condition *sdomain.SearchSeries) (*sdomain.PageResult[*sdomain.Series], error) {
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
		actorsMap, err := s.buildSeriesActorsMap(ctx, pageResult.List)
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
		subtitlesMap, err := s.buildSeriesSubtitlesMap(ctx, pageResult.List)
		if err != nil {
			return nil, err
		}
		for _, item := range pageResult.List {
			item.Actors = actorsMap[item.ID]
			item.Genres = genresMap[item.ID]
			item.Subtitles = subtitlesMap[item.ID]
			item.Keywords = keywordsMap[item.ID]
			item.Studios = studiosMap[item.ID]
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
