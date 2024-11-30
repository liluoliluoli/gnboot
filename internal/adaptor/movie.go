package adaptor

import (
	"context"
	"github.com/liluoliluoli/gnboot/api"
	"github.com/liluoliluoli/gnboot/api/movie"
	"github.com/liluoliluoli/gnboot/internal/common/utils/page_util"
	"github.com/liluoliluoli/gnboot/internal/service"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
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
		FilterByNextPlay: false,
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
		Page:             page_util.ToDomainPage(req.Page),
		Id:               req.Id,
		Type:             req.Type,
		FilterByNextPlay: false,
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

func (s *MovieProvider) NextToPlayMovies(ctx context.Context, req *movie.NextToPlayMoviesRequest) (*movie.SearchMovieResp, error) {
	condition := &sdomain.SearchMovie{
		Page:             page_util.ToDomainPage(req.Page),
		FilterByNextPlay: true,
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
