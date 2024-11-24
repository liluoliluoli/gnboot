package service

import (
	"context"
	"github.com/samber/lo"
	"gnboot/internal/common/constant"
	"gnboot/internal/conf"
	"gnboot/internal/repo"
	"gnboot/internal/repo/gen"
	"gnboot/internal/service/sdomain"
	"gnboot/internal/utils/cache_util"
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
	cache                    sdomain.Cache[*sdomain.Movie]
}

func NewMovieService(c *conf.Bootstrap,
	movieRepo *repo.MovieRepo,
	genreRepo *repo.GenreRepo, videoGenreMappingRepo *repo.VideoGenreMappingRepo,
	actorRepo *repo.ActorRepo, videoActorMappingRepo *repo.VideoActorMappingRepo,
	studioRepo *repo.StudioRepo, videoStudioMappingRepo *repo.VideoStudioMappingRepo,
	keywordRepo *repo.KeywordRepo, videoKeywordMappingRepo *repo.VideoKeywordMappingRepo,
	videoSubtitleMappingRepo *repo.VideoSubtitleMappingRepo) *MovieService {
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

func (s *MovieService) Get(ctx context.Context, id int64) (*sdomain.Movie, error) {
	return s.cache.Get(ctx, cache_util.GetCacheActionName(id), func(action string, ctx context.Context) (*sdomain.Movie, error) {
		return s.get(ctx, id)
	})
}

func (s *MovieService) get(ctx context.Context, id int64) (*sdomain.Movie, error) {
	item, err := s.movieRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (s *MovieService) Page(ctx context.Context, condition *sdomain.FindMovie) (*sdomain.PageResult[*sdomain.Movie], error) {
	rp, err := s.cache.GetPage(ctx, cache_util.GetCacheActionName(condition), func(action string, ctx context.Context) (*sdomain.PageResult[*sdomain.Movie], error) {
		return s.page(ctx, condition)
	})
	if err != nil {
		return nil, err
	}
	return rp, nil
}

func (s *MovieService) page(ctx context.Context, condition *sdomain.FindMovie) (*sdomain.PageResult[*sdomain.Movie], error) {
	pageResult, err := s.movieRepo.Page(ctx, condition)
	if err != nil {
		return nil, err
	}
	if pageResult != nil && len(pageResult.List) != 0 {
		for _, item := range pageResult.List {
			//补齐genre
			genreMappings, err := s.videoGenreMappingRepo.FindByVideoIdAndType(ctx, lo.Map(pageResult.List, func(item *sdomain.Movie, index int) int64 {
				return item.ID
			}), constant.VideoType_movie)
			if err != nil {
				return nil, err
			}
			genres, err := s.genreRepo.FindByIds(ctx, lo.Map(genreMappings, func(item *sdomain.VideoGenreMapping, index int) int64 {
				return item.GenreId
			}))
			if err != nil {
				return nil, err
			}
			item.Genres = genres
			//补齐actor
			actorMappings, err := s.videoActorMappingRepo.FindByVideoIdAndType(ctx, lo.Map(pageResult.List, func(item *sdomain.Movie, index int) int64 {
				return item.ID
			}), constant.VideoType_movie)
			if err != nil {
				return nil, err
			}
			actors, err := s.actorRepo.FindByIds(ctx, lo.Map(actorMappings, func(item *sdomain.VideoActorMapping, index int) int64 {
				return item.ActorId
			}))
			if err != nil {
				return nil, err
			}
			item.Actors = actors
			//补齐keyword
			keywordMappings, err := s.videoKeywordMappingRepo.FindByVideoIdAndType(ctx, lo.Map(pageResult.List, func(item *sdomain.Movie, index int) int64 {
				return item.ID
			}), constant.VideoType_movie)
			if err != nil {
				return nil, err
			}
			keywords, err := s.keywordRepo.FindByIds(ctx, lo.Map(keywordMappings, func(item *sdomain.VideoKeywordMapping, index int) int64 {
				return item.KeywordId
			}))
			if err != nil {
				return nil, err
			}
			item.Keywords = keywords
			//补齐studio
			studioMappings, err := s.videoStudioMappingRepo.FindByVideoIdAndType(ctx, lo.Map(pageResult.List, func(item *sdomain.Movie, index int) int64 {
				return item.ID
			}), constant.VideoType_movie)
			if err != nil {
				return nil, err
			}
			studios, err := s.studioRepo.FindByIds(ctx, lo.Map(studioMappings, func(item *sdomain.VideoStudioMapping, index int) int64 {
				return item.StudioId
			}))
			if err != nil {
				return nil, err
			}
			item.Studios = studios
			//补齐subtitle
			subtitleMappings, err := s.videoSubtitleMappingRepo.FindByVideoIdAndType(ctx, lo.Map(pageResult.List, func(item *sdomain.Movie, index int) int64 {
				return item.ID
			}), constant.VideoType_movie)
			if err != nil {
				return nil, err
			}
			item.Subtitles = subtitleMappings
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
