<template>
  <div class="auth-wrapper">
    <div class="auth-card glass-panel animate-fade-up">
      <!-- Logo -->
      <div class="auth-logo">
        <img src="../assets/logo.webp" alt="AlgoPlatform Logo" width="175" height="175" class="img-logo" />
        <span class="brand-name">AlgoPlatform</span>
      </div>

      <h2>{{ ui.t('loginTitle') }}</h2>
      <p class="subtitle">{{ ui.t('loginSubtitle') }}</p>

      <form @submit.prevent="handleLogin" class="auth-form">
        <div class="form-group">
          <label>{{ ui.t('loginEmail') }}</label>
          <input type="email" v-model="email" :placeholder="ui.t('loginEmailPlaceholder')" required />
        </div>
        <div class="form-group">
          <label>{{ ui.t('loginPassword') }}</label>
          <input type="password" v-model="password" :placeholder="ui.t('loginPasswordPlaceholder')" required />
        </div>

        <div class="error-msg" v-if="error">{{ error }}</div>

        <button type="submit" class="btn btn-primary full-width" :disabled="loading">
          {{ loading ? ui.t('loginLoading') : ui.t('loginSubmit') }}
        </button>
      </form>

      <div class="auth-footer">
        <p>{{ ui.t('loginFooter') }} <router-link to="/register">{{ ui.t('loginFooterLink') }}</router-link></p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import client from '../api/client'
import { useAuthStore } from '../stores/auth'
import { useUIStore } from '../stores/ui'

const router = useRouter()
const auth = useAuthStore()
const ui = useUIStore()

const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

const handleLogin = async () => {
  error.value = ''
  loading.value = true
  try {
    const res = await client.post('/login', { email: email.value, password: password.value })
    auth.setAuth(res.data.token, res.data.refresh_token, res.data.role)
    router.push('/problems')
  } catch (err) {
    error.value = err.response?.data || ui.t('loginFailed')
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
  padding: 2.5rem;
}

.auth-logo {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 1.5rem;
}

.brand-name {
  font-weight: 800;
  font-size: 1.35rem;
  letter-spacing: -0.03em;
  background: linear-gradient(90deg, #fff, #93c5fd);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

h2 {
  text-align: center;
  margin-bottom: 0.4rem;
  font-size: 1.5rem;
}

.subtitle {
  text-align: center;
  color: var(--text-muted);
  margin-bottom: 1.75rem;
  font-size: 0.9rem;
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 1.1rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.45rem;
}

label {
  font-size: 0.8rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--text-muted);
}

.full-width { width: 100%; margin-top: 0.75rem; }

.error-msg {
  color: #f87171;
  font-size: 0.875rem;
  text-align: center;
  background: rgba(248, 113, 113, 0.1);
  padding: 0.75rem;
  border-radius: 8px;
}

.auth-footer {
  margin-top: 1.75rem;
  text-align: center;
  font-size: 0.875rem;
  color: var(--text-muted);
}
</style>
