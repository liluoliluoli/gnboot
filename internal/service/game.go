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

func (s *GnbootService) CreateGnboot(ctx context.Context, req *gnboot.CreateGnbootRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "CreateGnboot")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.CreateGnboot{}
	copierx.Copy(&r, req)
	err = s.gnboot.Create(ctx, r)
	return
}

func (s *GnbootService) GetGnboot(ctx context.Context, req *gnboot.GetGnbootRequest) (rp *gnboot.GetGnbootReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "GetGnboot")
	defer span.End()
	rp = &gnboot.GetGnbootReply{}
	res, err := s.gnboot.Get(ctx, req.Id)
	if err != nil {
		return
	}
	copierx.Copy(&rp, res)
	return
}

func (s *GnbootService) FindGnboot(ctx context.Context, req *gnboot.FindGnbootRequest) (rp *gnboot.FindGnbootReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "FindGnboot")
	defer span.End()
	rp = &gnboot.FindGnbootReply{}
	rp.Page = &params.Page{}
	r := &biz.FindGnboot{}
	r.Page = page.Page{}
	copierx.Copy(&r, req)
	copierx.Copy(&r.Page, req.Page)
	res, err := s.gnboot.Find(ctx, r)
	if err != nil {
		return
	}
	copierx.Copy(&rp.Page, r.Page)
	copierx.Copy(&rp.List, res)
	return
}

func (s *GnbootService) UpdateGnboot(ctx context.Context, req *gnboot.UpdateGnbootRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "UpdateGnboot")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.UpdateGnboot{}
	copierx.Copy(&r, req)
	err = s.gnboot.Update(ctx, r)
	return
}

func (s *GnbootService) DeleteGnboot(ctx context.Context, req *params.IdsRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "DeleteGnboot")
	defer span.End()
	rp = &emptypb.Empty{}
	err = s.gnboot.Delete(ctx, utils.Str2Uint64Arr(req.Ids)...)
	return
}
