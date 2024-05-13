package util

import (
	"crypto/sha256"
	"encoding/base64"

	"golang.org/x/crypto/pbkdf2"
)

const MagicKeySalt = "tox"

func KDF(password string, salt []byte, keySize int) []byte {
	return pbkdf2.Key([]byte(password), salt, 1024, keySize, sha256.New)
}

func HashedPassword(password string) string {
	key := KDF(password, []byte(MagicKeySalt), 32)
	return base64.StdEncoding.EncodeToString(key)
}
