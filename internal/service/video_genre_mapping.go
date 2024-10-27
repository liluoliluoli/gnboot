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

func (s *GnbootService) CreateVideoGenreMapping(ctx context.Context, req *gnboot.CreateVideoGenreMappingRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "CreateVideoGenreMapping")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.VideoGenreMapping{}
	copierx.Copy(&r, req)
	err = s.videoGenreMapping.Create(ctx, r)
	return
}

func (s *GnbootService) GetVideoGenreMapping(ctx context.Context, req *gnboot.GetVideoGenreMappingRequest) (rp *gnboot.GetVideoGenreMappingReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "GetVideoGenreMapping")
	defer span.End()
	rp = &gnboot.GetVideoGenreMappingReply{}
	res, err := s.videoGenreMapping.Get(ctx, req.Id)
	if err != nil {
		return
	}
	copierx.Copy(&rp, res)
	return
}

func (s *GnbootService) FindVideoGenreMapping(ctx context.Context, req *gnboot.FindVideoGenreMappingRequest) (rp *gnboot.FindVideoGenreMappingReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "FindVideoGenreMapping")
	defer span.End()
	rp = &gnboot.FindVideoGenreMappingReply{}
	rp.Page = &params.Page{}
	r := &biz.FindVideoGenreMapping{}
	r.Page = page.Page{}
	copierx.Copy(&r, req)
	copierx.Copy(&r.Page, req.Page)
	res, err := s.videoGenreMapping.Find(ctx, r)
	if err != nil {
		return
	}
	copierx.Copy(&rp.Page, r.Page)
	copierx.Copy(&rp.List, res)
	return
}

func (s *GnbootService) UpdateVideoGenreMapping(ctx context.Context, req *gnboot.UpdateVideoGenreMappingRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "UpdateVideoGenreMapping")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.UpdateVideoGenreMapping{}
	copierx.Copy(&r, req)
	err = s.videoGenreMapping.Update(ctx, r)
	return
}

func (s *GnbootService) DeleteVideoGenreMapping(ctx context.Context, req *params.IdsRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "DeleteVideoGenreMapping")
	defer span.End()
	rp = &emptypb.Empty{}
	err = s.videoGenreMapping.Delete(ctx, utils.Str2Uint64Arr(req.Ids)...)
	return
}
