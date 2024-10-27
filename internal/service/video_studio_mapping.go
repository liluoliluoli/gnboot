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

func (s *GnbootService) CreateVideoStudioMapping(ctx context.Context, req *gnboot.CreateVideoStudioMappingRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "CreateVideoStudioMapping")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.VideoStudioMapping{}
	copierx.Copy(&r, req)
	err = s.videoStudioMapping.Create(ctx, r)
	return
}

func (s *GnbootService) GetVideoStudioMapping(ctx context.Context, req *gnboot.GetVideoStudioMappingRequest) (rp *gnboot.GetVideoStudioMappingReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "GetVideoStudioMapping")
	defer span.End()
	rp = &gnboot.GetVideoStudioMappingReply{}
	res, err := s.videoStudioMapping.Get(ctx, req.Id)
	if err != nil {
		return
	}
	copierx.Copy(&rp, res)
	return
}

func (s *GnbootService) FindVideoStudioMapping(ctx context.Context, req *gnboot.FindVideoStudioMappingRequest) (rp *gnboot.FindVideoStudioMappingReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "FindVideoStudioMapping")
	defer span.End()
	rp = &gnboot.FindVideoStudioMappingReply{}
	rp.Page = &params.Page{}
	r := &biz.FindVideoStudioMapping{}
	r.Page = page.Page{}
	copierx.Copy(&r, req)
	copierx.Copy(&r.Page, req.Page)
	res, err := s.videoStudioMapping.Find(ctx, r)
	if err != nil {
		return
	}
	copierx.Copy(&rp.Page, r.Page)
	copierx.Copy(&rp.List, res)
	return
}

func (s *GnbootService) UpdateVideoStudioMapping(ctx context.Context, req *gnboot.UpdateVideoStudioMappingRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "UpdateVideoStudioMapping")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.UpdateVideoStudioMapping{}
	copierx.Copy(&r, req)
	err = s.videoStudioMapping.Update(ctx, r)
	return
}

func (s *GnbootService) DeleteVideoStudioMapping(ctx context.Context, req *params.IdsRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "DeleteVideoStudioMapping")
	defer span.End()
	rp = &emptypb.Empty{}
	err = s.videoStudioMapping.Delete(ctx, utils.Str2Uint64Arr(req.Ids)...)
	return
}
