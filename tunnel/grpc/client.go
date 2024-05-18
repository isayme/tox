package grpc

import (
	"bytes"
	"context"
	"crypto/tls"
	"net/url"

	pool "github.com/isayme/go-grpcpool"
	"github.com/isayme/go-logger"
	"github.com/isayme/tox/proto"
	"github.com/isayme/tox/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type Client struct {
	tunnel string
	p      pool.Pool
}

func NewClient(tunnel string, password string) (*Client, error) {
	URL, err := url.Parse(tunnel)
	if err != nil {
		return nil, err
	}

	dialOptions := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithBackoffMaxDelay(BackoffMaxDelay),
		grpc.WithInitialWindowSize(InitialWindowSize),
		grpc.WithInitialConnWindowSize(InitialConnWindowSize),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(MaxSendMsgSize)),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MaxRecvMsgSize)),
		grpc.WithReadBufferSize(ReadBufferSize),
		grpc.WithWriteBufferSize(WriteBufferSize),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                KeepAliveTime,
			Timeout:             KeepAliveTimeout,
			PermitWithoutStream: true,
		}),
	}

	switch URL.Scheme {
	case "grpc":
		dialOptions = append(
			dialOptions,
			grpc.WithPerRPCCredentials(newJwtToken(password, false)),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
	case "grpcs":
		dialOptions = append(
			dialOptions,
			grpc.WithPerRPCCredentials(newJwtToken(password, true)),
			grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})),
		)
	}

	options := pool.Options{
		Dial: func(address string) (*grpc.ClientConn, error) {
			ctx, cancel := context.WithTimeout(context.Background(), DialTimeout)
			defer cancel()
			return grpc.DialContext(ctx, address, dialOptions...)
		},
		MaxIdle:              MaxIdle,
		MaxActive:            MaxActive,
		MaxConcurrentStreams: MaxConcurrentStreams,
		Reuse:                true,
	}

	p, err := pool.New(URL.Host, options)
	if err != nil {
		logger.Errorf("failed to new pool: %v", err)
		return nil, err
	}

	return &Client{
		tunnel: tunnel,
		p:      p,
	}, nil
}

func (t *Client) Connect(ctx context.Context) (util.LocalConn, error) {
	conn, err := t.p.Get()
	if err != nil {
		return nil, err
	}

	client := proto.NewTunnelClient(conn.Value())
	c, err := client.OnConnect(context.Background())
	if err != nil {
		conn.Close()
		return nil, err
	}

	return NewClientReadWriter(conn, c), nil
}

type clientReadWriter struct {
	conn   pool.Conn
	c      proto.Tunnel_OnConnectClient
	buffer *bytes.Buffer
}

func NewClientReadWriter(conn pool.Conn, c proto.Tunnel_OnConnectClient) *clientReadWriter {
	return &clientReadWriter{
		conn:   conn,
		c:      c,
		buffer: bytes.NewBuffer(nil),
	}
}

func (rw *clientReadWriter) Read(p []byte) (int, error) {
	if rw.buffer.Len() > 0 {
		return rw.buffer.Read(p)
	}

	d, err := rw.c.Recv()
	if err != nil {
		return 0, err
	}
	rw.buffer.Reset()
	rw.buffer.Write(d.Data)

	return rw.buffer.Read(p)
}

func (rw *clientReadWriter) Write(p []byte) (int, error) {
	d := &proto.Data{
		Data: p,
	}

	err := rw.c.Send(d)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (rw *clientReadWriter) Close() error {
	rw.c.CloseSend()
	return rw.conn.Close()
}

func (rw *clientReadWriter) CloseWrite() error {
	return rw.c.CloseSend()
}
