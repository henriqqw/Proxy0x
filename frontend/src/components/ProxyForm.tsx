import { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { Plus, ClipboardPaste, ChevronDown } from 'lucide-react'
import { AddProxy, AddProxiesBatch } from 'wailsjs/go/main/App'
import { useTranslation } from '../i18n/LanguageContext'

interface Props {
  onAdded: () => void
}

export default function ProxyForm({ onAdded }: Props) {
  const { t } = useTranslation()
  const [input, setInput] = useState('')
  const [expanded, setExpanded] = useState(false)
  const [batch, setBatch] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleAdd = async () => {
    if (!input.trim()) return
    setLoading(true)
    setError('')
    try {
      await AddProxy(input.trim())
      setInput('')
      onAdded()
    } catch (e: any) {
      setError(e?.message || t.failedAdd)
    } finally {
      setLoading(false)
    }
  }

  const handlePaste = async () => {
    try {
      const text = await navigator.clipboard.readText()
      if (text.includes('\n')) {
        setBatch(text)
        setExpanded(true)
      } else {
        setInput(text)
      }
    } catch {
      setError(t.clipboardDenied)
    }
  }

  const handleBatchAdd = async () => {
    if (!batch.trim()) return
    setLoading(true)
    setError('')
    try {
      await AddProxiesBatch(batch.trim())
      setBatch('')
      setExpanded(false)
      onAdded()
    } catch (e: any) {
      setError(e?.message || t.failedAddBatch)
    } finally {
      setLoading(false)
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') handleAdd()
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: -4 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.05 }}
      className="space-y-2"
    >
      <div className="flex gap-1.5">
        <div className="flex-1 relative">
          <input
            type="text"
            value={input}
            onChange={e => setInput(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder="host:port:user:pass"
            className="w-full h-8 px-3 pr-8 text-xs font-mono bg-surface-2 border border-border-default rounded-lg
                       text-text-primary placeholder:text-text-tertiary
                       focus:outline-none focus:border-border-strong focus:bg-surface-3 transition-all"
          />
          <button
            onClick={handlePaste}
            className="absolute right-1.5 top-1/2 -translate-y-1/2 p-1 rounded
                       text-text-tertiary hover:text-text-primary hover:bg-surface-4 transition-colors"
            title={t.paste}
            aria-label={t.paste}
          >
            <ClipboardPaste size={12} />
          </button>
        </div>
        <button
          onClick={handleAdd}
          disabled={loading || !input.trim()}
          aria-label={t.add}
          className="h-8 px-3 flex items-center gap-1.5 text-xs font-medium
                     bg-btn-accent border border-transparent rounded-lg
                     text-btn-accent-text hover:bg-btn-accent-hover font-semibold
                     disabled:opacity-30 disabled:cursor-not-allowed transition-all"
        >
          <Plus size={12} />
          {t.add}
        </button>
      </div>

      <button
        onClick={() => setExpanded(!expanded)}
        className="flex items-center gap-1 text-[10px] text-text-tertiary hover:text-text-secondary transition-colors"
      >
        <ChevronDown size={10} className={`transition-transform ${expanded ? 'rotate-180' : ''}`} />
        {t.batchImport}
      </button>

      <AnimatePresence>
        {expanded && (
          <motion.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: 'auto', opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.2 }}
            className="overflow-hidden"
          >
            <textarea
              value={batch}
              onChange={e => setBatch(e.target.value)}
              placeholder={t.batchPlaceholder}
              rows={4}
              className="w-full px-3 py-2 text-xs font-mono bg-surface-2 border border-border-default rounded-lg
                         text-text-primary placeholder:text-text-tertiary resize-none
                         focus:outline-none focus:border-border-strong focus:bg-surface-3 transition-all"
            />
            <button
              onClick={handleBatchAdd}
              disabled={loading || !batch.trim()}
              className="mt-1.5 h-7 px-3 text-xs font-medium bg-surface-4 border border-border-default rounded
                         text-text-secondary hover:text-text-primary hover:bg-surface-5
                         disabled:opacity-30 disabled:cursor-not-allowed transition-all"
            >
              {loading ? t.adding : t.importAll}
            </button>
          </motion.div>
        )}
      </AnimatePresence>

      <AnimatePresence>
        {error && (
          <motion.p
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="text-[10px] text-status-error"
          >
            {error}
          </motion.p>
        )}
      </AnimatePresence>
    </motion.div>
  )
}
