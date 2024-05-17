package websocket

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/isayme/tox/util"
	"golang.org/x/net/websocket"
)

type Server struct {
	tunnel   string
	password string
	handler  func(util.ServerConn)
}

func NewServer(tunnel string, password string) (*Server, error) {
	return &Server{
		tunnel:   tunnel,
		password: password,
	}, nil
}

func (s *Server) ListenAndServe(handler func(util.ServerConn)) error {
	URL, err := url.Parse(s.tunnel)
	if err != nil {
		return err
	}

	s.handler = handler

	http.Handle(URL.Path, websocket.Server{
		Handshake: s.handshakeWebsocket,
		Handler:   s.handleWebsocket,
	})

	addr := fmt.Sprintf(":%s", URL.Port())

	return http.ListenAndServe(addr, nil)
}

func (s *Server) ListenAndServeTLS(certFile, keyFile string, handler func(util.ServerConn)) error {
	URL, err := url.Parse(s.tunnel)
	if err != nil {
		return err
	}

	s.handler = handler

	http.Handle(URL.Path, websocket.Server{
		Handshake: s.handshakeWebsocket,
		Handler:   s.handleWebsocket,
	})

	addr := fmt.Sprintf(":%s", URL.Port())

	return http.ListenAndServeTLS(addr, certFile, keyFile, nil)
}

func (s *Server) handshakeWebsocket(config *websocket.Config, req *http.Request) error {
	var err error
	config.Origin, err = websocket.Origin(config, req)
	if err == nil && config.Origin == nil {
		return fmt.Errorf("null origin")
	}
	return err
}

func (s *Server) handleWebsocket(ws *websocket.Conn) {
	// defer ws.Close()

	token := ws.Request().Header.Get("token")
	if token != s.password {
		return
	}

	s.handler(&wsLocalConn{Conn: ws})
}
