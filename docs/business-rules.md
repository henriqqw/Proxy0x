# Business Rules

## Proxy Scoring

Each proxy tracks `Successes` and `Failures` counters, updated on every proxied request.

```
Score() = Successes * 2 - Failures
```

`BestProxy()` selects the active, connected proxy with the highest score. Ties are broken by iteration order (map traversal). A pool with no active proxies returns an error.

## Session Rotation

Residential proxy providers route traffic based on a session identifier embedded in the username (`user_session-<id>`). Calling `Rotate()`:

1. Generates a new 8-byte hex `SessionID` via `crypto/rand`.
2. The next request uses the new session, which the upstream provider maps to a different residential IP.
3. Verifies the new IP via `api.ipify.org` through the proxy.

## Auto-Rotation

Configured in Settings as an interval in seconds. Presets: 1, 3, 5, 10 minutes. Custom values: 1–120 minutes.

- A background goroutine (`autoRotateLoop`) sleeps for the configured interval, then rotates all active proxies.
- Setting the interval to `0` disables auto-rotation.
- The interval is persisted in `config.json` and restored on startup.

## Health Check Loop

A background goroutine (`healthCheckLoop`) runs every `HealthCheckInterval` seconds (default: 180s / 3 min).

For each active proxy:
1. Fetch `api.ipify.org` through the proxy with a 10-second timeout.
2. If the fetch fails, mark the proxy as `error` and emit a Windows notification.
3. If the fetch succeeds, update `CurrentIP` and `LastCheck`.

## Connection Lifecycle

### Connect

1. Set `Active = true`, `Status = connecting`.
2. Start a local proxy server (`LocalServer`) on `127.0.0.1` that forwards traffic to the upstream proxy.
3. Set the Windows system proxy (registry) to point to the local server.
4. Verify the external IP via `api.ipify.org` through the proxy.
5. On success: `Status = connected`, record `ConnectedAt` timestamp, fetch geolocation from `ip-api.com`.
6. On failure: `Status = error`, log the error, clear system proxy.

### Disconnect

1. Close all tracked connections (instant, no hanging sockets).
2. Stop the local proxy server.
3. Clear the Windows system proxy registry entries.
4. Set `Active = false`, `Status = offline`.
5. Accumulate elapsed time into `TotalUptime`.

## Credential Encryption

Uses AES-256-GCM with a per-machine key stored at `%AppData%/Proxy0x/.key`.

- **On save**: plaintext credentials are encrypted and stored with an `enc:` prefix followed by base64-encoded ciphertext.
- **On load**: values with the `enc:` prefix are decrypted; plaintext values pass through unchanged (backward compatibility).
- Empty strings are not encrypted.

## Proxy Persistence

Proxies are stored in `%AppData%/Proxy0x/proxies.json` as a JSON array.

- Auto-saved after every add, remove, or update operation.
- On load, all proxies have `Status` reset to `offline` and `Active` set to `false` (no stale connected state after a crash).
- Credentials are encrypted at rest; decrypted transparently on load.

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
