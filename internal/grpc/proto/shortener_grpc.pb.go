// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: proto/shortener.proto

package proto

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
	Shortener_CreateShortening_FullMethodName      = "/shortener.Shortener/CreateShortening"
	Shortener_GetFullString_FullMethodName         = "/shortener.Shortener/GetFullString"
	Shortener_CreateShorteningBatch_FullMethodName = "/shortener.Shortener/CreateShorteningBatch"
	Shortener_GetUserAllShortenings_FullMethodName = "/shortener.Shortener/GetUserAllShortenings"
	Shortener_DeleteRecord_FullMethodName          = "/shortener.Shortener/DeleteRecord"
	Shortener_GetStats_FullMethodName              = "/shortener.Shortener/GetStats"
)

// ShortenerClient is the client API for Shortener service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ShortenerClient interface {
	CreateShortening(ctx context.Context, in *CreateShorteningRequest, opts ...grpc.CallOption) (*ShorteningResponse, error)
	GetFullString(ctx context.Context, in *LongURLRequest, opts ...grpc.CallOption) (*LongURLResponse, error)
	CreateShorteningBatch(ctx context.Context, in *CreateShorteningBatchRequest, opts ...grpc.CallOption) (*ShorteningBatchResponse, error)
	GetUserAllShortenings(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*ShorteningBatchResponse, error)
	DeleteRecord(ctx context.Context, in *DeleteRecordRequest, opts ...grpc.CallOption) (*None, error)
	GetStats(ctx context.Context, in *None, opts ...grpc.CallOption) (*GetStatsResponse, error)
}

type shortenerClient struct {
	cc grpc.ClientConnInterface
}

func NewShortenerClient(cc grpc.ClientConnInterface) ShortenerClient {
	return &shortenerClient{cc}
}

func (c *shortenerClient) CreateShortening(ctx context.Context, in *CreateShorteningRequest, opts ...grpc.CallOption) (*ShorteningResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ShorteningResponse)
	err := c.cc.Invoke(ctx, Shortener_CreateShortening_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) GetFullString(ctx context.Context, in *LongURLRequest, opts ...grpc.CallOption) (*LongURLResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LongURLResponse)
	err := c.cc.Invoke(ctx, Shortener_GetFullString_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) CreateShorteningBatch(ctx context.Context, in *CreateShorteningBatchRequest, opts ...grpc.CallOption) (*ShorteningBatchResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ShorteningBatchResponse)
	err := c.cc.Invoke(ctx, Shortener_CreateShorteningBatch_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) GetUserAllShortenings(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*ShorteningBatchResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ShorteningBatchResponse)
	err := c.cc.Invoke(ctx, Shortener_GetUserAllShortenings_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) DeleteRecord(ctx context.Context, in *DeleteRecordRequest, opts ...grpc.CallOption) (*None, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(None)
	err := c.cc.Invoke(ctx, Shortener_DeleteRecord_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) GetStats(ctx context.Context, in *None, opts ...grpc.CallOption) (*GetStatsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetStatsResponse)
	err := c.cc.Invoke(ctx, Shortener_GetStats_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ShortenerServer is the server API for Shortener service.
// All implementations must embed UnimplementedShortenerServer
// for forward compatibility.
type ShortenerServer interface {
	CreateShortening(context.Context, *CreateShorteningRequest) (*ShorteningResponse, error)
	GetFullString(context.Context, *LongURLRequest) (*LongURLResponse, error)
	CreateShorteningBatch(context.Context, *CreateShorteningBatchRequest) (*ShorteningBatchResponse, error)
	GetUserAllShortenings(context.Context, *UserID) (*ShorteningBatchResponse, error)
	DeleteRecord(context.Context, *DeleteRecordRequest) (*None, error)
	GetStats(context.Context, *None) (*GetStatsResponse, error)
	mustEmbedUnimplementedShortenerServer()
}

// UnimplementedShortenerServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedShortenerServer struct{}

func (UnimplementedShortenerServer) CreateShortening(context.Context, *CreateShorteningRequest) (*ShorteningResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateShortening not implemented")
}
func (UnimplementedShortenerServer) GetFullString(context.Context, *LongURLRequest) (*LongURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFullString not implemented")
}
func (UnimplementedShortenerServer) CreateShorteningBatch(context.Context, *CreateShorteningBatchRequest) (*ShorteningBatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateShorteningBatch not implemented")
}
func (UnimplementedShortenerServer) GetUserAllShortenings(context.Context, *UserID) (*ShorteningBatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserAllShortenings not implemented")
}
func (UnimplementedShortenerServer) DeleteRecord(context.Context, *DeleteRecordRequest) (*None, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteRecord not implemented")
}
func (UnimplementedShortenerServer) GetStats(context.Context, *None) (*GetStatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStats not implemented")
}
func (UnimplementedShortenerServer) mustEmbedUnimplementedShortenerServer() {}
func (UnimplementedShortenerServer) testEmbeddedByValue()                   {}

// UnsafeShortenerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ShortenerServer will
// result in compilation errors.
type UnsafeShortenerServer interface {
	mustEmbedUnimplementedShortenerServer()
}

func RegisterShortenerServer(s grpc.ServiceRegistrar, srv ShortenerServer) {
	// If the following call pancis, it indicates UnimplementedShortenerServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Shortener_ServiceDesc, srv)
}

func _Shortener_CreateShortening_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateShorteningRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).CreateShortening(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Shortener_CreateShortening_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).CreateShortening(ctx, req.(*CreateShorteningRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_GetFullString_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LongURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).GetFullString(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Shortener_GetFullString_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).GetFullString(ctx, req.(*LongURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_CreateShorteningBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateShorteningBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).CreateShorteningBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Shortener_CreateShorteningBatch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).CreateShorteningBatch(ctx, req.(*CreateShorteningBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_GetUserAllShortenings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).GetUserAllShortenings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Shortener_GetUserAllShortenings_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).GetUserAllShortenings(ctx, req.(*UserID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_DeleteRecord_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRecordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).DeleteRecord(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Shortener_DeleteRecord_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).DeleteRecord(ctx, req.(*DeleteRecordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_GetStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(None)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).GetStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Shortener_GetStats_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).GetStats(ctx, req.(*None))
	}
	return interceptor(ctx, in, info, handler)
}

// Shortener_ServiceDesc is the grpc.ServiceDesc for Shortener service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Shortener_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "shortener.Shortener",
	HandlerType: (*ShortenerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateShortening",
			Handler:    _Shortener_CreateShortening_Handler,
		},
		{
			MethodName: "GetFullString",
			Handler:    _Shortener_GetFullString_Handler,
		},
		{
			MethodName: "CreateShorteningBatch",
			Handler:    _Shortener_CreateShorteningBatch_Handler,
		},
		{
			MethodName: "GetUserAllShortenings",
			Handler:    _Shortener_GetUserAllShortenings_Handler,
		},
		{
			MethodName: "DeleteRecord",
			Handler:    _Shortener_DeleteRecord_Handler,
		},
		{
			MethodName: "GetStats",
			Handler:    _Shortener_GetStats_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/shortener.proto",
}
