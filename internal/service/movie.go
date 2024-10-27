package service

import (
	"context"

	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/page"
	"github.com/go-cinch/common/proto/params"
	"github.com/go-cinch/common/utils"
	"gnboot/api/gnboot"
	"gnboot/internal/biz"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *GnbootService) CreateMovie(ctx context.Context, req *gnboot.CreateMovieRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "CreateMovie")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.CreateMovie{}
	copierx.Copy(&r, req)
	err = s.movie.Create(ctx, r)
	return
}

func (s *GnbootService) GetMovie(ctx context.Context, req *gnboot.GetMovieRequest) (rp *gnboot.GetMovieReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "GetMovie")
	defer span.End()
	rp = &gnboot.GetMovieReply{}
	res, err := s.movie.Get(ctx, req.Id)
	if err != nil {
		return
	}
	copierx.Copy(&rp, res)
	return
}

func (s *GnbootService) FindMovie(ctx context.Context, req *gnboot.FindMovieRequest) (rp *gnboot.FindMovieReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "FindMovie")
	defer span.End()
	rp = &gnboot.FindMovieReply{}
	rp.Page = &params.Page{}
	r := &biz.FindMovie{}
	r.Page = page.Page{}
	copierx.Copy(&r, req)
	copierx.Copy(&r.Page, req.Page)
	res, err := s.movie.Find(ctx, r)
	if err != nil {
		return
	}
	copierx.Copy(&rp.Page, r.Page)
	copierx.Copy(&rp.List, res)
	return
}

func (s *GnbootService) UpdateMovie(ctx context.Context, req *gnboot.UpdateMovieRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "UpdateMovie")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.UpdateMovie{}
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
