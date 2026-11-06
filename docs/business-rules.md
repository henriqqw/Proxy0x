# Business Rules

## Proxy Selection (MCP)

Only one proxy is active at a time, so the MCP server's `ActiveProxy()` returns the connected proxy, or an error when none is connected. There is no scoring.

## Session Rotation

Residential proxy providers route traffic based on a session identifier embedded in the username (`user_session-<id>`). Calling `Rotate()`:

1. Generates a candidate 6-byte hex `SessionID` via `crypto/rand`.
2. Verifies the exit IP through the proxy with the candidate session (see **Exit IP verification**).
3. On success: the session switches to the candidate, the old HTTP client is closed, the geolocation is refreshed and the local proxy switches new connections to the new session in place (same port, open browser connections kept, no registry change).
4. On failure: the previous session and status are kept (traffic keeps flowing on the old session) and an error is returned.

## Exit IP verification

Connect, rotate, health check and speed test verify the exit IP by fetching `api.ipify.org?format=json` through the proxy. A check succeeds only with HTTP 200 and a body whose `ip` field parses as an IP address; anything else (non-200, malformed JSON, empty IP, network error) is a failure. Each request times out after 20s. Network calls (verification and geolocation) never run while holding the proxy-list lock.

## Auto-Rotation

Configured in Settings as an interval in seconds. Presets: 1, 3, 5, 10 minutes. Custom values: 1–120 minutes.

- A background goroutine (`autoRotateLoop`) sleeps for the configured interval, then rotates all active proxies.
- Setting the interval to `0` disables auto-rotation.
- The interval is persisted in `config.json` and restored on startup.

## Health Check Loop

A background goroutine (`healthCheckLoop`) runs every `HealthCheckInterval` seconds (default: 180s / 3 min).

For each active proxy:
1. Verify the exit IP (see **Exit IP verification**).
2. If verification fails, mark the proxy as `error` and emit a Windows notification.
3. If it succeeds, update `Latency`, `LastCheck` and, when the IP changed, `CurrentIP` and geolocation.

## Connection Lifecycle

Only one proxy is active at a time. Connecting a proxy disconnects the currently active one first.

### Connect

1. Set `Status = connecting`.
2. Verify the exit IP through the proxy (see **Exit IP verification**). On failure: `Status = error`, IP and geolocation cleared, nothing else changes.
3. Fetch geolocation from `ipwho.is`.
4. Start a local proxy server (`LocalServer`) on `127.0.0.1` and point the Windows system proxy at it. If the local server cannot start, the connect fails: `Status = error`, system proxy untouched.
5. On success: `Status = connected`, `Active = true`, record `ConnectedAt`.

### Disconnect

1. Stop the local proxy server (closes every tracked connection, returns in under a second).
2. Clear the Windows system proxy.
3. Set `Active = false`, `Status = offline`, clear IP, latency and geolocation.
4. Accumulate elapsed time into `TotalUptime`.

### Remove

Removing the active proxy disconnects it first (local server stopped, system proxy cleared), then deletes it and closes its HTTP client.

## Credential Encryption

Uses AES-256-GCM with a per-machine key stored at `%AppData%/Proxy0x/.key`.

- **On save**: plaintext credentials are encrypted and stored with an `enc:` prefix followed by base64-encoded ciphertext.
- **On load**: values with the `enc:` prefix are decrypted; plaintext values pass through unchanged (backward compatibility).
- Empty strings are not encrypted.

## Proxy Persistence

Proxies are stored in `%AppData%/Proxy0x/proxies.json` as a JSON array.

- Auto-saved after every add, remove, or update operation.
- Writes are atomic: the new content goes to `proxies.json.tmp` and is renamed over `proxies.json`; the previous version is kept as `proxies.json.bak`. `config.json` is written the same way.
- On load, all proxies have `Status` reset to `offline` and `Active` set to `false` (no stale connected state after a crash).
- Credentials are encrypted at rest; decrypted transparently on load.
- If `proxies.json` exists but cannot be read (invalid JSON or credentials that fail to decrypt), the store becomes **read-only** for the session: nothing is saved over the file, mutating calls return an error and a `store-error` event is emitted. A missing file is a normal first run.
- Only the desktop app writes `proxies.json`. The MCP server reads it and never saves (see `docs/mcp.md`).

## Log Ring Buffer

Connection events are stored in a ring buffer capped at **200 entries**. Newest entries are prepended (index 0 = most recent). When the buffer is full, the oldest entry is discarded.

Each `LogEntry` records: timestamp, action type, proxy ID, host, and detail string.

## Proxy Format

```
[protocol://]host:port[:user:pass]
```

- Protocol prefix is optional. Supported: `http://` (default), `socks5://`.
- Minimum: `host:port`. With credentials: `host:port:user:pass`.
- `ParseProxy()` extracts fields and generates a random `ID` and `SessionID`.

## Group Management

- Each proxy has an optional `Group` string tag.
- Groups are derived dynamically from the proxy list (no separate group store).
- Operations: `ConnectGroup(name)`, `DisconnectGroup(name)`, filter by group in the UI.

## Speed Test

`BenchmarkProxy()` measures two metrics:

1. **Connect latency** (ms): time to establish a TCP connection through the proxy.
2. **Download speed** (Mbps): download a 10 MB test file through the proxy, calculate throughput.

Results are returned as a map with `connectTime` and `downloadSpeed` keys.
