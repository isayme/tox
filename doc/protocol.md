# tox Frame Protocol

## Overview

tox wraps raw TCP streams in a binary frame protocol for transmission over tunnel transports (gRPC, HTTP/2, HTTP/3, WebSocket). This is implemented in `util/conn.go` as the `ToxConnection` type.

## Wire Format

Each frame consists of:

```
┌──────────────────────┬──────────────────────────┐
│   4 bytes (uint32)   │   N bytes                 │
│   big-endian length  │   msgpack-encoded Frame   │
└──────────────────────┴──────────────────────────┘
```

The payload is a msgpack-serialized `Frame`:

```go
type Frame struct {
    Cmd  uint8
    Data []byte
}
```

## Commands

| Value | Name | Description |
|-------|------|-------------|
| `0x01` | DATA | Carries a chunk of the TCP stream |
| `0x02` | CLOSE_WRITE | Signals half-close (no more data from sender) |

## Half-Close Semantics

`ToxConn` implements a `CloseWrite()` method that sends a `CLOSE_WRITE` frame. This signals to the remote peer that the local side has finished sending data, while the connection remains open for reading. This mirrors TCP's half-close behavior.

A subsequent `Write()` after `CloseWrite()` returns an error (`broken pipe`). Multiple calls to `CloseWrite()` are safe (idempotent via `atomic.Bool`).

Reading after the remote peer sends `CLOSE_WRITE` returns `io.EOF`.

## Buffering

`ToxConnection` uses a `bufferpool` for both the 4-byte length header and the msgpack frame payload. Reads are buffered: a full frame is read from the underlying connection, and `Read()` calls drain from the internal buffer until it's empty, then the next frame is read.

## gRPC-specific Framing

gRPC tunnel (`tunnel/grpc/`) does NOT use the msgpack frame protocol. Instead, it uses protobuf messages on a bidirectional gRPC stream:

```protobuf
service Tunnel {
  rpc OnConnect(stream Data) returns (stream Data);
}

message Data {
  bytes data = 1;
}
```

The gRPC client and server connectors implement `io.ReadWriteCloser` directly over the protobuf stream, reading/writing `Data` messages.
