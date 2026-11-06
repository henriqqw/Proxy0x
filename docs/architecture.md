# Architecture

## Overview

Proxy0x is a Windows desktop proxy manager built with [Wails v2](https://wails.io/). The Go backend manages proxy connections, the Windows system proxy and credential encryption; a React frontend runs inside a WebView2 window. A separate MCP server binary lets AI agents use the active proxy through the Model Context Protocol, reading (never writing) the app's proxy store.

## System diagram

```
┌──────────────────────────────────────────────────────────────┐
│                       Wails WebView2                         │
│  React 18 + TypeScript + Tailwind · i18n (16 languages)      │
│                 │ bindings (copies, no password) · events    │
│  ┌──────────────▼───────────────────────────────────────┐    │
│  │                Go backend (package main)             │    │
│  │  App: CRUD · connect/rotate · health · failover      │    │
│  │  bulk check · DNS leak · MCP client config · tray    │    │
│  └──┬──────────┬──────────┬──────────┬──────────┬───────┘    │
│  ┌──▼─────┐ ┌──▼─────┐ ┌──▼─────┐ ┌──▼─────┐ ┌──▼────────┐   │
│  │internal│ │internal│ │internal│ │internal│ │ internal  │   │
│  │/proxy  │ │/crypto │ │/config │ │/fsutil │ │/killswitch│   │
│  │LocalSrv│ │DPAPI + │ │config  │ │atomic  │ │firewall   │   │
│  │SysProxy│ │AES-GCM │ │.json   │ │writes  │ │(UAC)      │   │
│  └───┬────┘ └───┬────┘ └────────┘ └────────┘ └───────────┘   │
└──────┼──────────┼────────────────────────────────────────────┘
       │          │
       ▼          ▼
 HKCU Internet   %AppData%/Proxy0x  ◄──── read-only ────┐
 Settings        proxies.json, .key                     │
 (+ backup file) config.json, sysproxy-backup.json      │
                                                        │
  Separate binary (cmd/proxy0x-mcp, stdio)              │
┌───────────────────────────────────────────┐           │
│ internal/mcpserver  9 tools · 3 res · 2 pr│           │
│ pkg/proxycore  Pool · Store · SSRF · Fetch├───────────┘
│   Store uses internal/crypto (same key)   │
└───────────────────────────────────────────┘
```

## Module layout

| Path | Contents |
|------|----------|
| `main.go` | Wails entry point (420×680, frameless, transparent) |
| `app.go` | `App` struct, bindings, connection lifecycle, health/rotation loops, tools |
| `tray.go` | Tray: show, disconnect/rotate the active proxy, language, quit |
| `mcp_integration.go` | MCP entry in Claude Desktop / Codex / Cursor / Windsurf configs; health check |
| `internal/proxy/server.go` | `LocalServer`: local proxy on `127.0.0.1:0`, connection tracking, `UpdateUpstream` |
| `internal/proxy/system.go` | `SystemProxy`: WinINET registry with backup/restore and crash recovery |
| `internal/crypto` | AES-256-GCM, key protected with DPAPI; used by the app **and** the MCP store |
| `internal/config` | `Config` and data file paths |
| `internal/fsutil` | `WriteFileAtomic` (temp + rename, `.bak`) |
| `internal/killswitch` | Windows Firewall rules via elevated script |
| `pkg/proxycore` | MCP-side pool (single active proxy), read-only store loader, SSRF guard, fetcher, HTTP client |
| `internal/mcpserver` | MCP tools/resources/prompts (`proxy_add`/`proxy_remove` refuse: read-only) |
| `cmd/proxy0x-mcp` | MCP binary; loads `%AppData%/Proxy0x` and serves stdio |

## Data flows

### Connect

```
UI ConnectProxy(id)
  → disconnect current active proxy
  → verifyExitIP through upstream (HTTP 200 + valid IP)   [no lock held]
  → ipwho.is geolocation                                  [no lock held]
  → LocalServer on 127.0.0.1:<port>
  → SystemProxy.Set(port): back up user's settings, point WinINET at 127.0.0.1:<port>
  → status=connected → event proxy-updated (copy, no password)

Browser ──► 127.0.0.1:<port> ──► upstream proxy (HTTP CONNECT / SOCKS5, user_session-<id>) ──► site
```

### Rotate (manual or auto)

```
candidate session → verifyExitIP via candidate client
  ok   → switch session, LocalServer.UpdateUpstream(cfg)   (same port, open tunnels kept, registry untouched)
  fail → keep previous session and status
```

### Disconnect / shutdown / crash recovery

```
Disconnect: LocalServer.Stop() (<1s, all tracked conns closed) → SystemProxy.Clear() (restore user's values)
Shutdown:   kill switch off (if on) → same as disconnect
Startup:    SystemProxy.Recover() restores only if the registry still holds Proxy0x's address
```

### Health and failover

```
every HealthCheckInterval (default 180s): verifyExitIP for the active proxy
  fail on a connected proxy → status=error, event proxy-dropped, tryFailover (if enabled)
  ok                        → status=connected, refresh latency / IP / geo
```

### MCP

```
MCP client → proxy0x-mcp.exe (stdio) → tool call
  proxy_fetch → Pool.ActiveProxy() → SSRF validation (fail closed) → request via proxy → JSON-RPC result
  proxy_add / proxy_remove → error: read-only
```

### Shared state

Only the desktop app writes `proxies.json`. The MCP server loads it at startup and keeps connect/rotate state in memory. Both use `internal/crypto`, so they share one key file and one format, and the MCP server can never replace the app's DPAPI-protected key.

## External services

| Service | Purpose |
|---------|---------|
| `api.ipify.org` | Exit IP (through proxy) and direct IP |
| `ipwho.is` | IP geolocation |
| `speed.cloudflare.com` | Speed test (100 KB) |
| `flagcdn.com` | Country flags (frontend) |
| `fonts.googleapis.com` | Inter font (frontend) |

## Build and release

```bash
cd frontend && npm install && cd ..
wails dev                                            # hot reload
wails build -ldflags "-s -w"                         # build/bin/Proxy0x.exe
go build -o build/bin/proxy0x-mcp.exe ./cmd/proxy0x-mcp
iscc /DMyAppVersion=1.2.3 build\installer.iss        # build/installer/Proxy0x-Setup-1.2.3.exe
```

The installer script uses `SourceDir=..`, so its paths are relative to the repository root. It ships `Proxy0x.exe` and `proxy0x-mcp.exe`.

**Release:** `git tag vX.Y.Z && git push origin vX.Y.Z`. `.github/workflows/release.yml` builds the app and the MCP server, runs vet and tests, builds the installer and publishes `Proxy0x-Setup-X.Y.Z.exe`, `Proxy0x.exe` and `proxy0x-mcp.exe`. Pushing to `main` publishes nothing; `workflow_dispatch` builds the artifacts without releasing.

## Testing

```bash
go vet ./...
go test ./...                    # needs frontend/dist (run a wails/npm build first)
cd frontend && npx tsc --noEmit
```

Tests run without the Wails runtime, the real registry, the network or `%AppData%`: `App` exposes injection points (`emit`, `sysproxy`, `newLocalServer`, `dataFile`, URL variables) and `SystemProxy` takes a `Registry` interface. The race detector needs `CGO_ENABLED=1`.

## Security

See `docs/security.md`: DPAPI-protected key, atomic writes and the read-only store, system proxy backup, SSRF rules, credential stripping, and the kill switch caveats.
