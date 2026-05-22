# Unit Tests for util/ and socks5/ Packages

## Overview

Add unit tests for `util/` and `socks5/` packages using the `github.com/stretchr/testify` framework. The project currently has zero test files.

## Dependencies

- Add `github.com/stretchr/testify` to `go.mod` (currently only in `go.sum` as indirect)

## Test Plan by File

### util/kdf_test.go

- `TestKDF` — verify deterministic PBKDF2 output for known inputs
- `TestHashedPassword` — verify base64-encoded hash is consistent

### util/nonce_test.go

- `TestNextNonce` — increment, carry, overflow wrap-around

### util/jwt_test.go

- `TestGenerateAndValidateJwtToken` — round-trip: generate token then validate it
- `TestValidateJwtToken_InvalidToken` — reject malformed tokens

### util/url_test.go

- `TestFormatURL` — table-driven: add default path, add default port per scheme

### util/json_test.go

- `TestStringify` — marshal struct/slice/map to JSON string

### util/option_test.go

- `TestWithPassword` — verify password is hashed via KDF
- `TestWithTunnel` / `TestWithLocalAddress` / etc — each option sets its field
- `TestToToxOptions` — combine multiple options

### util/conn_test.go

- `TestToxConnection_Read` — read DATA frame, CLOSE_WRITE returns EOF
- `TestToxConnection_Write` — write encodes frame correctly
- `TestToxConnection_CloseWrite` — sends CLOSE_WRITE, second call is no-op
- `TestToxConnection_Close` — delegates to underlying closer

Mock: hand-written `mockReadWriteCloser` implementing `io.ReadWriteCloser` with configurable read/write behavior. The mock must simulate the binary frame format ([4-byte len][msgpack Frame]).

### util/net_test.go

- `TestTimeoutConn_Read` — deadline is set when timeout > 0
- `TestTimeoutConn_Write` — deadline is set when timeout > 0
- `TestTimeoutConn_ZeroTimeout` — no deadline set when timeout == 0

Mock: hand-written `mockNetConn` implementing `net.Conn`.

### util/io_test.go

- `TestCopyBuffer` — copies data from reader to writer

### util/time_test.go

- `TestNowInMills` — returns millisecond-precision timestamp

### socks5/request_test.go (socks5 package)

- `TestNegotiate_ValidDomain` — valid SOCKS5 handshake with domain address
- `TestNegotiate_ValidIPv4` — valid SOCKS5 handshake with IPv4 address
- `TestNegotiate_ValidIPv6` — valid SOCKS5 handshake with IPv6 address
- `TestNegotiate_BadVersion` — rejects non-SOCKS5 version
- `TestNegotiate_UnsupportedCommand` — rejects non-CONNECT commands

Mock: hand-written `mockToxConn` implementing `util.ToxConn` (reader + writer + closer + closeWriter).

## Files to Skip

- `util/version.go` — only prints to stdout, low value
- `util/profiling.go` — side-effect only (starts HTTP server), low value
- `conf/`, `tunnel/`, `h3conn/`, `cmd/` — out of scope for this round

## Mock Strategy

No testify/mock — hand-write minimal mock structs. This keeps tests portable and avoids complexity, since the interfaces are small (1-4 methods each).
