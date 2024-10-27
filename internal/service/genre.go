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

func (s *GnbootService) CreateGenre(ctx context.Context, req *gnboot.CreateGenreRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "CreateGenre")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.Genre{}
	copierx.Copy(&r, req)
	err = s.genre.Create(ctx, r)
	return
}

func (s *GnbootService) GetGenre(ctx context.Context, req *gnboot.GetGenreRequest) (rp *gnboot.GetGenreReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "GetGenre")
	defer span.End()
	rp = &gnboot.GetGenreReply{}
	res, err := s.genre.Get(ctx, req.Id)
	if err != nil {
		return
	}
	copierx.Copy(&rp, res)
	return
}

func (s *GnbootService) FindGenre(ctx context.Context, req *gnboot.FindGenreRequest) (rp *gnboot.FindGenreReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "FindGenre")
	defer span.End()
	rp = &gnboot.FindGenreReply{}
	rp.Page = &params.Page{}
	r := &biz.FindGenre{}
	r.Page = page.Page{}
	copierx.Copy(&r, req)
	copierx.Copy(&r.Page, req.Page)
	res, err := s.genre.Find(ctx, r)
	if err != nil {
		return
	}
	copierx.Copy(&rp.Page, r.Page)
	copierx.Copy(&rp.List, res)
	return
}

func (s *GnbootService) UpdateGenre(ctx context.Context, req *gnboot.UpdateGenreRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "UpdateGenre")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.UpdateGenre{}
	copierx.Copy(&r, req)
	err = s.genre.Update(ctx, r)
	return
}

func (s *GnbootService) DeleteGenre(ctx context.Context, req *params.IdsRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "DeleteGenre")
	defer span.End()
	rp = &emptypb.Empty{}
	err = s.genre.Delete(ctx, utils.Str2Uint64Arr(req.Ids)...)
	return
}
