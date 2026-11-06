import { useState, useEffect } from 'react'
import { Minus, Square, X, Settings, Sun, Moon } from 'lucide-react'
import { MinimizeWindow, CloseWindow, ToggleMaximize, GetTheme, SetTheme } from 'wailsjs/go/main/App'
import { EventsOn } from 'wailsjs/runtime/runtime'
import { useTranslation } from '../i18n/LanguageContext'

interface Props {
  onOpenSettings: () => void
}

export default function TitleBar({ onOpenSettings }: Props) {
  const { t } = useTranslation()
  const [theme, setThemeState] = useState<'dark' | 'light'>('dark')

  useEffect(() => {
    GetTheme().then(saved => {
      const th = saved === 'light' ? 'light' : 'dark'
      setThemeState(th)
      applyTheme(th)
    }).catch(() => {})

    const off = EventsOn('theme-changed', (next: string) => {
      const th = next === 'light' ? 'light' as const : 'dark' as const
      setThemeState(th)
      applyTheme(th)
    })
    return () => { off() }
  }, [])

  const toggleTheme = () => {
    const next = theme === 'dark' ? 'light' : 'dark'
    setThemeState(next)
    applyTheme(next)
    SetTheme(next).catch(() => {})
  }

  return (
    <div
      className="flex items-center justify-between h-10 px-3 border-b border-border-default shrink-0"
      style={{ '--wails-draggable': 'drag', background: 'var(--surface-1)' } as React.CSSProperties}
    >
      <div className="flex items-center">
        <img src="/logo.png" alt="Proxy0x" className="h-4 logo-titlebar" draggable={false} />
      </div>
      <div className="flex items-center gap-0.5 no-drag">
        <button
          onClick={toggleTheme}
          className="p-1.5 rounded hover:bg-surface-4 transition-colors text-text-tertiary hover:text-text-secondary"
          aria-label={t.toggleTheme}
          title={t.toggleTheme}
        >
          {theme === 'dark' ? <Sun size={12} /> : <Moon size={12} />}
        </button>
        <button
          onClick={onOpenSettings}
          className="p-1.5 rounded hover:bg-surface-4 transition-colors text-text-tertiary hover:text-text-secondary"
          aria-label={t.settings}
          title={t.settings}
        >
          <Settings size={12} />
        </button>
        <button
          onClick={() => MinimizeWindow()}
          className="p-1.5 rounded hover:bg-surface-4 transition-colors text-text-tertiary hover:text-text-secondary"
          aria-label={t.minimize}
          title={t.minimize}
        >
          <Minus size={12} />
        </button>
        <button
          onClick={() => ToggleMaximize()}
          className="p-1.5 rounded hover:bg-surface-4 transition-colors text-text-tertiary hover:text-text-secondary"
          aria-label={t.maximize}
          title={t.maximize}
        >
          <Square size={10} />
        </button>
        <button
          onClick={() => CloseWindow()}
          className="p-1.5 rounded hover:bg-surface-4 transition-colors text-text-tertiary hover:text-status-error"
          aria-label={t.closeWindow}
          title={t.closeWindow}
        >
          <X size={12} />
        </button>
      </div>
    </div>
  )
}

function applyTheme(theme: 'dark' | 'light') {
  if (theme === 'light') {
    document.documentElement.setAttribute('data-theme', 'light')
  } else {
    document.documentElement.removeAttribute('data-theme')
  }
}
