# Architecture

## Overview

Proxy0x is a Windows desktop proxy manager built with [Wails v2](https://wails.io/). The Go backend manages proxy connections, system proxy settings, and encryption, while a React frontend provides the UI inside a WebView2 window. A separate MCP server binary enables AI agents to use the proxy pool via the Model Context Protocol.

## System Diagram

```
┌──────────────────────────────────────────────────────────┐
│                    Wails WebView2                        │
│  ┌────────────────────────────────────────────────────┐  │
│  │         React 18 + TypeScript + Tailwind           │  │
│  │   TitleBar · ProxyCard · ProxyForm · Settings      │  │
│  │   @dnd-kit · Lucide · i18n (16 languages)          │  │
│  └──────────────────┬─────────────────────────────────┘  │
│                     │ Wails bindings (auto-generated)    │
│  ┌──────────────────▼─────────────────────────────────┐  │
│  │               Go Backend (main pkg)                │  │
│  │  App struct · Proxy CRUD · Connect/Disconnect      │  │
│  │  Rotate · Health Check · Auto-Rotation · Benchmark │  │
│  │  MCP Config · Logging · Persistence                │  │
│  └──┬──────────────┬──────────────┬───────────────────┘  │
│     │              │              │                       │
│  ┌──▼───────┐ ┌────▼────┐ ┌──────▼──────┐               │
│  │ internal │ │internal │ │  internal   │               │
│  │ /proxy   │ │/crypto  │ │  /config    │               │
│  │ LocalSrv │ │AES-GCM  │ │  JSON cfg   │               │
│  │ SysProxy │ │         │ │            │               │
│  └──────────┘ └─────────┘ └────────────┘               │
│                                                          │
│  ┌────────────────────────────────────────────────────┐  │
│  │            System Tray (energye/systray)           │  │
│  │     Language menu · Connect/Disconnect All · Quit  │  │
│  └────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────┘
          │              │              │
          ▼              ▼              ▼
    api.ipify.org   ipwho.is      Windows Registry
    (IP lookup)    (geolocation)  (system proxy)


  Separate binary (cmd/proxy0x-mcp):
┌─────────────────────────────────────────────────┐
│           proxy0x-mcp.exe (stdio)               │
│  ┌───────────────────────────────────────────┐  │
│  │         internal/mcpserver                │  │
│  │  9 tools · 3 resources · 2 prompts        │  │
│  └───────────────┬───────────────────────────┘  │
│  ┌───────────────▼───────────────────────────┐  │
│  │          pkg/proxycore                    │  │
│  │  Pool · Store · Client · SSRF · Fetcher   │  │
│  └───────────────────────────────────────────┘  │
└─────────────────────────────────────────────────┘
          │
          ▼
    %APPDATA%/Proxy0x/proxies.json
    (shared encrypted proxy store)
```

## Module Layout

### Root Package (`main`)

| File | Purpose |
|------|---------|
| `main.go` | Wails entry point, window config (420×680, frameless, transparent) |
| `app.go` | `App` struct, all Wails-bound methods, proxy lifecycle, connection logic |
| `tray.go` | System tray with language submenu, global proxy controls |
| `mcp_integration.go` | Claude Desktop MCP config install/uninstall/health check |

### `internal/proxy`

| File | Purpose |
|------|---------|
| `server.go` | `LocalServer` — local HTTP/SOCKS5 proxy on `127.0.0.1:0`, CONNECT tunneling, connection tracking |
| `system.go` | Windows system proxy via registry (`HKCU\...\Internet Settings`), `wininet.dll` refresh |

### `internal/crypto`

| File | Purpose |
|------|---------|
| `crypto.go` | `Store` — AES-256-GCM encrypt/decrypt, key at `%APPDATA%/Proxy0x/.key`, `enc:` prefix convention |

### `internal/config`

| File | Purpose |
|------|---------|
| `config.go` | `Config` struct (language, theme, intervals), JSON persistence at `%APPDATA%/Proxy0x/config.json` |

### `pkg/proxycore` (shared library)

Standalone proxy logic used by the MCP server binary. No Wails dependency.

| File | Purpose |
|------|---------|
| `types.go` | `ProxyEntry`, `PoolStats`, `FetchResult`, `LogEntry` — shared types with JSON tags |
| `pool.go` | Thread-safe `Pool` — CRUD, connect/disconnect/rotate, health check, score-based selection, logging |
| `store.go` | JSON file persistence with AES-GCM encryption (same scheme as `internal/crypto`) |
| `client.go` | `BuildProxyURL`, `BuildHTTPClient` — HTTP and SOCKS5 client construction with session injection |
| `ssrf.go` | SSRF validation — blocks private CIDRs, metadata IPs, suspicious hostnames |
| `fetcher.go` | `Fetcher` — URL fetch through proxy with SSRF validation on target and each redirect hop |

### `internal/mcpserver`

MCP server implementation using the `mcp-go` SDK.

| File | Purpose |
|------|---------|
| `server.go` | `MCPServer` factory — wraps Pool + Fetcher + mcp-go server instance |
| `tools.go` | 9 MCP tools: proxy_fetch, proxy_list, proxy_status, proxy_connect, proxy_disconnect, proxy_rotate, proxy_add, proxy_remove, proxy_health |
| `resources.go` | 3 MCP resources: `proxy://status`, `proxy://pool`, `proxy://logs` |
| `prompts.go` | 2 MCP prompts: `scrape_workflow`, `geo_check` |

### `cmd/proxy0x-mcp`

Entry point for the standalone MCP server binary. Reads `%APPDATA%/Proxy0x`, loads the shared proxy store, and runs `server.ServeStdio()`.

## Data Flow

### UI → Backend

```
React component
  → calls Wails binding (e.g. ConnectProxy)
    → Go method on App struct
      → modifies proxies map, creates HTTP client
      → starts local proxy server
      → sets Windows system proxy via registry
      → emits Wails event (proxy-updated)
        → React state update via EventsOn listener
```

### MCP Server → Proxy Pool

```
Claude Desktop
  → spawns proxy0x-mcp.exe (stdio)
    → JSON-RPC initialize handshake
    → tool call (e.g. proxy_fetch)
      → MCPServer handler
        → Pool.BestProxy() (score-based selection)
        → SSRF validation on target URL
        → HTTP request through proxy
        → response back via JSON-RPC
```

### Shared State

Both binaries read/write the same proxy store at `%APPDATA%/Proxy0x/proxies.json` using the same AES-GCM encryption scheme and key file. They do not run simultaneously on the store — the Wails app manages the live pool, while the MCP server loads it on startup.

## External Services

| Service | Protocol | Purpose | Rate Limit |
|---------|----------|---------|------------|
| `api.ipify.org` | HTTPS | Exit IP detection | None documented |
| `ipwho.is` | HTTPS | IP geolocation (country, city, ISP, ASN) | 10K/month free |
| `speed.cloudflare.com` | HTTPS | Speed benchmark (100KB download) | None documented |
| `flagcdn.com` | HTTPS | Country flag images | CDN, unlimited |

## Build Pipeline

### Development

```bash
cd frontend && npm install && cd ..
wails dev
```

Vite HMR for frontend changes. Go recompiles on save.

### Production

```bash
wails build -ldflags "-s -w"
```

Output: `build/bin/Proxy0x.exe`

### MCP Server

```bash
go build -o proxy0x-mcp.exe ./cmd/proxy0x-mcp/
```

Place alongside `Proxy0x.exe` for the Settings → MCP integration button to find it.

### Installer

Requires [Inno Setup 6](https://jrsoftware.org/isinfo.php):

```bash
iscc installer.iss
```

Output: `build/installer/Proxy0x-Setup-1.0.0.exe`

Features: 14-language installer UI, custom install path, desktop shortcut, start-on-boot, uninstaller.

## Testing

```bash
go test ./pkg/proxycore/...       # Pool, Store, SSRF tests (10 test functions)
go test ./internal/mcpserver/...  # MCP server integration tests (9 test functions)
go vet ./...                      # Static analysis
cd frontend && npx tsc --noEmit   # TypeScript type checking
```

Race detector requires `CGO_ENABLED=1` (Windows default is CGO off).

## Security Considerations

- Credentials encrypted at rest with AES-256-GCM (per-machine key)
- SSRF protection blocks private CIDRs, metadata endpoints, localhost, .local/.internal hostnames
- Redirect hops validated against SSRF rules (max 10 hops)
- MCP proxy_list tool strips credentials from output
- System proxy bypass includes localhost and private ranges
- Body size limited to 100KB on proxy_fetch to prevent memory exhaustion
