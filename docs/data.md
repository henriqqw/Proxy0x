# Data Models & Persistence

## Models

### ProxyEntry (`app.go`)

| Field | Type | JSON | Notes |
|-------|------|------|-------|
| ID | string | `id` | 8 random bytes, hex |
| Name | string | `name` | Custom display name (optional) |
| Host | string | `host` | |
| Port | string | `port` | |
| Username | string | `username` | Encrypted at rest |
| Password | string | `password` | Encrypted at rest; always `""` in bindings and events |
| Protocol | string | `protocol` | `http` or `socks5` |
| Group | string | `group` | Optional tag |
| Status | ProxyStatus | `status` | `offline` / `connecting` / `connected` / `error` (runtime) |
| CurrentIP | string | `currentIp` | Exit IP (runtime) |
| Location | *GeoInfo | `location` | From `ipwho.is` |
| Latency | int64 | `latency` | ms (runtime) |
| LastCheck | string | `lastCheck` | RFC 3339 |
| Active | bool | `active` | The single connected proxy (runtime) |
| SessionID | string | `sessionId` | 6 random bytes, hex; appended as `_session-<id>` to the username |
| Order | int | `order` | Display order |
| ConnectedAt | string | `connectedAt` | RFC 3339 |
| TotalUptime | int64 | `totalUptime` | Seconds, accumulated on disconnect |

Runtime fields (`status`, `active`, `currentIp`, `latency`) are reset on save and on load.

### GeoInfo

`country`, `countryCode`, `region`, `city`, `isp`, `org`, `lat`, `lon`, `timezone`, `as` (e.g. `AS13335`). Mapped from the `ipwho.is` response.

### Config (`internal/config`)

| Field | JSON | Default |
|-------|------|---------|
| Language | `language` | `en` |
| Theme | `theme` | `dark` |
| AutoRotateInterval | `autoRotateInterval` | `0` (off), seconds |
| HealthCheckInterval | `healthCheckInterval` | `180`, seconds (no UI or binding; edit config.json to change) |
| KillSwitch | `killSwitch` | `false` |
| AutoFailover | `autoFailover` | `false` |

Other API types (`LogEntry`, `BulkCheckResult`, `DNSLeakResult`, `MCPStatus`, `MCPClientInfo`): `docs/api.md`. MCP-side types (`pkg/proxycore`: `ProxyEntry` without `name`, `PoolStats {total, active, healthy, errored, avgLatencyMs}`, `FetchResult`): `docs/mcp.md`.

## Files

All under `%AppData%` (`os.UserConfigDir()`):

| Path | Format | Written by | Purpose |
|------|--------|-----------|---------|
| `Proxy0x/proxies.json` | JSON array of `ProxyEntry` | App only (atomic) | Proxy store; credentials as `enc:...` |
| `Proxy0x/proxies.json.bak` | same | App | Previous version, kept on every save |
| `Proxy0x/config.json` (+ `.bak`) | JSON `Config` | App (atomic) | Settings |
| `Proxy0x/.key` | `dpapi:<base64>` | App or MCP (first use) | AES-256 key protected with DPAPI (current user) |
| `Proxy0x/.key.legacy` | 32 raw bytes | App (once, on migration) | Pre-DPAPI key, kept as a backup |
| `Proxy0x/sysproxy-backup.json` | JSON (below) | App, while connected | User's Windows proxy settings to restore |
| `Proxy0x/killswitch-backup.json` | JSON `{"domain": "...", "private": "...", "public": "..."}`, each `<inbound>,<outbound>` as netsh prints it (lowercase) | App, while the kill switch is on | Firewall policy per profile to restore on disable |
| MCP client config + `<file>.proxy0x.bak` | JSON | App (MCP integration) | `mcpServers.proxy0x` entry and backup |

`sysproxy-backup.json`:

```json
{
  "previous": { "enable": true, "server": "corp.example:8080", "override": "*.corp;<local>" },
  "ours": "127.0.0.1:51234"
}
```

An empty `server`/`override` means the value was absent and is deleted on restore.

## Encryption

```
enc:<base64(nonce || ciphertext)>      AES-256-GCM, 12-byte random nonce
```

- Empty strings are stored as-is; values without `enc:` are read as plaintext (backward compatibility).
- Key: 32 random bytes, stored as `dpapi:` + base64 of the DPAPI-protected key. If DPAPI is unavailable when the key is created, it falls back to raw bytes.
- Migration: a 32-byte raw `.key` is copied to `.key.legacy` and rewritten in DPAPI form.
- A `dpapi:` key that cannot be unprotected (e.g. profile moved to another machine) is an error, not replaced; the store then goes read-only. A `.key` that is neither `dpapi:` nor 32 bytes (corrupted) is also an error and is never replaced (restore from `.key.legacy`).
- Implementation: `internal/crypto`. `pkg/proxycore` (MCP) uses the same package, so both processes share one key and one format.

## Load and save rules

- **Save** (app only): encrypt credentials (an encryption failure aborts the save), strip runtime fields, write `proxies.json.tmp`, keep `proxies.json.bak`, rename over `proxies.json`.
- **Load**: a missing file is a normal first run. Invalid JSON or undecryptable credentials put the store in read-only mode (no save will overwrite the file; mutations return errors; `store-error` event).
- **MCP server**: loads `proxies.json` at startup and never saves; add/remove tools return an error.
