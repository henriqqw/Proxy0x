import { useTranslation } from '../i18n/LanguageContext'

export default function EmptyState() {
  const { t } = useTranslation()

  return (
    <div className="flex flex-col items-center justify-center py-16 gap-4 animate-fade-in">
      <img
        src="/icon-64.png"
        alt="Proxy0x"
        className="w-14 h-14 opacity-30 mascot-empty animate-breathe"
        draggable={false}
      />
      <div className="text-center space-y-1">
        <p className="text-xs text-text-tertiary">{t.noProxy}</p>
        <p className="text-[10px] text-text-tertiary/70">{t.addProxyHint}</p>
      </div>
    </div>
  )
}
