# Proxy0x

<p align="center"><img src="assets/branding/proxy0x-banner.png" alt="Proxy0x — gerenciador de proxies" width="100%"></p>

<p align="center">
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.26-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go 1.26"></a>
  <a href="https://wails.io/"><img src="https://img.shields.io/badge/Wails-2.16-DF0000?style=flat-square" alt="Wails 2.16"></a>
  <a href="https://react.dev/"><img src="https://img.shields.io/badge/React-18-149ECA?style=flat-square&logo=react&logoColor=white" alt="React 18"></a>
  <a href="https://www.typescriptlang.org/"><img src="https://img.shields.io/badge/TypeScript-5.6-3178C6?style=flat-square&logo=typescript&logoColor=white" alt="TypeScript 5.6"></a>
  <a href="https://github.com/henriqqw/Proxy0x/releases"><img src="https://img.shields.io/badge/Windows%20x64-desktop-4b5563?style=flat-square" alt="Windows x64"></a>
  <a href="https://github.com/henriqqw/Proxy0x/releases"><img src="https://img.shields.io/badge/binary-~10MB-2ea44f?style=flat-square" alt="Single binary ~10MB"></a>
</p>
<p align="center">
  <a href="https://github.com/henriqqw/Proxy0x/releases">⬇ Download</a>
  · <a href="https://github.com/henriqqw/Proxy0x/issues">🐛 Reportar problema</a>
</p>

Proxy manager desktop para Windows. Binário único de ~10MB sem dependências externas — adicione proxies residenciais, rotacione IPs por sessão, monitore latência e geolocalização em tempo real. Interface dark monochrome com system tray, janela frameless e animações fluidas. Construído com **Go**, **Wails** e **React**.

## O que dá para fazer

| Área | Recursos |
|---|---|
| Gerenciamento | Adicionar proxy única, colar múltiplas (batch), remover |
| Conexão | Conectar, desconectar, reiniciar proxy |
| Rotação | Rotacionar IP via session ID — cada rotação gera um novo IP sem trocar o servidor |
| Monitoramento | IP atual, latência em ms, status (connected / offline / connecting / error) |
| Geolocalização | País com bandeira, cidade, região, ISP, timezone, AS, organização |
| Health check | Verificação automática a cada 60s para proxies ativas |
| System tray | Minimizar para bandeja, restaurar com clique, menu de contexto |
| Interface | Frameless, dark monochrome, loading skeletons, animações com Framer Motion |

## Formato de proxy suportado

```
host:port:user:pass
```

Proxies residenciais com targeting geográfico (DataImpulse, etc.) são suportadas. A rotação de IP funciona adicionando `_session-{id}` ao username — cada session ID único recebe um IP diferente do pool.

## Usar

Baixe `ProxyUI.exe` da [release](https://github.com/henriqqw/Proxy0x/releases) e execute. Não precisa de instalação, runtime ou dependência externa.

1. Cole uma proxy no campo de input ou use **Batch import** para várias de uma vez.
2. Clique **Connect** para conectar e ver o IP, localização e latência.
3. Use **Rotate** para obter um novo IP mantendo o mesmo servidor.
4. Use **Restart** para desconectar e reconectar.
5. Feche a janela — o app vai para o system tray. Clique no ícone para restaurar; use **Quit** no menu do tray para sair.

## Arquitetura

```
┌─────────────────────────────────────────┐
│             Wails WebView2              │
│  ┌───────────────────────────────────┐  │
│  │   React + TypeScript + Tailwind   │  │
│  │   Framer Motion · Lucide Icons    │  │
│  └──────────────┬────────────────────┘  │
│                 │ bindings               │
│  ┌──────────────▼────────────────────┐  │
│  │          Go Backend               │  │
│  │  Proxy Manager · Health Checker   │  │
│  │  IP Fetcher · Geo Lookup          │  │
│  └──────────────┬────────────────────┘  │
│                 │                        │
│  ┌──────────────▼────────────────────┐  │
│  │         System Tray               │  │
│  │      energye/systray              │  │
│  └───────────────────────────────────┘  │
└─────────────────────────────────────────┘
          │                    │
          ▼                    ▼
   api.ipify.org         ip-api.com
   (IP lookup)          (geolocation)
```

## Stack

| Camada | Tecnologia | Responsabilidade |
|---|---|---|
| Desktop | Go 1.26, Wails 2.16 | Janela frameless, bindings, eventos, system tray |
| Interface | React 18, TypeScript 5.6, Tailwind CSS 3 | Componentes, tema dark, layout responsivo |
| Animações | Framer Motion 11 | Transições, loading states, layout animations |
| Ícones | Lucide React | Iconografia monochrome consistente |
| Tray | energye/systray | Ícone na bandeja do sistema com menu de contexto |
| Rede | net/http (stdlib) | Proxy HTTP client, health check, IP e geo lookup |

## Desenvolvimento

### Pré-requisitos

- Go 1.22+ (`go.mod`)
- Node.js 18+ e npm
- Wails CLI v2 (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

### Dev mode com hot reload

```bash
cd frontend
npm install
cd ..
wails dev
```

O frontend atualiza instantaneamente via Vite HMR. Mudanças no Go recompilam e reiniciam o app automaticamente.

### Build de produção

```bash
wails build
```

O binário fica em `build/bin/ProxyUI.exe`.

## Estrutura do projeto

```
proxy_ui/
├── main.go              # Entry point, configuração Wails
├── app.go               # Proxy manager, health checker, geo lookup
├── tray.go              # System tray integration
├── wails.json           # Configuração do projeto Wails
├── build/
│   ├── bin/             # Binário de produção
│   └── windows/         # Ícone, manifesto, metadados
└── frontend/
    ├── src/
    │   ├── App.tsx              # Root component
    │   └── components/
    │       ├── TitleBar.tsx      # Frameless window controls
    │       ├── ProxyCard.tsx     # Proxy card com controles e detalhes
    │       ├── ProxyForm.tsx     # Input e batch import
    │       ├── StatusBadge.tsx   # Status indicator animado
    │       ├── MyIPBadge.tsx     # IP real do usuário
    │       ├── FlagIcon.tsx      # Bandeira emoji por country code
    │       └── EmptyState.tsx    # Estado vazio
    ├── wailsjs/                  # Bindings Go ↔ TypeScript (auto-gerados)
    └── tailwind.config.js        # Paleta monochrome
```

## Serviços externos

| Serviço | Uso | Limite |
|---|---|---|
| [api.ipify.org](https://www.ipify.org/) | Descobrir IP atual através da proxy | Sem limite documentado |
| [ip-api.com](http://ip-api.com/) | Geolocalização do IP (país, cidade, ISP, AS) | 45 req/min (free tier) |

## Limitações conhecidas

- Proxies são mantidas em memória — reiniciar o app perde a lista. Persistência em arquivo JSON é o próximo passo.
- Apenas protocolo HTTP proxy (não SOCKS5).
- Testado com proxies DataImpulse; outros provedores que usem o formato `host:port:user:pass` devem funcionar.
- O health check usa `api.ipify.org`, que pode ser bloqueado por algumas redes.
