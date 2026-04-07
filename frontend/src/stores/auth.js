// src/stores/auth.js
// Простой реактивный store для auth-состояния.
// Используется вместо прямого обращения к localStorage в компонентах.
import { reactive, readonly } from 'vue'

const state = reactive({
  token: localStorage.getItem('token') || null,
})

export function useAuthStore() {
  const isAuthenticated = () => !!state.token

  const setToken = (token) => {
    state.token = token
    localStorage.setItem('token', token)
  }

  const clearToken = () => {
    state.token = null
    localStorage.removeItem('token')
  }

  return {
    state: readonly(state),
    isAuthenticated,
    setToken,
    clearToken,
  }
}
