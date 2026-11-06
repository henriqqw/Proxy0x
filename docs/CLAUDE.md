# CLAUDE.md

## Project

Proxy0x — Windows desktop proxy manager.  
Stack: Go 1.25 (backend) + React 18 / TypeScript 5.6 / Tailwind CSS 3 (frontend), connected via Wails v2 bindings.

## Build & Test

```bash
wails dev                    # Dev mode with hot reload
wails build -ldflags "-s -w" # Production build → build/bin/Proxy0x.exe
go test ./pkg/proxycore/... ./internal/mcpserver/...  # Run tests
go vet ./...                 # Lint
cd frontend && npx tsc --noEmit  # Type-check frontend
```

## Directory Layout

| Path | Purpose |
|------|---------|
| `*.go` (root) | Wails-bound App struct: proxy management, tray, MCP integration |
| `internal/proxy/` | Local HTTP/SOCKS5 proxy server, Windows system proxy (registry) |
| `internal/crypto/` | AES-256-GCM credential encryption |
| `internal/config/` | App configuration (language, theme, intervals) |
| `internal/mcpserver/` | MCP server: 9 tools, 3 resources, 2 prompts |
| `pkg/proxycore/` | Standalone proxy logic: pool, store, SSRF guard, fetcher, client |
| `cmd/proxy0x-mcp/` | MCP server binary (stdio transport) |
| `frontend/src/components/` | React UI components |
| `frontend/src/i18n/` | 16-language translation system |
| `frontend/wailsjs/` | Auto-generated Go↔TypeScript bindings |

## Conventions

- **i18n**: every UI string goes in `frontend/src/i18n/translations.ts`, in ALL 16 language blocks.
- **Wails bindings**: exported `App` methods auto-generate TypeScript in `frontend/wailsjs/go/main/App.*`.
- **Security**: credentials encrypted at rest (AES-GCM), SSRF validation on proxy_fetch, never expose passwords in MCP output.
- **Style**: dark-first Tailwind theme, all colors as CSS custom properties on `:root`. Frameless window, 420×680px.
- **Tests**: `pkg/proxycore/` and `internal/mcpserver/` have test suites. Race detector needs `CGO_ENABLED=1`.

## Data Files

All under `%AppData%` (Windows `os.UserConfigDir()`):

- `Proxy0x/proxies.json` — encrypted proxy store
- `Proxy0x/.key` — AES-256 encryption key
- `Proxy0x/config.json` — app settings
- `Claude/claude_desktop_config.json` — MCP integration config
