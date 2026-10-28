import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { api } from '@/lib/api/client'
import type { Device } from '@/types'

export function useDevices() {
  return useQuery({
    queryKey: ['devices'],
    queryFn: () => api.get<Device[]>('/api/devices'),
  })
}

export function useSaveDiscovered() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: { id: string; address: string; name: string; driver?: string }) =>
      api.post('/api/discover/save', data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['devices'] })
      qc.invalidateQueries({ queryKey: ['discover'] })
    },
  })
}

export function useKeypress() {
  return useMutation({
    mutationFn: ({ key, device }: { key: string; device?: string | null }) =>
      api.post(`/api/keypress/${encodeURIComponent(key)}${device ? `?device=${encodeURIComponent(device)}` : ''}`),
  })
}

export function useApps(device?: string | null) {
  return useQuery({
    queryKey: ['apps', device],
    queryFn: () => api.get<any[]>(`/api/apps${device ? `?device=${encodeURIComponent(device)}` : ''}`),
    enabled: !!device,
  })
}

export function useLaunch() {
  return useMutation({
    mutationFn: ({ appId, device }: { appId: string; device?: string | null }) =>
      api.post(`/api/launch/${encodeURIComponent(appId)}${device ? `?device=${encodeURIComponent(device)}` : ''}`),
  })
}
