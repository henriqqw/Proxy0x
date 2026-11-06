interface Props {
  status: 'offline' | 'connecting' | 'connected' | 'error'
}

const config = {
  offline:    { label: 'Offline',     color: 'bg-status-offline',    text: 'text-text-tertiary' },
  connecting: { label: 'Connecting',  color: 'bg-status-connecting', text: 'text-status-connecting' },
  connected:  { label: 'Connected',   color: 'bg-status-online',     text: 'text-status-online' },
  error:      { label: 'Error',       color: 'bg-status-error',      text: 'text-status-error' },
}

export default function StatusBadge({ status }: Props) {
  const c = config[status]

  return (
    <div className={`flex items-center gap-1.5 ${c.text}`}>
      <div
        className={`w-1.5 h-1.5 rounded-full ${c.color} ${
          status === 'connecting' ? 'animate-pulse' : ''
        }`}
      />
      <span className="text-[10px] font-medium uppercase tracking-wider">{c.label}</span>
    </div>
  )
}
