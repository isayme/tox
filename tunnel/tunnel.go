package tunnel

import (
	"context"
	"fmt"
	"net/url"

	"github.com/isayme/go-logger"
	"github.com/isayme/tox/tunnel/grpc"
	"github.com/isayme/tox/tunnel/h2"
	"github.com/isayme/tox/tunnel/websocket"
	"github.com/isayme/tox/util"
)

type Client interface {
	Connect(context.Context) (util.LocalConn, error)
}

type Server interface {
	ListenAndServe(handler func(util.ServerConn)) error
	ListenAndServeTLS(certFile, keyFile string, handler func(util.ServerConn)) error
}

func NewClient(tunnel string, password string) (Client, error) {
	URL, err := url.Parse(tunnel)
	if err != nil {
		return nil, err
	}

	password = util.HashedPassword(password)

	logger.Infof("tunnel: %s", tunnel)

	switch URL.Scheme {
	case "grpc", "grpcs":
		return grpc.NewClient(tunnel, password)
	case "http2", "h2":
		return h2.NewClient(tunnel, password)
	case "ws", "wss":
		return websocket.NewClient(tunnel, password)
	}

	return nil, fmt.Errorf("not supported schema: %s", URL.Scheme)
}

func NewServer(tunnel string, password string) (Server, error) {
	URL, err := url.Parse(tunnel)
	if err != nil {
		return nil, err
	}

	password = util.HashedPassword(password)

	switch URL.Scheme {
	case "grpc", "grpcs":
		return grpc.NewServer(tunnel, password)
	case "http2", "h2":
		return h2.NewServer(tunnel, password)
	case "ws", "wss":
		return websocket.NewServer(tunnel, password)
	}
	return nil, fmt.Errorf("not supported schema: %s", URL.Scheme)
}
