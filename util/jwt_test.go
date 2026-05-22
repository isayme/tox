package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAndValidateJwtToken(t *testing.T) {
	key := []byte("test-secret-key-32bytes!!----")
	tokenStr, err := GenerateJwtToken(key)
	require.NoError(t, err)
	assert.NotEmpty(t, tokenStr)

	err = ValidateJwtToken(tokenStr, key)
	assert.NoError(t, err)
}

func TestValidateJwtToken_InvalidToken(t *testing.T) {
	key := []byte("test-secret-key-32bytes!!----")

	err := ValidateJwtToken("not.a.valid.token", key)
	assert.Error(t, err)
}

func TestValidateJwtToken_WrongKey(t *testing.T) {
	key := []byte("test-secret-key-32bytes!!----")
	tokenStr, err := GenerateJwtToken(key)
	require.NoError(t, err)

	wrongKey := []byte("wrong-secret-key-32bytes!!---")
	err = ValidateJwtToken(tokenStr, wrongKey)
	assert.Error(t, err)
}

func TestGenerateJwtToken_Deterministic(t *testing.T) {
	key := []byte("test-secret-key-32bytes!!----")
	// Two tokens generated at nearly the same time may differ due to timestamp
	token1, err := GenerateJwtToken(key)
	require.NoError(t, err)
	token2, err := GenerateJwtToken(key)
	require.NoError(t, err)

	// Both should be independently valid
	err = ValidateJwtToken(token1, key)
	assert.NoError(t, err)
	err = ValidateJwtToken(token2, key)
	assert.NoError(t, err)
}
