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

- A background goroutine (`autoRotateLoop`) polls every 10s and rotates all active proxies once the configured interval has elapsed.
- Setting the interval to `0` disables auto-rotation.
- The interval is persisted in `config.json` and restored on startup.

## Health Check Loop

A background goroutine (`healthCheckLoop`) polls every 30s and checks active proxies every `HealthCheckInterval` seconds (default: 180s / 3 min), 1.5s apart.

For each active proxy:
1. Verify the exit IP (see **Exit IP verification**).
2. If verification fails, mark the proxy as `error`. If it was `connected`, log a drop, emit `proxy-dropped` (Windows notification) and run **Auto-Failover**. The system proxy keeps pointing at the local server (fail-closed): traffic stops until the user acts or failover connects another proxy.
3. If it succeeds, set `Status = connected` again (recovery after a transient failure), update `Latency`, `LastCheck` and, when the IP changed, `CurrentIP` and geolocation.

## Auto-Failover

Toggle in Settings → Security (`config.json` → `autoFailover`, off by default). When a connected proxy drops (step 2 above):

1. Candidates: every proxy except the failed one that is not active.
2. Pick the candidate with the lowest positive `Latency`; if none has a latency, the first candidate found (map order, not deterministic).
3. Log `failover` and call `ConnectProxy` on it in the background, which disconnects the failed proxy first.

There is no retry limit or backoff: if the new proxy also drops, the next health check fails over again.

## Kill Switch

Toggle in Settings → Security. Sets the Windows Firewall default outbound action to block (the previous policy is saved) and allows only Proxy0x, loopback and DNS, after one UAC prompt. The toggle shows the real firewall state. Turning it off or quitting the app restores the previous policy (UAC again); it is not re-enabled automatically on the next start. If the app crashes, the block stays until it is turned off; uninstalling removes it. Details and limitations: `docs/security.md` → Kill switch.

## Connect Fastest

`ConnectFastest()` verifies the exit IP of every inactive proxy concurrently (10 at a time), takes the one with the lowest latency and connects it. Fails with "no proxy responded" when none passes verification.

## Bulk Check

`BulkCheckProxies()` verifies every proxy concurrently (10 at a time) and returns per proxy: alive, latency, exit IP, country and an anonymity label — `transparent` when the exit IP equals the user's real IP, `elite` when different, `unknown` when either IP is missing. It updates each proxy's `Latency` and `LastCheck` and emits `bulk-check-progress` per result.

## DNS Leak Test

`CheckDNSLeak()` compares the active proxy's exit IP with the user's direct IP: equal → leaked; no active proxy → leaked; a failed lookup → not leaked (unknown). It does not inspect DNS resolvers.

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

Uses AES-256-GCM with a per-user key stored at `%AppData%/Proxy0x/.key`, itself protected with Windows DPAPI (see `docs/security.md` for format and migration).

- **On save**: plaintext credentials are encrypted and stored with an `enc:` prefix followed by base64-encoded ciphertext. An encryption failure aborts the save (plaintext is never written).
- **On load**: values with the `enc:` prefix are decrypted; plaintext values pass through unchanged (backward compatibility). A decryption failure puts the store in read-only mode.
- Empty strings are not encrypted.
- The desktop app and the MCP server use the same implementation (`internal/crypto`).

## Proxy Persistence

Proxies are stored in `%AppData%/Proxy0x/proxies.json` as a JSON array.

- Auto-saved after every add, remove, or update operation.
- Writes are atomic: the new content goes to `proxies.json.tmp` and is renamed over `proxies.json`; the previous version is kept as `proxies.json.bak`. `config.json` is written the same way.
- On load, all proxies have `Status` reset to `offline` and `Active` set to `false` (no stale connected state after a crash).
- Credentials are encrypted at rest; decrypted transparently on load.
- If `proxies.json` exists but cannot be read (invalid JSON or credentials that fail to decrypt), the store becomes **read-only** for the session: nothing is saved over the file, mutating calls return an error and a `store-error` event is emitted. A missing file is a normal first run.
- Only the desktop app writes `proxies.json`. The MCP server reads it and never saves (see `docs/mcp.md`).

## Log Ring Buffer

Connection events are kept in memory, capped at **200 entries** (oldest discarded). `GetConnectionLog()` returns them newest first. Each `LogEntry` records: timestamp, action type, proxy ID, host, and detail string. Every entry is also emitted as `log-entry`; the current UI does not display the log.

## Proxy Format

Single add and batch add (`AddProxy`, `AddProxiesBatch`):

```
[protocol://]host:port[:user:pass]
```

- Protocol prefix is optional. Supported: `http://` (default), `socks5://`.
- `parseProxy()` extracts fields and generates a random `ID` and 6-byte `SessionID`.

Import in Settings (`AddProxiesMultiFormat`) also accepts, one per line (blank lines and `#` comments ignored; invalid lines skipped silently):

- URIs: `socks5://user:pass@host:port`, `http://…`, `https://…` (port defaults to `8080`). `socks4://` is accepted but treated as **SOCKS5**.
- CSV: `host,port[,user,pass[,protocol]]` where protocol `socks5`/`socks4` → SOCKS5, anything else → HTTP.

## Proxy Name

Each proxy has an optional display `Name` (`SetProxyName`); the card shows it instead of `host:port` when set. Persisted in `proxies.json`.

## Group Management

- Each proxy has an optional `Group` string tag (`SetProxyGroup`).
- Groups are derived dynamically from the proxy list (no separate group store).
- The UI filters the list by group. If the selected group disappears (its last proxy removed or regrouped) the filter resets to all proxies. There are no group-wide connect/disconnect operations: only one proxy can be active.

## Speed Test

`BenchmarkProxy()` returns a map with:

1. `connectMs` — exit-IP verification time through the proxy (see **Exit IP verification**).
2. `downloadMs` / `downloadMbps` — download 100 KB from `speed.cloudflare.com` through the proxy; `downloadMs = -1` and `downloadMbps = 0` when the download fails.
