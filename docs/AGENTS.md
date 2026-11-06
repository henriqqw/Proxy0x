# Proxy0x — Agent Instructions

## Overview

Proxy0x is a Windows desktop proxy manager built with **Wails v2** (Go 1.26 backend + React 18 frontend). It manages HTTP and SOCKS5 proxies, routes traffic of apps that honor the Windows system proxy through a local proxy server, and ships a read-only MCP server for AI agents.

## Commands

```bash
# Development (hot reload)
cd frontend && npm install && cd ..
wails dev

# Production build
wails build -ldflags "-s -w"                                         # build/bin/Proxy0x.exe
go build -ldflags "-s -w" -o build/bin/proxy0x-mcp.exe ./cmd/proxy0x-mcp

# Checks (go test needs frontend/dist: run wails build or `npm run build` in frontend first)
go vet ./...
go test ./...
cd frontend && npx tsc --noEmit

# Installer (Inno Setup 6, both binaries in build/bin)
iscc /DMyAppVersion=1.2.3 build\installer.iss

# Release: push a version tag; pushing main publishes nothing
git tag v1.2.3 && git push origin v1.2.3
```

## Module Layout

```
proxy_ui/
├── main.go, app.go, tray.go, mcp_integration.go   # Wails app (main package)
├── internal/
│   ├── proxy/       # Local proxy server; Windows system proxy (backup/restore via Registry interface)
│   ├── crypto/      # DPAPI-protected key + AES-256-GCM (shared by app and MCP)
│   ├── config/      # App configuration (atomic writes)
│   ├── fsutil/      # WriteFileAtomic (temp + rename, .bak)
│   ├── killswitch/  # Windows Firewall rules via elevated cmd
│   └── mcpserver/   # MCP tools, resources, prompts
├── pkg/proxycore/   # MCP-side pool, store (uses internal/crypto), SSRF guard, fetcher, client
├── cmd/proxy0x-mcp/ # MCP server binary (stdio)
├── build/installer.iss            # Inno Setup script (SourceDir=.., version via /DMyAppVersion)
├── .github/workflows/release.yml  # Tag-triggered release
├── docs/ (+ docs/audit/)          # Specs and audit reports
└── frontend/src/                  # components/, i18n/, App.tsx
```

## Key Conventions

### Specs first

`docs/` is the source of truth (business rules, security, backend, MCP). Update the spec before or together with a behavior change.

### Wails bindings

Exported methods on `App` are exposed to the frontend. Bindings regenerate during `wails dev`/`wails build` (`frontend/wailsjs/go/main/App.*`, `models.ts`). Import in React: `import { MethodName } from 'wailsjs/go/main/App'`.

### Backend invariants (enforced by tests in `app_test.go`)

- **Events:** use `a.emitEvent(...)` / `a.emitProxy(p)`, never `runtime.EventsEmit` directly (crashes outside Wails, untestable).
- **UI copies:** bindings and events return `a.view(p)` / `snapshot(p)` copies taken under `a.mu`; the password is blanked. Never hand the live `*ProxyEntry` to Wails.
- **Locks:** read `ProxyEntry` fields under `a.mu`; never do network I/O while holding `a.mu`. Use `a.copiesOf(...)` to get unlocked copies for network checks. The HTTP client map has its own `clientsMu`; go through `getOrCreateClient` / `closeClient`.
- **Exit IP:** use `verifyExitIP(client)` (HTTP 200 + valid IP). Do not add new `client.Get(ipify)` copies.
- **Local proxy:** `setLocalProxy` returns an error (connect fails cleanly); rotation uses `updateLocalProxy` (in place, same port).
- **Persistence:** mutating methods call `checkStoreWritable()` first and return `saveProxies()` errors. Writes go through `fsutil.WriteFileAtomic`. If `proxies.json` can't be loaded the store is read-only.
- **System proxy:** only through `a.sysproxy` (`Set`/`Clear`/`Recover`), which backs up and restores the user's settings.
- **Seams for tests:** `a.dataFile`, `a.emit`, `a.sysproxy`, `a.newLocalServer`, and package vars `ipCheckURL`, `geoURLFormat`, `myIPURL`, `newSOCKS5`. Tests must never touch the real registry or `%APPDATA%`.

### Internationalization

Every user-facing string goes in `frontend/src/i18n/translations.ts`, in **all 16 languages** (en, pt-BR, pt-PT, es, fr, de, it, ru, ja, ko, zh-CN, ar, hi, tr, th, pl): add the key to the `Translations` interface, then to every language object (`tsc` enforces completeness).

### Security rules

- **Never expose credentials:** MCP `proxy_list` / `proxy://pool` strip username and password; the desktop UI never receives the password.
- **SSRF:** `proxy_fetch` URLs must pass `ValidateURL()` (fails closed on DNS errors, trailing-dot normalization, private/reserved/NAT64/6to4 ranges, redirect hops re-validated). It does not prevent DNS rebinding through the upstream proxy.
- **Encryption:** `internal/crypto` is the only implementation; `pkg/proxycore` delegates to it. Never write a second key format.
- **MCP is read-only:** it must never save `proxies.json`.

### Testing

`go test ./...` — packages with tests: root (`app_test.go`, `mcp_integration_test.go`), `internal/proxy`, `internal/crypto`, `internal/fsutil`, `internal/mcpserver`, `pkg/proxycore`. Fix bugs test-first. The race detector needs cgo (`CGO_ENABLED=1` + a C compiler).

## Architecture Notes

- The Wails app and the MCP server are **separate binaries** sharing `%AppData%/Proxy0x/proxies.json` and `.key`; only the app writes `proxies.json`.
- The MCP server speaks newline-delimited JSON-RPC 2.0 over stdio and is launched by the MCP client.
- The app can configure Claude Desktop, Codex, Cursor and Windsurf from Settings; client configs are backed up as `<file>.proxy0x.bak` and never overwritten when unparseable.
- System proxy: `HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings`, previous values saved in `%AppData%/Proxy0x/sysproxy-backup.json`.
