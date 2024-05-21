package quic

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/isayme/tox/h3conn"
	"github.com/isayme/tox/util"
	"github.com/quic-go/quic-go/http3"
)

type Client struct {
	serverAddr string
	opts       util.ToxOptions
	h3Client   *h3conn.Client
}

func NewClient(opts util.ToxOptions) (*Client, error) {
	URL, err := url.Parse(opts.Tunnel)
	if err != nil {
		return nil, err
	}

	switch URL.Scheme {
	case "quic", "http3":
		URL.Scheme = "https"
	}

	headers := http.Header{}
	password := opts.Password
	headers.Add("token", password)

	return &Client{
		serverAddr: URL.String(),
		opts:       opts,
		h3Client: &h3conn.Client{
			Method: http.MethodPost,
			Client: &http.Client{
				Transport: &http3.RoundTripper{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: opts.InsecureSkipVerify,
					},
				},
			},
			Header: headers,
		},
	}, nil
}

func (t *Client) Connect(ctx context.Context) (util.ToxConn, error) {
	remote, resp, err := t.h3Client.Connect(ctx, t.serverAddr)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		remote.Close()
		return nil, fmt.Errorf("h3: bad status code: %d", resp.StatusCode)
	}

	return util.NewToxConnection(newHttp2ClientConnection(remote, resp)), nil
}

type http3ClientConnection struct {
	*h3conn.Conn
	resp *http.Response
}

func newHttp2ClientConnection(conn *h3conn.Conn, resp *http.Response) io.ReadWriteCloser {
	return &http3ClientConnection{
		Conn: conn,
		resp: resp,
	}
}

func (conn *http3ClientConnection) Close() error {
	conn.Conn.Close()
	return conn.resp.Body.Close()
}
