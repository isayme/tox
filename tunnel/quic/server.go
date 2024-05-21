package quic

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/isayme/go-logger"
	"github.com/isayme/tox/h3conn"
	"github.com/isayme/tox/util"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

type Server struct {
	opts util.ToxOptions
}

func NewServer(opts util.ToxOptions) (*Server, error) {
	return &Server{
		opts: opts,
	}, nil
}

func (s *Server) ListenAndServe(handler func(util.ToxConn)) error {
	URL, err := url.Parse(s.opts.Tunnel)
	if err != nil {
		return err
	}

	certFile := s.opts.CertFile
	keyFile := s.opts.KeyFile
	if certFile == "" || keyFile == "" {
		return fmt.Errorf("certFile and keyFile required for http3(quic) protocol")
	}

	addr := fmt.Sprintf(":%s", URL.Port())
	server := &http3.Server{
		Addr: addr,
		QUICConfig: &quic.Config{
			HandshakeIdleTimeout: s.opts.ConnectTimeout,
			MaxIdleTimeout:       s.opts.Timeout,
		},
		StreamHijacker: nil,
	}

	http.HandleFunc(URL.Path, func(rw http.ResponseWriter, req *http.Request) {
		if req.Header.Get("token") != s.opts.Password {
			http.Error(rw, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		conn, err := h3conn.Accept(rw, req)
		if err != nil {
			logger.Infof("failed creating connection from %s: %s", req.RemoteAddr, err)
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		handler(util.NewToxConnection(conn))
	})

	return server.ListenAndServeTLS(certFile, keyFile)
}

func handleConnection(conn quic.Connection, handler func(util.ToxConn)) {
	defer conn.CloseWithError(0, "")

	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			logger.Infof("AcceptStream fail, err: %v", err)
			break
		}

		go handleStream(stream, handler)
	}
}

func handleStream(stream quic.Stream, handler func(util.ToxConn)) {
	handler(util.NewToxConnection(stream))
}
