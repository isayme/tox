package middleware

import "io"

type readWriter struct {
	io.Reader
	io.Writer
}

func newReadWriter(r io.Reader, w io.Writer) readWriter {
	return readWriter{
		Reader: r,
		Writer: w,
	}
}
