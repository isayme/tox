package util

import (
	"net"
	"time"
)

// timeoutConn net.Conn with Read/Write timeout. from https://qiita.com/kwi/items/b38d6273624ad3f6ae79
type timeoutConn struct {
	net.Conn
	timeout time.Duration
}

// NewTimeoutConn create timeout conn
func NewTimeoutConn(conn net.Conn, timeout time.Duration) *timeoutConn {
	return &timeoutConn{
		Conn:    conn,
		timeout: timeout,
	}
}

func (c *timeoutConn) Read(p []byte) (n int, err error) {
	if c.timeout > 0 {
		if err := c.Conn.SetDeadline(time.Now().Add(c.timeout)); err != nil {
			return 0, err
		}
	}

	return c.Conn.Read(p)
}

func (c *timeoutConn) Write(p []byte) (n int, err error) {
	if c.timeout > 0 {
		if err := c.Conn.SetDeadline(time.Now().Add(c.timeout)); err != nil {
			return 0, err
		}
	}

	return c.Conn.Write(p)
}
