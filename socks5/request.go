package socks5

import (
	"encoding/binary"
	"io"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/isayme/go-bufferpool"
	"github.com/isayme/go-logger"
	"github.com/isayme/tox/conf"
	"github.com/isayme/tox/util"
	"github.com/pkg/errors"
)

type Request struct {
	rw util.ServerConn

	cmd  byte
	atyp byte
	addr string
}

func NewRequest(conn util.ServerConn) *Request {
	return &Request{
		rw: conn,
	}
}

func (r *Request) Handle() error {
	if err := r.negotiate(); err != nil {
		return err
	}
	return r.handleRequest()
}

func (r *Request) negotiate() error {
	buf := bufferpool.Get(256)
	defer bufferpool.Put(buf)

	// version
	_, err := io.ReadFull(r.rw, buf[:1])
	if err != nil {
		return errors.Errorf("read version fail: %s", err)
	}

	if buf[0] != Version {
		return errors.New("not socks5 protocol")
	}

	// methods
	_, err = io.ReadFull(r.rw, buf[:1])
	if err != nil {
		return errors.Errorf("read nmethods fail: %s", err)
	}
	nMethods := buf[0]
	if nMethods < 1 {
		return errors.Errorf("nmethods not valid: %d", nMethods)
	}

	_, err = io.ReadFull(r.rw, buf[:nMethods])
	if err != nil {
		return errors.Errorf("read nmethods fail: %s", err)
	}

	_, err = r.rw.Write([]byte{Version, MethodNone})
	if err != nil {
		return errors.Errorf("write accepet method fail: %s", err)
	}

	_, err = io.ReadFull(r.rw, buf[:4])
	if err != nil {
		return errors.Errorf("read adrress fail: %s", err)
	}
	r.cmd = buf[1]
	r.atyp = buf[3]

	var reply = []byte{Version, 0, 0, r.atyp}

	switch r.cmd {
	case CmdConnect:
	default:
		return errors.Errorf("not support cmd: %d", r.cmd)
	}

	switch r.atyp {
	case AddressTypeDomain:
		_, err = io.ReadFull(r.rw, buf[:1])
		if err != nil {
			return errors.Errorf("read adrress fail: %s", err)
		}
		domainLen := buf[0]
		reply = append(reply, buf[0])

		_, err = io.ReadFull(r.rw, buf[:domainLen])
		if err != nil {
			return errors.Errorf("read domain fail: %s", err)
		}
		reply = append(reply, buf[:domainLen]...)

		domain := string(buf[:domainLen])
		r.addr = domain
	case AddressTypeIPV4:
		_, err = io.ReadFull(r.rw, buf[:net.IPv4len])
		if err != nil {
			return errors.Errorf("read ipv4 fail: %s", err)
		}

		reply = append(reply, buf[:net.IPv4len]...)

		ip := net.IP(buf[:net.IPv4len])
		r.addr = ip.String()
	case AddressTypeIPV6:
		_, err = io.ReadFull(r.rw, buf[:net.IPv6len])
		if err != nil {
			return errors.Errorf("read ipv6 fail: %s", err)
		}

		reply = append(reply, buf[:net.IPv6len]...)

		ip := net.IP(buf[:net.IPv6len])
		r.addr = ip.String()
	default:
		return errors.Errorf("not support adrress type: %d", r.atyp)
	}

	_, err = io.ReadFull(r.rw, buf[:2])
	if err != nil {
		return errors.Errorf("read port fail: %s", err)
	}
	reply = append(reply, buf[:2]...)
	port := binary.BigEndian.Uint16(buf[:2])

	_, err = r.rw.Write(reply)
	if err != nil {
		return errors.Errorf("reply request fail: %s", err)
	}

	r.addr = net.JoinHostPort(r.addr, strconv.Itoa(int(port)))

	logger.Infow("new socks5 request", "cmd", r.cmd, "atyp", r.atyp, "addr", r.addr)
	return nil
}

func (r *Request) handleRequest() error {
	remote, err := net.DialTimeout("tcp", r.addr, time.Second*5)
	if err != nil {
		logger.Infow("net.Dial fail", "err", err, "addr", r.addr)
		return err
	}
	defer remote.Close()

	config := conf.Get()

	remoteTcpConn, _ := remote.(*net.TCPConn)
	remote = util.NewTimeoutConn(remote, time.Duration(config.Timeout)*time.Second)

	logger.Infow("connect ok", "addr", r.addr)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		var err error
		var n int64
		n, err = util.CopyBuffer(r.rw, remote)
		r.rw.CloseWrite()
		logger.Debugw("copy from remote end", "n", n, "err", err)
	}()

	go func() {
		defer wg.Done()

		var err error
		var n int64
		n, err = util.CopyBuffer(remote, r.rw)
		logger.Debugw("copy from client end", "n", n, "err", err)
		remoteTcpConn.CloseWrite()
	}()

	wg.Wait()

	logger.Infow("handle end", "addr", r.addr)

	return nil
}
