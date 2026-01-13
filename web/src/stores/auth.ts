import { defineStore } from 'pinia'
import { setAccessToken } from '@/api/client'

export interface AuthUser {
  id: string
  name: string
}

interface AuthState {
  user: AuthUser | null
  accessToken: string | null
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    user: null,
    accessToken: null
  }),
  getters: {
    isAuthenticated: (state) => !!state.accessToken
  },
  actions: {
    setToken(token: string | null) {
      this.accessToken = token
      setAccessToken(token)
    },
    logout() {
      this.user = null
      this.setToken(null)
    }
  }
})

