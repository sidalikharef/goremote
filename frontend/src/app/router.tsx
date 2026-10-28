import { createBrowserRouter, Navigate } from 'react-router-dom'
import { AppShell } from '@/components/layout/AppShell'
import { RemotePage } from './pages/RemotePage'
import { DevicesPage } from './pages/DevicesPage'
import { DiscoveryPage } from './pages/DiscoveryPage'
import { SettingsPage } from './pages/SettingsPage'
import { AuthPage } from './pages/AuthPage'

export const router = createBrowserRouter([
  {
    path: '/login',
    element: <AuthPage />,
  },
  {
    path: '/',
    element: <AppShell />,
    children: [
      { index: true, element: <RemotePage /> },
      { path: 'devices', element: <DevicesPage /> },
      { path: 'discovery', element: <DiscoveryPage /> },
      { path: 'settings', element: <SettingsPage /> },
    ],
  },
  { path: '*', element: <Navigate to="/" replace /> },
])
