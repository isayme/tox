# Configuration Guide

## Config File

Configuration is loaded from a YAML file. The path is specified via the `CONF_FILE_PATH` environment variable.

Example (`config/example.config.yaml`):

```yaml
logger:
  level: debug

tunnel: grpcs://127.0.0.1:1443
localAddress: 0.0.0.0:1080
password: yourpassword
timeout: 120
connectTimeout: 5
# certFile: ./testdata/server.pem
# keyFile: ./testdata/server.key
# insecureSkipVerify: true
```

## Fields

### Common (client and server)

| Field | Type | Description |
|-------|------|-------------|
| `tunnel` | string | Tunnel URL. Scheme determines the transport protocol (required) |
| `password` | string | Authentication password, hashed with PBKDF2 before use (required) |
| `timeout` | int | I/O timeout in seconds |

### Server-only

| Field | Type | Description |
|-------|------|-------------|
| `certFile` | string | TLS certificate file path |
| `keyFile` | string | TLS key file path |

### Client-only

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `localAddress` | string | `:1080` | Local SOCKS5 listen address |
| `connectTimeout` | int | `3` | TCP connect timeout (seconds) |
| `insecureSkipVerify` | bool | `false` | Skip TLS certificate verification |

## Tunnel URL Schemes

The `tunnel` field is a URL that selects both the transport protocol and the server address.

| Scheme | Protocol | TLS | Server cert required |
|--------|----------|-----|---------------------|
| `grpc://host:port` | gRPC | No | No |
| `grpcs://host:port` | gRPC | Yes | Optional |
| `http2://host:port` | HTTP/2 | Yes | Yes |
| `h2://host:port` | HTTP/2 | Yes | Yes |
| `quic://host:port` | HTTP/3 (QUIC) | Yes | Yes |
| `http3://host:port` | HTTP/3 (QUIC) | Yes | Yes |
| `ws://host:port` | WebSocket | No | No |
| `wss://host:port` | WebSocket | Yes | Optional |

Default ports are assigned automatically: port 80 for `ws` and `grpc`, port 443 for all TLS-enabled schemes.

## Authentication

The password is hashed using PBKDF2 with a hardcoded salt (`"tox"`), 1024 iterations, producing a 32-byte key encoded as base64.

- **gRPC**: The hashed password is sent as a per-RPC metadata entry (`token`)
- **WebSocket**: Sent as an HTTP header (`token`) during handshake
- **HTTP/2, HTTP/3**: Sent as an HTTP header (`token`)

JWT support is available in `util/jwt.go` but is not wired into any tunnel transport.
