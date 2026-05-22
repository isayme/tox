package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToToxOptions_Empty(t *testing.T) {
	opts := ToToxOptions(nil)
	assert.Empty(t, opts.Password)
	assert.Empty(t, opts.Tunnel)
	assert.Empty(t, opts.LocalAddress)
}

func TestWithPassword(t *testing.T) {
	opts := ToToxOptions([]ToxOption{WithPassword("secret")})
	assert.NotEmpty(t, opts.Password)
	// password should be hashed, not raw
	assert.NotEqual(t, "secret", opts.Password)
}

func TestWithTunnel(t *testing.T) {
	opts := ToToxOptions([]ToxOption{WithTunnel("grpc://example.com")})
	assert.Equal(t, "grpc://example.com", opts.Tunnel)
}

func TestWithLocalAddress(t *testing.T) {
	opts := ToToxOptions([]ToxOption{WithLocalAddress(":1080")})
	assert.Equal(t, ":1080", opts.LocalAddress)
}

func TestWithCertFile(t *testing.T) {
	opts := ToToxOptions([]ToxOption{WithCertFile("/path/cert.pem")})
	assert.Equal(t, "/path/cert.pem", opts.CertFile)
}

func TestWithKeyFile(t *testing.T) {
	opts := ToToxOptions([]ToxOption{WithKeyFile("/path/key.pem")})
	assert.Equal(t, "/path/key.pem", opts.KeyFile)
}

func TestWithConnectTimeout(t *testing.T) {
	opts := ToToxOptions([]ToxOption{WithConnectTimeout(5 * time.Second)})
	assert.Equal(t, 5*time.Second, opts.ConnectTimeout)
}

func TestWithTimeout(t *testing.T) {
	opts := ToToxOptions([]ToxOption{WithTimeout(10 * time.Second)})
	assert.Equal(t, 10*time.Second, opts.Timeout)
}

func TestWithInsecureSkipVerify(t *testing.T) {
	opts := ToToxOptions([]ToxOption{WithInsecureSkipVerify(true)})
	assert.True(t, opts.InsecureSkipVerify)

	opts2 := ToToxOptions([]ToxOption{WithInsecureSkipVerify(false)})
	assert.False(t, opts2.InsecureSkipVerify)
}

func TestToToxOptions_Combined(t *testing.T) {
	opts := ToToxOptions([]ToxOption{
		WithTunnel("grpcs://example.com"),
		WithLocalAddress(":8080"),
		WithTimeout(3 * time.Second),
		WithInsecureSkipVerify(true),
	})

	assert.Equal(t, "grpcs://example.com", opts.Tunnel)
	assert.Equal(t, ":8080", opts.LocalAddress)
	assert.Equal(t, 3*time.Second, opts.Timeout)
	assert.True(t, opts.InsecureSkipVerify)
}
