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

func (s *GnbootService) CreateVideoActorMapping(ctx context.Context, req *gnboot.CreateVideoActorMappingRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "CreateVideoActorMapping")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.VideoActorMapping{}
	copierx.Copy(&r, req)
	err = s.videoActorMapping.Create(ctx, r)
	return
}

func (s *GnbootService) GetVideoActorMapping(ctx context.Context, req *gnboot.GetVideoActorMappingRequest) (rp *gnboot.GetVideoActorMappingReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "GetVideoActorMapping")
	defer span.End()
	rp = &gnboot.GetVideoActorMappingReply{}
	res, err := s.videoActorMapping.Get(ctx, req.Id)
	if err != nil {
		return
	}
	copierx.Copy(&rp, res)
	return
}

func (s *GnbootService) FindVideoActorMapping(ctx context.Context, req *gnboot.FindVideoActorMappingRequest) (rp *gnboot.FindVideoActorMappingReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "FindVideoActorMapping")
	defer span.End()
	rp = &gnboot.FindVideoActorMappingReply{}
	rp.Page = &params.Page{}
	r := &biz.FindVideoActorMapping{}
	r.Page = page.Page{}
	copierx.Copy(&r, req)
	copierx.Copy(&r.Page, req.Page)
	res, err := s.videoActorMapping.Find(ctx, r)
	if err != nil {
		return
	}
	copierx.Copy(&rp.Page, r.Page)
	copierx.Copy(&rp.List, res)
	return
}

func (s *GnbootService) UpdateVideoActorMapping(ctx context.Context, req *gnboot.UpdateVideoActorMappingRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "UpdateVideoActorMapping")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.UpdateVideoActorMapping{}
	copierx.Copy(&r, req)
	err = s.videoActorMapping.Update(ctx, r)
	return
}

func (s *GnbootService) DeleteVideoActorMapping(ctx context.Context, req *params.IdsRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "DeleteVideoActorMapping")
	defer span.End()
	rp = &emptypb.Empty{}
	err = s.videoActorMapping.Delete(ctx, utils.Str2Uint64Arr(req.Ids)...)
	return
}
