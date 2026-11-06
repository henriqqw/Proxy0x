# Proxy0x Design System

## Identity

**App:** Proxy0x — lightweight desktop proxy manager  
**Stack:** React 18 + TypeScript + Tailwind CSS + Wails (Go)  
**Theme:** Dark-first neutral with green accent  
**Font:** Inter (UI), SF Mono / Cascadia Code (monospace)

---

## Color Tokens

All colors are CSS custom properties defined in `index.css`, consumed via Tailwind in `tailwind.config.js`.

### Surfaces

Layered neutral scale. Higher number = more elevated.

| Token        | Dark           | Light          | Usage                              |
|-------------|----------------|----------------|-------------------------------------|
| `surface-0` | `#191919`      | `#f5f5f5`      | App background, root                |
| `surface-1` | `#202020`      | `#eeeeee`      | Cards, sections, title bar          |
| `surface-2` | `#252525`      | `#e5e5e5`      | Inputs, secondary buttons, textareas|
| `surface-3` | `#2f2f2f`      | `#d9d9d9`      | Hover states, selected items        |
| `surface-4` | `#373737`      | `#cccccc`      | Elevated buttons, icon button hover |
| `surface-5` | `#434343`      | `#b0b0b0`      | Active/selected presets             |
| `surface-6` | `#505050`      | `#999999`      | Protocol badge (HTTP)               |

### Borders

| Token            | Dark                      | Light                     | Usage                    |
|-----------------|---------------------------|---------------------------|--------------------------|
| `border-subtle`  | `rgba(255,255,255,0.06)` | `rgba(0,0,0,0.08)`       | Dividers, separators     |
| `border-default` | `rgba(255,255,255,0.09)` | `rgba(0,0,0,0.13)`       | Card borders, inputs     |
| `border-strong`  | `rgba(255,255,255,0.14)` | `rgba(0,0,0,0.20)`       | Focused inputs, selected |

### Text

| Token            | Dark       | Light      | Usage                        |
|-----------------|------------|------------|-------------------------------|
| `text-primary`   | `#ebebeb` | `#1a1a1a`  | Headings, IP addresses, names|
| `text-secondary`  | `#999999` | `#444444`  | Labels, button text, metadata|
| `text-tertiary`   | `#6b6b6b` | `#666666`  | Hints, descriptions, disabled|

### Accent

Fixed across themes — the brand green.

| Token            | Value          | Usage                                 |
|-----------------|----------------|----------------------------------------|
| `accent`         | `#7CDD4F`     | Status online, checkmarks, logo accent |
| `accent-dark`    | `#5BBF35`     | Hover/pressed accent                   |
| `accent-muted`   | `#7CDD4F40`  | Accent backgrounds (40% opacity)       |
| `accent-subtle`  | `#7CDD4F20`  | Accent tint (20% opacity)              |

### Button Accent (theme-aware)

Inverted color pair for primary CTA buttons.

| Token              | Dark       | Light      | Usage              |
|-------------------|------------|------------|--------------------|
| `btn-accent`       | `#ebebeb` | `#1a1a1a`  | Connect button bg  |
| `btn-accent-text`  | `#191919` | `#f5f5f5`  | Connect button text|
| `btn-accent-hover` | `#d4d4d4` | `#333333`  | Connect hover bg   |

### Status

| Token               | Value      | Usage            |
|---------------------|------------|------------------|
| `status-online`      | `#7CDD4F` | Connected state  |
| `status-offline`     | `#6b6b6b` | Offline state    |
| `status-connecting`  | `#d4a024` | Connecting state |
| `status-error`       | `#cc4444` | Error state      |

### Scrollbar

| Token              | Dark       | Light      |
|-------------------|------------|------------|
| `scrollbar-thumb`  | `#373737` | `#cccccc`  |
| `scrollbar-hover`  | `#434343` | `#b0b0b0`  |

---

## Typography

### Font Stack

```
UI:   Inter, -apple-system, BlinkMacSystemFont, SF Pro Display, Segoe UI, sans-serif
Mono: SF Mono, Cascadia Code, Consolas, monospace
```

### Scale

| Element              | Size        | Weight     | Extra                        |
|---------------------|-------------|------------|-------------------------------|
| App title (About)    | `text-lg`   | `800`      | `tracking-tight`             |
| Section title        | `text-xs`   | `600`      | —                             |
| Page header          | `text-sm`   | `600`      | —                             |
| Button label         | `10px`      | `500-600`  | —                             |
| Input text           | `text-xs`   | `400`      | `font-mono`                  |
| Badge (S5/HTTP)      | `8px`       | `700`      | `uppercase tracking-wider`   |
| Status label         | `10px`      | `500`      | `uppercase tracking-wider`   |
| Metadata / hints     | `10px`      | `400`      | `text-text-tertiary`         |
| Version badge        | `10px`      | `400`      | `font-mono`                  |

---

## Spacing

Based on Tailwind's 4px grid. The app is compact — desktop-sized, not mobile-first.

| Context              | Value       | Tailwind        |
|---------------------|-------------|-----------------|
| Page padding         | `16px`      | `px-4`          |
| Section internal     | `12px`      | `p-3`           |
| Between sections     | `16px`      | `space-y-4`     |
| Card padding         | `12px`      | `px-3 pt-3 pb-2`|
| Gap between items    | `6-8px`     | `gap-1.5-2`     |
| Button height        | `28-32px`   | `h-7` / `h-8`   |
| Icon button padding  | `6px`       | `p-1.5`         |

---

## Border Radius

| Context              | Value   | Tailwind     |
|---------------------|---------|--------------|
| Default              | `8px`   | `rounded`    |
| Cards / sections     | `12px`  | `rounded-xl` |
| App window           | `10px`  | `rounded-[10px]` |
| Badges / tags        | `4px`   | `rounded`    |
| Buttons              | `6px`   | `rounded-md` |
| Status dot           | `50%`   | `rounded-full` |

---

## Components

### Protocol Toggle (ProxyForm)

Segmented toggle, side by side, no gap, no border on wrapper.

| State         | Background          | Text              |
|--------------|---------------------|-------------------|
| S5 active     | `bg-btn-accent`     | `text-btn-accent-text` |
| HTTP active   | `bg-[#505050]`      | `text-[#d4d4d4]`  |
| Inactive      | `bg-[#505050]/25`   | `text-text-tertiary` |

Size: `h-8 px-2.5 text-[9px] font-bold uppercase tracking-wider`  
Default protocol: `socks5`

### Protocol Badge (ProxyCard)

Same color rules as toggle. Compact inline badge.

```
px-1.5 py-0.5 text-[8px] font-bold uppercase tracking-wider rounded
```

### Status Badge

Dot + uppercase label.

- Dot: `w-1.5 h-1.5 rounded-full` with status color
- Connecting: `animate-pulse`
- Label: `text-[10px] font-medium uppercase tracking-wider`

### Control Buttons (ProxyCard)

Two variants:

**Primary (Connect):**
```
h-7 px-2.5 bg-btn-accent border-transparent text-btn-accent-text
hover:bg-btn-accent-hover font-semibold rounded
```

**Secondary (Rotate, Speed Test):**
```
h-7 px-2.5 bg-transparent border-border-strong text-text-secondary
hover:text-text-primary hover:bg-surface-2 rounded
```

### Section (SettingsPage)

Container for settings groups.

```
rounded-xl bg-surface-1 border border-border-default p-3 space-y-2.5
```

Header: icon (`text-text-tertiary`) + title (`text-xs font-semibold`)

### Toggle Switch (Auto Rotation)

```
w-9 h-5 rounded-full
Active:   bg-[#68E045]  (exception: uses old green, not #7CDD4F)
Inactive: bg-surface-4
Knob:     w-4 h-4 rounded-full bg-white shadow
```

### Language Selector

Grid `grid-cols-2 gap-1`. Checkmark style.

| State      | Style                                                    |
|-----------|----------------------------------------------------------|
| Selected   | `bg-surface-3 text-text-primary border-border-strong font-medium` + Check icon in accent |
| Unselected | `bg-surface-2 text-text-secondary border-transparent`    |

### Interval Presets

Flat pills with Surface-5 active state.

| State      | Style                                                    |
|-----------|----------------------------------------------------------|
| Selected   | `bg-surface-5 text-text-primary font-semibold`           |
| Unselected | `bg-surface-2 text-text-secondary hover:bg-surface-3`    |

Custom mode: input with `-`/`+` buttons, `bg-surface-2 border-border-strong`.

### Title Bar

```
h-10 px-3 bg-surface-1 border-b border-border-default
```

Window control buttons: `p-1.5 rounded hover:bg-surface-4 text-text-tertiary`  
Close hover: `hover:text-red-400`

### Empty State

Centered vertically, mascot icon at `opacity-0.3` with `animate-breathe`.  
Text: `text-xs text-text-tertiary` + hint at `text-[10px] text-text-tertiary/70`.

### Input Fields

```
h-8 px-3 text-xs font-mono
bg-surface-2 border border-border-default rounded-lg
text-text-primary placeholder:text-text-tertiary
focus:border-border-strong focus:bg-surface-3
```

### Textarea

```
px-3 py-2 text-xs leading-relaxed font-mono
bg-surface-2 border border-border-default rounded
text-text-primary resize-none
focus:border-border-strong
```

With floating copy button: `absolute top-2 right-2 p-1.5 rounded bg-surface-3 hover:bg-surface-4`

### About Section (Hero)

Centered layout: logo `h-12`, description, version badge (`bg-surface-3 font-mono`), action buttons.

**GitHub button:**
```
px-2.5 py-1.5 rounded-md text-[10px] font-medium
bg-surface-2 border border-border-default text-text-secondary
hover:bg-surface-3 hover:text-text-primary
```

**Report Bug button (Surface elevated):**
```
px-2.5 py-1.5 rounded-md text-[10px] font-medium
bg-surface-4 border border-border-strong text-text-secondary
shadow-[0_2px_6px_rgba(0,0,0,0.3)]
hover:bg-surface-5 hover:text-text-primary
```
Bug icon: `text-red-400`

---

## Animations

| Name            | Duration | Easing     | Usage                          |
|----------------|----------|------------|--------------------------------|
| `fade-in`       | `0.2s`   | `ease-out` | Elements appearing             |
| `slide-down`    | `0.2s`   | `ease-out` | Form entrance, dropdowns       |
| `expand`        | `0.2s`   | `ease-out` | Collapsible sections           |
| `breathe`       | `3s`     | `ease-in-out infinite` | Empty state mascot  |
| `pulse-slow`    | `3s`     | `cubic-bezier` | Connecting indicator    |

All transitions on interactive elements: `transition-all` or `transition-colors`.

---

## Icons

**Library:** Lucide React  
**Default size:** `12-13px` for inline, `16px` for standalone  
**Color:** `text-text-tertiary` default, `text-text-secondary` on hover

Common icons used:
- `Power`, `X`, `RotateCw`, `Gauge` — proxy controls
- `Globe`, `MapPin`, `Building2`, `Activity` — metadata
- `ArrowLeft`, `Settings`, `Sun`, `Moon` — navigation
- `Github`, `Bug` — about section
- `Check`, `Minus`, `Plus` — form controls
- `Download`, `Upload`, `ClipboardCopy` — import/export
- `GripVertical` — drag handle

---

## Theming

Dark is the default. Light is toggled via `data-theme="light"` on `<html>`.

### Rules

1. All colors go through CSS custom properties — never hardcode theme-specific values
2. Exception: accent green `#7CDD4F` is fixed across themes
3. Exception: toggle switch uses `#68E045` (old green, intentional)
4. In light mode, logo images get `filter: brightness(0)` via `.logo-titlebar` / `.logo-footer` classes
5. `btn-accent` pair inverts (white-on-dark → dark-on-light) for contrast

### Scrollbar

Custom webkit scrollbar: `4px` wide, transparent track, `surface-4/5` thumb.

### Window

Root `#root` has `border-radius: 10px`, `overflow: hidden`, and `border: 1px solid border-default` for the Wails frameless window look.

---

## Layout Patterns

### Main View
```
flex flex-col h-full
├── TitleBar (h-10, shrink-0)
├── Content (flex-1, overflow-y-auto, px-4)
│   ├── MyIPBadge
│   ├── ProxyForm
│   └── ProxyCard[] (drag-and-drop via @dnd-kit)
└── Footer logo (opacity-0.08)
```

### Settings View
```
flex flex-col h-full bg-surface-0
├── Header (back button + title)
└── Scrollable content (px-4 pb-6 space-y-4)
    ├── Section: Language
    ├── Section: Auto Rotation
    ├── Section: Import/Export
    └── Section: About
```

---

## Do / Don't

| Do                                          | Don't                                       |
|--------------------------------------------|----------------------------------------------|
| Use surface scale for elevation             | Use arbitrary gray values                    |
| Keep text at 10-12px for UI density         | Go below 8px or above 14px for body text     |
| Use `transition-all` on interactive         | Skip transitions on state changes            |
| Use `border-border-default` on containers  | Use visible borders as decoration            |
| Use `font-mono` for IPs, ports, versions   | Use mono for labels or buttons               |
| Use accent green for positive/active states | Use accent for warnings or errors            |
| Use `rounded` (8px) for buttons             | Mix `rounded-lg` and `rounded` on same level |
| Keep buttons at `h-7` or `h-8`             | Make buttons taller than 32px                |
