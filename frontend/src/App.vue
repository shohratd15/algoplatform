<template>
  <div id="app">
    <nav class="navbar">
      <div class="nav-container">
        <router-link to="/" class="logo">
          <span class="logo-icon">▲</span>
          <span class="logo-text">AlgoPlatform</span>
        </router-link>
        
        <div class="nav-links">
          <router-link to="/problems" class="nav-link">Problems</router-link>
          <router-link to="/login" class="nav-link" v-if="!isAuthenticated">Login</router-link>
          <router-link to="/register" class="btn btn-primary btn-sm" v-if="!isAuthenticated">Sign Up</router-link>
          <a href="#" class="nav-link" v-if="isAuthenticated" @click.prevent="logout">Logout</a>
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
import { ref, watchEffect } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const isAuthenticated = ref(!!localStorage.getItem('token'))

watchEffect(() => {
  isAuthenticated.value = !!localStorage.getItem('token')
})

const logout = () => {
  localStorage.removeItem('token')
  isAuthenticated.value = false
  router.push('/login')
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
  text-shadow: 0 0 10px var(--accent-glow);
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

.nav-link:hover {
  color: #fff;
}

.btn-sm {
  padding: 0.5rem 1rem;
  font-size: 0.9rem;
}

.main-content {
  flex: 1;
  margin-top: 70px; /* offset navbar */
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

/* Page transitions */
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
