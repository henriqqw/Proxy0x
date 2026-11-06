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

Toggled in Settings → Security. Windows Firewall applies block rules before allow rules, so a "block everything" rule would also block Proxy0x. The kill switch instead changes each profile's **default outbound action** to block and adds allow rules for what must keep working.

Enabling (one UAC prompt):
1. Reads the current policy of the domain, private and public profiles (`netsh advfirewall show <profile>profile firewallpolicy`) and saves it to `%AppData%/Proxy0x/killswitch-backup.json`. An existing backup is kept, so enabling twice never records Proxy0x's own `blockoutbound` as the user's policy.
2. Adds the allow rules:

   | Rule | Effect |
   |------|--------|
   | `Proxy0x-KillSwitch-AllowApp` | outbound for the Proxy0x executable |
   | `Proxy0x-KillSwitch-AllowLocal` | outbound to `127.0.0.1` / `::1` (apps reach the local proxy) |
   | `Proxy0x-KillSwitch-AllowDNS-UDP` / `-TCP` | DNS Client service (`svchost.exe`, `dnscache`) to port 53, so Proxy0x can resolve upstream hostnames |

3. Sets each profile to `<its inbound policy>,blockoutbound`, last, so a failed rule never leaves the machine blocked without an exception for Proxy0x.

Disabling restores the saved policy first and deletes the rules last (including `Proxy0x-KillSwitch-Block`, left by versions up to 1.1.0), then removes the backup. Without a backup, outbound is reopened only on profiles that block it while Proxy0x's allow rule exists; a policy the user set themself is left alone. If no Proxy0x rule exists, Disable does nothing (no UAC prompt).

`GetKillSwitch` reports whether `Proxy0x-KillSwitch-AllowApp` exists, not the saved setting, so rules left by a crash show as on and can be turned off.

Lifecycle:
- **Quit:** the app disables the kill switch (UAC prompt at quit).
- **Crash or killed process:** the block stays on (fail closed). Start Proxy0x and turn it off, or run `Proxy0x.exe --disable-killswitch`.
- **Uninstall:** the uninstaller checks for `Proxy0x-KillSwitch-AllowApp` and, if present, runs `Proxy0x.exe --disable-killswitch` as the user (the policy backup is in their AppData); that asks for administrator rights. If the prompt is declined, remove it manually as administrator: restore the policy with `netsh advfirewall set allprofiles firewallpolicy blockinbound,allowoutbound` (Windows default) and delete the `Proxy0x-KillSwitch-*` rules.

Elevation (`RunElevated`): the whole command is passed on the command line of `%SystemRoot%\System32\cmd.exe /S /C "…"` through `ShellExecuteExW("runas")`; no script file is written, so nothing can be swapped between write and execution. The call waits for the process (60 s limit) and fails on a non-zero exit code or a declined UAC prompt (`ErrCancelled`). Values that reach the command are checked: policy strings must match `(blockinboundalways|blockinbound|allowinbound),(blockoutbound|allowoutbound)` (the backup file is user-writable), and program paths containing `"` or `%` are rejected.

Limitations:
- **Existing allow rules still apply.** Any program that already has an enabled outbound allow rule in Windows Firewall keeps network access: Windows services, Store apps and some games or launchers. Ordinary desktop apps without such a rule (Chrome, Edge, Firefox, …) are blocked unless they go through the local proxy.
- **DNS queries** from the DNS Client service are allowed, so hostnames looked up while the kill switch is on can reach the configured DNS server directly.
- **Firewall must be on.** If Windows Firewall is turned off for a profile, the policy has no effect there.
- **Program path:** the allow rule names the executable's path at the time it was enabled. Moving Proxy0x while the kill switch is on blocks it too; turn the kill switch off and on again.

`EnableUWPSupport` uses the same `RunElevated` to add `CheckNetIsolation LoopbackExempt` entries for Edge, Windows Search, Microsoft Store and Paint, so those UWP apps may use the local proxy on `127.0.0.1`. Packages that are not installed are skipped.

## DNS leak test

`CheckDNSLeak` compares the exit IP seen through the active proxy with the user's direct public IP (both from `api.ipify.org`). Equal IPs → `leaked: true`; no active proxy → `leaked: true`; if either lookup fails → `leaked: false` (unknown). Despite the name it does not inspect DNS resolvers — it verifies that traffic actually leaves through the proxy.

## MCP client configuration writes

`mcp_integration.go` edits other applications' config files to register the MCP server: Claude Desktop (`%AppData%\Claude\claude_desktop_config.json`), Codex (`~\.codex\config.json`), Cursor (`~\.cursor\mcp.json`), Windsurf (`~\.codeium\windsurf\mcp_config.json`). An unparseable file is never overwritten; the previous version is saved as `<file>.proxy0x.bak` and writes are atomic. Hypothesis: Codex reads TOML (`~\.codex\config.toml`), so the Codex entry may have no effect.

## Implementation locations

| Component | File |
|-----------|------|
| AES-GCM + DPAPI key (app and MCP) | `internal/crypto/crypto.go` (used by `pkg/proxycore/store.go`) |
| Kill switch | `internal/killswitch/killswitch.go` (tests: `killswitch_test.go`) |
| Elevated commands (kill switch, UWP) | `internal/killswitch/elevate.go` |
| Kill switch cleanup on uninstall | `build/installer.iss` (`[Code]`), `main.go` (`--disable-killswitch`) |
| MCP client config writes | `mcp_integration.go` |
| SSRF validation | `pkg/proxycore/ssrf.go` |
| SSRF tests | `pkg/proxycore/ssrf_test.go` |
| Redirect validation | `pkg/proxycore/fetcher.go` |
| Credential stripping | `internal/mcpserver/tools.go` |
| System proxy | `internal/proxy/system.go` |
| Connection tracking | `internal/proxy/server.go` |
