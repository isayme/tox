package websocket

import (
	"context"
	"crypto/tls"
	"net/url"

	"github.com/isayme/tox/util"
	"golang.org/x/net/websocket"
)

type Client struct {
	wsConfig *websocket.Config
}

func NewClient(opts util.ToxOptions) (*Client, error) {
	URL, err := url.Parse(opts.Tunnel)
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

	wsConfig, err := websocket.NewConfig(opts.Tunnel, origin)
	if err != nil {
		return nil, err
	}

	wsConfig.Header.Set("token", opts.Password)

	wsConfig.TlsConfig = &tls.Config{
		InsecureSkipVerify: opts.InsecureSkipVerify,
	}

	return &Client{
		wsConfig: wsConfig,
	}, nil
}

func (t *Client) Connect(ctx context.Context) (util.ToxConn, error) {
	ws, err := websocket.DialConfig(t.wsConfig)
	if err != nil {
		return nil, err
	}
	return util.NewToxConnection(ws), nil
}
