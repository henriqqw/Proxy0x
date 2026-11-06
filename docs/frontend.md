# Frontend

React 18 + TypeScript + Tailwind CSS frontend for the Proxy0x desktop app. Bundled with Vite and rendered inside a Wails WebView2 window.

## Stack

| Technology | Version | Role |
|------------|---------|------|
| React | 18 | UI framework |
| TypeScript | 5.6 | Type safety |
| Tailwind CSS | 3 | Utility-first styling via CSS custom properties |
| Vite | 6 | Dev server with HMR, production bundler |
| @dnd-kit | — | Drag-and-drop proxy list reordering |
| Lucide React | — | Monochrome icon set |
| flagcdn.com | — | Country flag images (lazy-loaded) |

## Component Tree

```
App.tsx
├── LanguageProvider              (i18n context)
│   └── AppContent
│       ├── TitleBar              (window controls, theme toggle, language, settings gear)
│       ├── MyIPBadge             (user's real IP with show/hide)
│       ├── ProxyForm             (input field, paste button, batch import)
│       ├── DndContext             (@dnd-kit drag-and-drop wrapper)
│       │   └── SortableContext
│       │       └── ProxyCard[]   (sortable proxy cards with stats, speed test, groups)
│       ├── EmptyState            (illustration when no proxies exist)
│       └── SettingsPage          (language, auto rotation, import/export, MCP, about)
```

Additional components:
- `StatusBadge` — animated dot indicator (offline / connecting / connected / error)
- `FlagIcon` — country flag from flagcdn.com CDN
- `ContextMenu` — right-click language selector

## Page Routing

Simple state-based routing in `AppContent`:
- `page === 'proxies'` — main proxy list view
- `page === 'settings'` — settings page

No React Router. The settings gear icon in `TitleBar` toggles the page. Mouse back/forward buttons (browser navigation) are handled with `popstate` events for in-app navigation.

## Wails Bindings

Auto-generated TypeScript bindings in `frontend/wailsjs/`:
- `wailsjs/go/main/App` — all exported `App` methods (Connect, Disconnect, Rotate, etc.)
- `wailsjs/go/models` — TypeScript classes for Go structs (`ProxyEntry`, `LogEntry`, `MCPStatus`, etc.)
- `wailsjs/runtime/runtime` — Wails runtime functions (`EventsOn`, `EventsEmit`, etc.)

Wails events used:
- `proxy-updated` — single proxy state changed
- `proxies-changed` — proxy list changed (add/remove/batch)
- `proxy-dropped` — health check detected a dropped proxy
- `log-entry` — new connection log entry
- `language-changed` — language selection changed
- `theme-changed` — theme toggled
- `config-changed` — settings changed

## i18n System

### Translations (`src/i18n/translations.ts`)

- `languages` array — 16 supported languages with code, display name, and flag country code
- `Translations` interface — all translatable string keys
- `translations` record — keyed by language code, each entry implements the full interface
- Exported as default for consumption by `LanguageContext`

### LanguageContext (`src/i18n/LanguageContext.tsx`)

React context providing:
- `lang` — current language code
- `t` — current translation object (typed as `Translations`)
- `setLang(code)` — switch language and persist to Go backend

Language is persisted via `SetLanguage()` Wails binding and loaded on mount via `GetLanguage()`.

### Supported Languages

EN, PT-BR, PT-PT, ES, FR, DE, IT, RU, JA, KO, ZH-CN, AR, HI, TR, TH, PL

## Design Tokens

All colors use CSS custom properties defined in `src/index.css`, mapped to Tailwind via `tailwind.config.js`.

### Surfaces (layered depth)

| Token | Dark | Light | Usage |
|-------|------|-------|-------|
| `surface-0` | `#191919` | `#f5f5f5` | App background |
| `surface-1` | `#202020` | `#eeeeee` | Cards, sections |
| `surface-2` | `#252525` | `#e5e5e5` | Inputs, secondary buttons |
| `surface-3` | `#2f2f2f` | `#d9d9d9` | Hover states |
| `surface-4` | `#373737` | `#cccccc` | Elevated buttons |
| `surface-5` | `#434343` | `#b0b0b0` | Active/selected presets |
| `surface-6` | `#505050` | `#999999` | Protocol badges |

### Text

| Token | Dark | Usage |
|-------|------|-------|
| `text-primary` | `#f5f5f5` | Headings, values |
| `text-secondary` | `#a0a0a0` | Labels, descriptions |
| `text-tertiary` | `#6b6b6b` | Hints, timestamps |

### Borders

| Token | Dark | Usage |
|-------|------|-------|
| `border-subtle` | `#ffffff08` | Card inner borders |
| `border-default` | `#ffffff12` | Standard borders |
| `border-strong` | `#ffffff20` | Active/focus borders |

### Accent

- Primary: `#7CDD4F` (green)
- Dark: `#5BBF35`
- Muted: `#7CDD4F40`
- Subtle: `#7CDD4F20`

### Status Colors

- Online: `#7CDD4F`
- Offline: `#6b6b6b`
- Connecting: `#d4a024`
- Error: `#cc4444`

## Typography

- UI font: Inter (Google Fonts fallback: system UI stack)
- Monospace: SF Mono / Cascadia Code / Consolas

## Drag and Drop

Uses `@dnd-kit/core` and `@dnd-kit/sortable`:
- `DndContext` with `closestCenter` collision detection
- `PointerSensor` with 5px activation distance (prevents accidental drags)
- On drag end, `arrayMove` reorders the local array, then `ReorderProxies(ids)` persists to Go
- Each `ProxyCard` is wrapped in a `useSortable` hook

## Features

### Frameless Window
- Custom title bar with drag region (`--wails-draggable`)
- Minimize, maximize toggle, close (to tray) buttons
- Transparent window background (Wails config: `StartHidden: false`, `Frameless: true`)

### Mouse Navigation
- `mouseup` event handler on `App.tsx` for mouse back/forward buttons (button 3/4)
- Navigates between proxy list and settings page

### Loading Skeletons
- Shown while `GetProxies()` is pending on mount
- Animated placeholder cards matching the ProxyCard layout

### Live Uptime Timer
- Connected proxies show a ticking uptime counter
- Updates every second via `setInterval` in `ProxyCard`
- Combines persisted `TotalUptime` with live session duration

### Settings Page Sections
1. **Language** — grid of 16 language buttons with flags
2. **Auto Rotation** — toggle + preset/custom interval picker
3. **Import/Export** — export to clipboard, import from textarea
4. **Claude MCP** — toggle to install/uninstall MCP config with health check
5. **About** — logo, version, GitHub/bug report links
