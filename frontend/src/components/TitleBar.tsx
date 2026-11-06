import { Minus, Square, X } from 'lucide-react'
import { MinimizeWindow, CloseWindow, ToggleMaximize } from 'wailsjs/go/main/App'

export default function TitleBar() {
  return (
    <div
      className="flex items-center justify-between h-10 px-3 bg-surface-1/80 border-b border-border-default shrink-0 backdrop-blur-sm"
      style={{ '--wails-draggable': 'drag' } as React.CSSProperties}
    >
      <div className="flex items-center gap-2.5">
        <svg className="w-4 h-4 text-text-secondary" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
          <circle cx="16" cy="16" r="7" stroke="currentColor" strokeWidth="1.5" fill="none"/>
          <circle cx="16" cy="16" r="2.5" fill="currentColor"/>
          <path d="M16 9V6M16 26v-3M23 16h3M6 16h3" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round"/>
        </svg>
        <span className="text-[11px] font-medium text-text-tertiary tracking-widest uppercase">
          Proxy<span className="text-text-secondary">UI</span>
        </span>
      </div>
      <div className="flex items-center gap-0.5 no-drag">
        <button
          onClick={() => MinimizeWindow()}
          className="p-1.5 rounded hover:bg-surface-4 transition-colors text-text-tertiary hover:text-text-secondary"
        >
          <Minus size={12} />
        </button>
        <button
          onClick={() => ToggleMaximize()}
          className="p-1.5 rounded hover:bg-surface-4 transition-colors text-text-tertiary hover:text-text-secondary"
        >
          <Square size={10} />
        </button>
        <button
          onClick={() => CloseWindow()}
          className="p-1.5 rounded hover:bg-surface-4 transition-colors text-text-tertiary hover:text-red-400"
        >
          <X size={12} />
        </button>
      </div>
    </div>
  )
}
