// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: milkyway/services/v1/messages.proto

package servicesv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Msg_CreateService_FullMethodName     = "/milkyway.services.v1.Msg/CreateService"
	Msg_UpdateService_FullMethodName     = "/milkyway.services.v1.Msg/UpdateService"
	Msg_ActivateService_FullMethodName   = "/milkyway.services.v1.Msg/ActivateService"
	Msg_DeactivateService_FullMethodName = "/milkyway.services.v1.Msg/DeactivateService"
	Msg_UpdateParams_FullMethodName      = "/milkyway.services.v1.Msg/UpdateParams"
)

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Msg defines the services module's gRPC message service.
type MsgClient interface {
	// CreateService defines the operation for registering a new service.
	CreateService(ctx context.Context, in *MsgCreateService, opts ...grpc.CallOption) (*MsgCreateServiceResponse, error)
	// UpdateService defines the operation for updating an existing service.
	UpdateService(ctx context.Context, in *MsgUpdateService, opts ...grpc.CallOption) (*MsgUpdateServiceResponse, error)
	// ActivateService defines the operation for activating an existing
	// service.
	ActivateService(ctx context.Context, in *MsgActivateService, opts ...grpc.CallOption) (*MsgActivateServiceResponse, error)
	// DeactivateService defines the operation for deactivating an existing
	// service.
	DeactivateService(ctx context.Context, in *MsgDeactivateService, opts ...grpc.CallOption) (*MsgDeactivateServiceResponse, error)
	// UpdateParams defines a (governance) operation for updating the module
	// parameters.
	// The authority defaults to the x/gov module account.
	UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error)
}

type msgClient struct {
	cc grpc.ClientConnInterface
}

func NewMsgClient(cc grpc.ClientConnInterface) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) CreateService(ctx context.Context, in *MsgCreateService, opts ...grpc.CallOption) (*MsgCreateServiceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MsgCreateServiceResponse)
	err := c.cc.Invoke(ctx, Msg_CreateService_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateService(ctx context.Context, in *MsgUpdateService, opts ...grpc.CallOption) (*MsgUpdateServiceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MsgUpdateServiceResponse)
	err := c.cc.Invoke(ctx, Msg_UpdateService_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) ActivateService(ctx context.Context, in *MsgActivateService, opts ...grpc.CallOption) (*MsgActivateServiceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MsgActivateServiceResponse)
	err := c.cc.Invoke(ctx, Msg_ActivateService_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) DeactivateService(ctx context.Context, in *MsgDeactivateService, opts ...grpc.CallOption) (*MsgDeactivateServiceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MsgDeactivateServiceResponse)
	err := c.cc.Invoke(ctx, Msg_DeactivateService_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MsgUpdateParamsResponse)
	err := c.cc.Invoke(ctx, Msg_UpdateParams_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
// All implementations must embed UnimplementedMsgServer
// for forward compatibility.
//
// Msg defines the services module's gRPC message service.
type MsgServer interface {
	// CreateService defines the operation for registering a new service.
	CreateService(context.Context, *MsgCreateService) (*MsgCreateServiceResponse, error)
	// UpdateService defines the operation for updating an existing service.
	UpdateService(context.Context, *MsgUpdateService) (*MsgUpdateServiceResponse, error)
	// ActivateService defines the operation for activating an existing
	// service.
	ActivateService(context.Context, *MsgActivateService) (*MsgActivateServiceResponse, error)
	// DeactivateService defines the operation for deactivating an existing
	// service.
	DeactivateService(context.Context, *MsgDeactivateService) (*MsgDeactivateServiceResponse, error)
	// UpdateParams defines a (governance) operation for updating the module
	// parameters.
	// The authority defaults to the x/gov module account.
	UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error)
	mustEmbedUnimplementedMsgServer()
}

// UnimplementedMsgServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMsgServer struct{}

func (UnimplementedMsgServer) CreateService(context.Context, *MsgCreateService) (*MsgCreateServiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateService not implemented")
}
func (UnimplementedMsgServer) UpdateService(context.Context, *MsgUpdateService) (*MsgUpdateServiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateService not implemented")
}
func (UnimplementedMsgServer) ActivateService(context.Context, *MsgActivateService) (*MsgActivateServiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ActivateService not implemented")
}
func (UnimplementedMsgServer) DeactivateService(context.Context, *MsgDeactivateService) (*MsgDeactivateServiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeactivateService not implemented")
}
func (UnimplementedMsgServer) UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateParams not implemented")
}
func (UnimplementedMsgServer) mustEmbedUnimplementedMsgServer() {}
func (UnimplementedMsgServer) testEmbeddedByValue()             {}

// UnsafeMsgServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MsgServer will
// result in compilation errors.
type UnsafeMsgServer interface {
	mustEmbedUnimplementedMsgServer()
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	// If the following call pancis, it indicates UnimplementedMsgServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Msg_ServiceDesc, srv)
}

func _Msg_CreateService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateService)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreateService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_CreateService_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreateService(ctx, req.(*MsgCreateService))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateService)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_UpdateService_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateService(ctx, req.(*MsgUpdateService))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_ActivateService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgActivateService)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ActivateService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_ActivateService_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ActivateService(ctx, req.(*MsgActivateService))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_DeactivateService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgDeactivateService)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).DeactivateService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_DeactivateService_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).DeactivateService(ctx, req.(*MsgDeactivateService))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_UpdateParams_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateParams(ctx, req.(*MsgUpdateParams))
	}
	return interceptor(ctx, in, info, handler)
}

// Msg_ServiceDesc is the grpc.ServiceDesc for Msg service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "milkyway.services.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateService",
			Handler:    _Msg_CreateService_Handler,
		},
		{
			MethodName: "UpdateService",
			Handler:    _Msg_UpdateService_Handler,
		},
		{
			MethodName: "ActivateService",
			Handler:    _Msg_ActivateService_Handler,
		},
		{
			MethodName: "DeactivateService",
			Handler:    _Msg_DeactivateService_Handler,
		},
		{
			MethodName: "UpdateParams",
			Handler:    _Msg_UpdateParams_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "milkyway/services/v1/messages.proto",
}