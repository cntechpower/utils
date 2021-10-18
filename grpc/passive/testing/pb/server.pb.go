// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pb/server.proto

package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type PingReq struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PingReq) Reset()         { *m = PingReq{} }
func (m *PingReq) String() string { return proto.CompactTextString(m) }
func (*PingReq) ProtoMessage()    {}
func (*PingReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_server_95f8eaab15ade0d8, []int{0}
}
func (m *PingReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PingReq.Unmarshal(m, b)
}
func (m *PingReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PingReq.Marshal(b, m, deterministic)
}
func (dst *PingReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PingReq.Merge(dst, src)
}
func (m *PingReq) XXX_Size() int {
	return xxx_messageInfo_PingReq.Size(m)
}
func (m *PingReq) XXX_DiscardUnknown() {
	xxx_messageInfo_PingReq.DiscardUnknown(m)
}

var xxx_messageInfo_PingReq proto.InternalMessageInfo

type PingResp struct {
	HostName             string   `protobuf:"bytes,1,opt,name=HostName,proto3" json:"HostName,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PingResp) Reset()         { *m = PingResp{} }
func (m *PingResp) String() string { return proto.CompactTextString(m) }
func (*PingResp) ProtoMessage()    {}
func (*PingResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_server_95f8eaab15ade0d8, []int{1}
}
func (m *PingResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PingResp.Unmarshal(m, b)
}
func (m *PingResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PingResp.Marshal(b, m, deterministic)
}
func (dst *PingResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PingResp.Merge(dst, src)
}
func (m *PingResp) XXX_Size() int {
	return xxx_messageInfo_PingResp.Size(m)
}
func (m *PingResp) XXX_DiscardUnknown() {
	xxx_messageInfo_PingResp.DiscardUnknown(m)
}

var xxx_messageInfo_PingResp proto.InternalMessageInfo

func (m *PingResp) GetHostName() string {
	if m != nil {
		return m.HostName
	}
	return ""
}

func init() {
	proto.RegisterType((*PingReq)(nil), "pb.PingReq")
	proto.RegisterType((*PingResp)(nil), "pb.PingResp")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ServiceClient is the client API for Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ServiceClient interface {
	Ping(ctx context.Context, in *PingReq, opts ...grpc.CallOption) (*PingResp, error)
}

type serviceClient struct {
	cc *grpc.ClientConn
}

func NewServiceClient(cc *grpc.ClientConn) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) Ping(ctx context.Context, in *PingReq, opts ...grpc.CallOption) (*PingResp, error) {
	out := new(PingResp)
	err := c.cc.Invoke(ctx, "/pb.Service/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServiceServer is the server API for Service service.
type ServiceServer interface {
	Ping(context.Context, *PingReq) (*PingResp, error)
}

func RegisterServiceServer(s *grpc.Server, srv ServiceServer) {
	s.RegisterService(&_Service_serviceDesc, srv)
}

func _Service_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Service/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).Ping(ctx, req.(*PingReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _Service_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Service",
	HandlerType: (*ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Service_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pb/server.proto",
}

func init() { proto.RegisterFile("pb/server.proto", fileDescriptor_server_95f8eaab15ade0d8) }

var fileDescriptor_server_95f8eaab15ade0d8 = []byte{
	// 123 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2f, 0x48, 0xd2, 0x2f,
	0x4e, 0x2d, 0x2a, 0x4b, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a, 0x48, 0x52,
	0xe2, 0xe4, 0x62, 0x0f, 0xc8, 0xcc, 0x4b, 0x0f, 0x4a, 0x2d, 0x54, 0x52, 0xe3, 0xe2, 0x80, 0x30,
	0x8b, 0x0b, 0x84, 0xa4, 0xb8, 0x38, 0x3c, 0xf2, 0x8b, 0x4b, 0xfc, 0x12, 0x73, 0x53, 0x25, 0x18,
	0x15, 0x18, 0x35, 0x38, 0x83, 0xe0, 0x7c, 0x23, 0x3d, 0x2e, 0xf6, 0xe0, 0xd4, 0xa2, 0xb2, 0xcc,
	0xe4, 0x54, 0x21, 0x65, 0x2e, 0x16, 0x90, 0x16, 0x21, 0x6e, 0xbd, 0x82, 0x24, 0x3d, 0xa8, 0x39,
	0x52, 0x3c, 0x08, 0x4e, 0x71, 0x81, 0x12, 0x43, 0x12, 0x1b, 0xd8, 0x36, 0x63, 0x40, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x75, 0xce, 0xb8, 0xf1, 0x80, 0x00, 0x00, 0x00,
}
