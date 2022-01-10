package stream

import (
	"crypto/cipher"
	"io"

	"github.com/isayme/go-bufferpool"
	"github.com/isayme/tox/util"
)

type streamReader struct {
	password string

	keySize int
	ivSize  int

	reader io.Reader

	newCipher func([]byte, []byte) (cipher.Stream, error)

	stream cipher.Stream
}

func NewReader(reader io.Reader, password string, keySize, ivSize int, newCipher func([]byte, []byte) (cipher.Stream, error)) *streamReader {
	return &streamReader{
		password:  password,
		keySize:   keySize,
		ivSize:    ivSize,
		reader:    reader,
		newCipher: newCipher,
	}
}

func (r *streamReader) Read(p []byte) (n int, err error) {
	if r.stream == nil {
		iv := bufferpool.Get(r.ivSize)
		defer bufferpool.Put(iv)
		if _, err = io.ReadFull(r.reader, iv); err != nil {
			return 0, err
		}

		key := util.KDF(r.password, iv, r.keySize)
		s, err := r.newCipher(key, iv)
		if err != nil {
			return 0, err
		}

		r.stream = s
	}

	n, err = r.reader.Read(p)
	r.stream.XORKeyStream(p, p[0:n])
	return n, err
}
