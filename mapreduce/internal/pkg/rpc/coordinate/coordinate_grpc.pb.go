// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: coordinate.proto

package coordinate

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

// CoordinateClient is the client API for Coordinate service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CoordinateClient interface {
	Register(ctx context.Context, in *ClientInfo, opts ...grpc.CallOption) (*Response, error)
	Heartbeat(ctx context.Context, in *HeartbeatRequest, opts ...grpc.CallOption) (*Response, error)
}

type coordinateClient struct {
	cc grpc.ClientConnInterface
}

func NewCoordinateClient(cc grpc.ClientConnInterface) CoordinateClient {
	return &coordinateClient{cc}
}

func (c *coordinateClient) Register(ctx context.Context, in *ClientInfo, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/coordinate.Coordinate/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *coordinateClient) Heartbeat(ctx context.Context, in *HeartbeatRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/coordinate.Coordinate/Heartbeat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CoordinateServer is the server API for Coordinate service.
// All implementations must embed UnimplementedCoordinateServer
// for forward compatibility
type CoordinateServer interface {
	Register(context.Context, *ClientInfo) (*Response, error)
	Heartbeat(context.Context, *HeartbeatRequest) (*Response, error)
	mustEmbedUnimplementedCoordinateServer()
}

// UnimplementedCoordinateServer must be embedded to have forward compatible implementations.
type UnimplementedCoordinateServer struct {
}

func (UnimplementedCoordinateServer) Register(context.Context, *ClientInfo) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedCoordinateServer) Heartbeat(context.Context, *HeartbeatRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Heartbeat not implemented")
}
func (UnimplementedCoordinateServer) mustEmbedUnimplementedCoordinateServer() {}

// UnsafeCoordinateServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CoordinateServer will
// result in compilation errors.
type UnsafeCoordinateServer interface {
	mustEmbedUnimplementedCoordinateServer()
}

func RegisterCoordinateServer(s grpc.ServiceRegistrar, srv CoordinateServer) {
	s.RegisterService(&Coordinate_ServiceDesc, srv)
}

func _Coordinate_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClientInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CoordinateServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/coordinate.Coordinate/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CoordinateServer).Register(ctx, req.(*ClientInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Coordinate_Heartbeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HeartbeatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CoordinateServer).Heartbeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/coordinate.Coordinate/Heartbeat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CoordinateServer).Heartbeat(ctx, req.(*HeartbeatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Coordinate_ServiceDesc is the grpc.ServiceDesc for Coordinate service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Coordinate_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "coordinate.Coordinate",
	HandlerType: (*CoordinateServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _Coordinate_Register_Handler,
		},
		{
			MethodName: "Heartbeat",
			Handler:    _Coordinate_Heartbeat_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "coordinate.proto",
}
