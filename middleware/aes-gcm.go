package middleware

import (
	"io"

	"github.com/isayme/tox/aead"
)

func NewAes128Gcm(rw io.ReadWriter, password string) io.ReadWriter {
	r := aead.NewReader(rw, password, 16, aead.NewAesGcmCipher)
	w := aead.NewWriter(rw, password, 16, aead.NewAesGcmCipher)
	return newReadWriter(r, w)
}

func NewAes256Gcm(rw io.ReadWriter, password string) io.ReadWriter {
	r := aead.NewReader(rw, password, 32, aead.NewAesGcmCipher)
	w := aead.NewWriter(rw, password, 32, aead.NewAesGcmCipher)
	return newReadWriter(r, w)
}
