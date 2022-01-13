package tunnel

import (
	"context"
	"fmt"
	"io"
	"net/url"

	"github.com/isayme/tox/tunnel/grpc"
	"github.com/isayme/tox/tunnel/h2"
	"github.com/isayme/tox/tunnel/quic"
	"github.com/isayme/tox/tunnel/websocket"
)

type Client interface {
	Connect(context.Context) (io.ReadWriteCloser, error)
}

type Server interface {
	ListenAndServeTLS(certFile, keyFile string, handler func(io.ReadWriter)) error
}

func NewClient(tunnel string) (Client, error) {
	URL, err := url.Parse(tunnel)
	if err != nil {
		return nil, err
	}

	switch URL.Scheme {
	case "quic", "http3":
		return quic.NewClient(tunnel)
	case "grpc", "grpcs":
		return grpc.NewClient(tunnel)
	case "http2", "h2":
		return h2.NewClient(tunnel)
	case "ws", "wss":
		return websocket.NewClient(tunnel)
	}

	return nil, fmt.Errorf("not supported schema: %s", URL.Scheme)
}

func NewServer(tunnel string) (Server, error) {
	URL, err := url.Parse(tunnel)
	if err != nil {
		return nil, err
	}

	switch URL.Scheme {
	case "quic", "http3":
		return quic.NewServer(tunnel)
	case "grpc", "grpcs":
		return grpc.NewServer(tunnel)
	case "http2", "h2":
		return h2.NewServer(tunnel)
	case "ws", "wss":
		return websocket.NewServer(tunnel)
	}
	return nil, fmt.Errorf("not supported schema: %s", URL.Scheme)
}
