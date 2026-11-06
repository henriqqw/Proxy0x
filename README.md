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

Windows desktop proxy manager. Route all browser traffic through HTTP or SOCKS5 proxies via a local proxy server, rotate IPs by session, run speed tests, organize proxies by group, and monitor latency and geolocation in real time. Kill switch blocks all traffic if a proxy drops. Auto-failover connects to the next fastest proxy automatically. DPAPI-protected credential encryption. DNS leak test, bulk proxy checker, UWP/Store app support, and multi-format import. Includes an **MCP server** for AI agent integration with Claude Desktop, Cursor, Codex and Windsurf. 16 languages, frameless dark UI with system tray. Built with **Go**, **Wails** and **React**.

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
| `Proxy0x-Setup-1.0.0.exe` | Windows installer — language selector, install location, startup option |
| `Proxy0x.exe` | Standalone portable executable |

Download from the [latest release](https://github.com/henriqqw/Proxy0x/releases/latest).

### Installer features

- 14 language options for the installer UI
- Custom install location
- Desktop shortcut (optional)
- Start on Windows boot (optional)
- Full uninstaller via Add/Remove Programs

## Features

| Area | Details |
|------|---------|
| Protocols | HTTP and SOCKS5 proxy support |
| System proxy | Local proxy server routes all browser traffic through the selected proxy |
| Proxy management | Add single, batch import, multi-format import (URI, CSV, `host:port:user:pass`), drag-to-reorder, remove, import/export JSON |
| Proxy naming | User-defined labels for proxies — clickable name in card header and expanded details |
| Proxy groups | Tag proxies with groups, filter by group |
| Connection | Connect, disconnect (instant), rotate IP via session ID |
| Connect fastest | Auto-connect to the lowest-latency proxy in the pool |
| Auto rotation | Configurable interval with presets (1m, 3m, 5m, 10m) or custom minutes |
| Auto-failover | On proxy drop, automatically connect to the next best proxy by latency |
| Kill switch | Block all outbound traffic if proxy drops (Windows Firewall rules with UAC elevation) |
| DNS leak test | Compare proxy IP vs direct IP to detect DNS leaks |
| Bulk checker | Test all proxies concurrently for latency, country, and anonymity level |
| UWP support | Enable proxy for Microsoft Edge and Store apps (loopback exemption) |
| Speed test | Measure connect latency (ms) and download speed (Mbps) per proxy |
| Monitoring | Current IP, latency, live uptime timer, connection status badge, copy IP to clipboard |
| Geolocation | Country flag, city, region, ISP, timezone, AS number, organization |
| Health check | Automatic health checks every 3 minutes, Windows notification on proxy drop |
| Connection stats | Connected-at time, total uptime per proxy (persisted across sessions) |
| Encryption | AES-GCM credential encryption with DPAPI-protected key (tied to Windows user profile) |
| Persistence | Proxies saved to `%APPDATA%/Proxy0x/proxies.json`, auto-restored on startup |
| MCP server | AI agent integration — 9 tools, 3 resources, 2 prompts, SSRF protection |
| MCP clients | One-click setup for Claude Desktop, Cursor, Codex and Windsurf |
| i18n | 16 languages — EN, PT-BR, PT-PT, ES, FR, DE, IT, RU, JA, KO, ZH-CN, AR, HI, TR, TH, PL |
| System tray | Minimize to tray, language submenu, connect/disconnect/rotate all, quit |
| Interface | Frameless window, drag-to-reorder, loading skeletons, mouse back/forward navigation |

## Proxy formats

```
host:port:user:pass                    # Standard format
socks5://user:pass@host:port           # URI format (http://, socks5://)
host,port,user,pass,socks5             # CSV format
```

All three formats are supported in the multi-format import (Settings → Import). Residential proxies with geographic targeting (DataImpulse, etc.) are supported. IP rotation works by appending `_session-{id}` to the username — each unique session ID gets a different IP from the pool.

## Usage

1. Paste a proxy in `host:port:user:pass`, URI (`socks5://user:pass@host:port`), or CSV format — or use **Batch Import** for multiple at once.
2. Click **Connect** — the app starts a local proxy server and routes your browser traffic through it. Or use **Connect Fastest** in Settings → Tools to auto-connect to the lowest-latency proxy.
3. See the proxy's real IP, country flag, city, ISP and latency. Click the IP to copy it to clipboard.
4. Give proxies custom **names** by clicking the name field in the card header or expanded details.
5. Use **Rotate** to get a new IP, or enable **Auto Rotation** in settings with a preset or custom interval.
6. Enable **Kill Switch** in Settings → Security to block all traffic if the proxy drops.
7. Enable **Auto-Failover** to automatically connect to the next best proxy on drop.
8. Run **DNS Leak Test** or **Bulk Proxy Check** in Settings → Tools to verify security and performance.
9. Enable **UWP Support** for Microsoft Edge and Store apps.
10. Run **Speed Test** to measure connect latency and download speed.
11. Assign proxies to **groups** for filtering (click the group label in the expanded card).
12. Change language via the globe button in the title bar.
13. Close the window — the app minimizes to system tray. Click the tray icon to restore; use **Quit** to exit.

## MCP Integration

Proxy0x includes an MCP (Model Context Protocol) server that lets AI agents use your proxy pool. The MCP server runs as a separate binary (`proxy0x-mcp.exe`) communicating via stdio.

### One-click setup

Open **Settings** in the app and toggle any MCP client — **Claude Desktop**, **Codex**, **Cursor** or **Windsurf**. This writes the config automatically and runs a health check to verify it works. Restart the client after connecting.

### Manual setup

1. Build the MCP server:
   ```bash
   go build -o proxy0x-mcp.exe ./cmd/proxy0x-mcp/
   ```

2. Add to `%APPDATA%/Claude/claude_desktop_config.json`:
   ```json
   {
     "mcpServers": {
       "proxy0x": {
         "command": "C:\\path\\to\\proxy0x-mcp.exe",
         "args": []
       }
     }
   }
   ```

3. Restart Claude Desktop.

### Available tools

| Tool | Description |
|------|-------------|
| `proxy_fetch` | Fetch a URL through the best proxy (SSRF-protected) |
| `proxy_list` | List all proxies (credentials stripped) |
| `proxy_status` | Pool stats overview |
| `proxy_connect` | Connect a proxy by ID |
| `proxy_disconnect` | Disconnect a proxy |
| `proxy_rotate` | Rotate session for a new IP |
| `proxy_add` | Add a proxy (`host:port[:user:pass]`) |
| `proxy_remove` | Remove a proxy by ID |
| `proxy_health` | Health check all active proxies |

See [docs/mcp.md](docs/mcp.md) for full MCP documentation.

## How it works

Proxy0x runs a local HTTP/SOCKS5 proxy server on `127.0.0.1`. When you connect a proxy, the app:

1. Starts a local listener that accepts browser connections
2. Sets the Windows system proxy to point to this local server
3. Forwards all traffic through the selected upstream proxy
4. Tracks all connections for instant disconnect (no hanging)

```
Browser ──► Local Proxy (127.0.0.1) ──► Upstream Proxy ──► Internet
              │                              │
              │ HTTP CONNECT / SOCKS5        │ Residential IP
              │ Connection tracking          │ Session rotation
              ▼                              ▼
         Windows System Proxy          api.ipify.org (IP)
         (auto-configured)             ip-api.com (geo)
```

## Architecture

```
┌──────────────────────────────────────────────────┐
│              Wails WebView2                      │
│  ┌────────────────────────────────────────────┐  │
│  │    React + TypeScript + Tailwind           │  │
│  │    Lucide · @dnd-kit · i18n (16 lang)      │  │
│  └──────────────────┬─────────────────────────┘  │
│                     │ bindings                    │
│  ┌──────────────────▼─────────────────────────┐  │
│  │           Go Backend                       │  │
│  │  Proxy Manager · Local Proxy Server        │  │
│  │  HTTP & SOCKS5 · Health Checker            │  │
│  │  Speed Test · Auto Rotation · Failover     │  │
│  │  Kill Switch · DNS Leak · Bulk Checker     │  │
│  │  DPAPI + AES-GCM Crypto · Config Store     │  │
│  └──────────────────┬─────────────────────────┘  │
│                     │                             │
│  ┌──────────────────▼─────────────────────────┐  │
│  │          System Tray                       │  │
│  │    energye/systray · Language menu         │  │
│  └────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────┘
         │              │              │
         ▼              ▼              ▼
   api.ipify.org   ip-api.com    Windows System
   (IP lookup)    (geolocation)  Proxy / Firewall

┌──────────────────────────────────────────────────┐
│            MCP Server (stdio)                    │
│  proxy0x-mcp.exe — separate binary               │
│  9 tools · 3 resources · 2 prompts              │
│  SSRF protection · Score-based rotation          │
│  Shares encrypted proxy store with the app       │
└──────────────────────────────────────────────────┘
         │
         ▼
   Claude Desktop · Cursor · Codex · Windsurf
```

## Stack

| Layer | Technology | Role |
|-------|------------|------|
| Desktop | Go 1.26, Wails 2.16 | Frameless window, bindings, events, system tray |
| Frontend | React 18, TypeScript 5.6, Tailwind CSS 3 | Components, theme system, responsive layout |
| Drag & drop | @dnd-kit | Sortable proxy list with grip handles |
| Icons | Lucide React | Consistent monochrome iconography |
| Flags | flagcdn.com | Country flag images with lazy loading |
| Tray | energye/systray | System tray icon with language and control submenus |
| Crypto | DPAPI + AES-GCM | DPAPI-protected key (crypt32.dll), AES-GCM credential encryption at rest |
| Network | net/http, net (stdlib) | HTTP & SOCKS5 proxy, connection pooling, health check |
| MCP | mcp-go (mark3labs) | Model Context Protocol server for AI agent integration |
| Installer | Inno Setup 6 | Windows installer with i18n and startup option |

## Development

### Prerequisites

- Go 1.26+
- Node.js 18+ and npm
- Wails CLI v2 (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

### Dev mode with hot reload

```bash
cd frontend && npm install && cd ..
wails dev
```

Frontend updates instantly via Vite HMR. Go changes recompile and restart the app automatically.

### Production build

```bash
wails build -ldflags "-s -w"
```

Binary output: `build/bin/Proxy0x.exe`

### MCP server build

```bash
go build -o proxy0x-mcp.exe ./cmd/proxy0x-mcp/
```

### Tests

```bash
go test ./pkg/proxycore/... ./internal/mcpserver/...
```

### Installer build

Requires [Inno Setup 6](https://jrsoftware.org/isinfo.php).

```bash
iscc build/installer.iss
```

Generates `build/installer/Proxy0x-Setup-1.0.0.exe`.

## Project structure

```
Proxy0x/
├── main.go                      # Entry point, Wails config
├── app.go                       # App struct, all Wails-bound methods
├── tray.go                      # System tray with language submenu
├── mcp_integration.go           # MCP client config install/uninstall/health
├── internal/
│   ├── killswitch/
│   │   └── killswitch.go        # Windows Firewall kill switch (UAC elevation)
│   ├── proxy/
│   │   ├── server.go            # Local HTTP/SOCKS5 proxy server
│   │   └── system.go            # Windows system proxy (registry)
│   ├── crypto/
│   │   └── crypto.go            # DPAPI + AES-GCM credential encryption
│   ├── config/
│   │   └── config.go            # App configuration persistence
│   └── mcpserver/
│       ├── server.go            # MCP server factory
│       ├── tools.go             # 9 MCP tools
│       ├── resources.go         # 3 MCP resources
│       └── prompts.go           # 2 MCP prompts
├── pkg/proxycore/
│   ├── types.go                 # Shared proxy types
│   ├── pool.go                  # Thread-safe proxy pool
│   ├── client.go                # HTTP/SOCKS5 client builder
│   ├── store.go                 # Encrypted JSON persistence
│   ├── ssrf.go                  # SSRF URL validation
│   └── fetcher.go               # Proxy fetch with SSRF protection
├── cmd/proxy0x-mcp/
│   └── main.go                  # MCP server entry point (stdio)
├── docs/                        # Project documentation
├── wails.json                   # Wails project config
├── build/
│   ├── bin/                     # Production binary
│   ├── installer/               # Installer output
│   ├── installer.iss            # Inno Setup installer script
│   └── windows/                 # Icons, manifest, tray icon
├── assets/
│   ├── banner.png               # README banner
│   └── screenshots/             # README screenshots
└── frontend/
    ├── src/
    │   ├── App.tsx              # Root component, routing, DnD context
    │   ├── i18n/
    │   │   ├── translations.ts  # 16 language sets (83 keys each)
    │   │   └── LanguageContext.tsx
    │   └── components/
    │       ├── TitleBar.tsx      # Window controls, theme, language
    │       ├── ProxyCard.tsx     # Sortable proxy card with naming + copy IP
    │       ├── ProxyForm.tsx     # Input and batch import
    │       ├── SettingsPage.tsx  # Settings, security, tools, MCP integration
    │       ├── ContextMenu.tsx   # Right-click language selector
    │       ├── StatusBadge.tsx   # Status indicator
    │       ├── MyIPBadge.tsx     # Real IP display
    │       ├── FlagIcon.tsx      # Country flag
    │       └── EmptyState.tsx    # Empty state
    ├── wailsjs/                  # Auto-generated Go↔TS bindings
    └── tailwind.config.js        # Theme tokens
```

## Documentation

| Document | Description |
|----------|-------------|
| [Architecture](docs/architecture.md) | System architecture, data flow, module layout |
| [Backend](docs/backend.md) | Go backend: proxy manager, local server, health checks |
| [Frontend](docs/frontend.md) | React components, i18n, Tailwind theme |
| [API Reference](docs/api.md) | All Wails-exposed Go methods |
| [MCP Server](docs/mcp.md) | MCP tools, resources, prompts, SSRF protection |
| [Business Rules](docs/business-rules.md) | Scoring, rotation, health check, session logic |
| [Data Models](docs/data.md) | Types, persistence, file formats |
| [Security](docs/security.md) | Encryption, SSRF, credential handling |
| [Design System](docs/design-system.md) | Color tokens, typography, component patterns |
| [Market Research](docs/market-research.md) | Proxy market analysis and competitor research |
| [AGENTS.md](docs/AGENTS.md) | Instructions for AI agents working with this codebase |
| [CLAUDE.md](docs/CLAUDE.md) | Claude Code project configuration |

## External services

| Service | Usage | Limit |
|---------|-------|-------|
| [api.ipify.org](https://www.ipify.org/) | Discover current IP through proxy | No documented limit |
| [ip-api.com](http://ip-api.com/) | IP geolocation (country, city, ISP, AS) | 45 req/min (free tier) |
| [ipwho.is](https://ipwho.is/) | Geolocation for MCP server | No documented limit |
| [flagcdn.com](https://flagcdn.com/) | Country flag images | CDN, no limit |

## Requirements

- Windows 10/11
- WebView2 Runtime (pre-installed on Windows 10 21H2+ and Windows 11)
