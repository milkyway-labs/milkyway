// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: stride/stakeibc/query.proto

package stakeibc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Query_Params_FullMethodName                       = "/stride.stakeibc.Query/Params"
	Query_Validators_FullMethodName                   = "/stride.stakeibc.Query/Validators"
	Query_HostZone_FullMethodName                     = "/stride.stakeibc.Query/HostZone"
	Query_HostZoneAll_FullMethodName                  = "/stride.stakeibc.Query/HostZoneAll"
	Query_ModuleAddress_FullMethodName                = "/stride.stakeibc.Query/ModuleAddress"
	Query_InterchainAccountFromAddress_FullMethodName = "/stride.stakeibc.Query/InterchainAccountFromAddress"
	Query_EpochTracker_FullMethodName                 = "/stride.stakeibc.Query/EpochTracker"
	Query_EpochTrackerAll_FullMethodName              = "/stride.stakeibc.Query/EpochTrackerAll"
	Query_NextPacketSequence_FullMethodName           = "/stride.stakeibc.Query/NextPacketSequence"
	Query_AddressUnbondings_FullMethodName            = "/stride.stakeibc.Query/AddressUnbondings"
	Query_AllTradeRoutes_FullMethodName               = "/stride.stakeibc.Query/AllTradeRoutes"
)

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QueryClient interface {
	// Parameters queries the parameters of the module.
	Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error)
	// Queries a Validator by host zone.
	Validators(ctx context.Context, in *QueryGetValidatorsRequest, opts ...grpc.CallOption) (*QueryGetValidatorsResponse, error)
	// Queries a HostZone by id.
	HostZone(ctx context.Context, in *QueryGetHostZoneRequest, opts ...grpc.CallOption) (*QueryGetHostZoneResponse, error)
	// Queries a list of HostZone items.
	HostZoneAll(ctx context.Context, in *QueryAllHostZoneRequest, opts ...grpc.CallOption) (*QueryAllHostZoneResponse, error)
	// Queries a list of ModuleAddress items.
	ModuleAddress(ctx context.Context, in *QueryModuleAddressRequest, opts ...grpc.CallOption) (*QueryModuleAddressResponse, error)
	// QueryInterchainAccountFromAddress returns the interchain account for given
	// owner address on a given connection pair
	InterchainAccountFromAddress(ctx context.Context, in *QueryInterchainAccountFromAddressRequest, opts ...grpc.CallOption) (*QueryInterchainAccountFromAddressResponse, error)
	// Queries a EpochTracker by index.
	EpochTracker(ctx context.Context, in *QueryGetEpochTrackerRequest, opts ...grpc.CallOption) (*QueryGetEpochTrackerResponse, error)
	// Queries a list of EpochTracker items.
	EpochTrackerAll(ctx context.Context, in *QueryAllEpochTrackerRequest, opts ...grpc.CallOption) (*QueryAllEpochTrackerResponse, error)
	// Queries the next packet sequence for one for a given channel
	NextPacketSequence(ctx context.Context, in *QueryGetNextPacketSequenceRequest, opts ...grpc.CallOption) (*QueryGetNextPacketSequenceResponse, error)
	// Queries an address's unbondings
	AddressUnbondings(ctx context.Context, in *QueryAddressUnbondings, opts ...grpc.CallOption) (*QueryAddressUnbondingsResponse, error)
	// Queries all trade routes
	AllTradeRoutes(ctx context.Context, in *QueryAllTradeRoutes, opts ...grpc.CallOption) (*QueryAllTradeRoutesResponse, error)
}

type queryClient struct {
	cc grpc.ClientConnInterface
}

func NewQueryClient(cc grpc.ClientConnInterface) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error) {
	out := new(QueryParamsResponse)
	err := c.cc.Invoke(ctx, Query_Params_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Validators(ctx context.Context, in *QueryGetValidatorsRequest, opts ...grpc.CallOption) (*QueryGetValidatorsResponse, error) {
	out := new(QueryGetValidatorsResponse)
	err := c.cc.Invoke(ctx, Query_Validators_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) HostZone(ctx context.Context, in *QueryGetHostZoneRequest, opts ...grpc.CallOption) (*QueryGetHostZoneResponse, error) {
	out := new(QueryGetHostZoneResponse)
	err := c.cc.Invoke(ctx, Query_HostZone_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) HostZoneAll(ctx context.Context, in *QueryAllHostZoneRequest, opts ...grpc.CallOption) (*QueryAllHostZoneResponse, error) {
	out := new(QueryAllHostZoneResponse)
	err := c.cc.Invoke(ctx, Query_HostZoneAll_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) ModuleAddress(ctx context.Context, in *QueryModuleAddressRequest, opts ...grpc.CallOption) (*QueryModuleAddressResponse, error) {
	out := new(QueryModuleAddressResponse)
	err := c.cc.Invoke(ctx, Query_ModuleAddress_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) InterchainAccountFromAddress(ctx context.Context, in *QueryInterchainAccountFromAddressRequest, opts ...grpc.CallOption) (*QueryInterchainAccountFromAddressResponse, error) {
	out := new(QueryInterchainAccountFromAddressResponse)
	err := c.cc.Invoke(ctx, Query_InterchainAccountFromAddress_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) EpochTracker(ctx context.Context, in *QueryGetEpochTrackerRequest, opts ...grpc.CallOption) (*QueryGetEpochTrackerResponse, error) {
	out := new(QueryGetEpochTrackerResponse)
	err := c.cc.Invoke(ctx, Query_EpochTracker_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) EpochTrackerAll(ctx context.Context, in *QueryAllEpochTrackerRequest, opts ...grpc.CallOption) (*QueryAllEpochTrackerResponse, error) {
	out := new(QueryAllEpochTrackerResponse)
	err := c.cc.Invoke(ctx, Query_EpochTrackerAll_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) NextPacketSequence(ctx context.Context, in *QueryGetNextPacketSequenceRequest, opts ...grpc.CallOption) (*QueryGetNextPacketSequenceResponse, error) {
	out := new(QueryGetNextPacketSequenceResponse)
	err := c.cc.Invoke(ctx, Query_NextPacketSequence_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) AddressUnbondings(ctx context.Context, in *QueryAddressUnbondings, opts ...grpc.CallOption) (*QueryAddressUnbondingsResponse, error) {
	out := new(QueryAddressUnbondingsResponse)
	err := c.cc.Invoke(ctx, Query_AddressUnbondings_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) AllTradeRoutes(ctx context.Context, in *QueryAllTradeRoutes, opts ...grpc.CallOption) (*QueryAllTradeRoutesResponse, error) {
	out := new(QueryAllTradeRoutesResponse)
	err := c.cc.Invoke(ctx, Query_AllTradeRoutes_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
// All implementations must embed UnimplementedQueryServer
// for forward compatibility
type QueryServer interface {
	// Parameters queries the parameters of the module.
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	// Queries a Validator by host zone.
	Validators(context.Context, *QueryGetValidatorsRequest) (*QueryGetValidatorsResponse, error)
	// Queries a HostZone by id.
	HostZone(context.Context, *QueryGetHostZoneRequest) (*QueryGetHostZoneResponse, error)
	// Queries a list of HostZone items.
	HostZoneAll(context.Context, *QueryAllHostZoneRequest) (*QueryAllHostZoneResponse, error)
	// Queries a list of ModuleAddress items.
	ModuleAddress(context.Context, *QueryModuleAddressRequest) (*QueryModuleAddressResponse, error)
	// QueryInterchainAccountFromAddress returns the interchain account for given
	// owner address on a given connection pair
	InterchainAccountFromAddress(context.Context, *QueryInterchainAccountFromAddressRequest) (*QueryInterchainAccountFromAddressResponse, error)
	// Queries a EpochTracker by index.
	EpochTracker(context.Context, *QueryGetEpochTrackerRequest) (*QueryGetEpochTrackerResponse, error)
	// Queries a list of EpochTracker items.
	EpochTrackerAll(context.Context, *QueryAllEpochTrackerRequest) (*QueryAllEpochTrackerResponse, error)
	// Queries the next packet sequence for one for a given channel
	NextPacketSequence(context.Context, *QueryGetNextPacketSequenceRequest) (*QueryGetNextPacketSequenceResponse, error)
	// Queries an address's unbondings
	AddressUnbondings(context.Context, *QueryAddressUnbondings) (*QueryAddressUnbondingsResponse, error)
	// Queries all trade routes
	AllTradeRoutes(context.Context, *QueryAllTradeRoutes) (*QueryAllTradeRoutesResponse, error)
	mustEmbedUnimplementedQueryServer()
}

// UnimplementedQueryServer must be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (UnimplementedQueryServer) Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}
func (UnimplementedQueryServer) Validators(context.Context, *QueryGetValidatorsRequest) (*QueryGetValidatorsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Validators not implemented")
}
func (UnimplementedQueryServer) HostZone(context.Context, *QueryGetHostZoneRequest) (*QueryGetHostZoneResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HostZone not implemented")
}
func (UnimplementedQueryServer) HostZoneAll(context.Context, *QueryAllHostZoneRequest) (*QueryAllHostZoneResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HostZoneAll not implemented")
}
func (UnimplementedQueryServer) ModuleAddress(context.Context, *QueryModuleAddressRequest) (*QueryModuleAddressResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ModuleAddress not implemented")
}
func (UnimplementedQueryServer) InterchainAccountFromAddress(context.Context, *QueryInterchainAccountFromAddressRequest) (*QueryInterchainAccountFromAddressResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InterchainAccountFromAddress not implemented")
}
func (UnimplementedQueryServer) EpochTracker(context.Context, *QueryGetEpochTrackerRequest) (*QueryGetEpochTrackerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EpochTracker not implemented")
}
func (UnimplementedQueryServer) EpochTrackerAll(context.Context, *QueryAllEpochTrackerRequest) (*QueryAllEpochTrackerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EpochTrackerAll not implemented")
}
func (UnimplementedQueryServer) NextPacketSequence(context.Context, *QueryGetNextPacketSequenceRequest) (*QueryGetNextPacketSequenceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NextPacketSequence not implemented")
}
func (UnimplementedQueryServer) AddressUnbondings(context.Context, *QueryAddressUnbondings) (*QueryAddressUnbondingsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddressUnbondings not implemented")
}
func (UnimplementedQueryServer) AllTradeRoutes(context.Context, *QueryAllTradeRoutes) (*QueryAllTradeRoutesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AllTradeRoutes not implemented")
}
func (UnimplementedQueryServer) mustEmbedUnimplementedQueryServer() {}

// UnsafeQueryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QueryServer will
// result in compilation errors.
type UnsafeQueryServer interface {
	mustEmbedUnimplementedQueryServer()
}

func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	s.RegisterService(&Query_ServiceDesc, srv)
}

func _Query_Params_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Params(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Params_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Validators_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetValidatorsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Validators(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Validators_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Validators(ctx, req.(*QueryGetValidatorsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_HostZone_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetHostZoneRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).HostZone(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_HostZone_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).HostZone(ctx, req.(*QueryGetHostZoneRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_HostZoneAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAllHostZoneRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).HostZoneAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_HostZoneAll_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).HostZoneAll(ctx, req.(*QueryAllHostZoneRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_ModuleAddress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryModuleAddressRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ModuleAddress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_ModuleAddress_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ModuleAddress(ctx, req.(*QueryModuleAddressRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_InterchainAccountFromAddress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryInterchainAccountFromAddressRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).InterchainAccountFromAddress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_InterchainAccountFromAddress_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).InterchainAccountFromAddress(ctx, req.(*QueryInterchainAccountFromAddressRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_EpochTracker_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetEpochTrackerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).EpochTracker(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_EpochTracker_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).EpochTracker(ctx, req.(*QueryGetEpochTrackerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_EpochTrackerAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAllEpochTrackerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).EpochTrackerAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_EpochTrackerAll_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).EpochTrackerAll(ctx, req.(*QueryAllEpochTrackerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_NextPacketSequence_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetNextPacketSequenceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).NextPacketSequence(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_NextPacketSequence_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).NextPacketSequence(ctx, req.(*QueryGetNextPacketSequenceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_AddressUnbondings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAddressUnbondings)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).AddressUnbondings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_AddressUnbondings_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).AddressUnbondings(ctx, req.(*QueryAddressUnbondings))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_AllTradeRoutes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAllTradeRoutes)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).AllTradeRoutes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_AllTradeRoutes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).AllTradeRoutes(ctx, req.(*QueryAllTradeRoutes))
	}
	return interceptor(ctx, in, info, handler)
}

// Query_ServiceDesc is the grpc.ServiceDesc for Query service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Query_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "stride.stakeibc.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "Validators",
			Handler:    _Query_Validators_Handler,
		},
		{
			MethodName: "HostZone",
			Handler:    _Query_HostZone_Handler,
		},
		{
			MethodName: "HostZoneAll",
			Handler:    _Query_HostZoneAll_Handler,
		},
		{
			MethodName: "ModuleAddress",
			Handler:    _Query_ModuleAddress_Handler,
		},
		{
			MethodName: "InterchainAccountFromAddress",
			Handler:    _Query_InterchainAccountFromAddress_Handler,
		},
		{
			MethodName: "EpochTracker",
			Handler:    _Query_EpochTracker_Handler,
		},
		{
			MethodName: "EpochTrackerAll",
			Handler:    _Query_EpochTrackerAll_Handler,
		},
		{
			MethodName: "NextPacketSequence",
			Handler:    _Query_NextPacketSequence_Handler,
		},
		{
			MethodName: "AddressUnbondings",
			Handler:    _Query_AddressUnbondings_Handler,
		},
		{
			MethodName: "AllTradeRoutes",
			Handler:    _Query_AllTradeRoutes_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "stride/stakeibc/query.proto",
}