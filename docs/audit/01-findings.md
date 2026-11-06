# 01 — Achados da auditoria

Commit auditado: `0d5768e` · Modo desta execução: **AUDITAR** (campo de modo veio em branco; nenhuma linha de código foi alterada)

## Resumo

| Severidade | Qtde |
|---|---|
| Crítica | 0 |
| Alta | 6 |
| Média | 19 |
| Baixa | 14 |

**O que mais importa, em uma frase cada:**

1. **CORE-01/02** — O proxy local trava o "Stop" por até 5 min e não atende a 2ª requisição de uma conexão keep-alive em sites HTTP. **Reproduzido** com teste.
2. **CORE-03** — Uma escrita em mapa sob `RLock` pode derrubar o app com `fatal error: concurrent map writes`. A janela é o health check rodando junto com um speed test. Depois do crash, o proxy do sistema fica ligado e o usuário fica sem internet.
3. **MCP-01** — "Instalar MCP" apaga todas as preferências do Claude Desktop se o JSON não parsear. O WIP atual repete o padrão para 4 clientes.
4. **MCP-02** — O app e o servidor MCP sobrescrevem o `proxies.json` um do outro.
5. **OPS-01** — Todo push na `main` publica uma release para usuários. A "Latest" hoje é a v1.0.11 (NSIS), não o instalador documentado.

Classificação usada: **Fato verificado** (visto no código ou reproduzido) · **Hipótese** (mecanismo provável, não reproduzido) · **Suposição** (depende de decisão ou intenção).

## Índice

| ID | Título | Sev. | Tipo | Classificação |
|---|---|---|---|---|
| CORE-01 | Stop/Disconnect/encerrar travam até 5 min (conexão upstream HTTP não rastreada) | Alta | Bug | Fato verificado (reproduzido) |
| CORE-02 | Proxy local atende só 1 requisição por conexão; keep-alive trava | Alta | Bug | Fato verificado (reproduzido) |
| CORE-03 | Escrita concorrente em mapa sob RLock → crash do processo | Alta | Concorrência | Fato (código) / crash = Hipótese |
| MCP-01 | Instalar MCP apaga o config do Claude Desktop se o JSON falhar | Alta | Perda de dados | Fato verificado |
| MCP-02 | App e MCP sobrescrevem o `proxies.json` um do outro | Alta | Concorrência / perda de dados | Fato (mecanismo) |
| OPS-01 | CI publica release a cada push na main; "Latest" ≠ artefato documentado | Alta | Operação / supply chain | Fato verificado |
| CORE-04 | Remover o proxy ativo mantém o tráfego passando por ele | Média | Bug | Fato verificado |
| CORE-05 | Config de proxy do usuário sobrescrita, não restaurada e forçada a OFF | Média | Bug / perda de config | Fato verificado |
| CORE-06 | Falha ao subir o LocalServer deixa ponteiro para servidor parado → panic | Média | Bug | Fato verificado |
| CORE-07 | Cada rotação reinicia o proxy local e derruba todas as conexões | Média | Robustez | Fato verificado |
| CORE-08 | "Connect All" da bandeja dispara conexões concorrentes num modelo de 1 ativo | Média | Bug / lógica | Fato (código) |
| CORE-09 | `*ProxyEntry` serializado enquanto é mutado (data race) | Média | Concorrência | Hipótese |
| CORE-10 | Escritas não atômicas e erros ignorados → arquivo truncado → lista vazia | Média | Robustez / perda de dados | Fato (código) |
| SEC-01 | Chave AES fica ao lado do texto cifrado | Média | Segurança | Fato verificado |
| SEC-02 | Senhas em texto puro vão para o webview e para o clipboard | Média | Segurança | Fato verificado |
| SEC-03 | Guarda SSRF falha aberta em erro de DNS; proteção contra rebinding não existe | Média | Segurança | Fato verificado |
| SEC-04 | 14 vulnerabilidades alcançáveis na stdlib (go1.26.5) | Média | Dependências | Fato (govulncheck) |
| PRIV-01 | Vazamentos do IP real fora do proxy (Google Fonts, flagcdn, UDP/WebRTC) | Média | Privacidade | Fato / Hipótese |
| FE-01 | Filtro de grupo vira beco sem saída (lista vazia sem botão de voltar) | Média | Bug | Fato verificado |
| FE-02 | Erros engolidos em toda a UI (`catch {}`) | Média | Robustez / UX | Fato verificado |
| OPS-02 | Binário MCP não é distribuído; integração sempre falha para quem instala | Média | Empacotamento | Fato verificado |
| ARCH-01 | Núcleo duplicado entre app e `pkg/proxycore`, já divergindo | Média | DRY | Fato verificado |
| ARCH-02 | Score/BestProxy do pool MCP é lógica morta; contadores perdidos ao salvar | Média | Lógica / over-engineering | Fato verificado |
| TEST-01 | Núcleo do app desktop sem nenhum teste; lacunas nos testes MCP | Média | Testes | Fato verificado |
| DOC-01 | Specs (`docs/`) contradizem o código e se contradizem | Média | Documentação / SDD | Fato verificado |
| CORE-11 | Sem timeout: `GetMyIP`, dial SOCKS5 | Baixa | Robustez | Fato verificado |
| CORE-12 | Status HTTP e JSON não checados → "conectado" com IP vazio | Baixa | Robustez | Fato verificado |
| CORE-13 | Falha na rotação deixa estado inconsistente | Baixa | Bug | Fato verificado |
| CORE-14 | `Accept()` com erro entra em loop ocupado (100% CPU) | Baixa | Robustez | Fato verificado |
| CORE-15 | Erro no SOCKS5 cai para conexão direta (vazaria o IP real) | Baixa | Privacidade | Hipótese (alcançabilidade) |
| CORE-16 | Health check falho mantém o sistema apontando para upstream morto | Baixa | Comportamento | Suposição |
| FE-03 | Contador de uptime não avança | Baixa | Bug | Fato verificado |
| FE-04 | Excluir sem confirmação e sem `await` | Baixa | UX | Fato verificado |
| FE-05 | Input de intervalo custom salva a cada tecla e não pode ser apagado | Baixa | UX | Fato verificado |
| FE-06 | Textos fixos em inglês e botões sem rótulo acessível | Baixa | i18n / a11y | Fato verificado |
| FE-07 | Links `target="_blank"` no Wails; cores fixas; máscara de IP só IPv4 | Baixa | Bug / design system | Fato / Hipótese |
| DEAD-01 | Bindings, evento e componente sem uso | Baixa | Código morto | Fato verificado |
| OPS-03 | Instalador: textos romanizados, versão fixa, sem checagem WebView2, deixa credenciais | Baixa | Empacotamento | Fato verificado |
| MCP-03 | Health check do MCP: janela de console e stdin fechado cedo; exemplo com caminho de dev | Baixa | UX / higiene | Hipótese |

---

## Alta

### CORE-01 · Stop/Disconnect/encerrar travam até 5 min — Alta · Bug
- **Local:** `internal/proxy/server.go:250-292` (`handleHTTP`), `:62-73` (`Stop`)
- **Evidência:** em `handleHTTP` o `upConn` é aberto (`:263`, `:272`) mas **não** passa por `trackConn` (só o CONNECT rastreia, `:185`, `:194`). `Stop()` fecha só as conexões rastreadas e depois faz `wg.Wait()`. A goroutine fica presa em `io.Copy(clientConn, upConn)` (`:292`) até o `SetReadDeadline` de 5 min.
- **Prova:** teste no scratchpad com upstream falso keep-alive → `Stop still blocked after 5s`.
- **Impacto:** o botão "Stop" fica girando, o Disconnect pelo tray trava e o fechamento do app atrasa. É o mesmo sintoma já corrigido para HTTPS, voltando por sites HTTP. `clearLocalProxyIf` segura `lpMu` durante a espera e bloqueia connect/rotate.
- **Correção:** rastrear `upConn` em `handleHTTP` nos dois ramos (mesmo padrão de `handleConnect`).
- **Esforço:** P · **Risco:** Baixo

### CORE-02 · Só a 1ª requisição de cada conexão é atendida — Alta · Bug
- **Local:** `internal/proxy/server.go:159-169` (`handleConn` lê um único request), `:250-292`
- **Evidência:** `http.ReadRequest(br)` é chamado uma vez. A requisição é encaminhada sem `Connection: close` e o upstream mantém a conexão aberta, então `io.Copy` não retorna. O navegador reaproveita o socket para a próxima requisição e ninguém a lê.
- **Prova:** teste no scratchpad → `second request on same keep-alive connection got no response: i/o timeout`.
- **Impacto:** sites HTTP sem TLS carregam a primeira resposta e travam nos recursos seguintes. Também alimenta o CORE-01.
- **Correção (opções):** (a) forçar `Connection: close` na requisição encaminhada e fechar depois da resposta. É a mais simples e a recomendada. (b) Loop lendo várias requisições de `br` com `http.ReadResponse` delimitando cada resposta.
- **Esforço:** P · **Risco:** Médio (muda o comportamento do proxy HTTP; precisa de teste)

### CORE-03 · Escrita concorrente em mapa sob RLock — Alta · Concorrência
- **Local:** `app.go:535` (`a.clients[clientKey] = client` dentro de `getOrCreateClient`), chamado sob `a.mu.RLock()` em `CheckHealth` (`app.go:747-755`) e `BenchmarkProxy` (`app.go:939-946`). Mesmo padrão em `pkg/proxycore/pool.go:344-351`.
- **Evidência:** `RLock` permite vários leitores ao mesmo tempo, e os dois caminhos escrevem no mapa.
- **Impacto:** se o health check (goroutine a cada 3 min) coincidir com um Speed Test, o runtime do Go aborta com `fatal error: concurrent map writes`. Isso não é recuperável. Como não há limpeza no crash, o proxy do sistema continua apontando para uma porta morta e o usuário fica sem internet até reabrir o app.
- **Correção:** criar o client sob `Lock`, ou usar um `clientsMu` dedicado.
- **Esforço:** P · **Risco:** Baixo · **Classificação:** escrita sob RLock = Fato verificado. Crash = Hipótese (o `-race` não roda nesta máquina; dá para provar com um teste de estresse).

### MCP-01 · "Instalar MCP" apaga o config do Claude Desktop — Alta · Perda de dados
- **Local:** `mcp_integration.go:76-101` (HEAD)
- **Evidência:** `json.Unmarshal(data, &cfg)` tem o erro ignorado. Se falhar, segue com `cfg = make(map[string]any)` e `os.WriteFile(configPath, out, 0o644)`.
- **Impacto:** esse arquivo guarda **todas as preferências do Claude Desktop** (~95 KB nesta máquina), não só `mcpServers`. Um parse falho (arquivo sendo escrito pelo Claude no mesmo instante, BOM, comentário) faz o arquivo inteiro ser trocado por `{"mcpServers":{...}}`. Contradiz `docs/mcp.md` ("Preserves any existing entries"). O WIP não commitado repete o padrão para Codex, Cursor e Windsurf.
- **Correção:** abortar em erro de parse, criar `.bak` antes de escrever, gravar com arquivo temporário + rename.
- **Esforço:** P · **Risco:** Baixo

### MCP-02 · App e MCP sobrescrevem o `proxies.json` um do outro — Alta · Concorrência / perda de dados
- **Local:** `pkg/proxycore/store.go:16-19, 143-170`; `pkg/proxycore/pool.go:114-135`; `app.go:236-258`
- **Evidência:** o `Store.mu` só vale dentro do processo. Cada processo carrega o arquivo uma vez e depois regrava sua cópia inteira em memória.
- **Impacto:** um proxy adicionado no app some quando o MCP salva, e vice-versa. Como a escrita não é atômica (CORE-10), duas escritas simultâneas podem deixar o JSON truncado, e aí o app carrega lista vazia.
- **Correção (opções):** (a) o app vira o único dono do arquivo e o MCP fala com ele (IPC). É a mais limpa, porém maior. (b) Lock de arquivo entre processos + recarregar, mesclar e salvar + rename atômico. (c) MCP só leitura. É a mais simples, mas perde o `proxy_add`. Recomendo **(c) agora e (a) depois**.
- **Esforço:** M · **Risco:** Médio · **Classificação:** mecanismo = Fato verificado; intercalação real = Hipótese

### OPS-01 · CI publica release a cada push na `main` — Alta · Operação / supply chain
- **Local:** `.github/workflows/release.yml:3-5` (`on: push: branches: [main]`), `:66` (`wails@latest`), `:83-103`
- **Evidência:** foram 7 releases automáticas em 1 dia (v1.0.5–v1.0.11), inclusive para commits só de docs. A v1.0.0 manual foi publicada **depois** da v1.0.7. A "Latest" é a v1.0.11 com `Proxy0x-win64-installer.exe` (NSIS), enquanto o README documenta `Proxy0x-Setup-1.0.0.exe` (Inno Setup). O CI também gera 386 e arm64, que nunca foram testados. `wails@latest` e as actions por tag não estão fixados.
- **Impacto:** build não revisado chega direto a usuários, os artefatos não batem com a documentação, e o versionamento fica confuso (a v1.0.0 "mais nova" que a v1.0.7). Há risco de cadeia de suprimentos com o `@latest`.
- **Correção:** disparar só por push de tag `v*` ou `workflow_dispatch`, fixar a versão do Wails, decidir um instalador canônico (Inno **ou** NSIS) e gerá-lo no CI, e ajustar o README.
- **Esforço:** M · **Risco:** Médio (muda o processo de release; **decisão sua**)

---

## Média

### CORE-04 · Remover o proxy ativo mantém o tráfego passando por ele — Média · Bug
- **Local:** `app.go:351-363`
- **Evidência:** `RemoveProxy` apaga o proxy do mapa sem chamar `DisconnectProxy`. Além disso, `delete(a.clients, id)` (`:359`) usa a chave errada: a chave real é `id + ":" + sessionID` (`:488`).
- **Impacto:** o proxy local e o registro continuam roteando pelo proxy excluído, sem card na UI para parar. O transporte HTTP vaza.
- **Correção:** se estiver ativo, desconectar antes; fechar o client pela chave correta.
- **Esforço:** P · **Risco:** Baixo

### CORE-05 · Config de proxy do usuário é sobrescrita e não restaurada — Média · Bug / perda de configuração
- **Local:** `internal/proxy/system.go:15-28`; `app.go:104` (`ClearSystemProxy()` no startup)
- **Evidência:** `SetSystemProxy` sobrescreve `ProxyServer`/`ProxyOverride` sem guardar os valores anteriores. `ClearSystemProxy` força `ProxyEnable=0`, e isso também roda **a cada inicialização**.
- **Impacto:** quem usa proxy corporativo ou de outra ferramenta perde a configuração só por abrir o app. Somado ao CORE-03, um crash deixa o sistema apontando para uma porta morta.
- **Correção:** salvar os 3 valores anteriores ao conectar e restaurar ao desconectar. No startup, só limpar se `ProxyServer` apontar para `127.0.0.1` (resíduo nosso).
- **Esforço:** P · **Risco:** Baixo

### CORE-06 · Falha ao subir o LocalServer → panic no próximo Stop — Média · Bug
- **Local:** `app.go:560-575`; `internal/proxy/server.go:63`
- **Evidência:** `a.localProxy.Stop()` roda, mas `a.localProxy` só é substituído se `NewLocalServer` der certo. Em erro, o ponteiro antigo (já parado) fica guardado. O próximo `Stop()` faz `close(ls.stopCh)` de novo → `panic: close of closed channel`.
- **Impacto:** crash raro, já que exige falha de `net.Listen` em 127.0.0.1. O erro também é engolido e o usuário vê "connected" sem proxy de sistema.
- **Correção:** zerar `a.localProxy` logo depois do `Stop()` e propagar o erro.
- **Esforço:** P · **Risco:** Baixo

### CORE-07 · Cada rotação reinicia o proxy local — Média · Robustez
- **Local:** `app.go:733` (`RotateProxy` → `setLocalProxy`); `app.go:560-575`
- **Evidência:** a rotação para o servidor e sobe outro numa porta nova, regravando o registro.
- **Impacto:** com a auto-rotação em 1 min, todas as conexões do navegador (downloads, streams, websockets) caem a cada minuto.
- **Correção:** `LocalServer.UpdateUpstream(cfg)` sob o `mu` que já existe, mantendo a porta e o registro.
- **Esforço:** P/M · **Risco:** Médio

### CORE-08 · "Connect All" da bandeja contradiz o modelo de 1 proxy ativo — Média · Bug / lógica
- **Local:** `tray.go:70-82` (`go a.ConnectProxy(id)` na `:80`)
- **Evidência:** `ConnectProxy` desconecta qualquer outro ativo (`app.go:588-598`), mas o tray dispara N conexões em paralelo.
- **Impacto:** as conexões competem, o estado final fica indeterminado, o proxy local troca várias vezes e pode sobrar mais de um card marcado como ativo. O mesmo vale para `ConnectGroup` (que não tem uso).
- **Correção:** remover "Connect All" (e `ConnectGroup`) ou trocar por "Conectar o primeiro".
- **Esforço:** P · **Risco:** Baixo · **Classificação:** código = Fato; efeito exato = Hipótese

### CORE-09 · `*ProxyEntry` serializado enquanto é mutado — Média · Concorrência
- **Local:** `app.go:322, 612, 627, 653, 699, 735, 794` (`EventsEmit(..., p)` com o ponteiro vivo); `app.go:366` (`GetProxies` devolve ponteiros); `pkg/proxycore/pool.go:137-155`, `fetcher.go:103-114`
- **Evidência:** os campos mudam sob `a.mu` enquanto o Wails faz o marshal do mesmo ponteiro fora do lock.
- **Impacto:** data race, ou seja, JSON inconsistente no frontend. Na prática, um status de um instante e um IP de outro.
- **Correção:** emitir e devolver cópias feitas sob o lock.
- **Esforço:** M · **Risco:** Baixo · **Classificação:** Hipótese (sem `-race` nesta máquina)

### CORE-10 · Escritas não atômicas e erros ignorados — Média · Robustez / perda de dados
- **Local:** `app.go:257`; `internal/config/config.go:53`; `pkg/proxycore/store.go:170`; `mcp_integration.go:101,130`; `cmd/proxy0x-mcp/main.go:25-27`; `pkg/proxycore/pool.go:166,182`
- **Evidência:** `os.WriteFile` direto, com o retorno ignorado. O `loadProxies` (`app.go:210-216`) devolve lista vazia em silêncio se o JSON estiver corrompido. O MCP segue com pool vazio depois de um erro de load.
- **Impacto:** um crash ou queda de energia no meio de uma escrita → JSON truncado → no próximo start a lista vem vazia → o próximo save **grava a lista vazia por cima**, e todos os proxies se perdem.
- **Correção:** escrever num temporário + rename, manter um `.bak`, recusar salvar quando o load falhou, e propagar erros para a UI.
- **Esforço:** P · **Risco:** Baixo · **Classificação:** código = Fato; perda real = Hipótese

### SEC-01 · Chave AES fica ao lado do texto cifrado — Média · Segurança
- **Local:** `internal/crypto/crypto.go:22-43`; `pkg/proxycore/store.go` (mesmo esquema)
- **Evidência:** a chave fica em `%APPDATA%/Proxy0x/.key` e o texto cifrado em `%APPDATA%/Proxy0x/proxies.json`. O `0600` não tem efeito real em NTFS.
- **Impacto:** qualquer coisa que leia a pasta (malware, backup, sincronização do OneDrive) leva os dois. A "criptografia em repouso" que README e specs prometem protege muito pouco. Se a chave corromper, ela é regenerada em silêncio (`:31-41`) e as credenciais ficam ilegíveis.
- **Correção:** proteger a chave com DPAPI (`CryptProtectData`, escopo do usuário, via `golang.org/x/sys/windows`, que já é dependência indireta), migrando a chave antiga. App e MCP precisam mudar juntos.
- **Esforço:** M · **Risco:** Médio (migração de dados)

### SEC-02 · Senhas em texto puro no webview e no clipboard — Média · Segurança
- **Local:** `frontend/wailsjs/go/models.ts:80` (`password: string`); `SettingsPage.tsx:52-55, 147-151` (HEAD)
- **Evidência:** todo `GetProxies` e todo evento `proxy-updated` levam a senha para o JS. O export coloca `host:port:user:pass` num textarea e no clipboard, sem aviso.
- **Impacto:** amplia a superfície de exposição (o histórico e a sincronização do clipboard do Windows guardam o texto). O card não renderiza a senha.
- **Correção:** DTO sem senha para a UI, aviso antes de exportar, limpar o clipboard depois de N segundos.
- **Esforço:** M · **Risco:** Médio (o Go continua precisando da senha)

### SEC-03 · Guarda SSRF falha aberta — Média · Segurança
- **Local:** `pkg/proxycore/ssrf.go:74-75`, `:27-37`; `fetcher.go:79`; `docs/security.md` ("prevents DNS rebinding")
- **Evidência:** `ips, err := net.LookupIP(host); if err == nil {...}`: se o DNS local falhar, a URL passa. Um ponto final em `metadata.google.internal.` escapa da checagem de sufixo. A validação resolve uma vez e a conexão resolve de novo (TOCTOU), então a proteção contra rebinding não existe. Faltam 100.64/10, 198.18/15, 240/4, 64:ff9b::/96 e 2002::/16.
- **Impacto:** é possível contornar a guarda. O impacto real é menor porque a saída passa pelo proxy remoto, não pela rede local.
- **Correção:** rejeitar em erro de lookup, normalizar o ponto final, adicionar as faixas e corrigir a doc (ou fixar o IP validado no dial).
- **Esforço:** P · **Risco:** Baixo

### SEC-04 · 14 vulnerabilidades alcançáveis na stdlib — Média · Dependências
- **Local:** toolchain go1.26.5; `release.yml:56` (`go-version: '1.25'`, versão flutuante)
- **Evidência:** `govulncheck` aponta GO-2026-6090 (`crypto/tls`), GO-2026-5026 (`net/http`), GO-2026-5972 (`encoding/asn1`) e outras 11, todas corrigidas em go1.26.6.
- **Impacto:** o app faz TLS através de proxies de terceiros, um cenário com potencial de MITM ativo.
- **Correção:** `toolchain go1.26.6` no `go.mod` e a mesma versão fixa no CI.
- **Esforço:** P · **Risco:** Baixo

### PRIV-01 · Vazamentos do IP real fora do proxy — Média · Privacidade
- **Local:** `frontend/index.html:8-10` (Google Fonts); `FlagIcon.tsx:11` (flagcdn.com); README ("Route all browser traffic")
- **Evidência:** recursos remotos carregados em toda inicialização, sem CSP. O proxy de sistema (WinINET) cobre só apps que o respeitam e não cobre UDP/WebRTC.
- **Impacto:** o Google e o flagcdn veem o IP real antes de conectar, e o WebRTC pode expor o IP real com o proxy ligado. A promessa do README é maior do que o app entrega.
- **Correção:** empacotar a fonte e as bandeiras localmente, adicionar CSP, e ajustar o README para "tráfego HTTP/HTTPS de apps que usam o proxy do sistema".
- **Esforço:** M · **Risco:** Baixo · **Classificação:** recursos remotos = Fato. WebRTC/UDP = Hipótese (comportamento do navegador).

### FE-01 · Filtro de grupo vira beco sem saída — Média · Bug
- **Local:** `frontend/src/App.tsx:122, 141, 166`
- **Evidência:** a barra de grupos só aparece com `groups.length > 0`, mas `activeGroup` continua setado depois que o último proxy do grupo sai.
- **Impacto:** a lista fica vazia mesmo com proxies cadastrados, sem botão "All". Só reiniciando o app.
- **Correção:** resetar `activeGroup` quando ele não estiver mais em `groups`.
- **Esforço:** P · **Risco:** Baixo

### FE-02 · Erros engolidos em toda a UI — Média · Robustez / UX
- **Local:** `ProxyCard.tsx:71`; `SettingsPage.tsx:32, 36, 48, 58, 69`; `TitleBar.tsx:18, 32`
- **Evidência:** `try { await fn() } catch {}`
- **Impacto:** falhas de speed test, import, export, config e MCP não dão nenhum retorno ao usuário.
- **Correção:** mostrar a mensagem inline ou num toast.
- **Esforço:** M · **Risco:** Baixo

### OPS-02 · Binário MCP não é distribuído — Média · Empacotamento
- **Local:** `installer.iss:82` (só `Proxy0x.exe`); `release.yml` (não gera `proxy0x-mcp.exe`); `mcp_integration.go` (`mcpExePath()` espera o binário ao lado do app)
- **Impacto:** o botão de integração MCP nas Settings sempre falha com "binary not found" para quem instalou o app.
- **Correção:** compilar `cmd/proxy0x-mcp` e incluir no instalador e no CI, ou esconder o recurso até lá.
- **Esforço:** P · **Risco:** Baixo

### ARCH-01 · Núcleo duplicado entre app e `pkg/proxycore` — Média · DRY
- **Local:** `pkg/proxycore/pool.go:78-112`, `client.go`, `store.go:37-111` × `app.go` / `internal/crypto`
- **Evidência:** o parse, o sufixo `_session-`, o formato AES e o client HTTP estão copiados e já divergem (o timeout de geo só existe em uma das cópias).
- **Impacto:** toda correção deste relatório precisa ser feita duas vezes, e o formato em disco compartilhado pode divergir.
- **Correção:** o app passa a usar `pkg/proxycore`. É uma refatoração separada, **depois** das correções.
- **Esforço:** G · **Risco:** Médio

### ARCH-02 · Score/BestProxy do pool MCP é lógica morta — Média · Lógica / over-engineering
- **Local:** `pkg/proxycore/pool.go:204-215, 418-440`; `internal/mcpserver/tools.go:16`
- **Evidência:** `Connect` desconecta os outros, então só existe 0 ou 1 candidato. O `ProxyEntry` do app não tem `Successes`/`Failures`, e cada save do app apaga esses contadores.
- **Impacto:** as descrições das tools e a doc prometem um "pool rotativo por score" que não existe.
- **Correção:** remover o score (YAGNI) ou permitir vários ativos. Recomendo remover.
- **Esforço:** M · **Risco:** Médio

### TEST-01 · Núcleo do app sem testes — Média · Testes
- **Local:** `internal/proxy`, `internal/crypto`, `internal/config`, `app.go`, frontend: 0 testes. MCP: sem httptest de fetch, sem teste de redirect, de DNS que falha aberto, de ponto final ou de concorrência no store. `mcp_integration` escreve no `%APPDATA%` real, sem injeção de caminho.
- **Impacto:** os dois bugs Altos do proxy local passaram sem nenhum teste pegar.
- **Correção:** começar pelos testes de regressão de CORE-01/02/03 (já escritos no scratchpad) e por round-trip de crypto e config.
- **Esforço:** M · **Risco:** Baixo

### DOC-01 · Specs contradizem o código — Média · Documentação / SDD
- **Evidência (amostra):**
  - `business-rules.md:42-47`: o connect configura o proxy de sistema **antes** de verificar e usa `ip-api.com`. O código verifica primeiro e usa `ipwho.is` (`backend.md` está correto).
  - `business-rules.md:17`: SessionID de 8 bytes; o código usa 6 (`app.go:269`).
  - `business-rules.md:34`: timeout de 10 s; o código usa 20 s (`app.go:533`).
  - `business-rules.md:99-102`: arquivo de 10 MB e chaves `connectTime`/`downloadSpeed`; o código baixa 100 KB e devolve `connectMs`/`downloadMbps`.
  - `security.md:99`: "net.Conn wrapper intercepts Close()" não existe, e "no hanging connections" é falso (CORE-01).
  - `backend.md:191-194`: tray "with flag icons" e "event listener", que não existem.
  - `mcp.md`: parâmetro `headers` inexistente; argumentos dos prompts com outros nomes.
  - README: "live uptime timer" (FE-03) e "all browser traffic" (PRIV-01).
- **Impacto:** sob SDD a spec é a fonte da verdade, e hoje ela está errada. Qualquer correção feita "pela spec" introduziria bugs.
- **Correção:** atualizar as specs para o comportamento correto **antes** da fase 3 (é parte do plano).
- **Esforço:** M · **Risco:** Baixo

---

## Baixa

### CORE-11 · Sem timeout — Baixa · Robustez
- **Local:** `app.go:981` (`http.Get` com o client padrão, sem timeout); `internal/proxy/server.go:146` e `app.go` (ramo SOCKS5 do `getOrCreateClient`: `dialer.Dial` sem timeout nem ctx); `pkg/proxycore/client.go:50-52`
- **Impacto:** o badge "Your IP" e os dials SOCKS5 podem travar indefinidamente.
- **Correção:** client com `Timeout`; `netproxy.FromURL`/`ContextDialer` com `net.Dialer{Timeout}`. · **Esforço:** P · **Risco:** Baixo

### CORE-12 · Status HTTP e JSON não checados — Baixa · Robustez
- **Local:** `app.go` em Connect, Rotate e CheckHealth (`json.Unmarshal(body, &ipResult)` sem checar `resp.StatusCode` nem o erro)
- **Impacto:** um 429/5xx do ipify marca o proxy como "conectado" com IP vazio e geo nula.
- **Correção:** exigir 200 e um IP válido. · **Esforço:** P · **Risco:** Baixo

### CORE-13 · Falha na rotação deixa estado inconsistente — Baixa · Bug
- **Local:** `app.go:694, 706-711`
- **Evidência:** o `SessionID` muda antes da verificação. Se ela falhar, o status vira `error`, mas o proxy local continua na sessão antiga e o tráfego segue funcionando.
- **Correção:** só trocar a sessão depois de verificar; em falha, manter a anterior. · **Esforço:** P · **Risco:** Baixo

### CORE-14 · `Accept()` com erro → loop ocupado — Baixa · Robustez
- **Local:** `internal/proxy/server.go:76-85`
- **Evidência:** qualquer erro que não seja de parada cai em `continue`, sem backoff (ex.: EMFILE).
- **Correção:** backoff exponencial curto, como no `net/http.Server`. · **Esforço:** P · **Risco:** Baixo

### CORE-15 · Erro no SOCKS5 cai para conexão direta — Baixa · Privacidade
- **Local:** `app.go` (ramo SOCKS5 de `getOrCreateClient`: `if err != nil { dialer = netproxy.Direct }`); `pkg/proxycore/client.go:45-48`
- **Impacto:** se alcançável, o "IP do proxy" exibido seria o IP real do usuário. Hoje é latente, porque o `x/net` praticamente não retorna esse erro.
- **Correção:** retornar o erro. · **Esforço:** P · **Risco:** Baixo · **Classificação:** Hipótese (alcançabilidade)

### CORE-16 · Health check falho mantém upstream morto — Baixa · Comportamento
- **Local:** `app.go:746-775`
- **Evidência:** em falha, o app marca `error` e notifica, mas o proxy do sistema continua ativo.
- **Impacto:** o usuário fica sem internet até agir. Pode ser intencional (fail-closed protege a privacidade), mas não está documentado.
- **Correção:** **decisão sua**: manter fail-closed e documentar, ou oferecer o auto-disconnect. · **Esforço:** P · **Risco:** Baixo · **Classificação:** Suposição

### FE-03 · Uptime congelado — Baixa · Bug
- **Local:** `ProxyCard.tsx:24-37, 216-221` (não há `setInterval` em `src/`)
- **Impacto:** o contador só muda a cada `proxy-updated` (~3 min), o que contradiz `docs/frontend.md:159` e o README.
- **Correção:** tick de 1 s enquanto conectado. · **Esforço:** P · **Risco:** Baixo

### FE-04 · Excluir sem confirmação e sem `await` — Baixa · UX
- **Local:** `ProxyCard.tsx:125`
- **Impacto:** um clique errado apaga o proxy, inclusive o ativo (ver CORE-04).
- **Correção:** `await` + confirmar ou desfazer. · **Esforço:** P · **Risco:** Baixo

### FE-05 · Input de intervalo custom — Baixa · UX
- **Local:** `SettingsPage.tsx:344-347, 384-388` (HEAD)
- **Evidência:** digitar "15" salva 1 e depois 15; o campo vazio vira NaN e volta ao valor anterior.
- **Correção:** estado local em string, com commit no blur ou no Enter. · **Esforço:** P · **Risco:** Baixo

### FE-06 · Textos fixos em inglês e botões sem rótulo — Baixa · i18n / a11y
- **Local:** `App.tsx:85, 89, 150` ("All", "Proxy … dropped"); `ProxyCard.tsx:102, 127, 230`; `TitleBar.tsx:48, 56, 59-75` (minimizar, maximizar e fechar sem `aria-label`); `MyIPBadge.tsx:47`
- **Impacto:** 15 idiomas mostram inglês em alguns pontos, o que viola a convenção de i18n do `docs/CLAUDE.md`.
- **Correção:** mover para `translations.ts` e rotular os botões. · **Esforço:** P · **Risco:** Baixo

### FE-07 · Links, cores e máscara de IP — Baixa · Bug / design system
- **Local:** `SettingsPage.tsx:267-268, 278-279` (`target="_blank"`: no Wails o caminho é `BrowserOpenURL`, **Hipótese** de que hoje não abre nada); `ProxyCard.tsx:115`, `TitleBar.tsx:73`, `SettingsPage.tsx:254` (cores fixas contra `docs/design-system.md:318`); `MyIPBadge.tsx:21-23` (máscara só para IPv4, IPv6 aparece inteiro)
- **Correção:** `BrowserOpenURL`, tokens de cor, mascarar qualquer endereço que não seja IPv4. · **Esforço:** P · **Risco:** Baixo

### DEAD-01 · Código morto — Baixa · Código morto
- **Local:** bindings nunca importados pelo frontend: `ConnectGroup`, `DisconnectGroup`, `RestartProxy`, `GetConnectionLog`, `Get/SetHealthCheckInterval`, `QuitApp`, `ShowWindow`. Evento `log-entry` sem ouvinte. `ContextMenu.tsx` não é montado (mas `docs/frontend.md:36` o lista).
- **Correção:** remover o que não tem consumidor (o `CheckHealth` é usado internamente; só deixaria de ser exportado). · **Esforço:** P · **Risco:** Baixo

### OPS-03 · Instalador — Baixa · Empacotamento
- **Local:** `installer.iss:54-59` (ru/ja/ko/ar/th aparecem romanizados: "Zapuskat' Proxy0x…", "kidou ji ni…"); `:2` (versão fixa `1.0.0`); sem checagem do WebView2; o desinstalador deixa `%APPDATA%/Proxy0x` (credenciais + chave)
- **Correção:** gravar o `.iss` em UTF-8 com BOM e textos nativos, passar a versão via `/D`, checar o WebView2 e oferecer limpeza de dados na desinstalação. · **Esforço:** P · **Risco:** Baixo

### MCP-03 · Detalhes da integração MCP — Baixa · UX / higiene
- **Local:** `mcp_integration.go:141, 160` (HEAD: sem `CREATE_NO_WINDOW`, `stdin.Close()` antes da resposta); `cmd/proxy0x-mcp/mcp-config.example.json:4` (`D:\\dev\\proxy_ui\\…`)
- **Impacto:** uma janela de console pode piscar e o health check pode dar falso negativo. O exemplo confunde quem for usar.
- **Correção:** esconder a janela, manter o stdin aberto até a resposta e usar um placeholder no exemplo. · **Esforço:** P · **Risco:** Baixo · **Classificação:** Hipótese

---

## Cobertura

**Revisado:** `main.go`, `app.go`, `tray.go`, `internal/proxy/*`, `internal/crypto`, `internal/config`, `mcp_integration.go` (HEAD), `cmd/proxy0x-mcp`, `internal/mcpserver/*`, `pkg/proxycore/*`, `frontend/src/**` (HEAD), `frontend/index.html`, `installer.iss`, `.github/workflows/release.yml`, `wails.json`, e as specs `docs/business-rules.md`, `docs/security.md`, `docs/backend.md`, `docs/mcp.md`, `docs/frontend.md` (trechos), `docs/design-system.md` (trechos).

**NÃO revisado:**
- WIP não commitado (`SettingsPage.tsx`, `translations.ts`, `mcp_integration.go`): só compilado e checado por tipos. Observação: o WIP mira `~/.codex/config.json`, mas o Codex usa `config.toml` (Hipótese).
- Comportamento em runtime no WebView2: nenhum teste manual nesta fase.
- Conteúdo e qualidade das traduções nos 16 idiomas.
- `docs/architecture.md`, `docs/api.md`, `docs/data.md`, `docs/AGENTS.md`, `docs/market-research.md`.
- `EmptyState.tsx`, `StatusBadge.tsx`, `main.tsx`, `index.css` (só passada rápida).
- Internos do `mcp-go` (concorrência de dispatch).
- Data races sob `-race`: precisa de cgo/gcc.
- Assinatura de código dos binários e SmartScreen.

## Artefatos de prova

Os testes de reprodução de CORE-01 e CORE-02 estão em `C:\Users\caosdev\AppData\Local\Temp\claude\D--dev-proxy-ui\f6273a25-1151-47c1-bdc5-a27867c75ce3\scratchpad\repro\server_test.go` (cópia de `internal/proxy/server.go` + upstream falso). Eles viram os testes de regressão da fase 3.
