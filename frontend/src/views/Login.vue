<template>
  <div class="auth-wrapper">
    <div class="auth-card glass-panel animate-fade-up">
      <h2>Welcome Back</h2>
      <p class="subtitle">Enter your credentials to continue</p>

      <form @submit.prevent="handleLogin" class="auth-form">
        <div class="form-group">
          <label>Email</label>
          <input type="email" v-model="email" placeholder="coder@example.com" required />
        </div>

        <div class="form-group">
          <label>Password</label>
          <input type="password" v-model="password" placeholder="••••••••" required />
        </div>

        <div class="error-msg" v-if="error">{{ error }}</div>

        <button type="submit" class="btn btn-primary full-width" :disabled="loading">
          {{ loading ? 'Authenticating...' : 'Sign In' }}
        </button>
      </form>

      <div class="auth-footer">
        <p>New to AlgoPlatform? <router-link to="/register">Create an account</router-link></p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import client from '../api/client'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const auth = useAuthStore()

const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

const handleLogin = async () => {
  error.value = ''
  loading.value = true

  try {
    const res = await client.post('/login', {
      email: email.value,
      password: password.value,
    })

    auth.setToken(res.data.token)
    router.push('/problems')
  } catch (err) {
    error.value = err.response?.data || 'Login failed. Please check credentials.'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: calc(100vh - 70px);
  padding: 2rem;
}

.auth-card {
  width: 100%;
  max-width: 420px;
  padding: 3rem 2.5rem;
}

h2 { text-align: center; margin-bottom: 0.5rem; }

.subtitle {
  text-align: center;
  color: var(--text-muted);
  margin-bottom: 2rem;
  font-size: 0.95rem;
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

label {
  font-size: 0.85rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--text-muted);
}

.full-width {
  width: 100%;
  margin-top: 1rem;
}

.error-msg {
  color: #ff5f56;
  font-size: 0.9rem;
  text-align: center;
  background: rgba(255, 95, 86, 0.1);
  padding: 0.75rem;
  border-radius: 8px;
}

.auth-footer {
  margin-top: 2rem;
  text-align: center;
  font-size: 0.9rem;
  color: var(--text-muted);
}
</style>
