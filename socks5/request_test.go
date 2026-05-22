package socks5

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"testing"

	"github.com/isayme/tox/conf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockToxConn implements util.ToxConn for testing SOCKS5 negotiation.
type mockToxConn struct {
	readBuf  *bytes.Buffer
	writeBuf *bytes.Buffer
	closed   bool
}

func newMockToxConn() *mockToxConn {
	return &mockToxConn{
		readBuf:  bytes.NewBuffer(nil),
		writeBuf: bytes.NewBuffer(nil),
	}
}

func (m *mockToxConn) Read(p []byte) (int, error) {
	return m.readBuf.Read(p)
}

func (m *mockToxConn) Write(p []byte) (int, error) {
	return m.writeBuf.Write(p)
}

func (m *mockToxConn) Close() error {
	m.closed = true
	return nil
}

func (m *mockToxConn) CloseWrite() error {
	return nil
}

func (m *mockToxConn) writeBytes(b []byte) {
	m.readBuf.Write(b)
}

// writeSOCKS5Greeting writes a valid SOCKS5 greeting: VER=5, NMETHODS=1, METHOD=0
func (m *mockToxConn) writeSOCKS5Greeting() {
	m.writeBytes([]byte{Version, 1, MethodNone})
}

// writeSOCKS5Request writes a SOCKS5 CONNECT request with the given address type and destination.
func (m *mockToxConn) writeSOCKS5RequestDomain(domain string, port uint16) {
	buf := make([]byte, 0, 5+len(domain)+2)
	buf = append(buf, Version, CmdConnect, 0x00, AddressTypeDomain, byte(len(domain)))
	buf = append(buf, []byte(domain)...)
	portBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(portBytes, port)
	buf = append(buf, portBytes...)
	m.writeBytes(buf)
}

func (m *mockToxConn) writeSOCKS5RequestIPv4(ip net.IP, port uint16) {
	buf := make([]byte, 0, 5+net.IPv4len+2)
	buf = append(buf, Version, CmdConnect, 0x00, AddressTypeIPV4)
	buf = append(buf, ip.To4()...)
	portBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(portBytes, port)
	buf = append(buf, portBytes...)
	m.writeBytes(buf)
}

func (m *mockToxConn) writeSOCKS5RequestIPv6(ip net.IP, port uint16) {
	buf := make([]byte, 0, 5+net.IPv6len+2)
	buf = append(buf, Version, CmdConnect, 0x00, AddressTypeIPV6)
	buf = append(buf, ip.To16()...)
	portBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(portBytes, port)
	buf = append(buf, portBytes...)
	m.writeBytes(buf)
}

func testConfig() *conf.Config {
	c := &conf.Config{}
	c.Default()
	return c
}

func TestNegotiate_Domain(t *testing.T) {
	mock := newMockToxConn()
	mock.writeSOCKS5Greeting()
	mock.writeSOCKS5RequestDomain("example.com", 443)

	req := NewRequest(testConfig(), mock)
	err := req.negotiate()
	require.NoError(t, err)

	assert.Equal(t, byte(CmdConnect), req.cmd)
	assert.Equal(t, byte(AddressTypeDomain), req.atyp)
	assert.Equal(t, "example.com:443", req.addr)
}

func TestNegotiate_IPv4(t *testing.T) {
	mock := newMockToxConn()
	mock.writeSOCKS5Greeting()
	mock.writeSOCKS5RequestIPv4(net.IPv4(127, 0, 0, 1), 8080)

	req := NewRequest(testConfig(), mock)
	err := req.negotiate()
	require.NoError(t, err)

	assert.Equal(t, byte(CmdConnect), req.cmd)
	assert.Equal(t, byte(AddressTypeIPV4), req.atyp)
	assert.Equal(t, "127.0.0.1:8080", req.addr)
}

func TestNegotiate_IPv6(t *testing.T) {
	mock := newMockToxConn()
	mock.writeSOCKS5Greeting()
	mock.writeSOCKS5RequestIPv6(net.ParseIP("::1"), 9090)

	req := NewRequest(testConfig(), mock)
	err := req.negotiate()
	require.NoError(t, err)

	assert.Equal(t, byte(CmdConnect), req.cmd)
	assert.Equal(t, byte(AddressTypeIPV6), req.atyp)
	assert.Equal(t, "[::1]:9090", req.addr)
}

func TestNegotiate_BadVersion(t *testing.T) {
	mock := newMockToxConn()
	mock.writeBytes([]byte{4, 1, MethodNone}) // version 4, not 5

	req := NewRequest(testConfig(), mock)
	err := req.negotiate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not socks5 protocol")
}

func TestNegotiate_UnsupportedCommand(t *testing.T) {
	mock := newMockToxConn()
	mock.writeSOCKS5Greeting()

	// Send a UDP ASSOCIATE request (0x03) instead of CONNECT
	buf := []byte{Version, CmdUDPAssociate, 0x00, AddressTypeIPV4}
	ip := net.IPv4(10, 0, 0, 1).To4()
	buf = append(buf, ip...)
	portBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(portBytes, 53)
	buf = append(buf, portBytes...)
	mock.writeBytes(buf)

	req := NewRequest(testConfig(), mock)
	err := req.negotiate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not support cmd")
}

func TestNegotiate_InvalidAddressType(t *testing.T) {
	mock := newMockToxConn()
	mock.writeSOCKS5Greeting()

	// Address type 0x02 is not supported
	buf := []byte{Version, CmdConnect, 0x00, 0x02}
	// still need some dummy address data
	buf = append(buf, net.IPv4(1, 2, 3, 4).To4()...)
	portBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(portBytes, 80)
	buf = append(buf, portBytes...)
	mock.writeBytes(buf)

	req := NewRequest(testConfig(), mock)
	err := req.negotiate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not support adrress type")
}

func TestNegotiate_ReplyFormat(t *testing.T) {
	mock := newMockToxConn()
	mock.writeSOCKS5Greeting()
	mock.writeSOCKS5RequestDomain("example.com", 443)

	req := NewRequest(testConfig(), mock)
	err := req.negotiate()
	require.NoError(t, err)

	// Verify the reply sent back to client
	reply := mock.writeBuf.Bytes()

	// Method selection reply: first 2 bytes
	require.GreaterOrEqual(t, len(reply), 2)
	assert.Equal(t, []byte{Version, MethodNone}, reply[:2])

	// CONNECT reply: after method negotiation
	// Format: [VER, REP, RSV, ATYP, BND.ADDR, BND.PORT]
	connectReply := reply[2:]
	require.GreaterOrEqual(t, len(connectReply), 4)
	assert.Equal(t, byte(Version), connectReply[0])
	assert.Equal(t, byte(0x00), connectReply[1]) // REP = success
}

func TestNegotiate_NoMethods(t *testing.T) {
	mock := newMockToxConn()
	// Version OK but 0 methods
	mock.writeBytes([]byte{Version, 0})

	req := NewRequest(testConfig(), mock)
	err := req.negotiate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nmethods not valid")
}

func TestHandle_CloseOnExit(t *testing.T) {
	mock := newMockToxConn()
	mock.writeSOCKS5Greeting()
	mock.writeSOCKS5RequestDomain("192.0.2.1", 9999) // TEST-NET, will fail to dial

	c := testConfig()
	c.ConnectTimeout = 1 // 1 second timeout instead of default 3s

	req := NewRequest(c, mock)
	err := req.Handle()
	assert.Error(t, err) // dial fails
	assert.True(t, mock.closed)
}

func TestNegotiate_ConnectionError(t *testing.T) {
	mock := &errorToxConn{err: errors.New("connection reset")}
	req := NewRequest(testConfig(), mock)
	err := req.negotiate()
	assert.Error(t, err)
}

// errorToxConn returns an error on every read.
type errorToxConn struct {
	err error
}

func (m *errorToxConn) Read(p []byte) (int, error)  { return 0, m.err }
func (m *errorToxConn) Write(p []byte) (int, error) { return 0, nil }
func (m *errorToxConn) Close() error                { return nil }
func (m *errorToxConn) CloseWrite() error           { return nil }
