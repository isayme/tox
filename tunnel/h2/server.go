package h2

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/isayme/go-logger"
	"github.com/posener/h2conn"
)

type Server struct {
	tunnel string
}

func NewServer(tunnel string) (*Server, error) {
	return &Server{
		tunnel: tunnel,
	}, nil
}

func (s *Server) ListenAndServeTLS(certFile, keyFile string, handler func(io.ReadWriter)) error {
	URL, err := url.Parse(s.tunnel)
	if err != nil {
		return err
	}

	http.HandleFunc(URL.Path, func(rw http.ResponseWriter, req *http.Request) {
		conn, err := h2conn.Accept(rw, req)
		if err != nil {
			logger.Infof("failed creating connection from %s: %s", req.RemoteAddr, err)
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		handler(conn)
	})

	addr := fmt.Sprintf(":%s", URL.Port())
	return http.ListenAndServeTLS(addr, certFile, keyFile, nil)
}
