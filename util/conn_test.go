package util

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmihailenco/msgpack/v5"
)

// mockReadWriteCloser simulates the underlying connection with the frame protocol.
type mockReadWriteCloser struct {
	readBuf  *bytes.Buffer
	writeBuf *bytes.Buffer
	closed   bool
}

func newMockConn() *mockReadWriteCloser {
	return &mockReadWriteCloser{
		readBuf:  bytes.NewBuffer(nil),
		writeBuf: bytes.NewBuffer(nil),
	}
}

func (m *mockReadWriteCloser) Read(p []byte) (int, error) {
	return m.readBuf.Read(p)
}

func (m *mockReadWriteCloser) Write(p []byte) (int, error) {
	return m.writeBuf.Write(p)
}

func (m *mockReadWriteCloser) Close() error {
	m.closed = true
	return nil
}

func (m *mockReadWriteCloser) writeFrameToReadBuf(cmd uint8, data []byte) {
	frame := Frame{Cmd: cmd, Data: data}
	encoded, _ := msgpack.Marshal(frame)
	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, uint32(len(encoded)))
	m.readBuf.Write(lenBuf)
	m.readBuf.Write(encoded)
}

func (m *mockReadWriteCloser) readFrameFromWriteBuf() (*Frame, error) {
	lenBuf := make([]byte, 4)
	_, err := io.ReadFull(m.writeBuf, lenBuf)
	if err != nil {
		return nil, err
	}
	dataLen := binary.BigEndian.Uint32(lenBuf)
	data := make([]byte, dataLen)
	_, err = io.ReadFull(m.writeBuf, data)
	if err != nil {
		return nil, err
	}
	var frame Frame
	err = msgpack.Unmarshal(data, &frame)
	if err != nil {
		return nil, err
	}
	return &frame, nil
}

func TestToxConnection_ReadData(t *testing.T) {
	mock := newMockConn()
	conn := NewToxConnection(mock)

	mock.writeFrameToReadBuf(COMMAND_DATA, []byte("hello"))
	buf := make([]byte, 128)
	n, err := conn.Read(buf)
	require.NoError(t, err)
	assert.Equal(t, "hello", string(buf[:n]))
}

func TestToxConnection_ReadCloseWrite(t *testing.T) {
	mock := newMockConn()
	conn := NewToxConnection(mock)

	mock.writeFrameToReadBuf(COMMAND_CLOSE_WRITE, nil)
	buf := make([]byte, 128)
	n, err := conn.Read(buf)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, 0, n)
}

func TestToxConnection_ReadBuffered(t *testing.T) {
	mock := newMockConn()
	conn := NewToxConnection(mock)

	// write a long frame, then read in small chunks
	mock.writeFrameToReadBuf(COMMAND_DATA, []byte("hello world"))
	buf := make([]byte, 5)
	n, err := conn.Read(buf)
	require.NoError(t, err)
	assert.Equal(t, "hello", string(buf[:n]))

	// remaining bytes from buffer
	n, err = conn.Read(buf)
	require.NoError(t, err)
	assert.Equal(t, " worl", string(buf[:n]))
}

func TestToxConnection_Write(t *testing.T) {
	mock := newMockConn()
	conn := NewToxConnection(mock)

	n, err := conn.Write([]byte("test-data"))
	require.NoError(t, err)
	assert.Equal(t, 9, n)

	frame, err := mock.readFrameFromWriteBuf()
	require.NoError(t, err)
	assert.Equal(t, uint8(COMMAND_DATA), frame.Cmd)
	assert.Equal(t, []byte("test-data"), frame.Data)
}

func TestToxConnection_CloseWrite(t *testing.T) {
	mock := newMockConn()
	conn := NewToxConnection(mock)

	err := conn.CloseWrite()
	require.NoError(t, err)

	frame, err := mock.readFrameFromWriteBuf()
	require.NoError(t, err)
	assert.Equal(t, uint8(COMMAND_CLOSE_WRITE), frame.Cmd)
}

func TestToxConnection_CloseWrite_Idempotent(t *testing.T) {
	mock := newMockConn()
	conn := NewToxConnection(mock)

	err := conn.CloseWrite()
	require.NoError(t, err)

	err = conn.CloseWrite()
	require.NoError(t, err)

	// only one frame written
	frame, err := mock.readFrameFromWriteBuf()
	require.NoError(t, err)
	assert.Equal(t, uint8(COMMAND_CLOSE_WRITE), frame.Cmd)

	_, err = mock.readFrameFromWriteBuf()
	assert.Error(t, err) // second frame should not exist
}

func TestToxConnection_WriteAfterCloseWrite(t *testing.T) {
	mock := newMockConn()
	conn := NewToxConnection(mock)

	err := conn.CloseWrite()
	require.NoError(t, err)

	_, err = conn.Write([]byte("data"))
	assert.Error(t, err)
}

func TestToxConnection_Close(t *testing.T) {
	mock := newMockConn()
	conn := NewToxConnection(mock)

	assert.False(t, mock.closed)
	err := conn.Close()
	require.NoError(t, err)
	assert.True(t, mock.closed)
}
