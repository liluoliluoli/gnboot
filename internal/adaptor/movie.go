package adaptor

import (
	"context"
	"github.com/go-cinch/common/proto/params"
	"gnboot/api/movie"
	"gnboot/internal/service/sdomain"
	"gnboot/internal/utils/page_util"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *GnbootService) CreateMovie(ctx context.Context, req *movie.CreateMovieRequest) (*emptypb.Empty, error) {
	err := s.movie.Create(ctx, (&sdomain.CreateMovie{}).ConvertFromDto(req))
	return &emptypb.Empty{}, err
}

func (s *GnbootService) GetMovie(ctx context.Context, req *movie.GetMovieRequest) (*movie.GetMovieResp, error) {
	res, err := s.movie.Get(ctx, int64(req.Id))
	if err != nil {
		return nil, err
	}
	return res.ConvertFromDto(), nil
}

func (s *GnbootService) FindMovie(ctx context.Context, req *movie.FindMovieRequest) (*movie.FindMovieResp, error) {
	condition := &sdomain.FindMovie{
		Page:   lo.FromPtr(page_util.ToDomainPage(req.Page)),
		Search: req.Search,
		Sort: &sdomain.Sort{
			Filter: lo.TernaryF(req.Sort != nil, func() *string {
				return req.Sort.Filter
			}, func() *string {
				return nil
			}),
			Type: lo.TernaryF(req.Sort != nil, func() *string {
				return req.Sort.Type
			}, func() *string {
				return nil
			}),
			Direction: lo.TernaryF(req.Sort != nil, func() *string {
				return req.Sort.Direction
			}, func() *string {
				return nil
			}),
		},
	}
	res, err := s.movie.Page(ctx, condition)
	if err != nil {
		return nil, err
	}
	return &movie.FindMovieResp{
		Page: page_util.ToAdaptorPage(res.Page),
		List: lo.Map(res.List, func(item *sdomain.Movie, index int) *movie.MovieResp {
			return &movie.MovieResp{
				Id:            int64(item.ID),
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
			}
		}),
	}, nil
}

func (s *GnbootService) UpdateMovie(ctx context.Context, req *movie.UpdateMovieRequest) (*emptypb.Empty, error) {
	err := s.movie.Update(ctx, (&sdomain.UpdateMovie{}).ConvertFromDto(req))
	return &emptypb.Empty{}, err
}

func (s *GnbootService) DeleteMovie(ctx context.Context, req *params.IdsRequest) (*emptypb.Empty, error) {
	err := s.movie.Delete(ctx, lo.Map(strings.Split(req.Ids, ","), func(item string, index int) int64 {
		id, err := strconv.ParseInt(item, 10, 64)
		if err != nil {
			return 0
		}
		return id
	})...)
	return &emptypb.Empty{}, err
}
