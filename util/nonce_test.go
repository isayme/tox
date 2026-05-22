package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNextNonce(t *testing.T) {
	// basic increment
	nonce := []byte{0x00, 0x00, 0x00, 0x00}
	NextNonce(nonce)
	assert.Equal(t, []byte{0x01, 0x00, 0x00, 0x00}, nonce)

	NextNonce(nonce)
	assert.Equal(t, []byte{0x02, 0x00, 0x00, 0x00}, nonce)
}

func TestNextNonce_Carry(t *testing.T) {
	nonce := []byte{0xff, 0x00, 0x00, 0x00}
	NextNonce(nonce)
	assert.Equal(t, []byte{0x00, 0x01, 0x00, 0x00}, nonce)
}

func TestNextNonce_MultiCarry(t *testing.T) {
	nonce := []byte{0xff, 0xff, 0x00, 0x00}
	NextNonce(nonce)
	assert.Equal(t, []byte{0x00, 0x00, 0x01, 0x00}, nonce)
}

func TestNextNonce_Overflow(t *testing.T) {
	nonce := []byte{0xff, 0xff, 0xff, 0xff}
	NextNonce(nonce)
	assert.Equal(t, []byte{0x00, 0x00, 0x00, 0x00}, nonce)
}

func TestNextNonce_EmptyNonce(t *testing.T) {
	nonce := []byte{}
	NextNonce(nonce)
	assert.Empty(t, nonce)
}
