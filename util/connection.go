package util

import "io"

type Connection struct {
	io.Reader
	io.Writer
}

func NewConnection(r io.Reader, w io.Writer) Connection {
	return Connection{
		Reader: r,
		Writer: w,
	}
}
