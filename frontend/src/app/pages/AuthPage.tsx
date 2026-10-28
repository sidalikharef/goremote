import { useState } from 'react'
import { useAuthStatus, useEnroll, useLogin } from '@/hooks/useAuth'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { useNavigate } from 'react-router-dom'

export function AuthPage() {
  const { data: status } = useAuthStatus()
  const enroll = useEnroll()
  const login = useLogin()
  const navigate = useNavigate()

  const [enrollForm, setEnrollForm] = useState({ username: '', password: '', bootstrap_token: '' })
  const [loginForm, setLoginForm] = useState({ username: '', password: '' })

  const handleEnroll = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await enroll.mutateAsync(enrollForm)
      await login.mutateAsync({ username: enrollForm.username, password: enrollForm.password })
      navigate('/')
    } catch (err) {
      // ignore
    }
  }

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await login.mutateAsync(loginForm)
      navigate('/')
    } catch (err) {
      // ignore
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center px-4">
      <div className="w-full max-w-md space-y-6">
        <div className="text-center">
          <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-accent/20">
            <div className="h-4 w-4 rounded-full bg-accent" />
          </div>
          <h1 className="text-2xl font-semibold">GoRemote</h1>
          <p className="text-sm text-muted">Sign in to control your TVs</p>
        </div>
        {status?.enrollment_available && (
          <Card>
            <CardHeader>
              <CardTitle>First-time setup</CardTitle>
              <CardDescription>Create the first admin account (loopback only)</CardDescription>
            </CardHeader>
            <CardContent>
              <form onSubmit={handleEnroll} className="space-y-3">
                <Input placeholder="Username" value={enrollForm.username} onChange={(e) => setEnrollForm({ ...enrollForm, username: e.target.value })} />
                <Input type="password" placeholder="Password (min 8)" value={enrollForm.password} onChange={(e) => setEnrollForm({ ...enrollForm, password: e.target.value })} />
                <Input placeholder="Bootstrap token" value={enrollForm.bootstrap_token} onChange={(e) => setEnrollForm({ ...enrollForm, bootstrap_token: e.target.value })} />
                <Button type="submit" className="w-full" disabled={enroll.isPending}>
                  {enroll.isPending ? 'Creating...' : 'Create account'}
                </Button>
              </form>
            </CardContent>
          </Card>
        )}
        {status?.has_users && !status?.enrollment_available && (
          <Card>
            <CardHeader>
              <CardTitle>Sign in</CardTitle>
              <CardDescription>Use your GoRemote credentials</CardDescription>
            </CardHeader>
            <CardContent>
              <form onSubmit={handleLogin} className="space-y-3">
                <Input placeholder="Username" value={loginForm.username} onChange={(e) => setLoginForm({ ...loginForm, username: e.target.value })} />
                <Input type="password" placeholder="Password" value={loginForm.password} onChange={(e) => setLoginForm({ ...loginForm, password: e.target.value })} />
                <Button type="submit" className="w-full" disabled={login.isPending}>
                  {login.isPending ? 'Signing in...' : 'Sign in'}
                </Button>
              </form>
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  )
}
