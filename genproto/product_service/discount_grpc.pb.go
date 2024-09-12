// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.1
// source: submodule/product_service/discount.proto

package product_service

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
	DiscountService_CreateDiscount_FullMethodName = "/product_service.DiscountService/CreateDiscount"
	DiscountService_GetDiscount_FullMethodName    = "/product_service.DiscountService/GetDiscount"
	DiscountService_UpdateDiscount_FullMethodName = "/product_service.DiscountService/UpdateDiscount"
	DiscountService_DeleteDiscount_FullMethodName = "/product_service.DiscountService/DeleteDiscount"
	DiscountService_ListDiscounts_FullMethodName  = "/product_service.DiscountService/ListDiscounts"
)

// DiscountServiceClient is the client API for DiscountService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// DiscountService defines the gRPC service for managing discounts.
type DiscountServiceClient interface {
	CreateDiscount(ctx context.Context, in *CreateDiscountRequest, opts ...grpc.CallOption) (*CreateDiscountResponse, error)
	GetDiscount(ctx context.Context, in *GetDiscountRequest, opts ...grpc.CallOption) (*GetDiscountResponse, error)
	UpdateDiscount(ctx context.Context, in *UpdateDiscountRequest, opts ...grpc.CallOption) (*UpdateDiscountResponse, error)
	DeleteDiscount(ctx context.Context, in *DeleteDiscountRequest, opts ...grpc.CallOption) (*DeleteDiscountResponse, error)
	ListDiscounts(ctx context.Context, in *ListDiscountsRequest, opts ...grpc.CallOption) (*ListDiscountsResponse, error)
}

type discountServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDiscountServiceClient(cc grpc.ClientConnInterface) DiscountServiceClient {
	return &discountServiceClient{cc}
}

func (c *discountServiceClient) CreateDiscount(ctx context.Context, in *CreateDiscountRequest, opts ...grpc.CallOption) (*CreateDiscountResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateDiscountResponse)
	err := c.cc.Invoke(ctx, DiscountService_CreateDiscount_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *discountServiceClient) GetDiscount(ctx context.Context, in *GetDiscountRequest, opts ...grpc.CallOption) (*GetDiscountResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetDiscountResponse)
	err := c.cc.Invoke(ctx, DiscountService_GetDiscount_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *discountServiceClient) UpdateDiscount(ctx context.Context, in *UpdateDiscountRequest, opts ...grpc.CallOption) (*UpdateDiscountResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateDiscountResponse)
	err := c.cc.Invoke(ctx, DiscountService_UpdateDiscount_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *discountServiceClient) DeleteDiscount(ctx context.Context, in *DeleteDiscountRequest, opts ...grpc.CallOption) (*DeleteDiscountResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteDiscountResponse)
	err := c.cc.Invoke(ctx, DiscountService_DeleteDiscount_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *discountServiceClient) ListDiscounts(ctx context.Context, in *ListDiscountsRequest, opts ...grpc.CallOption) (*ListDiscountsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListDiscountsResponse)
	err := c.cc.Invoke(ctx, DiscountService_ListDiscounts_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DiscountServiceServer is the server API for DiscountService service.
// All implementations must embed UnimplementedDiscountServiceServer
// for forward compatibility.
//
// DiscountService defines the gRPC service for managing discounts.
type DiscountServiceServer interface {
	CreateDiscount(context.Context, *CreateDiscountRequest) (*CreateDiscountResponse, error)
	GetDiscount(context.Context, *GetDiscountRequest) (*GetDiscountResponse, error)
	UpdateDiscount(context.Context, *UpdateDiscountRequest) (*UpdateDiscountResponse, error)
	DeleteDiscount(context.Context, *DeleteDiscountRequest) (*DeleteDiscountResponse, error)
	ListDiscounts(context.Context, *ListDiscountsRequest) (*ListDiscountsResponse, error)
	mustEmbedUnimplementedDiscountServiceServer()
}

// UnimplementedDiscountServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedDiscountServiceServer struct{}

func (UnimplementedDiscountServiceServer) CreateDiscount(context.Context, *CreateDiscountRequest) (*CreateDiscountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateDiscount not implemented")
}
func (UnimplementedDiscountServiceServer) GetDiscount(context.Context, *GetDiscountRequest) (*GetDiscountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDiscount not implemented")
}
func (UnimplementedDiscountServiceServer) UpdateDiscount(context.Context, *UpdateDiscountRequest) (*UpdateDiscountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateDiscount not implemented")
}
func (UnimplementedDiscountServiceServer) DeleteDiscount(context.Context, *DeleteDiscountRequest) (*DeleteDiscountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteDiscount not implemented")
}
func (UnimplementedDiscountServiceServer) ListDiscounts(context.Context, *ListDiscountsRequest) (*ListDiscountsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListDiscounts not implemented")
}
func (UnimplementedDiscountServiceServer) mustEmbedUnimplementedDiscountServiceServer() {}
func (UnimplementedDiscountServiceServer) testEmbeddedByValue()                         {}

// UnsafeDiscountServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DiscountServiceServer will
// result in compilation errors.
type UnsafeDiscountServiceServer interface {
	mustEmbedUnimplementedDiscountServiceServer()
}

func RegisterDiscountServiceServer(s grpc.ServiceRegistrar, srv DiscountServiceServer) {
	// If the following call pancis, it indicates UnimplementedDiscountServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&DiscountService_ServiceDesc, srv)
}

func _DiscountService_CreateDiscount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateDiscountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscountServiceServer).CreateDiscount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DiscountService_CreateDiscount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscountServiceServer).CreateDiscount(ctx, req.(*CreateDiscountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DiscountService_GetDiscount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDiscountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscountServiceServer).GetDiscount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DiscountService_GetDiscount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscountServiceServer).GetDiscount(ctx, req.(*GetDiscountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DiscountService_UpdateDiscount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateDiscountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscountServiceServer).UpdateDiscount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DiscountService_UpdateDiscount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscountServiceServer).UpdateDiscount(ctx, req.(*UpdateDiscountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DiscountService_DeleteDiscount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteDiscountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscountServiceServer).DeleteDiscount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DiscountService_DeleteDiscount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscountServiceServer).DeleteDiscount(ctx, req.(*DeleteDiscountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DiscountService_ListDiscounts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListDiscountsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscountServiceServer).ListDiscounts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DiscountService_ListDiscounts_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscountServiceServer).ListDiscounts(ctx, req.(*ListDiscountsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DiscountService_ServiceDesc is the grpc.ServiceDesc for DiscountService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DiscountService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "product_service.DiscountService",
	HandlerType: (*DiscountServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateDiscount",
			Handler:    _DiscountService_CreateDiscount_Handler,
		},
		{
			MethodName: "GetDiscount",
			Handler:    _DiscountService_GetDiscount_Handler,
		},
		{
			MethodName: "UpdateDiscount",
			Handler:    _DiscountService_UpdateDiscount_Handler,
		},
		{
			MethodName: "DeleteDiscount",
			Handler:    _DiscountService_DeleteDiscount_Handler,
		},
		{
			MethodName: "ListDiscounts",
			Handler:    _DiscountService_ListDiscounts_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "submodule/product_service/discount.proto",
}
