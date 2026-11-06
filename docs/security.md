# Security

## Credential encryption

Proxy credentials (username and password) are encrypted at rest using **AES-256-GCM**.

### Key management

- A 32-byte random AES key is generated on first use and stored at `%AppData%/Proxy0x/.key`, protected with **Windows DPAPI** (`CryptProtectData`, current-user scope, `CRYPTPROTECT_UI_FORBIDDEN`). File format: `dpapi:` + base64(DPAPI blob).
- **Migration:** a legacy raw 32-byte key file is copied to `.key.legacy` and rewritten in the `dpapi:` format. If DPAPI fails during migration, the raw file is kept and used.
- **Unreadable DPAPI key** (e.g. profile moved to another machine/user): loading fails with an error that points to `.key.legacy`; the key is not regenerated, so the store goes read-only instead of losing credentials.
- If DPAPI is unavailable when a new key is created, the key is written raw (fallback).
- A `.key` file that is neither `dpapi:` nor 32 raw bytes (corrupted) is an error and is never replaced; restore it from `.key.legacy`.
- One implementation for both processes: the MCP server's store (`pkg/proxycore/store.go`) delegates to `internal/crypto`, so it reads the DPAPI key and never creates or rewrites `.key`. (Before this, the MCP copy only understood raw keys and would overwrite the DPAPI key, making every saved credential unreadable.)
- File permissions (`0600`) have no real effect on NTFS; protection comes from DPAPI.

### Encryption format

```
"enc:" + base64(nonce || ciphertext)
```

- 12-byte random nonce prepended to the ciphertext
- Standard library `crypto/aes` + `crypto/cipher` (GCM mode)
- Empty strings are stored as empty (not encrypted)
- Plaintext values (without `enc:` prefix) are returned as-is for backward compatibility

### Store lifecycle

1. **Save:** encrypt username and password, strip runtime state (status, currentIp, active), write JSON atomically (temp file + rename, previous version kept as `.bak`)
2. **Load:** read JSON, decrypt credentials, reset all proxies to `offline` / `active: false`. A decryption failure is a load error: the store goes read-only instead of re-encrypting ciphertext as if it were plaintext.
3. **Writers:** only the desktop app saves `proxies.json`; the MCP server is read-only.

## SSRF protection

The MCP server's `proxy_fetch` tool validates every URL before making a request through the proxy pool.

### Blocked schemes

Only `http` and `https` are allowed. `ftp`, `file`, `javascript`, and schemeless URLs are rejected.

### Blocked IP ranges

| Range | Reason |
|-------|--------|
| `10.0.0.0/8` | RFC 1918 private |
| `172.16.0.0/12` | RFC 1918 private |
| `192.168.0.0/16` | RFC 1918 private |
| `127.0.0.0/8` | Loopback |
| `0.0.0.0/8` | "This" network |
| `169.254.0.0/16` | Link-local |
| `::1/128` | IPv6 loopback |
| `fc00::/7` | IPv6 ULA |
| `fe80::/10` | IPv6 link-local |
| `100.64.0.0/10` | Carrier-grade NAT (RFC 6598) |
| `198.18.0.0/15` | Benchmarking (RFC 2544) |
| `224.0.0.0/4`, `240.0.0.0/4` | Multicast, reserved |
| `64:ff9b::/96` | NAT64 (embeds IPv4) |
| `2002::/16` | 6to4 (embeds IPv4) |

IPv4-mapped IPv6 addresses (`::ffff:a.b.c.d`) are checked as their IPv4 form.

### Cloud metadata IPs

| IP | Service |
|----|---------|
| `169.254.169.254` | AWS EC2 / GCP / Azure metadata |
| `169.254.170.2` | AWS ECS task metadata |
| `100.100.100.200` | Alibaba Cloud metadata |
| `fd00:ec2::254` | AWS EC2 IPv6 metadata |

### Blocked hostnames

- `localhost`
- `metadata`
- Any hostname ending in `.local`, `.localhost`, `.internal`
- `metadata.google.internal`

Hostnames are lowercased and a trailing dot is removed before these checks (`metadata.google.internal.` is blocked).

### DNS resolution check

After hostname validation, the URL's hostname is resolved via DNS. Every returned IP address is checked against the blocked IP ranges above. **If resolution fails, the URL is rejected** (fail closed).

Limitation: this is a point-in-time check. The request itself is sent through the upstream proxy, which resolves the name again, so it does not prevent DNS rebinding. The practical impact is bounded because egress happens from the proxy provider's network, not the user's.

### Redirect validation

HTTP redirects are followed up to a maximum of **10 hops**. Each redirect target URL is fully SSRF-validated (scheme, hostname, IP, DNS) before following. This prevents open-redirect bypasses.

### Body size limit

Response bodies from `proxy_fetch` are capped at **100 KB**. Larger responses are truncated and flagged with a `truncated: true` field.

## Credential stripping

The `proxy_list` MCP tool and the `proxy://pool` MCP resource strip `username` and `password` fields from output. AI agents see proxy hosts and ports but never credentials.

The desktop app never sends the proxy password to the webview: every proxy returned by a binding (`GetProxies`, `ConnectProxy`, …) or emitted in a `proxy-updated` event is a copy taken under the lock with `password` blanked. The Go side keeps the password for connecting and persisting. `ExportProxies` is the only path that returns credentials, on explicit user action.

## Windows system proxy

When a proxy is connected, the app:

1. Starts a local HTTP/SOCKS5 server on `127.0.0.1`
2. Saves the user's current `ProxyEnable`, `ProxyServer` and `ProxyOverride` (under `HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings`) to `%AppData%/Proxy0x/sysproxy-backup.json`, together with the address it is about to set. An existing backup is kept (a rotation or reconnect never overwrites the user's original values with ours).
3. Points the system proxy at the local server
4. Traffic of applications that honor the Windows (WinINET) proxy routes through the local server → upstream proxy

On disconnect or app shutdown the saved values are restored exactly (including a corporate proxy the user had before) and the backup file is deleted.

On startup (crash recovery): if the registry still points at the address recorded in the backup, the saved values are restored. Settings the user changed in the meantime are left alone (the stale backup is discarded). Without a backup, only a leftover from older versions (`127.0.0.1:*` with Proxy0x's exact bypass list) is disabled; any other proxy — including another local tool such as `127.0.0.1:8888` — is never touched.

Tracked TCP connections are force-closed for instant disconnect.

## Connection tracking

The local proxy server tracks every TCP connection it handles — the client side and the upstream side, for both CONNECT tunnels and plain HTTP requests — via `trackConn`/`untrackConn`. `Stop()` closes every tracked connection before waiting for its goroutines, so it returns in under a second regardless of in-flight traffic.

## Kill switch (`internal/killswitch`)

Toggled in Settings → Security. Enabling runs an elevated script (UAC prompt) that adds three Windows Firewall rules:

| Rule | Effect |
|------|--------|
| `Proxy0x-KillSwitch-Block` | block all outbound traffic |
| `Proxy0x-KillSwitch-AllowApp` | allow outbound for the Proxy0x executable |
| `Proxy0x-KillSwitch-AllowLocal` | allow outbound to `127.0.0.1` / `::1` |

Disabling deletes the three rules (another UAC prompt). On app shutdown, if the kill switch is on, the app disables it (UAC prompt at quit) and stores `killSwitch=false`; it is not re-enabled on the next start.

Known risks (not verified on a real machine yet):
- **Rule precedence (unverified):** in Windows Defender Firewall, block rules take precedence over allow rules. If that applies here, the allow-app rule does not exempt Proxy0x and enabling the kill switch cuts *all* internet access, including the proxy itself.
- **Rules persist:** if the app is killed, crashes or is uninstalled while the kill switch is on, the firewall rules stay and outbound traffic stays blocked until they are removed (Disable in the app, or `netsh advfirewall firewall delete rule name=Proxy0x-KillSwitch-Block`). The uninstaller runs without elevation and does not remove them.
- **Elevated script handling:** the script is written to a fixed path (`%TEMP%\proxy0x-fw.bat`) and executed elevated via `ShellExecute("runas")`. `ShellExecute` returns before the script runs and the file is deleted right away, so the script may be gone before `cmd.exe` reads it (the toggle can silently do nothing), and another process running as the user could replace the file between write and elevated execution.

`EnableUWPSupport` uses the same elevated-script mechanism to add `CheckNetIsolation LoopbackExempt` entries for Edge, Windows Search, Microsoft Store and Paint, so those UWP apps may use the local proxy on `127.0.0.1`.

## DNS leak test

`CheckDNSLeak` compares the exit IP seen through the active proxy with the user's direct public IP (both from `api.ipify.org`). Equal IPs → `leaked: true`; no active proxy → `leaked: true`; if either lookup fails → `leaked: false` (unknown). Despite the name it does not inspect DNS resolvers — it verifies that traffic actually leaves through the proxy.

## MCP client configuration writes

`mcp_integration.go` edits other applications' config files to register the MCP server: Claude Desktop (`%AppData%\Claude\claude_desktop_config.json`), Codex (`~\.codex\config.json`), Cursor (`~\.cursor\mcp.json`), Windsurf (`~\.codeium\windsurf\mcp_config.json`). An unparseable file is never overwritten; the previous version is saved as `<file>.proxy0x.bak` and writes are atomic. Hypothesis: Codex reads TOML (`~\.codex\config.toml`), so the Codex entry may have no effect.

## Implementation locations

| Component | File |
|-----------|------|
| AES-GCM + DPAPI key (app and MCP) | `internal/crypto/crypto.go` (used by `pkg/proxycore/store.go`) |
| Kill switch / UWP elevation | `internal/killswitch/killswitch.go` |
| MCP client config writes | `mcp_integration.go` |
| SSRF validation | `pkg/proxycore/ssrf.go` |
| SSRF tests | `pkg/proxycore/ssrf_test.go` |
| Redirect validation | `pkg/proxycore/fetcher.go` |
| Credential stripping | `internal/mcpserver/tools.go` |
| System proxy | `internal/proxy/system.go` |
| Connection tracking | `internal/proxy/server.go` |
