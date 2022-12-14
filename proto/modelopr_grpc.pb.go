// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: modelopr.proto

package modelpb

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

// DatalakeSvcClient is the client API for DatalakeSvc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DatalakeSvcClient interface {
	GetDataFromUUID(ctx context.Context, in *UUIDExchangeRequest, opts ...grpc.CallOption) (*UUIDExchangeResponse, error)
}

type datalakeSvcClient struct {
	cc grpc.ClientConnInterface
}

func NewDatalakeSvcClient(cc grpc.ClientConnInterface) DatalakeSvcClient {
	return &datalakeSvcClient{cc}
}

func (c *datalakeSvcClient) GetDataFromUUID(ctx context.Context, in *UUIDExchangeRequest, opts ...grpc.CallOption) (*UUIDExchangeResponse, error) {
	out := new(UUIDExchangeResponse)
	err := c.cc.Invoke(ctx, "/da.DatalakeSvc/getDataFromUUID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DatalakeSvcServer is the server API for DatalakeSvc service.
// All implementations must embed UnimplementedDatalakeSvcServer
// for forward compatibility
type DatalakeSvcServer interface {
	GetDataFromUUID(context.Context, *UUIDExchangeRequest) (*UUIDExchangeResponse, error)
	mustEmbedUnimplementedDatalakeSvcServer()
}

// UnimplementedDatalakeSvcServer must be embedded to have forward compatible implementations.
type UnimplementedDatalakeSvcServer struct {
}

func (UnimplementedDatalakeSvcServer) GetDataFromUUID(context.Context, *UUIDExchangeRequest) (*UUIDExchangeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDataFromUUID not implemented")
}
func (UnimplementedDatalakeSvcServer) mustEmbedUnimplementedDatalakeSvcServer() {}

// UnsafeDatalakeSvcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DatalakeSvcServer will
// result in compilation errors.
type UnsafeDatalakeSvcServer interface {
	mustEmbedUnimplementedDatalakeSvcServer()
}

func RegisterDatalakeSvcServer(s grpc.ServiceRegistrar, srv DatalakeSvcServer) {
	s.RegisterService(&DatalakeSvc_ServiceDesc, srv)
}

func _DatalakeSvc_GetDataFromUUID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UUIDExchangeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatalakeSvcServer).GetDataFromUUID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/da.DatalakeSvc/getDataFromUUID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatalakeSvcServer).GetDataFromUUID(ctx, req.(*UUIDExchangeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DatalakeSvc_ServiceDesc is the grpc.ServiceDesc for DatalakeSvc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DatalakeSvc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "da.DatalakeSvc",
	HandlerType: (*DatalakeSvcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "getDataFromUUID",
			Handler:    _DatalakeSvc_GetDataFromUUID_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "modelopr.proto",
}
