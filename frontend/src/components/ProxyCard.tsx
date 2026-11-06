import { useState } from 'react'
import {
  Power, X, RotateCw,
  Trash2, Copy, MapPin, Clock, Wifi, Server,
  ChevronDown, Activity, Building2, GripVertical, Timer,
  Gauge, Tag
} from 'lucide-react'
import {
  ConnectProxy, DisconnectProxy, RotateProxy,
  RemoveProxy, BenchmarkProxy, SetProxyGroup,
} from 'wailsjs/go/main/App'
import { main } from 'wailsjs/go/models'
import { useSortable } from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'
import { useTranslation } from '../i18n/LanguageContext'
import StatusBadge from './StatusBadge'
import FlagIcon from './FlagIcon'

interface Props {
  proxy: main.ProxyEntry
  onRemoved: () => void
}

function formatUptime(totalSeconds: number, connectedAt?: string): string {
  let seconds = totalSeconds
  if (connectedAt) {
    const start = new Date(connectedAt).getTime()
    if (!isNaN(start)) {
      seconds += Math.floor((Date.now() - start) / 1000)
    }
  }
  if (seconds < 60) return `${seconds}s`
  const mins = Math.floor(seconds / 60)
  if (mins < 60) return `${mins}m ${seconds % 60}s`
  const hrs = Math.floor(mins / 60)
  return `${hrs}h ${mins % 60}m`
}

export default function ProxyCard({ proxy, onRemoved }: Props) {
  const { t } = useTranslation()
  const [busy, setBusy] = useState<string | null>(null)
  const [expanded, setExpanded] = useState(false)
  const [copied, setCopied] = useState(false)
  const [bench, setBench] = useState<{connectMs?: number; downloadMbps?: number} | null>(null)
  const [editGroup, setEditGroup] = useState(false)
  const [groupInput, setGroupInput] = useState(proxy.group || '')

  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: proxy.id })

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    zIndex: isDragging ? 50 : undefined,
    opacity: isDragging ? 0.5 : undefined,
  }

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
    <div
      ref={setNodeRef}
      style={style}
      className={`rounded-lg border transition-all animate-fade-in ${
        isActive
          ? 'bg-surface-2 border-accent/20 shadow-[0_1px_12px_rgba(104,224,69,0.06)]'
          : 'bg-surface-1 border-border-default hover:border-border-strong'
      }`}
    >
      {/* Header */}
      <div className="px-3 pt-3 pb-2">
        <div className="flex items-center justify-between mb-2">
          <div className="flex items-center gap-2 min-w-0 flex-1">
            <button
              {...attributes}
              {...listeners}
              className="cursor-grab active:cursor-grabbing p-0.5 -ml-1 text-text-tertiary/30 hover:text-text-tertiary transition-colors touch-none"
              tabIndex={-1}
              aria-label="Drag to reorder"
            >
              <GripVertical size={12} />
            </button>
            <StatusBadge status={status} />
            <span className="text-[10px] font-mono text-text-tertiary truncate">
              {proxy.host}:{proxy.port}
            </span>
          </div>
          <div className="flex items-center gap-1">
            <span className={`px-1.5 py-0.5 text-[8px] font-bold uppercase tracking-wider rounded ${
              proxy.protocol === 'socks5'
                ? 'bg-btn-accent text-btn-accent-text'
                : 'bg-[#505050] text-[#d4d4d4]'
            }`}>
              {proxy.protocol === 'socks5' ? 'S5' : 'HTTP'}
            </span>
            {proxy.group && (
              <span className="px-1.5 py-0.5 text-[8px] font-medium rounded bg-accent/10 text-accent border border-accent/20 truncate max-w-[60px]">
                {proxy.group}
              </span>
            )}
            <button
              onClick={() => { RemoveProxy(proxy.id); onRemoved() }}
              className="p-1 rounded text-text-tertiary/40 hover:text-status-error hover:bg-surface-4 transition-all"
              aria-label="Remove proxy"
            >
              <Trash2 size={11} />
            </button>
          </div>
        </div>

        {/* IP Display */}
        {proxy.currentIp && (
          <div className="flex items-center gap-2 mb-2">
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
                  copied ? 'text-accent' : 'text-text-tertiary opacity-0 group-hover:opacity-100'
                }`}
              />
            </button>
            {proxy.latency > 0 && (
              <span className="ml-auto flex items-center gap-1 text-[10px] text-text-tertiary">
                <Activity size={9} />
                {proxy.latency}ms
              </span>
            )}
          </div>
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
              label={t.connect}
              onClick={() => action('connect', () => ConnectProxy(proxy.id))}
              busy={busy === 'connect'}
              accent
            />
          ) : (
            <ControlBtn
              icon={<X size={12} />}
              label={t.stop}
              onClick={() => action('disconnect', () => DisconnectProxy(proxy.id))}
              busy={busy === 'disconnect'}
            />
          )}
          <ControlBtn
            icon={<RotateCw size={12} />}
            label={t.rotate}
            onClick={() => action('rotate', () => RotateProxy(proxy.id))}
            busy={busy === 'rotate'}
            disabled={!isActive && !isConnecting}
          />
          <ControlBtn
            icon={<Gauge size={12} />}
            label={t.speedTest}
            onClick={() => action('bench', async () => {
              const r = await BenchmarkProxy(proxy.id)
              setBench(r)
            })}
            busy={busy === 'bench'}
            disabled={!isActive && !isConnecting}
          />

          {isActive && (proxy.totalUptime > 0 || proxy.connectedAt) && (
            <span className="ml-auto flex items-center gap-1 text-[10px] text-text-tertiary">
              <Timer size={9} />
              {formatUptime(proxy.totalUptime || 0, proxy.connectedAt)}
            </span>
          )}
        </div>
      </div>

      {/* Expandable details */}
      <button
        onClick={() => setExpanded(!expanded)}
        className="w-full flex items-center justify-center py-1.5 border-t border-border-subtle
                   text-text-tertiary/40 hover:text-text-tertiary hover:bg-surface-3/60 transition-all"
        aria-label={expanded ? 'Collapse details' : 'Expand details'}
      >
        <ChevronDown size={10} className={`transition-transform ${expanded ? 'rotate-180' : ''}`} />
      </button>

      {expanded && (
        <div className="px-3 pb-3 space-y-1.5 text-[10px] text-text-tertiary border-t border-border-subtle animate-expand">
          <DetailRow icon={<Server size={9} />} label={t.host} value={`${proxy.host}:${proxy.port}`} />
          <DetailRow icon={<Wifi size={9} />} label={t.protocol} value={(proxy.protocol || 'http').toUpperCase()} />
          <DetailRow icon={<Wifi size={9} />} label={t.session} value={proxy.sessionId || '--'} />
          <div className="flex items-center justify-between pt-1.5">
            <span className="flex items-center gap-1.5 text-text-tertiary/70"><Tag size={9} /> {t.group}</span>
            {editGroup ? (
              <input
                autoFocus
                value={groupInput}
                onChange={e => setGroupInput(e.target.value)}
                onBlur={() => { SetProxyGroup(proxy.id, groupInput.trim()); setEditGroup(false) }}
                onKeyDown={e => { if (e.key === 'Enter') { SetProxyGroup(proxy.id, groupInput.trim()); setEditGroup(false) } }}
                className="w-20 h-5 px-1.5 text-[10px] font-mono bg-surface-2 border border-border-strong rounded text-text-primary focus:outline-none text-right"
              />
            ) : (
              <button onClick={() => setEditGroup(true)} className="font-mono text-text-secondary hover:text-text-primary transition-colors">
                {proxy.group || '--'}
              </button>
            )}
          </div>
          {bench && (
            <>
              <DetailRow icon={<Gauge size={9} />} label={t.connectTime} value={`${bench.connectMs ?? '--'}ms`} />
              <DetailRow icon={<Gauge size={9} />} label={t.downloadSpeed} value={bench.downloadMbps != null ? `${bench.downloadMbps.toFixed(1)} Mbps` : '--'} />
            </>
          )}
          {proxy.connectedAt && (
            <DetailRow icon={<Clock size={9} />} label={t.connectedAt} value={new Date(proxy.connectedAt).toLocaleTimeString()} />
          )}
          {(proxy.totalUptime > 0 || proxy.connectedAt) && (
            <DetailRow icon={<Timer size={9} />} label={t.totalUptime} value={formatUptime(proxy.totalUptime || 0, proxy.connectedAt)} />
          )}
          {proxy.location?.timezone && (
            <DetailRow icon={<Clock size={9} />} label={t.timezone} value={proxy.location.timezone} />
          )}
          {proxy.location?.org && (
            <DetailRow icon={<Building2 size={9} />} label={t.org} value={proxy.location.org} />
          )}
          {proxy.location?.as && (
            <DetailRow icon={<Activity size={9} />} label={t.asNumber} value={proxy.location.as} />
          )}
          {proxy.lastCheck && (
            <DetailRow icon={<Clock size={9} />} label={t.lastCheck} value={new Date(proxy.lastCheck).toLocaleTimeString()} />
          )}
        </div>
      )}
    </div>
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
      aria-label={label}
      className={`h-7 px-2.5 flex items-center gap-1.5 text-[10px] font-medium rounded
                  border transition-all
                  disabled:opacity-20 disabled:cursor-not-allowed
                  ${accent
                    ? 'bg-btn-accent border-transparent text-btn-accent-text hover:bg-btn-accent-hover font-semibold'
                    : 'bg-transparent border-border-strong text-text-secondary hover:text-text-primary hover:bg-surface-2'
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
