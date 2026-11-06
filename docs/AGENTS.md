# Proxy0x — Agent Instructions

## Overview

Proxy0x is a Windows desktop proxy manager built with **Wails v2** (Go 1.26 backend + React 18 frontend). It manages HTTP and SOCKS5 proxies, routes browser traffic through a local proxy server, and includes an MCP server for AI agent integration.

## Commands

```bash
# Development (hot reload)
cd frontend && npm install && cd ..
wails dev

# Production build
wails build -ldflags "-s -w"
# Output: build/bin/Proxy0x.exe

# Build MCP server
go build -o build/bin/proxy0x-mcp.exe ./cmd/proxy0x-mcp/

# Run tests
go test ./pkg/proxycore/... ./internal/mcpserver/...

# Lint
go vet ./...

# Frontend type-check
cd frontend && npx tsc --noEmit
```

## Module Layout

```
proxy_ui/
├── main.go, app.go, tray.go, mcp_integration.go   # Wails app (main package)
├── internal/
│   ├── proxy/      # Local proxy server, Windows system proxy (registry)
│   ├── crypto/     # AES-GCM credential encryption
│   ├── config/     # App configuration (load, save, defaults)
│   └── mcpserver/  # MCP server: tools, resources, prompts
├── pkg/proxycore/  # Reusable proxy logic (pool, store, SSRF, fetcher, client)
├── cmd/proxy0x-mcp/  # MCP server binary entry point (stdio transport)
└── frontend/src/
    ├── components/  # React components (ProxyCard, SettingsPage, etc.)
    ├── i18n/        # Translations (16 languages) and LanguageContext
    └── App.tsx      # Root component
```

## Key Conventions

### Wails Bindings

Exported methods on the `App` struct in the main package are automatically exposed to the frontend as TypeScript functions. After adding or changing a method signature:

- Bindings are regenerated automatically during `wails dev`.
- Generated files: `frontend/wailsjs/go/main/App.js`, `App.d.ts`, and `models.ts`.
- Import in React: `import { MethodName } from 'wailsjs/go/main/App'`.

### Internationalization (i18n)

Every user-facing string must be translated in **all 16 languages** in `frontend/src/i18n/translations.ts`:

1. Add the key to the `Translations` interface.
2. Add the value to all 16 language objects (en, pt-BR, pt-PT, es, fr, de, it, ru, ja, ko, zh-CN, ar, hi, tr, th, pl).

### Security Rules

- **Never expose proxy credentials** in MCP tool output. The `proxy_list` tool strips `Username` and `Password` before returning.
- **SSRF protection** is mandatory for `proxy_fetch`. All URLs must pass `ValidateURL()` which blocks private IPs, metadata endpoints, and non-HTTP schemes. Redirect hops are also validated.
- **Credential encryption**: proxy credentials are encrypted with AES-256-GCM at rest. The encryption is transparent — `Store.Save()` encrypts, `Store.Load()` decrypts.

### Testing

Test suites exist in:
- `pkg/proxycore/` — pool, store, SSRF validation (35 tests)
- `internal/mcpserver/` — MCP tool handlers, credential stripping (9 tests)

Run with: `go test ./pkg/proxycore/... ./internal/mcpserver/...`

Race detector requires CGO (`CGO_ENABLED=1`); tests run fine without it.

## Architecture Notes

- The Wails app and the MCP server are **separate binaries** that share the same encrypted proxy file (`%AppData%/Proxy0x/proxies.json`).
- The MCP server communicates via **stdio** (newline-delimited JSON-RPC 2.0) and is launched by Claude Desktop as a child process.
- The Wails app can auto-configure Claude Desktop's MCP config via the Settings page toggle.
- System proxy is set via Windows registry (`HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings`).
