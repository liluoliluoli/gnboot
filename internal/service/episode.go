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

func (s *GnbootService) CreateEpisode(ctx context.Context, req *gnboot.CreateEpisodeRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "CreateEpisode")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.Episode{}
	copierx.Copy(&r, req)
	err = s.episode.Create(ctx, r)
	return
}

func (s *GnbootService) GetEpisode(ctx context.Context, req *gnboot.GetEpisodeRequest) (rp *gnboot.GetEpisodeReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "GetEpisode")
	defer span.End()
	rp = &gnboot.GetEpisodeReply{}
	res, err := s.episode.Get(ctx, req.Id)
	if err != nil {
		return
	}
	copierx.Copy(&rp, res)
	return
}

func (s *GnbootService) FindEpisode(ctx context.Context, req *gnboot.FindEpisodeRequest) (rp *gnboot.FindEpisodeReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "FindEpisode")
	defer span.End()
	rp = &gnboot.FindEpisodeReply{}
	rp.Page = &params.Page{}
	r := &biz.FindEpisode{}
	r.Page = page.Page{}
	copierx.Copy(&r, req)
	copierx.Copy(&r.Page, req.Page)
	res, err := s.episode.Find(ctx, r)
	if err != nil {
		return
	}
	copierx.Copy(&rp.Page, r.Page)
	copierx.Copy(&rp.List, res)
	return
}

func (s *GnbootService) UpdateEpisode(ctx context.Context, req *gnboot.UpdateEpisodeRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "UpdateEpisode")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.UpdateEpisode{}
	copierx.Copy(&r, req)
	err = s.episode.Update(ctx, r)
	return
}

func (s *GnbootService) DeleteEpisode(ctx context.Context, req *params.IdsRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "DeleteEpisode")
	defer span.End()
	rp = &emptypb.Empty{}
	err = s.episode.Delete(ctx, utils.Str2Uint64Arr(req.Ids)...)
	return
}
