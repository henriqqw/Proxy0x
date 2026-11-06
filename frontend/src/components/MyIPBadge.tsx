import { useState } from 'react'
import { Globe, RefreshCw, Eye, EyeOff } from 'lucide-react'
import { useTranslation } from '../i18n/LanguageContext'

interface Props {
  ip: string
  onRefresh: () => void
}

export default function MyIPBadge({ ip, onRefresh }: Props) {
  const { t } = useTranslation()
  const [spinning, setSpinning] = useState(false)
  const [visible, setVisible] = useState(false)

  const handleRefresh = () => {
    setSpinning(true)
    onRefresh()
    setTimeout(() => setSpinning(false), 1000)
  }

  // IPv4 keeps its first octet; anything else (IPv6) is fully masked.
  const masked = !ip
    ? '...'
    : /^\d+\.\d+\.\d+\.\d+$/.test(ip)
      ? ip.replace(/^(\d+)\.(.+)$/, (_, first, rest) => `${first}.${rest.replace(/\d/g, '●')}`)
      : ip.replace(/[0-9a-fA-F]/g, '●')

  return (
    <div className="flex items-center justify-between py-1.5 px-3 rounded-lg bg-surface-1 border border-border-default animate-slide-down">
      <div className="flex items-center gap-2">
        <Globe size={12} className="text-accent/50" />
        <span className="text-[10px] text-text-tertiary uppercase tracking-wider">{t.yourIp}</span>
      </div>
      <div className="flex items-center gap-1.5">
        <span className="text-xs font-mono text-text-secondary select-none">
          {visible ? (ip || '...') : masked}
        </span>
        <button
          onClick={() => setVisible(!visible)}
          className="p-1 rounded hover:bg-surface-4 transition-colors text-text-tertiary hover:text-text-secondary"
          title={visible ? t.hideIp : t.showIp}
          aria-label={visible ? t.hideIp : t.showIp}
        >
          {visible ? <Eye size={11} /> : <EyeOff size={11} />}
        </button>
        <button
          onClick={handleRefresh}
          className="p-1 rounded hover:bg-surface-4 transition-colors text-text-tertiary hover:text-text-secondary"
          aria-label={t.refreshIp}
          title={t.refreshIp}
        >
          <RefreshCw size={10} className={spinning ? 'animate-spin' : ''} />
        </button>
      </div>
    </div>
  )
}
