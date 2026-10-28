export type Device = {
  id: string
  name: string
  driver: string
  state: string
  address?: string
}

export type DiscoveredDevice = {
  id: string
  name: string
  address: string
  driver_hint?: string
  location?: string
  usn?: string
}

export type App = {
  id: string
  name: string
}

export type AuthStatus = {
  has_users: boolean
  enrollment_available: boolean
}
