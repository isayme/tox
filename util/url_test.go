package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "ws with default path and port",
			input:    "ws://example.com",
			expected: "ws://example.com:80/",
		},
		{
			name:     "ws with path",
			input:    "ws://example.com/foo",
			expected: "ws://example.com:80/foo",
		},
		{
			name:     "grpc with default path and port",
			input:    "grpc://example.com",
			expected: "grpc://example.com:80/",
		},
		{
			name:     "wss with default path and port",
			input:    "wss://example.com",
			expected: "wss://example.com:443/",
		},
		{
			name:     "grpcs with default path and port",
			input:    "grpcs://example.com",
			expected: "grpcs://example.com:443/",
		},
		{
			name:     "http2 with default path and port",
			input:    "http2://example.com",
			expected: "http2://example.com:443/",
		},
		{
			name:     "already has port",
			input:    "ws://example.com:8080",
			expected: "ws://example.com:8080/",
		},
		{
			name:     "already has path and port",
			input:    "wss://example.com:8443/foo",
			expected: "wss://example.com:8443/foo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FormatURL(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatURL_Invalid(t *testing.T) {
	_, err := FormatURL("://bad")
	assert.Error(t, err)
}
