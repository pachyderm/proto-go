// Code generated by protoc-gen-go.
// source: http/protohttp.proto
// DO NOT EDIT!

package protohttp

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto.ProtoPackageIsVersion1

type BasicAuth struct {
	Username string `protobuf:"bytes,1,opt,name=username" json:"username,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password" json:"password,omitempty"`
}

func (m *BasicAuth) Reset()                    { *m = BasicAuth{} }
func (m *BasicAuth) String() string            { return proto.CompactTextString(m) }
func (*BasicAuth) ProtoMessage()               {}
func (*BasicAuth) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func init() {
	proto.RegisterType((*BasicAuth)(nil), "protohttp.BasicAuth")
}

var fileDescriptor0 = []byte{
	// 94 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x12, 0xc9, 0x28, 0x29, 0x29,
	0xd0, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0x07, 0xb1, 0xf4, 0xc0, 0x2c, 0x21, 0x4e, 0xb8, 0x80, 0x92,
	0x3e, 0x17, 0xa7, 0x53, 0x62, 0x71, 0x66, 0xb2, 0x63, 0x69, 0x49, 0x86, 0x90, 0x00, 0x17, 0x47,
	0x69, 0x71, 0x6a, 0x51, 0x5e, 0x62, 0x6e, 0xaa, 0x04, 0xa3, 0x02, 0xa3, 0x06, 0x27, 0x48, 0xa4,
	0x20, 0xb1, 0xb8, 0xb8, 0x3c, 0xbf, 0x28, 0x45, 0x82, 0x09, 0x24, 0x92, 0xc4, 0x06, 0xd6, 0x6b,
	0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0xf6, 0x38, 0x07, 0xe5, 0x5a, 0x00, 0x00, 0x00,
}
