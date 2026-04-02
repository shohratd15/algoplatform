<template>
  <div class="auth-wrapper">
    <div class="auth-card glass-panel animate-fade-up">
      <h2>Join AlgoPlatform</h2>
      <p class="subtitle">Start your algorithmic mastery journey</p>
      
      <form @submit.prevent="handleRegister" class="auth-form">
        <div class="form-group">
          <label>Username</label>
          <input type="text" v-model="username" placeholder="coder_ninja" required minlength="3">
        </div>
        
        <div class="form-group">
          <label>Email</label>
          <input type="email" v-model="email" placeholder="ninja@example.com" required>
        </div>
        
        <div class="form-group">
          <label>Password</label>
          <input type="password" v-model="password" placeholder="Min 8 characters" required minlength="8">
        </div>
        
        <div class="error-msg" v-if="error">{{ error }}</div>
        
        <button type="submit" class="btn btn-primary full-width" :disabled="loading">
          {{ loading ? 'Creating Account...' : 'Sign Up' }}
        </button>
      </form>
      
      <div class="auth-footer">
        <p>Already have an account? <router-link to="/login">Sign In</router-link></p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

const router = useRouter()
const username = ref('')
const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

const handleRegister = async () => {
  error.value = ''
  loading.value = true
  
  try {
    await axios.post('/api/register', {
      username: username.value,
      email: email.value,
      password: password.value,
      role: 'user'
    })
    
    // Automatically redirect to login
    router.push('/login')
  } catch (err) {
    error.value = err.response?.data || 'Failed to register account.'
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

h2 {
  text-align: center;
  margin-bottom: 0.5rem;
}

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
