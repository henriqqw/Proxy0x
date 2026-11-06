import { useState, useEffect, useCallback } from 'react'
import { AnimatePresence } from 'framer-motion'
import { EventsOn } from 'wailsjs/runtime/runtime'
import { GetProxies, GetMyIP, ReorderProxies } from 'wailsjs/go/main/App'
import { main } from 'wailsjs/go/models'
import {
  DndContext, closestCenter, PointerSensor, useSensor, useSensors,
  type DragEndEvent,
} from '@dnd-kit/core'
import {
  SortableContext, verticalListSortingStrategy, arrayMove,
} from '@dnd-kit/sortable'
import { LanguageProvider } from './i18n/LanguageContext'
import TitleBar from './components/TitleBar'
import ProxyForm from './components/ProxyForm'
import ProxyCard from './components/ProxyCard'
import EmptyState from './components/EmptyState'
import MyIPBadge from './components/MyIPBadge'
import ContextMenu from './components/ContextMenu'

function AppContent() {
  const [proxies, setProxies] = useState<main.ProxyEntry[]>([])
  const [loading, setLoading] = useState(true)
  const [myIP, setMyIP] = useState<string>('')

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 5 } })
  )

  const fetchProxies = useCallback(async () => {
    try {
      const list = await GetProxies()
      setProxies(list || [])
    } catch {
      // not in wails runtime (dev mode)
    } finally {
      setLoading(false)
    }
  }, [])

  const fetchMyIP = useCallback(async () => {
    try {
      const ip = await GetMyIP()
      setMyIP(ip)
    } catch {
      setMyIP('--')
    }
  }, [])

  useEffect(() => {
    fetchProxies()
    fetchMyIP()

    const offUpdated = EventsOn('proxy-updated', (proxy: main.ProxyEntry) => {
      setProxies(prev => {
        const idx = prev.findIndex(p => p.id === proxy.id)
        if (idx >= 0) {
          const next = [...prev]
          next[idx] = proxy
          return next
        }
        return [...prev, proxy]
      })
    })

    const offChanged = EventsOn('proxies-changed', () => {
      fetchProxies()
    })

    const offDropped = EventsOn('proxy-dropped', (data: { host: string }) => {
      if ('Notification' in window && Notification.permission === 'granted') {
        new Notification('Proxy0x', { body: `Proxy ${data.host} dropped`, icon: '/icon-64.png' })
      } else if ('Notification' in window && Notification.permission !== 'denied') {
        Notification.requestPermission().then(p => {
          if (p === 'granted') {
            new Notification('Proxy0x', { body: `Proxy ${data.host} dropped`, icon: '/icon-64.png' })
          }
        })
      }
    })

    return () => {
      offUpdated()
      offChanged()
      offDropped()
    }
  }, [fetchProxies])

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event
    if (!over || active.id === over.id) return

    setProxies(prev => {
      const oldIndex = prev.findIndex(p => p.id === active.id)
      const newIndex = prev.findIndex(p => p.id === over.id)
      const reordered = arrayMove(prev, oldIndex, newIndex)
      ReorderProxies(reordered.map(p => p.id))
      return reordered
    })
  }

  return (
    <div className="h-full flex flex-col bg-surface-0">
      <TitleBar />
      <div className="flex-1 overflow-y-auto px-4 pt-2 pb-3 space-y-2">
        <MyIPBadge ip={myIP} onRefresh={fetchMyIP} />
        <ProxyForm onAdded={fetchProxies} />
        {loading ? (
          <SkeletonList />
        ) : proxies.length === 0 ? (
          <EmptyState />
        ) : (
          <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
            <SortableContext items={proxies.map(p => p.id)} strategy={verticalListSortingStrategy}>
              <AnimatePresence mode="popLayout">
                {proxies.map(proxy => (
                  <ProxyCard key={proxy.id} proxy={proxy} onRemoved={fetchProxies} />
                ))}
              </AnimatePresence>
            </SortableContext>
          </DndContext>
        )}
      </div>
      <div className="shrink-0 flex items-center justify-center py-3 border-t border-border-subtle">
        <img src="/logo.png" alt="Proxy0x" className="h-7 opacity-[0.12] logo-footer" draggable={false} />
      </div>
      <ContextMenu />
    </div>
  )
}

export default function App() {
  return (
    <LanguageProvider>
      <AppContent />
    </LanguageProvider>
  )
}

function SkeletonList() {
  return (
    <>
      {[0, 1, 2].map(i => (
        <div key={i} className="animate-pulse rounded-lg bg-surface-1 border border-border-subtle p-4 space-y-3">
          <div className="flex items-center gap-3">
            <div className="w-2 h-2 rounded-full bg-surface-5" />
            <div className="h-3 w-32 rounded bg-surface-4" />
            <div className="ml-auto h-3 w-16 rounded bg-surface-3" />
          </div>
          <div className="h-3 w-48 rounded bg-surface-3" />
          <div className="flex gap-2">
            <div className="h-7 w-16 rounded bg-surface-4" />
            <div className="h-7 w-16 rounded bg-surface-4" />
            <div className="h-7 w-16 rounded bg-surface-4" />
          </div>
        </div>
      ))}
    </>
  )
}
