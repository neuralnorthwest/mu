// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: complete/v1/hello.proto

package v1

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

// CompleteServiceClient is the client API for CompleteService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CompleteServiceClient interface {
	// Hello is a simple RPC that returns a greeting.
	Hello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloResponse, error)
}

type completeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCompleteServiceClient(cc grpc.ClientConnInterface) CompleteServiceClient {
	return &completeServiceClient{cc}
}

func (c *completeServiceClient) Hello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloResponse, error) {
	out := new(HelloResponse)
	err := c.cc.Invoke(ctx, "/complete.v1.CompleteService/Hello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CompleteServiceServer is the server API for CompleteService service.
// All implementations must embed UnimplementedCompleteServiceServer
// for forward compatibility
type CompleteServiceServer interface {
	// Hello is a simple RPC that returns a greeting.
	Hello(context.Context, *HelloRequest) (*HelloResponse, error)
	mustEmbedUnimplementedCompleteServiceServer()
}

// UnimplementedCompleteServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCompleteServiceServer struct {
}

func (UnimplementedCompleteServiceServer) Hello(context.Context, *HelloRequest) (*HelloResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Hello not implemented")
}
func (UnimplementedCompleteServiceServer) mustEmbedUnimplementedCompleteServiceServer() {}

// UnsafeCompleteServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CompleteServiceServer will
// result in compilation errors.
type UnsafeCompleteServiceServer interface {
	mustEmbedUnimplementedCompleteServiceServer()
}

func RegisterCompleteServiceServer(s grpc.ServiceRegistrar, srv CompleteServiceServer) {
	s.RegisterService(&CompleteService_ServiceDesc, srv)
}

func _CompleteService_Hello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompleteServiceServer).Hello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/complete.v1.CompleteService/Hello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompleteServiceServer).Hello(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CompleteService_ServiceDesc is the grpc.ServiceDesc for CompleteService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CompleteService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "complete.v1.CompleteService",
	HandlerType: (*CompleteServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Hello",
			Handler:    _CompleteService_Hello_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "complete/v1/hello.proto",
}