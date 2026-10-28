import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { useDiscover, useDiscoverResults } from '@/hooks/useDiscovery'
import { useSaveDiscovered } from '@/hooks/useDevices'
import { Search, Radar } from 'lucide-react'

export function DiscoveryPage() {
  const discover = useDiscover()
  const results = useDiscoverResults().data || []
  const save = useSaveDiscovered()

  const handleScan = () => {
    discover.mutate()
  }

  const handleAdd = (d: any) => {
    save.mutate({
      id: d.id || d.address,
      name: d.name || 'TV',
      address: d.address,
      driver: d.driver_hint || 'roku',
    })
  }

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
          <Button onClick={handleScan} disabled={discover.isPending}>
            {discover.isPending ? 'Scanning...' : 'Scan for TVs'}
          </Button>
          {discover.isPending && <div className="text-sm text-muted">Scanning your network...</div>}
          {!discover.isPending && results.length === 0 && discover.isSuccess && (
            <div className="rounded-lg border border-border p-4 text-sm">
              <div className="font-medium">No TVs found</div>
              <div className="mt-1 text-muted">Make sure your TV is powered on and connected to the same network. Try again.</div>
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
            {results.map((d: any) => (
              <div key={d.id} className="flex items-center justify-between rounded-lg border border-border px-4 py-3">
                <div className="flex items-center gap-3">
                  <div className="flex h-9 w-9 items-center justify-center rounded-full bg-surface-elevated">
                    <Radar className="h-4 w-4" />
                  </div>
                  <div>
                    <div className="font-medium">{d.name}</div>
                    <div className="text-xs text-muted">
                      {d.address} {d.driver_hint ? `• ${d.driver_hint}` : ''}
                    </div>
                  </div>
                </div>
                <Button size="sm" onClick={() => handleAdd(d)} disabled={save.isPending}>Add TV</Button>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
