package middleware

import (
	"io"

	"github.com/isayme/go-toh2/aead"
)

func NewAes128Gcm(rw io.ReadWriter, password string) io.ReadWriter {
	r := aead.NewAeadReader(rw, password, 16, aead.NewAesGcmCipher)
	w := aead.NewAeadWriter(rw, password, 16, aead.NewAesGcmCipher)
	return newReadWriter(r, w)
}

func NewAes256Gcm(rw io.ReadWriter, password string) io.ReadWriter {
	r := aead.NewAeadReader(rw, password, 32, aead.NewAesGcmCipher)
	w := aead.NewAeadWriter(rw, password, 32, aead.NewAesGcmCipher)
	return newReadWriter(r, w)
}
