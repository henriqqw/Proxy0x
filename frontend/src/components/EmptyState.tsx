import { motion } from 'framer-motion'

export default function EmptyState() {
  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      className="flex flex-col items-center justify-center py-16 gap-4"
    >
      <motion.img
        src="/icon-64.png"
        alt="Proxy0x"
        className="w-14 h-14 opacity-30"
        draggable={false}
        animate={{ opacity: [0.2, 0.35, 0.2] }}
        transition={{ duration: 3, repeat: Infinity, ease: 'easeInOut' }}
      />
      <div className="text-center space-y-1">
        <p className="text-xs text-text-tertiary">Nenhuma proxy adicionada</p>
        <p className="text-[10px] text-text-tertiary/50">
          Adicione uma proxy acima para começar
        </p>
      </div>
    </motion.div>
  )
}
