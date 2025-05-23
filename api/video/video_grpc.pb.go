// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.29.3
// source: proto/video.proto

package video

import (
	context "context"
	api "github.com/liluoliluoli/gnboot/api"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	VideoRemoteService_CreateVideo_FullMethodName   = "/gnboot.VideoRemoteService/CreateVideo"
	VideoRemoteService_GetVideo_FullMethodName      = "/gnboot.VideoRemoteService/GetVideo"
	VideoRemoteService_SearchVideo_FullMethodName   = "/gnboot.VideoRemoteService/SearchVideo"
	VideoRemoteService_UpdateVideo_FullMethodName   = "/gnboot.VideoRemoteService/UpdateVideo"
	VideoRemoteService_DeleteVideo_FullMethodName   = "/gnboot.VideoRemoteService/DeleteVideo"
	VideoRemoteService_PageFavorites_FullMethodName = "/gnboot.VideoRemoteService/PageFavorites"
)

// VideoRemoteServiceClient is the client API for VideoRemoteService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VideoRemoteServiceClient interface {
	CreateVideo(ctx context.Context, in *CreateVideoRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetVideo(ctx context.Context, in *GetVideoRequest, opts ...grpc.CallOption) (*Video, error)
	SearchVideo(ctx context.Context, in *SearchVideoRequest, opts ...grpc.CallOption) (*SearchVideoResp, error)
	UpdateVideo(ctx context.Context, in *UpdateVideoRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DeleteVideo(ctx context.Context, in *api.IdsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	PageFavorites(ctx context.Context, in *PageFavoritesRequest, opts ...grpc.CallOption) (*SearchVideoResp, error)
}

type videoRemoteServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewVideoRemoteServiceClient(cc grpc.ClientConnInterface) VideoRemoteServiceClient {
	return &videoRemoteServiceClient{cc}
}

func (c *videoRemoteServiceClient) CreateVideo(ctx context.Context, in *CreateVideoRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, VideoRemoteService_CreateVideo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoRemoteServiceClient) GetVideo(ctx context.Context, in *GetVideoRequest, opts ...grpc.CallOption) (*Video, error) {
	out := new(Video)
	err := c.cc.Invoke(ctx, VideoRemoteService_GetVideo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoRemoteServiceClient) SearchVideo(ctx context.Context, in *SearchVideoRequest, opts ...grpc.CallOption) (*SearchVideoResp, error) {
	out := new(SearchVideoResp)
	err := c.cc.Invoke(ctx, VideoRemoteService_SearchVideo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoRemoteServiceClient) UpdateVideo(ctx context.Context, in *UpdateVideoRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, VideoRemoteService_UpdateVideo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoRemoteServiceClient) DeleteVideo(ctx context.Context, in *api.IdsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, VideoRemoteService_DeleteVideo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoRemoteServiceClient) PageFavorites(ctx context.Context, in *PageFavoritesRequest, opts ...grpc.CallOption) (*SearchVideoResp, error) {
	out := new(SearchVideoResp)
	err := c.cc.Invoke(ctx, VideoRemoteService_PageFavorites_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VideoRemoteServiceServer is the server API for VideoRemoteService service.
// All implementations must embed UnimplementedVideoRemoteServiceServer
// for forward compatibility
type VideoRemoteServiceServer interface {
	CreateVideo(context.Context, *CreateVideoRequest) (*emptypb.Empty, error)
	GetVideo(context.Context, *GetVideoRequest) (*Video, error)
	SearchVideo(context.Context, *SearchVideoRequest) (*SearchVideoResp, error)
	UpdateVideo(context.Context, *UpdateVideoRequest) (*emptypb.Empty, error)
	DeleteVideo(context.Context, *api.IdsRequest) (*emptypb.Empty, error)
	PageFavorites(context.Context, *PageFavoritesRequest) (*SearchVideoResp, error)
	mustEmbedUnimplementedVideoRemoteServiceServer()
}

// UnimplementedVideoRemoteServiceServer must be embedded to have forward compatible implementations.
type UnimplementedVideoRemoteServiceServer struct {
}

func (UnimplementedVideoRemoteServiceServer) CreateVideo(context.Context, *CreateVideoRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateVideo not implemented")
}
func (UnimplementedVideoRemoteServiceServer) GetVideo(context.Context, *GetVideoRequest) (*Video, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVideo not implemented")
}
func (UnimplementedVideoRemoteServiceServer) SearchVideo(context.Context, *SearchVideoRequest) (*SearchVideoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchVideo not implemented")
}
func (UnimplementedVideoRemoteServiceServer) UpdateVideo(context.Context, *UpdateVideoRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateVideo not implemented")
}
func (UnimplementedVideoRemoteServiceServer) DeleteVideo(context.Context, *api.IdsRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteVideo not implemented")
}
func (UnimplementedVideoRemoteServiceServer) PageFavorites(context.Context, *PageFavoritesRequest) (*SearchVideoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PageFavorites not implemented")
}
func (UnimplementedVideoRemoteServiceServer) mustEmbedUnimplementedVideoRemoteServiceServer() {}

// UnsafeVideoRemoteServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VideoRemoteServiceServer will
// result in compilation errors.
type UnsafeVideoRemoteServiceServer interface {
	mustEmbedUnimplementedVideoRemoteServiceServer()
}

func RegisterVideoRemoteServiceServer(s grpc.ServiceRegistrar, srv VideoRemoteServiceServer) {
	s.RegisterService(&VideoRemoteService_ServiceDesc, srv)
}

func _VideoRemoteService_CreateVideo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateVideoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoRemoteServiceServer).CreateVideo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VideoRemoteService_CreateVideo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoRemoteServiceServer).CreateVideo(ctx, req.(*CreateVideoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VideoRemoteService_GetVideo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetVideoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoRemoteServiceServer).GetVideo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VideoRemoteService_GetVideo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoRemoteServiceServer).GetVideo(ctx, req.(*GetVideoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VideoRemoteService_SearchVideo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchVideoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoRemoteServiceServer).SearchVideo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VideoRemoteService_SearchVideo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoRemoteServiceServer).SearchVideo(ctx, req.(*SearchVideoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VideoRemoteService_UpdateVideo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateVideoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoRemoteServiceServer).UpdateVideo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VideoRemoteService_UpdateVideo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoRemoteServiceServer).UpdateVideo(ctx, req.(*UpdateVideoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VideoRemoteService_DeleteVideo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(api.IdsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoRemoteServiceServer).DeleteVideo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VideoRemoteService_DeleteVideo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoRemoteServiceServer).DeleteVideo(ctx, req.(*api.IdsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VideoRemoteService_PageFavorites_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PageFavoritesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoRemoteServiceServer).PageFavorites(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VideoRemoteService_PageFavorites_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoRemoteServiceServer).PageFavorites(ctx, req.(*PageFavoritesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// VideoRemoteService_ServiceDesc is the grpc.ServiceDesc for VideoRemoteService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var VideoRemoteService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gnboot.VideoRemoteService",
	HandlerType: (*VideoRemoteServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateVideo",
			Handler:    _VideoRemoteService_CreateVideo_Handler,
		},
		{
			MethodName: "GetVideo",
			Handler:    _VideoRemoteService_GetVideo_Handler,
		},
		{
			MethodName: "SearchVideo",
			Handler:    _VideoRemoteService_SearchVideo_Handler,
		},
		{
			MethodName: "UpdateVideo",
			Handler:    _VideoRemoteService_UpdateVideo_Handler,
		},
		{
			MethodName: "DeleteVideo",
			Handler:    _VideoRemoteService_DeleteVideo_Handler,
		},
		{
			MethodName: "PageFavorites",
			Handler:    _VideoRemoteService_PageFavorites_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/video.proto",
}
