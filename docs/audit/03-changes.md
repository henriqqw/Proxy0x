# 03 — Registro de mudanças (fase 3)

Branch: `fix/audit-2026-11` (worktree `.claude/worktrees/fix-audit`) · Base: `9df0d8c` · 38 commits · **Nada foi enviado ao remoto nem mesclado na `main`.**

## Resultado contra a linha de base

| Verificação | Linha de base (`00-baseline.md`) | Agora |
|---|---|---|
| `go build ./...` / `go vet ./...` | OK | OK |
| `go test ./...` | 2 pacotes com teste (mcpserver, proxycore) | 5 pacotes com teste, todos verdes: + raiz (`app`), `internal/proxy`, `internal/fsutil` |
| Testes novos | — | 37 funções de teste novas e 4 reescritas para o contrato novo (fora os subtests); cada correção tem teste que falhava antes (exceções abaixo) |
| `govulncheck` | 14 vulnerabilidades alcançáveis (stdlib) | **0 alcançáveis** |
| `npx tsc --noEmit` | OK | OK |
| `npm audit --omit=dev` | 0 | 0 |
| Bindings Wails (`wails generate module`) | — | conteúdo idêntico (só diferença de final de linha, descartada) |

Nenhuma regressão encontrada. Nenhum item precisou ser revertido.

## Itens tratados

| ID | Commit(s) | Teste que falhava antes → passa depois | Observação |
|---|---|---|---|
| SEC-04 | `4cd237b` | `govulncheck` 14 → 0 | Toolchain **go1.26.10** (o 1.26.6 do plano não bastava) e `x/net` v0.60.0, que exige a diretiva `go 1.26`. CI lê a versão do `go.mod`. |
| CORE-01 | `08be13b` | `TestStopReturnsQuicklyWithPlainHTTPInFlight` | Reproduzido: `Stop` preso > 1 s. |
| CORE-02 | `dbe983d` | `TestPlainHTTPClosesClientAfterResponse` (content-length e chunked) | Proxy local encaminha com `Connection: close`. |
| CORE-14 | `bb70513` | `TestAcceptErrorsBackOff` | Antes: 2,3 milhões de `Accept` em 100 ms. |
| CORE-11 | `61043d1`, `e37913c` | `TestSilentUpstreamTimesOut` (http, socks5), `TestSOCKS5DialHonorsContext`, `TestGetMyIPTimesOut` | O teste do proxy local travou 127 s com o código antigo. |
| CORE-03 | `49b58e2` | `TestGetOrCreateClientIsSafeForConcurrentReaders` | **Era Hipótese, agora Fato verificado**: `fatal error: concurrent map read and map write`. O teste equivalente em `pkg/proxycore` foi escrito depois da correção (não foi visto falhar). |
| CORE-06 | `3c94951` | `TestLocalServerStartFailureFailsConnect`, `TestRotateWithLocalServerFailureDoesNotPanicOnDisconnect` | Reproduzido: `panic: close of closed channel`. Depois do CORE-07 a rotação não recria mais o servidor, então o 2º teste passou a cobrir só "não há panic nem proxy sobrando". |
| CORE-04 | `62887fc` | `TestRemoveActiveProxyDisconnectsFirst` | |
| CORE-12 | `a82b112` | `TestConnectRequiresHTTP200AndValidIP` (5 casos), `TestHealthCheckFailsOnNon200` | As 3 cópias da checagem viraram `verifyExitIP`. |
| X-01 (novo) | `251a609` | `TestHealthCheckDoesNotHoldLockDuringGeoLookup` | **Achado durante a correção**: o health check fazia geolocalização (rede) com o lock de escrita preso; `GetProxies` ficou 350 ms bloqueado no teste. |
| X-02 (novo) | `4403e24` | `TestHealthCheckRecoversStatusAfterTransientFailure` | **Achado durante a correção**: depois de uma falha, o proxy ficava `error` para sempre. |
| CORE-13 | `57cc662` | `TestFailedRotationKeepsPreviousSession` | |
| CORE-15 | `ebc9684` | `TestSOCKS5ErrorNeverFallsBackToDirect` (app) | No `pkg/proxycore` a mesma correção foi **verificada só por leitura**: não há como forçar o erro do `x/net` sem um ponto de injeção. |
| CORE-08 | `6eec2b6` | — | Bandeja: **sem teste automatizado**; verificação manual pendente. |
| CORE-10 | `c0cb120` | `TestUnreadableStoreBecomesReadOnly` (2 casos), `TestSaveIsAtomicWithBackup`, testes de `fsutil` | **Achado extra X-03**: se a criptografia falhasse, a senha ia **em texto puro** para o disco; e uma credencial que não descriptografava era re-cifrada como se fosse texto. As duas coisas foram corrigidas no mesmo commit. |
| MCP-01 | `bbdff76` | `TestAddProxy0xRefusesToOverwriteInvalidJSON`, `TestAddAndRemoveProxy0xPreserveOtherKeysWithBackup` | Backup em `<arquivo>.proxy0x.bak`. |
| MCP-02 | `cc3eff8` | `TestPoolNeverWritesStore`, `TestHandleProxyAddIsReadOnly`, `TestHandleProxyRemoveIsReadOnly` | ⚠ **Contrato das tools mudou** (decisão D2): `proxy_add`/`proxy_remove` agora retornam erro. Os testes antigos desse contrato foram substituídos. |
| CORE-05 | `2d54ab1` (refatoração), `e2314c3` | 6 testes em `system_test.go` (5 falhavam) | O registro agora é acessado via `x/sys/windows/registry` (antes `reg.exe`). **Verificação manual pendente** no registro real. |
| CORE-07 | `cfd1e0d` | `TestUpdateUpstreamKeepsOpenTunnelsAndPort`, `TestRotationKeepsLocalServerAndRegistry` | |
| CORE-09 | `878b0c7` | `TestBindingsAndEventsReturnCopies` | Sem `-race` (não há cgo nesta máquina): o teste prova que os valores são cópias, não a ausência total de corridas. |
| SEC-02 | `501cd53` | `TestPasswordNeverReachesTheUI` | **Só a parte do backend.** O aviso e a limpeza do clipboard no export ficam para o lote de frontend. |
| SEC-03 | `328f278` (refatoração), `03f4a85` | `TestValidateURLFailsClosed` (13 casos; 9 falhavam) | Os testes de SSRF não dependem mais de DNS real. A doc agora assume que a checagem **não impede DNS rebinding** através do proxy. |
| ARCH-02 | `0f24ce0` | `TestPoolActiveProxyReturnsConnected` | O campo `score` saiu da saída de `proxy_list` e `proxy://pool` (mudança pequena de contrato de saída). |
| MCP-03 | `eb6001e` (refatoração), `fd1fe23`, `39c3f7e` | `TestCheckMCPHealthWaitsForTheResponse` | Reproduzido: falso "no response" com um servidor que encerra no EOF. |
| DOC-01 | `b2569f8`, `c3a1c44`, `8b5a8a7`, `8d9c8aa`, `703e317`, `4d76405` | — | Specs atualizadas **antes** de cada lote (SDD). |

## Adiados, com justificativa

| ID | Motivo |
|---|---|
| SEC-01 (DPAPI) | Implementado pelo **dev2** na `main` (`d6c1e0d`). Nesta branch: correção do vet e a chave compartilhada com o MCP (ver "Integração"). |
| CORE-16 (fail-closed) | Definido pelo kill switch e pelo auto failover do dev2. |
| Lote 7 — FE-01…FE-07, PRIV-01, parte de frontend do SEC-02 | O dev2 está editando o frontend e as traduções agora. Aplicar em paralelo garantiria conflito. Retomar depois que ele commitar. |
| DEAD-01 | O dev2 pode passar a usar bindings que hoje parecem mortos (logs/histórico, grupos, `RestartProxy` no failover). Reavaliar depois do commit dele. |
| ARCH-01 | Esforço G; vai para o roadmap (modo EVOLUIR). |
| Caminho do Codex (`config.json` → `config.toml`) | Exige escrever TOML (código novo ou dependência). Fica registrado como **Hipótese** para quem mantém a integração (dev2). |

## Integração com a `main` (features do dev2) e Lote 5

| Item | Commit | Teste | Observação |
|---|---|---|---|
| Merge `origin/main` | `3a9beb1`, `c47eb5e` | suíte completa | Único conflito textual: imports do `app.go`. |
| X-04 (novo, **crítico**) Servidor MCP sobrescrevia a chave DPAPI | `7f41f86` | `TestStoreReadsTheAppsDPAPIKeyWithoutReplacingIt` | Reproduzido: `cipher: message authentication failed`. O store do MCP tinha sua própria cópia do código de chave: via o `.key` DPAPI como "inválido", gerava uma chave nova e gravava por cima, e o app perdia todas as senhas salvas. Agora o MCP usa `internal/crypto` (passo 1 do SEC-01 / ARCH-01 para crypto). |
| `go vet` quebrado pelo DPAPI | `3974286` | `TestDPAPIRoundTrip`, `TestStoreEncryptDecryptWithDPAPIKey` | `possible misuse of unsafe.Pointer`. Trocado pelos wrappers do `x/sys/windows`, mesmo formato de chave. |
| X-05 Features novas contornavam as correções | `de56583` | `TestSetProxyNameReturnsCopyAndRespectsReadOnly`, `TestAddProxiesMultiFormatRespectsReadOnly`, `TestBulkCheckRequiresValidExitIP` | Checagens de IP sem status (um 429 contava como vivo), `http.Get` sem timeout, `runtime.EventsEmit` direto (derruba o processo fora do Wails), ponteiros vivos com senha, modo só leitura ignorado. |
| OPS-01/02/03 | `0d34283` | ISCC local gerou `Proxy0x-Setup-0.0.0-test.exe` (7 MB) | Release só por tag `v*`. Actions fixadas por SHA, Wails v2.16.0, instalador Inno com `proxy0x-mcp.exe`. Textos nativos (UTF-8 com BOM). `SourceDir=..`, porque o `.iss` foi movido para `build/` e todos os caminhos relativos quebraram. O workflow em si **ainda não rodou no GitHub**; roda na primeira tag. |

### Riscos nas features novas (não corrigidos; decisão de produto)

- **Kill switch:** cria uma regra de *bloqueio* de saída mais regras de *permissão* no Windows Firewall. Pela documentação do firewall, regras de bloqueio têm precedência sobre as de permissão, então o kill switch pode bloquear a internet inteira, inclusive o próprio Proxy0x. **Hipótese não testada**: testar exige mexer no firewall. O script elevado é gravado num caminho fixo do `%TEMP%` e apagado logo depois do `ShellExecute`, que é assíncrono: pode ser apagado antes de rodar, ou trocado por outro processo nesse intervalo. Se o app morrer com o kill switch ligado, as regras ficam no sistema.
- **"DNS leak test":** compara o IP de saída do proxy com o IP direto. Não testa o resolvedor DNS.
- **Import:** `socks4://` é tratado como SOCKS5.
- **Integração Codex:** escreve em `~/.codex/config.json`, mas o Codex usa `config.toml` (Hipótese).

## Verificação manual pendente (não feita nesta fase)

1. `wails dev` na branch: conectar, rotacionar (a porta do proxy local deve continuar a mesma), parar (deve ser instantâneo) e remover o proxy ativo.
2. Navegar num site **HTTP** (sem TLS) com várias imagens: deve carregar inteiro (CORE-02).
3. Com um proxy corporativo configurado no Windows: conectar e desconectar no Proxy0x → o proxy corporativo volta (CORE-05).
4. Bandeja: "Disconnect" e "Rotate" (CORE-08).
5. Settings → MCP → health check: sem janela de console piscando (MCP-03).

## Riscos residuais

- **Merge com o trabalho do dev2:** os dois mexem no `app.go`, ele com métodos novos e eu com os existentes. Mudanças de semântica que ele precisa respeitar (já avisado): `setLocalProxy` retorna erro; `verifyExitIP` é a checagem única; o mapa de clients tem mutex próprio; eventos via `a.emitEvent`; bindings devolvem cópias sem senha.
- `RemoveProxy` e `ReorderProxies` agora retornam `error` (em TS continua `Promise<void>`). O frontend ainda não trata a rejeição no modo só leitura. Isso entra no lote de frontend.
- O evento `store-error` ainda não tem ouvinte na UI, então o modo só leitura fica invisível para o usuário até o lote de frontend.
- Sem detector de race nesta máquina.
- As releases automáticas v1.0.5–v1.0.14 continuam publicadas (nada foi apagado); a próxima tag passa a ser a "Latest".
