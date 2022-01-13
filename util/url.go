package util

import (
	"net"
	"net/url"
)

func FormatURL(s string) (string, error) {
	URL, err := url.Parse(s)
	if err != nil {
		return "", err
	}

	if URL.Path == "" {
		URL.Path = "/"
	}

	if URL.Port() == "" {
		var port string
		switch URL.Scheme {
		case "ws", "grpc":
			port = "80"
		default:
			port = "443"
		}
		URL.Host = net.JoinHostPort(URL.Host, port)
	}

	return URL.String(), nil
}
