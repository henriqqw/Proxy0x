import { useState, useRef, useEffect } from 'react'
import { Minus, Square, X, Globe, Sun, Moon } from 'lucide-react'
import { MinimizeWindow, CloseWindow, ToggleMaximize, GetTheme, SetTheme } from 'wailsjs/go/main/App'
import { EventsOn } from 'wailsjs/runtime/runtime'
import { languages } from '../i18n/translations'
import { useTranslation } from '../i18n/LanguageContext'

export default function TitleBar() {
  const { lang, setLang } = useTranslation()
  const [showLang, setShowLang] = useState(false)
  const [theme, setThemeState] = useState<'dark' | 'light'>('dark')
  const dropdownRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    GetTheme().then(t => {
      const th = t === 'light' ? 'light' : 'dark'
      setThemeState(th)
      applyTheme(th)
    }).catch(() => {})

    const off = EventsOn('theme-changed', (t: string) => {
      const th = t === 'light' ? 'light' as const : 'dark' as const
      setThemeState(th)
      applyTheme(th)
    })
    return () => { off() }
  }, [])

  useEffect(() => {
    if (!showLang) return
    const handle = (e: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(e.target as Node)) {
        setShowLang(false)
      }
    }
    document.addEventListener('mousedown', handle)
    return () => document.removeEventListener('mousedown', handle)
  }, [showLang])

  const toggleTheme = () => {
    const next = theme === 'dark' ? 'light' : 'dark'
    setThemeState(next)
    applyTheme(next)
    SetTheme(next).catch(() => {})
  }

  return (
    <div
      className="flex items-center justify-between h-10 px-3 bg-surface-1/80 border-b border-border-default shrink-0 backdrop-blur-sm"
      style={{ '--wails-draggable': 'drag' } as React.CSSProperties}
    >
      <div className="flex items-center">
        <img src="/logo.png" alt="Proxy0x" className="h-4 logo-titlebar" draggable={false} />
      </div>
      <div className="flex items-center gap-0.5 no-drag">
        <button
          onClick={toggleTheme}
          className="p-1.5 rounded hover:bg-surface-4 transition-colors text-text-tertiary hover:text-text-secondary"
          aria-label="Toggle theme"
        >
          {theme === 'dark' ? <Sun size={12} /> : <Moon size={12} />}
        </button>
        <div className="relative" ref={dropdownRef}>
          <button
            onClick={() => setShowLang(!showLang)}
            className="p-1.5 rounded hover:bg-surface-4 transition-colors text-text-tertiary hover:text-text-secondary"
            aria-label="Change language"
          >
            <Globe size={12} />
          </button>
          {showLang && (
            <div className="absolute right-0 top-full mt-1 w-[180px] bg-surface-2 border border-border-default rounded-lg shadow-xl py-1 z-50 max-h-[320px] overflow-y-auto">
              {languages.map(l => (
                <button
                  key={l.code}
                  onClick={() => { setLang(l.code); setShowLang(false) }}
                  className={`w-full flex items-center gap-2 px-3 py-1.5 text-xs transition-colors
                    ${lang === l.code
                      ? 'text-text-primary bg-surface-4'
                      : 'text-text-secondary hover:text-text-primary hover:bg-surface-3'
                    }`}
                >
                  <img
                    src={`https://flagcdn.com/w20/${l.flag}.png`}
                    alt={l.code}
                    className="w-4 h-3 object-cover rounded-[1px]"
                    loading="lazy"
                  />
                  <span className="flex-1 text-left">{l.name}</span>
                  {lang === l.code && (
                    <span className="text-accent text-[10px]">&#10003;</span>
                  )}
                </button>
              ))}
            </div>
          )}
        </div>
        <button
          onClick={() => MinimizeWindow()}
          className="p-1.5 rounded hover:bg-surface-4 transition-colors text-text-tertiary hover:text-text-secondary"
        >
          <Minus size={12} />
        </button>
        <button
          onClick={() => ToggleMaximize()}
          className="p-1.5 rounded hover:bg-surface-4 transition-colors text-text-tertiary hover:text-text-secondary"
        >
          <Square size={10} />
        </button>
        <button
          onClick={() => CloseWindow()}
          className="p-1.5 rounded hover:bg-surface-4 transition-colors text-text-tertiary hover:text-red-400"
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
