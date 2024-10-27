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

func (s *GnbootService) CreateStudio(ctx context.Context, req *gnboot.CreateStudioRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "CreateStudio")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.Studio{}
	copierx.Copy(&r, req)
	err = s.studio.Create(ctx, r)
	return
}

func (s *GnbootService) GetStudio(ctx context.Context, req *gnboot.GetStudioRequest) (rp *gnboot.GetStudioReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "GetStudio")
	defer span.End()
	rp = &gnboot.GetStudioReply{}
	res, err := s.studio.Get(ctx, req.Id)
	if err != nil {
		return
	}
	copierx.Copy(&rp, res)
	return
}

func (s *GnbootService) FindStudio(ctx context.Context, req *gnboot.FindStudioRequest) (rp *gnboot.FindStudioReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "FindStudio")
	defer span.End()
	rp = &gnboot.FindStudioReply{}
	rp.Page = &params.Page{}
	r := &biz.FindStudio{}
	r.Page = page.Page{}
	copierx.Copy(&r, req)
	copierx.Copy(&r.Page, req.Page)
	res, err := s.studio.Find(ctx, r)
	if err != nil {
		return
	}
	copierx.Copy(&rp.Page, r.Page)
	copierx.Copy(&rp.List, res)
	return
}

func (s *GnbootService) UpdateStudio(ctx context.Context, req *gnboot.UpdateStudioRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "UpdateStudio")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.UpdateStudio{}
	copierx.Copy(&r, req)
	err = s.studio.Update(ctx, r)
	return
}

func (s *GnbootService) DeleteStudio(ctx context.Context, req *params.IdsRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "DeleteStudio")
	defer span.End()
	rp = &emptypb.Empty{}
	err = s.studio.Delete(ctx, utils.Str2Uint64Arr(req.Ids)...)
	return
}
