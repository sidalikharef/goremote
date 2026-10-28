import { useState } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useDevices, useSaveDiscovered } from '@/hooks/useDevices'
import { useVidaaRegister, useVidaaPairStart, useVidaaPairSubmit, useVidaaPairCancel, useVidaaUnpair } from '@/hooks/useVidaa'
import { Tv, PlusCircle } from 'lucide-react'

export function DevicesPage() {
  const { data: devices = [], refetch } = useDevices()
  const save = useSaveDiscovered()
  const vidaaRegister = useVidaaRegister()
  const vidaaPairStart = useVidaaPairStart()
  const vidaaPairSubmit = useVidaaPairSubmit()
  const vidaaPairCancel = useVidaaPairCancel()
  const vidaaUnpair = useVidaaUnpair()

  const [manual, setManual] = useState({ id: '', name: '', address: '', driver: 'roku' })
  const [vidaa, setVidaa] = useState({ id: '', name: '', address: '', enable: false })
  const [pin, setPin] = useState('')

  const addManual = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await save.mutateAsync({ id: manual.id || manual.address, name: manual.name || 'TV', address: manual.address, driver: manual.driver })
      setManual({ id: '', name: '', address: '', driver: 'roku' })
      refetch()
    } catch (err) {
      // ignore
    }
  }

  const registerVidaa = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await vidaaRegister.mutateAsync({ id: vidaa.id || vidaa.address, name: vidaa.name || 'VIDAA TV', address: vidaa.address, enable: vidaa.enable })
      setVidaa({ id: '', name: '', address: '', enable: false })
      refetch()
    } catch (err) {
      // ignore
    }
  }

  return (
    <div className="grid gap-6 lg:grid-cols-2">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Tv className="h-5 w-5" />
            My TVs
          </CardTitle>
          <CardDescription>Saved devices in this GoRemote instance</CardDescription>
        </CardHeader>
        <CardContent className="space-y-3">
          {devices.length === 0 && <div className="text-sm text-muted">No devices registered yet.</div>}
          {devices.map((d) => (
            <div key={d.id} className="flex items-center justify-between rounded-lg border border-border px-4 py-3">
              <div>
                <div className="font-medium">{d.name}</div>
                <div className="text-xs text-muted">
                  {d.driver} • {d.state} • {d.id}
                </div>
              </div>
            </div>
          ))}
        </CardContent>
      </Card>
      <div className="space-y-6">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <PlusCircle className="h-5 w-5" />
              Add TV manually
            </CardTitle>
            <CardDescription>Register a device by IP address</CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={addManual} className="space-y-3">
              <Input placeholder="Device name" value={manual.name} onChange={(e) => setManual({ ...manual, name: e.target.value })} />
              <Input placeholder="IP address" value={manual.address} onChange={(e) => setManual({ ...manual, address: e.target.value })} />
              <select value={manual.driver} onChange={(e) => setManual({ ...manual, driver: e.target.value })} className="h-10 w-full rounded-lg border border-border bg-surface px-3 text-sm">
                <option value="roku">Roku</option>
                <option value="vidaa">VIDAA (experimental)</option>
              </select>
              <Button type="submit" className="w-full">Add device</Button>
            </form>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Experimental VIDAA</CardTitle>
            <CardDescription>Hisense VIDAA support - unverified. Enable explicitly.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-3">
            <form onSubmit={registerVidaa} className="space-y-3">
              <Input placeholder="Device ID" value={vidaa.id} onChange={(e) => setVidaa({ ...vidaa, id: e.target.value })} />
              <Input placeholder="Name" value={vidaa.name} onChange={(e) => setVidaa({ ...vidaa, name: e.target.value })} />
              <Input placeholder="IP Address" value={vidaa.address} onChange={(e) => setVidaa({ ...vidaa, address: e.target.value })} />
              <label className="flex items-center gap-2 text-sm text-muted">
                <input type="checkbox" checked={vidaa.enable} onChange={(e) => setVidaa({ ...vidaa, enable: e.target.checked })} />
                Enable experimental
              </label>
              <Button type="submit" variant="secondary" className="w-full">Register VIDAA</Button>
            </form>
            <div className="grid grid-cols-2 gap-2">
              <Button variant="secondary" onClick={() => vidaaPairStart.mutate(undefined)}>Start Pairing</Button>
              <Button variant="secondary" onClick={() => vidaaPairCancel.mutate()}>Cancel</Button>
            </div>
            <div className="flex gap-2">
              <Input placeholder="PIN" maxLength={4} value={pin} onChange={(e) => setPin(e.target.value)} />
              <Button variant="secondary" onClick={() => vidaaPairSubmit.mutate({ pin })}>Submit</Button>
            </div>
            <Button variant="ghost" className="w-full" onClick={() => vidaaUnpair.mutate()}>Unpair</Button>
            <div className="text-xs text-muted">PINs are not stored in browser storage. Pairing responses depend on backend state.</div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
