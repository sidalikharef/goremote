import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { MemoryRouter } from 'react-router-dom'
import { AuthPage } from './AuthPage'

describe('AuthPage', () => {
  it('renders login title', () => {
    const qc = new QueryClient()
    render(
      <QueryClientProvider client={qc}>
        <MemoryRouter>
          <AuthPage />
        </MemoryRouter>
      </QueryClientProvider>
    )
    expect(screen.getByText('GoRemote')).toBeInTheDocument()
  })
})
