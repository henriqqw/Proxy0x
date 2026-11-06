# Backend

Go backend for the Proxy0x desktop app. All exported methods on the `App` struct are exposed to the frontend via Wails bindings.

## App Struct

```go
type App struct {
    ctx        context.Context
    proxies    map[string]*ProxyEntry
    clients    map[string]*http.Client
    mu         sync.RWMutex
    lpMu       sync.Mutex
    localProxy *proxy.LocalServer
    stopCh     chan struct{}
    logs       []LogEntry
    logMu      sync.RWMutex
    crypt      *crypto.Store
}
```

- `proxies` — in-memory proxy pool keyed by random hex ID
- `clients` — HTTP clients keyed by `id:sessionId`, one per proxy+session pair
- `localProxy` — the single active local proxy server instance
- `stopCh` — signals background goroutines to stop on shutdown
- `crypt` — AES-GCM credential encryption via `internal/crypto`

## Lifecycle

### `startup(ctx)`
1. Stores the Wails context
2. Clears any leftover Windows system proxy settings
3. Loads proxies from `%APPDATA%/Proxy0x/proxies.json` (decrypts credentials)
4. Launches system tray in a goroutine
5. Launches health-check loop in a goroutine
6. Launches auto-rotation loop in a goroutine

### `shutdown(ctx)`
1. Stops the local proxy server if running
2. Clears Windows system proxy
3. Closes `stopCh` to signal background goroutines

## Proxy Management

### CRUD

| Method | Description |
|--------|-------------|
| `AddProxy(raw string)` | Parse `host:port[:user:pass]` format, assign random ID and session ID, persist |
| `AddProxiesBatch(raw string)` | Parse multiple proxies (newline-separated), persist all at once |
| `RemoveProxy(id)` | Delete from map, close HTTP client, persist |
| `GetProxies()` | Return all proxies sorted by `Order` field |
| `ReorderProxies(ids)` | Update `Order` field for each proxy by array index |
| `ExportProxies()` | Return all proxies as `host:port:user:pass` text |

### Proxy Format Parsing

Supported formats:
- `host:port` — HTTP proxy without auth
- `host:port:user:pass` — HTTP proxy with auth
- `socks5://host:port` — SOCKS5 proxy without auth
- `socks5://host:port:user:pass` — SOCKS5 with auth
- `http://host:port` — explicit HTTP scheme

### Groups

Proxies can be tagged with a group name. Groups are derived dynamically from the proxy set (no separate storage).

- `GetGroups()` — return sorted unique group names
- `SetProxyGroup(id, group)` — assign a group to a proxy
- `ConnectGroup(group)` / `DisconnectGroup(group)` — batch operations

## Connection Flow

### `ConnectProxy(id)`

1. Disconnect any currently active proxy
2. Set status to `connecting`, create HTTP client with proxy transport
3. Hit `https://api.ipify.org?format=json` to verify connectivity and get the exit IP
4. Fetch geolocation from `https://ipwho.is/{ip}`
5. Set status to `connected`, record IP, latency, geo, `connectedAt` timestamp
6. Start local proxy server and set Windows system proxy

### `DisconnectProxy(id)`

1. Accumulate uptime into `TotalUptime`
2. Reset status to `offline`, clear IP/latency/geo
3. Close the HTTP client transport
4. Stop local proxy server if it was for this proxy
5. Clear Windows system proxy

### `RotateProxy(id)`

1. Close old HTTP client (keyed by old session ID)
2. Generate a new random session ID
3. Create new HTTP client (new session ID → new upstream IP)
4. Verify new IP via ipify, fetch geo
5. Restart local proxy server with new upstream config

Session rotation works by appending `_session-{id}` to the proxy username. Each unique session ID gets a different exit IP from residential proxy pools.

## Local Proxy Server

Defined in `internal/proxy/server.go`. The `LocalServer` struct:

- Binds to `127.0.0.1:0` (random available port)
- Accepts incoming connections and routes them through the upstream proxy
- Supports HTTP CONNECT tunneling (for HTTPS) and plain HTTP forwarding
- Plain HTTP: one request per client connection. The request is forwarded with `Connection: close` and the client connection is closed after the response body, so browsers open a new connection for the next request (no keep-alive through the local proxy)
- Supports SOCKS5 upstream via `golang.org/x/net/proxy`; upstream dials (HTTP and SOCKS5) time out after 15s
- Tracks all client and upstream connections; `Stop()` closes them all and returns in under a second
- Temporary `Accept` errors back off (5ms doubling up to 1s) instead of spinning

### System Proxy (`internal/proxy/system.go`)

Uses Windows registry to configure system-wide proxy:
- `SetSystemProxy(port)` — sets `ProxyServer`, `ProxyOverride`, `ProxyEnable` in `HKCU\...\Internet Settings`
- `ClearSystemProxy()` — sets `ProxyEnable` to `0`
- Calls `InternetSetOptionW` via `wininet.dll` to refresh proxy settings immediately

Bypass list: `localhost;127.*;10.*;192.168.*;<local>`

## Health Check Loop

Runs every 30 seconds (polls). Actual checks happen at configurable intervals (default 180s).

1. Collect all active proxy IDs
2. For each, call `CheckHealth(id)` — hit ipify to verify connectivity
3. If a connected proxy becomes unreachable, emit `proxy-dropped` event (triggers Windows notification)
4. 1.5s delay between checks to avoid bursting

## Auto-Rotation Loop

Runs every 10 seconds (polls). Actual rotations happen at configurable intervals (0 = disabled).

1. If `AutoRotateInterval > 0` and enough time has passed since last rotation
2. Rotate all active proxies
3. Log the batch rotation

## Speed Benchmark

`BenchmarkProxy(id)` returns `map[string]interface{}`:

1. **Connect latency** — time to reach `api.ipify.org` through the proxy (ms)
2. **Download speed** — fetch 100KB from `speed.cloudflare.com/__down`, compute Mbps

## Logging

Ring buffer of 200 `LogEntry` structs (action, proxyID, host, detail, timestamp). Newest entries replace oldest. Emitted via Wails events for real-time UI updates.

## Persistence

Proxies stored at `%APPDATA%/Proxy0x/proxies.json`:
- On save: credentials encrypted with AES-GCM, status/active/IP fields stripped
- On load: credentials decrypted, status reset to `offline`, active reset to `false`

## Configuration

`internal/config/config.go` — JSON file at `%APPDATA%/Proxy0x/config.json`:

```go
type Config struct {
    Language            string
    Theme               string
    AutoRotateInterval  int    // seconds, 0 = disabled
    HealthCheckInterval int    // seconds, minimum 60, default 180
}
```

## Encryption

`internal/crypto/crypto.go` — AES-256-GCM with:
- 256-bit key stored at `%APPDATA%/Proxy0x/.key` (auto-generated on first use)
- Encrypted values prefixed with `enc:` followed by base64(nonce + ciphertext)
- Empty strings pass through unencrypted
- Non-prefixed strings treated as plaintext (backward compatibility)

## MCP Integration

`mcp_integration.go` — methods for Claude Desktop config management:

| Method | Description |
|--------|-------------|
| `GetMCPStatus()` | Check if MCP binary exists and if proxy0x is in Claude's config |
| `InstallMCPConfig()` | Write proxy0x entry to `%APPDATA%/Claude/claude_desktop_config.json` |
| `UninstallMCPConfig()` | Remove proxy0x entry from Claude's config |
| `CheckMCPHealth()` | Spawn MCP binary, send JSON-RPC `initialize` request, verify response |

## System Tray

`tray.go` — uses `energye/systray`:
- Custom icon from `build/windows/tray.ico`
- Language submenu (16 languages with flag icons)
- Connect All / Disconnect All / Rotate All
- Show / Quit
- Updates language in real time via Wails event listener

## External Services

| Service | Purpose | Used by |
|---------|---------|---------|
| `api.ipify.org` | Get exit IP through proxy | Connect, Rotate, Health Check, Benchmark |
| `ipwho.is` | IP geolocation (country, city, ISP, ASN, timezone) | Connect, Rotate |
| `speed.cloudflare.com` | Download speed measurement (100KB file) | Benchmark |
