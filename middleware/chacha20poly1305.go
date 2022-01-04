package middleware

import (
	"io"

	"github.com/isayme/go-toh2/aead"
)

func NewChacha20Poly1305(rw io.ReadWriter, password string) io.ReadWriter {
	r := aead.NewAeadReader(rw, password, 32, aead.NewChacha20Poly1305Cipher)
	w := aead.NewAeadWriter(rw, password, 32, aead.NewChacha20Poly1305Cipher)
	return newReadWriter(r, w)
}

func NewXChacha20Poly1305(rw io.ReadWriter, password string) io.ReadWriter {
	r := aead.NewAeadReader(rw, password, 32, aead.NewChacha20Poly1305Cipher)
	w := aead.NewAeadWriter(rw, password, 32, aead.NewChacha20Poly1305Cipher)
	return newReadWriter(r, w)
}
