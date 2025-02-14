// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: milkyway/assets/v1/query.proto

package types

import (
	context "context"
	fmt "fmt"
	query "github.com/cosmos/cosmos-sdk/types/query"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
	_ "github.com/cosmos/gogoproto/gogoproto"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
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

// QueryAssetsRequest is the request type for the Query/Assets RPC method.
type QueryAssetsRequest struct {
	// Ticker defines an optional filter parameter to query assets with the given
	// ticker.
	Ticker string `protobuf:"bytes,1,opt,name=ticker,proto3" json:"ticker,omitempty"`
	// Pagination defines an optional pagination for the request.
	Pagination *query.PageRequest `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryAssetsRequest) Reset()         { *m = QueryAssetsRequest{} }
func (m *QueryAssetsRequest) String() string { return proto.CompactTextString(m) }
func (*QueryAssetsRequest) ProtoMessage()    {}
func (*QueryAssetsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2d64a9f39441c4d1, []int{0}
}
func (m *QueryAssetsRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryAssetsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryAssetsRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryAssetsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAssetsRequest.Merge(m, src)
}
func (m *QueryAssetsRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryAssetsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAssetsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAssetsRequest proto.InternalMessageInfo

func (m *QueryAssetsRequest) GetTicker() string {
	if m != nil {
		return m.Ticker
	}
	return ""
}

func (m *QueryAssetsRequest) GetPagination() *query.PageRequest {
	if m != nil {
		return m.Pagination
	}
	return nil
}

// QueryAssetsResponse is the response type for the Query/Assets RPC method.
type QueryAssetsResponse struct {
	// Assets represents all the assets registered.
	Assets []Asset `protobuf:"bytes,1,rep,name=assets,proto3" json:"assets"`
	// Pagination defines the pagination in the response.
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryAssetsResponse) Reset()         { *m = QueryAssetsResponse{} }
func (m *QueryAssetsResponse) String() string { return proto.CompactTextString(m) }
func (*QueryAssetsResponse) ProtoMessage()    {}
func (*QueryAssetsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_2d64a9f39441c4d1, []int{1}
}
func (m *QueryAssetsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryAssetsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryAssetsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryAssetsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAssetsResponse.Merge(m, src)
}
func (m *QueryAssetsResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryAssetsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAssetsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAssetsResponse proto.InternalMessageInfo

func (m *QueryAssetsResponse) GetAssets() []Asset {
	if m != nil {
		return m.Assets
	}
	return nil
}

func (m *QueryAssetsResponse) GetPagination() *query.PageResponse {
	if m != nil {
		return m.Pagination
	}
	return nil
}

// QueryAssetRequest is the request type for the Query/Asset RPC method.
type QueryAssetRequest struct {
	// Denom is the token denomination for which the ticker is to be queried.
	Denom string `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty"`
}

func (m *QueryAssetRequest) Reset()         { *m = QueryAssetRequest{} }
func (m *QueryAssetRequest) String() string { return proto.CompactTextString(m) }
func (*QueryAssetRequest) ProtoMessage()    {}
func (*QueryAssetRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2d64a9f39441c4d1, []int{2}
}
func (m *QueryAssetRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryAssetRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryAssetRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryAssetRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAssetRequest.Merge(m, src)
}
func (m *QueryAssetRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryAssetRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAssetRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAssetRequest proto.InternalMessageInfo

func (m *QueryAssetRequest) GetDenom() string {
	if m != nil {
		return m.Denom
	}
	return ""
}

// QueryAssetResponse is the response type for the Query/Asset RPC method.
type QueryAssetResponse struct {
	// Asset is the asset associated with the token denomination.
	Asset Asset `protobuf:"bytes,1,opt,name=asset,proto3" json:"asset"`
}

func (m *QueryAssetResponse) Reset()         { *m = QueryAssetResponse{} }
func (m *QueryAssetResponse) String() string { return proto.CompactTextString(m) }
func (*QueryAssetResponse) ProtoMessage()    {}
func (*QueryAssetResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_2d64a9f39441c4d1, []int{3}
}
func (m *QueryAssetResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryAssetResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryAssetResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryAssetResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAssetResponse.Merge(m, src)
}
func (m *QueryAssetResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryAssetResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAssetResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAssetResponse proto.InternalMessageInfo

func (m *QueryAssetResponse) GetAsset() Asset {
	if m != nil {
		return m.Asset
	}
	return Asset{}
}

func init() {
	proto.RegisterType((*QueryAssetsRequest)(nil), "milkyway.assets.v1.QueryAssetsRequest")
	proto.RegisterType((*QueryAssetsResponse)(nil), "milkyway.assets.v1.QueryAssetsResponse")
	proto.RegisterType((*QueryAssetRequest)(nil), "milkyway.assets.v1.QueryAssetRequest")
	proto.RegisterType((*QueryAssetResponse)(nil), "milkyway.assets.v1.QueryAssetResponse")
}

func init() { proto.RegisterFile("milkyway/assets/v1/query.proto", fileDescriptor_2d64a9f39441c4d1) }

var fileDescriptor_2d64a9f39441c4d1 = []byte{
	// 464 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x52, 0xbf, 0x6e, 0xd3, 0x40,
	0x1c, 0xf6, 0x05, 0x6c, 0xa9, 0xd7, 0xa9, 0x47, 0x85, 0x82, 0x55, 0xb9, 0x91, 0x05, 0x69, 0x88,
	0xc4, 0x9d, 0x9c, 0x8a, 0x01, 0x89, 0x85, 0x0e, 0x30, 0x20, 0x24, 0xf0, 0xc8, 0x76, 0x4e, 0x4f,
	0xc6, 0xaa, 0xed, 0x73, 0x73, 0x17, 0x83, 0x85, 0xba, 0x74, 0x65, 0x41, 0x62, 0xe5, 0x01, 0x18,
	0x79, 0x8c, 0x8e, 0x95, 0x58, 0x98, 0x10, 0x4a, 0x90, 0x78, 0x0d, 0xe4, 0xbb, 0x73, 0x93, 0x28,
	0x26, 0x59, 0xa2, 0xbb, 0xfc, 0xbe, 0xef, 0xf7, 0xfd, 0xf1, 0x41, 0x2f, 0x4b, 0xd2, 0xb3, 0xea,
	0x3d, 0xad, 0x08, 0x15, 0x82, 0x49, 0x41, 0xca, 0x80, 0x9c, 0x4f, 0xd9, 0xa4, 0xc2, 0xc5, 0x84,
	0x4b, 0x8e, 0x50, 0x33, 0xc7, 0x7a, 0x8e, 0xcb, 0xc0, 0xdd, 0xa3, 0x59, 0x92, 0x73, 0xa2, 0x7e,
	0x35, 0xcc, 0x1d, 0x8e, 0xb9, 0xc8, 0xb8, 0x20, 0x11, 0x15, 0x4c, 0xf3, 0x49, 0x19, 0x44, 0x4c,
	0xd2, 0x80, 0x14, 0x34, 0x4e, 0x72, 0x2a, 0x13, 0x9e, 0x1b, 0xec, 0x7e, 0xcc, 0x63, 0xae, 0x8e,
	0xa4, 0x3e, 0x99, 0x7f, 0x0f, 0x62, 0xce, 0xe3, 0x94, 0x11, 0x5a, 0x24, 0x84, 0xe6, 0x39, 0x97,
	0x8a, 0x22, 0xcc, 0xf4, 0xb0, 0xc5, 0x66, 0xc6, 0x4f, 0x59, 0x6a, 0x00, 0xbe, 0x84, 0xe8, 0x4d,
	0x2d, 0xfb, 0x4c, 0x8d, 0x43, 0x76, 0x3e, 0x65, 0x42, 0xa2, 0xbb, 0xd0, 0x91, 0xc9, 0xf8, 0x8c,
	0x4d, 0xba, 0xa0, 0x07, 0x06, 0x3b, 0xa1, 0xb9, 0xa1, 0xe7, 0x10, 0x2e, 0x6c, 0x75, 0x3b, 0x3d,
	0x30, 0xd8, 0x1d, 0xf5, 0xb1, 0xce, 0x80, 0xeb, 0x0c, 0x58, 0x77, 0x60, 0x32, 0xe0, 0xd7, 0x34,
	0x66, 0x66, 0x67, 0xb8, 0xc4, 0xf4, 0xbf, 0x02, 0x78, 0x67, 0x45, 0x56, 0x14, 0x3c, 0x17, 0x0c,
	0x3d, 0x85, 0x8e, 0xf6, 0xd9, 0x05, 0xbd, 0x5b, 0x83, 0xdd, 0xd1, 0x3d, 0xbc, 0x5e, 0x23, 0x56,
	0x9c, 0x93, 0x9d, 0xab, 0x5f, 0x87, 0xd6, 0xb7, 0xbf, 0xdf, 0x87, 0x20, 0x34, 0x1c, 0xf4, 0xa2,
	0xc5, 0xdd, 0xd1, 0x56, 0x77, 0x5a, 0x7a, 0xc5, 0xde, 0x43, 0xb8, 0xb7, 0x70, 0xd7, 0x74, 0xb2,
	0x0f, 0xed, 0x53, 0x96, 0xf3, 0xcc, 0x54, 0xa2, 0x2f, 0xfe, 0xcb, 0xe5, 0xfe, 0x6e, 0x72, 0x3c,
	0x86, 0xb6, 0xf2, 0xa4, 0xb0, 0x1b, 0x63, 0xdc, 0xae, 0x63, 0x84, 0x1a, 0x3d, 0xfa, 0xd4, 0x81,
	0xb6, 0xda, 0x86, 0x2e, 0xa0, 0xa3, 0xab, 0x41, 0xfd, 0x36, 0xee, 0xfa, 0x27, 0x73, 0x8f, 0xb6,
	0xe2, 0xb4, 0x37, 0xdf, 0xbf, 0xfc, 0xf1, 0xe7, 0x4b, 0xe7, 0x00, 0xb9, 0xa4, 0xe5, 0x6d, 0x98,
	0x26, 0x2f, 0x01, 0xb4, 0x15, 0x0d, 0x3d, 0xd8, 0xbc, 0xb6, 0x51, 0xef, 0x6f, 0x83, 0x19, 0xf1,
	0xa1, 0x12, 0xbf, 0x8f, 0xfc, 0xff, 0x8b, 0x93, 0x8f, 0xaa, 0xd9, 0x8b, 0x93, 0x57, 0x57, 0x33,
	0x0f, 0x5c, 0xcf, 0x3c, 0xf0, 0x7b, 0xe6, 0x81, 0xcf, 0x73, 0xcf, 0xba, 0x9e, 0x7b, 0xd6, 0xcf,
	0xb9, 0x67, 0xbd, 0x3d, 0x8e, 0x13, 0xf9, 0x6e, 0x1a, 0xe1, 0x31, 0xcf, 0x6e, 0xf6, 0x3c, 0x4a,
	0x69, 0x24, 0x16, 0x5b, 0xcb, 0x27, 0xe4, 0x43, 0xb3, 0x50, 0x56, 0x05, 0x13, 0x91, 0xa3, 0x1e,
	0xfc, 0xf1, 0xbf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x57, 0x4c, 0x98, 0x6b, 0xba, 0x03, 0x00, 0x00,
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
	// Assets defined a gRPC query method that returns all assets registered.
	Assets(ctx context.Context, in *QueryAssetsRequest, opts ...grpc.CallOption) (*QueryAssetsResponse, error)
	// Asset defines a gRPC query method that returns the asset associated with
	// the given token denomination.
	Asset(ctx context.Context, in *QueryAssetRequest, opts ...grpc.CallOption) (*QueryAssetResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Assets(ctx context.Context, in *QueryAssetsRequest, opts ...grpc.CallOption) (*QueryAssetsResponse, error) {
	out := new(QueryAssetsResponse)
	err := c.cc.Invoke(ctx, "/milkyway.assets.v1.Query/Assets", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Asset(ctx context.Context, in *QueryAssetRequest, opts ...grpc.CallOption) (*QueryAssetResponse, error) {
	out := new(QueryAssetResponse)
	err := c.cc.Invoke(ctx, "/milkyway.assets.v1.Query/Asset", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Assets defined a gRPC query method that returns all assets registered.
	Assets(context.Context, *QueryAssetsRequest) (*QueryAssetsResponse, error)
	// Asset defines a gRPC query method that returns the asset associated with
	// the given token denomination.
	Asset(context.Context, *QueryAssetRequest) (*QueryAssetResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) Assets(ctx context.Context, req *QueryAssetsRequest) (*QueryAssetsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Assets not implemented")
}
func (*UnimplementedQueryServer) Asset(ctx context.Context, req *QueryAssetRequest) (*QueryAssetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Asset not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_Assets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAssetsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Assets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/milkyway.assets.v1.Query/Assets",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Assets(ctx, req.(*QueryAssetsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Asset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAssetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Asset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/milkyway.assets.v1.Query/Asset",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Asset(ctx, req.(*QueryAssetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var Query_serviceDesc = _Query_serviceDesc
var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "milkyway.assets.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Assets",
			Handler:    _Query_Assets_Handler,
		},
		{
			MethodName: "Asset",
			Handler:    _Query_Asset_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "milkyway/assets/v1/query.proto",
}

func (m *QueryAssetsRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryAssetsRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryAssetsRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
	if len(m.Ticker) > 0 {
		i -= len(m.Ticker)
		copy(dAtA[i:], m.Ticker)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Ticker)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryAssetsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryAssetsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryAssetsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
	if len(m.Assets) > 0 {
		for iNdEx := len(m.Assets) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Assets[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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

func (m *QueryAssetRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryAssetRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryAssetRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryAssetResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryAssetResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryAssetResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Asset.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
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
func (m *QueryAssetsRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Ticker)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	if m.Pagination != nil {
		l = m.Pagination.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryAssetsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Assets) > 0 {
		for _, e := range m.Assets {
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

func (m *QueryAssetRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryAssetResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Asset.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryAssetsRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryAssetsRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryAssetsRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ticker", wireType)
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
			m.Ticker = string(dAtA[iNdEx:postIndex])
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
func (m *QueryAssetsResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryAssetsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryAssetsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Assets", wireType)
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
			m.Assets = append(m.Assets, Asset{})
			if err := m.Assets[len(m.Assets)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *QueryAssetRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryAssetRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryAssetRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
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
			m.Denom = string(dAtA[iNdEx:postIndex])
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
func (m *QueryAssetResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryAssetResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryAssetResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Asset", wireType)
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
			if err := m.Asset.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
