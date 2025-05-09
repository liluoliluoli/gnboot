// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v5.29.3
// source: proto/appversion.proto

package appversion

import (
	_ "github.com/google/gnostic/openapiv3"
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

type AppVersion struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id            int32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	VersionCode   string `protobuf:"bytes,2,opt,name=versionCode,proto3" json:"versionCode,omitempty"`
	VersionName   string `protobuf:"bytes,3,opt,name=versionName,proto3" json:"versionName,omitempty"`
	PublishedTime int32  `protobuf:"varint,4,opt,name=publishedTime,proto3" json:"publishedTime,omitempty"`
	ForceUpdate   bool   `protobuf:"varint,5,opt,name=forceUpdate,proto3" json:"forceUpdate,omitempty"`
	ApkUrl        string `protobuf:"bytes,6,opt,name=apkUrl,proto3" json:"apkUrl,omitempty"`
	Remark        string `protobuf:"bytes,7,opt,name=remark,proto3" json:"remark,omitempty"`
}

func (x *AppVersion) Reset() {
	*x = AppVersion{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_appversion_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppVersion) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppVersion) ProtoMessage() {}

func (x *AppVersion) ProtoReflect() protoreflect.Message {
	mi := &file_proto_appversion_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppVersion.ProtoReflect.Descriptor instead.
func (*AppVersion) Descriptor() ([]byte, []int) {
	return file_proto_appversion_proto_rawDescGZIP(), []int{0}
}

func (x *AppVersion) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *AppVersion) GetVersionCode() string {
	if x != nil {
		return x.VersionCode
	}
	return ""
}

func (x *AppVersion) GetVersionName() string {
	if x != nil {
		return x.VersionName
	}
	return ""
}

func (x *AppVersion) GetPublishedTime() int32 {
	if x != nil {
		return x.PublishedTime
	}
	return 0
}

func (x *AppVersion) GetForceUpdate() bool {
	if x != nil {
		return x.ForceUpdate
	}
	return false
}

func (x *AppVersion) GetApkUrl() string {
	if x != nil {
		return x.ApkUrl
	}
	return ""
}

func (x *AppVersion) GetRemark() string {
	if x != nil {
		return x.Remark
	}
	return ""
}

type GetLastAppVersionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetLastAppVersionRequest) Reset() {
	*x = GetLastAppVersionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_appversion_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetLastAppVersionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetLastAppVersionRequest) ProtoMessage() {}

func (x *GetLastAppVersionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_appversion_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetLastAppVersionRequest.ProtoReflect.Descriptor instead.
func (*GetLastAppVersionRequest) Descriptor() ([]byte, []int) {
	return file_proto_appversion_proto_rawDescGZIP(), []int{1}
}

var File_proto_appversion_proto protoreflect.FileDescriptor

var file_proto_appversion_proto_rawDesc = []byte{
	0x0a, 0x16, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x70, 0x70, 0x76, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x67, 0x6e, 0x62, 0x6f, 0x6f, 0x74,
	0x1a, 0x1b, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x33, 0x2f, 0x61, 0x6e, 0x6e, 0x6f,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x62, 0x65,
	0x68, 0x61, 0x76, 0x69, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d,
	0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd8, 0x01, 0x0a, 0x0a, 0x41,
	0x70, 0x70, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x76, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x24, 0x0a,
	0x0d, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x65, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x0d, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x65, 0x64, 0x54,
	0x69, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x70, 0x6b, 0x55, 0x72, 0x6c, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x70, 0x6b, 0x55, 0x72, 0x6c, 0x12, 0x16, 0x0a,
	0x06, 0x72, 0x65, 0x6d, 0x61, 0x72, 0x6b, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72,
	0x65, 0x6d, 0x61, 0x72, 0x6b, 0x22, 0x1a, 0x0a, 0x18, 0x47, 0x65, 0x74, 0x4c, 0x61, 0x73, 0x74,
	0x41, 0x70, 0x70, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x32, 0x86, 0x01, 0x0a, 0x17, 0x41, 0x70, 0x70, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x6b, 0x0a,
	0x0e, 0x47, 0x65, 0x74, 0x4c, 0x61, 0x73, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12,
	0x20, 0x2e, 0x67, 0x6e, 0x62, 0x6f, 0x6f, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x4c, 0x61, 0x73, 0x74,
	0x41, 0x70, 0x70, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x12, 0x2e, 0x67, 0x6e, 0x62, 0x6f, 0x6f, 0x74, 0x2e, 0x41, 0x70, 0x70, 0x56, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x23, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1d, 0x12, 0x1b, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2f, 0x67, 0x65, 0x74, 0x4c,
	0x61, 0x73, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x42, 0x7d, 0xba, 0x47, 0x4b, 0x12,
	0x13, 0x0a, 0x0a, 0x61, 0x70, 0x70, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x32, 0x05, 0x31,
	0x2e, 0x30, 0x2e, 0x30, 0x2a, 0x22, 0x3a, 0x20, 0x0a, 0x1e, 0x0a, 0x0a, 0x42, 0x65, 0x61, 0x72,
	0x65, 0x72, 0x41, 0x75, 0x74, 0x68, 0x12, 0x10, 0x0a, 0x0e, 0x0a, 0x04, 0x68, 0x74, 0x74, 0x70,
	0x2a, 0x06, 0x62, 0x65, 0x61, 0x72, 0x65, 0x72, 0x32, 0x10, 0x0a, 0x0e, 0x0a, 0x0a, 0x42, 0x65,
	0x61, 0x72, 0x65, 0x72, 0x41, 0x75, 0x74, 0x68, 0x12, 0x00, 0x5a, 0x2d, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6c, 0x69, 0x6c, 0x75, 0x6f, 0x6c, 0x69, 0x6c, 0x75,
	0x6f, 0x6c, 0x69, 0x2f, 0x67, 0x6e, 0x62, 0x6f, 0x6f, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61,
	0x70, 0x70, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_proto_appversion_proto_rawDescOnce sync.Once
	file_proto_appversion_proto_rawDescData = file_proto_appversion_proto_rawDesc
)

func file_proto_appversion_proto_rawDescGZIP() []byte {
	file_proto_appversion_proto_rawDescOnce.Do(func() {
		file_proto_appversion_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_appversion_proto_rawDescData)
	})
	return file_proto_appversion_proto_rawDescData
}

var file_proto_appversion_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_appversion_proto_goTypes = []interface{}{
	(*AppVersion)(nil),               // 0: gnboot.AppVersion
	(*GetLastAppVersionRequest)(nil), // 1: gnboot.GetLastAppVersionRequest
}
var file_proto_appversion_proto_depIdxs = []int32{
	1, // 0: gnboot.AppVersionRemoteService.GetLastVersion:input_type -> gnboot.GetLastAppVersionRequest
	0, // 1: gnboot.AppVersionRemoteService.GetLastVersion:output_type -> gnboot.AppVersion
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_appversion_proto_init() }
func file_proto_appversion_proto_init() {
	if File_proto_appversion_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_appversion_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppVersion); i {
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
		file_proto_appversion_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetLastAppVersionRequest); i {
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
			RawDescriptor: file_proto_appversion_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_appversion_proto_goTypes,
		DependencyIndexes: file_proto_appversion_proto_depIdxs,
		MessageInfos:      file_proto_appversion_proto_msgTypes,
	}.Build()
	File_proto_appversion_proto = out.File
	file_proto_appversion_proto_rawDesc = nil
	file_proto_appversion_proto_goTypes = nil
	file_proto_appversion_proto_depIdxs = nil
}
