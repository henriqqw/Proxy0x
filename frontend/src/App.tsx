import { useState, useEffect, useCallback } from 'react'
import { AnimatePresence } from 'framer-motion'
import { EventsOn } from 'wailsjs/runtime/runtime'
import { GetProxies, GetMyIP } from 'wailsjs/go/main/App'
import { main } from 'wailsjs/go/models'
import TitleBar from './components/TitleBar'
import ProxyForm from './components/ProxyForm'
import ProxyCard from './components/ProxyCard'
import EmptyState from './components/EmptyState'
import MyIPBadge from './components/MyIPBadge'

export default function App() {
  const [proxies, setProxies] = useState<main.ProxyEntry[]>([])
  const [loading, setLoading] = useState(true)
  const [myIP, setMyIP] = useState<string>('')

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

    return () => {
      offUpdated()
      offChanged()
    }
  }, [fetchProxies])

  return (
    <div className="h-full flex flex-col bg-surface-0">
      <TitleBar />
      <div className="flex-1 overflow-y-auto px-4 pt-3 pb-4 space-y-3">
        <MyIPBadge ip={myIP} onRefresh={fetchMyIP} />
        <ProxyForm onAdded={fetchProxies} />
        <AnimatePresence mode="popLayout">
          {loading ? (
            <SkeletonList />
          ) : proxies.length === 0 ? (
            <EmptyState />
          ) : (
            proxies.map(proxy => (
              <ProxyCard key={proxy.id} proxy={proxy} onRemoved={fetchProxies} />
            ))
          )}
        </AnimatePresence>
      </div>
    </div>
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
