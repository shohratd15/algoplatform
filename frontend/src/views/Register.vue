<template>
  <div class="auth-wrapper">
    <div class="auth-card glass-panel animate-fade-up">
      <div class="auth-logo">
        <svg width="52" height="52" viewBox="0 0 120 120" fill="none" xmlns="http://www.w3.org/2000/svg">
          <defs>
            <linearGradient id="lg-reg" x1="0%" y1="0%" x2="100%" y2="100%">
              <stop offset="0%" stop-color="#1d4ed8"/>
              <stop offset="100%" stop-color="#0ea5e9"/>
            </linearGradient>
          </defs>
          <line x1="60" y1="8" x2="10" y2="100" stroke="url(#lg-reg)" stroke-width="9" stroke-linecap="round"/>
          <line x1="60" y1="8" x2="110" y2="100" stroke="url(#lg-reg)" stroke-width="9" stroke-linecap="round"/>
          <line x1="10" y1="100" x2="110" y2="100" stroke="url(#lg-reg)" stroke-width="9" stroke-linecap="round"/>
          <ellipse cx="60" cy="72" rx="40" ry="16" stroke="url(#lg-reg)" stroke-width="7" fill="none" stroke-linecap="round"/>
          <circle cx="60" cy="8" r="8" fill="url(#lg-reg)"/>
          <circle cx="10" cy="100" r="8" fill="url(#lg-reg)"/>
          <circle cx="110" cy="100" r="8" fill="url(#lg-reg)"/>
        </svg>
        <span class="brand-name">AlgoPlatform</span>
      </div>

      <h2>{{ ui.t('registerTitle') }}</h2>
      <p class="subtitle">{{ ui.t('registerSubtitle') }}</p>

      <form @submit.prevent="handleRegister" class="auth-form">
        <div class="form-group">
          <label>{{ ui.t('registerUsername') }}</label>
          <input type="text" v-model="username" :placeholder="ui.t('registerUsernamePlaceholder')" required minlength="3" />
        </div>
        <div class="form-group">
          <label>{{ ui.t('registerEmail') }}</label>
          <input type="email" v-model="email" :placeholder="ui.t('registerEmailPlaceholder')" required />
        </div>
        <div class="form-group">
          <label>{{ ui.t('registerPassword') }}</label>
          <input type="password" v-model="password" :placeholder="ui.t('registerPasswordPlaceholder')" required minlength="8" />
        </div>

        <div class="error-msg" v-if="error">{{ error }}</div>

        <button type="submit" class="btn btn-primary full-width" :disabled="loading">
          {{ loading ? ui.t('registerLoading') : ui.t('registerSubmit') }}
        </button>
      </form>

      <div class="auth-footer">
        <p>{{ ui.t('registerFooter') }} <router-link to="/login">{{ ui.t('registerFooterLink') }}</router-link></p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import client from '../api/client'
import { useUIStore } from '../stores/ui'

const router = useRouter()
const ui = useUIStore()

const username = ref('')
const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

const handleRegister = async () => {
  error.value = ''
  loading.value = true
  try {
    await client.post('/register', { username: username.value, email: email.value, password: password.value })
    router.push('/login')
  } catch (err) {
    error.value = err.response?.data || ui.t('registerFailed')
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
.auth-card { width: 100%; max-width: 420px; padding: 2.5rem; }
.auth-logo { display: flex; flex-direction: column; align-items: center; gap: 0.5rem; margin-bottom: 1.5rem; }
.brand-name {
  font-weight: 800; font-size: 1.35rem; letter-spacing: -0.03em;
  background: linear-gradient(90deg, #fff, #93c5fd);
  -webkit-background-clip: text; background-clip: text; -webkit-text-fill-color: transparent;
}
h2 { text-align: center; margin-bottom: 0.4rem; font-size: 1.5rem; }
.subtitle { text-align: center; color: var(--text-muted); margin-bottom: 1.75rem; font-size: 0.9rem; }
.auth-form { display: flex; flex-direction: column; gap: 1.1rem; }
.form-group { display: flex; flex-direction: column; gap: 0.45rem; }
label { font-size: 0.8rem; font-weight: 600; text-transform: uppercase; letter-spacing: 0.06em; color: var(--text-muted); }
.full-width { width: 100%; margin-top: 0.75rem; }
.error-msg { color: #f87171; font-size: 0.875rem; text-align: center; background: rgba(248,113,113,0.1); padding: 0.75rem; border-radius: 8px; }
.auth-footer { margin-top: 1.75rem; text-align: center; font-size: 0.875rem; color: var(--text-muted); }
</style>
