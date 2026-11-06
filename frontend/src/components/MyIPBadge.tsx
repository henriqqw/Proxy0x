import { useState } from 'react'
import { motion } from 'framer-motion'
import { Globe, RefreshCw, Eye, EyeOff } from 'lucide-react'

interface Props {
  ip: string
  onRefresh: () => void
}

export default function MyIPBadge({ ip, onRefresh }: Props) {
  const [spinning, setSpinning] = useState(false)
  const [visible, setVisible] = useState(false)

  const handleRefresh = () => {
    setSpinning(true)
    onRefresh()
    setTimeout(() => setSpinning(false), 1000)
  }

  const masked = ip ? ip.replace(/\d/g, '●') : '...'

  return (
    <motion.div
      initial={{ opacity: 0, y: -4 }}
      animate={{ opacity: 1, y: 0 }}
      className="flex items-center justify-between py-2 px-3 rounded-lg bg-surface-1 border border-border-default"
    >
      <div className="flex items-center gap-2">
        <Globe size={12} className="text-accent/50" />
        <span className="text-[10px] text-text-tertiary uppercase tracking-wider">Your IP</span>
      </div>
      <div className="flex items-center gap-1.5">
        <span className="text-xs font-mono text-text-secondary select-none">
          {visible ? (ip || '...') : masked}
        </span>
        <button
          onClick={() => setVisible(!visible)}
          className="p-1 rounded hover:bg-surface-4 transition-colors text-text-tertiary hover:text-text-secondary"
          title={visible ? 'Ocultar IP' : 'Mostrar IP'}
        >
          <motion.div
            key={visible ? 'eye' : 'eyeoff'}
            initial={{ scale: 0.7, opacity: 0 }}
            animate={{ scale: 1, opacity: 1 }}
            transition={{ duration: 0.15 }}
          >
            {visible ? <Eye size={11} /> : <EyeOff size={11} />}
          </motion.div>
        </button>
        <button
          onClick={handleRefresh}
          className="p-1 rounded hover:bg-surface-4 transition-colors text-text-tertiary hover:text-text-secondary"
        >
          <RefreshCw size={10} className={spinning ? 'animate-spin' : ''} />
        </button>
      </div>
    </motion.div>
  )
}
