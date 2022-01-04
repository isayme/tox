package util

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/posener/h2conn"
)

var defaultTransport http.RoundTripper = &http.Transport{
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
}

var H2Client = &h2conn.Client{
	Method: http.MethodPost,
	Client: &http.Client{Transport: defaultTransport},
}
