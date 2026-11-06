<p align="center">
  <img src="assets/banner.png" alt="Proxy0x — Proxy Manager" width="600">
</p>

<p align="center">
  <a href="https://github.com/henriqqw/Proxy0x/releases"><img src="https://img.shields.io/github/v/release/henriqqw/Proxy0x?style=flat-square&color=68E045" alt="Release"></a>
  <a href="https://github.com/henriqqw/Proxy0x/releases"><img src="https://img.shields.io/github/downloads/henriqqw/Proxy0x/total?style=flat-square&color=4b5563" alt="Downloads"></a>
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.26-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go 1.26"></a>
  <a href="https://wails.io/"><img src="https://img.shields.io/badge/Wails-2.16-DF0000?style=flat-square" alt="Wails 2.16"></a>
  <a href="https://react.dev/"><img src="https://img.shields.io/badge/React-18-149ECA?style=flat-square&logo=react&logoColor=white" alt="React 18"></a>
</p>
<p align="center">
  <a href="https://github.com/henriqqw/Proxy0x/releases/latest">Download</a>
  · <a href="#mcp-integration">MCP Integration</a>
  · <a href="https://github.com/henriqqw/Proxy0x/issues">Report issue</a>
</p>

Windows desktop proxy manager. Routes the traffic of applications that use the Windows system proxy (browsers included) through HTTP or SOCKS5 proxies via a local proxy server, rotates IPs by session, runs speed tests, organizes proxies by group and monitors latency and geolocation. Optional kill switch, auto-failover, DPAPI-protected credential encryption, exit-IP check, bulk proxy checker, UWP/Store app support and multi-format import. Includes a read-only **MCP server** for AI agents (Claude Desktop, Cursor, Codex, Windsurf). 16 languages, frameless dark UI with system tray. Built with **Go**, **Wails** and **React**.

## Screenshots

<p align="center">
  <img src="assets/screenshots/connected.png" alt="Connected proxy with geolocation" width="200">
  <img src="assets/screenshots/empty.png" alt="Empty state" width="200">
  <img src="assets/screenshots/settings.png" alt="Settings — language, rotation, import/export" width="200">
  <img src="assets/screenshots/tools-mcp.png" alt="Security, tools and MCP integration" width="200">
</p>

## Downloads

| File | Description |
|------|-------------|
| `Proxy0x-Setup-<version>.exe` | Windows installer (recommended) — includes the app and the MCP server |
| `Proxy0x.exe` | Portable desktop app |
| `proxy0x-mcp.exe` | MCP server binary (place it next to `Proxy0x.exe` for the one-click MCP setup) |

Download from the [latest release](https://github.com/henriqqw/Proxy0x/releases/latest). Windows x64 only.

### Installer features

- 14 language options for the installer UI
- Per-user install by default (no admin needed), custom install location
- Installs `Proxy0x.exe` and `proxy0x-mcp.exe`
- Desktop shortcut (optional), start with Windows (optional)
- Uninstaller via Add/Remove Programs; it asks whether to also delete saved proxies and settings (`AppData\Roaming\Proxy0x`)

## Features

| Area | Details |
|------|---------|
| Protocols | HTTP and SOCKS5 upstream proxies |
| System proxy | Local proxy server on `127.0.0.1`; the Windows (WinINET) system proxy points at it. Your previous proxy settings are backed up and restored on disconnect, and recovered after a crash |
| Proxy management | Add single, batch import, multi-format import (URI, CSV, `host:port:user:pass`), drag-to-reorder, remove, export as `host:port:user:pass` text |
| Proxy naming | User-defined labels shown on the card |
| Proxy groups | Tag proxies with a group, filter by group |
| Connection | One active proxy at a time. Connect, instant disconnect, rotate IP via session ID (rotation keeps the local port and open connections) |
| Connect fastest | Tests every inactive proxy and connects the lowest-latency one |
| Auto rotation | Presets (1, 3, 5, 10 min) or a custom interval in minutes |
| Auto-failover | When the active proxy fails a health check, connects the inactive proxy with the lowest last-known latency |
| Kill switch | While enabled, Windows Firewall blocks outbound traffic by default and allows only Proxy0x, loopback and DNS (requires UAC). Your previous firewall policy is restored on disable, on app exit and on uninstall. See [Known limitations](#known-limitations) |
| Exit IP check ("DNS leak test") | Compares the exit IP through the proxy with your direct IP |
| Bulk checker | Tests all proxies concurrently (10 at a time): latency, exit IP, country, transparent/elite |
| UWP support | Loopback exemption so Edge/Store apps can reach the local proxy (requires UAC) |
| Speed test | Connect latency (ms) and download speed (Mbps, 100 KB sample) |
| Monitoring | Exit IP (copy to clipboard), latency, status badge, connected-at time, total uptime |
| Geolocation | Country flag, city, region, ISP, timezone, AS number, organization |
| Health check | Every 3 minutes for the active proxy; Windows notification when it drops; status recovers when it answers again |
| Exit IP verification | Connect, rotate, health check and speed test only succeed with HTTP 200 and a valid IP |
| Encryption | Credentials encrypted with AES-256-GCM; the key is protected with Windows DPAPI (tied to your Windows user) and shared by the app and the MCP server |
| Persistence | `%APPDATA%/Proxy0x/proxies.json`, written atomically (previous version kept as `.bak`). If the file can't be read, the store becomes read-only instead of being overwritten |
| MCP server | Read-only access for AI agents — 9 tools, 3 resources, 2 prompts, SSRF protection |
| MCP clients | One-click setup for Claude Desktop, Cursor, Codex and Windsurf (backs up the client's config, never overwrites an unparseable one) |
| i18n | 16 languages — EN, PT-BR, PT-PT, ES, FR, DE, IT, RU, JA, KO, ZH-CN, AR, HI, TR, TH, PL |
| System tray | Show, Disconnect / Rotate the active proxy, language submenu, Quit |
| Interface | Frameless window, drag-to-reorder, loading skeletons, mouse back/forward navigation |

## Proxy formats

```
host:port                              # No auth
host:port:user:pass                    # Standard format
socks5://host:port:user:pass           # Explicit protocol (http:// also accepted)
socks5://user:pass@host:port           # URI format (multi-format import)
host,port,user,pass,socks5             # CSV format (multi-format import)
```

The single-proxy field and batch import accept the first three forms; **Settings → Import** also accepts URI and CSV lines. Residential proxies with geographic targeting (DataImpulse, etc.) are supported: IP rotation appends `_session-{id}` to the username, and each session ID maps to a different IP from the provider's pool.

## Usage

1. Paste a proxy (`host:port:user:pass`) or use **Batch Import** for many at once.
2. Click **Connect** — the app verifies the exit IP, starts the local proxy server and points the Windows system proxy at it. Or use **Connect Fastest** in Settings → Tools.
3. See the exit IP, country flag, city, ISP and latency. Click the IP to copy it.
4. Name a proxy by clicking its name on the card; assign a **group** in the expanded card.
5. Use **Rotate** for a new IP, or enable **Auto Rotation** in Settings.
6. Optional: enable **Kill Switch** and **Auto-Failover** in Settings → Security.
7. Run the **exit IP check**, **Bulk Check** or **Speed Test** in Settings → Tools.
8. Enable **UWP Support** if Edge/Store apps should use the proxy.
9. Close the window — the app keeps running in the system tray. Use **Quit** in the tray to exit (this restores your previous system proxy settings).

## MCP Integration

Proxy0x ships an MCP (Model Context Protocol) server, `proxy0x-mcp.exe`, that lets AI agents fetch URLs through your proxies. It runs as a separate process over stdio and **only reads** `proxies.json`: proxies are added and removed in the desktop app.

### One-click setup

Open **Settings** and toggle an MCP client — **Claude Desktop**, **Codex**, **Cursor** or **Windsurf**. Proxy0x adds a `proxy0x` entry to the client's config (keeping a `.proxy0x.bak` backup and refusing to touch a config it can't parse), then runs a health check. Restart the client afterwards. Requires `proxy0x-mcp.exe` next to `Proxy0x.exe` (the installer does this).

### Manual setup

Add to `%APPDATA%/Claude/claude_desktop_config.json` (see `cmd/proxy0x-mcp/mcp-config.example.json`):

```json
{
  "mcpServers": {
    "proxy0x": {
      "command": "C:\\Path\\To\\Proxy0x\\proxy0x-mcp.exe",
      "args": []
    }
  }
}
```

Then restart Claude Desktop.

### Available tools

| Tool | Description |
|------|-------------|
| `proxy_fetch` | Fetch a URL through the active proxy, or a given `proxy_id` (SSRF-protected, body capped at 100 KB) |
| `proxy_list` | List proxies (credentials stripped) |
| `proxy_status` | Pool stats |
| `proxy_connect` | Connect a proxy by ID (MCP server's own in-memory state) |
| `proxy_disconnect` | Disconnect a proxy |
| `proxy_rotate` | New session ID for a new IP |
| `proxy_health` | Health check of active proxies |
| `proxy_add` / `proxy_remove` | Not available: the server is read-only and returns an error pointing to the app |

See [docs/mcp.md](docs/mcp.md) for the full reference.

## How it works

```
App (WinINET proxy aware) ──► Local proxy 127.0.0.1:<port> ──► Upstream proxy ──► Internet
                                │                                │
                                │ HTTP CONNECT tunnels           │ Residential IP
                                │ plain HTTP (Connection: close) │ Session rotation
                                ▼                                ▼
                         Windows system proxy             api.ipify.org (exit IP)
                         (backed up + restored)           ipwho.is (geolocation)
```

When you connect a proxy, Proxy0x:

1. Verifies the exit IP through the upstream proxy (HTTP 200 + valid IP)
2. Starts a local listener on `127.0.0.1` (random port)
3. Saves your current Windows proxy settings, then points the system proxy at the local listener
4. Forwards CONNECT tunnels and plain HTTP requests to the upstream (HTTP or SOCKS5), with 15 s dial/handshake timeouts
5. Tracks every connection, so Disconnect returns in under a second and restores your previous proxy settings

Rotation swaps the upstream session in place: same local port, open connections untouched.

## Architecture

```
┌──────────────────────────────────────────────────┐
│              Wails WebView2                      │
│  ┌────────────────────────────────────────────┐  │
│  │    React + TypeScript + Tailwind           │  │
│  │    Lucide · @dnd-kit · i18n (16 lang)      │  │
│  └──────────────────┬─────────────────────────┘  │
│                     │ bindings (copies, no passwords)
│  ┌──────────────────▼─────────────────────────┐  │
│  │           Go backend (app.go)              │  │
│  │  Proxy manager · Health check · Rotation   │  │
│  │  Failover · Bulk check · Speed test        │  │
│  │  internal/proxy   local server, system proxy│ │
│  │  internal/crypto  DPAPI + AES-256-GCM      │  │
│  │  internal/killswitch  Windows Firewall     │  │
│  │  internal/config · internal/fsutil         │  │
│  └──────────────────┬─────────────────────────┘  │
│                     │                             │
│  ┌──────────────────▼─────────────────────────┐  │
│  │          System tray (energye/systray)     │  │
│  └────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────┘
         │              │              │
         ▼              ▼              ▼
   api.ipify.org    ipwho.is     Windows registry
   (exit IP)      (geolocation)  / Firewall

┌──────────────────────────────────────────────────┐
│            MCP server (stdio)                    │
│  proxy0x-mcp.exe — separate binary, read-only    │
│  9 tools · 3 resources · 2 prompts · SSRF guard  │
│  pkg/proxycore, same crypto as the app           │
└──────────────────────────────────────────────────┘
         │
         ▼
   Claude Desktop · Cursor · Codex · Windsurf
```

## Stack

| Layer | Technology | Role |
|-------|------------|------|
| Desktop | Go 1.26 (toolchain go1.26.10), Wails 2.16 | Frameless window, bindings, events, system tray |
| Frontend | React 18, TypeScript 5.6, Tailwind CSS 3, Vite 6 | Components, theme, layout |
| Drag & drop | @dnd-kit | Sortable proxy list |
| Icons | Lucide React | Iconography |
| Flags | country-flag-icons | Country flags bundled with the app (no remote requests) |
| Tray | energye/systray | Tray icon and menu |
| Crypto | AES-256-GCM + Windows DPAPI (`golang.org/x/sys/windows`) | Credential encryption at rest |
| Network | `net`, `net/http`, `golang.org/x/net/proxy` | Local proxy, HTTP/SOCKS5 clients |
| Registry | `golang.org/x/sys/windows/registry` | System proxy backup/restore |
| MCP | mcp-go (mark3labs) | MCP server |
| Installer | Inno Setup 6 | Windows installer |
| CI | GitHub Actions | Tag-triggered release build |

## Development

### Prerequisites

- Go 1.26+ (the `go.mod` toolchain line downloads go1.26.10 automatically)
- Node.js 18+ and npm
- Wails CLI v2.16 (`go install github.com/wailsapp/wails/v2/cmd/wails@v2.16.0`)

### Dev mode with hot reload

```bash
cd frontend && npm install && cd ..
wails dev
```

### Production build

```bash
wails build -ldflags "-s -w"
go build -ldflags "-s -w" -o build/bin/proxy0x-mcp.exe ./cmd/proxy0x-mcp
```

Output: `build/bin/Proxy0x.exe` and `build/bin/proxy0x-mcp.exe`.

### Tests and checks

```bash
go vet ./...
go test ./...
cd frontend && npx tsc --noEmit
```

Packages with tests: root (`app`), `internal/proxy`, `internal/crypto`, `internal/fsutil`, `internal/mcpserver`, `pkg/proxycore`. The root package embeds `frontend/dist`, so run `wails build` (or `npm run build` in `frontend/`) once before `go test ./...`. The race detector needs cgo (`CGO_ENABLED=1` and a C compiler).

### Installer build

Requires [Inno Setup 6](https://jrsoftware.org/isinfo.php) and both binaries in `build/bin/`.

```bash
iscc /DMyAppVersion=1.2.3 build\installer.iss
```

Generates `build/installer/Proxy0x-Setup-1.2.3.exe`. Without `/DMyAppVersion` the version is `0.0.0-dev`.

### Releasing

Releases are published only by pushing a version tag:

```bash
git tag v1.2.3
git push origin v1.2.3
```

The `Release` workflow builds the app and the MCP server, runs `go vet` and `go test`, builds the installer and publishes `Proxy0x-Setup-<version>.exe`, `Proxy0x.exe` and `proxy0x-mcp.exe`. Pushing to `main` publishes nothing. Running the workflow manually (`workflow_dispatch`) builds the artifacts without publishing.

## Project structure

```
Proxy0x/
├── main.go                      # Entry point, Wails config
├── app.go                       # App struct and all Wails-bound methods
├── tray.go                      # System tray
├── mcp_integration.go           # MCP client config install/uninstall/health check
├── internal/
│   ├── proxy/
│   │   ├── server.go            # Local HTTP/SOCKS5 proxy server, connection tracking
│   │   └── system.go            # Windows system proxy: backup, set, restore, crash recovery
│   ├── crypto/crypto.go         # DPAPI-protected key + AES-256-GCM
│   ├── config/config.go         # App configuration (atomic writes)
│   ├── fsutil/fsutil.go         # Atomic file writes with .bak
│   ├── killswitch/killswitch.go # Windows Firewall kill switch, UAC elevation
│   └── mcpserver/               # MCP tools, resources, prompts
├── pkg/proxycore/               # MCP-side proxy logic: pool, store, client, SSRF, fetcher
├── cmd/proxy0x-mcp/             # MCP server entry point (stdio) + example config
├── build/
│   ├── installer.iss            # Inno Setup script
│   ├── bin/                     # Build output (ignored)
│   ├── installer/               # Installer output
│   └── windows/                 # Icons, manifest
├── .github/workflows/release.yml  # Tag-triggered release
├── docs/                        # Project documentation (+ docs/audit/: audit reports)
├── assets/                      # Banner, screenshots
└── frontend/
    ├── src/
    │   ├── App.tsx              # Root component, routing, DnD context
    │   ├── i18n/                # translations.ts (16 languages), LanguageContext
    │   └── components/          # TitleBar, ProxyCard, ProxyForm, SettingsPage,
    │                            # StatusBadge, MyIPBadge, FlagIcon, EmptyState
    ├── wailsjs/                 # Auto-generated Go↔TS bindings
    └── tailwind.config.js       # Theme tokens
```

## Documentation

| Document | Description |
|----------|-------------|
| [Architecture](docs/architecture.md) | System architecture, data flow, module layout |
| [Backend](docs/backend.md) | Go backend: proxy manager, local server, health checks |
| [Frontend](docs/frontend.md) | React components, i18n, Tailwind theme |
| [API Reference](docs/api.md) | Wails-exposed Go methods |
| [MCP Server](docs/mcp.md) | MCP tools, resources, prompts, SSRF protection |
| [Business Rules](docs/business-rules.md) | Connection lifecycle, rotation, health check, persistence |
| [Data Models](docs/data.md) | Types, persistence, file formats |
| [Security](docs/security.md) | Encryption, SSRF, system proxy, credential handling |
| [Design System](docs/design-system.md) | Color tokens, typography, component patterns |
| [Audit](docs/audit/) | Baseline, findings, remediation plan and change log |
| [Market Research](docs/market-research.md) | Proxy market analysis |
| [AGENTS.md](docs/AGENTS.md) | Instructions for AI agents |
| [CLAUDE.md](docs/CLAUDE.md) | Claude Code project configuration |

## External services

| Service | Usage | Seen IP |
|---------|-------|---------|
| [api.ipify.org](https://www.ipify.org/) | Exit IP through the proxy (connect, rotate, health check, speed test) | Proxy IP |
| [api.ipify.org](https://www.ipify.org/) | Direct IP for "My IP" and the exit-IP check | Real IP (by design) |
| [ipwho.is](https://ipwho.is/) | Geolocation for the app and the MCP server, sent through the proxy | Proxy IP |
| [speed.cloudflare.com](https://speed.cloudflare.com/) | Speed test download sample, through the proxy | Proxy IP |

Country flags and fonts are bundled or system-provided (no remote UI requests; production builds enforce this with a Content-Security-Policy). All requests use HTTPS.

## Known limitations

- **System proxy scope:** only applications that honor the Windows (WinINET) system proxy are routed. UDP traffic (e.g. WebRTC, QUIC) and apps with their own network stack are not covered and can reveal your real IP.
- **Kill switch scope:** programs that already have their own outbound allow rule in Windows Firewall (Windows services, Store apps, some games and launchers) keep network access, and DNS lookups are allowed. If Proxy0x crashes while the switch is on, the block stays until you turn it off in the app (fail closed). Details: [docs/security.md](docs/security.md#kill-switch-internalkillswitch).
- **"DNS leak test"** compares the exit IP through the proxy with your direct IP. It is an exit-IP check, not a DNS resolver test.
- **Single active proxy:** connecting one proxy disconnects the previous one.
- Windows x64 only.

## Requirements

- Windows 10/11 (x64)
- WebView2 Runtime (pre-installed on Windows 10 21H2+ and Windows 11)
