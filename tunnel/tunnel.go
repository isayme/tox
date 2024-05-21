package tunnel

import (
	"context"
	"fmt"
	"net/url"

	"github.com/isayme/go-logger"
	"github.com/isayme/tox/tunnel/grpc"
	"github.com/isayme/tox/tunnel/h2"
	"github.com/isayme/tox/tunnel/quic"
	"github.com/isayme/tox/tunnel/websocket"
	"github.com/isayme/tox/util"
)

type Client interface {
	Connect(context.Context) (util.ToxConn, error)
}

type Server interface {
	ListenAndServe(handler func(util.ToxConn)) error
}

func NewClient(opts util.ToxOptions) (Client, error) {
	URL, err := url.Parse(opts.Tunnel)
	if err != nil {
		return nil, err
	}

	logger.Infof("tunnel: %s", opts.Tunnel)

	switch URL.Scheme {
	case "grpc", "grpcs":
		return grpc.NewClient(opts)
	case "http2", "h2":
		return h2.NewClient(opts)
	case "ws", "wss":
		return websocket.NewClient(opts)
	case "quic", "http3":
		return quic.NewClient(opts)
	}

	return nil, fmt.Errorf("not supported schema: %s", URL.Scheme)
}

func NewServer(opts util.ToxOptions) (Server, error) {
	URL, err := url.Parse(opts.Tunnel)
	if err != nil {
		return nil, err
	}

	switch URL.Scheme {
	case "grpc", "grpcs":
		return grpc.NewServer(opts)
	case "http2", "h2":
		return h2.NewServer(opts)
	case "ws", "wss":
		return websocket.NewServer(opts)
	case "quic", "http3":
		return quic.NewServer(opts)
	}
	return nil, fmt.Errorf("not supported schema: %s", URL.Scheme)
}
