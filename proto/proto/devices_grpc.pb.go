// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.11
// source: proto/devices.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// DevicesClient is the client API for Devices service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DevicesClient interface {
	AddDevice(ctx context.Context, in *AddDeviceReq, opts ...grpc.CallOption) (*Device, error)
	ListSpecificDeviceForUser(ctx context.Context, in *ListDevicesReq, opts ...grpc.CallOption) (*ListDevicesRes, error)
	DeleteDevice(ctx context.Context, in *DeleteDeviceReq, opts ...grpc.CallOption) (*emptypb.Empty, error)
	ListAllDevices(ctx context.Context, in *ListAllDevicesReq, opts ...grpc.CallOption) (*ListAllDevicesRes, error)
}

type devicesClient struct {
	cc grpc.ClientConnInterface
}

func NewDevicesClient(cc grpc.ClientConnInterface) DevicesClient {
	return &devicesClient{cc}
}

func (c *devicesClient) AddDevice(ctx context.Context, in *AddDeviceReq, opts ...grpc.CallOption) (*Device, error) {
	out := new(Device)
	err := c.cc.Invoke(ctx, "/proto.Devices/AddDevice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *devicesClient) ListSpecificDeviceForUser(ctx context.Context, in *ListDevicesReq, opts ...grpc.CallOption) (*ListDevicesRes, error) {
	out := new(ListDevicesRes)
	err := c.cc.Invoke(ctx, "/proto.Devices/ListSpecificDeviceForUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *devicesClient) DeleteDevice(ctx context.Context, in *DeleteDeviceReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/proto.Devices/DeleteDevice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *devicesClient) ListAllDevices(ctx context.Context, in *ListAllDevicesReq, opts ...grpc.CallOption) (*ListAllDevicesRes, error) {
	out := new(ListAllDevicesRes)
	err := c.cc.Invoke(ctx, "/proto.Devices/ListAllDevices", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DevicesServer is the server API for Devices service.
// All implementations must embed UnimplementedDevicesServer
// for forward compatibility
type DevicesServer interface {
	AddDevice(context.Context, *AddDeviceReq) (*Device, error)
	ListSpecificDeviceForUser(context.Context, *ListDevicesReq) (*ListDevicesRes, error)
	DeleteDevice(context.Context, *DeleteDeviceReq) (*emptypb.Empty, error)
	ListAllDevices(context.Context, *ListAllDevicesReq) (*ListAllDevicesRes, error)
	mustEmbedUnimplementedDevicesServer()
}

// UnimplementedDevicesServer must be embedded to have forward compatible implementations.
type UnimplementedDevicesServer struct {
}

func (UnimplementedDevicesServer) AddDevice(context.Context, *AddDeviceReq) (*Device, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddDevice not implemented")
}
func (UnimplementedDevicesServer) ListSpecificDeviceForUser(context.Context, *ListDevicesReq) (*ListDevicesRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSpecificDeviceForUser not implemented")
}
func (UnimplementedDevicesServer) DeleteDevice(context.Context, *DeleteDeviceReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteDevice not implemented")
}
func (UnimplementedDevicesServer) ListAllDevices(context.Context, *ListAllDevicesReq) (*ListAllDevicesRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAllDevices not implemented")
}
func (UnimplementedDevicesServer) mustEmbedUnimplementedDevicesServer() {}

// UnsafeDevicesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DevicesServer will
// result in compilation errors.
type UnsafeDevicesServer interface {
	mustEmbedUnimplementedDevicesServer()
}

func RegisterDevicesServer(s grpc.ServiceRegistrar, srv DevicesServer) {
	s.RegisterService(&Devices_ServiceDesc, srv)
}

func _Devices_AddDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddDeviceReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DevicesServer).AddDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Devices/AddDevice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DevicesServer).AddDevice(ctx, req.(*AddDeviceReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Devices_ListSpecificDeviceForUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListDevicesReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DevicesServer).ListSpecificDeviceForUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Devices/ListSpecificDeviceForUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DevicesServer).ListSpecificDeviceForUser(ctx, req.(*ListDevicesReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Devices_DeleteDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteDeviceReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DevicesServer).DeleteDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Devices/DeleteDevice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DevicesServer).DeleteDevice(ctx, req.(*DeleteDeviceReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Devices_ListAllDevices_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListAllDevicesReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DevicesServer).ListAllDevices(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Devices/ListAllDevices",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DevicesServer).ListAllDevices(ctx, req.(*ListAllDevicesReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Devices_ServiceDesc is the grpc.ServiceDesc for Devices service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Devices_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Devices",
	HandlerType: (*DevicesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddDevice",
			Handler:    _Devices_AddDevice_Handler,
		},
		{
			MethodName: "ListSpecificDeviceForUser",
			Handler:    _Devices_ListSpecificDeviceForUser_Handler,
		},
		{
			MethodName: "DeleteDevice",
			Handler:    _Devices_DeleteDevice_Handler,
		},
		{
			MethodName: "ListAllDevices",
			Handler:    _Devices_ListAllDevices_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/devices.proto",
}
