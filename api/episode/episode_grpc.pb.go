// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.28.3
// source: proto/episode.proto

package episode

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	EpisodeRemoteService_GetEpisode_FullMethodName = "/gnboot.EpisodeRemoteService/GetEpisode"
)

// EpisodeRemoteServiceClient is the client API for EpisodeRemoteService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EpisodeRemoteServiceClient interface {
	GetEpisode(ctx context.Context, in *GetEpisodeRequest, opts ...grpc.CallOption) (*EpisodeResp, error)
}

type episodeRemoteServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEpisodeRemoteServiceClient(cc grpc.ClientConnInterface) EpisodeRemoteServiceClient {
	return &episodeRemoteServiceClient{cc}
}

func (c *episodeRemoteServiceClient) GetEpisode(ctx context.Context, in *GetEpisodeRequest, opts ...grpc.CallOption) (*EpisodeResp, error) {
	out := new(EpisodeResp)
	err := c.cc.Invoke(ctx, EpisodeRemoteService_GetEpisode_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EpisodeRemoteServiceServer is the server API for EpisodeRemoteService service.
// All implementations must embed UnimplementedEpisodeRemoteServiceServer
// for forward compatibility
type EpisodeRemoteServiceServer interface {
	GetEpisode(context.Context, *GetEpisodeRequest) (*EpisodeResp, error)
	mustEmbedUnimplementedEpisodeRemoteServiceServer()
}

// UnimplementedEpisodeRemoteServiceServer must be embedded to have forward compatible implementations.
type UnimplementedEpisodeRemoteServiceServer struct {
}

func (UnimplementedEpisodeRemoteServiceServer) GetEpisode(context.Context, *GetEpisodeRequest) (*EpisodeResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEpisode not implemented")
}
func (UnimplementedEpisodeRemoteServiceServer) mustEmbedUnimplementedEpisodeRemoteServiceServer() {}

// UnsafeEpisodeRemoteServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EpisodeRemoteServiceServer will
// result in compilation errors.
type UnsafeEpisodeRemoteServiceServer interface {
	mustEmbedUnimplementedEpisodeRemoteServiceServer()
}

func RegisterEpisodeRemoteServiceServer(s grpc.ServiceRegistrar, srv EpisodeRemoteServiceServer) {
	s.RegisterService(&EpisodeRemoteService_ServiceDesc, srv)
}

func _EpisodeRemoteService_GetEpisode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEpisodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EpisodeRemoteServiceServer).GetEpisode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EpisodeRemoteService_GetEpisode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EpisodeRemoteServiceServer).GetEpisode(ctx, req.(*GetEpisodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// EpisodeRemoteService_ServiceDesc is the grpc.ServiceDesc for EpisodeRemoteService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EpisodeRemoteService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gnboot.EpisodeRemoteService",
	HandlerType: (*EpisodeRemoteServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetEpisode",
			Handler:    _EpisodeRemoteService_GetEpisode_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/episode.proto",
}