import { motion } from 'framer-motion'
import { Unplug } from 'lucide-react'

export default function EmptyState() {
  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      className="flex flex-col items-center justify-center py-16 gap-3"
    >
      <Unplug size={28} className="text-text-tertiary" strokeWidth={1.5} />
      <p className="text-xs text-text-tertiary">No proxies added yet</p>
      <p className="text-[10px] text-text-tertiary/60">
        Add a proxy above to get started
      </p>
    </motion.div>
  )
}
