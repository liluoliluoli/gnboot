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

func (s *GnbootService) CreateUser(ctx context.Context, req *gnboot.CreateUserRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "CreateUser")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.User{}
	copierx.Copy(&r, req)
	err = s.user.Create(ctx, r)
	return
}

func (s *GnbootService) GetUser(ctx context.Context, req *gnboot.GetUserRequest) (rp *gnboot.GetUserReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "GetUser")
	defer span.End()
	rp = &gnboot.GetUserReply{}
	res, err := s.user.Get(ctx, req.Id)
	if err != nil {
		return
	}
	copierx.Copy(&rp, res)
	return
}

func (s *GnbootService) FindUser(ctx context.Context, req *gnboot.FindUserRequest) (rp *gnboot.FindUserReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "FindUser")
	defer span.End()
	rp = &gnboot.FindUserReply{}
	rp.Page = &params.Page{}
	r := &biz.FindUser{}
	r.Page = page.Page{}
	copierx.Copy(&r, req)
	copierx.Copy(&r.Page, req.Page)
	res, err := s.user.Find(ctx, r)
	if err != nil {
		return
	}
	copierx.Copy(&rp.Page, r.Page)
	copierx.Copy(&rp.List, res)
	return
}

func (s *GnbootService) UpdateUser(ctx context.Context, req *gnboot.UpdateUserRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "UpdateUser")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.UpdateUser{}
	copierx.Copy(&r, req)
	err = s.user.Update(ctx, r)
	return
}

func (s *GnbootService) DeleteUser(ctx context.Context, req *params.IdsRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "DeleteUser")
	defer span.End()
	rp = &emptypb.Empty{}
	err = s.user.Delete(ctx, utils.Str2Uint64Arr(req.Ids)...)
	return
}
