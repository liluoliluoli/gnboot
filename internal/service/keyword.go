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

func (s *GnbootService) CreateKeyword(ctx context.Context, req *gnboot.CreateKeywordRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "CreateKeyword")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.Keyword{}
	copierx.Copy(&r, req)
	err = s.keyword.Create(ctx, r)
	return
}

func (s *GnbootService) GetKeyword(ctx context.Context, req *gnboot.GetKeywordRequest) (rp *gnboot.GetKeywordReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "GetKeyword")
	defer span.End()
	rp = &gnboot.GetKeywordReply{}
	res, err := s.keyword.Get(ctx, req.Id)
	if err != nil {
		return
	}
	copierx.Copy(&rp, res)
	return
}

func (s *GnbootService) FindKeyword(ctx context.Context, req *gnboot.FindKeywordRequest) (rp *gnboot.FindKeywordReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "FindKeyword")
	defer span.End()
	rp = &gnboot.FindKeywordReply{}
	rp.Page = &params.Page{}
	r := &biz.FindKeyword{}
	r.Page = page.Page{}
	copierx.Copy(&r, req)
	copierx.Copy(&r.Page, req.Page)
	res, err := s.keyword.Find(ctx, r)
	if err != nil {
		return
	}
	copierx.Copy(&rp.Page, r.Page)
	copierx.Copy(&rp.List, res)
	return
}

func (s *GnbootService) UpdateKeyword(ctx context.Context, req *gnboot.UpdateKeywordRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "UpdateKeyword")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.UpdateKeyword{}
	copierx.Copy(&r, req)
	err = s.keyword.Update(ctx, r)
	return
}

func (s *GnbootService) DeleteKeyword(ctx context.Context, req *params.IdsRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "DeleteKeyword")
	defer span.End()
	rp = &emptypb.Empty{}
	err = s.keyword.Delete(ctx, utils.Str2Uint64Arr(req.Ids)...)
	return
}
