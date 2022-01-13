package websocket

import (
	"context"
	"io"
	"net/url"

	"golang.org/x/net/websocket"
)

type Client struct {
	tunnel string
	origin string
}

func NewClient(tunnel string) (*Client, error) {
	URL, err := url.Parse(tunnel)
	if err != nil {
		return nil, err
	}
	switch URL.Scheme {
	case "ws":
		URL.Scheme = "http"
	case "wss":
		URL.Scheme = "https"
	}

	return &Client{
		tunnel: tunnel,
		origin: URL.String(),
	}, nil
}

func (t *Client) Connect(ctx context.Context) (io.ReadWriteCloser, error) {
	ws, err := websocket.Dial(t.tunnel, "", t.origin)
	if err != nil {
		return nil, err
	}

	return ws, nil
}
