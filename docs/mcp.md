# MCP Server

Proxy0x exposes its proxy pool to AI agents via the [Model Context Protocol](https://modelcontextprotocol.io/) (MCP). The server is a standalone Go binary that communicates over **stdio** using newline-delimited JSON-RPC 2.0.

## Architecture

```
Claude Desktop ──stdio──► proxy0x-mcp.exe ──► %AppData%/Proxy0x/proxies.json
                                              (shared with the Wails app)
```

Both the Wails desktop app and the MCP server read and write the same AES-256-GCM encrypted proxy store at `%AppData%/Proxy0x/proxies.json`. The encryption key lives at `%AppData%/Proxy0x/.key` and is created on first use.

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
| `proxy_fetch` | `url` (required), `method`, `headers`, `body` | Fetch a URL through the best proxy. SSRF-validated. Body capped at 100 KB. |
| `proxy_list` | — | List all proxies. Credentials are stripped from output. |
| `proxy_status` | — | Pool stats: total, active, connected counts. |
| `proxy_connect` | `id` | Connect a proxy by ID. Verifies IP via api.ipify.org. |
| `proxy_disconnect` | `id` | Disconnect a proxy by ID. |
| `proxy_rotate` | `id` | Generate a new session ID for the proxy, getting a different IP. |
| `proxy_add` | `proxy` (host:port[:user:pass]) | **Read-only server:** always returns an error asking to add proxies in the Proxy0x app. |
| `proxy_remove` | `id` | **Read-only server:** always returns an error asking to remove proxies in the Proxy0x app. |

The MCP server never writes `proxies.json`. The desktop app owns the file; two processes rewriting it would silently undo each other's changes. Connect, disconnect and rotate only change the MCP server's in-memory state.
| `proxy_health` | — | Health-check all active proxies. Returns per-proxy results. |

### Proxy selection

`proxy_fetch` uses `BestProxy()` which selects the active proxy with the highest score:

```
score = successes × 2 − failures
```

Each successful fetch increments `successes`; each failure increments `failures`.

### SSRF protection

Every URL passed to `proxy_fetch` is validated before the request is made:

- **Scheme:** only `http` and `https` allowed
- **Blocked hostnames:** `localhost`, `metadata`, and any name ending in `.local`, `.localhost`, `.internal`
- **Blocked IPs:** private ranges (10/8, 172.16/12, 192.168/16, 127/8), link-local (169.254/16), zeroes (0/8), IPv6 loopback (::1), ULA (fc00::/7), link-local v6 (fe80::/10)
- **Cloud metadata:** 169.254.169.254, 169.254.170.2, 100.100.100.200, fd00:ec2::254
- **DNS resolution:** hostnames are resolved and every resulting IP is checked
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
| `scrape_workflow` | `target_url`, `data_needed` | Guided web scraping workflow: connect → fetch → extract → rotate |
| `geo_check` | `target_url`, `countries` | Geographic access verification: test URL from proxies in different locations |

## One-click integration

The Wails app provides a toggle in **Settings → Claude MCP** that:

1. **`InstallMCPConfig()`** — writes the MCP server entry into `%AppData%/Claude/claude_desktop_config.json`:
   ```json
   {
     "mcpServers": {
       "proxy0x": {
         "command": "C:\\...\\proxy0x-mcp.exe",
         "args": []
       }
     }
   }
   ```
   Preserves any existing entries in the config file. If the existing file is not valid JSON the install aborts with an error and the file is left untouched. Before writing, a backup is saved as `<file>.proxy0x.bak`, and the write is atomic (temp file + rename).

2. **`CheckMCPHealth()`** — spawns the MCP binary, sends a JSON-RPC `initialize` request, and verifies a valid response within 5 seconds.

3. **`UninstallMCPConfig()`** — removes the `proxy0x` key from `mcpServers`.

After toggling on, restart Claude Desktop to pick up the new config.

## Building

```bash
go build -o proxy0x-mcp.exe ./cmd/proxy0x-mcp/
```

Place the binary next to `Proxy0x.exe` so the one-click integration can find it.

## Testing

```bash
go test ./pkg/proxycore/...
go test ./internal/mcpserver/...
```

Tests cover: SSRF validation (27 URL cases + 8 IP cases), pool CRUD, encrypted persistence, score-based selection, tool handlers (add, list, status, remove, fetch-SSRF, health), and credential stripping.
