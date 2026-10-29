import { useNavigate, useOutletContext } from 'react-router-dom'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { useDevices, useKeypress } from '@/hooks/useDevices'
import { ChevronUp, ChevronDown, ChevronLeft, ChevronRight, Home, ArrowLeft, Power, Volume2, VolumeX } from 'lucide-react'
import { useState } from 'react'
import { cn } from '@/lib/utils/cn'

export function RemotePage() {
  const navigate = useNavigate()
  const { selectedDevice } = useOutletContext<{ selectedDevice: string }>()
  const { data: devices = [] } = useDevices()
  const keypress = useKeypress()
  const [lastKey, setLastKey] = useState<string | null>(null)

  const device = devices.find((d) => d.id === selectedDevice)

  const send = (key: string) => {
    setLastKey(key)
    keypress.mutate({ key, device: selectedDevice || null })
  }

  if (!devices.length || !selectedDevice) {
    return (
      <div className="flex flex-col items-center justify-center gap-6 py-12 text-center">
        <div>
          <h1 className="text-2xl font-semibold">Connect your first TV</h1>
          <p className="mt-2 text-muted">Find compatible TVs on your local network or add one manually.</p>
        </div>
        <div className="flex flex-wrap items-center justify-center gap-3">
          <Button onClick={() => navigate('/discovery')}>Scan for TVs</Button>
          <Button variant="secondary" onClick={() => navigate('/devices')}>Add TV manually</Button>
        </div>
        <Card className="w-full max-w-md">
          <CardHeader className="text-left">
            <CardTitle>Quick steps</CardTitle>
            <CardDescription>Get started in minutes</CardDescription>
          </CardHeader>
          <CardContent className="space-y-3 text-left text-sm">
            <div>1. Power on your TV.</div>
            <div>2. Connect it to the same network.</div>
            <div>3. Scan and select your device.</div>
            <div className="pt-2 text-xs text-muted">Roku discovery is supported. Experimental VIDAA devices may require manual setup.</div>
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className="grid gap-6 lg:grid-cols-[1.2fr_1fr]">
      <Card>
        <CardHeader>
          <CardTitle>Remote Control</CardTitle>
          <CardDescription>
            {device?.name} • {device?.driver} • {device?.state}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="mx-auto flex max-w-sm flex-col items-center gap-6">
            <div className="flex w-full items-center justify-between">
              <Button variant="secondary" size="sm" onClick={() => send('Power')}>
                <Power className="h-4 w-4" />
                Power
              </Button>
              <Button variant="secondary" size="sm" onClick={() => send('Home')}>
                <Home className="h-4 w-4" />
                Home
              </Button>
              <Button variant="secondary" size="sm" onClick={() => send('Back')}>
                <ArrowLeft className="h-4 w-4" />
                Back
              </Button>
            </div>
            <div className="relative flex h-56 w-56 items-center justify-center rounded-full border border-border bg-surface-elevated shadow-inner">
              <button className="absolute top-6 flex h-12 w-12 items-center justify-center rounded-full hover:bg-surface" onClick={() => send('Up')} aria-label="Up">
                <ChevronUp className="h-6 w-6" />
              </button>
              <button className="absolute bottom-6 flex h-12 w-12 items-center justify-center rounded-full hover:bg-surface" onClick={() => send('Down')} aria-label="Down">
                <ChevronDown className="h-6 w-6" />
              </button>
              <button className="absolute left-6 flex h-12 w-12 items-center justify-center rounded-full hover:bg-surface" onClick={() => send('Left')} aria-label="Left">
                <ChevronLeft className="h-6 w-6" />
              </button>
              <button className="absolute right-6 flex h-12 w-12 items-center justify-center rounded-full hover:bg-surface" onClick={() => send('Right')} aria-label="Right">
                <ChevronRight className="h-6 w-6" />
              </button>
              <button className="flex h-20 w-20 items-center justify-center rounded-full bg-accent text-white shadow-lg hover:bg-accent/90" onClick={() => send('Select')} aria-label="Select">
                OK
              </button>
            </div>
            <div className="grid w-full grid-cols-3 gap-3">
              <Button variant="secondary" onClick={() => send('Rev')}>Rewind</Button>
              <Button variant="secondary" onClick={() => send('Play')}>Play/Pause</Button>
              <Button variant="secondary" onClick={() => send('Fwd')}>Fast Forward</Button>
            </div>
            <div className="grid w-full grid-cols-3 gap-3">
              <Button variant="secondary" onClick={() => send('VolumeDown')}>
                <Volume2 className="mr-2 h-4 w-4" />
                Vol -
              </Button>
              <Button variant="secondary" onClick={() => send('VolumeMute')}>
                <VolumeX className="h-4 w-4" />
              </Button>
              <Button variant="secondary" onClick={() => send('VolumeUp')}>Vol +</Button>
            </div>
            {lastKey && <div className="text-xs text-muted">Last command: {lastKey}</div>}
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <CardTitle>Device Info</CardTitle>
          <CardDescription>Selected television</CardDescription>
        </CardHeader>
        <CardContent className="space-y-2 text-sm">
          <div className="flex justify-between">
            <span className="text-muted">Name</span>
            <span>{device?.name}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-muted">Driver</span>
            <span className="capitalize">{device?.driver}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-muted">Status</span>
            <span className={cn(device?.state === 'connected' ? 'text-success' : 'text-muted')}>{device?.state}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-muted">ID</span>
            <span className="text-xs text-muted">{device?.id}</span>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
