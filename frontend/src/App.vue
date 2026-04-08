<template>
  <div id="app">
    <nav class="navbar">
      <div class="nav-container">
        <router-link to="/" class="logo">
          <!-- AlgoPlatform logo mark -->
          <img src="./assets/logo.webp" alt="AlgoPlatform Logo" width="50" height="50" class="img-logo" />
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
      <p>{{ ui.t('footerText') }}</p>
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
  background: rgba(10, 14, 26, 0.85);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
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
  gap: 0.65rem;
  text-decoration: none;
}

.logo-text {
  font-weight: 800;
  font-size: 1.2rem;
  letter-spacing: -0.03em;
  background: linear-gradient(90deg, #fff, #93c5fd);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.nav-links {
  display: flex;
  align-items: center;
  gap: 1.75rem;
}

.nav-link {
  color: var(--text-muted);
  font-weight: 500;
  font-size: 0.95rem;
  transition: color 0.2s;
}

.nav-link:hover { color: #fff; text-shadow: none; }

.btn-sm {
  padding: 0.45rem 1rem;
  font-size: 0.88rem;
}

.lang-select {
  background: rgba(0, 0, 0, 0.3);
  color: #fff;
  border: 1px solid var(--glass-border);
  border-radius: 6px;
  padding: 0.35rem 0.6rem;
  cursor: pointer;
  outline: none;
  transition: border-color 0.2s;
}

.lang-select:hover {
  border-color: rgba(59, 130, 246, 0.4);
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
  font-size: 0.875rem;
  border-top: 1px solid var(--glass-border);
  margin-top: auto;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.25s ease, transform 0.25s ease;
}
.fade-enter-from { opacity: 0; transform: translateY(8px); }
.fade-leave-to   { opacity: 0; transform: translateY(-8px); }
</style>
