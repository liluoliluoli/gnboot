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

func (s *GnbootService) CreateVideoSubtitleMapping(ctx context.Context, req *gnboot.CreateVideoSubtitleMappingRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "CreateVideoSubtitleMapping")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.VideoSubtitleMapping{}
	copierx.Copy(&r, req)
	err = s.videoSubtitleMapping.Create(ctx, r)
	return
}

func (s *GnbootService) GetVideoSubtitleMapping(ctx context.Context, req *gnboot.GetVideoSubtitleMappingRequest) (rp *gnboot.GetVideoSubtitleMappingReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "GetVideoSubtitleMapping")
	defer span.End()
	rp = &gnboot.GetVideoSubtitleMappingReply{}
	res, err := s.videoSubtitleMapping.Get(ctx, req.Id)
	if err != nil {
		return
	}
	copierx.Copy(&rp, res)
	return
}

func (s *GnbootService) FindVideoSubtitleMapping(ctx context.Context, req *gnboot.FindVideoSubtitleMappingRequest) (rp *gnboot.FindVideoSubtitleMappingReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "FindVideoSubtitleMapping")
	defer span.End()
	rp = &gnboot.FindVideoSubtitleMappingReply{}
	rp.Page = &params.Page{}
	r := &biz.FindVideoSubtitleMapping{}
	r.Page = page.Page{}
	copierx.Copy(&r, req)
	copierx.Copy(&r.Page, req.Page)
	res, err := s.videoSubtitleMapping.Find(ctx, r)
	if err != nil {
		return
	}
	copierx.Copy(&rp.Page, r.Page)
	copierx.Copy(&rp.List, res)
	return
}

func (s *GnbootService) UpdateVideoSubtitleMapping(ctx context.Context, req *gnboot.UpdateVideoSubtitleMappingRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "UpdateVideoSubtitleMapping")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.UpdateVideoSubtitleMapping{}
	copierx.Copy(&r, req)
	err = s.videoSubtitleMapping.Update(ctx, r)
	return
}

func (s *GnbootService) DeleteVideoSubtitleMapping(ctx context.Context, req *params.IdsRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "DeleteVideoSubtitleMapping")
	defer span.End()
	rp = &emptypb.Empty{}
	err = s.videoSubtitleMapping.Delete(ctx, utils.Str2Uint64Arr(req.Ids)...)
	return
}
