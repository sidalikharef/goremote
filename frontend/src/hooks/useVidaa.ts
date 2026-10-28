import { useMutation } from '@tanstack/react-query'
import { api } from '@/lib/api/client'

export function useVidaaRegister() {
  return useMutation({
    mutationFn: (data: { id: string; name: string; address: string; enable: boolean }) =>
      api.post('/api/devices/vidaa/register', data),
  })
}

export function useVidaaPairStart() {
  return useMutation({
    mutationFn: (data?: unknown) => api.post('/api/devices/vidaa/pair/start', data),
  })
}

export function useVidaaPairSubmit() {
  return useMutation({
    mutationFn: (data: { pin: string }) => api.post('/api/devices/vidaa/pair/submit', data),
  })
}

export function useVidaaPairStatus() {
  return useMutation({
    mutationFn: () => api.get('/api/devices/vidaa/pair/status'),
  })
}

export function useVidaaPairCancel() {
  return useMutation({
    mutationFn: () => api.post('/api/devices/vidaa/pair/cancel'),
  })
}

export function useVidaaUnpair() {
  return useMutation({
    mutationFn: () => api.post('/api/devices/vidaa/unpair'),
  })
}
