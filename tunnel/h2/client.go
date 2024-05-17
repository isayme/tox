package h2

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/isayme/tox/util"
	"github.com/posener/h2conn"
)

type Client struct {
	tunnel   string
	h2Client *h2conn.Client
}

func NewClient(tunnel string, password string) (*Client, error) {
	URL, err := url.Parse(tunnel)
	if err != nil {
		return nil, err
	}
	switch URL.Scheme {
	case "h2", "http2", "https":
		URL.Scheme = "https"
	}

	headers := http.Header{}
	if password != "" {
		headers.Add("token", password)
	}

	return &Client{
		tunnel: URL.String(),
		h2Client: &h2conn.Client{
			Method: http.MethodPost,
			Client: &http.Client{
				Transport: &http.Transport{
					DialContext: (&net.Dialer{
						Timeout:   5 * time.Second,
						KeepAlive: 30 * time.Second,
					}).DialContext,
					TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
					ForceAttemptHTTP2:     true,
					MaxIdleConns:          100,
					IdleConnTimeout:       90 * time.Second,
					TLSHandshakeTimeout:   10 * time.Second,
					ExpectContinueTimeout: 10 * time.Second,
				},
			},
			Header: headers,
		},
	}, nil
}

func (t *Client) Connect(ctx context.Context) (util.LocalConn, error) {
	remote, resp, err := t.h2Client.Connect(ctx, t.tunnel)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		remote.Close()
		return nil, fmt.Errorf("h2: bad status code: %d", resp.StatusCode)
	}

	return &h2LocalConn{
		Conn: remote,
	}, nil
}

type h2LocalConn struct {
	*h2conn.Conn
}

func (conn *h2LocalConn) CloseWrite() error {
	return conn.Close()
}
