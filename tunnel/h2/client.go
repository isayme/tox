package h2

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/isayme/go-logger"
	"github.com/isayme/tox/util"
)

type Client struct {
	tunnel string
}

func NewClient(tunnel string) (*Client, error) {
	URL, err := url.Parse(tunnel)
	if err != nil {
		return nil, err
	}
	switch URL.Scheme {
	case "h2", "http2", "https":
		URL.Scheme = "https"
	default:
		logger.Panicw("URL.Scheme invalid", "address", tunnel, "sceham", URL.Scheme)
	}

	return &Client{
		tunnel: URL.String(),
	}, nil
}

func (t *Client) Connect(ctx context.Context) (io.ReadWriteCloser, error) {
	remote, resp, err := util.H2Client.Connect(ctx, t.tunnel)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		remote.Close()
		return nil, fmt.Errorf("h2: bad status code: %d", resp.StatusCode)
	}

	return remote, nil
}
