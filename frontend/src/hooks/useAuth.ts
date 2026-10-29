import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { api } from '@/lib/api/client'
import type { AuthStatus } from '@/types'

export function useAuthStatus() {
  return useQuery({
    queryKey: ['auth', 'status'],
    queryFn: () => api.get<AuthStatus>('/api/auth/status'),
  })
}

export function useEnroll() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: { username: string; password: string; bootstrap_token: string }) =>
      api.post('/api/auth/enroll', data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['auth', 'status'] })
    },
  })
}

export function useLogin() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: { username: string; password: string }) => api.post('/api/auth/login', data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['auth', 'status'] })
    },
  })
}

export function useLogout() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: () => api.post('/api/auth/logout'),
    onSuccess: () => {
      qc.clear()
    },
    onError: () => {
      qc.clear()
    },
  })
}
