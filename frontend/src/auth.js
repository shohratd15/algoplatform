// src/auth.js
// Общее реактивное состояние авторизации.
// Используем один shared ref — все компоненты импортируют его напрямую,
// поэтому изменения мгновенно отражаются в navbar и router guard.
// Pinia не нужна: один глобальный ref с сохранением в localStorage достаточен.
import { ref, computed } from 'vue'

export const authToken = ref(localStorage.getItem('token'))

export const isAuthenticated = computed(() => !!authToken.value)

export function setToken(token) {
  authToken.value = token
  if (token) {
    localStorage.setItem('token', token)
  } else {
    localStorage.removeItem('token')
  }
}

export function clearToken() {
  setToken(null)
}
