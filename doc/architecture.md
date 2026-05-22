# tox Architecture

## Overview

tox is a TCP-over-tunnel proxy. It tunnels TCP streams between a local client and a remote server over various application-layer protocols.

## Components

```
                    ┌─────────────────────────────┐
                    │        tox local             │
                    │                              │
  ┌──────────┐      │  ┌────────┐    ┌─────────┐  │
  │ client   │◄─TCP─┤  │ SOCKS5 │───►│ tunnel  │  │
  │ app      │      │  │ handler│    │ client  │  │
  └──────────┘      │  └────────┘    └────┬────┘  │
                    │                     │       │
                    └─────────────────────┼───────┘
                                          │
                          tunnel protocol │ (gRPC/HTTP2/HTTP3/WS)
                                          │
                    ┌─────────────────────┼───────┐
                    │        tox server   │       │
                    │                     ▼       │
                    │  ┌────────┐    ┌─────────┐  │
                    │  │ SOCKS5 │◄───┤ tunnel  │  │
                    │  │ request│    │ server  │  │
                    │  └───┬────┘    └─────────┘  │
                    │      │                      │
                    └──────┼──────────────────────┘
                           │ TCP
                           ▼
                     ┌──────────┐
                     │ target   │
                     │ server   │
                     └──────────┘
```

## Code Structure

```
main.go          # entry point → cmd.Execute()
cmd/
  root.go        # cobra root + local/server subcommands
  server.go      # server startup
  local.go       # local (client) startup
conf/
  config.go      # YAML config loading, defaults
tunnel/
  tunnel.go      # Client/Server interfaces + dispatch by URL scheme
  grpc/          # gRPC tunnel (bidirectional stream + JWT)
  h2/            # HTTP/2 tunnel (h2conn full-duplex)
  quic/          # HTTP/3 tunnel (quic-go + h3conn full-duplex)
  websocket/     # WebSocket tunnel
h3conn/          # HTTP/3 bidirectional connection (client + server accepting)
socks5/          # SOCKS5 handshake + target TCP dial
util/
  conn.go        # ToxConn: frame protocol over raw I/O
  option.go      # Functional options pattern
  kdf.go         # PBKDF2 password hashing
  jwt.go         # JWT token generation/validation
  net.go         # TimeoutConn wrapper
  io.go          # CopyBuffer with pool
  url.go         # URL normalization (default ports, path)
proto/           # protobuf schema for gRPC tunnel
config/
  example.config.yaml
```

## Data Flow

1. Client app connects to `tox local` via SOCKS5 (default `:1080`)
2. `local` performs SOCKS5 handshake, extracts target address
3. `local` establishes a tunnel connection to `tox server`
4. `local` begins bidirectional copy: client ↔ tunnel
5. `server` accepts tunnel connection, reads SOCKS5 request from stream
6. `server` dials the target TCP address
7. `server` begins bidirectional copy: tunnel ↔ target

## Tunnel Interface

Every tunnel protocol implements the same two interfaces defined in `tunnel/tunnel.go`:

```go
type Client interface {
    Connect(context.Context) (util.ToxConn, error)
}

type Server interface {
    ListenAndServe(handler func(util.ToxConn)) error
}
```

The dispatch in `NewClient`/`NewServer` selects the implementation based on the URL scheme of the `tunnel` config field.
