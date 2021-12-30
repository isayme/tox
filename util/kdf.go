package util

import (
	"crypto/sha256"

	"golang.org/x/crypto/pbkdf2"
)

func KDF(password string, salt []byte, keySize int) []byte {
	return pbkdf2.Key([]byte(password), salt, 1024, keySize, sha256.New)
}
