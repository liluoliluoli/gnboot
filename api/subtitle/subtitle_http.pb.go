// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.8.1
// - protoc             v5.28.3
// source: proto/subtitle.proto

package subtitle

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationSubtitleRemoteServiceFindGenre = "/gnboot.SubtitleRemoteService/FindGenre"

type SubtitleRemoteServiceHTTPServer interface {
	FindGenre(context.Context, *FindSubtitleRequest) (*FindSubtitleResp, error)
}

func RegisterSubtitleRemoteServiceHTTPServer(s *http.Server, srv SubtitleRemoteServiceHTTPServer) {
	r := s.Route("/")
	r.POST("/subtitle/query/all", _SubtitleRemoteService_FindGenre0_HTTP_Handler(srv))
}

func _SubtitleRemoteService_FindGenre0_HTTP_Handler(srv SubtitleRemoteServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in FindSubtitleRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationSubtitleRemoteServiceFindGenre)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.FindGenre(ctx, req.(*FindSubtitleRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*FindSubtitleResp)
		return ctx.Result(200, reply)
	}
}

type SubtitleRemoteServiceHTTPClient interface {
	FindGenre(ctx context.Context, req *FindSubtitleRequest, opts ...http.CallOption) (rsp *FindSubtitleResp, err error)
}

type SubtitleRemoteServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewSubtitleRemoteServiceHTTPClient(client *http.Client) SubtitleRemoteServiceHTTPClient {
	return &SubtitleRemoteServiceHTTPClientImpl{client}
}

func (c *SubtitleRemoteServiceHTTPClientImpl) FindGenre(ctx context.Context, in *FindSubtitleRequest, opts ...http.CallOption) (*FindSubtitleResp, error) {
	var out FindSubtitleResp
	pattern := "/subtitle/query/all"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationSubtitleRemoteServiceFindGenre))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}