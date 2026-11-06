import { useState, useEffect, useRef } from 'react'
import { ArrowLeft, RotateCw, Activity, Download, Upload, Globe, ClipboardCopy, Check, Minus, Plus, Github, Bug } from 'lucide-react'
import {
  GetAutoRotateInterval, SetAutoRotateInterval,
  ExportProxies, AddProxiesBatch,
} from 'wailsjs/go/main/App'
import { languages } from '../i18n/translations'
import { useTranslation } from '../i18n/LanguageContext'

interface Props {
  onBack: () => void
  onImported: () => void
}

export default function SettingsPage({ onBack, onImported }: Props) {
  const { lang, t, setLang } = useTranslation()
  const [autoRotate, setAutoRotate] = useState(0)
  const [exportText, setExportText] = useState('')
  const [copyDone, setCopyDone] = useState(false)
  const [importMode, setImportMode] = useState(false)
  const [importText, setImportText] = useState('')
  const [importing, setImporting] = useState(false)


  useEffect(() => {
    GetAutoRotateInterval().then(setAutoRotate).catch(() => {})
  }, [])

  const rotatePresets = [
    { label: '1 min', value: 60 },
    { label: '3 min', value: 180 },
    { label: '5 min', value: 300 },
    { label: '10 min', value: 600 },
  ]

  const setRotate = (v: number) => {
    setAutoRotate(v)
    SetAutoRotateInterval(v).catch(() => {})
  }

  const handleExport = async () => {
    try {
      const text = await ExportProxies()
      setExportText(text)
      await navigator.clipboard.writeText(text)
      setCopyDone(true)
      setTimeout(() => setCopyDone(false), 2000)
    } catch {}
  }

  const handleImport = async () => {
    if (!importText.trim()) return
    setImporting(true)
    try {
      await AddProxiesBatch(importText.trim())
      setImportText('')
      setImportMode(false)
      onImported()
    } catch {}
    setImporting(false)
  }

  return (
    <div className="h-full flex flex-col bg-surface-0">
      {/* Header */}
      <div className="flex items-center gap-2 px-4 pt-4 pb-2">
        <button
          onClick={onBack}
          className="p-1.5 rounded-lg hover:bg-surface-3 transition-colors text-text-secondary hover:text-text-primary"
          aria-label={t.back}
        >
          <ArrowLeft size={16} />
        </button>
        <h1 className="text-sm font-semibold text-text-primary">{t.settings}</h1>
      </div>

      <div className="flex-1 overflow-y-auto px-4 pb-6 space-y-4">
        {/* Language */}
        <Section icon={<Globe size={13} />} title={t.language}>
          <div className="grid grid-cols-2 gap-1">
            {languages.map(l => (
              <button
                key={l.code}
                onClick={() => setLang(l.code)}
                className={`flex items-center gap-2 px-2.5 py-2 rounded-lg text-xs transition-all
                  ${lang === l.code
                    ? 'bg-surface-3 text-text-primary border border-border-strong font-medium'
                    : 'bg-surface-2 text-text-secondary hover:text-text-primary hover:bg-surface-3 border border-transparent'
                  }`}
              >
                <img src={`https://flagcdn.com/w20/${l.flag.toLowerCase()}.png`} alt={l.flag} width={14} height={10} className="inline-block object-contain" loading="lazy" />
                <span className="truncate flex-1">{l.name}</span>
                {lang === l.code && <Check size={12} className="text-accent shrink-0" />}
              </button>
            ))}
          </div>
        </Section>

        {/* Auto Rotation */}
        <Section icon={<RotateCw size={13} />} title={t.autoRotation} desc={t.autoRotationDesc}>
          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <span className="text-xs text-text-secondary">
                {autoRotate > 0 ? t.enabled : t.disabled}
              </span>
              <Toggle active={autoRotate > 0} onToggle={() => setRotate(autoRotate > 0 ? 0 : 300)} />
            </div>
            {autoRotate > 0 && (
              <IntervalPicker presets={rotatePresets} value={autoRotate} onChange={setRotate} />
            )}
          </div>
        </Section>

        {/* Import / Export */}
        <Section icon={<Download size={13} />} title={`${t.exportProxies} / ${t.importProxies}`}>
          <div className="space-y-2">
            <div className="flex gap-2">
              <button
                onClick={handleExport}
                className="flex-1 h-8 flex items-center justify-center gap-1.5 text-xs font-medium
                           bg-surface-2 border border-border-default rounded
                           text-text-secondary hover:text-text-primary hover:bg-surface-3 transition-all"
              >
                <Download size={12} />
                {t.exportProxies}
              </button>
              <button
                onClick={() => setImportMode(!importMode)}
                className="flex-1 h-8 flex items-center justify-center gap-1.5 text-xs font-medium
                           bg-surface-2 border border-border-default rounded
                           text-text-secondary hover:text-text-primary hover:bg-surface-3 transition-all"
              >
                <Upload size={12} />
                {t.importProxies}
              </button>
            </div>
            {exportText && !importMode && (
              <div className="relative">
                <textarea
                  readOnly
                  value={exportText}
                  rows={4}
                  className="w-full px-3 py-2 pr-9 text-[10px] leading-relaxed font-mono bg-surface-2 border border-border-default rounded
                             text-text-secondary resize-none focus:outline-none"
                />
                <button
                  onClick={() => {
                    navigator.clipboard.writeText(exportText)
                    setCopyDone(true)
                    setTimeout(() => setCopyDone(false), 1500)
                  }}
                  className="absolute top-2 right-2 p-1.5 rounded bg-surface-3 hover:bg-surface-4 transition-colors"
                  title={t.copied}
                >
                  {copyDone
                    ? <Check size={11} className="text-accent" />
                    : <ClipboardCopy size={11} className="text-text-tertiary" />
                  }
                </button>
              </div>
            )}
            {importMode && (
              <div className="space-y-2">
                <textarea
                  value={importText}
                  onChange={e => setImportText(e.target.value)}
                  placeholder="host:port:user:pass"
                  rows={4}
                  className="w-full px-3 py-2 text-xs leading-relaxed font-mono bg-surface-2 border border-border-default rounded
                             text-text-primary placeholder:text-text-tertiary resize-none
                             focus:outline-none focus:border-border-strong transition-all"
                />
                <button
                  onClick={handleImport}
                  disabled={importing || !importText.trim()}
                  className="h-7 px-4 text-xs font-medium bg-btn-accent border border-transparent rounded
                             text-btn-accent-text hover:bg-btn-accent-hover font-semibold
                             disabled:opacity-30 disabled:cursor-not-allowed transition-all"
                >
                  {importing ? '...' : t.importAll}
                </button>
              </div>
            )}
          </div>
        </Section>

        {/* About */}
        <Section icon={<Activity size={13} />} title={t.about}>
          <div className="text-center py-2">
            <img src="/logo.png" alt="Proxy0x" className="h-12 mx-auto mb-3 logo-titlebar" draggable={false} />
            <p className="text-[10px] text-text-tertiary mt-1">{t.aboutDesc}</p>
            <div className="flex items-center justify-center gap-2 mt-1">
              <span className="px-2 py-0.5 rounded bg-surface-3 text-[10px] font-mono text-text-tertiary">v1.0.0</span>
            </div>
            <div className="flex gap-1.5 justify-center mt-4">
              <a
                href="https://github.com/henriqqw/proxy0x"
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center gap-1.5 px-2.5 py-1.5 rounded-md text-[10px] font-medium
                           bg-surface-2 border border-border-default text-text-secondary
                           hover:bg-surface-3 hover:text-text-primary transition-all"
              >
                <Github size={12} />
                GitHub
              </a>
              <a
                href="https://github.com/henriqqw/proxy0x/issues/new"
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center gap-1.5 px-2.5 py-1.5 rounded-md text-[10px] font-medium
                           bg-surface-4 border border-border-strong text-text-secondary shadow-[0_2px_6px_rgba(0,0,0,0.3)]
                           hover:bg-surface-5 hover:text-text-primary transition-all"
              >
                <Bug size={12} className="text-red-400" />
                Report Bug
              </a>
            </div>
          </div>
        </Section>
      </div>
    </div>
  )
}

function Section({ icon, title, desc, children }: {
  icon: React.ReactNode
  title: string
  desc?: string
  children: React.ReactNode
}) {
  return (
    <div className="rounded-xl bg-surface-1 border border-border-default p-3 space-y-2.5">
      <div>
        <div className="flex items-center gap-2 text-text-primary">
          <span className="text-text-tertiary">{icon}</span>
          <span className="text-xs font-semibold">{title}</span>
        </div>
        {desc && <p className="text-[10px] text-text-tertiary mt-0.5 ml-[21px]">{desc}</p>}
      </div>
      {children}
    </div>
  )
}

function formatInterval(seconds: number): string {
  if (seconds < 60) return `${seconds}s`
  const m = Math.floor(seconds / 60)
  const s = seconds % 60
  if (m >= 60) {
    const h = Math.floor(m / 60)
    const rm = m % 60
    return rm > 0 ? `${h}h ${rm}min` : `${h}h`
  }
  return s > 0 ? `${m}min ${s}s` : `${m} min`
}

function IntervalPicker({ presets, value, onChange }: {
  presets: { label: string; value: number }[]
  value: number
  onChange: (v: number) => void
}) {
  const [customMode, setCustomMode] = useState(false)
  const [customMin, setCustomMin] = useState(1)
  const isPreset = presets.some(p => p.value === value)
  const inputRef = useRef<HTMLInputElement>(null)

  const openCustom = () => {
    const currentMin = Math.max(1, Math.round(value / 60))
    setCustomMin(currentMin)
    setCustomMode(true)
  }

  const applyCustom = (min: number) => {
    const clamped = Math.max(1, Math.min(120, min))
    setCustomMin(clamped)
    onChange(clamped * 60)
  }

  const step = (delta: number) => {
    const next = Math.max(1, Math.min(120, customMin + delta))
    setCustomMin(next)
    onChange(next * 60)
  }

  return (
    <div className="flex flex-wrap items-center gap-1.5">
      {presets.map(p => (
        <button
          key={p.value}
          onClick={() => { onChange(p.value); setCustomMode(false) }}
          className={`px-2.5 py-1.5 rounded-lg text-[11px] font-medium transition-all
            ${value === p.value
              ? 'bg-surface-5 text-text-primary font-semibold border border-transparent'
              : 'bg-surface-2 text-text-secondary hover:text-text-primary hover:bg-surface-3 border border-transparent'
            }`}
        >
          {p.label}
        </button>
      ))}
      {customMode ? (
        <div className="flex items-center h-7 rounded-lg bg-surface-2 border border-border-strong overflow-hidden">
          <button
            onClick={() => step(-1)}
            className="w-7 h-full flex items-center justify-center text-text-tertiary hover:text-text-primary hover:bg-surface-3 transition-colors"
          >
            <Minus size={10} />
          </button>
          <input
            ref={inputRef}
            autoFocus
            type="text"
            inputMode="numeric"
            value={customMin}
            onChange={e => {
              const v = parseInt(e.target.value)
              if (!isNaN(v)) applyCustom(v)
            }}
            onKeyDown={e => {
              if (e.key === 'ArrowUp') { e.preventDefault(); step(1) }
              if (e.key === 'ArrowDown') { e.preventDefault(); step(-1) }
              if (e.key === 'Escape') setCustomMode(false)
            }}
            className="w-8 h-full text-[11px] font-mono text-text-primary text-center bg-transparent
                       focus:outline-none appearance-none [&::-webkit-inner-spin-button]:appearance-none [&::-webkit-outer-spin-button]:appearance-none"
          />
          <span className="text-[10px] text-text-tertiary pr-2">min</span>
          <button
            onClick={() => step(1)}
            className="w-7 h-full flex items-center justify-center text-text-tertiary hover:text-text-primary hover:bg-surface-3 transition-colors"
          >
            <Plus size={10} />
          </button>
        </div>
      ) : (
        <button
          onClick={openCustom}
          className={`px-2.5 py-1.5 rounded-lg text-[11px] font-medium transition-all
            ${!isPreset
              ? 'bg-accent/10 text-accent border border-accent/20'
              : 'bg-surface-2 text-text-secondary hover:text-text-primary hover:bg-surface-3 border border-transparent'
            }`}
        >
          {!isPreset ? formatInterval(value) : 'Custom'}
        </button>
      )}
    </div>
  )
}

function Toggle({ active, onToggle }: { active: boolean; onToggle: () => void }) {
  return (
    <button
      onClick={onToggle}
      className={`relative w-9 h-5 rounded-full transition-colors ${
        active ? 'bg-[#68E045]' : 'bg-surface-4'
      }`}
      aria-checked={active}
      role="switch"
    >
      <span
        className={`absolute top-0.5 left-0.5 w-4 h-4 rounded-full bg-white shadow transition-transform ${
          active ? 'translate-x-4' : 'translate-x-0'
        }`}
      />
    </button>
  )
}
