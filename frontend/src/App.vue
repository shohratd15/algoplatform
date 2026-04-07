<template>
  <div id="app">
    <nav class="navbar">
      <div class="nav-container">
        <router-link to="/" class="logo">
          <span class="logo-icon">▲</span>
          <span class="logo-text">AlgoPlatform</span>
        </router-link>

        <div class="nav-links">
          <router-link to="/problems" class="nav-link">{{ ui.t('navProblems') }}</router-link>
          <select class="lang-select" :value="ui.state.locale" @change="onLocaleChange">
            <option value="en">EN</option>
            <option value="ru">RU</option>
            <option value="tm">TM</option>
          </select>
          <template v-if="!auth.state.token">
            <router-link to="/login" class="nav-link">{{ ui.t('navLogin') }}</router-link>
            <router-link to="/register" class="btn btn-primary btn-sm">{{ ui.t('navSignup') }}</router-link>
          </template>
          <template v-else>
            <router-link v-if="auth.state.role === 'admin'" to="/admin/problems" class="nav-link">{{ ui.t('navAdmin') }}</router-link>
            <a href="#" class="nav-link" @click.prevent="logout">{{ ui.t('navLogout') }}</a>
          </template>
        </div>
      </div>
    </nav>

    <main class="main-content">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>

    <footer class="footer">
      <p>&copy; 2026 AlgoPlatform — Learn Algorithms efficiently.</p>
    </footer>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useAuthStore } from './stores/auth'
import { useUIStore } from './stores/ui'

const router = useRouter()
const auth = useAuthStore()
const ui = useUIStore()

const logout = () => {
  auth.clearToken()
  router.push('/login')
}

const onLocaleChange = (e) => {
  ui.setLocale(e.target.value)
}
</script>

<style scoped>
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 70px;
  background: rgba(15, 17, 21, 0.8);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border-bottom: 1px solid var(--glass-border);
  z-index: 100;
  display: flex;
  align-items: center;
}

.nav-container {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.logo {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.logo-icon {
  color: var(--accent);
  font-size: 1.5rem;
}

.logo-text {
  font-weight: 800;
  font-size: 1.25rem;
  letter-spacing: -0.02em;
  color: #fff;
}

.nav-links {
  display: flex;
  align-items: center;
  gap: 2rem;
}

.nav-link {
  color: var(--text-muted);
  font-weight: 500;
  font-size: 0.95rem;
}

.nav-link:hover { color: #fff; }

.btn-sm {
  padding: 0.5rem 1rem;
  font-size: 0.9rem;
}

.lang-select {
  background: rgba(0, 0, 0, 0.25);
  color: #fff;
  border: 1px solid var(--glass-border);
  border-radius: 6px;
  padding: 0.35rem 0.5rem;
}

.main-content {
  flex: 1;
  margin-top: 70px;
  display: flex;
  flex-direction: column;
}

.footer {
  text-align: center;
  padding: 2rem;
  color: var(--text-muted);
  font-size: 0.9rem;
  border-top: 1px solid var(--glass-border);
  margin-top: auto;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.fade-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
