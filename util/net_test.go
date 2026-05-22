package util

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// mockNetConn implements net.Conn for testing timeout behavior.
type mockNetConn struct {
	deadline time.Time
}

func (m *mockNetConn) Read(b []byte) (n int, err error)  { return 0, nil }
func (m *mockNetConn) Write(b []byte) (n int, err error)  { return len(b), nil }
func (m *mockNetConn) Close() error                       { return nil }
func (m *mockNetConn) LocalAddr() net.Addr                { return nil }
func (m *mockNetConn) RemoteAddr() net.Addr               { return nil }
func (m *mockNetConn) SetDeadline(t time.Time) error      { m.deadline = t; return nil }
func (m *mockNetConn) SetReadDeadline(t time.Time) error  { m.deadline = t; return nil }
func (m *mockNetConn) SetWriteDeadline(t time.Time) error { m.deadline = t; return nil }

func TestTimeoutConn_Read(t *testing.T) {
	mock := &mockNetConn{}
	timeout := 3 * time.Second
	conn := NewTimeoutConn(mock, timeout)

	before := time.Now()
	conn.Read(nil)
	// deadline should be set to roughly now + timeout
	assert.WithinDuration(t, before.Add(timeout), mock.deadline, 100*time.Millisecond)
}

func TestTimeoutConn_Write(t *testing.T) {
	mock := &mockNetConn{}
	timeout := 3 * time.Second
	conn := NewTimeoutConn(mock, timeout)

	before := time.Now()
	conn.Write(nil)
	assert.WithinDuration(t, before.Add(timeout), mock.deadline, 100*time.Millisecond)
}

func TestTimeoutConn_ZeroTimeout(t *testing.T) {
	mock := &mockNetConn{}
	conn := NewTimeoutConn(mock, 0)

	conn.Read(nil)
	assert.True(t, mock.deadline.IsZero())

	conn.Write(nil)
	assert.True(t, mock.deadline.IsZero())
}

func TestNewTimeoutConn(t *testing.T) {
	mock := &mockNetConn{}
	conn := NewTimeoutConn(mock, 5*time.Second)
	assert.NotNil(t, conn)
	assert.Equal(t, mock, conn.Conn)
}
