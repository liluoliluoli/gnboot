// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.8.1
// - protoc             v5.28.3
// source: proto/actor.proto

package actor

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

const OperationActorRemoteServiceFindGenre = "/gnboot.ActorRemoteService/FindGenre"

type ActorRemoteServiceHTTPServer interface {
	FindGenre(context.Context, *FindActorRequest) (*FindActorResp, error)
}

func RegisterActorRemoteServiceHTTPServer(s *http.Server, srv ActorRemoteServiceHTTPServer) {
	r := s.Route("/")
	r.POST("/actor/query/all", _ActorRemoteService_FindGenre0_HTTP_Handler(srv))
}

func _ActorRemoteService_FindGenre0_HTTP_Handler(srv ActorRemoteServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in FindActorRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationActorRemoteServiceFindGenre)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.FindGenre(ctx, req.(*FindActorRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*FindActorResp)
		return ctx.Result(200, reply)
	}
}

type ActorRemoteServiceHTTPClient interface {
	FindGenre(ctx context.Context, req *FindActorRequest, opts ...http.CallOption) (rsp *FindActorResp, err error)
}

type ActorRemoteServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewActorRemoteServiceHTTPClient(client *http.Client) ActorRemoteServiceHTTPClient {
	return &ActorRemoteServiceHTTPClientImpl{client}
}

func (c *ActorRemoteServiceHTTPClientImpl) FindGenre(ctx context.Context, in *FindActorRequest, opts ...http.CallOption) (*FindActorResp, error) {
	var out FindActorResp
	pattern := "/actor/query/all"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationActorRemoteServiceFindGenre))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}