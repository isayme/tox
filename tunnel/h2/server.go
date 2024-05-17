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
	tunnel   string
	password string
}

func NewServer(tunnel string, password string) (*Server, error) {
	return &Server{
		tunnel:   tunnel,
		password: password,
	}, nil
}

func (s *Server) ListenAndServe(handler func(util.ServerConn)) error {
	return fmt.Errorf("tls required for http2 protocol")
}

func (s *Server) ListenAndServeTLS(certFile, keyFile string, handler func(util.ServerConn)) error {
	URL, err := url.Parse(s.tunnel)
	if err != nil {
		return err
	}

	http.HandleFunc(URL.Path, func(rw http.ResponseWriter, req *http.Request) {
		if req.Header.Get("token") != s.password {
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

		handler(&h2LocalConn{Conn: conn})
	})

	addr := fmt.Sprintf(":%s", URL.Port())
	return http.ListenAndServeTLS(addr, certFile, keyFile, nil)
}
