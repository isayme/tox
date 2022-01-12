package grpc

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"net/url"

	"github.com/isayme/go-logger"
	"github.com/isayme/tox/proto"
	"github.com/shimingyah/pool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

type Client struct {
	tunnel string
	p      pool.Pool
}

func NewClient(tunnel string) (*Client, error) {
	URL, err := url.Parse(tunnel)
	if err != nil {
		return nil, err
	}

	dialOptions := []grpc.DialOption{
		grpc.WithBackoffMaxDelay(pool.BackoffMaxDelay),
		grpc.WithInitialWindowSize(pool.InitialWindowSize),
		grpc.WithInitialConnWindowSize(pool.InitialConnWindowSize),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(2 << 30)),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(2 << 30)),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                pool.KeepAliveTime,
			Timeout:             pool.KeepAliveTimeout,
			PermitWithoutStream: true,
		}),
	}

	switch URL.Scheme {
	case "grpc":
		dialOptions = append(dialOptions, grpc.WithInsecure())
	case "grpcs":
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})))
	default:
		logger.Panicw("URL.Scheme invalid", "address", tunnel, "sceham", URL.Scheme)
	}

	options := pool.DefaultOptions
	options.Dial = func(address string) (*grpc.ClientConn, error) {
		ctx, cancel := context.WithTimeout(context.Background(), pool.DialTimeout)
		defer cancel()
		return grpc.DialContext(ctx, address, dialOptions...)
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

func (t *Client) Connect(ctx context.Context) (io.ReadWriteCloser, error) {
	conn, err := t.p.Get()
	if err != nil {
		return nil, err
	}

	client := proto.NewTunnelClient(conn.Value())
	c, err := client.OnConnect(ctx)
	if err != nil {
		return nil, err
	}

	return NewClientReadWriter(c), nil
}

type clientReadWriter struct {
	c      proto.Tunnel_OnConnectClient
	buffer *bytes.Buffer
}

func NewClientReadWriter(c proto.Tunnel_OnConnectClient) *clientReadWriter {
	return &clientReadWriter{
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
	return rw.c.CloseSend()
}
