package adaptor

import (
	"context"
	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/proto/params"
	"github.com/go-cinch/common/utils"
	"gnboot/api/movie"
	"gnboot/internal/utils/page_util"

	"github.com/samber/lo"
	"gnboot/internal/service"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *GnbootService) CreateMovie(ctx context.Context, req *movie.CreateMovieRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "CreateMovie")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &service.CreateMovie{}
	copierx.Copy(&r, req)
	err = s.movie.Create(ctx, r)
	return
}

func (s *GnbootService) GetMovie(ctx context.Context, req *movie.GetMovieRequest) (rp *movie.GetMovieResp, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "GetMovie")
	defer span.End()
	rp = &movie.GetMovieResp{}
	res, err := s.movie.Get(ctx, req.Id)
	if err != nil {
		return
	}
	copierx.Copy(&rp, res)
	return
}

func (s *GnbootService) FindMovie(ctx context.Context, req *movie.FindMovieRequest) (*movie.FindMovieResp, error) {
	condition := &service.FindMovie{
		Page:   lo.FromPtr(page_util.ToDomainPage(req.Page)),
		Search: req.Search,
		Sort: &service.Sort{
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
	res, err := s.movie.Find(ctx, condition)
	if err != nil {
		return nil, err
	}
	return &movie.FindMovieResp{
		Page: page_util.ToAdaptorPage(condition.Page),
		List: lo.Map(res, func(item *service.Movie, index int) *movie.MovieResp {
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
			}
		}),
	}, nil
}

func (s *GnbootService) UpdateMovie(ctx context.Context, req *movie.UpdateMovieRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "UpdateMovie")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &service.UpdateMovie{}
	copierx.Copy(&r, req)
	err = s.movie.Update(ctx, r)
	return
}

func (s *GnbootService) DeleteMovie(ctx context.Context, req *params.IdsRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "DeleteMovie")
	defer span.End()
	rp = &emptypb.Empty{}
	err = s.movie.Delete(ctx, utils.Str2Uint64Arr(req.Ids)...)
	return
}
