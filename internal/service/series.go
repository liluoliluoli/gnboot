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

func (s *GnbootService) CreateSeries(ctx context.Context, req *gnboot.CreateSeriesRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "CreateSeries")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.Series{}
	copierx.Copy(&r, req)
	err = s.series.Create(ctx, r)
	return
}

func (s *GnbootService) GetSeries(ctx context.Context, req *gnboot.GetSeriesRequest) (rp *gnboot.GetSeriesReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "GetSeries")
	defer span.End()
	rp = &gnboot.GetSeriesReply{}
	res, err := s.series.Get(ctx, req.Id)
	if err != nil {
		return
	}
	copierx.Copy(&rp, res)
	return
}

func (s *GnbootService) FindSeries(ctx context.Context, req *gnboot.FindSeriesRequest) (rp *gnboot.FindSeriesReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "FindSeries")
	defer span.End()
	rp = &gnboot.FindSeriesReply{}
	rp.Page = &params.Page{}
	r := &biz.FindSeries{}
	r.Page = page.Page{}
	copierx.Copy(&r, req)
	copierx.Copy(&r.Page, req.Page)
	res, err := s.series.Find(ctx, r)
	if err != nil {
		return
	}
	copierx.Copy(&rp.Page, r.Page)
	copierx.Copy(&rp.List, res)
	return
}

func (s *GnbootService) UpdateSeries(ctx context.Context, req *gnboot.UpdateSeriesRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "UpdateSeries")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.UpdateSeries{}
	copierx.Copy(&r, req)
	err = s.series.Update(ctx, r)
	return
}

func (s *GnbootService) DeleteSeries(ctx context.Context, req *params.IdsRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "DeleteSeries")
	defer span.End()
	rp = &emptypb.Empty{}
	err = s.series.Delete(ctx, utils.Str2Uint64Arr(req.Ids)...)
	return
}
