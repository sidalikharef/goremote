import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { MemoryRouter } from 'react-router-dom'
import { DiscoveryPage } from './DiscoveryPage'

describe('DiscoveryPage', () => {
  it('renders scan button', () => {
    const qc = new QueryClient()
    render(
      <QueryClientProvider client={qc}>
        <MemoryRouter>
          <DiscoveryPage />
        </MemoryRouter>
      </QueryClientProvider>
    )
    expect(screen.getByText('Scan for TVs')).toBeInTheDocument()
  })
})
