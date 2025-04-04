// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v6.30.1
// source: writy.proto

package libwrity

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
	WrityService_Set_FullMethodName   = "/WrityService/Set"
	WrityService_Get_FullMethodName   = "/WrityService/Get"
	WrityService_Del_FullMethodName   = "/WrityService/Del"
	WrityService_Keys_FullMethodName  = "/WrityService/Keys"
	WrityService_Flush_FullMethodName = "/WrityService/Flush"
)

// WrityServiceClient is the client API for WrityService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type WrityServiceClient interface {
	Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*Empty, error)
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	Del(ctx context.Context, in *DelRequest, opts ...grpc.CallOption) (*Empty, error)
	Keys(ctx context.Context, in *KeysRequest, opts ...grpc.CallOption) (*KeysResponse, error)
	Flush(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
}

type writyServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewWrityServiceClient(cc grpc.ClientConnInterface) WrityServiceClient {
	return &writyServiceClient{cc}
}

func (c *writyServiceClient) Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, WrityService_Set_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *writyServiceClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, WrityService_Get_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *writyServiceClient) Del(ctx context.Context, in *DelRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, WrityService_Del_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *writyServiceClient) Keys(ctx context.Context, in *KeysRequest, opts ...grpc.CallOption) (*KeysResponse, error) {
	out := new(KeysResponse)
	err := c.cc.Invoke(ctx, WrityService_Keys_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *writyServiceClient) Flush(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, WrityService_Flush_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WrityServiceServer is the server API for WrityService service.
// All implementations must embed UnimplementedWrityServiceServer
// for forward compatibility
type WrityServiceServer interface {
	Set(context.Context, *SetRequest) (*Empty, error)
	Get(context.Context, *GetRequest) (*GetResponse, error)
	Del(context.Context, *DelRequest) (*Empty, error)
	Keys(context.Context, *KeysRequest) (*KeysResponse, error)
	Flush(context.Context, *Empty) (*Empty, error)
	mustEmbedUnimplementedWrityServiceServer()
}

// UnimplementedWrityServiceServer must be embedded to have forward compatible implementations.
type UnimplementedWrityServiceServer struct {
}

func (UnimplementedWrityServiceServer) Set(context.Context, *SetRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
}
func (UnimplementedWrityServiceServer) Get(context.Context, *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedWrityServiceServer) Del(context.Context, *DelRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Del not implemented")
}
func (UnimplementedWrityServiceServer) Keys(context.Context, *KeysRequest) (*KeysResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Keys not implemented")
}
func (UnimplementedWrityServiceServer) Flush(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Flush not implemented")
}
func (UnimplementedWrityServiceServer) mustEmbedUnimplementedWrityServiceServer() {}

// UnsafeWrityServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WrityServiceServer will
// result in compilation errors.
type UnsafeWrityServiceServer interface {
	mustEmbedUnimplementedWrityServiceServer()
}

func RegisterWrityServiceServer(s grpc.ServiceRegistrar, srv WrityServiceServer) {
	s.RegisterService(&WrityService_ServiceDesc, srv)
}

func _WrityService_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WrityServiceServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WrityService_Set_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WrityServiceServer).Set(ctx, req.(*SetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WrityService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WrityServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WrityService_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WrityServiceServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WrityService_Del_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WrityServiceServer).Del(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WrityService_Del_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WrityServiceServer).Del(ctx, req.(*DelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WrityService_Keys_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeysRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WrityServiceServer).Keys(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WrityService_Keys_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WrityServiceServer).Keys(ctx, req.(*KeysRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WrityService_Flush_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WrityServiceServer).Flush(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WrityService_Flush_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WrityServiceServer).Flush(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// WrityService_ServiceDesc is the grpc.ServiceDesc for WrityService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var WrityService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "WrityService",
	HandlerType: (*WrityServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Set",
			Handler:    _WrityService_Set_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _WrityService_Get_Handler,
		},
		{
			MethodName: "Del",
			Handler:    _WrityService_Del_Handler,
		},
		{
			MethodName: "Keys",
			Handler:    _WrityService_Keys_Handler,
		},
		{
			MethodName: "Flush",
			Handler:    _WrityService_Flush_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "writy.proto",
}
