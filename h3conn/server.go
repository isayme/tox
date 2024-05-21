package h3conn

import (
	"fmt"
	"io"
	"net/http"
)

// ErrHTTP3NotSupported is returned by Accept if the client connection does not
// support HTTP3 connection.
// The server than can response to the client with an HTTP1.1 as he wishes.
var ErrHTTP3NotSupported = fmt.Errorf("HTTP3 not supported")

// Server can "accept" an http3 connection to obtain a read/write object
// for full duplex communication with a client.
type Server struct {
	StatusCode int
}

// Accept is used on a server http.Handler to extract a full-duplex communication object with the client.
// See h2conn.Accept documentation for more info.
func (u *Server) Accept(w http.ResponseWriter, r *http.Request) (*Conn, error) {
	if !r.ProtoAtLeast(3, 0) {
		return nil, ErrHTTP3NotSupported
	}
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, ErrHTTP3NotSupported
	}

	c, ctx := newConn(r.Context(), r.Body, &flushWrite{w: w, f: flusher})

	// Update the request context with the connection context.
	// If the connection is closed by the server, it will also notify everything that waits on the request context.
	*r = *r.WithContext(ctx)

	w.WriteHeader(u.StatusCode)
	flusher.Flush()

	return c, nil
}

var defaultUpgrader = Server{
	StatusCode: http.StatusOK,
}

// Accept is used on a server http.Handler to extract a full-duplex communication object with the client.
// The server connection will be closed when the http handler function will return.
// If the client does not support HTTP3, an ErrHTTP3NotSupported is returned.
//
// Usage:
//
//	     func (w http.ResponseWriter, r *http.Request) {
//	         conn, err := h2conn.Accept(w, r)
//	         if err != nil {
//			        log.Printf("Failed creating http3 connection: %s", err)
//			        http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
//			        return
//		        }
//	         // use conn
//	     }
func Accept(w http.ResponseWriter, r *http.Request) (*Conn, error) {
	return defaultUpgrader.Accept(w, r)
}

type flushWrite struct {
	w io.Writer
	f http.Flusher
}

func (w *flushWrite) Write(data []byte) (int, error) {
	n, err := w.w.Write(data)
	w.f.Flush()
	return n, err
}

func (w *flushWrite) Close() error {
	// Currently server side close of connection is not supported in Go.
	// The server closes the connection when the http.Handler function returns.
	// We use connection context and cancel function as a work-around.
	return nil
}
