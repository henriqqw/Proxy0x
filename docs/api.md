# Wails API Reference

All public methods on the `App` struct in `app.go` are exposed to the React frontend via Wails bindings. TypeScript declarations are auto-generated at `frontend/wailsjs/go/main/App.d.ts`.

## Proxy CRUD

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `AddProxy` | `raw string` | `ProxyEntry` | Parse `host:port[:user:pass]` and add to pool. Supports `socks5://` prefix. |
| `AddProxiesBatch` | `text string` | `[]ProxyEntry` | Add multiple proxies, one per line. |
| `RemoveProxy` | `id string` | — | Remove proxy by ID. Disconnects first if active. |
| `GetProxies` | — | `[]ProxyEntry` | List all proxies, ordered by `Order` field. |
| `ReorderProxies` | `ids []string` | — | Set display order from an array of proxy IDs. |

## Connection

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `ConnectProxy` | `id string` | `ProxyEntry` | Start local proxy server, set Windows system proxy, verify IP via api.ipify.org, fetch geolocation. |
| `DisconnectProxy` | `id string` | `ProxyEntry` | Stop local proxy, clear system proxy, close tracked connections. |
| `RotateProxy` | `id string` | `ProxyEntry` | Generate new session ID → new IP. Reconnects automatically. |
| `RestartProxy` | `id string` | `ProxyEntry` | Disconnect then reconnect (full restart cycle). |

## Groups

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `GetGroups` | — | `[]string` | Distinct group names across all proxies. |
| `SetProxyGroup` | `id string, group string` | `ProxyEntry` | Assign a proxy to a group (empty string removes). |
| `ConnectGroup` | `group string` | — | Connect all proxies in a group. |
| `DisconnectGroup` | `group string` | — | Disconnect all proxies in a group. |

## Monitoring

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `CheckHealth` | `id string` | `ProxyEntry` | Test proxy connectivity, update latency and status. |
| `BenchmarkProxy` | `id string` | `map[string]any` | Speed test: measure connect latency (ms) and download speed (Mbps). |
| `GetMyIP` | — | `string` | Fetch current public IP via api.ipify.org (direct, not through proxy). |
| `GetConnectionLog` | — | `[]LogEntry` | Last 200 log entries (newest first). |

## Settings

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `GetAutoRotateInterval` | — | `int` | Current auto-rotate interval in seconds (0 = disabled). |
| `SetAutoRotateInterval` | `seconds int` | — | Set auto-rotate interval. 0 disables. |
| `GetHealthCheckInterval` | — | `int` | Current health check interval in seconds. |
| `SetHealthCheckInterval` | `seconds int` | — | Set health check interval. |
| `GetLanguage` | — | `string` | Current language code (e.g. `en`, `pt-BR`). |
| `SetLanguage` | `code string` | — | Change UI language. |
| `GetTheme` | — | `string` | Current theme (`dark` or `light`). |
| `SetTheme` | `theme string` | — | Change theme. |

## Security

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `GetKillSwitch` | — | `bool` | Whether the kill switch (firewall rules) is active. |
| `SetKillSwitch` | `enabled bool` | — | Enable/disable kill switch. Blocks all traffic except through the proxy app. Requires UAC elevation. |
| `GetAutoFailover` | — | `bool` | Whether auto-failover is enabled. |
| `SetAutoFailover` | `enabled bool` | — | Enable/disable auto-failover. On proxy drop, connects to the next best proxy automatically. |

## Tools

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `CheckDNSLeak` | `id string` | `DNSLeakResult` | Compare exit IP through proxy vs direct IP. Detects DNS leaks. |
| `BulkCheckProxies` | — | `[]BulkCheckResult` | Test all proxies in parallel (10 workers). Returns latency, anonymity level, status. |
| `ConnectFastest` | — | `string` | Run health checks on all proxies, connect the one with lowest latency. Returns proxy ID. |
| `EnableUWPSupport` | — | — | Run `CheckNetIsolation` to exempt UWP/Store apps from loopback restrictions. Requires UAC elevation. |

## Proxy Naming

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `SetProxyName` | `id string, name string` | `ProxyEntry` | Assign a custom display name to a proxy. |

## Import / Export

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `ExportProxies` | — | `string` | All proxies as `host:port:user:pass` text, one per line. |
| `AddProxiesBatch` | `text string` | `[]ProxyEntry` | Import from multi-line text (legacy, one per line). |
| `AddProxiesMultiFormat` | `text string` | `[]ProxyEntry` | Import from multi-format text: `socks5://user:pass@host:port` URI, CSV, or `host:port:user:pass`. |

## MCP Integration

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `GetMCPStatus` | — | `MCPStatus` | Check if MCP is configured and binary exists. |
| `InstallMCPConfig` | — | — | Write proxy0x entry to Claude Desktop config. |
| `UninstallMCPConfig` | — | — | Remove proxy0x entry from Claude Desktop config. |
| `CheckMCPHealth` | — | `MCPStatus` | Spawn MCP binary, send initialize request, verify response. |

### MCPStatus

```go
type MCPStatus struct {
    Installed   bool   `json:"installed"`
    ExePath     string `json:"exePath"`
    ExeExists   bool   `json:"exeExists"`
    ConfigPath  string `json:"configPath"`
    HealthOK    bool   `json:"healthOk"`
    HealthError string `json:"healthError,omitempty"`
}
```

## Window

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `MinimizeWindow` | — | — | Minimize to taskbar. |
| `CloseWindow` | — | — | Hide window (stays in system tray). |
| `ShowWindow` | — | — | Restore from tray. |
| `ToggleMaximize` | — | — | Toggle maximize/restore. |
| `QuitApp` | — | — | Stop local proxy, clear system proxy, exit. |

## Types

### ProxyEntry

```go
type ProxyEntry struct {
    ID          string      `json:"id"`
    Host        string      `json:"host"`
    Port        string      `json:"port"`
    Username    string      `json:"username"`
    Password    string      `json:"password"`
    Protocol    string      `json:"protocol"`       // "http" or "socks5"
    Group       string      `json:"group"`
    Name        string      `json:"name"`            // custom display name
    Status      ProxyStatus `json:"status"`          // offline|connecting|connected|error
    CurrentIP   string      `json:"currentIp"`
    Location    *GeoInfo    `json:"location"`
    Latency     int64       `json:"latency"`         // ms
    LastCheck   string      `json:"lastCheck"`
    Active      bool        `json:"active"`
    SessionID   string      `json:"sessionId"`
    Order       int         `json:"order"`
    ConnectedAt string      `json:"connectedAt"`
    TotalUptime int64       `json:"totalUptime"`     // seconds
}
```

### DNSLeakResult

```go
type DNSLeakResult struct {
    ProxyIP  string `json:"proxyIp"`
    DirectIP string `json:"directIp"`
    Leaking  bool   `json:"leaking"`
}
```

### BulkCheckResult

```go
type BulkCheckResult struct {
    ID        string `json:"id"`
    Host      string `json:"host"`
    Port      string `json:"port"`
    Status    string `json:"status"`    // "ok" or "fail"
    Latency   int64  `json:"latency"`   // ms
    Anonymity string `json:"anonymity"` // "transparent", "anonymous", "elite"
    Error     string `json:"error,omitempty"`
}
```

### LogEntry

```go
type LogEntry struct {
    Time    string `json:"time"`
    Action  string `json:"action"`
    ProxyID string `json:"proxyId"`
    Host    string `json:"host"`
    Detail  string `json:"detail"`
}
```

## Frontend bindings

Wails auto-generates TypeScript bindings at:
- `frontend/wailsjs/go/main/App.js` — runtime stubs
- `frontend/wailsjs/go/main/App.d.ts` — type declarations
- `frontend/wailsjs/go/models.ts` — Go struct types as TS classes

Import in React:
```typescript
import { ConnectProxy, GetProxies } from 'wailsjs/go/main/App'
```
