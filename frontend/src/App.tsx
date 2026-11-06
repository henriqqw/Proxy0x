import { useState, useEffect, useCallback } from 'react'
import { EventsOn } from 'wailsjs/runtime/runtime'
import { GetProxies, GetMyIP, ReorderProxies, GetGroups } from 'wailsjs/go/main/App'
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
import SettingsPage from './components/SettingsPage'

type Page = 'proxies' | 'settings'

function AppContent() {
  const [page, setPage] = useState<Page>('proxies')
  const [proxies, setProxies] = useState<main.ProxyEntry[]>([])
  const [loading, setLoading] = useState(true)
  const [myIP, setMyIP] = useState<string>('')
  const [groups, setGroups] = useState<string[]>([])
  const [activeGroup, setActiveGroup] = useState<string>('')

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

  const fetchGroups = useCallback(async () => {
    try {
      const g = await GetGroups()
      setGroups(g || [])
    } catch {}
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
    fetchGroups()

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
      fetchGroups()
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

    const onMouseNav = (e: MouseEvent) => {
      if (e.button === 3) setPage('proxies')
      if (e.button === 4) setPage('settings')
    }
    window.addEventListener('mouseup', onMouseNav)

    return () => {
      offUpdated()
      offChanged()
      offDropped()
      window.removeEventListener('mouseup', onMouseNav)
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

  const filtered = activeGroup ? proxies.filter(p => p.group === activeGroup) : proxies

  if (page === 'settings') {
    return (
      <>
        <SettingsPage
          onBack={() => setPage('proxies')}
          onImported={() => { fetchProxies(); setPage('proxies') }}
        />
      </>
    )
  }

  return (
    <div className="h-full flex flex-col bg-surface-0">
      <TitleBar onOpenSettings={() => setPage('settings')} />
      <div className="flex-1 overflow-y-auto px-4 pt-2 pb-3 space-y-2">
        <MyIPBadge ip={myIP} onRefresh={fetchMyIP} />
        <ProxyForm onAdded={() => { fetchProxies(); fetchGroups() }} />
        {groups.length > 0 && (
          <div className="flex gap-1 flex-wrap">
            <button
              onClick={() => setActiveGroup('')}
              className={`px-2 py-1 text-[10px] rounded-md border transition-all ${
                activeGroup === ''
                  ? 'bg-accent/10 text-accent border-accent/20'
                  : 'bg-surface-2 text-text-tertiary border-border-default hover:text-text-secondary'
              }`}
            >All</button>
            {groups.map(g => (
              <button
                key={g}
                onClick={() => setActiveGroup(g)}
                className={`px-2 py-1 text-[10px] rounded-md border transition-all ${
                  activeGroup === g
                    ? 'bg-accent/10 text-accent border-accent/20'
                    : 'bg-surface-2 text-text-tertiary border-border-default hover:text-text-secondary'
                }`}
              >{g}</button>
            ))}
          </div>
        )}
        {loading ? (
          <SkeletonList />
        ) : proxies.length === 0 ? (
          <EmptyState />
        ) : (
          <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
            <SortableContext items={filtered.map(p => p.id)} strategy={verticalListSortingStrategy}>
                {filtered.map(proxy => (
                  <ProxyCard key={proxy.id} proxy={proxy} onRemoved={() => { fetchProxies(); fetchGroups() }} />
                ))}
            </SortableContext>
          </DndContext>
        )}
      </div>
      <div className="shrink-0 flex items-center justify-center py-3 border-t border-border-subtle">
        <img src="/logo.png" alt="Proxy0x" className="h-7 opacity-[0.12] logo-footer" draggable={false} />
      </div>
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
