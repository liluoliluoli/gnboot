// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v5.29.3
// source: proto/episode.proto

package episode

import (
	_ "github.com/google/gnostic/openapiv3"
	_ "github.com/liluoliluoli/gnboot/api/actor"
	subtitle "github.com/liluoliluoli/gnboot/api/subtitle"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/emptypb"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Episode struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           int32                `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	VideoId      int32                `protobuf:"varint,2,opt,name=videoId,proto3" json:"videoId,omitempty"`
	Episode      int32                `protobuf:"varint,3,opt,name=episode,proto3" json:"episode,omitempty"`
	EpisodeTitle string               `protobuf:"bytes,4,opt,name=episodeTitle,proto3" json:"episodeTitle,omitempty"`
	Url          string               `protobuf:"bytes,5,opt,name=url,proto3" json:"url,omitempty"`
	Platform     string               `protobuf:"bytes,6,opt,name=platform,proto3" json:"platform,omitempty"`
	Ext          string               `protobuf:"bytes,7,opt,name=ext,proto3" json:"ext,omitempty"`
	Duration     int32                `protobuf:"varint,8,opt,name=duration,proto3" json:"duration,omitempty"`
	Subtitles    []*subtitle.Subtitle `protobuf:"bytes,9,rep,name=subtitles,proto3" json:"subtitles,omitempty"`
}

func (x *Episode) Reset() {
	*x = Episode{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_episode_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Episode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Episode) ProtoMessage() {}

func (x *Episode) ProtoReflect() protoreflect.Message {
	mi := &file_proto_episode_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Episode.ProtoReflect.Descriptor instead.
func (*Episode) Descriptor() ([]byte, []int) {
	return file_proto_episode_proto_rawDescGZIP(), []int{0}
}

func (x *Episode) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Episode) GetVideoId() int32 {
	if x != nil {
		return x.VideoId
	}
	return 0
}

func (x *Episode) GetEpisode() int32 {
	if x != nil {
		return x.Episode
	}
	return 0
}

func (x *Episode) GetEpisodeTitle() string {
	if x != nil {
		return x.EpisodeTitle
	}
	return ""
}

func (x *Episode) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *Episode) GetPlatform() string {
	if x != nil {
		return x.Platform
	}
	return ""
}

func (x *Episode) GetExt() string {
	if x != nil {
		return x.Ext
	}
	return ""
}

func (x *Episode) GetDuration() int32 {
	if x != nil {
		return x.Duration
	}
	return 0
}

func (x *Episode) GetSubtitles() []*subtitle.Subtitle {
	if x != nil {
		return x.Subtitles
	}
	return nil
}

type GetEpisodeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetEpisodeRequest) Reset() {
	*x = GetEpisodeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_episode_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetEpisodeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetEpisodeRequest) ProtoMessage() {}

func (x *GetEpisodeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_episode_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetEpisodeRequest.ProtoReflect.Descriptor instead.
func (*GetEpisodeRequest) Descriptor() ([]byte, []int) {
	return file_proto_episode_proto_rawDescGZIP(), []int{1}
}

func (x *GetEpisodeRequest) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

var File_proto_episode_proto protoreflect.FileDescriptor

var file_proto_episode_proto_rawDesc = []byte{
	0x0a, 0x13, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x70, 0x69, 0x73, 0x6f, 0x64, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x67, 0x6e, 0x62, 0x6f, 0x6f, 0x74, 0x1a, 0x1b, 0x6f,
	0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x33, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x62, 0x65, 0x68, 0x61, 0x76,
	0x69, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x14, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73,
	0x75, 0x62, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xfd, 0x01, 0x0a, 0x07, 0x45, 0x70, 0x69, 0x73, 0x6f, 0x64, 0x65, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07,
	0x76, 0x69, 0x64, 0x65, 0x6f, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x76,
	0x69, 0x64, 0x65, 0x6f, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x70, 0x69, 0x73, 0x6f, 0x64,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x65, 0x70, 0x69, 0x73, 0x6f, 0x64, 0x65,
	0x12, 0x22, 0x0a, 0x0c, 0x65, 0x70, 0x69, 0x73, 0x6f, 0x64, 0x65, 0x54, 0x69, 0x74, 0x6c, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x65, 0x70, 0x69, 0x73, 0x6f, 0x64, 0x65, 0x54,
	0x69, 0x74, 0x6c, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f,
	0x72, 0x6d, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f,
	0x72, 0x6d, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x78, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x65, 0x78, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x2e, 0x0a, 0x09, 0x73, 0x75, 0x62, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x73, 0x18, 0x09, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x67, 0x6e, 0x62, 0x6f, 0x6f, 0x74, 0x2e, 0x53, 0x75, 0x62,
	0x74, 0x69, 0x74, 0x6c, 0x65, 0x52, 0x09, 0x73, 0x75, 0x62, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x73,
	0x22, 0x23, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x45, 0x70, 0x69, 0x73, 0x6f, 0x64, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x02, 0x69, 0x64, 0x32, 0x6d, 0x0a, 0x14, 0x45, 0x70, 0x69, 0x73, 0x6f, 0x64, 0x65,
	0x52, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x55, 0x0a,
	0x0a, 0x47, 0x65, 0x74, 0x45, 0x70, 0x69, 0x73, 0x6f, 0x64, 0x65, 0x12, 0x19, 0x2e, 0x67, 0x6e,
	0x62, 0x6f, 0x6f, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x45, 0x70, 0x69, 0x73, 0x6f, 0x64, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x67, 0x6e, 0x62, 0x6f, 0x6f, 0x74, 0x2e,
	0x45, 0x70, 0x69, 0x73, 0x6f, 0x64, 0x65, 0x22, 0x1b, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x15, 0x3a,
	0x01, 0x2a, 0x22, 0x10, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x70, 0x69, 0x73, 0x6f, 0x64, 0x65,
	0x2f, 0x67, 0x65, 0x74, 0x42, 0x77, 0xba, 0x47, 0x48, 0x12, 0x10, 0x0a, 0x07, 0x65, 0x70, 0x69,
	0x73, 0x6f, 0x64, 0x65, 0x32, 0x05, 0x31, 0x2e, 0x30, 0x2e, 0x30, 0x2a, 0x22, 0x3a, 0x20, 0x0a,
	0x1e, 0x0a, 0x0a, 0x42, 0x65, 0x61, 0x72, 0x65, 0x72, 0x41, 0x75, 0x74, 0x68, 0x12, 0x10, 0x0a,
	0x0e, 0x0a, 0x04, 0x68, 0x74, 0x74, 0x70, 0x2a, 0x06, 0x62, 0x65, 0x61, 0x72, 0x65, 0x72, 0x32,
	0x10, 0x0a, 0x0e, 0x0a, 0x0a, 0x42, 0x65, 0x61, 0x72, 0x65, 0x72, 0x41, 0x75, 0x74, 0x68, 0x12,
	0x00, 0x5a, 0x2a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6c, 0x69,
	0x6c, 0x75, 0x6f, 0x6c, 0x69, 0x6c, 0x75, 0x6f, 0x6c, 0x69, 0x2f, 0x67, 0x6e, 0x62, 0x6f, 0x6f,
	0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x70, 0x69, 0x73, 0x6f, 0x64, 0x65, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_episode_proto_rawDescOnce sync.Once
	file_proto_episode_proto_rawDescData = file_proto_episode_proto_rawDesc
)

func file_proto_episode_proto_rawDescGZIP() []byte {
	file_proto_episode_proto_rawDescOnce.Do(func() {
		file_proto_episode_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_episode_proto_rawDescData)
	})
	return file_proto_episode_proto_rawDescData
}

var file_proto_episode_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_episode_proto_goTypes = []interface{}{
	(*Episode)(nil),           // 0: gnboot.Episode
	(*GetEpisodeRequest)(nil), // 1: gnboot.GetEpisodeRequest
	(*subtitle.Subtitle)(nil), // 2: gnboot.Subtitle
}
var file_proto_episode_proto_depIdxs = []int32{
	2, // 0: gnboot.Episode.subtitles:type_name -> gnboot.Subtitle
	1, // 1: gnboot.EpisodeRemoteService.GetEpisode:input_type -> gnboot.GetEpisodeRequest
	0, // 2: gnboot.EpisodeRemoteService.GetEpisode:output_type -> gnboot.Episode
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_episode_proto_init() }
func file_proto_episode_proto_init() {
	if File_proto_episode_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_episode_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Episode); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_episode_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetEpisodeRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_episode_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_episode_proto_goTypes,
		DependencyIndexes: file_proto_episode_proto_depIdxs,
		MessageInfos:      file_proto_episode_proto_msgTypes,
	}.Build()
	File_proto_episode_proto = out.File
	file_proto_episode_proto_rawDesc = nil
	file_proto_episode_proto_goTypes = nil
	file_proto_episode_proto_depIdxs = nil
}
