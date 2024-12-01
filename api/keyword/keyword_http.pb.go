// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.8.1
// - protoc             v5.28.3
// source: proto/keyword.proto

package keyword

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

const OperationKeywordRemoteServiceFindGenre = "/gnboot.KeywordRemoteService/FindGenre"

type KeywordRemoteServiceHTTPServer interface {
	FindGenre(context.Context, *FindKeywordRequest) (*FindKeywordResp, error)
}

func RegisterKeywordRemoteServiceHTTPServer(s *http.Server, srv KeywordRemoteServiceHTTPServer) {
	r := s.Route("/")
	r.POST("/keyword/query/all", _KeywordRemoteService_FindGenre0_HTTP_Handler(srv))
}

func _KeywordRemoteService_FindGenre0_HTTP_Handler(srv KeywordRemoteServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in FindKeywordRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationKeywordRemoteServiceFindGenre)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.FindGenre(ctx, req.(*FindKeywordRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*FindKeywordResp)
		return ctx.Result(200, reply)
	}
}

type KeywordRemoteServiceHTTPClient interface {
	FindGenre(ctx context.Context, req *FindKeywordRequest, opts ...http.CallOption) (rsp *FindKeywordResp, err error)
}

type KeywordRemoteServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewKeywordRemoteServiceHTTPClient(client *http.Client) KeywordRemoteServiceHTTPClient {
	return &KeywordRemoteServiceHTTPClientImpl{client}
}

func (c *KeywordRemoteServiceHTTPClientImpl) FindGenre(ctx context.Context, in *FindKeywordRequest, opts ...http.CallOption) (*FindKeywordResp, error) {
	var out FindKeywordResp
	pattern := "/keyword/query/all"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationKeywordRemoteServiceFindGenre))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
