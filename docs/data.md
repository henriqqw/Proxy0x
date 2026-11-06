# Data Models & Persistence

## Models

### ProxyEntry

Primary data model for a proxy in the pool.

| Field | Type | JSON | Description |
|-------|------|------|-------------|
| ID | string | `id` | Random 8-byte hex, generated on add |
| Host | string | `host` | Proxy hostname or IP |
| Port | string | `port` | Proxy port |
| Username | string | `username` | Auth username (encrypted at rest) |
| Password | string | `password` | Auth password (encrypted at rest) |
| Protocol | string | `protocol` | `http` or `socks5` |
| Group | string | `group` | Optional group tag |
| Status | ProxyStatus | `status` | `offline` / `connecting` / `connected` / `error` |
| CurrentIP | string | `currentIp` | External IP when connected |
| Location | *GeoInfo | `location` | Geolocation data, nil when unknown |
| Latency | int64 | `latency` | Last measured latency in ms |
| LastCheck | string | `lastCheck` | ISO 8601 timestamp of last health check |
| Active | bool | `active` | Whether the proxy is selected for use |
| SessionID | string | `sessionId` | Current rotation session (8-byte hex) |
| Order | int | `order` | Display order for drag-to-reorder |
| ConnectedAt | string | `connectedAt` | ISO 8601 timestamp when connected |
| TotalUptime | int64 | `totalUptime` | Accumulated uptime in seconds |

### GeoInfo

IP geolocation data fetched from `ip-api.com`.

| Field | Type | JSON |
|-------|------|------|
| Country | string | `country` |
| CountryCode | string | `countryCode` |
| Region | string | `region` |
| City | string | `city` |
| ISP | string | `isp` |
| Org | string | `org` |
| Lat | float64 | `lat` |
| Lon | float64 | `lon` |
| Timezone | string | `timezone` |
| AS | string | `as` |

### LogEntry

Connection event log entry.

| Field | Type | JSON |
|-------|------|------|
| Time | string | `time` |
| Action | string | `action` |
| ProxyID | string | `proxyId` |
| Host | string | `host` |
| Detail | string | `detail` |

### MCPStatus

Status of the MCP (Model Context Protocol) integration with Claude Desktop.

| Field | Type | JSON |
|-------|------|------|
| Installed | bool | `installed` |
| ExePath | string | `exePath` |
| ExeExists | bool | `exeExists` |
| ConfigPath | string | `configPath` |
| HealthOK | bool | `healthOk` |
| HealthError | string | `healthError` |

### PoolStats

Summary statistics for the proxy pool.

| Field | Type | JSON |
|-------|------|------|
| Total | int | `total` |
| Active | int | `active` |
| Connected | int | `connected` |

### Config

Application configuration persisted between sessions.

| Field | Type | JSON | Default |
|-------|------|------|---------|
| Language | string | `language` | `en` |
| Theme | string | `theme` | `dark` |
| AutoRotateInterval | int | `autoRotateInterval` | `0` (disabled) |
| HealthCheckInterval | int | `healthCheckInterval` | `180` (3 min) |

## File Locations

All under `%AppData%` (`os.UserConfigDir()` on Windows).

| Path | Format | Description |
|------|--------|-------------|
| `Proxy0x/proxies.json` | JSON array | Encrypted proxy store |
| `Proxy0x/.key` | Raw bytes (32B) | AES-256 encryption key, auto-generated |
| `Proxy0x/config.json` | JSON object | App settings |
| `Claude/claude_desktop_config.json` | JSON object | Claude Desktop MCP server config |

## Encryption Format

Credentials are stored with an `enc:` prefix:

```
enc:<base64(nonce + ciphertext)>
```

- Algorithm: AES-256-GCM
- Key: 32-byte random key at `Proxy0x/.key`, created on first use and reused.
- Nonce: 12-byte random, prepended to ciphertext.
- Empty strings are stored as-is (not encrypted).
- Values without the `enc:` prefix are treated as plaintext (backward compatibility).

## MCP Server Data Sharing

The Wails app and the MCP server (`proxy0x-mcp.exe`) are separate processes that share the same data files:

- Both read/write `Proxy0x/proxies.json` using the same encryption key.
- The MCP server loads proxies on startup and operates on its own in-memory pool.
- No IPC between the two processes; changes made by one are visible to the other on next load.
