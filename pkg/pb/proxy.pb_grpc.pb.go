// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// ProxyHttpServiceClient is the client API for ProxyHttpService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProxyHttpServiceClient interface {
	ProxyCloud2Edge(ctx context.Context, opts ...grpc.CallOption) (ProxyHttpService_ProxyCloud2EdgeClient, error)
	ProxyEdge2Cloud(ctx context.Context, opts ...grpc.CallOption) (ProxyHttpService_ProxyEdge2CloudClient, error)
}

type proxyHttpServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewProxyHttpServiceClient(cc grpc.ClientConnInterface) ProxyHttpServiceClient {
	return &proxyHttpServiceClient{cc}
}

func (c *proxyHttpServiceClient) ProxyCloud2Edge(ctx context.Context, opts ...grpc.CallOption) (ProxyHttpService_ProxyCloud2EdgeClient, error) {
	stream, err := c.cc.NewStream(ctx, &_ProxyHttpService_serviceDesc.Streams[0], "/pb.ProxyHttpService/ProxyCloud2Edge", opts...)
	if err != nil {
		return nil, err
	}
	x := &proxyHttpServiceProxyCloud2EdgeClient{stream}
	return x, nil
}

type ProxyHttpService_ProxyCloud2EdgeClient interface {
	Send(*Response) error
	Recv() (*Request, error)
	grpc.ClientStream
}

type proxyHttpServiceProxyCloud2EdgeClient struct {
	grpc.ClientStream
}

func (x *proxyHttpServiceProxyCloud2EdgeClient) Send(m *Response) error {
	return x.ClientStream.SendMsg(m)
}

func (x *proxyHttpServiceProxyCloud2EdgeClient) Recv() (*Request, error) {
	m := new(Request)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *proxyHttpServiceClient) ProxyEdge2Cloud(ctx context.Context, opts ...grpc.CallOption) (ProxyHttpService_ProxyEdge2CloudClient, error) {
	stream, err := c.cc.NewStream(ctx, &_ProxyHttpService_serviceDesc.Streams[1], "/pb.ProxyHttpService/ProxyEdge2Cloud", opts...)
	if err != nil {
		return nil, err
	}
	x := &proxyHttpServiceProxyEdge2CloudClient{stream}
	return x, nil
}

type ProxyHttpService_ProxyEdge2CloudClient interface {
	Send(*Request) error
	Recv() (*Response, error)
	grpc.ClientStream
}

type proxyHttpServiceProxyEdge2CloudClient struct {
	grpc.ClientStream
}

func (x *proxyHttpServiceProxyEdge2CloudClient) Send(m *Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *proxyHttpServiceProxyEdge2CloudClient) Recv() (*Response, error) {
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ProxyHttpServiceServer is the server API for ProxyHttpService service.
// All implementations must embed UnimplementedProxyHttpServiceServer
// for forward compatibility
type ProxyHttpServiceServer interface {
	ProxyCloud2Edge(ProxyHttpService_ProxyCloud2EdgeServer) error
	ProxyEdge2Cloud(ProxyHttpService_ProxyEdge2CloudServer) error
	mustEmbedUnimplementedProxyHttpServiceServer()
}

// UnimplementedProxyHttpServiceServer must be embedded to have forward compatible implementations.
type UnimplementedProxyHttpServiceServer struct {
}

func (UnimplementedProxyHttpServiceServer) ProxyCloud2Edge(ProxyHttpService_ProxyCloud2EdgeServer) error {
	return status.Errorf(codes.Unimplemented, "method ProxyCloud2Edge not implemented")
}
func (UnimplementedProxyHttpServiceServer) ProxyEdge2Cloud(ProxyHttpService_ProxyEdge2CloudServer) error {
	return status.Errorf(codes.Unimplemented, "method ProxyEdge2Cloud not implemented")
}
func (UnimplementedProxyHttpServiceServer) mustEmbedUnimplementedProxyHttpServiceServer() {}

// UnsafeProxyHttpServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProxyHttpServiceServer will
// result in compilation errors.
type UnsafeProxyHttpServiceServer interface {
	mustEmbedUnimplementedProxyHttpServiceServer()
}

func RegisterProxyHttpServiceServer(s grpc.ServiceRegistrar, srv ProxyHttpServiceServer) {
	s.RegisterService(&_ProxyHttpService_serviceDesc, srv)
}

func _ProxyHttpService_ProxyCloud2Edge_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ProxyHttpServiceServer).ProxyCloud2Edge(&proxyHttpServiceProxyCloud2EdgeServer{stream})
}

type ProxyHttpService_ProxyCloud2EdgeServer interface {
	Send(*Request) error
	Recv() (*Response, error)
	grpc.ServerStream
}

type proxyHttpServiceProxyCloud2EdgeServer struct {
	grpc.ServerStream
}

func (x *proxyHttpServiceProxyCloud2EdgeServer) Send(m *Request) error {
	return x.ServerStream.SendMsg(m)
}

func (x *proxyHttpServiceProxyCloud2EdgeServer) Recv() (*Response, error) {
	m := new(Response)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _ProxyHttpService_ProxyEdge2Cloud_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ProxyHttpServiceServer).ProxyEdge2Cloud(&proxyHttpServiceProxyEdge2CloudServer{stream})
}

type ProxyHttpService_ProxyEdge2CloudServer interface {
	Send(*Response) error
	Recv() (*Request, error)
	grpc.ServerStream
}

type proxyHttpServiceProxyEdge2CloudServer struct {
	grpc.ServerStream
}

func (x *proxyHttpServiceProxyEdge2CloudServer) Send(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *proxyHttpServiceProxyEdge2CloudServer) Recv() (*Request, error) {
	m := new(Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _ProxyHttpService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.ProxyHttpService",
	HandlerType: (*ProxyHttpServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ProxyCloud2Edge",
			Handler:       _ProxyHttpService_ProxyCloud2Edge_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "ProxyEdge2Cloud",
			Handler:       _ProxyHttpService_ProxyEdge2Cloud_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proxy.pb",
}
