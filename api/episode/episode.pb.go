// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v5.28.3
// source: proto/episode.proto

package episode

import (
	_ "github.com/google/gnostic/openapiv3"
	actor "github.com/liluoliluoli/gnboot/api/actor"
	_ "github.com/liluoliluoli/gnboot/api/genre"
	_ "github.com/liluoliluoli/gnboot/api/keyword"
	_ "github.com/liluoliluoli/gnboot/api/studio"
	subtitle "github.com/liluoliluoli/gnboot/api/subtitle"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type EpisodeResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id                 int32                    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Episode            int32                    `protobuf:"varint,2,opt,name=episode,proto3" json:"episode,omitempty"`
	Url                string                   `protobuf:"bytes,3,opt,name=url,proto3" json:"url,omitempty"`
	Download           bool                     `protobuf:"varint,4,opt,name=download,proto3" json:"download,omitempty"`
	Ext                string                   `protobuf:"bytes,5,opt,name=ext,proto3" json:"ext,omitempty"`
	FileSize           int32                    `protobuf:"varint,6,opt,name=fileSize,proto3" json:"fileSize,omitempty"`
	Subtitles          []*subtitle.SubtitleResp `protobuf:"bytes,7,rep,name=subtitles,proto3" json:"subtitles,omitempty"`
	LastPlayedPosition int32                    `protobuf:"varint,8,opt,name=lastPlayedPosition,proto3" json:"lastPlayedPosition,omitempty"`
	LastPlayedTime     *timestamppb.Timestamp   `protobuf:"bytes,9,opt,name=lastPlayedTime,proto3" json:"lastPlayedTime,omitempty"`
	SkipIntro          int32                    `protobuf:"varint,10,opt,name=skipIntro,proto3" json:"skipIntro,omitempty"`
	SkipEnding         int32                    `protobuf:"varint,11,opt,name=skipEnding,proto3" json:"skipEnding,omitempty"`
	Title              string                   `protobuf:"bytes,12,opt,name=title,proto3" json:"title,omitempty"`
	Poster             string                   `protobuf:"bytes,13,opt,name=poster,proto3" json:"poster,omitempty"`
	Logo               string                   `protobuf:"bytes,14,opt,name=logo,proto3" json:"logo,omitempty"`
	AirDate            *timestamppb.Timestamp   `protobuf:"bytes,15,opt,name=airDate,proto3" json:"airDate,omitempty"`
	Overview           string                   `protobuf:"bytes,16,opt,name=overview,proto3" json:"overview,omitempty"`
	Favorite           bool                     `protobuf:"varint,17,opt,name=favorite,proto3" json:"favorite,omitempty"`
	SeasonId           int32                    `protobuf:"varint,18,opt,name=seasonId,proto3" json:"seasonId,omitempty"`
	Season             int32                    `protobuf:"varint,19,opt,name=season,proto3" json:"season,omitempty"`
	SeasonTitle        string                   `protobuf:"bytes,20,opt,name=seasonTitle,proto3" json:"seasonTitle,omitempty"`
	SeriesTitle        string                   `protobuf:"bytes,21,opt,name=seriesTitle,proto3" json:"seriesTitle,omitempty"`
	Actors             []*actor.ActorResp       `protobuf:"bytes,22,rep,name=actors,proto3" json:"actors,omitempty"`
	Filename           string                   `protobuf:"bytes,23,opt,name=filename,proto3" json:"filename,omitempty"`
}

func (x *EpisodeResp) Reset() {
	*x = EpisodeResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_episode_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EpisodeResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EpisodeResp) ProtoMessage() {}

func (x *EpisodeResp) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use EpisodeResp.ProtoReflect.Descriptor instead.
func (*EpisodeResp) Descriptor() ([]byte, []int) {
	return file_proto_episode_proto_rawDescGZIP(), []int{0}
}

func (x *EpisodeResp) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *EpisodeResp) GetEpisode() int32 {
	if x != nil {
		return x.Episode
	}
	return 0
}

func (x *EpisodeResp) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *EpisodeResp) GetDownload() bool {
	if x != nil {
		return x.Download
	}
	return false
}

func (x *EpisodeResp) GetExt() string {
	if x != nil {
		return x.Ext
	}
	return ""
}

func (x *EpisodeResp) GetFileSize() int32 {
	if x != nil {
		return x.FileSize
	}
	return 0
}

func (x *EpisodeResp) GetSubtitles() []*subtitle.SubtitleResp {
	if x != nil {
		return x.Subtitles
	}
	return nil
}

func (x *EpisodeResp) GetLastPlayedPosition() int32 {
	if x != nil {
		return x.LastPlayedPosition
	}
	return 0
}

func (x *EpisodeResp) GetLastPlayedTime() *timestamppb.Timestamp {
	if x != nil {
		return x.LastPlayedTime
	}
	return nil
}

func (x *EpisodeResp) GetSkipIntro() int32 {
	if x != nil {
		return x.SkipIntro
	}
	return 0
}

func (x *EpisodeResp) GetSkipEnding() int32 {
	if x != nil {
		return x.SkipEnding
	}
	return 0
}

func (x *EpisodeResp) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *EpisodeResp) GetPoster() string {
	if x != nil {
		return x.Poster
	}
	return ""
}

func (x *EpisodeResp) GetLogo() string {
	if x != nil {
		return x.Logo
	}
	return ""
}

func (x *EpisodeResp) GetAirDate() *timestamppb.Timestamp {
	if x != nil {
		return x.AirDate
	}
	return nil
}

func (x *EpisodeResp) GetOverview() string {
	if x != nil {
		return x.Overview
	}
	return ""
}

func (x *EpisodeResp) GetFavorite() bool {
	if x != nil {
		return x.Favorite
	}
	return false
}

func (x *EpisodeResp) GetSeasonId() int32 {
	if x != nil {
		return x.SeasonId
	}
	return 0
}

func (x *EpisodeResp) GetSeason() int32 {
	if x != nil {
		return x.Season
	}
	return 0
}

func (x *EpisodeResp) GetSeasonTitle() string {
	if x != nil {
		return x.SeasonTitle
	}
	return ""
}

func (x *EpisodeResp) GetSeriesTitle() string {
	if x != nil {
		return x.SeriesTitle
	}
	return ""
}

func (x *EpisodeResp) GetActors() []*actor.ActorResp {
	if x != nil {
		return x.Actors
	}
	return nil
}

func (x *EpisodeResp) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
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
	0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67,
	0x65, 0x6e, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x12, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x73, 0x74, 0x75, 0x64, 0x69, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x13,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x14, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x75, 0x62, 0x74, 0x69,
	0x74, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe8, 0x05, 0x0a,
	0x0b, 0x45, 0x70, 0x69, 0x73, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x73, 0x70, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07,
	0x65, 0x70, 0x69, 0x73, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x65,
	0x70, 0x69, 0x73, 0x6f, 0x64, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x6f, 0x77, 0x6e,
	0x6c, 0x6f, 0x61, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x64, 0x6f, 0x77, 0x6e,
	0x6c, 0x6f, 0x61, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x78, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x65, 0x78, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x53, 0x69,
	0x7a, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x53, 0x69,
	0x7a, 0x65, 0x12, 0x32, 0x0a, 0x09, 0x73, 0x75, 0x62, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x73, 0x18,
	0x07, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6e, 0x62, 0x6f, 0x6f, 0x74, 0x2e, 0x53,
	0x75, 0x62, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x52, 0x09, 0x73, 0x75, 0x62,
	0x74, 0x69, 0x74, 0x6c, 0x65, 0x73, 0x12, 0x2e, 0x0a, 0x12, 0x6c, 0x61, 0x73, 0x74, 0x50, 0x6c,
	0x61, 0x79, 0x65, 0x64, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x12, 0x6c, 0x61, 0x73, 0x74, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x64, 0x50, 0x6f,
	0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x42, 0x0a, 0x0e, 0x6c, 0x61, 0x73, 0x74, 0x50, 0x6c,
	0x61, 0x79, 0x65, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0e, 0x6c, 0x61, 0x73, 0x74,
	0x50, 0x6c, 0x61, 0x79, 0x65, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x6b,
	0x69, 0x70, 0x49, 0x6e, 0x74, 0x72, 0x6f, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x73,
	0x6b, 0x69, 0x70, 0x49, 0x6e, 0x74, 0x72, 0x6f, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x6b, 0x69, 0x70,
	0x45, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x73, 0x6b,
	0x69, 0x70, 0x45, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c,
	0x65, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x70, 0x6f, 0x73, 0x74, 0x65, 0x72, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x70, 0x6f, 0x73, 0x74, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x6c, 0x6f, 0x67, 0x6f, 0x18, 0x0e,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6c, 0x6f, 0x67, 0x6f, 0x12, 0x34, 0x0a, 0x07, 0x61, 0x69,
	0x72, 0x44, 0x61, 0x74, 0x65, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x61, 0x69, 0x72, 0x44, 0x61, 0x74, 0x65,
	0x12, 0x1a, 0x0a, 0x08, 0x6f, 0x76, 0x65, 0x72, 0x76, 0x69, 0x65, 0x77, 0x18, 0x10, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x6f, 0x76, 0x65, 0x72, 0x76, 0x69, 0x65, 0x77, 0x12, 0x1a, 0x0a, 0x08,
	0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x18, 0x11, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08,
	0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x61, 0x73,
	0x6f, 0x6e, 0x49, 0x64, 0x18, 0x12, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x73, 0x65, 0x61, 0x73,
	0x6f, 0x6e, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x18, 0x13,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x73, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12, 0x20, 0x0a, 0x0b,
	0x73, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x14, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x73, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x20,
	0x0a, 0x0b, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x15, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x54, 0x69, 0x74, 0x6c, 0x65,
	0x12, 0x29, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x18, 0x16, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x11, 0x2e, 0x67, 0x6e, 0x62, 0x6f, 0x6f, 0x74, 0x2e, 0x41, 0x63, 0x74, 0x6f, 0x72, 0x52,
	0x65, 0x73, 0x70, 0x52, 0x06, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x66,
	0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x17, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66,
	0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x23, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x45, 0x70,
	0x69, 0x73, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x32, 0x6f, 0x0a, 0x14,
	0x45, 0x70, 0x69, 0x73, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x57, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x45, 0x70, 0x69, 0x73, 0x6f,
	0x64, 0x65, 0x12, 0x19, 0x2e, 0x67, 0x6e, 0x62, 0x6f, 0x6f, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x45,
	0x70, 0x69, 0x73, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e,
	0x67, 0x6e, 0x62, 0x6f, 0x6f, 0x74, 0x2e, 0x45, 0x70, 0x69, 0x73, 0x6f, 0x64, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x22, 0x19, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x13, 0x12, 0x11, 0x2f, 0x65, 0x70, 0x69,
	0x73, 0x6f, 0x64, 0x65, 0x2f, 0x71, 0x75, 0x65, 0x72, 0x79, 0x2f, 0x69, 0x64, 0x42, 0x9d, 0x01,
	0xba, 0x47, 0x6e, 0x12, 0x36, 0x0a, 0x0f, 0x65, 0x70, 0x69, 0x73, 0x6f, 0x64, 0x65, 0x20, 0x61,
	0x64, 0x61, 0x70, 0x74, 0x6f, 0x72, 0x12, 0x1c, 0x54, 0x68, 0x69, 0x73, 0x20, 0x69, 0x73, 0x20,
	0x65, 0x70, 0x69, 0x73, 0x6f, 0x64, 0x65, 0x20, 0x61, 0x64, 0x61, 0x70, 0x74, 0x6f, 0x72, 0x20,
	0x64, 0x6f, 0x63, 0x73, 0x32, 0x05, 0x31, 0x2e, 0x30, 0x2e, 0x30, 0x2a, 0x22, 0x3a, 0x20, 0x0a,
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
	(*EpisodeResp)(nil),           // 0: gnboot.EpisodeResp
	(*GetEpisodeRequest)(nil),     // 1: gnboot.GetEpisodeRequest
	(*subtitle.SubtitleResp)(nil), // 2: gnboot.SubtitleResp
	(*timestamppb.Timestamp)(nil), // 3: google.protobuf.Timestamp
	(*actor.ActorResp)(nil),       // 4: gnboot.ActorResp
}
var file_proto_episode_proto_depIdxs = []int32{
	2, // 0: gnboot.EpisodeResp.subtitles:type_name -> gnboot.SubtitleResp
	3, // 1: gnboot.EpisodeResp.lastPlayedTime:type_name -> google.protobuf.Timestamp
	3, // 2: gnboot.EpisodeResp.airDate:type_name -> google.protobuf.Timestamp
	4, // 3: gnboot.EpisodeResp.actors:type_name -> gnboot.ActorResp
	1, // 4: gnboot.EpisodeRemoteService.GetEpisode:input_type -> gnboot.GetEpisodeRequest
	0, // 5: gnboot.EpisodeRemoteService.GetEpisode:output_type -> gnboot.EpisodeResp
	5, // [5:6] is the sub-list for method output_type
	4, // [4:5] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_proto_episode_proto_init() }
func file_proto_episode_proto_init() {
	if File_proto_episode_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_episode_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EpisodeResp); i {
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
