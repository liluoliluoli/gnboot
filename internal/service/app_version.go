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

func (s *GnbootService) CreateAppVersion(ctx context.Context, req *gnboot.CreateAppVersionRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "CreateAppVersion")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.AppVersion{}
	copierx.Copy(&r, req)
	err = s.appVersion.Create(ctx, r)
	return
}

func (s *GnbootService) GetAppVersion(ctx context.Context, req *gnboot.GetAppVersionRequest) (rp *gnboot.GetAppVersionReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "GetAppVersion")
	defer span.End()
	rp = &gnboot.GetAppVersionReply{}
	res, err := s.appVersion.Get(ctx, req.Id)
	if err != nil {
		return
	}
	copierx.Copy(&rp, res)
	return
}

func (s *GnbootService) FindAppVersion(ctx context.Context, req *gnboot.FindAppVersionRequest) (rp *gnboot.FindAppVersionReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "FindAppVersion")
	defer span.End()
	rp = &gnboot.FindAppVersionReply{}
	rp.Page = &params.Page{}
	r := &biz.FindAppVersion{}
	r.Page = page.Page{}
	copierx.Copy(&r, req)
	copierx.Copy(&r.Page, req.Page)
	res, err := s.appVersion.Find(ctx, r)
	if err != nil {
		return
	}
	copierx.Copy(&rp.Page, r.Page)
	copierx.Copy(&rp.List, res)
	return
}

func (s *GnbootService) UpdateAppVersion(ctx context.Context, req *gnboot.UpdateAppVersionRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "UpdateAppVersion")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.UpdateAppVersion{}
	copierx.Copy(&r, req)
	err = s.appVersion.Update(ctx, r)
	return
}

func (s *GnbootService) DeleteAppVersion(ctx context.Context, req *params.IdsRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "DeleteAppVersion")
	defer span.End()
	rp = &emptypb.Empty{}
	err = s.appVersion.Delete(ctx, utils.Str2Uint64Arr(req.Ids)...)
	return
}
