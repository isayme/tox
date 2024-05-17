package websocket

import (
	"context"
	"crypto/tls"
	"net/url"

	"github.com/isayme/tox/util"
	"golang.org/x/net/websocket"
)

type Client struct {
	config *websocket.Config
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
	origin := URL.String()

	wsc, err := websocket.NewConfig(tunnel, origin)
	if err != nil {
		return nil, err
	}
	if password != "" {
		wsc.Header.Set("token", password)
	}

	wsc.TlsConfig = &tls.Config{InsecureSkipVerify: true}

	return &Client{
		config: wsc,
	}, nil
}

func (t *Client) Connect(ctx context.Context) (util.LocalConn, error) {
	ws, err := websocket.DialConfig(t.config)
	if err != nil {
		return nil, err
	}
	return &wsLocalConn{Conn: ws}, nil
}

type wsLocalConn struct {
	*websocket.Conn
}

func (conn *wsLocalConn) CloseWrite() error {
	return conn.Conn.Close()
}
