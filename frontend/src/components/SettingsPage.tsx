import { useState, useEffect, useRef } from 'react'
import { ArrowLeft, RotateCw, Activity, Download, Upload, Globe, ClipboardCopy, Check, Minus, Plus, Github, Bug, Bot, Loader2, CheckCircle2, XCircle, AlertTriangle } from 'lucide-react'
import {
  GetAutoRotateInterval, SetAutoRotateInterval,
  ExportProxies, AddProxiesBatch,
  GetMCPStatus, InstallMCPClient, UninstallMCPClient, CheckMCPHealth,
} from 'wailsjs/go/main/App'
import { languages } from '../i18n/translations'
import { useTranslation } from '../i18n/LanguageContext'

const mcpMeta: Record<string, { bg: string; icon: JSX.Element }> = {
  'claude-desktop': {
    bg: 'linear-gradient(135deg, #E87C56, #D4693E)',
    icon: <svg width="12" height="12" viewBox="0 0 24 24" fill="white"><path d="m4.7144 15.9555 4.7174-2.6471.079-.2307-.079-.1275h-.2307l-.7893-.0486-2.6956-.0729-2.3375-.0971-2.2646-.1214-.5707-.1215-.5343-.7042.0546-.3522.4797-.3218.686.0608 1.5179.1032 2.2767.1578 1.6514.0972 2.4468.255h.3886l.0546-.1579-.1336-.0971-.1032-.0972L6.973 9.8356l-2.55-1.6879-1.3356-.9714-.7225-.4918-.3643-.4614-.1578-1.0078.6557-.7225.8803.0607.2246.0607.8925.686 1.9064 1.4754 2.4893 1.8336.3643.3035.1457-.1032.0182-.0728-.164-.2733-1.3539-2.4467-1.445-2.4893-.6435-1.032-.17-.6194c-.0607-.255-.1032-.4674-.1032-.7285L6.287.1335 6.6997 0l.9957.1336.419.3642.6192 1.4147 1.0018 2.2282 1.5543 3.0296.4553.8985.2429.8318.091.255h.1579v-.1457l.1275-1.706.2368-2.0947.2307-2.6957.0789-.7589.3764-.9107.7468-.4918.5828.2793.4797.686-.0668.4433-.2853 1.8517-.5586 2.9021-.3643 1.9429h.2125l.2429-.2429.9835-1.3053 1.6514-2.0643.7286-.8196.85-.9046.5464-.4311h1.0321l.759 1.1293-.34 1.1657-1.0625 1.3478-.8804 1.1414-1.2628 1.7-.7893 1.36.0729.1093.1882-.0183 2.8535-.607 1.5421-.2794 1.8396-.3157.8318.3886.091.3946-.3278.8075-1.967.4857-2.3072.4614-3.4364.8136-.0425.0304.0486.0607 1.5482.1457.6618.0364h1.621l3.0175.2247.7892.522.4736.6376-.079.4857-1.2142.6193-1.6393-.3886-3.825-.9107-1.3113-.3279h-.1822v.1093l1.0929 1.0686 2.0035 1.8092 2.5075 2.3314.1275.5768-.3218.4554-.34-.0486-2.2039-1.6575-.85-.7468-1.9246-1.621h-.1275v.17l.4432.6496 2.3436 3.5214.1214 1.0807-.17.3521-.6071.2125-.6679-.1214-1.3721-1.9246L14.38 17.959l-1.1414-1.9428-.1397.079-.674 7.2552-.3156.3703-.7286.2793-.6071-.4614-.3218-.7468.3218-1.4753.3886-1.9246.3157-1.53.2853-1.9004.17-.6314-.0121-.0425-.1397.0182-1.4328 1.9672-2.1796 2.9446-1.7243 1.8456-.4128.164-.7164-.3704.0667-.6618.4008-.5889 2.386-3.0357 1.4389-1.882.929-1.0868-.0062-.1579h-.0546l-6.3385 4.1164-1.1293.1457-.4857-.4554.0608-.7467.2307-.2429 1.9064-1.3114Z"/></svg>,
  },
  codex: {
    bg: 'linear-gradient(135deg, #171717, #2a2a2a)',
    icon: <svg width="12" height="12" viewBox="0 0 24 24" fill="white"><path d="M22.282 9.821a5.985 5.985 0 0 0-.516-4.91 6.046 6.046 0 0 0-6.51-2.9A6.065 6.065 0 0 0 4.981 4.18a5.985 5.985 0 0 0-3.998 2.9 6.046 6.046 0 0 0 .743 7.097 5.98 5.98 0 0 0 .51 4.911 6.051 6.051 0 0 0 6.515 2.9A5.985 5.985 0 0 0 13.26 24a6.056 6.056 0 0 0 5.772-4.206 5.99 5.99 0 0 0 3.997-2.9 6.056 6.056 0 0 0-.747-7.073zM13.26 22.43a4.476 4.476 0 0 1-2.876-1.04l.14-.08 4.778-2.759a.795.795 0 0 0 .393-.681v-6.737l2.02 1.168a.071.071 0 0 1 .038.052v5.583a4.504 4.504 0 0 1-4.494 4.494zM3.6 18.304a4.47 4.47 0 0 1-.535-3.014l.142.085 4.783 2.759a.771.771 0 0 0 .78 0l5.843-3.369v2.332a.08.08 0 0 1-.033.062l-4.83 2.786a4.499 4.499 0 0 1-6.15-1.64zM2.34 7.896a4.485 4.485 0 0 1 2.366-1.973V11.6a.766.766 0 0 0 .388.677l5.814 3.354-2.02 1.168a.076.076 0 0 1-.071 0l-4.83-2.786A4.504 4.504 0 0 1 2.34 7.872zm16.597 3.855L13.104 8.364l2.02-1.164a.076.076 0 0 1 .071 0l4.83 2.791a4.494 4.494 0 0 1-.676 8.105v-5.678a.79.79 0 0 0-.412-.676zm2.01-3.023-.141-.085-4.774-2.782a.776.776 0 0 0-.785 0L9.409 9.23V6.897a.066.066 0 0 1 .028-.061l4.83-2.787a4.5 4.5 0 0 1 6.68 4.66zM8.307 12.863l-2.02-1.164a.08.08 0 0 1-.038-.057V6.074a4.5 4.5 0 0 1 7.375-3.453l-.14.08-4.779 2.758a.795.795 0 0 0-.393.681zm1.097-2.361 2.603-1.506 2.603 1.506v3.012l-2.603 1.506-2.603-1.506z"/></svg>,
  },
  cursor: {
    bg: 'linear-gradient(135deg, #2a2a2a, #1a1a1a)',
    icon: <svg width="12" height="12" viewBox="0 0 24 24" fill="white"><path d="M11.503.131 1.891 5.678a.84.84 0 0 0-.42.726v11.188c0 .3.162.575.42.724l9.609 5.55a1 1 0 0 0 .998 0l9.61-5.55a.84.84 0 0 0 .42-.724V6.404a.84.84 0 0 0-.42-.726L12.497.131a1.01 1.01 0 0 0-.996 0M2.657 6.338h18.55c.263 0 .43.287.297.515L12.23 22.918c-.062.107-.229.064-.229-.06V12.335a.59.59 0 0 0-.295-.51l-9.11-5.257c-.109-.063-.064-.23.061-.23"/></svg>,
  },
  windsurf: {
    bg: 'linear-gradient(135deg, #09B6A2, #07A090)',
    icon: <svg width="12" height="12" viewBox="0 0 24 24" fill="white"><path d="M23.55 5.067c-1.2038-.002-2.1806.973-2.1806 2.1765v4.8676c0 .972-.8035 1.7594-1.7597 1.7594-.568 0-1.1352-.286-1.4718-.7659l-4.9713-7.1003c-.4125-.5896-1.0837-.941-1.8103-.941-1.1334 0-2.1533.9635-2.1533 2.153v4.8957c0 .972-.7969 1.7594-1.7596 1.7594-.57 0-1.1363-.286-1.4728-.7658L.4076 5.1598C.2822 4.9798 0 5.0688 0 5.2882v4.2452c0 .2147.0656.4228.1884.599l5.4748 7.8183c.3234.462.8006.8052 1.3509.9298 1.3771.313 2.6446-.747 2.6446-2.0977v-4.893c0-.972.7875-1.7593 1.7596-1.7593h.003a1.798 1.798 0 0 1 1.4718.7658l4.9723 7.0994c.4135.5905 1.05.941 1.8093.941 1.1587 0 2.1515-.9645 2.1515-2.153v-4.8948c0-.972.7875-1.7594 1.7596-1.7594h.194a.22.22 0 0 0 .2204-.2202v-4.622a.22.22 0 0 0-.2203-.2203Z"/></svg>,
  },
}

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
  const [mcpClients, setMcpClients] = useState<{ id: string; name: string; detected: boolean; installed: boolean }[]>([])
  const [mcpExeExists, setMcpExeExists] = useState(true)
  const [mcpLoading, setMcpLoading] = useState<string | null>(null)
  const [mcpHealth, setMcpHealth] = useState<'idle' | 'checking' | 'ok' | 'fail'>('idle')
  const [mcpError, setMcpError] = useState('')
  const [mcpDone, setMcpDone] = useState(false)

  useEffect(() => {
    GetAutoRotateInterval().then(setAutoRotate).catch(() => {})
    GetMCPStatus().then(s => {
      setMcpExeExists(s.exeExists)
      setMcpClients((s.clients || []).map((c: any) => ({ id: c.id, name: c.name, detected: c.detected, installed: c.installed })))
    }).catch(() => {})
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

        {/* MCP Integration */}
        <Section icon={<Bot size={13} />} title={t.mcpIntegration} desc={t.mcpIntegrationDesc}>
          <div className="grid grid-cols-2 gap-1.5">
            {mcpClients.map(client => {
              const meta = mcpMeta[client.id]
              const disabled = !client.detected || !mcpExeExists || mcpLoading !== null
              return (
                <div key={client.id} className={`flex items-center justify-between px-2.5 py-2 rounded-[10px] transition-all ${
                  client.installed
                    ? 'bg-accent/[0.06] border border-accent/[0.15]'
                    : 'bg-surface-2 border border-transparent'
                } ${!client.detected ? 'opacity-40' : ''}`}>
                  <div className="flex items-center gap-1.5 min-w-0">
                    <div className="w-[18px] h-[18px] rounded-[5px] shrink-0 flex items-center justify-center"
                         style={{ background: meta?.bg || '#666' }}>
                      {meta?.icon}
                    </div>
                    <div className="min-w-0">
                      <div className="text-[10px] font-semibold text-text-primary leading-tight">{client.name}</div>
                      <div className={`text-[8px] leading-tight truncate ${
                        !client.detected ? 'text-text-tertiary'
                        : client.installed ? 'text-accent'
                        : 'text-text-tertiary'
                      }`}>
                        {!client.detected ? t.mcpNotInstalled
                          : client.installed ? t.mcpConnected
                          : t.mcpNotConnected}
                      </div>
                    </div>
                  </div>
                  <Toggle
                    active={client.installed}
                    onToggle={async () => {
                      if (disabled) return
                      setMcpLoading(client.id)
                      setMcpHealth('idle')
                      setMcpError('')
                      try {
                        if (client.installed) {
                          await UninstallMCPClient(client.id)
                          setMcpClients(prev => prev.map(c => c.id === client.id ? { ...c, installed: false } : c))
                        } else {
                          await InstallMCPClient(client.id)
                          setMcpClients(prev => prev.map(c => c.id === client.id ? { ...c, installed: true } : c))
                          setMcpDone(true)
                          setTimeout(() => setMcpDone(false), 3000)
                          setMcpHealth('checking')
                          const result = await CheckMCPHealth()
                          setMcpHealth(result.healthOk ? 'ok' : 'fail')
                          if (!result.healthOk && result.healthError) setMcpError(result.healthError)
                        }
                      } catch (e: any) {
                        setMcpHealth('fail')
                        setMcpError(e?.message || t.mcpHealthFail)
                      }
                      setMcpLoading(null)
                    }}
                  />
                </div>
              )
            })}
          </div>
          {(mcpHealth !== 'idle' || mcpDone) && (
            <div className="mt-2 space-y-1 text-center">
              {mcpHealth === 'checking' && (
                <div className="flex items-center justify-center gap-1.5">
                  <Loader2 size={11} className="text-text-tertiary animate-spin" />
                  <span className="text-[10px] text-text-tertiary">{t.mcpChecking}</span>
                </div>
              )}
              {mcpHealth === 'ok' && (
                <div className="flex items-center justify-center gap-1.5">
                  <CheckCircle2 size={11} className="text-[#68E045]" />
                  <span className="text-[10px] text-[#68E045]">{t.mcpHealthOk}</span>
                </div>
              )}
              {mcpHealth === 'fail' && (
                <div className="flex items-center justify-center gap-1.5">
                  <XCircle size={11} className="text-red-400" />
                  <span className="text-[10px] text-red-400">{mcpError || t.mcpHealthFail}</span>
                </div>
              )}
              {mcpDone && (
                <p className="text-[10px] text-amber-400">{t.mcpRestartClaude}</p>
              )}
            </div>
          )}
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
