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
	opts     util.ToxOptions
	h2Client *h2conn.Client
}

func NewClient(opts util.ToxOptions) (*Client, error) {
	URL, err := url.Parse(opts.Tunnel)
	if err != nil {
		return nil, err
	}

	switch URL.Scheme {
	case "h2", "http2", "https":
		URL.Scheme = "https"
	}

	headers := http.Header{}
	password := opts.Password
	headers.Add("token", password)

	return &Client{
		tunnel: URL.String(),
		h2Client: &h2conn.Client{
			Method: http.MethodPost,
			Client: &http.Client{
				Transport: &http.Transport{
					DialContext: (&net.Dialer{
						Timeout:   opts.ConnectTimeout,
						KeepAlive: 30 * time.Second,
					}).DialContext,
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: opts.InsecureSkipVerify,
					},
					ForceAttemptHTTP2:     true,
					MaxIdleConns:          100,
					IdleConnTimeout:       opts.Timeout,
					TLSHandshakeTimeout:   opts.ConnectTimeout,
					ExpectContinueTimeout: opts.ConnectTimeout,
				},
			},
			Header: headers,
		},
	}, nil
}

func (t *Client) Connect(ctx context.Context) (util.ToxConn, error) {
	remote, resp, err := t.h2Client.Connect(ctx, t.tunnel)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		remote.Close()
		return nil, fmt.Errorf("h2: bad status code: %d", resp.StatusCode)
	}

	return util.NewToxConnection(remote), nil
}
