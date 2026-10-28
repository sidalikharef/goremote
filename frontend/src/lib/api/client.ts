const defaultHeaders: HeadersInit = {
  'Content-Type': 'application/json',
}

async function request<T>(input: RequestInfo | URL, init: RequestInit = {}): Promise<T> {
  const res = await fetch(input, {
    credentials: 'include',
    headers: {
      ...defaultHeaders,
      ...(init.headers || {}),
    },
    ...init,
  })

  if (res.status === 401) {
    throw new Error('unauthorized')
  }
  if (!res.ok) {
    throw new Error(`request_failed:${res.status}`)
  }
  const contentType = res.headers.get('content-type') || ''
  if (contentType.includes('application/json')) {
    return (await res.json()) as T
  }
  return undefined as T
}

export const api = {
  get: <T>(url: string) => request<T>(url, { method: 'GET' }),
  post: <T>(url: string, body?: unknown) =>
    request<T>(url, {
      method: 'POST',
      body: body !== undefined ? JSON.stringify(body) : undefined,
    }),
  del: <T>(url: string) => request<T>(url, { method: 'DELETE' }),
}
