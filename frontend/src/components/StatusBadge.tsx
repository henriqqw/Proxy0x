import { motion } from 'framer-motion'

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
      <motion.div
        className={`w-1.5 h-1.5 rounded-full ${c.color}`}
        animate={status === 'connecting' ? { opacity: [1, 0.3, 1] } : {}}
        transition={status === 'connecting' ? { duration: 1.2, repeat: Infinity } : {}}
      />
      <span className="text-[10px] font-medium uppercase tracking-wider">{c.label}</span>
    </div>
  )
}
