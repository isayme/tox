package h2

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/isayme/go-logger"
	"github.com/isayme/tox/util"
	"github.com/posener/h2conn"
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
	certFile := s.opts.CertFile
	keyFile := s.opts.KeyFile
	if certFile == "" || keyFile == "" {
		return fmt.Errorf("certFile and keyFile required for http2 protocol")
	}

	URL, err := url.Parse(s.opts.Tunnel)
	if err != nil {
		return err
	}

	http.HandleFunc(URL.Path, func(rw http.ResponseWriter, req *http.Request) {
		if req.Header.Get("token") != s.opts.Password {
			http.Error(rw, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		conn, err := h2conn.Accept(rw, req)
		if err != nil {
			logger.Infof("failed creating connection from %s: %s", req.RemoteAddr, err)
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		// defer conn.Close()

		handler(util.NewToxConnection(conn))
	})

	addr := fmt.Sprintf(":%s", URL.Port())
	return http.ListenAndServeTLS(addr, certFile, keyFile, nil)
}
