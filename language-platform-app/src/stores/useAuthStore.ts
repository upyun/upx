import { create } from 'zustand'
import { apiRequest } from '@/lib/api'
import { useSettingsStore } from '@/stores/useSettingsStore'
import type { Locale } from '@/i18n/strings'

export type AuthUser = {
  id: string
  email: string
  uiLocale: Locale
}

type AuthState = {
  user: AuthUser | null
  loading: boolean
  error: string | null
  register: (payload: { email: string; password: string; uiLocale: Locale }) => Promise<void>
  login: (payload: { email: string; password: string }) => Promise<void>
  logout: () => Promise<void>
  loadMe: () => Promise<void>
  clearError: () => void
}

export const useAuthStore = create<AuthState>((set, get) => ({
  user: null,
  loading: false,
  error: null,
  clearError: () => set({ error: null }),
  register: async ({ email, password, uiLocale }) => {
    set({ loading: true, error: null })
    try {
      const res = await apiRequest<{ success: true; user: AuthUser }>('/api/auth/register', {
        method: 'POST',
        json: { email, password, uiLocale },
      })
      set({ user: res.user, loading: false })
      useSettingsStore.getState().setUiLocale(res.user.uiLocale)
    } catch (e) {
      set({ loading: false, error: e instanceof Error ? e.message : 'Registration failed' })
    }
  },
  login: async ({ email, password }) => {
    set({ loading: true, error: null })
    try {
      const res = await apiRequest<{ success: true; user: AuthUser }>('/api/auth/login', {
        method: 'POST',
        json: { email, password },
      })
      set({ user: res.user, loading: false })
      useSettingsStore.getState().setUiLocale(res.user.uiLocale)
    } catch (e) {
      set({ loading: false, error: e instanceof Error ? e.message : 'Login failed' })
    }
  },
  logout: async () => {
    set({ loading: true, error: null })
    try {
      await apiRequest<{ success: true }>('/api/auth/logout', { method: 'POST' })
      set({ user: null, loading: false })
    } catch (e) {
      set({ loading: false, error: e instanceof Error ? e.message : 'Logout failed' })
    }
  },
  loadMe: async () => {
    if (get().loading) return
    set({ loading: true, error: null })
    try {
      const res = await apiRequest<{ success: true; user: AuthUser }>('/api/auth/me')
      set({ user: res.user, loading: false })
      useSettingsStore.getState().setUiLocale(res.user.uiLocale)
    } catch {
      set({ user: null, loading: false, error: null })
    }
  },
}))

