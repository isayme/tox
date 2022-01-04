package middleware

import "io"

func NewNoop(rw io.ReadWriter, password string) io.ReadWriter {
	return rw
}
