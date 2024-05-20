package util

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/isayme/go-bufferpool"
)

const (
	COMMAND_DATA        = 0x01
	COMMAND_CLOSE_WRITE = 0x02
)

var errBrokenPipe = fmt.Errorf("broken pipe")
var closeWriteData = []byte{COMMAND_CLOSE_WRITE}

type ToxConn interface {
	io.Reader
	io.Writer
	io.Closer
	CloseWrite() error
}

type ToxConnection struct {
	conn       io.ReadWriteCloser
	buffer     *bytes.Buffer
	closeWrite bool
}

func NewToxConnection(conn io.ReadWriteCloser) ToxConn {
	return &ToxConnection{
		conn:   conn,
		buffer: nil,
	}
}

func (conn *ToxConnection) Read(p []byte) (int, error) {
	if conn.buffer != nil && conn.buffer.Len() > 0 {
		return conn.buffer.Read(p)
	}

	if conn.buffer != nil {
		conn.buffer.Reset()
		bufferpool.Put(conn.buffer.Bytes())
		conn.buffer = nil
	}

	var cmd byte
	var len int

	// read data cmd
	{
		buf := bufferpool.Get(1)

		_, err := io.ReadFull(conn.conn, buf)
		if err != nil {
			bufferpool.Put(buf)
			return 0, err
		}
		cmd = buf[0]
		bufferpool.Put(buf)
	}

	if cmd != COMMAND_DATA {
		return 0, io.EOF
	}

	// read data length
	{
		buf := bufferpool.Get(4)
		_, err := io.ReadFull(conn.conn, buf)
		if err != nil {
			bufferpool.Put(buf)
			return 0, err
		}
		len = int(binary.BigEndian.Uint32(buf))
		bufferpool.Put(buf)
	}

	{
		// read data
		buf := bufferpool.Get(len)

		_, err := io.ReadFull(conn.conn, buf)
		if err != nil {
			return 0, err
		}

		conn.buffer = bytes.NewBuffer(buf)
		return conn.buffer.Read(p)
	}
}

func (conn *ToxConnection) Write(p []byte) (int, error) {
	if conn.closeWrite {
		return 0, errBrokenPipe
	}

	// write cmd
	{
		buf := bufferpool.Get(1)
		buf[0] = COMMAND_DATA
		n, err := conn.conn.Write(buf)
		bufferpool.Put(buf)
		if err != nil {
			return 0, err
		}
		if n != 1 {
			return 0, io.ErrShortWrite
		}
	}

	// write data length
	{
		buf := bufferpool.Get(4)
		binary.BigEndian.PutUint32(buf, uint32(len(p)))
		n, err := conn.conn.Write(buf)
		bufferpool.Put(buf)
		if err != nil {
			return 0, err
		}
		if n != 4 {
			return 0, io.ErrShortWrite
		}
	}

	// write data
	return conn.conn.Write(p)
}

func (conn *ToxConnection) CloseWrite() error {
	conn.closeWrite = true
	_, err := conn.conn.Write(closeWriteData)
	return err
}

func (conn *ToxConnection) Close() error {
	return conn.conn.Close()
}
