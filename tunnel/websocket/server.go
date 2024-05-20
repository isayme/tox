package websocket

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/isayme/tox/util"
	"golang.org/x/net/websocket"
)

type Server struct {
	opts    util.ToxOptions
	handler func(util.ToxConn)
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

	s.handler = handler

	http.Handle(URL.Path, websocket.Server{
		Handshake: s.handshakeWebsocket,
		Handler:   s.handleWebsocket,
	})

	addr := fmt.Sprintf(":%s", URL.Port())

	certKey := s.opts.CertFile
	keyKey := s.opts.KeyFile
	if certKey != "" && keyKey != "" {
		return http.ListenAndServeTLS(addr, certKey, keyKey, nil)
	} else {
		return http.ListenAndServe(addr, nil)
	}
}

func (s *Server) handshakeWebsocket(config *websocket.Config, req *http.Request) error {
	var err error
	config.Origin, err = websocket.Origin(config, req)
	if err == nil && config.Origin == nil {
		return fmt.Errorf("null origin")
	}

	token := req.Header.Get("token")
	if token != s.opts.Password {
		return fmt.Errorf("invalid password")
	}

	return err
}

func (s *Server) handleWebsocket(ws *websocket.Conn) {
	s.handler(util.NewToxConnection(ws))
}
