package util

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"sync/atomic"

	"github.com/isayme/go-bufferpool"
	"github.com/vmihailenco/msgpack/v5"
)

const (
	COMMAND_DATA        = 0x01
	COMMAND_CLOSE_WRITE = 0x02

	PACK_DATA_BUF_SIZE = 4
)

var errBrokenPipe = fmt.Errorf("broken pipe")
var closeWriteData = []byte{COMMAND_CLOSE_WRITE}

type Frame struct {
	Cmd  uint8
	Data []byte
}

type ToxConn interface {
	io.Reader
	io.Writer
	io.Closer
	CloseWrite() error
}

type ToxConnection struct {
	conn       io.ReadWriteCloser
	buffer     *bytes.Buffer
	closeWrite atomic.Bool
}

func NewToxConnection(conn io.ReadWriteCloser) ToxConn {
	return &ToxConnection{
		conn:   conn,
		buffer: nil,
	}
}

func (conn *ToxConnection) readFrame() (*Frame, error) {
	var len int

	// read data length
	{
		buf := bufferpool.Get(PACK_DATA_BUF_SIZE)
		_, err := io.ReadFull(conn.conn, buf)
		if err != nil {
			bufferpool.Put(buf)
			return nil, err
		}

		len = int(binary.BigEndian.Uint32(buf))
		bufferpool.Put(buf)
	}

	// read data
	buf := bufferpool.Get(len)

	_, err := io.ReadFull(conn.conn, buf)
	if err != nil {
		bufferpool.Put(buf)
		return nil, err
	}

	var frame Frame
	err = msgpack.Unmarshal(buf, &frame)
	bufferpool.Put(buf)
	if err != nil {
		return nil, err
	}

	return &frame, nil
}

func (conn *ToxConnection) Read(p []byte) (int, error) {
	if conn.buffer != nil && conn.buffer.Len() > 0 {
		return conn.buffer.Read(p)
	}

	if conn.buffer != nil {
		conn.buffer.Reset()
		conn.buffer = nil
	}

	frame, err := conn.readFrame()
	if err != nil {
		return 0, err
	}

	switch frame.Cmd {
	case COMMAND_DATA:
		conn.buffer = bytes.NewBuffer(frame.Data)
		return conn.buffer.Read(p)
	case COMMAND_CLOSE_WRITE:
		fallthrough
	default:
		return 0, io.EOF
	}
}

func (conn *ToxConnection) writeFrame(cmd uint8, p []byte) error {
	frame := Frame{
		Cmd:  cmd,
		Data: p,
	}
	data, err := msgpack.Marshal(frame)
	if err != nil {
		return err
	}

	// write data length
	{
		buf := bufferpool.Get(PACK_DATA_BUF_SIZE)
		binary.BigEndian.PutUint32(buf, uint32(len(data)))
		n, err := conn.conn.Write(buf)
		bufferpool.Put(buf)
		if err != nil {
			return err
		}
		if n != PACK_DATA_BUF_SIZE {
			return io.ErrShortWrite
		}
	}

	_, err = conn.conn.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (conn *ToxConnection) Write(p []byte) (int, error) {
	if conn.closeWrite.Load() {
		return 0, errBrokenPipe
	}

	err := conn.writeFrame(COMMAND_DATA, p)
	if err != nil {
		return 0, err
	}

	return len(p), nil
}

func (conn *ToxConnection) CloseWrite() error {
	if !conn.closeWrite.CompareAndSwap(false, true) {
		return nil
	}

	err := conn.writeFrame(COMMAND_CLOSE_WRITE, nil)
	if err != nil {
		return err
	}
	return nil
}

func (conn *ToxConnection) Close() error {
	return conn.conn.Close()
}
