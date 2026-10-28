import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { api } from '@/lib/api/client'
import type { DiscoveredDevice } from '@/types'

export function useDiscover() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: () => api.post<{ devices: DiscoveredDevice[] }>('/api/discover/start'),
    onSuccess: (data) => {
      qc.setQueryData(['discover', 'results'], data.devices || [])
    },
  })
}

export function useDiscoverResults() {
  return useQuery({
    queryKey: ['discover', 'results'],
    queryFn: () => [] as DiscoveredDevice[],
    initialData: [] as DiscoveredDevice[],
  })
}
