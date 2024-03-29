// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: policy/query.proto

package types

import (
	bytes "bytes"
	context "context"
	fmt "fmt"
	query "github.com/cosmos/cosmos-sdk/types/query"
	_ "github.com/gogo/protobuf/gogoproto"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	github_com_tendermint_tendermint_libs_bytes "github.com/tendermint/tendermint/libs/bytes"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// QueryRegoRequest is the request type for the Query/Code RPC method
type QueryRegoRequest struct {
	RegoId uint64 `protobuf:"varint,1,opt,name=rego_id,json=regoId,proto3" json:"rego_id,omitempty"`
}

func (m *QueryRegoRequest) Reset()         { *m = QueryRegoRequest{} }
func (m *QueryRegoRequest) String() string { return proto.CompactTextString(m) }
func (*QueryRegoRequest) ProtoMessage()    {}
func (*QueryRegoRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a7e8b43f6e481a8b, []int{0}
}
func (m *QueryRegoRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryRegoRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryRegoRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryRegoRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryRegoRequest.Merge(m, src)
}
func (m *QueryRegoRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryRegoRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryRegoRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryRegoRequest proto.InternalMessageInfo

// QueryRegoResponse is the response type for the Query/Rego RPC method
type QueryRegoResponse struct {
	*RegoInfoResponse `protobuf:"bytes,1,opt,name=rego_info,json=regoInfo,proto3,embedded=rego_info" json:""`
	Data              []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data"`
}

func (m *QueryRegoResponse) Reset()         { *m = QueryRegoResponse{} }
func (m *QueryRegoResponse) String() string { return proto.CompactTextString(m) }
func (*QueryRegoResponse) ProtoMessage()    {}
func (*QueryRegoResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a7e8b43f6e481a8b, []int{1}
}
func (m *QueryRegoResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryRegoResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryRegoResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryRegoResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryRegoResponse.Merge(m, src)
}
func (m *QueryRegoResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryRegoResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryRegoResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryRegoResponse proto.InternalMessageInfo

// RegoInfoResponse contains rego meta data from RegoInfo
type RegoInfoResponse struct {
	RegoID      uint64                                               `protobuf:"varint,1,opt,name=rego_id,json=regoId,proto3" json:"id"`
	Creator     string                                               `protobuf:"bytes,2,opt,name=creator,proto3" json:"creator,omitempty"`
	RegoHash    github_com_tendermint_tendermint_libs_bytes.HexBytes `protobuf:"bytes,3,opt,name=rego_hash,json=regoHash,proto3,casttype=github.com/tendermint/tendermint/libs/bytes.HexBytes" json:"rego_hash,omitempty"`
	Source      string                                               `protobuf:"bytes,4,opt,name=source,proto3" json:"source,omitempty"`
	EntryPoints []string                                             `protobuf:"bytes,5,rep,name=entryPoints,proto3" json:"entryPoints,omitempty"`
}

func (m *RegoInfoResponse) Reset()         { *m = RegoInfoResponse{} }
func (m *RegoInfoResponse) String() string { return proto.CompactTextString(m) }
func (*RegoInfoResponse) ProtoMessage()    {}
func (*RegoInfoResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a7e8b43f6e481a8b, []int{2}
}
func (m *RegoInfoResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RegoInfoResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RegoInfoResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RegoInfoResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegoInfoResponse.Merge(m, src)
}
func (m *RegoInfoResponse) XXX_Size() int {
	return m.Size()
}
func (m *RegoInfoResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RegoInfoResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RegoInfoResponse proto.InternalMessageInfo

// QueryRegosRequest is the request type for the Query/Codes RPC method
type QueryRegosRequest struct {
	// pagination defines an optional pagination for the request.
	Pagination *query.PageRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryRegosRequest) Reset()         { *m = QueryRegosRequest{} }
func (m *QueryRegosRequest) String() string { return proto.CompactTextString(m) }
func (*QueryRegosRequest) ProtoMessage()    {}
func (*QueryRegosRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a7e8b43f6e481a8b, []int{3}
}
func (m *QueryRegosRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryRegosRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryRegosRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryRegosRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryRegosRequest.Merge(m, src)
}
func (m *QueryRegosRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryRegosRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryRegosRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryRegosRequest proto.InternalMessageInfo

// QueryCodesResponse is the response type for the Query/Regos RPC method
type QueryRegosResponse struct {
	RegoInfos []RegoInfoResponse `protobuf:"bytes,1,rep,name=rego_infos,json=regoInfos,proto3" json:"rego_infos"`
	// pagination defines the pagination in the response.
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryRegosResponse) Reset()         { *m = QueryRegosResponse{} }
func (m *QueryRegosResponse) String() string { return proto.CompactTextString(m) }
func (*QueryRegosResponse) ProtoMessage()    {}
func (*QueryRegosResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a7e8b43f6e481a8b, []int{4}
}
func (m *QueryRegosResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryRegosResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryRegosResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryRegosResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryRegosResponse.Merge(m, src)
}
func (m *QueryRegosResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryRegosResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryRegosResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryRegosResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*QueryRegoRequest)(nil), "StylusFrost.policy.policy.QueryRegoRequest")
	proto.RegisterType((*QueryRegoResponse)(nil), "StylusFrost.policy.policy.QueryRegoResponse")
	proto.RegisterType((*RegoInfoResponse)(nil), "StylusFrost.policy.policy.RegoInfoResponse")
	proto.RegisterType((*QueryRegosRequest)(nil), "StylusFrost.policy.policy.QueryRegosRequest")
	proto.RegisterType((*QueryRegosResponse)(nil), "StylusFrost.policy.policy.QueryRegosResponse")
}

func init() { proto.RegisterFile("policy/query.proto", fileDescriptor_a7e8b43f6e481a8b) }

var fileDescriptor_a7e8b43f6e481a8b = []byte{
	// 607 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0x31, 0x6f, 0xd3, 0x40,
	0x14, 0xf6, 0xa5, 0x6e, 0x9a, 0x5c, 0x2b, 0x14, 0x4e, 0x55, 0x31, 0x51, 0x65, 0x47, 0x19, 0x20,
	0xd0, 0xe2, 0x53, 0x0b, 0x03, 0x62, 0x8c, 0x50, 0x69, 0x17, 0x14, 0x8c, 0x58, 0x60, 0x40, 0xe7,
	0xe4, 0xea, 0x58, 0x4a, 0x7c, 0xae, 0xef, 0x82, 0x6a, 0x10, 0x03, 0x2c, 0x0c, 0x2c, 0x48, 0xfc,
	0x01, 0x46, 0x36, 0xfe, 0x46, 0xc6, 0x48, 0x2c, 0x4c, 0x16, 0x24, 0x0c, 0x28, 0x03, 0x3f, 0x80,
	0x09, 0xf9, 0x7c, 0x2e, 0x26, 0x12, 0x10, 0x96, 0xe4, 0x5e, 0xee, 0xfb, 0xde, 0xfb, 0xde, 0x7b,
	0x5f, 0x0e, 0xa2, 0x90, 0x0d, 0xfc, 0x6e, 0x8c, 0x4f, 0x46, 0x34, 0x8a, 0xed, 0x30, 0x62, 0x82,
	0xa1, 0x8b, 0xf7, 0x45, 0x3c, 0x18, 0xf1, 0x83, 0x88, 0x71, 0x61, 0x67, 0xf7, 0xea, 0xab, 0xbe,
	0xe9, 0x31, 0x8f, 0x49, 0x14, 0x4e, 0x4f, 0x19, 0xa1, 0xbe, 0xed, 0x31, 0xe6, 0x0d, 0x28, 0x26,
	0xa1, 0x8f, 0x49, 0x10, 0x30, 0x41, 0x84, 0xcf, 0x02, 0xae, 0x6e, 0xaf, 0x76, 0x19, 0x1f, 0x32,
	0x8e, 0x5d, 0xc2, 0x69, 0x56, 0x07, 0x3f, 0xd9, 0x73, 0xa9, 0x20, 0x7b, 0x38, 0x24, 0x9e, 0x1f,
	0x48, 0xb0, 0xc2, 0xe6, 0x72, 0x44, 0x1c, 0x52, 0xc5, 0x6f, 0xee, 0xc0, 0xda, 0xbd, 0x94, 0xe5,
	0x50, 0x8f, 0x39, 0xf4, 0x64, 0x44, 0xb9, 0x40, 0x17, 0xe0, 0x5a, 0x44, 0x3d, 0xf6, 0xd8, 0xef,
	0x19, 0xa0, 0x01, 0x5a, 0xba, 0x53, 0x4e, 0xc3, 0xa3, 0x5e, 0xf3, 0x35, 0x80, 0xe7, 0x0b, 0x68,
	0x1e, 0xb2, 0x80, 0x53, 0xe4, 0xc0, 0x6a, 0x06, 0x0f, 0x8e, 0x99, 0x24, 0xac, 0xef, 0xef, 0xd8,
	0x7f, 0xec, 0xd2, 0x4e, 0xb9, 0x47, 0xc1, 0xf1, 0x19, 0xbf, 0x5d, 0x99, 0x24, 0x16, 0x98, 0x27,
	0x96, 0xe6, 0x54, 0x22, 0x75, 0x87, 0xb6, 0xa1, 0xde, 0x23, 0x82, 0x18, 0xa5, 0x06, 0x68, 0x6d,
	0xb4, 0x2b, 0xf3, 0xc4, 0x92, 0xb1, 0x23, 0x3f, 0x6f, 0xe9, 0xdf, 0xde, 0x59, 0xa0, 0xf9, 0x1d,
	0xc0, 0xda, 0x62, 0x32, 0x74, 0x65, 0x41, 0x7b, 0xbb, 0x36, 0x4d, 0xac, 0xb2, 0x84, 0xdd, 0x9e,
	0x27, 0x56, 0xc9, 0xef, 0xe5, 0xdd, 0x20, 0x03, 0xae, 0x75, 0x23, 0x4a, 0x04, 0x8b, 0x64, 0x99,
	0xaa, 0x93, 0x87, 0xe8, 0x81, 0xea, 0xa8, 0x4f, 0x78, 0xdf, 0x58, 0x91, 0x12, 0x6e, 0xfe, 0x48,
	0xac, 0x1b, 0x9e, 0x2f, 0xfa, 0x23, 0xd7, 0xee, 0xb2, 0x21, 0x16, 0x34, 0xe8, 0xd1, 0x68, 0xe8,
	0x07, 0xa2, 0x78, 0x1c, 0xf8, 0x2e, 0xc7, 0x6e, 0x2c, 0x28, 0xb7, 0x0f, 0xe9, 0x69, 0x3b, 0x3d,
	0x64, 0x4d, 0x1d, 0x12, 0xde, 0x47, 0x5b, 0xb0, 0xcc, 0xd9, 0x28, 0xea, 0x52, 0x43, 0x97, 0xf5,
	0x54, 0x84, 0x1a, 0x70, 0x9d, 0x06, 0x22, 0x8a, 0x3b, 0xcc, 0x0f, 0x04, 0x37, 0x56, 0x1b, 0x2b,
	0xad, 0xaa, 0x53, 0xfc, 0x49, 0x35, 0xfc, 0xa8, 0x30, 0x7d, 0x9e, 0x2f, 0xeb, 0x00, 0xc2, 0x5f,
	0x8b, 0x56, 0xe3, 0xbf, 0x64, 0x67, 0xae, 0xb0, 0x53, 0x57, 0xd8, 0x99, 0xfb, 0x94, 0x2b, 0xec,
	0x0e, 0xf1, 0xa8, 0xe2, 0x3a, 0x05, 0x66, 0xf3, 0x03, 0x80, 0xa8, 0x98, 0x5d, 0xcd, 0xb3, 0x03,
	0xe1, 0xd9, 0x72, 0xb9, 0x01, 0x1a, 0x2b, 0xff, 0xbb, 0x5d, 0x7d, 0x9c, 0x6e, 0xb6, 0x9a, 0x6f,
	0x96, 0xa3, 0x3b, 0xbf, 0x09, 0x2e, 0x49, 0xc1, 0x97, 0xff, 0x29, 0x38, 0xcb, 0x56, 0x54, 0xbc,
	0xff, 0xaa, 0x04, 0x57, 0xa5, 0x62, 0xf4, 0x02, 0x40, 0x3d, 0x2d, 0x8c, 0xfe, 0xa6, 0x6c, 0xd1,
	0xe6, 0xf5, 0xdd, 0xe5, 0xc0, 0x59, 0xe5, 0xa6, 0xf9, 0xf2, 0xe3, 0xd7, 0xb7, 0x25, 0x03, 0x6d,
	0x61, 0xf5, 0x2f, 0x4a, 0x3b, 0xc2, 0xcf, 0x94, 0xd9, 0x9e, 0xa3, 0xa7, 0x70, 0x55, 0x4e, 0x0e,
	0x2d, 0x95, 0x36, 0x5f, 0x5f, 0xfd, 0xda, 0x92, 0x68, 0xa5, 0x62, 0x53, 0xaa, 0x38, 0x87, 0x36,
	0x8a, 0x2a, 0xda, 0x77, 0xc7, 0x5f, 0x4c, 0xed, 0xfd, 0xd4, 0xd4, 0xc6, 0x53, 0x13, 0x4c, 0xa6,
	0x26, 0xf8, 0x3c, 0x35, 0xc1, 0x9b, 0x99, 0xa9, 0x4d, 0x66, 0xa6, 0xf6, 0x69, 0x66, 0x6a, 0x0f,
	0x77, 0x0b, 0xd6, 0x2d, 0x14, 0xcc, 0xb3, 0x9c, 0xe2, 0xe2, 0xd3, 0xe0, 0x96, 0xe5, 0xdb, 0x70,
	0xfd, 0x67, 0x00, 0x00, 0x00, 0xff, 0xff, 0x75, 0xc2, 0x28, 0x81, 0xc0, 0x04, 0x00, 0x00,
}

func (this *QueryRegoResponse) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*QueryRegoResponse)
	if !ok {
		that2, ok := that.(QueryRegoResponse)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !this.RegoInfoResponse.Equal(that1.RegoInfoResponse) {
		return false
	}
	if !bytes.Equal(this.Data, that1.Data) {
		return false
	}
	return true
}
func (this *RegoInfoResponse) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*RegoInfoResponse)
	if !ok {
		that2, ok := that.(RegoInfoResponse)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.RegoID != that1.RegoID {
		return false
	}
	if this.Creator != that1.Creator {
		return false
	}
	if !bytes.Equal(this.RegoHash, that1.RegoHash) {
		return false
	}
	if this.Source != that1.Source {
		return false
	}
	if len(this.EntryPoints) != len(that1.EntryPoints) {
		return false
	}
	for i := range this.EntryPoints {
		if this.EntryPoints[i] != that1.EntryPoints[i] {
			return false
		}
	}
	return true
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	// Rego gets the binary code and metadata for a singe rego code
	Rego(ctx context.Context, in *QueryRegoRequest, opts ...grpc.CallOption) (*QueryRegoResponse, error)
	// Regos gets the metadata for all stored rego codes
	Regos(ctx context.Context, in *QueryRegosRequest, opts ...grpc.CallOption) (*QueryRegosResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Rego(ctx context.Context, in *QueryRegoRequest, opts ...grpc.CallOption) (*QueryRegoResponse, error) {
	out := new(QueryRegoResponse)
	err := c.cc.Invoke(ctx, "/StylusFrost.policy.policy.Query/Rego", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Regos(ctx context.Context, in *QueryRegosRequest, opts ...grpc.CallOption) (*QueryRegosResponse, error) {
	out := new(QueryRegosResponse)
	err := c.cc.Invoke(ctx, "/StylusFrost.policy.policy.Query/Regos", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Rego gets the binary code and metadata for a singe rego code
	Rego(context.Context, *QueryRegoRequest) (*QueryRegoResponse, error)
	// Regos gets the metadata for all stored rego codes
	Regos(context.Context, *QueryRegosRequest) (*QueryRegosResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) Rego(ctx context.Context, req *QueryRegoRequest) (*QueryRegoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Rego not implemented")
}
func (*UnimplementedQueryServer) Regos(ctx context.Context, req *QueryRegosRequest) (*QueryRegosResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Regos not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_Rego_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryRegoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Rego(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/StylusFrost.policy.policy.Query/Rego",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Rego(ctx, req.(*QueryRegoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Regos_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryRegosRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Regos(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/StylusFrost.policy.policy.Query/Regos",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Regos(ctx, req.(*QueryRegosRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "StylusFrost.policy.policy.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Rego",
			Handler:    _Query_Rego_Handler,
		},
		{
			MethodName: "Regos",
			Handler:    _Query_Regos_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "policy/query.proto",
}

func (m *QueryRegoRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryRegoRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryRegoRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.RegoId != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.RegoId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *QueryRegoResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryRegoResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryRegoResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Data) > 0 {
		i -= len(m.Data)
		copy(dAtA[i:], m.Data)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Data)))
		i--
		dAtA[i] = 0x12
	}
	if m.RegoInfoResponse != nil {
		{
			size, err := m.RegoInfoResponse.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *RegoInfoResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RegoInfoResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RegoInfoResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.EntryPoints) > 0 {
		for iNdEx := len(m.EntryPoints) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.EntryPoints[iNdEx])
			copy(dAtA[i:], m.EntryPoints[iNdEx])
			i = encodeVarintQuery(dAtA, i, uint64(len(m.EntryPoints[iNdEx])))
			i--
			dAtA[i] = 0x2a
		}
	}
	if len(m.Source) > 0 {
		i -= len(m.Source)
		copy(dAtA[i:], m.Source)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Source)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.RegoHash) > 0 {
		i -= len(m.RegoHash)
		copy(dAtA[i:], m.RegoHash)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.RegoHash)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0x12
	}
	if m.RegoID != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.RegoID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *QueryRegosRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryRegosRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryRegosRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Pagination != nil {
		{
			size, err := m.Pagination.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryRegosResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryRegosResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryRegosResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Pagination != nil {
		{
			size, err := m.Pagination.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.RegoInfos) > 0 {
		for iNdEx := len(m.RegoInfos) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.RegoInfos[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintQuery(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryRegoRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.RegoId != 0 {
		n += 1 + sovQuery(uint64(m.RegoId))
	}
	return n
}

func (m *QueryRegoResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.RegoInfoResponse != nil {
		l = m.RegoInfoResponse.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	l = len(m.Data)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *RegoInfoResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.RegoID != 0 {
		n += 1 + sovQuery(uint64(m.RegoID))
	}
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	l = len(m.RegoHash)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	l = len(m.Source)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	if len(m.EntryPoints) > 0 {
		for _, s := range m.EntryPoints {
			l = len(s)
			n += 1 + l + sovQuery(uint64(l))
		}
	}
	return n
}

func (m *QueryRegosRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Pagination != nil {
		l = m.Pagination.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryRegosResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.RegoInfos) > 0 {
		for _, e := range m.RegoInfos {
			l = e.Size()
			n += 1 + l + sovQuery(uint64(l))
		}
	}
	if m.Pagination != nil {
		l = m.Pagination.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryRegoRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryRegoRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryRegoRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RegoId", wireType)
			}
			m.RegoId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RegoId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryRegoResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryRegoResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryRegoResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RegoInfoResponse", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.RegoInfoResponse == nil {
				m.RegoInfoResponse = &RegoInfoResponse{}
			}
			if err := m.RegoInfoResponse.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Data = append(m.Data[:0], dAtA[iNdEx:postIndex]...)
			if m.Data == nil {
				m.Data = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *RegoInfoResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RegoInfoResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RegoInfoResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RegoID", wireType)
			}
			m.RegoID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RegoID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Creator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RegoHash", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RegoHash = append(m.RegoHash[:0], dAtA[iNdEx:postIndex]...)
			if m.RegoHash == nil {
				m.RegoHash = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Source", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Source = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EntryPoints", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.EntryPoints = append(m.EntryPoints, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryRegosRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryRegosRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryRegosRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pagination", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pagination == nil {
				m.Pagination = &query.PageRequest{}
			}
			if err := m.Pagination.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryRegosResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryRegosResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryRegosResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RegoInfos", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RegoInfos = append(m.RegoInfos, RegoInfoResponse{})
			if err := m.RegoInfos[len(m.RegoInfos)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pagination", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pagination == nil {
				m.Pagination = &query.PageResponse{}
			}
			if err := m.Pagination.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)
