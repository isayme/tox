package websocket

import (
	"context"
	"crypto/tls"
	"io"
	"net/url"

	"golang.org/x/net/websocket"
)

type Client struct {
	tunnel string
	origin string
}

func NewClient(tunnel string, password string) (*Client, error) {
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
	c, err := websocket.NewConfig(t.tunnel, t.origin)
	if err != nil {
		return nil, err
	}
	c.TlsConfig = &tls.Config{InsecureSkipVerify: true}
	ws, err := websocket.DialConfig(c)
	if err != nil {
		return nil, err
	}
	return ws, nil
}
