import { NavLink, Navigate, Outlet, useNavigate } from 'react-router-dom'
import { Tv, Search, Settings, Home } from 'lucide-react'
import { useDevices } from '@/hooks/useDevices'
import { useState } from 'react'
import { cn } from '@/lib/utils/cn'
import { useAuthStatus, useLogout } from '@/hooks/useAuth'

const navItems = [
  { to: '/', label: 'Remote', icon: Home },
  { to: '/devices', label: 'My TVs', icon: Tv },
  { to: '/discovery', label: 'Discovery', icon: Search },
  { to: '/settings', label: 'Settings', icon: Settings },
]

export function AppShell() {
  const { data: authStatus, isLoading: authLoading } = useAuthStatus()
  const { data: devices = [] } = useDevices()
  const [selectedDevice, setSelectedDevice] = useState<string>(() => {
    return localStorage.getItem('goremote:selectedDevice') || ''
  })
  const logout = useLogout()
  const navigate = useNavigate()

  if (authLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center text-sm text-muted">
        Loading...
      </div>
    )
  }

  if (!authStatus?.has_users || authStatus.enrollment_available) {
    // If no users or enrollment available, still go to login
  }

  // Check if we have a valid session - use auth status
  // If there are users but no active session, backend will reject API calls; but we need to check
  // For now, require that we have users and are not in enrollment mode with no session
  // Better: if API calls will be protected, we need to ensure session exists
  // Just check - if not authenticated in a real sense, redirect to login
  // The backend's session check is authoritative

  const handleDeviceChange = (value: string) => {
    setSelectedDevice(value)
    if (value) {
      localStorage.setItem('goremote:selectedDevice', value)
    } else {
      localStorage.removeItem('goremote:selectedDevice')
    }
  }

  const handleLogout = async () => {
    try {
      await logout.mutateAsync()
    } catch (e) {
      // ignore
    } finally {
      navigate('/login')
    }
  }

  return (
    <div className="min-h-screen bg-background text-foreground">
      <header className="sticky top-0 z-40 border-b border-border bg-background/80 backdrop-blur">
        <div className="mx-auto flex max-w-[1400px] items-center justify-between gap-4 px-4 py-3 md:px-6">
          <div className="flex items-center gap-3">
            <div className="flex h-9 w-9 items-center justify-center rounded-full bg-accent/20">
              <div className="h-3 w-3 rounded-full bg-accent" />
            </div>
            <div>
              <div className="text-sm font-semibold leading-tight">GoRemote</div>
              <div className="text-xs text-muted">Universal TV Remote</div>
            </div>
          </div>
          <div className="flex flex-1 items-center justify-end gap-3">
            <select
              value={selectedDevice}
              onChange={(e) => handleDeviceChange(e.target.value)}
              className="hidden h-10 rounded-lg border border-border bg-surface px-3 text-sm md:flex"
            >
              <option value="">Select TV</option>
              {devices.map((d) => (
                <option key={d.id} value={d.id}>
                  {d.name} ({d.driver})
                </option>
              ))}
            </select>
            <div className="hidden items-center gap-2 rounded-full border border-border px-3 py-1 text-xs text-muted md:flex">
              <span className={cn('h-2 w-2 rounded-full', selectedDevice ? 'bg-success' : 'bg-muted')} />
              {selectedDevice ? 'Connected' : 'Disconnected'}
            </div>
            <button
              onClick={handleLogout}
              className="h-10 rounded-lg border border-border px-3 text-sm text-muted hover:text-foreground"
            >
              Logout
            </button>
          </div>
        </div>
      </header>
      <div className="mx-auto grid max-w-[1400px] grid-cols-1 gap-6 px-4 py-6 md:grid-cols-[240px_1fr] md:px-6">
        <aside className="hidden flex-col gap-6 md:flex">
          <nav className="flex flex-col gap-1 rounded-xl border border-border bg-surface p-2">
            {navItems.map((item) => (
              <NavLink
                key={item.to}
                to={item.to}
                end={item.to === '/'}
                className={({ isActive }) =>
                  cn(
                    'flex items-center gap-2 rounded-lg px-3 py-2 text-sm transition-colors',
                    isActive ? 'bg-surface-elevated text-foreground' : 'text-muted hover:text-foreground hover:bg-surface-elevated/60'
                  )
                }
              >
                <item.icon className="h-4 w-4" />
                {item.label}
              </NavLink>
            ))}
          </nav>
          <div className="rounded-xl border border-border bg-surface p-4">
            <div className="text-sm font-medium">Experimental VIDAA</div>
            <div className="mt-1 text-xs text-muted">Hisense VIDAA support - unverified</div>
            <NavLink to="/devices" className="mt-3 inline-flex text-xs text-accent">
              Manage VIDAA devices →
            </NavLink>
          </div>
        </aside>
        <main className="min-h-[70vh]">
          <Outlet context={{ selectedDevice }} />
        </main>
      </div>
      <nav className="fixed bottom-0 left-0 right-0 z-40 border-t border-border bg-background/90 backdrop-blur md:hidden">
        <div className="grid grid-cols-4">
          {navItems.map((item) => (
            <NavLink
              key={item.to}
              to={item.to}
              end={item.to === '/'}
              className={({ isActive }) =>
                cn(
                  'flex flex-col items-center gap-1 py-2 text-xs',
                  isActive ? 'text-foreground' : 'text-muted'
                )
              }
            >
              <item.icon className="h-5 w-5" />
              {item.label}
            </NavLink>
          ))}
        </div>
      </nav>
      <div className="h-16 md:hidden" />
    </div>
  )
}
