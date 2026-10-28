import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { MemoryRouter, Outlet, Route, Routes } from 'react-router-dom'
import { RemotePage } from './RemotePage'

describe('RemotePage', () => {
  it('shows first-run when no devices', () => {
    const qc = new QueryClient()
    render(
      <QueryClientProvider client={qc}>
        <MemoryRouter initialEntries={['/']}>
          <Routes>
            <Route element={<Outlet context={{ selectedDevice: '' }} />}>
              <Route path="/" element={<RemotePage />} />
            </Route>
          </Routes>
        </MemoryRouter>
      </QueryClientProvider>
    )
    expect(screen.getByText('Connect your first TV')).toBeInTheDocument()
  })
})
