import { create } from 'zustand'
import { persist } from 'zustand/middleware'

interface User {
  id: string
  email: string
  name?: string
}

interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  login: (email: string, password: string) => Promise<void>
  setAuth: (user: User, token: string) => void
  logout: () => void
}

const API_BASE_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:3002'

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      login: async (email: string, password: string): Promise<void> => {
        const response = await fetch(`${API_BASE_URL}/api/auth/login`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ email, password })
        })

        if (!response.ok) {
          const error = await response.json() as { error: string }
          throw new Error(error.error ?? 'Login failed')
        }

        const data = await response.json() as { user: User, token: string }
        set({
          user: data.user,
          token: data.token,
          isAuthenticated: true
        })
      },
      setAuth: (user: User, token: string): void =>
        set({
          user,
          token,
          isAuthenticated: true
        }),
      logout: (): void =>
        set({
          user: null,
          token: null,
          isAuthenticated: false
        })
    }),
    {
      name: 'auth-storage'
    }
  )
)
