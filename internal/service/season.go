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

func (s *GnbootService) CreateSeason(ctx context.Context, req *gnboot.CreateSeasonRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "CreateSeason")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.Season{}
	copierx.Copy(&r, req)
	err = s.season.Create(ctx, r)
	return
}

func (s *GnbootService) GetSeason(ctx context.Context, req *gnboot.GetSeasonRequest) (rp *gnboot.GetSeasonReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "GetSeason")
	defer span.End()
	rp = &gnboot.GetSeasonReply{}
	res, err := s.season.Get(ctx, req.Id)
	if err != nil {
		return
	}
	copierx.Copy(&rp, res)
	return
}

func (s *GnbootService) FindSeason(ctx context.Context, req *gnboot.FindSeasonRequest) (rp *gnboot.FindSeasonReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "FindSeason")
	defer span.End()
	rp = &gnboot.FindSeasonReply{}
	rp.Page = &params.Page{}
	r := &biz.FindSeason{}
	r.Page = page.Page{}
	copierx.Copy(&r, req)
	copierx.Copy(&r.Page, req.Page)
	res, err := s.season.Find(ctx, r)
	if err != nil {
		return
	}
	copierx.Copy(&rp.Page, r.Page)
	copierx.Copy(&rp.List, res)
	return
}

func (s *GnbootService) UpdateSeason(ctx context.Context, req *gnboot.UpdateSeasonRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "UpdateSeason")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.UpdateSeason{}
	copierx.Copy(&r, req)
	err = s.season.Update(ctx, r)
	return
}

func (s *GnbootService) DeleteSeason(ctx context.Context, req *params.IdsRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "DeleteSeason")
	defer span.End()
	rp = &emptypb.Empty{}
	err = s.season.Delete(ctx, utils.Str2Uint64Arr(req.Ids)...)
	return
}
