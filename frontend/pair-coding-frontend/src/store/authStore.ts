import { create } from 'zustand'
import { persist } from 'zustand/middleware'

type AuthState = {
  isLoggedIn: boolean
  isAdminLoggedIn: boolean
  login: () => void
  logout: () => void
  adminLogin: () => void
  adminLogout: () => void
  markUnauthorized: () => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      isLoggedIn: true,
      isAdminLoggedIn: false,
      login: () => set({ isLoggedIn: true }),
      logout: () => set({ isLoggedIn: false }),
      adminLogin: () => set({ isAdminLoggedIn: true }),
      adminLogout: () => set({ isAdminLoggedIn: false }),
      markUnauthorized: () => set({ isLoggedIn: false }),
    }),
    {
      name: 'fodecorces-auth',
    },
  ),
)

export function handleUnauthorized() {
  useAuthStore.getState().markUnauthorized()
}
