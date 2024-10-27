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

func (s *GnbootService) CreateVideoUserMapping(ctx context.Context, req *gnboot.CreateVideoUserMappingRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "CreateVideoUserMapping")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.VideoUserMapping{}
	copierx.Copy(&r, req)
	err = s.videoUserMapping.Create(ctx, r)
	return
}

func (s *GnbootService) GetVideoUserMapping(ctx context.Context, req *gnboot.GetVideoUserMappingRequest) (rp *gnboot.GetVideoUserMappingReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "GetVideoUserMapping")
	defer span.End()
	rp = &gnboot.GetVideoUserMappingReply{}
	res, err := s.videoUserMapping.Get(ctx, req.Id)
	if err != nil {
		return
	}
	copierx.Copy(&rp, res)
	return
}

func (s *GnbootService) FindVideoUserMapping(ctx context.Context, req *gnboot.FindVideoUserMappingRequest) (rp *gnboot.FindVideoUserMappingReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "FindVideoUserMapping")
	defer span.End()
	rp = &gnboot.FindVideoUserMappingReply{}
	rp.Page = &params.Page{}
	r := &biz.FindVideoUserMapping{}
	r.Page = page.Page{}
	copierx.Copy(&r, req)
	copierx.Copy(&r.Page, req.Page)
	res, err := s.videoUserMapping.Find(ctx, r)
	if err != nil {
		return
	}
	copierx.Copy(&rp.Page, r.Page)
	copierx.Copy(&rp.List, res)
	return
}

func (s *GnbootService) UpdateVideoUserMapping(ctx context.Context, req *gnboot.UpdateVideoUserMappingRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "UpdateVideoUserMapping")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.UpdateVideoUserMapping{}
	copierx.Copy(&r, req)
	err = s.videoUserMapping.Update(ctx, r)
	return
}

func (s *GnbootService) DeleteVideoUserMapping(ctx context.Context, req *params.IdsRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "DeleteVideoUserMapping")
	defer span.End()
	rp = &emptypb.Empty{}
	err = s.videoUserMapping.Delete(ctx, utils.Str2Uint64Arr(req.Ids)...)
	return
}
