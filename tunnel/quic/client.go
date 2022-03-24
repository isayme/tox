package quic

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/lucas-clemente/quic-go/http3"
	"github.com/polevpn/h3conn"
)

var h3client = h3conn.Client{
	RoundTripper: &http3.RoundTripper{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
}

type Client struct {
	tunnel string
}

func NewClient(tunnel string) (*Client, error) {
	URL, err := url.Parse(tunnel)
	if err != nil {
		return nil, err
	}
	switch URL.Scheme {
	case "http3", "quic":
		URL.Scheme = "https"
	}

	return &Client{
		tunnel: URL.String(),
	}, nil
}

func (t *Client) Connect(ctx context.Context) (io.ReadWriteCloser, error) {
	remote, resp, err := h3client.Connect(t.tunnel, time.Second*5, http.Header{})
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		remote.Close()
		return nil, fmt.Errorf("h3: bad status code: %d", resp.StatusCode)
	}

	return remote, nil
}
