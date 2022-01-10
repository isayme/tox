package middleware

import (
	"io"

	"github.com/isayme/tox/stream"
)

func NewAes256Cfb(rw io.ReadWriter, password string) io.ReadWriter {
	r := stream.NewReader(rw, password, 32, 16, stream.NewAesCfbReader)
	w := stream.NewWriter(rw, password, 32, 16, stream.NewAesCfbWriter)
	return newReadWriter(r, w)
}
