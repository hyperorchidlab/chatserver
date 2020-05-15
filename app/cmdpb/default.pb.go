// Code generated by protoc-gen-go. DO NOT EDIT.
// source: default.proto

package cmdpb // import "."

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type DefaultRequest struct {
	Reqid                int32    `protobuf:"varint,1,opt,name=reqid,proto3" json:"reqid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DefaultRequest) Reset()         { *m = DefaultRequest{} }
func (m *DefaultRequest) String() string { return proto.CompactTextString(m) }
func (*DefaultRequest) ProtoMessage()    {}
func (*DefaultRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_default_cb6cdcc1a5baef95, []int{0}
}
func (m *DefaultRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DefaultRequest.Unmarshal(m, b)
}
func (m *DefaultRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DefaultRequest.Marshal(b, m, deterministic)
}
func (dst *DefaultRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DefaultRequest.Merge(dst, src)
}
func (m *DefaultRequest) XXX_Size() int {
	return xxx_messageInfo_DefaultRequest.Size(m)
}
func (m *DefaultRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DefaultRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DefaultRequest proto.InternalMessageInfo

func (m *DefaultRequest) GetReqid() int32 {
	if m != nil {
		return m.Reqid
	}
	return 0
}

type DefaultRequestMsg struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DefaultRequestMsg) Reset()         { *m = DefaultRequestMsg{} }
func (m *DefaultRequestMsg) String() string { return proto.CompactTextString(m) }
func (*DefaultRequestMsg) ProtoMessage()    {}
func (*DefaultRequestMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_default_cb6cdcc1a5baef95, []int{1}
}
func (m *DefaultRequestMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DefaultRequestMsg.Unmarshal(m, b)
}
func (m *DefaultRequestMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DefaultRequestMsg.Marshal(b, m, deterministic)
}
func (dst *DefaultRequestMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DefaultRequestMsg.Merge(dst, src)
}
func (m *DefaultRequestMsg) XXX_Size() int {
	return xxx_messageInfo_DefaultRequestMsg.Size(m)
}
func (m *DefaultRequestMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_DefaultRequestMsg.DiscardUnknown(m)
}

var xxx_messageInfo_DefaultRequestMsg proto.InternalMessageInfo

func (m *DefaultRequestMsg) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type DefaultResp struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DefaultResp) Reset()         { *m = DefaultResp{} }
func (m *DefaultResp) String() string { return proto.CompactTextString(m) }
func (*DefaultResp) ProtoMessage()    {}
func (*DefaultResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_default_cb6cdcc1a5baef95, []int{2}
}
func (m *DefaultResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DefaultResp.Unmarshal(m, b)
}
func (m *DefaultResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DefaultResp.Marshal(b, m, deterministic)
}
func (dst *DefaultResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DefaultResp.Merge(dst, src)
}
func (m *DefaultResp) XXX_Size() int {
	return xxx_messageInfo_DefaultResp.Size(m)
}
func (m *DefaultResp) XXX_DiscardUnknown() {
	xxx_messageInfo_DefaultResp.DiscardUnknown(m)
}

var xxx_messageInfo_DefaultResp proto.InternalMessageInfo

func (m *DefaultResp) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type DefaultRequestIDMsg struct {
	Reqid                int32    `protobuf:"varint,1,opt,name=reqid,proto3" json:"reqid,omitempty"`
	Message              string   `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DefaultRequestIDMsg) Reset()         { *m = DefaultRequestIDMsg{} }
func (m *DefaultRequestIDMsg) String() string { return proto.CompactTextString(m) }
func (*DefaultRequestIDMsg) ProtoMessage()    {}
func (*DefaultRequestIDMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_default_cb6cdcc1a5baef95, []int{3}
}
func (m *DefaultRequestIDMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DefaultRequestIDMsg.Unmarshal(m, b)
}
func (m *DefaultRequestIDMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DefaultRequestIDMsg.Marshal(b, m, deterministic)
}
func (dst *DefaultRequestIDMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DefaultRequestIDMsg.Merge(dst, src)
}
func (m *DefaultRequestIDMsg) XXX_Size() int {
	return xxx_messageInfo_DefaultRequestIDMsg.Size(m)
}
func (m *DefaultRequestIDMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_DefaultRequestIDMsg.DiscardUnknown(m)
}

var xxx_messageInfo_DefaultRequestIDMsg proto.InternalMessageInfo

func (m *DefaultRequestIDMsg) GetReqid() int32 {
	if m != nil {
		return m.Reqid
	}
	return 0
}

func (m *DefaultRequestIDMsg) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*DefaultRequest)(nil), "DefaultRequest")
	proto.RegisterType((*DefaultRequestMsg)(nil), "DefaultRequestMsg")
	proto.RegisterType((*DefaultResp)(nil), "DefaultResp")
	proto.RegisterType((*DefaultRequestIDMsg)(nil), "DefaultRequestIDMsg")
}

func init() { proto.RegisterFile("default.proto", fileDescriptor_default_cb6cdcc1a5baef95) }

var fileDescriptor_default_cb6cdcc1a5baef95 = []byte{
	// 140 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4d, 0x49, 0x4d, 0x4b,
	0x2c, 0xcd, 0x29, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x57, 0x52, 0xe3, 0xe2, 0x73, 0x81, 0x08,
	0x04, 0xa5, 0x16, 0x96, 0xa6, 0x16, 0x97, 0x08, 0x89, 0x70, 0xb1, 0x16, 0xa5, 0x16, 0x66, 0xa6,
	0x48, 0x30, 0x2a, 0x30, 0x6a, 0xb0, 0x06, 0x41, 0x38, 0x4a, 0xba, 0x5c, 0x82, 0xa8, 0xea, 0x7c,
	0x8b, 0xd3, 0x85, 0x24, 0xb8, 0xd8, 0x73, 0x53, 0x8b, 0x8b, 0x13, 0xd3, 0x53, 0xc1, 0x8a, 0x39,
	0x83, 0x60, 0x5c, 0x25, 0x75, 0x2e, 0x6e, 0xb8, 0xf2, 0xe2, 0x02, 0x3c, 0x0a, 0x5d, 0xb9, 0x84,
	0x51, 0xcd, 0xf5, 0x74, 0x01, 0x99, 0x8c, 0xd5, 0x11, 0xc8, 0xc6, 0x30, 0xa1, 0x18, 0xe3, 0xc4,
	0x19, 0xc5, 0xae, 0x67, 0x9d, 0x9c, 0x9b, 0x52, 0x90, 0x94, 0xc4, 0x06, 0xf6, 0x98, 0x31, 0x20,
	0x00, 0x00, 0xff, 0xff, 0x85, 0xd3, 0xe5, 0x17, 0xe9, 0x00, 0x00, 0x00,
}