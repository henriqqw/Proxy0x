# 00 — Linha de base

Data: 2026-11-06 · Commit auditado: `0d5768e` (main) · Auditor: Claude (Opus 5.5)

## Estado do repositório no momento da auditoria

- Branch `main`, em dia com `origin/main`.
- Alterações **não commitadas de outra sessão** em: `frontend/src/components/SettingsPage.tsx`, `frontend/src/i18n/translations.ts`, `mcp_integration.go`, `frontend/wailsjs/go/*`. Pasta `.claude/` não rastreada.
- A auditoria cobre o **HEAD**. O WIP só foi compilado e checado por tipos, não foi revisado.
- Os relatórios ficam como arquivos novos em `docs/audit/`. Nenhuma troca de branch foi feita nesta pasta, para não interferir na sessão que está editando.

## Stack (verificada no código)

| Camada | Tecnologia | Fonte |
|---|---|---|
| Desktop | Go (diretiva `go 1.25.5`, toolchain local 1.26.5), Wails v2.16.0 | `go.mod`, `wails.json` |
| Frontend | React 18.3, TypeScript 5.6, Vite 6, Tailwind 3.4, @dnd-kit, lucide-react | `frontend/package.json` |
| Proxy | `net`, `net/http`, `golang.org/x/net/proxy` (SOCKS5) | `internal/proxy/server.go` |
| Tray | `github.com/energye/systray` | `tray.go` |
| MCP | `github.com/mark3labs/mcp-go` v1.5.0 (stdio) | `cmd/proxy0x-mcp`, `internal/mcpserver` |
| Instalador | Inno Setup 6 (manual) **e** NSIS via Wails (CI) | `installer.iss`, `.github/workflows/release.yml` |

## Pontos de entrada

1. `main.go`: app Wails. Faz bind de `*App` (todos os métodos exportados de `app.go` e `mcp_integration.go`).
2. `cmd/proxy0x-mcp/main.go`: servidor MCP standalone via stdio. Lê o mesmo `proxies.json` e `.key`.
3. `tray.go`: menu da bandeja (Connect/Disconnect/Rotate All, idioma, Quit).
4. `.github/workflows/release.yml`: **todo push na `main`** cria uma tag `v1.0.N+1` e publica uma release.
5. `installer.iss`: build manual do instalador Inno Setup.

## Fluxo de dados

```
React UI ──Wails bindings──► App (app.go)
                               ├─ http.Client por proxy+sessão ──► upstream ──► api.ipify.org / ipwho.is / speed.cloudflare.com
                               ├─ LocalServer 127.0.0.1:<porta aleatória> ◄── navegador (WinINET via registro HKCU)
                               │                                  └──► upstream HTTP CONNECT / SOCKS5
                               ├─ %APPDATA%/Proxy0x/proxies.json (credenciais AES-GCM) + .key + config.json
                               └─ eventos Wails: proxy-updated, proxies-changed, proxy-dropped, log-entry
Claude Desktop ──stdio──► proxy0x-mcp ──► pkg/proxycore (pool próprio, mesmo store/.key)
mcp_integration.go ──► escreve %APPDATA%/Claude/claude_desktop_config.json
```

## Comandos executados e resultados

| Comando | Resultado |
|---|---|
| `go build ./...` | OK (exit 0) |
| `go vet ./...` | OK, sem avisos |
| `go test ./...` | `internal/mcpserver` OK, `pkg/proxycore` OK. **Sem testes:** `proxy_ui` (raiz), `internal/proxy`, `internal/crypto`, `internal/config`, `cmd/proxy0x-mcp` |
| `go test -race ./...` | **Não executado**: exige cgo e não há `gcc` na máquina |
| `govulncheck@v1.1.4 ./...` | **14 vulnerabilidades alcançáveis na stdlib** (`crypto/tls`, `net/http`, `encoding/asn1`…) em go1.26.5, todas corrigidas em **go1.26.6**. 7 em pacotes e 5 em módulos não alcançáveis. |
| `npx tsc --noEmit` (frontend) | OK (exit 0, inclui o WIP) |
| `npm audit` / `npm audit --omit=dev` | 0 vulnerabilidades |
| `staticcheck` | Não instalado, não executado |
| Testes de reprodução (scratchpad, fora do repo) | 2 falhas que confirmam os achados A-01 e A-02 (ver 01-findings) |

## CI / releases (estado observado via `gh`)

- Tags: `v1.0.0`, `v1.0.5`–`v1.0.11`. A release **"Latest" é a `v1.0.11`** (NSIS, gerada pelo CI num commit só de docs), não a `v1.0.0` com o instalador Inno Setup publicada à mão.
- As últimas 8 execuções do workflow `Release` passaram (verde).

## Cobertura de testes (linha de base)

- Código rastreado do app desktop (raiz, `internal/proxy`, `internal/crypto`, `internal/config`, frontend): **0 testes**.
- `pkg/proxycore` e `internal/mcpserver`: suítes existentes passando.
- Essa é a régua de "sem regressão" para a fase 3.
