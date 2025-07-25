// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: order.proto

package pb

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
	OrderService_AddProduct_FullMethodName   = "/orderpb.OrderService/AddProduct"
	OrderService_DeleteDetail_FullMethodName = "/orderpb.OrderService/DeleteDetail"
	OrderService_GetUserCart_FullMethodName  = "/orderpb.OrderService/GetUserCart"
	OrderService_PlaceOrder_FullMethodName   = "/orderpb.OrderService/PlaceOrder"
)

// OrderServiceClient is the client API for OrderService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OrderServiceClient interface {
	AddProduct(ctx context.Context, in *AddProductRequest, opts ...grpc.CallOption) (*CommonResponse, error)
	DeleteDetail(ctx context.Context, in *DeleteDetailRequest, opts ...grpc.CallOption) (*CommonResponse, error)
	GetUserCart(ctx context.Context, in *UserCartRequest, opts ...grpc.CallOption) (*UserCartResponse, error)
	PlaceOrder(ctx context.Context, in *PlaceOrderRequest, opts ...grpc.CallOption) (*PlaceOrderResponse, error)
}

type orderServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewOrderServiceClient(cc grpc.ClientConnInterface) OrderServiceClient {
	return &orderServiceClient{cc}
}

func (c *orderServiceClient) AddProduct(ctx context.Context, in *AddProductRequest, opts ...grpc.CallOption) (*CommonResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CommonResponse)
	err := c.cc.Invoke(ctx, OrderService_AddProduct_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderServiceClient) DeleteDetail(ctx context.Context, in *DeleteDetailRequest, opts ...grpc.CallOption) (*CommonResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CommonResponse)
	err := c.cc.Invoke(ctx, OrderService_DeleteDetail_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderServiceClient) GetUserCart(ctx context.Context, in *UserCartRequest, opts ...grpc.CallOption) (*UserCartResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UserCartResponse)
	err := c.cc.Invoke(ctx, OrderService_GetUserCart_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderServiceClient) PlaceOrder(ctx context.Context, in *PlaceOrderRequest, opts ...grpc.CallOption) (*PlaceOrderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PlaceOrderResponse)
	err := c.cc.Invoke(ctx, OrderService_PlaceOrder_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrderServiceServer is the server API for OrderService service.
// All implementations should embed UnimplementedOrderServiceServer
// for forward compatibility.
type OrderServiceServer interface {
	AddProduct(context.Context, *AddProductRequest) (*CommonResponse, error)
	DeleteDetail(context.Context, *DeleteDetailRequest) (*CommonResponse, error)
	GetUserCart(context.Context, *UserCartRequest) (*UserCartResponse, error)
	PlaceOrder(context.Context, *PlaceOrderRequest) (*PlaceOrderResponse, error)
}

// UnimplementedOrderServiceServer should be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedOrderServiceServer struct{}

func (UnimplementedOrderServiceServer) AddProduct(context.Context, *AddProductRequest) (*CommonResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddProduct not implemented")
}
func (UnimplementedOrderServiceServer) DeleteDetail(context.Context, *DeleteDetailRequest) (*CommonResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteDetail not implemented")
}
func (UnimplementedOrderServiceServer) GetUserCart(context.Context, *UserCartRequest) (*UserCartResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserCart not implemented")
}
func (UnimplementedOrderServiceServer) PlaceOrder(context.Context, *PlaceOrderRequest) (*PlaceOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PlaceOrder not implemented")
}
func (UnimplementedOrderServiceServer) testEmbeddedByValue() {}

// UnsafeOrderServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OrderServiceServer will
// result in compilation errors.
type UnsafeOrderServiceServer interface {
	mustEmbedUnimplementedOrderServiceServer()
}

func RegisterOrderServiceServer(s grpc.ServiceRegistrar, srv OrderServiceServer) {
	// If the following call pancis, it indicates UnimplementedOrderServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&OrderService_ServiceDesc, srv)
}

func _OrderService_AddProduct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddProductRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServiceServer).AddProduct(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OrderService_AddProduct_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServiceServer).AddProduct(ctx, req.(*AddProductRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrderService_DeleteDetail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteDetailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServiceServer).DeleteDetail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OrderService_DeleteDetail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServiceServer).DeleteDetail(ctx, req.(*DeleteDetailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrderService_GetUserCart_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserCartRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServiceServer).GetUserCart(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OrderService_GetUserCart_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServiceServer).GetUserCart(ctx, req.(*UserCartRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrderService_PlaceOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PlaceOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServiceServer).PlaceOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OrderService_PlaceOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServiceServer).PlaceOrder(ctx, req.(*PlaceOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// OrderService_ServiceDesc is the grpc.ServiceDesc for OrderService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OrderService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "orderpb.OrderService",
	HandlerType: (*OrderServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddProduct",
			Handler:    _OrderService_AddProduct_Handler,
		},
		{
			MethodName: "DeleteDetail",
			Handler:    _OrderService_DeleteDetail_Handler,
		},
		{
			MethodName: "GetUserCart",
			Handler:    _OrderService_GetUserCart_Handler,
		},
		{
			MethodName: "PlaceOrder",
			Handler:    _OrderService_PlaceOrder_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "order.proto",
}
