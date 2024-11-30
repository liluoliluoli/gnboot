package service

import (
	"context"
	"github.com/liluoliluoli/gnboot/internal/common/constant"
	"github.com/liluoliluoli/gnboot/internal/common/utils/cache_util"
	"github.com/liluoliluoli/gnboot/internal/conf"
	"github.com/liluoliluoli/gnboot/internal/repo"
	"github.com/liluoliluoli/gnboot/internal/repo/gen"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/samber/lo"
	"time"
)

type MovieService struct {
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
	userRepo                 *repo.UserRepo
	videoUserMappingRepo     *repo.VideoUserMappingRepo
	cache                    sdomain.Cache[*sdomain.Movie]
}

func NewMovieService(c *conf.Bootstrap,
	movieRepo *repo.MovieRepo,
	genreRepo *repo.GenreRepo, videoGenreMappingRepo *repo.VideoGenreMappingRepo,
	actorRepo *repo.ActorRepo, videoActorMappingRepo *repo.VideoActorMappingRepo,
	studioRepo *repo.StudioRepo, videoStudioMappingRepo *repo.VideoStudioMappingRepo,
	keywordRepo *repo.KeywordRepo, videoKeywordMappingRepo *repo.VideoKeywordMappingRepo,
	videoSubtitleMappingRepo *repo.VideoSubtitleMappingRepo,
	userRepo *repo.UserRepo, videoUserMappingRepo *repo.VideoUserMappingRepo) *MovieService {
	return &MovieService{
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
		userRepo:                 userRepo,
		videoUserMappingRepo:     videoUserMappingRepo,
		cache:                    repo.NewCache[*sdomain.Movie](c, movieRepo.Data.Cache()),
	}
}

func (s *MovieService) Create(ctx context.Context, item *sdomain.CreateMovie) error {
	err := gen.Use(s.movieRepo.Data.DB(ctx)).Transaction(func(tx *gen.Query) error {
		err := s.cache.Flush(ctx, func(ctx context.Context) error {
			return s.movieRepo.Create(ctx, tx, item)
		})
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *MovieService) Get(ctx context.Context, id int64, userId int64) (*sdomain.Movie, error) {
	return s.cache.Get(ctx, cache_util.GetCacheActionName(id), func(action string, ctx context.Context) (*sdomain.Movie, error) {
		return s.get(ctx, id, userId)
	})
}

func (s *MovieService) get(ctx context.Context, id int64, userId int64) (*sdomain.Movie, error) {
	item, err := s.movieRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	genresMap, err := s.buildMovieGenresMap(ctx, []*sdomain.Movie{item})
	if err != nil {
		return nil, err
	}
	actorsMap, err := s.buildMovieActorsMap(ctx, []*sdomain.Movie{item})
	if err != nil {
		return nil, err
	}
	keywordsMap, err := s.buildMovieKeywordsMap(ctx, []*sdomain.Movie{item})
	if err != nil {
		return nil, err
	}
	studiosMap, err := s.buildMovieStudiosMap(ctx, []*sdomain.Movie{item})
	if err != nil {
		return nil, err
	}
	subtitlesMap, err := s.buildMovieSubtitlesMap(ctx, []*sdomain.Movie{item})
	if err != nil {
		return nil, err
	}
	if userId != 0 {
		moviePlayedMap, err := s.buildMovieLastPlayedInfo(ctx, userId, []*sdomain.Movie{item})
		if err != nil {
			return nil, err
		}
		item.LastPlayedPosition = lo.TernaryF(moviePlayedMap[item.ID] != nil, func() int32 {
			return moviePlayedMap[item.ID].LastPlayedPosition
		}, func() int32 {
			return 0
		})
		item.LastPlayedTime = lo.TernaryF(moviePlayedMap[item.ID] != nil, func() *time.Time {
			return lo.ToPtr(moviePlayedMap[item.ID].LastPlayedTime)
		}, func() *time.Time {
			return nil
		})
	}
	item.Actors = actorsMap[item.ID]
	item.Genres = genresMap[item.ID]
	item.Keywords = keywordsMap[item.ID]
	item.Studios = studiosMap[item.ID]
	item.Subtitles = subtitlesMap[item.ID]
	return item, nil
}

func (s *MovieService) Page(ctx context.Context, condition *sdomain.SearchMovie, userId int64) (*sdomain.PageResult[*sdomain.Movie], error) {
	rp, err := s.cache.Page(ctx, cache_util.GetCacheActionName(condition), func(action string, ctx context.Context) (*sdomain.PageResult[*sdomain.Movie], error) {
		return s.page(ctx, condition, userId)
	})
	if err != nil {
		return nil, err
	}
	return rp, nil
}

func (s *MovieService) page(ctx context.Context, condition *sdomain.SearchMovie, userId int64) (*sdomain.PageResult[*sdomain.Movie], error) {
	filterIds := make([]int64, 0)
	if condition.Id != 0 && condition.Type == "" {
		if condition.Type == constant.FilterType_genre {
			genreMappings, err := s.videoGenreMappingRepo.FindByGenreIdAndVideoType(ctx, condition.Id, constant.VideoType_movie)
			if err != nil {
				return nil, err
			}
			filterIds = append(filterIds, lo.Map(genreMappings, func(item *sdomain.VideoGenreMapping, index int) int64 {
				return item.VideoId
			})...)
		}
		if condition.Type == constant.FilterType_studio {
			studioMappings, err := s.videoStudioMappingRepo.FindByStudioIdAndVideoType(ctx, condition.Id, constant.VideoType_movie)
			if err != nil {
				return nil, err
			}
			filterIds = append(filterIds, lo.Map(studioMappings, func(item *sdomain.VideoStudioMapping, index int) int64 {
				return item.VideoId
			})...)
		}
		if condition.Type == constant.FilterType_keyword {
			keywordMappings, err := s.videoKeywordMappingRepo.FindByKeywordIdAndVideoType(ctx, condition.Id, constant.VideoType_movie)
			if err != nil {
				return nil, err
			}
			filterIds = append(filterIds, lo.Map(keywordMappings, func(item *sdomain.VideoKeywordMapping, index int) int64 {
				return item.VideoId
			})...)
		}
		if condition.Type == constant.FilterType_actor {
			actorMappings, err := s.videoActorMappingRepo.FindByActorIdAndVideoType(ctx, condition.Id, constant.VideoType_movie)
			if err != nil {
				return nil, err
			}
			filterIds = append(filterIds, lo.Map(actorMappings, func(item *sdomain.VideoActorMapping, index int) int64 {
				return item.VideoId
			})...)
		}
	}
	condition.FilterIds = filterIds
	pageResult, err := s.movieRepo.Page(ctx, condition)
	if err != nil {
		return nil, err
	}
	if pageResult != nil && len(pageResult.List) != 0 {
		genresMap, err := s.buildMovieGenresMap(ctx, pageResult.List)
		if err != nil {
			return nil, err
		}
		actorsMap, err := s.buildMovieActorsMap(ctx, pageResult.List)
		if err != nil {
			return nil, err
		}
		keywordsMap, err := s.buildMovieKeywordsMap(ctx, pageResult.List)
		if err != nil {
			return nil, err
		}
		studiosMap, err := s.buildMovieStudiosMap(ctx, pageResult.List)
		if err != nil {
			return nil, err
		}
		subtitlesMap, err := s.buildMovieSubtitlesMap(ctx, pageResult.List)
		if err != nil {
			return nil, err
		}

		var moviePlayedMap map[int64]*sdomain.VideoUserMapping
		if userId != 0 {
			moviePlayedMap, err = s.buildMovieLastPlayedInfo(ctx, userId, pageResult.List)
			if err != nil {
				return nil, err
			}
		}
		for _, item := range pageResult.List {
			item.Actors = actorsMap[item.ID]
			item.Genres = genresMap[item.ID]
			item.Subtitles = subtitlesMap[item.ID]
			item.Keywords = keywordsMap[item.ID]
			item.Studios = studiosMap[item.ID]
			if userId != 0 {
				item.LastPlayedPosition = lo.TernaryF(moviePlayedMap[item.ID] != nil, func() int32 {
					return moviePlayedMap[item.ID].LastPlayedPosition
				}, func() int32 {
					return 0
				})
				item.LastPlayedTime = lo.TernaryF(moviePlayedMap[item.ID] != nil, func() *time.Time {
					return lo.ToPtr(moviePlayedMap[item.ID].LastPlayedTime)
				}, func() *time.Time {
					return nil
				})
			}
		}
	}
	return pageResult, nil
}

func (s *MovieService) Update(ctx context.Context, item *sdomain.UpdateMovie) error {
	err := gen.Use(s.movieRepo.Data.DB(ctx)).Transaction(func(tx *gen.Query) error {
		err := s.cache.Flush(ctx, func(ctx context.Context) error {
			return s.movieRepo.Update(ctx, tx, item)
		})
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *MovieService) Delete(ctx context.Context, ids ...int64) error {
	err := gen.Use(s.movieRepo.Data.DB(ctx)).Transaction(func(tx *gen.Query) error {
		err := s.cache.Flush(ctx, func(ctx context.Context) error {
			return s.movieRepo.Delete(ctx, tx, ids...)
		})
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// 补齐genre
func (s *MovieService) buildMovieGenresMap(ctx context.Context, movies []*sdomain.Movie) (map[int64][]*sdomain.Genre, error) {
	movieIds := lo.Map(movies, func(item *sdomain.Movie, index int) int64 {
		return item.ID
	})
	genreMappings, err := s.videoGenreMappingRepo.FindByVideoIdAndType(ctx, movieIds, constant.VideoType_movie)
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
func (s *MovieService) buildMovieActorsMap(ctx context.Context, movies []*sdomain.Movie) (map[int64][]*sdomain.Actor, error) {
	movieIds := lo.Map(movies, func(item *sdomain.Movie, index int) int64 {
		return item.ID
	})
	actorMappings, err := s.videoActorMappingRepo.FindByVideoIdAndType(ctx, movieIds, constant.VideoType_movie)
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
func (s *MovieService) buildMovieKeywordsMap(ctx context.Context, movies []*sdomain.Movie) (map[int64][]*sdomain.Keyword, error) {
	movieIds := lo.Map(movies, func(item *sdomain.Movie, index int) int64 {
		return item.ID
	})
	keywordMappings, err := s.videoKeywordMappingRepo.FindByVideoIdAndType(ctx, movieIds, constant.VideoType_movie)
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
func (s *MovieService) buildMovieStudiosMap(ctx context.Context, movies []*sdomain.Movie) (map[int64][]*sdomain.Studio, error) {
	movieIds := lo.Map(movies, func(item *sdomain.Movie, index int) int64 {
		return item.ID
	})
	studioMappings, err := s.videoStudioMappingRepo.FindByVideoIdAndType(ctx, movieIds, constant.VideoType_movie)
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
func (s *MovieService) buildMovieSubtitlesMap(ctx context.Context, movies []*sdomain.Movie) (map[int64][]*sdomain.VideoSubtitleMapping, error) {
	movieIds := lo.Map(movies, func(item *sdomain.Movie, index int) int64 {
		return item.ID
	})
	subtitleMappings, err := s.videoSubtitleMappingRepo.FindByVideoIdAndType(ctx, movieIds, constant.VideoType_movie)
	if err != nil {
		return nil, err
	}
	return lo.GroupBy(subtitleMappings, func(item *sdomain.VideoSubtitleMapping) int64 {
		return item.VideoId
	}), nil
}

// 补齐最后播放信息
func (s *MovieService) buildMovieLastPlayedInfo(ctx context.Context, userId int64, movies []*sdomain.Movie) (map[int64]*sdomain.VideoUserMapping, error) {
	movieIds := lo.Map(movies, func(item *sdomain.Movie, index int) int64 {
		return item.ID
	})
	userMappings, err := s.videoUserMappingRepo.FindByUserIdAndVideoIdAndType(ctx, userId, movieIds, constant.VideoType_movie)
	if err != nil {
		return nil, err
	}
	return lo.SliceToMap(userMappings, func(item *sdomain.VideoUserMapping) (int64, *sdomain.VideoUserMapping) {
		return item.VideoId, item
	}), nil
}
