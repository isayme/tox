package grpc

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/url"

	"github.com/isayme/tox/proto"
	"github.com/isayme/tox/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Server struct {
	proto.UnimplementedTunnelServer

	handler func(util.ToxConn)

	opts util.ToxOptions
}

func NewServer(opts util.ToxOptions) (*Server, error) {
	return &Server{
		opts: opts,
	}, nil
}

func (s *Server) ListenAndServe(handler func(util.ToxConn)) error {
	URL, err := url.Parse(s.opts.Tunnel)
	if err != nil {
		return err
	}

	s.handler = handler

	addr := fmt.Sprintf(":%s", URL.Port())
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	options := make([]grpc.ServerOption, 0)

	certFile := s.opts.CertFile
	keyFile := s.opts.KeyFile
	if certFile != "" && keyFile != "" {
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			return err
		}

		options = append(options, grpc.Creds(creds))
	}

	options = append(options, grpc.MaxRecvMsgSize(MaxRecvMsgSize))
	options = append(options, grpc.MaxSendMsgSize(MaxSendMsgSize))
	options = append(options, grpc.ReadBufferSize(ReadBufferSize))
	options = append(options, grpc.WriteBufferSize(WriteBufferSize))
	options = append(options, grpc.ConnectionTimeout(s.opts.ConnectTimeout))

	grpcs := grpc.NewServer(options...)
	proto.RegisterTunnelServer(grpcs, s)

	return grpcs.Serve(l)
}

func (s *Server) OnConnect(stream proto.Tunnel_OnConnectServer) error {
	err := VerifyTokenFromContext(stream.Context(), s.opts.Password)
	if err != nil {
		return err
	}

	s.handler(util.NewToxConnection(NewGrpcServerConn(stream)))
	return nil
}

type GrpcServerConn struct {
	conn   proto.Tunnel_OnConnectServer
	buffer *bytes.Buffer
}

func NewGrpcServerConn(conn proto.Tunnel_OnConnectServer) io.ReadWriteCloser {
	return &GrpcServerConn{
		conn:   conn,
		buffer: bytes.NewBuffer(nil),
	}
}

func (conn *GrpcServerConn) Read(p []byte) (int, error) {
	if conn.buffer.Len() > 0 {
		return conn.buffer.Read(p)
	}

	d, err := conn.conn.Recv()
	if err != nil {
		return 0, err
	}
	conn.buffer.Reset()
	conn.buffer.Write(d.Data)

	return conn.buffer.Read(p)
}

func (conn *GrpcServerConn) Write(p []byte) (int, error) {
	d := &proto.Data{
		Data: p,
	}

	err := conn.conn.Send(d)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (conn *GrpcServerConn) Close() error {
	return nil
}
