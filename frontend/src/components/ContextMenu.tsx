import { useState, useEffect, useCallback } from 'react'
import { Globe } from 'lucide-react'
import { languages } from '../i18n/translations'
import { useTranslation } from '../i18n/LanguageContext'

export default function ContextMenu() {
  const { lang, t, setLang } = useTranslation()
  const [pos, setPos] = useState<{ x: number; y: number } | null>(null)

  const handleContext = useCallback((e: MouseEvent) => {
    e.preventDefault()
    const x = Math.min(e.clientX, window.innerWidth - 200)
    const y = Math.min(e.clientY, window.innerHeight - 400)
    setPos({ x, y })
  }, [])

  const handleClose = useCallback(() => setPos(null), [])

  useEffect(() => {
    document.addEventListener('contextmenu', handleContext)
    document.addEventListener('click', handleClose)
    return () => {
      document.removeEventListener('contextmenu', handleContext)
      document.removeEventListener('click', handleClose)
    }
  }, [handleContext, handleClose])

  if (!pos) return null

  return (
    <div
      className="fixed inset-0 z-[9999]"
      onClick={handleClose}
      onContextMenu={e => { e.preventDefault(); handleClose() }}
    >
      <div
        className="absolute border border-border-default rounded-lg shadow-xl py-1 w-[180px] overflow-hidden"
        style={{ left: pos.x, top: pos.y, background: 'var(--surface-2)' }}
        onClick={e => e.stopPropagation()}
      >
        <div className="flex items-center gap-2 px-3 py-1.5 text-[10px] text-text-tertiary uppercase tracking-wider">
          <Globe size={10} />
          {t.language}
        </div>
        <div className="border-t border-border-subtle my-0.5" />
        <div className="max-h-[320px] overflow-y-auto">
          {languages.map(l => (
            <button
              key={l.code}
              onClick={() => { setLang(l.code); handleClose() }}
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
      </div>
    </div>
  )
}
