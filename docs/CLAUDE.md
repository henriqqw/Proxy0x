# CLAUDE.md

## Project

Proxy0x — Windows desktop proxy manager.  
Stack: Go 1.26 (toolchain go1.26.10) + React 18 / TypeScript 5.6 / Tailwind CSS 3, connected via Wails v2 bindings. MCP server in a separate binary.

## Build & Test

```bash
wails dev                                    # Dev mode with hot reload
wails build -ldflags "-s -w"                 # build/bin/Proxy0x.exe
go build -o build/bin/proxy0x-mcp.exe ./cmd/proxy0x-mcp
go vet ./...                                 # Lint
go test ./...                                # Tests (needs frontend/dist from a wails/npm build)
cd frontend && npx tsc --noEmit              # Type-check frontend
iscc /DMyAppVersion=1.2.3 build\installer.iss  # Installer
```

Release: `git tag vX.Y.Z && git push origin vX.Y.Z` (CI publishes installer, portable exe and MCP exe). Pushing `main` publishes nothing.

## Directory Layout

| Path | Purpose |
|------|---------|
| `*.go` (root) | Wails-bound App struct: proxy management, tray, MCP client integration |
| `internal/proxy/` | Local HTTP/SOCKS5 proxy server; Windows system proxy with backup/restore |
| `internal/crypto/` | DPAPI-protected key + AES-256-GCM (single implementation, used by app and MCP) |
| `internal/config/` | App configuration |
| `internal/fsutil/` | Atomic file writes (temp + rename, `.bak`) |
| `internal/killswitch/` | Windows Firewall kill switch (UAC) |
| `internal/mcpserver/` | MCP server: 9 tools, 3 resources, 2 prompts |
| `pkg/proxycore/` | MCP-side pool, read-only store, SSRF guard, fetcher, client |
| `cmd/proxy0x-mcp/` | MCP server binary (stdio transport) |
| `build/installer.iss` | Inno Setup script |
| `.github/workflows/release.yml` | Tag-triggered release |
| `frontend/src/` | React components and 16-language i18n |
| `frontend/wailsjs/` | Auto-generated Go↔TypeScript bindings |
| `docs/`, `docs/audit/` | Specs (source of truth) and audit reports |

## Conventions

- **Specs first**: update `docs/` before or with behavior changes.
- **i18n**: every UI string goes in `frontend/src/i18n/translations.ts`, in ALL 16 language blocks.
- **Backend invariants** (see `docs/AGENTS.md`): `a.emitEvent` instead of `runtime.EventsEmit`; return `a.view(p)` copies (no password) to the UI; no network I/O under `a.mu`; `verifyExitIP` for exit-IP checks; `checkStoreWritable` + `fsutil.WriteFileAtomic` for persistence; system proxy only via `a.sysproxy`.
- **Security**: credentials encrypted at rest; never expose passwords to the webview or MCP output; SSRF validation on `proxy_fetch`; the MCP server never writes `proxies.json`.
- **Tests**: fix bugs test-first; tests must not touch the real registry or `%APPDATA%` (use the App seams). Race detector needs `CGO_ENABLED=1`.
- **Style**: dark-first Tailwind theme, colors as CSS custom properties on `:root`. Frameless window, 420×680px.

## Data Files

All under `%AppData%` (Windows `os.UserConfigDir()`):

- `Proxy0x/proxies.json` (+ `.bak`) — encrypted proxy store, written only by the app
- `Proxy0x/.key` — AES-256 key, protected with DPAPI (`.key.legacy` = pre-DPAPI backup)
- `Proxy0x/config.json` — app settings (language, theme, intervals, kill switch, failover)
- `Proxy0x/sysproxy-backup.json` — user's Windows proxy settings while connected
- `Proxy0x/killswitch-backup.json` — firewall policy per profile while the kill switch is on
- MCP client configs (e.g. `Claude/claude_desktop_config.json`) — backed up as `<file>.proxy0x.bak`
