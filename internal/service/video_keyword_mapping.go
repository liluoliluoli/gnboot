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

func (s *GnbootService) CreateVideoKeywordMapping(ctx context.Context, req *gnboot.CreateVideoKeywordMappingRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "CreateVideoKeywordMapping")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.VideoKeywordMapping{}
	copierx.Copy(&r, req)
	err = s.videoKeywordMapping.Create(ctx, r)
	return
}

func (s *GnbootService) GetVideoKeywordMapping(ctx context.Context, req *gnboot.GetVideoKeywordMappingRequest) (rp *gnboot.GetVideoKeywordMappingReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "GetVideoKeywordMapping")
	defer span.End()
	rp = &gnboot.GetVideoKeywordMappingReply{}
	res, err := s.videoKeywordMapping.Get(ctx, req.Id)
	if err != nil {
		return
	}
	copierx.Copy(&rp, res)
	return
}

func (s *GnbootService) FindVideoKeywordMapping(ctx context.Context, req *gnboot.FindVideoKeywordMappingRequest) (rp *gnboot.FindVideoKeywordMappingReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "FindVideoKeywordMapping")
	defer span.End()
	rp = &gnboot.FindVideoKeywordMappingReply{}
	rp.Page = &params.Page{}
	r := &biz.FindVideoKeywordMapping{}
	r.Page = page.Page{}
	copierx.Copy(&r, req)
	copierx.Copy(&r.Page, req.Page)
	res, err := s.videoKeywordMapping.Find(ctx, r)
	if err != nil {
		return
	}
	copierx.Copy(&rp.Page, r.Page)
	copierx.Copy(&rp.List, res)
	return
}

func (s *GnbootService) UpdateVideoKeywordMapping(ctx context.Context, req *gnboot.UpdateVideoKeywordMappingRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "UpdateVideoKeywordMapping")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.UpdateVideoKeywordMapping{}
	copierx.Copy(&r, req)
	err = s.videoKeywordMapping.Update(ctx, r)
	return
}

func (s *GnbootService) DeleteVideoKeywordMapping(ctx context.Context, req *params.IdsRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "DeleteVideoKeywordMapping")
	defer span.End()
	rp = &emptypb.Empty{}
	err = s.videoKeywordMapping.Delete(ctx, utils.Str2Uint64Arr(req.Ids)...)
	return
}
