package grpc

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"net/url"

	"github.com/isayme/tox/proto"
	"github.com/isayme/tox/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type Client struct {
	opts            util.ToxOptions
	grpcDialHost    string
	grpcDialOptions []grpc.DialOption
}

func NewClient(opts util.ToxOptions) (*Client, error) {
	URL, err := url.Parse(opts.Tunnel)
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
			grpc.WithPerRPCCredentials(newJwtToken(opts.Password, opts.InsecureSkipVerify)),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
	case "grpcs":
		dialOptions = append(
			dialOptions,
			grpc.WithPerRPCCredentials(newJwtToken(opts.Password, opts.InsecureSkipVerify)),
			grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
				InsecureSkipVerify: opts.InsecureSkipVerify,
			})),
		)
	}

	return &Client{
		opts:            opts,
		grpcDialHost:    URL.Host,
		grpcDialOptions: dialOptions,
	}, nil
}

func (t *Client) Connect(ctx context.Context) (util.ToxConn, error) {
	conn, err := grpc.DialContext(ctx, t.grpcDialHost, t.grpcDialOptions...)
	if err != nil {
		return nil, err
	}

	client := proto.NewTunnelClient(conn)
	c, err := client.OnConnect(ctx)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return util.NewToxConnection(NewGrpcClientConn(c, conn)), nil
}

type GrpcClientConn struct {
	c      proto.Tunnel_OnConnectClient
	conn   *grpc.ClientConn
	buffer *bytes.Buffer
}

func NewGrpcClientConn(c proto.Tunnel_OnConnectClient, conn *grpc.ClientConn) io.ReadWriteCloser {
	return &GrpcClientConn{
		c:      c,
		conn:   conn,
		buffer: bytes.NewBuffer(nil),
	}
}

func (conn *GrpcClientConn) Read(p []byte) (int, error) {
	if conn.buffer.Len() > 0 {
		return conn.buffer.Read(p)
	}

	d, err := conn.c.Recv()
	if err != nil {
		return 0, err
	}
	conn.buffer.Reset()
	conn.buffer.Write(d.Data)

	return conn.buffer.Read(p)
}

func (conn *GrpcClientConn) Write(p []byte) (int, error) {
	d := &proto.Data{
		Data: p,
	}

	err := conn.c.Send(d)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (conn *GrpcClientConn) Close() error {
	conn.c.CloseSend()
	return conn.conn.Close()
}
