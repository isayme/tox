package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKDF(t *testing.T) {
	key := KDF("hello", []byte("salt"), 32)
	assert.Len(t, key, 32)

	// deterministic
	key2 := KDF("hello", []byte("salt"), 32)
	assert.Equal(t, key, key2)

	// different password produces different key
	key3 := KDF("world", []byte("salt"), 32)
	assert.NotEqual(t, key, key3)

	// different salt produces different key
	key4 := KDF("hello", []byte("pepper"), 32)
	assert.NotEqual(t, key, key4)

	// different key size
	key5 := KDF("hello", []byte("salt"), 64)
	assert.Len(t, key5, 64)
}

func TestHashedPassword(t *testing.T) {
	h1 := HashedPassword("mypassword")
	assert.NotEmpty(t, h1)

	// deterministic
	h2 := HashedPassword("mypassword")
	assert.Equal(t, h1, h2)

	// different password
	h3 := HashedPassword("other")
	assert.NotEqual(t, h1, h3)
}
