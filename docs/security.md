# Security

## Credential encryption

Proxy credentials (username and password) are encrypted at rest using **AES-256-GCM**.

### Key management

- A 32-byte random key is generated on first use and stored at `%AppData%/Proxy0x/.key`
- The key file has `0600` permissions (owner read/write only)
- One key per machine — shared between the Wails app and the MCP server
- Implementation: `internal/crypto/crypto.go` (Wails app), `pkg/proxycore/store.go` (MCP server)

### Encryption format

```
"enc:" + base64(nonce || ciphertext)
```

- 12-byte random nonce prepended to the ciphertext
- Standard library `crypto/aes` + `crypto/cipher` (GCM mode)
- Empty strings are stored as empty (not encrypted)
- Plaintext values (without `enc:` prefix) are returned as-is for backward compatibility

### Store lifecycle

1. **Save:** encrypt username and password, strip runtime state (status, currentIp, active), write JSON
2. **Load:** read JSON, decrypt credentials, reset all proxies to `offline` / `active: false`

## SSRF protection

The MCP server's `proxy_fetch` tool validates every URL before making a request through the proxy pool.

### Blocked schemes

Only `http` and `https` are allowed. `ftp`, `file`, `javascript`, and schemeless URLs are rejected.

### Blocked IP ranges

| Range | Reason |
|-------|--------|
| `10.0.0.0/8` | RFC 1918 private |
| `172.16.0.0/12` | RFC 1918 private |
| `192.168.0.0/16` | RFC 1918 private |
| `127.0.0.0/8` | Loopback |
| `0.0.0.0/8` | "This" network |
| `169.254.0.0/16` | Link-local |
| `::1/128` | IPv6 loopback |
| `fc00::/7` | IPv6 ULA |
| `fe80::/10` | IPv6 link-local |

### Cloud metadata IPs

| IP | Service |
|----|---------|
| `169.254.169.254` | AWS EC2 / GCP / Azure metadata |
| `169.254.170.2` | AWS ECS task metadata |
| `100.100.100.200` | Alibaba Cloud metadata |
| `fd00:ec2::254` | AWS EC2 IPv6 metadata |

### Blocked hostnames

- `localhost`
- `metadata`
- Any hostname ending in `.local`, `.localhost`, `.internal`
- `metadata.google.internal`

### DNS resolution check

After hostname validation, the URL's hostname is resolved via DNS. Every returned IP address is checked against the blocked IP ranges above. This prevents DNS rebinding attacks where a hostname resolves to a private IP.

### Redirect validation

HTTP redirects are followed up to a maximum of **10 hops**. Each redirect target URL is fully SSRF-validated (scheme, hostname, IP, DNS) before following. This prevents open-redirect bypasses.

### Body size limit

Response bodies from `proxy_fetch` are capped at **100 KB**. Larger responses are truncated and flagged with a `truncated: true` field.

## Credential stripping

The `proxy_list` MCP tool and the `proxy://pool` MCP resource strip `username` and `password` fields from output. AI agents see proxy hosts and ports but never credentials.

## Windows system proxy

When a proxy is connected, the app:

1. Starts a local HTTP/SOCKS5 server on `127.0.0.1`
2. Sets `HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings\ProxyServer` and `ProxyEnable`
3. All browser traffic routes through the local server → upstream proxy

On disconnect or app shutdown:
- System proxy settings are cleared (`ProxyEnable = 0`)
- `clearSystemProxy()` is called in the `shutdown()` lifecycle hook
- Tracked TCP connections are force-closed for instant disconnect

## Connection tracking

The local proxy server tracks every TCP connection it handles — the client side and the upstream side, for both CONNECT tunnels and plain HTTP requests — via `trackConn`/`untrackConn`. `Stop()` closes every tracked connection before waiting for its goroutines, so it returns in under a second regardless of in-flight traffic.

## Implementation locations

| Component | File |
|-----------|------|
| AES-GCM encryption (Wails) | `internal/crypto/crypto.go` |
| AES-GCM encryption (MCP) | `pkg/proxycore/store.go` |
| SSRF validation | `pkg/proxycore/ssrf.go` |
| SSRF tests | `pkg/proxycore/ssrf_test.go` |
| Redirect validation | `pkg/proxycore/fetcher.go` |
| Credential stripping | `internal/mcpserver/tools.go` |
| System proxy | `internal/proxy/system.go` |
| Connection tracking | `internal/proxy/server.go` |
