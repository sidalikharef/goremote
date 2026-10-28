import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { useLogout } from '@/hooks/useAuth'
import { Info, Server, FlaskConical } from 'lucide-react'

export function SettingsPage() {
  const logout = useLogout()

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Info className="h-5 w-5" />
            Application information
          </CardTitle>
          <CardDescription>GoRemote local instance</CardDescription>
        </CardHeader>
        <CardContent className="space-y-2 text-sm">
          <div className="flex justify-between">
            <span className="text-muted">Version</span>
            <span>dev</span>
          </div>
          <div className="flex justify-between">
            <span className="text-muted">Mode</span>
            <span>Local network</span>
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Server className="h-5 w-5" />
            Server status
          </CardTitle>
          <CardDescription>Connectivity and safety</CardDescription>
        </CardHeader>
        <CardContent className="space-y-2 text-sm text-muted">
          <div>Binding: localhost only by default</div>
          <div>Session auth required for API operations</div>
          <div>Origin/Host checks enforced</div>
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <FlaskConical className="h-5 w-5" />
            Experimental features
          </CardTitle>
          <CardDescription>Use with caution</CardDescription>
        </CardHeader>
        <CardContent className="space-y-2 text-sm text-muted">
          <div>VIDAA: Manual registration only, requires explicit enable, unverified on hardware.</div>
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <CardTitle>Session</CardTitle>
          <CardDescription>Sign out of this session</CardDescription>
        </CardHeader>
        <CardContent>
          <Button variant="secondary" onClick={() => logout.mutate()}>Logout</Button>
        </CardContent>
      </Card>
    </div>
  )
}
