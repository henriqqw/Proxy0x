import { Minus, Square, X } from 'lucide-react'
import { MinimizeWindow, CloseWindow, ToggleMaximize } from 'wailsjs/go/main/App'

export default function TitleBar() {
  return (
    <div
      className="flex items-center justify-between h-10 px-3 bg-surface-1/80 border-b border-border-default shrink-0 backdrop-blur-sm"
      style={{ '--wails-draggable': 'drag' } as React.CSSProperties}
    >
      <div className="flex items-center">
        <img src="/logo.png" alt="Proxy0x" className="h-4" draggable={false} />
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
