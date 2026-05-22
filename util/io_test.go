package util

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopyBuffer(t *testing.T) {
	src := bytes.NewReader([]byte("hello world"))
	var dst bytes.Buffer

	n, err := CopyBuffer(&dst, src)
	require.NoError(t, err)
	assert.Equal(t, int64(11), n)
	assert.Equal(t, "hello world", dst.String())
}

func TestCopyBuffer_Empty(t *testing.T) {
	src := bytes.NewReader(nil)
	var dst bytes.Buffer

	n, err := CopyBuffer(&dst, src)
	require.NoError(t, err)
	assert.Equal(t, int64(0), n)
	assert.Empty(t, dst.String())
}

func TestCopyBuffer_LargeData(t *testing.T) {
	data := make([]byte, 1024*1024) // 1MB
	for i := range data {
		data[i] = byte(i % 256)
	}
	src := bytes.NewReader(data)
	var dst bytes.Buffer

	n, err := CopyBuffer(&dst, src)
	require.NoError(t, err)
	assert.Equal(t, int64(len(data)), n)
	assert.Equal(t, data, dst.Bytes())
}

func TestCopyBuffer_WriteError(t *testing.T) {
	src := bytes.NewReader([]byte("data"))
	dst := &failingWriter{}

	_, err := CopyBuffer(dst, src)
	assert.Error(t, err)
}

type failingWriter struct{}

func (w *failingWriter) Write(p []byte) (n int, err error) {
	return 0, io.ErrShortWrite
}
