# Backend

Go backend of the Proxy0x desktop app. Exported `App` methods are Wails bindings (full list with signatures and events: `docs/api.md`).

## Packages

| Package | Responsibility |
|---------|----------------|
| `main` (`app.go`, `tray.go`, `mcp_integration.go`, `main.go`) | `App` state, Wails bindings, lifecycle, tray, MCP client config |
| `internal/proxy` | `LocalServer` (local HTTP/CONNECT proxy) and `SystemProxy` (WinINET registry with backup/restore) |
| `internal/crypto` | AES-256-GCM credential encryption, key protected with DPAPI. Shared with the MCP server. |
| `internal/config` | `config.json` and data-file paths |
| `internal/fsutil` | `WriteFileAtomic` (temp file + rename, keeps `.bak`) |
| `internal/killswitch` | Windows Firewall rules through an elevated script |
| `internal/mcpserver`, `pkg/proxycore`, `cmd/proxy0x-mcp` | MCP server (see `docs/mcp.md`) |

## App struct

```go
type App struct {
    ctx        context.Context
    proxies    map[string]*ProxyEntry  // keyed by random hex ID; guarded by mu
    clients    map[string]*http.Client // keyed by "id:sessionId"; guarded by clientsMu
    clientsMu  sync.Mutex
    mu         sync.RWMutex
    lpMu       sync.Mutex              // guards localProxy
    localProxy *proxy.LocalServer      // the single local server
    stopCh     chan struct{}           // stops background loops
    logs       []LogEntry              // ring buffer, 200 entries
    logMu      sync.RWMutex
    crypt      *crypto.Store
    storeErr   error                   // non-nil: proxies.json unreadable → read-only

    // injected at the I/O boundaries (tests replace them)
    dataFile       string
    emit           func(ctx context.Context, name string, data ...interface{})
    sysproxy       systemProxy // Set(port) / Clear() / Recover()
    newLocalServer func(proxy.UpstreamConfig) (*proxy.LocalServer, error)
}
```

Rules the code follows:

- Network calls never run while holding `mu`.
- Bindings and events get copies (`view`/`snapshot`) with the password blanked.
- Events go through `a.emitEvent` (never `runtime.EventsEmit` directly) so the code runs outside Wails in tests.
- Exit-IP checks go through `verifyExitIP` (HTTP 200 + valid IP).

## Lifecycle

### `startup(ctx)`
1. Store the Wails context.
2. `sysproxy.Recover()`: undo a system proxy Proxy0x left behind after a crash, restoring the user's previous values.
3. Load `proxies.json`; on invalid JSON or undecryptable credentials the store becomes read-only and `store-error` is emitted.
4. Start the tray, the health-check loop and the auto-rotation loop.

### `shutdown(ctx)`
1. If the kill switch is on: remove its firewall rules (UAC prompt) and save the flag as off.
2. Stop the local server and restore the user's system proxy.
3. Close `stopCh`.

## Proxy management

- Formats: `host:port[:user:pass]`, optional `socks5://` / `http://` prefix (`AddProxy`, `AddProxiesBatch`); `AddProxiesMultiFormat` also accepts URIs with `user:pass@` and CSV (see `docs/api.md`).
- IDs are 8 random bytes (hex); session IDs 6 random bytes (hex).
- Groups are derived from the `group` field; no separate storage.
- Every mutation checks that the store is writable, saves, and returns copies.

## Connection flow

### `ConnectProxy(id)`
1. Disconnect the active proxy (only one can be active).
2. `status=connecting`; build/reuse the HTTP client for `id:session` (SOCKS5 setup errors fail the dial, never fall back to direct).
3. `verifyExitIP` → on failure `status=error`, IP/geo cleared.
4. Geolocation from `ipwho.is`.
5. `setLocalProxy`: start `LocalServer` and `sysproxy.Set(port)`; on failure `status=error`, nothing left running.
6. `status=connected`, `active=true`, `connectedAt` recorded.

### `DisconnectProxy(id)`
Accumulate uptime, reset runtime fields, close the HTTP client, stop the local server if it serves this proxy, restore the system proxy.

### `RotateProxy(id)`
Verify a candidate session through its own client; on success switch the session, close the old client and call `LocalServer.UpdateUpstream` (same port, open tunnels kept, registry untouched). On failure close the candidate client and keep the previous session/status.

Rotation appends `_session-{id}` to the username; residential providers map each session to a different exit IP.

### Auto-failover (`tryFailover`)
Called by `CheckHealth` when a connected proxy drops and `AutoFailover` is on: connects (asynchronously) the inactive proxy with the lowest known latency, or any proxy if none has a latency yet.

## Local proxy server (`internal/proxy/server.go`)

- Listens on `127.0.0.1:0`.
- CONNECT tunnels (HTTPS) and plain HTTP, upstream HTTP or SOCKS5.
- Plain HTTP: one request per client connection, forwarded with `Connection: close`; the response is relayed marked `close`.
- Upstream dial and handshake bounded by 15 s (`upstreamTimeout`).
- Every client and upstream connection is tracked; `Stop()` closes them all and returns in under a second.
- `UpdateUpstream(cfg)` switches new connections to a new upstream session.
- Accept errors back off 5 ms → 1 s.

## System proxy (`internal/proxy/system.go`)

`SystemProxy` works through a `Registry` interface (real implementation: `golang.org/x/sys/windows/registry`, key `HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings`):

- `Set(port)`: on first use saves `{previous: {enable, server, override}, ours: "127.0.0.1:<port>"}` to `sysproxy-backup.json`, then sets `ProxyEnable=1`, `ProxyServer=127.0.0.1:<port>`, `ProxyOverride=localhost;127.*;10.*;192.168.*;<local>`.
- `Clear()`: restores `previous` exactly and deletes the backup. Without a backup it only disables an old-version leftover.
- `Recover()` (startup): restores only if the registry still holds `ours`. If the user changed the settings since, the backup is discarded. Without a backup, only `127.0.0.1:*` with the exact bypass list above is disabled.
- `InternetSetOptionW` refreshes WinINET after each change.

## Kill switch (`internal/killswitch`)

`Enable()` writes a `.bat` to `%TEMP%\proxy0x-fw.bat` and runs it elevated (`ShellExecute runas`). The script adds three outbound rules: `Proxy0x-KillSwitch-Block` (block all), `-AllowApp` (allow Proxy0x.exe) and `-AllowLocal` (allow 127.0.0.1/::1). `Disable()` deletes them. See `docs/security.md` for the open issues with this approach.

## Background loops

- **Health check**: ticks every 30 s; checks active proxies every `HealthCheckInterval` (default 180 s), 1.5 s apart.
- **Auto-rotation**: ticks every 10 s; rotates active proxies every `AutoRotateInterval` seconds (0 = off).

## Persistence

- `proxies.json`: credentials encrypted on save, runtime fields (`status`, `active`, `currentIp`, `latency`) stripped. Written atomically with `.bak`. Encryption failure aborts the save (never writes plaintext).
- `config.json`: written atomically.
- Shapes and key format: `docs/data.md`.

## MCP client integration (`mcp_integration.go`)

Supports Claude Desktop, Codex, Cursor and Windsurf (paths in `docs/api.md`). Writes only `mcpServers.proxy0x`, refuses to touch an unparseable file, backs it up as `<file>.proxy0x.bak` and writes atomically. `CheckMCPHealth` starts `proxy0x-mcp.exe` without a console window and waits up to 5 s for the `initialize` result.

## System tray (`tray.go`)

Show · Disconnect (active proxy) · Rotate (active proxy) · Language submenu (16 languages) · Quit.

## External services

| Service | Purpose | Used by |
|---------|---------|---------|
| `api.ipify.org` | Exit IP (through proxy) and direct IP | Connect, rotate, health, benchmark, bulk check, connect fastest, DNS leak, `GetMyIP` |
| `ipwho.is` | Geolocation | Connect, rotate, health (on IP change), bulk check |
| `speed.cloudflare.com` | 100 KB download | Benchmark |
