# MCP Server

Proxy0x exposes its proxy pool to AI agents via the [Model Context Protocol](https://modelcontextprotocol.io/) (MCP). The server is a standalone Go binary that communicates over **stdio** using newline-delimited JSON-RPC 2.0.

## Architecture

```
MCP client ──stdio──► proxy0x-mcp.exe ──reads──► %AppData%/Proxy0x/proxies.json
(Claude Desktop,                                 (written only by the desktop app)
 Codex, Cursor,
 Windsurf)
```

The MCP server **reads** the desktop app's encrypted store at `%AppData%/Proxy0x/proxies.json` once at startup and never writes it. Decryption goes through `internal/crypto` — the same code the app uses — so it reads the DPAPI-protected `%AppData%/Proxy0x/.key` and never creates or rewrites it. Proxies added in the app after the MCP server started are not visible until the MCP client restarts the server.

### Package structure

| Package | Path | Role |
|---------|------|------|
| `proxycore` | `pkg/proxycore/` | Standalone proxy library — types, pool, store, client, fetcher, SSRF |
| `mcpserver` | `internal/mcpserver/` | MCP server factory, tool handlers, resources, prompts |
| `main` | `cmd/proxy0x-mcp/` | Entry point — loads pool from store, runs `server.ServeStdio()` |

**Dependency:** [`github.com/mark3labs/mcp-go`](https://github.com/mark3labs/mcp-go) — Go SDK for MCP server/client.

## Tools (9)

| Tool | Parameters | Description |
|------|------------|-------------|
| `proxy_fetch` | `url` (required), `method`, `body`, `proxy_id` | Fetch a URL through the active proxy (or `proxy_id`). SSRF-validated. Body capped at 100 KB. |
| `proxy_list` | — | List all proxies. Credentials are stripped from output. |
| `proxy_status` | — | Pool stats: total, active, connected counts. |
| `proxy_connect` | `id` | Connect a proxy by ID. Verifies IP via api.ipify.org. |
| `proxy_disconnect` | `id` | Disconnect a proxy by ID. |
| `proxy_rotate` | `id` | Generate a new session ID for the proxy, getting a different IP. |
| `proxy_add` | `proxy` (host:port[:user:pass]) | **Read-only server:** always returns an error asking to add proxies in the Proxy0x app. |
| `proxy_remove` | `id` | **Read-only server:** always returns an error asking to remove proxies in the Proxy0x app. |
| `proxy_health` | — | Health-check all active proxies. Returns per-proxy results. |

The MCP server never writes `proxies.json`. The desktop app owns the file; two processes rewriting it would silently undo each other's changes. Connect, disconnect and rotate only change the MCP server's in-memory state (they do not affect the desktop app's connection or the Windows system proxy).

`proxy_list` and `proxy://pool` strip `username` and `password` and no longer include a `score` field.

### Proxy selection

Only one proxy can be active at a time (connecting one disconnects the other). `proxy_fetch` uses `ActiveProxy()` — the connected proxy — unless `proxy_id` is given. There is no scoring.

### SSRF protection

Every URL passed to `proxy_fetch` is validated before the request is made (full rules in `docs/security.md`):

- **Scheme:** only `http` and `https` allowed
- **Blocked hostnames:** `localhost`, `metadata`, and any name ending in `.local`, `.localhost`, `.internal` (lowercased, trailing dot removed)
- **Blocked IPs:** private, loopback, link-local, CGNAT, benchmarking, multicast/reserved, NAT64, 6to4, IPv6 ULA/link-local, cloud metadata addresses
- **DNS resolution:** every resolved IP is checked; a failed lookup rejects the URL (fail closed). It does not prevent DNS rebinding, because the upstream proxy resolves the name again.
- **Redirect hops:** up to 10 redirects, each hop re-validated

## Resources (3)

| URI | MIME | Description |
|-----|------|-------------|
| `proxy://status` | `application/json` | Pool stats (same as `proxy_status` tool) |
| `proxy://pool` | `application/json` | Full proxy list (credentials stripped) |
| `proxy://logs` | `application/json` | Connection log (last 200 entries, newest first) |

## Prompts (2)

| Prompt | Arguments | Description |
|--------|-----------|-------------|
| `scrape_workflow` | `target_url` (required), `data_type` | Guided web scraping workflow: connect → fetch → extract → rotate |
| `geo_check` | `url` (required), `country` | Geographic access verification: test the URL from a proxy in the given country |

## One-click integration

**Settings → MCP Integration** shows one card per supported client. A client is "detected" when its folder exists:

| Client (`id`) | Config file written | Detected by |
|---------------|--------------------|-------------|
| Claude Desktop (`claude-desktop`) | `%AppData%\Claude\claude_desktop_config.json` | `%AppData%\Claude` |
| Codex (`codex`) | `~\.codex\config.json` | `~\.codex` |
| Cursor (`cursor`) | `~\.cursor\mcp.json` | `%AppData%\Cursor` or `~\.cursor` |
| Windsurf (`windsurf`) | `~\.codeium\windsurf\mcp_config.json` | `%AppData%\Windsurf` or `~\.codeium\windsurf` |

Hypothesis (not verified): Codex reads its MCP servers from TOML (`~\.codex\config.toml`), so the JSON entry written for Codex may have no effect.

Bindings:

1. **`GetMCPStatus()`** — path of `proxy0x-mcp.exe` (next to `Proxy0x.exe`), whether it exists, and each client's detected/installed state.
2. **`InstallMCPClient(id)`** — adds `mcpServers.proxy0x = { "command": "<path>\\proxy0x-mcp.exe", "args": [] }`, preserving every other key. If the existing file is not valid JSON the install aborts and the file is left untouched. Before writing, a backup is saved as `<file>.proxy0x.bak`, and the write is atomic (temp file + rename).
3. **`UninstallMCPClient(id)`** — removes `mcpServers.proxy0x`, with the same backup/atomic write.
4. **`CheckMCPHealth()`** — starts `proxy0x-mcp.exe` without a console window, sends a JSON-RPC `initialize` request, keeps stdin open until a reply with `result` arrives, and fails after 5 seconds.

Restart the MCP client after installing so it picks up the new config.

## Building

```bash
go build -ldflags "-s -w" -o build/bin/proxy0x-mcp.exe ./cmd/proxy0x-mcp
```

The installer and the GitHub release ship `proxy0x-mcp.exe` next to `Proxy0x.exe`, where the integration looks for it.

## Testing

```bash
go test ./pkg/proxycore/... ./internal/mcpserver/...
```

Tests cover: SSRF validation (scheme, hostnames, IP ranges, fail-closed DNS with an injected resolver), pool operations, the store reading the app's DPAPI key without replacing it, the pool never writing `proxies.json`, `ActiveProxy` selection, tool handlers (read-only add/remove, list, status, fetch-SSRF, health) and credential stripping. The integration code (`mcp_integration.go`) has tests in the root package for config preservation, refusal to overwrite invalid JSON and the health check.
