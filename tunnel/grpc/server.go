package grpc

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/url"

	"github.com/isayme/tox/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Server struct {
	proto.UnimplementedTunnelServer

	handler func(io.ReadWriter)
	tunnel  string
	key     []byte
}

func NewServer(tunnel string, password string) (*Server, error) {
	return &Server{
		tunnel: tunnel,
		key:    []byte(password),
	}, nil
}

func (s *Server) ListenAndServeTLS(certFile, keyFile string, handler func(io.ReadWriter)) error {
	URL, err := url.Parse(s.tunnel)
	if err != nil {
		return err
	}

	s.handler = handler

	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		return err
	}

	addr := fmt.Sprintf(":%s", URL.Port())
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	grpcs := grpc.NewServer(grpc.Creds(creds))
	proto.RegisterTunnelServer(grpcs, s)

	return grpcs.Serve(l)
}

func (s *Server) OnConnect(stream proto.Tunnel_OnConnectServer) error {
	err := VerifyTokenFromContext(stream.Context(), s.key)
	if err != nil {
		return err
	}

	rw := NewServerReadWriter(stream)

	s.handler(rw)
	return nil
}

type serverReadWriter struct {
	c      proto.Tunnel_OnConnectServer
	buffer *bytes.Buffer
}

func NewServerReadWriter(c proto.Tunnel_OnConnectServer) *serverReadWriter {
	return &serverReadWriter{
		c:      c,
		buffer: bytes.NewBuffer(nil),
	}
}

func (rw *serverReadWriter) Read(p []byte) (int, error) {
	if rw.buffer.Len() > 0 {
		return rw.buffer.Read(p)
	}

	d, err := rw.c.Recv()
	if err != nil {
		return 0, err
	}
	rw.buffer.Write(d.Data)

	return rw.buffer.Read(p)
}

func (rw *serverReadWriter) Write(p []byte) (int, error) {
	d := &proto.Data{
		Data: p,
	}

	err := rw.c.Send(d)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}
