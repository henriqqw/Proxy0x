# 02 — Plano de correção (SDD)

Base: `9df0d8c` (main, o WIP de MCP multi-cliente já está commitado) · Achados: `01-findings.md` (aprovados em 2026-11-06)

## Como a fase 3 vai rodar

- **Branch própria:** `fix/audit-2026-11`, numa worktree separada. A pasta principal e as outras sessões não são tocadas. Nada é enviado para o remoto sem a sua aprovação.
- **SDD:** cada lote começa com um commit `docs:` que corrige a spec daquele trecho (`docs/business-rules.md`, `security.md`, `backend.md`, `mcp.md`). O código só muda depois.
- **Um problema por commit.** Refatoração (por exemplo, abrir um ponto de injeção para teste) vai num commit separado da correção.
- **Ciclo por item:** teste que falha → menor correção na causa raiz → `go build`, `go vet`, `go test ./...`, `tsc` → comparar com a linha de base (`00-baseline.md`). Se algo regredir: reverter e registrar em `03-changes.md`.
- **Rollback padrão:** `git revert <commit>` do item. As exceções estão anotadas no próprio item.
- **Ao fim de cada lote:** a suíte inteira roda e eu mando um resumo curto. Os lotes seguintes não esperam aprovação, salvo os marcados com ⚠.

### Pontos de injeção (só nas fronteiras de I/O)

São necessários para testar sem rede, sem Wails e sem registro real. Cada um entra num commit de refatoração próprio, sem mudança de comportamento:

1. `app.go`: um campo `emit` (padrão `runtime.EventsEmit`) e as variáveis `ipCheckURL` e `geoURL`.
2. `internal/proxy/system.go`: interface mínima `registry` (Get/Set de 3 valores), com implementação real via `golang.org/x/sys/windows/registry`. Hoje esse módulo já está no `go.sum` como indireto, então não entra dependência nova. O `reg.exe` sai.

---

## Decisões pendentes (padrão = recomendação; confirme ou troque)

| # | Tema | Padrão adotado no plano | Alternativas |
|---|---|---|---|
| D1 | Processo de release (OPS-01) | CI só por tag `v*` + `workflow_dispatch`; gera o **Inno Setup**, o portátil e o `proxy0x-mcp.exe`; só amd64. As releases antigas não são apagadas; a próxima tag (**v1.1.0**) vira a "Latest". | (b) manter o NSIS do Wails e remover o `installer.iss` · (c) desligar o CI e fazer release manual |
| D2 | `proxies.json` compartilhado (MCP-02) | O **MCP vira só leitura**: `proxy_add`/`proxy_remove` passam a responder "gerencie proxies pelo app" e o MCP nunca salva. ⚠ **Muda o contrato das tools MCP.** | (b) lock de arquivo entre processos + recarregar/mesclar/salvar (M) · (c) o app é o dono e o MCP usa IPC (G) |
| D3 | Health check com falha (CORE-16) | Manter fail-closed (privacidade) e **documentar** na spec | Desconectar automaticamente e restaurar o proxy do sistema |
| D4 | Recursos remotos (PRIV-01) | Trocar a Google Fonts pela fonte do sistema (Segoe UI); manter flagcdn, liberado por uma CSP estrita; corrigir o texto do README | Baixar a Inter e as bandeiras para dentro do app (precisa baixar arquivos, ou seja, mais uma aprovação) |
| D5 | Código morto (DEAD-01) | Remover os bindings sem uso **e o ring buffer de logs** (nada o consome) | Manter os logs para uma futura tela de histórico |

---

## Lote 0 — Preparação · ~3 commits

### SEC-04 · Toolchain com correções de segurança
- **Spec:** compilar com Go ≥ 1.26.6 localmente e no CI.
- **Mudança:** `toolchain go1.26.6` no `go.mod`; o CI passa a usar `go-version-file: go.mod`.
- **Aceite:** `govulncheck ./...` com **0** vulnerabilidades alcançáveis na stdlib (linha de base: 14); build e testes verdes.
- **Teste:** o próprio `govulncheck` serve como teste. **Rollback:** revert.

### Commits de refatoração: pontos de injeção 1 e 2 (acima)
- **Aceite:** sem mudança de comportamento; build, vet e testes iguais à linha de base. Teste manual rápido: conectar, rotacionar e parar no `wails dev`.

---

## Lote 1 — Proxy local · ~5 commits · prioridade máxima

Começa trazendo para `internal/proxy/server_test.go` os 2 testes de reprodução do scratchpad (os dois falham hoje).

### CORE-01 · Stop não pode travar
- **Spec:** `LocalServer.Stop()` retorna em menos de 1 s com qualquer conexão em andamento (CONNECT ou HTTP simples), e fecha todas as conexões do cliente e do upstream.
- **Mudança:** `trackConn`/`untrackConn` no `upConn` dos dois ramos de `handleHTTP`.
- **Aceite:** `TestStopReturnsQuicklyWithPlainHTTPInFlight` passa (limite de 1 s).
- **Rollback:** revert.

### CORE-02 · Keep-alive HTTP
- **Spec:** cada requisição HTTP simples recebida tem resposta. O proxy local não sustenta keep-alive: encaminha com `Connection: close` e fecha a conexão do cliente depois da resposta, então o navegador abre outra.
- **Mudança:** em `handleHTTP`, `req.Close = true` e o header `Connection: close`; se a resposta vier keep-alive, a conexão é encerrada depois do corpo dela, delimitado com `http.ReadResponse`.
- **Aceite:** `TestKeepAliveSecondRequestGetsResponse` passa, reescrito para "a 2ª requisição recebe resposta, ou a conexão é fechada logo e o cliente reabre". Novo teste: o corpo da resposta chega inteiro com `Content-Length` e com `chunked`.
- **Rollback:** revert. **Risco:** médio, por isso também há verificação manual num site HTTP real.

### CORE-14 · Erro de Accept não pode virar loop ocupado
- **Spec:** um erro temporário de `Accept` espera com backoff (5 ms dobrando até 1 s).
- **Teste:** um listener falso que devolve erro N vezes; o teste checa que houve espera.
- **Rollback:** revert.

### CORE-11 (parte do proxy) · Timeout no dial SOCKS5
- **Spec:** o dial SOCKS5 tem timeout de 15 s, igual ao ramo HTTP.
- **Mudança:** `netproxy.SOCKS5(..., &net.Dialer{Timeout: 15s})`.
- **Teste:** um upstream que aceita e não responde; o dial falha dentro do limite.
- **Rollback:** revert.

---

## Lote 2 — Estado e concorrência do App · ~9 commits

### CORE-03 · Sem escrita em mapa sob RLock
- **Spec:** criar ou buscar o client HTTP é seguro para chamadas concorrentes.
- **Mudança:** `clientsMu sync.Mutex` dedicado em `getOrCreateClient`/`closeClient` (app). Mesma correção em `pkg/proxycore/pool.go:344-351`.
- **Teste:** 50 goroutines chamando `getOrCreateClient` com sessões diferentes. Hoje deve abortar com "concurrent map writes"; depois, passa.
- **Rollback:** revert.

### CORE-06 · Sem ponteiro para servidor parado
- **Spec:** depois de `Stop()`, `a.localProxy` vira `nil`. Uma falha em `NewLocalServer` volta como erro, e o proxy fica `error` em vez de "connected" sem proxy de sistema.
- **Teste:** `NewLocalServer` injetável que falha; depois, disconnect não dá panic.
- **Rollback:** revert.

### CORE-04 · Remover o proxy ativo
- **Spec:** remover um proxy ativo desconecta antes (para o proxy local e limpa o sistema). O client é fechado pela chave `id:sessão`.
- **Teste:** um proxy ativo com servidor local falso é removido; o teste checa que `localProxy == nil` e que `clients` está vazio.
- **Rollback:** revert.

### CORE-13 · Rotação que falha não altera o estado
- **Spec:** a sessão só muda depois que a verificação passa. Em falha, a sessão anterior fica, o status volta ao anterior e um erro é retornado.
- **Teste:** com `ipCheckURL` falso devolvendo 500, o teste checa que o `SessionID` não mudou.
- **Rollback:** revert.

### CORE-12 · Verificação exige 200 e IP válido
- **Spec:** connect, rotate e health check só aceitam HTTP 200 com IP parseável (`net.ParseIP`). Qualquer outra coisa é falha.
- **Teste:** um servidor falso com 429 e com corpo inválido; o resultado é `error`.
- **Rollback:** revert.

### CORE-08 · Bandeja coerente com o modelo de 1 proxy ativo
- **Spec:** o tray oferece Show, Disconnect, Rotate (do ativo), Language e Quit. Sai o "Connect All".
- **Teste:** manual, porque o tray não tem como ser testado de forma automatizada.
- **Rollback:** revert.

### CORE-11 (resto) · Timeouts
- **Spec:** `GetMyIP` usa um client com timeout de 10 s. O transport SOCKS5 do app usa `ContextDialer` com timeout.
- **Teste:** um servidor que nunca responde; a chamada falha dentro do limite.
- **Rollback:** revert.

### CORE-15 · SOCKS5 nunca cai para conexão direta
- **Spec:** erro ao criar o dialer SOCKS5 vira erro de conexão (app e proxycore).
- **Teste:** unitário da função que monta o transport.
- **Rollback:** revert.

---

## Lote 3 — Integridade de dados · ~5 commits

### CORE-10 · Escrita atômica e recusa de salvar depois de load falho
- **Spec:** `proxies.json` e `config.json` são escritos em `arquivo.tmp` + `os.Rename`, mantendo `arquivo.bak` da versão anterior. Se o load do `proxies.json` falhou por JSON inválido, o app **não salva** por cima e avisa a UI (evento `store-error`).
- **Mudança:** helper `internal/fsutil.WriteFileAtomic`. Há 4 consumidores concretos (app, config, proxycore store, mcp_integration), então o helper se justifica.
- **Teste:** grava, corrompe o arquivo, recarrega → erro + save recusado; o `.bak` continua íntegro.
- **Rollback:** revert. Os arquivos `.bak` e `.tmp` que sobrarem são inofensivos.

### MCP-01 · Nunca sobrescrever o config de outro app
- **Spec:** `addProxy0x`/`removeProxy0x` abortam com erro quando o JSON existente é inválido. Antes de escrever, criam `<arquivo>.proxy0x.bak` e gravam de forma atômica.
- **Local:** `mcp_integration.go:136-184`.
- **Teste:** arquivo temporário com JSON inválido → erro e arquivo intacto. JSON válido com outras chaves → as chaves são preservadas.
- **Rollback:** revert.

### MCP-02 · MCP só leitura (D2) ⚠ contrato
- **Spec:** o processo MCP nunca grava no `proxies.json`. `proxy_add`/`proxy_remove` retornam um erro explicativo. Connect e rotate do MCP só alteram a memória. A doc `mcp.md` é atualizada.
- **Teste:** uma tool add retorna erro e o arquivo não muda (mtime e conteúdo).
- **Rollback:** revert, que devolve as tools de escrita.
- **Migração:** quem usa essas tools passa a adicionar proxies pela UI. Hoje não há consumidor externo conhecido.

---

## Lote 4 — Proxy do sistema · ~4 commits

### CORE-05 · Preservar a configuração do usuário
- **Spec:** ao conectar, o app guarda os valores anteriores de `ProxyEnable`, `ProxyServer` e `ProxyOverride` em `%APPDATA%/Proxy0x/sysproxy-backup.json`. Ao desconectar, restaura. No startup, só limpa se `ProxyServer` apontar para `127.0.0.1:*` (resíduo de crash) e, nesse caso, restaura o backup.
- **Teste:** `registry` falso com proxy corporativo → connect/disconnect → valores originais de volta. Startup com resíduo → restaura. Startup sem resíduo → não mexe.
- **Rollback:** revert. O arquivo de backup que sobrar é inofensivo.

### CORE-07 · Rotação sem derrubar conexões
- **Spec:** `LocalServer.UpdateUpstream(cfg)` troca o upstream para **novas** conexões. A porta, o registro e os túneis abertos continuam. Rotate e auto-rotate usam isso em vez de reiniciar o servidor.
- **Teste:** um túnel aberto continua ativo depois do `UpdateUpstream`; uma conexão nova usa o upstream novo; a porta não muda.
- **Rollback:** revert.

### CORE-16 · Fail-closed documentado (D3)
- **Spec:** um health check falho mantém o roteamento e notifica. O comportamento entra em `business-rules.md`.
- Só doc, sem teste.

---

## Lote 5 — Release e empacotamento · ~4 commits · ⚠ precisa de push

### OPS-01 · Release só por tag (D1)
- **Spec:** push na `main` não publica nada. Push de tag `v*` gera uma release com `Proxy0x-Setup-<versão>.exe` (Inno), `Proxy0x.exe` e `proxy0x-mcp.exe`. O `workflow_dispatch` com `dry_run=true` só gera os artefatos do workflow.
- **Mudança:** gatilhos novos; remover o auto-tag; `wails` fixo em `v2.16.0`; actions fixadas por SHA; `choco install innosetup` + `iscc /DMyAppVersion=<tag>`.
- **Aceite:** um `workflow_dispatch` (dry-run) na branch passa e produz os 3 artefatos.
- **⚠ Exige push da branch para o remoto para rodar o dry-run: vou pedir antes.**
- **Rollback:** revert do workflow.

### OPS-02 · Distribuir o binário MCP
- **Spec:** o instalador e a release incluem `proxy0x-mcp.exe` ao lado do `Proxy0x.exe`.
- **Teste:** instalar num diretório temporário → `GetMCPStatus` acha o binário (manual).
- **Rollback:** revert.

### OPS-03 · Instalador
- **Spec:**
  - `.iss` em UTF-8 com BOM, com textos nativos em ru/ja/ko/ar/th.
  - Versão vinda de `/DMyAppVersion` (padrão `0.0.0-dev`).
  - Aviso se o WebView2 estiver ausente.
  - A desinstalação pergunta se deve apagar `%APPDATA%\Proxy0x`.
- **Teste:** compilar com `iscc` e abrir o instalador em japonês e russo (verificação visual).
- **Rollback:** revert.

---

## Lote 6 — Segurança e concorrência média · ~7 commits

### CORE-09 · Snapshots em vez de ponteiros vivos
- **Spec:** `GetProxies` e todo `emit` entregam **cópias** feitas sob `a.mu`.
- **Teste:** depois de `GetProxies`, alterar o proxy interno não muda a cópia. Race sob estresse fica como verificação parcial, já que não há `-race` nesta máquina (registrado como limitação).
- **Rollback:** revert.

### SEC-02 · Senha não vai para o webview
- **Spec:** o snapshot da UI zera `Password` (o disco e o Go continuam com ela). O export mostra um aviso e o clipboard é limpo depois de 30 s.
- **Teste:** o snapshot emitido tem `password == ""`; a persistência continua com a senha (round-trip).
- **Rollback:** revert.

### SEC-03 · SSRF falha fechada
- **Spec:**
  - Erro de DNS → URL rejeitada.
  - O host é normalizado sem o ponto final.
  - Novas faixas bloqueadas: 100.64/10, 198.18/15, 240/4, 64:ff9b::/96, 2002::/16.
  - A doc deixa de prometer proteção contra rebinding (o DNS é resolvido pelo proxy remoto) e explica o limite.
- **Teste:** casos de tabela em `ssrf_test.go` (host inexistente, `metadata.google.internal.`, cada faixa nova).
- **Rollback:** revert.

### SEC-01 · Chave protegida por DPAPI
- **Passo 1 (refatoração, DRY):** `pkg/proxycore/store.go` passa a usar `internal/crypto.Store` (mesmo formato e mesma chave; é o mesmo conhecimento de negócio).
- **Passo 2 (correção):**
  - **Spec:** a chave fica em `.key.dpapi`, protegida com `CryptProtectData` (escopo do usuário). Na migração, a `.key` antiga é lida, protegida e conferida com um round-trip, e depois renomeada para `.key.legacy`. Chave ilegível vira **erro** (não é regenerada em silêncio).
  - **Teste:** round-trip de proteger/desproteger; migração de `.key` → `.key.dpapi` com as credenciais existentes ainda legíveis; chave corrompida → erro.
  - **Rollback:** revert. **Atenção:** versões antigas do app não leem `.key.dpapi`. Para voltar, renomeie `.key.legacy` → `.key`. Isso vai documentado em `security.md`. O `.key.legacy` é removido numa versão futura.

### PRIV-01 · Menos vazamento fora do proxy (D4)
- **Spec:** sem Google Fonts (fonte do sistema), com CSP `default-src 'self'; img-src 'self' https://flagcdn.com data:; style-src 'self' 'unsafe-inline'`. O README descreve o alcance real (apps que respeitam o proxy do sistema; UDP/WebRTC fora).
- **Teste:** `index.html` sem domínio externo além do flagcdn; o app abre sem erro de CSP no console (verificação manual no `wails dev`).
- **Rollback:** revert.

---

## Lote 7 — Frontend · ~8 commits

| Item | Spec / aceite | Teste |
|---|---|---|
| FE-01 | `activeGroup` volta para `''` quando sai de `groups`; nunca há lista vazia com proxies existentes | Manual: cria um grupo, remove o último proxy dele → lista volta a "All" |
| FE-02 | Toda ação que falha mostra a mensagem no card/settings por 5 s | Manual: connect num proxy inválido mostra o erro |
| FE-03 | O uptime avança a cada 1 s enquanto conectado; o intervalo é limpo ao desmontar | Manual |
| FE-04 | Excluir pede confirmação inline ("Excluir? Sim/Não") e faz `await` | Manual |
| FE-05 | O intervalo custom guarda uma string local e grava no blur/Enter; o vazio é permitido durante a digitação | Manual |
| FE-06 | "All", "Proxy … dropped", os rótulos de botão e "min" passam para `translations.ts` nos 16 idiomas; os botões da janela ganham `aria-label` | `tsc` (a interface de traduções obriga todas as chaves) |
| FE-07 | Links via `BrowserOpenURL`; cores fixas trocadas por tokens; máscara para IPv6 | Manual |

O frontend não tem harness de teste hoje. Adicionar Vitest seria uma **dependência nova**, então fica para o roadmap. Até lá a verificação é manual no `wails dev` e cada item fica registrado em `03-changes.md`. Rollback: revert por item.

---

## Lote 8 — Limpeza · ~5 commits

- **DEAD-01 (D5):** remove `ConnectGroup`, `DisconnectGroup`, `RestartProxy`, `GetConnectionLog`, `Get/SetHealthCheckInterval`, `QuitApp`, `ShowWindow`, o ring buffer e o evento `log-entry`, e `ContextMenu.tsx`. `CheckHealth` deixa de ser exportado. Os bindings são regenerados. **Aceite:** `tsc` OK, `grep` sem referências restantes, specs atualizadas.
- **ARCH-02:** remove `Score`/`Successes`/`Failures`; `BestProxy` → `ActiveProxy` (devolve o proxy ativo). As descrições das tools são corrigidas. Os testes do score saem junto com o recurso. O contrato das tools fica igual (o fetch continua usando o proxy ativo).
- **MCP-03:** health check com `CREATE_NO_WINDOW` e stdin aberto até a resposta, com timeout de 5 s. O exemplo passa a usar `<caminho do Proxy0x>`. O caminho do Codex é corrigido para `config.toml`, se confirmado na doc do Codex; senão, vira Hipótese registrada.
- **DOC-01 (restante):** o que não entrou nos lotes anteriores (`frontend.md`, `mcp.md`, README).

---

## Fora deste plano

- **ARCH-01** (usar `pkg/proxycore` no app inteiro): esforço G. Vai para `04-roadmap.md` se o modo EVOLUIR for ativado. O SEC-01 passo 1 já adianta a parte de crypto.
- **TEST-01:** coberto aos poucos pelos testes de cada item. Harness de frontend: roadmap.

## Estimativa

~50 commits pequenos em 9 lotes. Ordem: 0 → 1 → 2 → 3 → 4 → 6 → 7 → 8, com o **5** quando você liberar o push. Se a `main` andar no meio do caminho, sincronizo a branch antes de cada lote.

## Critério de conclusão

- Altos: todos corrigidos com teste. Médios: corrigidos ou adiados com justificativa.
- `go build`, `go vet`, `go test ./...`, `tsc` e `npm audit` verdes; `govulncheck` sem vulnerabilidade alcançável na stdlib.
- Nenhuma regressão contra `00-baseline.md`.
- Specs alinhadas ao código.
- `03-changes.md` com cada commit, teste e verificação manual feita ou pendente.
