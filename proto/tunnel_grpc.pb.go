// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

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

// TunnelClient is the client API for Tunnel service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TunnelClient interface {
	OnConnect(ctx context.Context, opts ...grpc.CallOption) (Tunnel_OnConnectClient, error)
}

type tunnelClient struct {
	cc grpc.ClientConnInterface
}

func NewTunnelClient(cc grpc.ClientConnInterface) TunnelClient {
	return &tunnelClient{cc}
}

func (c *tunnelClient) OnConnect(ctx context.Context, opts ...grpc.CallOption) (Tunnel_OnConnectClient, error) {
	stream, err := c.cc.NewStream(ctx, &Tunnel_ServiceDesc.Streams[0], "/Tunnel/OnConnect", opts...)
	if err != nil {
		return nil, err
	}
	x := &tunnelOnConnectClient{stream}
	return x, nil
}

type Tunnel_OnConnectClient interface {
	Send(*Data) error
	Recv() (*Data, error)
	grpc.ClientStream
}

type tunnelOnConnectClient struct {
	grpc.ClientStream
}

func (x *tunnelOnConnectClient) Send(m *Data) error {
	return x.ClientStream.SendMsg(m)
}

func (x *tunnelOnConnectClient) Recv() (*Data, error) {
	m := new(Data)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TunnelServer is the server API for Tunnel service.
// All implementations must embed UnimplementedTunnelServer
// for forward compatibility
type TunnelServer interface {
	OnConnect(Tunnel_OnConnectServer) error
	mustEmbedUnimplementedTunnelServer()
}

// UnimplementedTunnelServer must be embedded to have forward compatible implementations.
type UnimplementedTunnelServer struct {
}

func (UnimplementedTunnelServer) OnConnect(Tunnel_OnConnectServer) error {
	return status.Errorf(codes.Unimplemented, "method OnConnect not implemented")
}
func (UnimplementedTunnelServer) mustEmbedUnimplementedTunnelServer() {}

// UnsafeTunnelServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TunnelServer will
// result in compilation errors.
type UnsafeTunnelServer interface {
	mustEmbedUnimplementedTunnelServer()
}

func RegisterTunnelServer(s grpc.ServiceRegistrar, srv TunnelServer) {
	s.RegisterService(&Tunnel_ServiceDesc, srv)
}

func _Tunnel_OnConnect_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(TunnelServer).OnConnect(&tunnelOnConnectServer{stream})
}

type Tunnel_OnConnectServer interface {
	Send(*Data) error
	Recv() (*Data, error)
	grpc.ServerStream
}

type tunnelOnConnectServer struct {
	grpc.ServerStream
}

func (x *tunnelOnConnectServer) Send(m *Data) error {
	return x.ServerStream.SendMsg(m)
}

func (x *tunnelOnConnectServer) Recv() (*Data, error) {
	m := new(Data)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Tunnel_ServiceDesc is the grpc.ServiceDesc for Tunnel service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Tunnel_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Tunnel",
	HandlerType: (*TunnelServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "OnConnect",
			Handler:       _Tunnel_OnConnect_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/tunnel.proto",
}
