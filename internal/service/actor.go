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

func (s *GnbootService) CreateActor(ctx context.Context, req *gnboot.CreateActorRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "CreateActor")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.Actor{}
	copierx.Copy(&r, req)
	err = s.actor.Create(ctx, r)
	return
}

func (s *GnbootService) GetActor(ctx context.Context, req *gnboot.GetActorRequest) (rp *gnboot.GetActorReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "GetActor")
	defer span.End()
	rp = &gnboot.GetActorReply{}
	res, err := s.actor.Get(ctx, req.Id)
	if err != nil {
		return
	}
	copierx.Copy(&rp, res)
	return
}

func (s *GnbootService) FindActor(ctx context.Context, req *gnboot.FindActorRequest) (rp *gnboot.FindActorReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "FindActor")
	defer span.End()
	rp = &gnboot.FindActorReply{}
	rp.Page = &params.Page{}
	r := &biz.FindActor{}
	r.Page = page.Page{}
	copierx.Copy(&r, req)
	copierx.Copy(&r.Page, req.Page)
	res, err := s.actor.Find(ctx, r)
	if err != nil {
		return
	}
	copierx.Copy(&rp.Page, r.Page)
	copierx.Copy(&rp.List, res)
	return
}

func (s *GnbootService) UpdateActor(ctx context.Context, req *gnboot.UpdateActorRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "UpdateActor")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.UpdateActor{}
	copierx.Copy(&r, req)
	err = s.actor.Update(ctx, r)
	return
}

func (s *GnbootService) DeleteActor(ctx context.Context, req *params.IdsRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "DeleteActor")
	defer span.End()
	rp = &emptypb.Empty{}
	err = s.actor.Delete(ctx, utils.Str2Uint64Arr(req.Ids)...)
	return
}
