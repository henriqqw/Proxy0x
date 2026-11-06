# Wails API Reference

Every exported method on `App` (`app.go`, `mcp_integration.go`) is bound to the React frontend. TypeScript declarations are generated in `frontend/wailsjs/go/main/App.d.ts`; a Go `error` return becomes a rejected Promise.

**Proxy values are copies.** Every `ProxyEntry` returned by a binding or sent in a `proxy-updated` event is a snapshot taken under the lock with `password` blanked. The password never reaches the webview; `ExportProxies` is the only call that returns credentials.

**Read-only store.** If `proxies.json` could not be loaded (invalid JSON or undecryptable credentials), mutating calls (`AddProxy`, `AddProxiesBatch`, `AddProxiesMultiFormat`, `RemoveProxy`, `ReorderProxies`, `SetProxyGroup`, `SetProxyName`) return an error and a `store-error` event is emitted.

## Proxy CRUD

| Method | Signature | Description |
|--------|-----------|-------------|
| `AddProxy` | `(raw string) (*ProxyEntry, error)` | Parse `[socks5://\|http://]host:port[:user:pass]`, add, persist. |
| `AddProxiesBatch` | `(raw string) ([]*ProxyEntry, error)` | One proxy per line (same format); invalid lines are skipped. |
| `AddProxiesMultiFormat` | `(raw string) ([]*ProxyEntry, error)` | Lines as URI (`socks5://`, `socks4://`, `http://`, `https://` with optional `user:pass@`), CSV `host,port[,user,pass[,protocol]]`, or `host:port[:user:pass]`. `#` lines are comments. `socks4` is stored as `socks5`. |
| `RemoveProxy` | `(id string) error` | Disconnects first if active (local server stopped, system proxy restored), then deletes and closes its HTTP client. |
| `GetProxies` | `() []*ProxyEntry` | All proxies sorted by `Order`. |
| `ReorderProxies` | `(ids []string) error` | Sets `Order` from the array index. |
| `ExportProxies` | `() string` | `host:port[:user:pass]` lines, **plaintext credentials**. |
| `SetProxyName` | `(id, name string) (*ProxyEntry, error)` | Custom display name. |

## Groups

| Method | Signature | Description |
|--------|-----------|-------------|
| `GetGroups` | `() []string` | Sorted distinct group names. |
| `SetProxyGroup` | `(id, group string) (*ProxyEntry, error)` | Empty string removes the group. |
| `ConnectGroup` | `(group string) error` | Connects each inactive proxy of the group in sequence. Only one proxy can be active, so the last successful one wins. Not used by the UI. |
| `DisconnectGroup` | `(group string) error` | Disconnects active proxies of the group. Not used by the UI. |

## Connection

| Method | Signature | Description |
|--------|-----------|-------------|
| `ConnectProxy` | `(id string) (*ProxyEntry, error)` | Disconnects the active proxy, verifies the exit IP (HTTP 200 + valid IP), fetches geolocation, starts the local server and points the system proxy at it. Any failure → `status=error`. |
| `DisconnectProxy` | `(id string) (*ProxyEntry, error)` | Stops the local server (returns in < 1 s), restores the user's system proxy, accumulates `totalUptime`. |
| `RotateProxy` | `(id string) (*ProxyEntry, error)` | Verifies a candidate session first; on success switches session and updates the local server in place (same port, open connections kept). On failure the previous session and status are kept and an error is returned. |
| `RestartProxy` | `(id string) (*ProxyEntry, error)` | Disconnect + connect. Not used by the UI. |
| `ConnectFastest` | `() (*ProxyEntry, error)` | Verifies every inactive proxy (10 in parallel), connects the one with the lowest latency. |

## Monitoring & tools

| Method | Signature | Description |
|--------|-----------|-------------|
| `CheckHealth` | `(id string) (*ProxyEntry, error)` | Exit IP check for an active proxy. Failure → `error` + `proxy-dropped` (and auto-failover if enabled); success restores `connected`. Also run by the health loop. |
| `BenchmarkProxy` | `(id string) (map[string]interface{}, error)` | `connectMs` (exit IP check) and `downloadMs`/`downloadMbps` (100 KB from speed.cloudflare.com). |
| `BulkCheckProxies` | `() []BulkCheckResult` | Checks all proxies (10 in parallel); emits `bulk-check-progress` per result; updates `latency`/`lastCheck`. |
| `CheckDNSLeak` | `() DNSLeakResult` | Compares the active proxy's exit IP with the direct public IP (`GetMyIP`). `leaked=true` when they match or when no proxy is active. This detects "traffic not going through the proxy", not DNS-resolver leaks. |
| `GetMyIP` | `() (string, error)` | Direct public IP via api.ipify.org, 10 s timeout. |
| `GetConnectionLog` | `() []LogEntry` | Last 200 log entries, newest first. Not used by the UI. |

## Settings & security

| Method | Signature | Description |
|--------|-----------|-------------|
| `GetLanguage` / `SetLanguage` | `() string` / `(lang string)` | UI language; emits `language-changed`. |
| `GetTheme` / `SetTheme` | `() string` / `(theme string)` | `dark`/`light`; emits `theme-changed`. |
| `GetAutoRotateInterval` / `SetAutoRotateInterval` | `() int` / `(seconds int)` | 0 disables; emits `config-changed`. |
| `GetHealthCheckInterval` / `SetHealthCheckInterval` | `() int` / `(seconds int)` | Min 60, default 180; emits `config-changed`. Not exposed in the UI. |
| `GetKillSwitch` | `() bool` | Live firewall state: whether the `Proxy0x-KillSwitch-AllowApp` rule exists. |
| `SetKillSwitch` | `(enabled bool) error` | Enables (outbound default block + allow rules, previous policy saved) or disables (policy restored, rules removed) through one elevated (UAC) command, waits for it, then saves the flag. Errors on a declined UAC prompt or a failed `netsh`. Disabled again on app shutdown. |
| `GetAutoFailover` / `SetAutoFailover` | `() bool` / `(enabled bool)` | When a connected proxy drops, connect the inactive proxy with the lowest known latency. |
| `EnableUWPSupport` | `() error` | Elevated `CheckNetIsolation LoopbackExempt` for Edge, Store, Search, Paint. |

## MCP integration

| Method | Signature | Description |
|--------|-----------|-------------|
| `GetMCPStatus` | `() MCPStatus` | Binary path/existence and, per client, detection and whether `proxy0x` is configured. |
| `InstallMCPClient` | `(clientID string) error` | Adds `mcpServers.proxy0x` to the client's config. Aborts on unparseable JSON; writes atomically with a `<file>.proxy0x.bak` backup. |
| `UninstallMCPClient` | `(clientID string) error` | Removes `mcpServers.proxy0x` (same safety rules). |
| `CheckMCPHealth` | `() MCPStatus` | Starts `proxy0x-mcp.exe` hidden, sends JSON-RPC `initialize`, keeps stdin open until a `result` arrives (5 s timeout). |

Client IDs and config files: `claude-desktop` (`%AppData%\Claude\claude_desktop_config.json`), `codex` (`~\.codex\config.json`), `cursor` (`~\.cursor\mcp.json`), `windsurf` (`~\.codeium\windsurf\mcp_config.json`).

## Window

| Method | Signature | Description |
|--------|-----------|-------------|
| `MinimizeWindow` | `()` | Minimize. |
| `CloseWindow` | `()` | Hide to tray. |
| `ShowWindow` | `()` | Restore from tray. Not used by the UI (the tray calls the runtime directly). |
| `ToggleMaximize` | `()` | Maximize/restore. |
| `QuitApp` | `()` | Stop local server, restore system proxy, quit. Not used by the UI. |

## Events (Go → frontend)

| Event | Payload | Emitted when |
|-------|---------|--------------|
| `proxy-updated` | `ProxyEntry` (copy, no password) | Any change to a single proxy (connect, disconnect, rotate, health, add, rename, group). |
| `proxies-changed` | — | List changed (batch add, multi-format import, remove). |
| `proxy-dropped` | `{id, host}` | Health check failed on a connected proxy. |
| `store-error` | `string` | `proxies.json` could not be loaded or saved; the store is read-only. |
| `bulk-check-progress` | `BulkCheckResult` | Each result of `BulkCheckProxies`. |
| `language-changed` | `string` | `SetLanguage`. |
| `theme-changed` | `string` | `SetTheme`. |
| `config-changed` | `Config` | `SetAutoRotateInterval`, `SetHealthCheckInterval`. |
| `log-entry` | `LogEntry` | Every log line. No listener in the UI. |

## Types

```go
type ProxyEntry struct {
    ID          string      `json:"id"`
    Name        string      `json:"name"`
    Host        string      `json:"host"`
    Port        string      `json:"port"`
    Username    string      `json:"username"`
    Password    string      `json:"password"`   // always "" in bindings/events
    Protocol    string      `json:"protocol"`   // "http" | "socks5"
    Group       string      `json:"group"`
    Status      ProxyStatus `json:"status"`     // offline | connecting | connected | error
    CurrentIP   string      `json:"currentIp"`
    Location    *GeoInfo    `json:"location"`
    Latency     int64       `json:"latency"`    // ms
    LastCheck   string      `json:"lastCheck"`  // RFC 3339
    Active      bool        `json:"active"`
    SessionID   string      `json:"sessionId"`  // 12 hex chars
    Order       int         `json:"order"`
    ConnectedAt string      `json:"connectedAt"`
    TotalUptime int64       `json:"totalUptime"` // seconds
}

type BulkCheckResult struct {
    ID          string `json:"id"`
    Host        string `json:"host"`        // host:port
    Alive       bool   `json:"alive"`       // exit IP verified (HTTP 200 + valid IP)
    Latency     int64  `json:"latency"`
    ExitIP      string `json:"exitIp"`
    Country     string `json:"country"`
    CountryCode string `json:"countryCode"`
    Anonymity   string `json:"anonymity"`   // "elite" (exit ≠ real IP) | "transparent" | "unknown"
    Error       string `json:"error,omitempty"`
}

type DNSLeakResult struct {
    ProxyIP   string `json:"proxyIp"`
    DNSIP     string `json:"dnsIp"`      // direct public IP
    DNSServer string `json:"dnsServer"`  // currently never set
    Leaked    bool   `json:"leaked"`
}

type MCPStatus struct {
    ExePath     string          `json:"exePath"`
    ExeExists   bool            `json:"exeExists"`
    Clients     []MCPClientInfo `json:"clients"`
    HealthOK    bool            `json:"healthOk"`
    HealthError string          `json:"healthError,omitempty"`
}

type MCPClientInfo struct {
    ID         string `json:"id"`
    Name       string `json:"name"`
    Detected   bool   `json:"detected"`
    Installed  bool   `json:"installed"`
    ConfigPath string `json:"configPath"`
}

type LogEntry struct {
    Time    string `json:"time"`
    Action  string `json:"action"`
    ProxyID string `json:"proxyId"`
    Host    string `json:"host"`
    Detail  string `json:"detail"`
}
```

`GeoInfo` and `Config`: see `docs/data.md`.

## Known gaps

- `DNSLeakResult.DNSServer` is declared but never filled.

## Frontend usage

```typescript
import { ConnectProxy, GetProxies } from 'wailsjs/go/main/App'
import { EventsOn } from 'wailsjs/runtime/runtime'
```
