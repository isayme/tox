# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

tox is a TCP-over-tunnel proxy that forwards TCP streams over gRPC, HTTP/2, HTTP/3 (QUIC), or WebSocket. It has two modes:

- **server** — accepts tunnel connections from clients and forwards traffic to the target TCP address (via SOCKS5)
- **local** — listens on a local TCP port, wraps incoming connections into a tunnel to the server

## Commands

```bash
# Build
go build ./...

# Run tests
go test ./...

# Run local client (SOCKS5 proxy on :1080)
CONF_FILE_PATH=./config/example.config.yaml go run main.go local

# Run server
CONF_FILE_PATH=./config/example.config.yaml go run main.go server

# Show version
go run main.go -v

# Enable pprof (exposed on :6060)
go run main.go --profiling server

# Regenerate protobuf/gRPC code
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  proto/tunnel.proto
```

## Architecture

```
┌──────────┐     tunnel (gRPC/HTTP2/HTTP3/WS)     ┌──────────┐
│  local   │ ◄──────────────────────────────────► │  server  │
│ (:1080)  │                                       │ (:1443)  │
└────┬─────┘                                       └────┬─────┘
     │ raw TCP                                          │ TCP
     ▼                                                  ▼
  client app                                       target addr
  (SOCKS5 client)
```

**Entry point** — `main.go` dispatches to `cmd.Execute()`. Commands are defined with `spf13/cobra` in `cmd/root.go`. The `local` subcommand calls `startLocal()`, `server` calls `startServer()`.

**Tunnel dispatch** — `tunnel/tunnel.go` defines `Client` and `Server` interfaces. `NewClient` / `NewServer` parse the tunnel URL scheme and return the correct implementation:

| Scheme | Package | Transport |
|--------|---------|-----------|
| `grpc` / `grpcs` | `tunnel/grpc` | gRPC bidirectional stream + per-RPC token auth |
| `http2` / `h2` | `tunnel/h2` | HTTP/2 full-duplex via `h2conn` |
| `ws` / `wss` | `tunnel/websocket` | WebSocket (`golang.org/x/net/websocket`) |
| `quic` / `http3` | `tunnel/quic` | HTTP/3 via `quic-go/http3`, custom `h3conn` for bidirectional I/O |

**Frame protocol** — `util/conn.go` wraps any `io.ReadWriteCloser` into `ToxConn`, adding a binary frame layer using msgpack:

```
[4-byte big-endian length][msgpack Frame{Cmd, Data}]
```

Commands: `0x01` (DATA), `0x02` (CLOSE_WRITE). `CloseWrite()` signals half-close to the remote side.

**SOCKS5** — `socks5/request.go` implements a subset of SOCKS5 (CONNECT command, no auth, IPv4/IPv6/domain address types).

**Authentication** — tunnel password is hashed with PBKDF2 (`util/kdf.go`) and verified per-protocol: gRPC uses per-RPC metadata with the hashed password as a `token` entry; WebSocket and HTTP/2/3 use an HTTP `token` header. JWT helper functions exist in `util/jwt.go` but are not wired into any transport currently.

**Configuration** — `conf/config.go` reads a YAML config file via `go-config` (path from env `CONF_FILE_PATH`). `Config.Default()` fills in defaults for `ConnectTimeout` (3s) and `LocalAddress` (:1080).

**Options pattern** — `util/option.go` uses the functional options pattern (`ToxOption` interface) to build a `ToxOptions` struct passed to tunnel constructors.

## Key dependencies

| Package | Usage |
|---------|-------|
| `google.golang.org/grpc` | gRPC tunnel transport |
| `github.com/quic-go/quic-go/http3` | HTTP/3 (QUIC) tunnel transport |
| `github.com/posener/h2conn` | HTTP/2 bidirectional connection |
| `golang.org/x/net/websocket` | WebSocket tunnel transport |
| `github.com/vmihailenco/msgpack/v5` | Frame serialization |
| `github.com/spf13/cobra` | CLI command structure |
| `github.com/isayme/go-config` | YAML config loading |
| `golang.org/x/crypto` | PBKDF2 password hashing |
| `gopkg.in/DataDog/dd-trace-go.v1` | Datadog profiling |

## Config file format

See `config/example.config.yaml`. `CONF_FILE_PATH` env var must be set. The `tunnel` field URL scheme determines the transport protocol. TLS-enabled schemes (`grpcs`, `wss`, `http2`, `quic`) may require `certFile`/`keyFile` on the server and support `insecureSkipVerify` on the client.
