package adaptor

import (
	"context"
	"github.com/liluoliluoli/gnboot/api"
	"github.com/liluoliluoli/gnboot/api/actor"
	"github.com/liluoliluoli/gnboot/api/genre"
	"github.com/liluoliluoli/gnboot/api/keyword"
	"github.com/liluoliluoli/gnboot/api/movie"
	"github.com/liluoliluoli/gnboot/api/studio"
	"github.com/liluoliluoli/gnboot/api/subtitle"
	"github.com/liluoliluoli/gnboot/internal/service"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/liluoliluoli/gnboot/internal/utils/page_util"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MovieProvider struct {
	movie.UnimplementedMovieRemoteServiceServer
	movie *service.MovieService
}

func NewMovieProvider(movie *service.MovieService) *MovieProvider {
	return &MovieProvider{movie: movie}
}

func (s *MovieProvider) CreateMovie(ctx context.Context, req *movie.CreateMovieRequest) (*emptypb.Empty, error) {
	err := s.movie.Create(ctx, (&sdomain.CreateMovie{}).ConvertFromDto(req))
	return &emptypb.Empty{}, err
}

func (s *MovieProvider) GetMovie(ctx context.Context, req *movie.GetMovieRequest) (*movie.MovieResp, error) {
	res, err := s.movie.Get(ctx, req.Id, 1)
	if err != nil {
		return nil, err
	}
	return res.ConvertToDto(), nil
}

func (s *MovieProvider) FindMovie(ctx context.Context, req *movie.FindMovieRequest) (*movie.SearchMovieResp, error) {
	condition := &sdomain.SearchMovie{
		Page:   page_util.ToDomainPage(req.Page),
		Search: lo.FromPtr(req.Search),
		Sort: &sdomain.Sort{
			Filter: lo.TernaryF(req.Sort != nil, func() string {
				return lo.FromPtr(req.Sort.Filter)
			}, func() string {
				return ""
			}),
			Type: lo.TernaryF(req.Sort != nil, func() string {
				return lo.FromPtr(req.Sort.Type)
			}, func() string {
				return ""
			}),
			Direction: lo.TernaryF(req.Sort != nil, func() string {
				return lo.FromPtr(req.Sort.Direction)
			}, func() string {
				return ""
			}),
		},
	}
	res, err := s.movie.Page(ctx, condition, 1)
	if err != nil {
		return nil, err
	}
	return &movie.SearchMovieResp{
		Page: page_util.ToAdaptorPage(res.Page),
		List: lo.Map(res.List, func(item *sdomain.Movie, index int) *movie.MovieResp {
			return item.ConvertToDto()
		}),
	}, nil
}

func (s *MovieProvider) FilterMovie(ctx context.Context, req *movie.FilterMovieRequest) (*movie.SearchMovieResp, error) {
	condition := &sdomain.SearchMovie{
		Page: page_util.ToDomainPage(req.Page),
		Id:   req.Id,
		Type: req.Type,
	}
	res, err := s.movie.Page(ctx, condition, 1)
	if err != nil {
		return nil, err
	}
	return &movie.SearchMovieResp{
		Page: page_util.ToAdaptorPage(res.Page),
		List: lo.Map(res.List, func(item *sdomain.Movie, index int) *movie.MovieResp {
			return &movie.MovieResp{
				Id:            item.ID,
				OriginalTitle: item.OriginalTitle,
				Status:        item.Status,
				VoteAverage:   item.VoteAverage,
				VoteCount:     item.VoteCount,
				Country:       item.Country,
				Trailer:       item.Trailer,
				Url:           item.URL,
				Downloaded:    item.Downloaded,
				FileSize:      item.FileSize,
				Filename:      item.Filename,
				Ext:           item.Ext,
				Genres: lo.Map(item.Genres, func(item *sdomain.Genre, index int) *genre.GenreResp {
					return item.ConvertToDto()
				}),
				Studios: lo.Map(item.Studios, func(item *sdomain.Studio, index int) *studio.StudioResp {
					return item.ConvertToDto()
				}),
				Keywords: lo.Map(item.Keywords, func(item *sdomain.Keyword, index int) *keyword.KeywordResp {
					return item.ConvertToDto()
				}),
				Actors: lo.Map(item.Actors, func(item *sdomain.Actor, index int) *actor.ActorResp {
					return item.ConvertToDto()
				}),
				Subtitles: lo.Map(item.Subtitles, func(item *sdomain.VideoSubtitleMapping, index int) *subtitle.SubtitleResp {
					return item.ConvertToDto()
				}),
			}
		}),
	}, nil
}

func (s *MovieProvider) UpdateMovie(ctx context.Context, req *movie.UpdateMovieRequest) (*emptypb.Empty, error) {
	err := s.movie.Update(ctx, (&sdomain.UpdateMovie{}).ConvertFromDto(req))
	return &emptypb.Empty{}, err
}

func (s *MovieProvider) DeleteMovie(ctx context.Context, req *api.IdsRequest) (*emptypb.Empty, error) {
	err := s.movie.Delete(ctx, lo.Map(strings.Split(req.Ids, ","), func(item string, index int) int64 {
		id, err := strconv.ParseInt(item, 10, 64)
		if err != nil {
			return 0
		}
		return id
	})...)
	return &emptypb.Empty{}, err
}
