// src/api/client.js
import axios from 'axios'
import { useAuthStore } from '../stores/auth'

const client = axios.create({
  baseURL: '/api',
})

client.interceptors.request.use((config) => {
  const { state } = useAuthStore()
  if (state.token) {
    config.headers.Authorization = `Bearer ${state.token}`
  }
  return config
})

client.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      const { clearToken } = useAuthStore()
      clearToken()
      import('../router').then(({ default: router }) => {
        router.push('/login')
      })
    }
    return Promise.reject(error)
  }
)

export default client
