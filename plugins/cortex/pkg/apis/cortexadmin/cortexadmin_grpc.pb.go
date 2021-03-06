// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - ragù               v0.2.3
// source: plugins/cortex/pkg/apis/cortexadmin/cortexadmin.proto

package cortexadmin

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

// CortexAdminClient is the client API for CortexAdmin service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CortexAdminClient interface {
	AllUserStats(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*UserIDStatsList, error)
	WriteMetrics(ctx context.Context, in *WriteRequest, opts ...grpc.CallOption) (*WriteResponse, error)
	Query(ctx context.Context, in *QueryRequest, opts ...grpc.CallOption) (*QueryResponse, error)
	QueryRange(ctx context.Context, in *QueryRangeRequest, opts ...grpc.CallOption) (*QueryResponse, error)
}

type cortexAdminClient struct {
	cc grpc.ClientConnInterface
}

func NewCortexAdminClient(cc grpc.ClientConnInterface) CortexAdminClient {
	return &cortexAdminClient{cc}
}

func (c *cortexAdminClient) AllUserStats(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*UserIDStatsList, error) {
	out := new(UserIDStatsList)
	err := c.cc.Invoke(ctx, "/cortexadmin.CortexAdmin/AllUserStats", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cortexAdminClient) WriteMetrics(ctx context.Context, in *WriteRequest, opts ...grpc.CallOption) (*WriteResponse, error) {
	out := new(WriteResponse)
	err := c.cc.Invoke(ctx, "/cortexadmin.CortexAdmin/WriteMetrics", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cortexAdminClient) Query(ctx context.Context, in *QueryRequest, opts ...grpc.CallOption) (*QueryResponse, error) {
	out := new(QueryResponse)
	err := c.cc.Invoke(ctx, "/cortexadmin.CortexAdmin/Query", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cortexAdminClient) QueryRange(ctx context.Context, in *QueryRangeRequest, opts ...grpc.CallOption) (*QueryResponse, error) {
	out := new(QueryResponse)
	err := c.cc.Invoke(ctx, "/cortexadmin.CortexAdmin/QueryRange", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CortexAdminServer is the server API for CortexAdmin service.
// All implementations must embed UnimplementedCortexAdminServer
// for forward compatibility
type CortexAdminServer interface {
	AllUserStats(context.Context, *emptypb.Empty) (*UserIDStatsList, error)
	WriteMetrics(context.Context, *WriteRequest) (*WriteResponse, error)
	Query(context.Context, *QueryRequest) (*QueryResponse, error)
	QueryRange(context.Context, *QueryRangeRequest) (*QueryResponse, error)
	mustEmbedUnimplementedCortexAdminServer()
}

// UnimplementedCortexAdminServer must be embedded to have forward compatible implementations.
type UnimplementedCortexAdminServer struct {
}

func (UnimplementedCortexAdminServer) AllUserStats(context.Context, *emptypb.Empty) (*UserIDStatsList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AllUserStats not implemented")
}
func (UnimplementedCortexAdminServer) WriteMetrics(context.Context, *WriteRequest) (*WriteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WriteMetrics not implemented")
}
func (UnimplementedCortexAdminServer) Query(context.Context, *QueryRequest) (*QueryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Query not implemented")
}
func (UnimplementedCortexAdminServer) QueryRange(context.Context, *QueryRangeRequest) (*QueryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryRange not implemented")
}
func (UnimplementedCortexAdminServer) mustEmbedUnimplementedCortexAdminServer() {}

// UnsafeCortexAdminServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CortexAdminServer will
// result in compilation errors.
type UnsafeCortexAdminServer interface {
	mustEmbedUnimplementedCortexAdminServer()
}

func RegisterCortexAdminServer(s grpc.ServiceRegistrar, srv CortexAdminServer) {
	s.RegisterService(&CortexAdmin_ServiceDesc, srv)
}

func _CortexAdmin_AllUserStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CortexAdminServer).AllUserStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cortexadmin.CortexAdmin/AllUserStats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CortexAdminServer).AllUserStats(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _CortexAdmin_WriteMetrics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WriteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CortexAdminServer).WriteMetrics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cortexadmin.CortexAdmin/WriteMetrics",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CortexAdminServer).WriteMetrics(ctx, req.(*WriteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CortexAdmin_Query_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CortexAdminServer).Query(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cortexadmin.CortexAdmin/Query",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CortexAdminServer).Query(ctx, req.(*QueryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CortexAdmin_QueryRange_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryRangeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CortexAdminServer).QueryRange(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cortexadmin.CortexAdmin/QueryRange",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CortexAdminServer).QueryRange(ctx, req.(*QueryRangeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CortexAdmin_ServiceDesc is the grpc.ServiceDesc for CortexAdmin service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CortexAdmin_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cortexadmin.CortexAdmin",
	HandlerType: (*CortexAdminServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AllUserStats",
			Handler:    _CortexAdmin_AllUserStats_Handler,
		},
		{
			MethodName: "WriteMetrics",
			Handler:    _CortexAdmin_WriteMetrics_Handler,
		},
		{
			MethodName: "Query",
			Handler:    _CortexAdmin_Query_Handler,
		},
		{
			MethodName: "QueryRange",
			Handler:    _CortexAdmin_QueryRange_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "plugins/cortex/pkg/apis/cortexadmin/cortexadmin.proto",
}
