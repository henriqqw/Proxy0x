import { useState } from 'react'
import { motion } from 'framer-motion'
import {
  Power, PowerOff, RotateCw, RefreshCw,
  Trash2, Copy, MapPin, Clock, Wifi, Server,
  ChevronDown, Activity, Building2
} from 'lucide-react'
import {
  ConnectProxy, DisconnectProxy, RotateProxy,
  RestartProxy, RemoveProxy
} from 'wailsjs/go/main/App'
import { main } from 'wailsjs/go/models'
import StatusBadge from './StatusBadge'
import FlagIcon from './FlagIcon'

interface Props {
  proxy: main.ProxyEntry
  onRemoved: () => void
}

export default function ProxyCard({ proxy, onRemoved }: Props) {
  const [busy, setBusy] = useState<string | null>(null)
  const [expanded, setExpanded] = useState(false)
  const [copied, setCopied] = useState(false)

  type Status = 'offline' | 'connecting' | 'connected' | 'error'
  const status = (proxy.status || 'offline') as Status
  const isActive = status === 'connected'
  const isConnecting = status === 'connecting'

  const action = async (name: string, fn: () => Promise<any>) => {
    setBusy(name)
    try { await fn() } catch {}
    setBusy(null)
  }

  const copyIP = () => {
    if (proxy.currentIp) {
      navigator.clipboard.writeText(proxy.currentIp)
      setCopied(true)
      setTimeout(() => setCopied(false), 1500)
    }
  }

  return (
    <motion.div
      layout
      initial={{ opacity: 0, y: 8 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, scale: 0.95 }}
      transition={{ duration: 0.2 }}
      className={`rounded-lg border transition-all ${
        isActive
          ? 'bg-surface-2 border-border-strong shadow-[0_1px_8px_rgba(0,0,0,0.4)]'
          : 'bg-surface-1 border-border-default hover:border-border-strong'
      }`}
    >
      {/* Header */}
      <div className="px-3 pt-3 pb-2">
        <div className="flex items-center justify-between mb-2">
          <div className="flex items-center gap-2 min-w-0 flex-1">
            <StatusBadge status={status} />
            <span className="text-[10px] font-mono text-text-tertiary truncate">
              {proxy.host}:{proxy.port}
            </span>
          </div>
          <button
            onClick={() => { RemoveProxy(proxy.id); onRemoved() }}
            className="p-1 rounded text-text-tertiary/40 hover:text-status-error hover:bg-surface-4 transition-all"
          >
            <Trash2 size={11} />
          </button>
        </div>

        {/* IP Display */}
        {proxy.currentIp && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className="flex items-center gap-2 mb-2"
          >
            {proxy.location?.countryCode && (
              <FlagIcon countryCode={proxy.location.countryCode} size={13} />
            )}
            <button
              onClick={copyIP}
              className="flex items-center gap-1.5 group"
            >
              <span className="text-sm font-mono text-text-primary font-medium">
                {proxy.currentIp}
              </span>
              <Copy
                size={10}
                className={`transition-colors ${
                  copied ? 'text-status-online' : 'text-text-tertiary opacity-0 group-hover:opacity-100'
                }`}
              />
            </button>
            {proxy.latency > 0 && (
              <span className="ml-auto flex items-center gap-1 text-[10px] text-text-tertiary">
                <Activity size={9} />
                {proxy.latency}ms
              </span>
            )}
          </motion.div>
        )}

        {/* Location */}
        {proxy.location && (
          <div className="flex flex-wrap items-center gap-x-3 gap-y-1 text-[10px] text-text-tertiary mb-2">
            <span className="flex items-center gap-1">
              <MapPin size={9} />
              {[proxy.location.city, proxy.location.region, proxy.location.country]
                .filter(Boolean).join(', ')}
            </span>
            {proxy.location.isp && (
              <span className="flex items-center gap-1">
                <Building2 size={9} />
                {proxy.location.isp}
              </span>
            )}
          </div>
        )}

        {/* Controls */}
        <div className="flex items-center gap-1">
          {!isActive && !isConnecting ? (
            <ControlBtn
              icon={<Power size={12} />}
              label="Connect"
              onClick={() => action('connect', () => ConnectProxy(proxy.id))}
              busy={busy === 'connect'}
              accent
            />
          ) : (
            <ControlBtn
              icon={<PowerOff size={12} />}
              label="Stop"
              onClick={() => action('disconnect', () => DisconnectProxy(proxy.id))}
              busy={busy === 'disconnect'}
            />
          )}
          <ControlBtn
            icon={<RotateCw size={12} />}
            label="Rotate"
            onClick={() => action('rotate', () => RotateProxy(proxy.id))}
            busy={busy === 'rotate'}
            disabled={!isActive && !isConnecting}
          />
          <ControlBtn
            icon={<RefreshCw size={12} />}
            label="Restart"
            onClick={() => action('restart', () => RestartProxy(proxy.id))}
            busy={busy === 'restart'}
            disabled={!isActive}
          />
        </div>
      </div>

      {/* Expandable details */}
      <button
        onClick={() => setExpanded(!expanded)}
        className="w-full flex items-center justify-center py-1.5 border-t border-border-subtle
                   text-text-tertiary/40 hover:text-text-tertiary hover:bg-surface-3/60 transition-all"
      >
        <ChevronDown size={10} className={`transition-transform ${expanded ? 'rotate-180' : ''}`} />
      </button>

      {expanded && (
        <motion.div
          initial={{ height: 0, opacity: 0 }}
          animate={{ height: 'auto', opacity: 1 }}
          className="px-3 pb-3 space-y-1.5 text-[10px] text-text-tertiary border-t border-border-subtle"
        >
          <DetailRow icon={<Server size={9} />} label="Host" value={`${proxy.host}:${proxy.port}`} />
          <DetailRow icon={<Wifi size={9} />} label="Session" value={proxy.sessionId || '--'} />
          {proxy.location?.timezone && (
            <DetailRow icon={<Clock size={9} />} label="Timezone" value={proxy.location.timezone} />
          )}
          {proxy.location?.org && (
            <DetailRow icon={<Building2 size={9} />} label="Org" value={proxy.location.org} />
          )}
          {proxy.location?.as && (
            <DetailRow icon={<Activity size={9} />} label="AS" value={proxy.location.as} />
          )}
          {proxy.lastCheck && (
            <DetailRow icon={<Clock size={9} />} label="Last check" value={new Date(proxy.lastCheck).toLocaleTimeString()} />
          )}
        </motion.div>
      )}
    </motion.div>
  )
}

function ControlBtn({
  icon, label, onClick, busy, disabled, accent,
}: {
  icon: React.ReactNode
  label: string
  onClick: () => void
  busy?: boolean
  disabled?: boolean
  accent?: boolean
}) {
  return (
    <button
      onClick={onClick}
      disabled={disabled || busy}
      title={label}
      className={`h-7 px-2.5 flex items-center gap-1.5 text-[10px] font-medium rounded
                  border transition-all
                  disabled:opacity-20 disabled:cursor-not-allowed
                  ${accent
                    ? 'bg-surface-4 border-border-strong text-text-primary hover:bg-surface-5'
                    : 'bg-surface-3 border-border-default text-text-secondary hover:text-text-primary hover:bg-surface-4'
                  }`}
    >
      <span className={busy ? 'animate-spin' : ''}>{icon}</span>
      {label}
    </button>
  )
}

function DetailRow({ icon, label, value }: { icon: React.ReactNode; label: string; value: string }) {
  return (
    <div className="flex items-center justify-between pt-1.5">
      <span className="flex items-center gap-1.5 text-text-tertiary/70">
        {icon} {label}
      </span>
      <span className="font-mono text-text-secondary">{value}</span>
    </div>
  )
}
