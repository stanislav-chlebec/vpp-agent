// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: model.proto

package model

import proto "github.com/gogo/protobuf/proto"
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
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type Interface struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Interface) Reset()         { *m = Interface{} }
func (m *Interface) String() string { return proto.CompactTextString(m) }
func (*Interface) ProtoMessage()    {}
func (*Interface) Descriptor() ([]byte, []int) {
	return fileDescriptor_model_ffb99e66e52ff107, []int{0}
}
func (m *Interface) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Interface.Unmarshal(m, b)
}
func (m *Interface) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Interface.Marshal(b, m, deterministic)
}
func (dst *Interface) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Interface.Merge(dst, src)
}
func (m *Interface) XXX_Size() int {
	return xxx_messageInfo_Interface.Size(m)
}
func (m *Interface) XXX_DiscardUnknown() {
	xxx_messageInfo_Interface.DiscardUnknown(m)
}

var xxx_messageInfo_Interface proto.InternalMessageInfo

func (m *Interface) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type Route struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	InterfaceName        string   `protobuf:"bytes,2,opt,name=interface_name,json=interfaceName,proto3" json:"interface_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Route) Reset()         { *m = Route{} }
func (m *Route) String() string { return proto.CompactTextString(m) }
func (*Route) ProtoMessage()    {}
func (*Route) Descriptor() ([]byte, []int) {
	return fileDescriptor_model_ffb99e66e52ff107, []int{1}
}
func (m *Route) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Route.Unmarshal(m, b)
}
func (m *Route) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Route.Marshal(b, m, deterministic)
}
func (dst *Route) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Route.Merge(dst, src)
}
func (m *Route) XXX_Size() int {
	return xxx_messageInfo_Route.Size(m)
}
func (m *Route) XXX_DiscardUnknown() {
	xxx_messageInfo_Route.DiscardUnknown(m)
}

var xxx_messageInfo_Route proto.InternalMessageInfo

func (m *Route) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Route) GetInterfaceName() string {
	if m != nil {
		return m.InterfaceName
	}
	return ""
}

func init() {
	proto.RegisterType((*Interface)(nil), "model.Interface")
	proto.RegisterType((*Route)(nil), "model.Route")
}

func init() { proto.RegisterFile("model.proto", fileDescriptor_model_ffb99e66e52ff107) }

var fileDescriptor_model_ffb99e66e52ff107 = []byte{
	// 105 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xce, 0xcd, 0x4f, 0x49,
	0xcd, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x73, 0x94, 0xe4, 0xb9, 0x38, 0x3d,
	0xf3, 0x4a, 0x52, 0x8b, 0xd2, 0x12, 0x93, 0x53, 0x85, 0x84, 0xb8, 0x58, 0xf2, 0x12, 0x73, 0x53,
	0x25, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0xc0, 0x6c, 0x25, 0x27, 0x2e, 0xd6, 0xa0, 0xfc, 0xd2,
	0x12, 0xac, 0x92, 0x42, 0xaa, 0x5c, 0x7c, 0x99, 0x30, 0xdd, 0xf1, 0x60, 0x59, 0x26, 0xb0, 0x2c,
	0x2f, 0x5c, 0xd4, 0x2f, 0x31, 0x37, 0x35, 0x89, 0x0d, 0x6c, 0xa5, 0x31, 0x20, 0x00, 0x00, 0xff,
	0xff, 0xc8, 0x22, 0xca, 0xe4, 0x81, 0x00, 0x00, 0x00,
}
