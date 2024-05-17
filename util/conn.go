package util

import "io"

type LocalConn interface {
	io.Reader
	io.Writer
	io.Closer
	CloseWrite() error
}

type ServerConn interface {
	io.Reader
	io.Writer
	CloseWrite() error
}
