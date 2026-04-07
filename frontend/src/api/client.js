// src/api/client.js
import axios from 'axios'
import { useAuthStore } from '../stores/auth'

const client = axios.create({
  baseURL: '/api',
})

let isRefreshing = false
let refreshQueue = []

const processQueue = (token) => {
  refreshQueue.forEach((resolve) => resolve(token))
  refreshQueue = []
}

client.interceptors.request.use((config) => {
  const { state } = useAuthStore()
  if (state.token) {
    config.headers.Authorization = `Bearer ${state.token}`
  }
  return config
})

client.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config
    const { state, clearToken, setAuth } = useAuthStore()

    if (error.response?.status === 401 && !originalRequest?._retry && state.refreshToken) {
      originalRequest._retry = true

      if (isRefreshing) {
        return new Promise((resolve) => {
          refreshQueue.push((newToken) => {
            originalRequest.headers.Authorization = `Bearer ${newToken}`
            resolve(client(originalRequest))
          })
        })
      }

      isRefreshing = true
      try {
        const res = await axios.post('/api/refresh', {
          refresh_token: state.refreshToken,
        })
        setAuth(res.data.token, res.data.refresh_token, res.data.role)
        processQueue(res.data.token)
        originalRequest.headers.Authorization = `Bearer ${res.data.token}`
        return client(originalRequest)
      } catch (refreshErr) {
        clearToken()
        import('../router').then(({ default: router }) => {
          router.push('/login')
        })
        return Promise.reject(refreshErr)
      } finally {
        isRefreshing = false
      }
    }

    if (error.response?.status === 401) {
      clearToken()
      import('../router').then(({ default: router }) => {
        router.push('/login')
      })
    }
    return Promise.reject(error)
  }
)

export default client
