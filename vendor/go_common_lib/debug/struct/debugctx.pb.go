// Code generated by protoc-gen-go. DO NOT EDIT.
// source: debugctx.proto

package proto_debugctx

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type DebugCtxMessage struct {
	Uid                  int32       `protobuf:"varint,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Timestamp            int32       `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	SupplementDatas      []*ListData `protobuf:"bytes,3,rep,name=supplement_datas,json=supplementDatas,proto3" json:"supplement_datas,omitempty"`
	MergeData            *ListData   `protobuf:"bytes,4,opt,name=merge_data,json=mergeData,proto3" json:"merge_data,omitempty"`
	ScatterDatas         []*ListData `protobuf:"bytes,5,rep,name=scatter_datas,json=scatterDatas,proto3" json:"scatter_datas,omitempty"`
	InterposeDatas       []*ListData `protobuf:"bytes,6,rep,name=interpose_datas,json=interposeDatas,proto3" json:"interpose_datas,omitempty"`
	OutputData           *ListData   `protobuf:"bytes,7,opt,name=output_data,json=outputData,proto3" json:"output_data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *DebugCtxMessage) Reset()         { *m = DebugCtxMessage{} }
func (m *DebugCtxMessage) String() string { return proto.CompactTextString(m) }
func (*DebugCtxMessage) ProtoMessage()    {}
func (*DebugCtxMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_f43e13008befc6e5, []int{0}
}

func (m *DebugCtxMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DebugCtxMessage.Unmarshal(m, b)
}
func (m *DebugCtxMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DebugCtxMessage.Marshal(b, m, deterministic)
}
func (m *DebugCtxMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DebugCtxMessage.Merge(m, src)
}
func (m *DebugCtxMessage) XXX_Size() int {
	return xxx_messageInfo_DebugCtxMessage.Size(m)
}
func (m *DebugCtxMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_DebugCtxMessage.DiscardUnknown(m)
}

var xxx_messageInfo_DebugCtxMessage proto.InternalMessageInfo

func (m *DebugCtxMessage) GetUid() int32 {
	if m != nil {
		return m.Uid
	}
	return 0
}

func (m *DebugCtxMessage) GetTimestamp() int32 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *DebugCtxMessage) GetSupplementDatas() []*ListData {
	if m != nil {
		return m.SupplementDatas
	}
	return nil
}

func (m *DebugCtxMessage) GetMergeData() *ListData {
	if m != nil {
		return m.MergeData
	}
	return nil
}

func (m *DebugCtxMessage) GetScatterDatas() []*ListData {
	if m != nil {
		return m.ScatterDatas
	}
	return nil
}

func (m *DebugCtxMessage) GetInterposeDatas() []*ListData {
	if m != nil {
		return m.InterposeDatas
	}
	return nil
}

func (m *DebugCtxMessage) GetOutputData() *ListData {
	if m != nil {
		return m.OutputData
	}
	return nil
}

type ListData struct {
	PluginName           string      `protobuf:"bytes,1,opt,name=plugin_name,json=pluginName,proto3" json:"plugin_name,omitempty"`
	LiveList             []*LiveInfo `protobuf:"bytes,2,rep,name=live_list,json=liveList,proto3" json:"live_list,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *ListData) Reset()         { *m = ListData{} }
func (m *ListData) String() string { return proto.CompactTextString(m) }
func (*ListData) ProtoMessage()    {}
func (*ListData) Descriptor() ([]byte, []int) {
	return fileDescriptor_f43e13008befc6e5, []int{1}
}

func (m *ListData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListData.Unmarshal(m, b)
}
func (m *ListData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListData.Marshal(b, m, deterministic)
}
func (m *ListData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListData.Merge(m, src)
}
func (m *ListData) XXX_Size() int {
	return xxx_messageInfo_ListData.Size(m)
}
func (m *ListData) XXX_DiscardUnknown() {
	xxx_messageInfo_ListData.DiscardUnknown(m)
}

var xxx_messageInfo_ListData proto.InternalMessageInfo

func (m *ListData) GetPluginName() string {
	if m != nil {
		return m.PluginName
	}
	return ""
}

func (m *ListData) GetLiveList() []*LiveInfo {
	if m != nil {
		return m.LiveList
	}
	return nil
}

type LiveInfo struct {
	Uid                  int32    `protobuf:"varint,1,opt,name=uid,proto3" json:"uid,omitempty"`
	LiveId               string   `protobuf:"bytes,2,opt,name=live_id,json=liveId,proto3" json:"live_id,omitempty"`
	Reason               string   `protobuf:"bytes,3,opt,name=reason,proto3" json:"reason,omitempty"`
	Distance             float32  `protobuf:"fixed32,4,opt,name=distance,proto3" json:"distance,omitempty"`
	Token                string   `protobuf:"bytes,5,opt,name=token,proto3" json:"token,omitempty"`
	FilterName           string   `protobuf:"bytes,6,opt,name=filter_name,json=filterName,proto3" json:"filter_name,omitempty"`
	Appearance           string   `protobuf:"bytes,7,opt,name=appearance,proto3" json:"appearance,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LiveInfo) Reset()         { *m = LiveInfo{} }
func (m *LiveInfo) String() string { return proto.CompactTextString(m) }
func (*LiveInfo) ProtoMessage()    {}
func (*LiveInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_f43e13008befc6e5, []int{2}
}

func (m *LiveInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LiveInfo.Unmarshal(m, b)
}
func (m *LiveInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LiveInfo.Marshal(b, m, deterministic)
}
func (m *LiveInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LiveInfo.Merge(m, src)
}
func (m *LiveInfo) XXX_Size() int {
	return xxx_messageInfo_LiveInfo.Size(m)
}
func (m *LiveInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_LiveInfo.DiscardUnknown(m)
}

var xxx_messageInfo_LiveInfo proto.InternalMessageInfo

func (m *LiveInfo) GetUid() int32 {
	if m != nil {
		return m.Uid
	}
	return 0
}

func (m *LiveInfo) GetLiveId() string {
	if m != nil {
		return m.LiveId
	}
	return ""
}

func (m *LiveInfo) GetReason() string {
	if m != nil {
		return m.Reason
	}
	return ""
}

func (m *LiveInfo) GetDistance() float32 {
	if m != nil {
		return m.Distance
	}
	return 0
}

func (m *LiveInfo) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *LiveInfo) GetFilterName() string {
	if m != nil {
		return m.FilterName
	}
	return ""
}

func (m *LiveInfo) GetAppearance() string {
	if m != nil {
		return m.Appearance
	}
	return ""
}

func init() {
	proto.RegisterType((*DebugCtxMessage)(nil), "mypackage.DebugCtxMessage")
	proto.RegisterType((*ListData)(nil), "mypackage.ListData")
	proto.RegisterType((*LiveInfo)(nil), "mypackage.LiveInfo")
}

func init() { proto.RegisterFile("debugctx.proto", fileDescriptor_f43e13008befc6e5) }

var fileDescriptor_f43e13008befc6e5 = []byte{
	// 405 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x52, 0x3d, 0x8f, 0xd3, 0x40,
	0x10, 0x55, 0x12, 0xe2, 0xc4, 0x13, 0x48, 0x4e, 0x0b, 0x82, 0x15, 0x42, 0x10, 0xa5, 0xba, 0xca,
	0x42, 0x81, 0x82, 0x02, 0x51, 0xc0, 0x35, 0x27, 0x01, 0x85, 0x4b, 0x24, 0x14, 0xcd, 0xd9, 0x73,
	0x66, 0x15, 0xef, 0x87, 0xbc, 0xe3, 0xe8, 0xf8, 0x77, 0xfc, 0x17, 0xfe, 0xc8, 0x69, 0x77, 0x73,
	0x49, 0x73, 0xee, 0xe6, 0x3d, 0xbf, 0x37, 0x33, 0xcf, 0xb3, 0xb0, 0xac, 0xe9, 0xa6, 0x6f, 0x2a,
	0xbe, 0x2b, 0x5c, 0x67, 0xd9, 0x8a, 0x5c, 0xff, 0x75, 0x58, 0xed, 0xb1, 0xa1, 0xcd, 0xff, 0x31,
	0xac, 0xae, 0xc2, 0xd7, 0x6f, 0x7c, 0xf7, 0x83, 0xbc, 0xc7, 0x86, 0xc4, 0x05, 0x4c, 0x7a, 0x55,
	0xcb, 0xd1, 0x7a, 0x74, 0x39, 0x2d, 0x43, 0x29, 0xde, 0x40, 0xce, 0x4a, 0x93, 0x67, 0xd4, 0x4e,
	0x8e, 0x23, 0x7f, 0x26, 0xc4, 0x17, 0xb8, 0xf0, 0xbd, 0x73, 0x2d, 0x69, 0x32, 0xbc, 0xab, 0x91,
	0xd1, 0xcb, 0xc9, 0x7a, 0x72, 0xb9, 0xd8, 0x3e, 0x2f, 0x4e, 0x93, 0x8a, 0xef, 0xca, 0xf3, 0x15,
	0x32, 0x96, 0xab, 0xb3, 0x38, 0x60, 0x2f, 0xb6, 0x00, 0x9a, 0xba, 0x86, 0xa2, 0x55, 0x3e, 0x59,
	0x8f, 0x86, 0x9c, 0x79, 0x94, 0x85, 0x52, 0x7c, 0x82, 0x67, 0xbe, 0x42, 0x66, 0xea, 0x8e, 0x03,
	0xa7, 0xc3, 0x03, 0x9f, 0x1e, 0x95, 0x69, 0xda, 0x67, 0x58, 0x29, 0xc3, 0xd4, 0x39, 0xeb, 0xe9,
	0xe8, 0xcd, 0x86, 0xbd, 0xcb, 0x93, 0x36, 0xb9, 0x3f, 0xc2, 0xc2, 0xf6, 0xec, 0xfa, 0x94, 0x53,
	0xce, 0x86, 0x97, 0x85, 0xa4, 0x0b, 0xf5, 0xe6, 0x37, 0xcc, 0x1f, 0x78, 0xf1, 0x0e, 0x16, 0xae,
	0xed, 0x1b, 0x65, 0x76, 0x06, 0x35, 0xc5, 0xbf, 0x9c, 0x97, 0x90, 0xa8, 0x9f, 0xa8, 0x49, 0xbc,
	0x87, 0xbc, 0x55, 0x07, 0xda, 0xb5, 0xca, 0xb3, 0x1c, 0x3f, 0xb2, 0xda, 0x81, 0xae, 0xcd, 0xad,
	0x2d, 0xe7, 0x41, 0x15, 0xda, 0x6e, 0xfe, 0x8d, 0x42, 0xff, 0x44, 0x3f, 0x72, 0xbd, 0x57, 0x30,
	0x8b, 0x0d, 0x55, 0x1d, 0x6f, 0x97, 0x97, 0x59, 0x80, 0xd7, 0xb5, 0x78, 0x09, 0x59, 0x47, 0xe8,
	0xad, 0x91, 0x93, 0xc4, 0x27, 0x24, 0x5e, 0xc3, 0xbc, 0x56, 0x9e, 0xd1, 0x54, 0x14, 0xcf, 0x31,
	0x2e, 0x4f, 0x58, 0xbc, 0x80, 0x29, 0xdb, 0x3d, 0x19, 0x39, 0x8d, 0x96, 0x04, 0x42, 0xa8, 0x5b,
	0xd5, 0x86, 0x6b, 0xc4, 0x50, 0x59, 0x0a, 0x95, 0xa8, 0x18, 0xea, 0x2d, 0x00, 0x3a, 0x47, 0xd8,
	0xc5, 0xa6, 0xb3, 0xf4, 0xfd, 0xcc, 0x7c, 0xdd, 0xc2, 0xba, 0xb2, 0xba, 0x50, 0x66, 0x4f, 0x85,
	0xa7, 0xee, 0xa0, 0x2a, 0x4a, 0x8f, 0xb5, 0xf8, 0x83, 0x6d, 0x5b, 0x84, 0x85, 0x7f, 0x2d, 0x13,
	0xf1, 0xf0, 0x98, 0x6f, 0xb2, 0x88, 0x3f, 0xdc, 0x07, 0x00, 0x00, 0xff, 0xff, 0x13, 0x91, 0x32,
	0x16, 0xdf, 0x02, 0x00, 0x00,
}
