// src/stores/auth.js
// Простой реактивный store для auth-состояния.
// Используется вместо прямого обращения к localStorage в компонентах.
import { reactive, readonly } from 'vue'

const state = reactive({
  token: localStorage.getItem('token') || null,
  refreshToken: localStorage.getItem('refresh_token') || null,
  role: localStorage.getItem('role') || null,
})

export function useAuthStore() {
  const isAuthenticated = () => !!state.token

  const setAuth = (token, refreshToken, role) => {
    state.token = token
    state.refreshToken = refreshToken || null
    state.role = role || null
    localStorage.setItem('token', token)
    if (refreshToken) localStorage.setItem('refresh_token', refreshToken)
    else localStorage.removeItem('refresh_token')
    if (role) localStorage.setItem('role', role)
    else localStorage.removeItem('role')
  }

  const clearToken = () => {
    state.token = null
    state.refreshToken = null
    state.role = null
    localStorage.removeItem('token')
    localStorage.removeItem('refresh_token')
    localStorage.removeItem('role')
  }

  return {
    state: readonly(state),
    isAuthenticated,
    setAuth,
    clearToken,
  }
}
