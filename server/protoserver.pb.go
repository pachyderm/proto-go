// Code generated by protoc-gen-go.
// source: server/protoserver.proto
// DO NOT EDIT!

package protoserver

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "go.pedge.io/pb/go/google/protobuf"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ServerStarted struct {
	Port     uint32 `protobuf:"varint,1,opt,name=port" json:"port,omitempty"`
	HttpPort uint32 `protobuf:"varint,2,opt,name=http_port,json=httpPort" json:"http_port,omitempty"`
}

func (m *ServerStarted) Reset()                    { *m = ServerStarted{} }
func (m *ServerStarted) String() string            { return proto.CompactTextString(m) }
func (*ServerStarted) ProtoMessage()               {}
func (*ServerStarted) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type ServerFinished struct {
	Error    string                    `protobuf:"bytes,1,opt,name=error" json:"error,omitempty"`
	Duration *google_protobuf.Duration `protobuf:"bytes,2,opt,name=duration" json:"duration,omitempty"`
}

func (m *ServerFinished) Reset()                    { *m = ServerFinished{} }
func (m *ServerFinished) String() string            { return proto.CompactTextString(m) }
func (*ServerFinished) ProtoMessage()               {}
func (*ServerFinished) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ServerFinished) GetDuration() *google_protobuf.Duration {
	if m != nil {
		return m.Duration
	}
	return nil
}

func init() {
	proto.RegisterType((*ServerStarted)(nil), "protoserver.ServerStarted")
	proto.RegisterType((*ServerFinished)(nil), "protoserver.ServerFinished")
}

func init() { proto.RegisterFile("server/protoserver.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 180 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x4c, 0x8e, 0xc1, 0xcb, 0x82, 0x40,
	0x10, 0xc5, 0xf1, 0xe3, 0x2b, 0x74, 0xc4, 0x0e, 0x4b, 0x07, 0x2b, 0x88, 0xf0, 0xd4, 0x69, 0x85,
	0xa2, 0x7b, 0x87, 0xe8, 0x1c, 0x7a, 0x8e, 0x50, 0xdc, 0x54, 0x08, 0x67, 0x19, 0xc7, 0xfe, 0xfe,
	0x68, 0x36, 0xa3, 0xdb, 0x9b, 0xf7, 0x7e, 0xbc, 0x79, 0x10, 0xf7, 0x86, 0x9e, 0x86, 0x52, 0x4b,
	0xc8, 0xe8, 0xb4, 0x16, 0xad, 0xc2, 0x1f, 0x6b, 0xb9, 0xae, 0x11, 0xeb, 0x87, 0x71, 0x58, 0x39,
	0xdc, 0xd3, 0x6a, 0xa0, 0x82, 0x5b, 0xec, 0x1c, 0x9c, 0x1c, 0x21, 0xca, 0x85, 0xcc, 0xb9, 0x20,
	0x36, 0x95, 0x52, 0xf0, 0x6f, 0x91, 0x38, 0xf6, 0x36, 0xde, 0x36, 0xca, 0x44, 0xab, 0x15, 0x04,
	0x0d, 0xb3, 0xbd, 0x49, 0xf0, 0x27, 0x81, 0xff, 0x36, 0x2e, 0x48, 0x9c, 0x5c, 0x61, 0xe6, 0x1a,
	0xce, 0x6d, 0xd7, 0xf6, 0x8d, 0xa9, 0xd4, 0x1c, 0x26, 0x86, 0x08, 0x49, 0x3a, 0x82, 0xcc, 0x1d,
	0xea, 0x00, 0xfe, 0xf8, 0x5b, 0x3a, 0xc2, 0xdd, 0x42, 0xbb, 0x71, 0x7a, 0x1c, 0xa7, 0x4f, 0x1f,
	0x20, 0xfb, 0xa2, 0xe5, 0x54, 0xc2, 0xfd, 0x2b, 0x00, 0x00, 0xff, 0xff, 0x99, 0xd1, 0x2e, 0xcb,
	0xf0, 0x00, 0x00, 0x00,
}
