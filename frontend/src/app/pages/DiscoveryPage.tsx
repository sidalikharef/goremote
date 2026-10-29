import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { useDiscover, useDiscoverResults } from '@/hooks/useDiscovery'
import { useSaveDiscovered } from '@/hooks/useDevices'
import { Search, Radar } from 'lucide-react'
import { useState } from 'react'

export function DiscoveryPage() {
  const discover = useDiscover()
  const results = useDiscoverResults().data || []
  const save = useSaveDiscovered()
  const [addingId, setAddingId] = useState<string | null>(null)
  const [addedId, setAddedId] = useState<string | null>(null)
  const [errorId, setErrorId] = useState<string | null>(null)

  const handleScan = () => {
    setAddedId(null)
    setErrorId(null)
    discover.mutate()
  }

  const handleAdd = async (d: any) => {
    const id = d.id || d.address
    setAddingId(id)
    setErrorId(null)
    try {
      await save.mutateAsync({
        id,
        name: d.name || 'TV',
        address: d.address,
        driver: d.driver_hint || d.driver || 'roku',
      })
      setAddedId(id)
    } catch (e) {
      setErrorId(id)
    } finally {
      setAddingId(null)
    }
  }

  const getDriver = (d: any) => d.driver_hint || d.driver || ''
  const getSource = (d: any) => d.discovery_source || d.source || ''

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Search className="h-5 w-5" />
            Discovery
          </CardTitle>
          <CardDescription>Find compatible TVs on your local network</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center gap-2">
            <Button onClick={handleScan} disabled={discover.isPending}>
              {discover.isPending ? 'Scanning...' : 'Scan for TVs'}
            </Button>
            {results.length > 0 && !discover.isPending && (
              <div className="text-xs text-muted">{results.length} candidate{results.length === 1 ? '' : 's'} found</div>
            )}
          </div>
          {discover.isPending && <div className="text-sm text-muted">Scanning your network...</div>}
          {!discover.isPending && results.length === 0 && discover.isSuccess && (
            <div className="rounded-lg border border-border p-4 text-sm">
              <div className="font-medium">No compatible TVs detected</div>
              <div className="mt-1 text-muted">Make sure your TV is powered on and connected to the same network. Try again or add it manually.</div>
            </div>
          )}
          {discover.isError && (
            <div className="rounded-lg border border-error/40 p-4 text-sm">
              <div className="font-medium text-error">Discovery failed</div>
              <div className="mt-1 text-muted">Check authentication, timeout, or network connectivity and retry.</div>
              <Button variant="secondary" className="mt-3" onClick={handleScan}>Retry</Button>
            </div>
          )}
          <div className="space-y-2">
            {results.map((d: any) => {
              const id = d.id || d.address
              const driver = getDriver(d)
              const source = getSource(d)
              const isAdding = addingId === id
              const isAdded = addedId === id
              const isError = errorId === id
              return (
                <div key={id} className="flex flex-col gap-2 rounded-lg border border-border px-4 py-3 md:flex-row md:items-center md:justify-between">
                  <div className="flex items-center gap-3">
                    <div className="flex h-9 w-9 items-center justify-center rounded-full bg-surface-elevated">
                      <Radar className="h-4 w-4" />
                    </div>
                    <div>
                      <div className="font-medium">{d.name || 'Unidentified TV candidate'}</div>
                      <div className="text-xs text-muted">
                        {d.manufacturer ? `${d.manufacturer}` : 'Unknown manufacturer'}{driver ? ` • ${driver}` : ''}{source ? ` • ${source}` : ''}{d.address ? ` • ${d.address}` : ''}
                      </div>
                      {d.experimental && (
                        <div className="text-xs text-amber-400">Experimental — pairing and compatibility unverified</div>
                      )}
                      {d.requires_pairing && driver === 'vidaa' && (
                        <div className="text-xs text-muted">Pairing required</div>
                      )}
                      {isError && <div className="text-xs text-error">Registration failed</div>}
                      {isAdded && <div className="text-xs text-success">Added to My TVs</div>}
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    {isAdded ? (
                      <Button size="sm" variant="secondary" disabled>Added</Button>
                    ) : (
                      <Button size="sm" onClick={() => handleAdd(d)} disabled={isAdding || save.isPending}>
                        {isAdding ? 'Adding...' : 'Add TV'}
                      </Button>
                    )}
                  </div>
                </div>
              )
            })}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
