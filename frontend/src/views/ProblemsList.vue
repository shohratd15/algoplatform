<template>
  <div class="problems-wrapper">
    <div class="header-section animate-fade-up">
      <h1>Algorithm Challenges</h1>
      <p class="subtitle">Select a problem to solve and improve your algorithmic skills.</p>
    </div>

    <div class="problem-list animate-fade-up" style="animation-delay: 0.1s">
      <div v-if="loading" class="loading">Loading challenges...</div>

      <div v-else-if="error" class="error-state glass-panel">
        <p>{{ error }}</p>
      </div>

      <div v-else-if="problems.length === 0" class="empty-state glass-panel">
        <p>No problems available yet.</p>
      </div>

      <div v-else class="grid">
        <router-link
          v-for="p in problems"
          :key="p.id"
          :to="'/problems/' + p.id"
          class="problem-card glass-panel"
        >
          <div class="card-top">
            <span class="diff-badge" :class="p.difficulty.toLowerCase()">
              {{ p.difficulty }}
            </span>
            <span class="id-badge">#{{ p.id }}</span>
          </div>
          <h3>{{ p.slug.replace(/-/g, ' ') }}</h3>
          <p class="date">Added on {{ new Date(p.created_at).toLocaleDateString() }}</p>

          <div class="card-action">
            <span>Solve Challenge</span>
            <span class="arrow">→</span>
          </div>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import client from '../api/client'

const problems = ref([])
const loading = ref(true)
const error = ref('')

const fetchProblems = async () => {
  loading.value = true
  error.value = ''

  try {
    const res = await client.get('/problems')
    problems.value = res.data || []
  } catch (err) {
    // 401 обрабатывается interceptor-ом — редирект на /login
    if (err.response?.status !== 401) {
      error.value = 'Failed to load problems. Please try again.'
    }
  } finally {
    loading.value = false
  }
}

onMounted(fetchProblems)
</script>

<style scoped>
.problems-wrapper {
  max-width: 1200px;
  margin: 0 auto;
  padding: 3rem 2rem;
  width: 100%;
}

.header-section {
  text-align: center;
  margin-bottom: 3rem;
}

.header-section h1 { margin-bottom: 0.5rem; }

.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 1.5rem;
}

.problem-card {
  display: flex;
  flex-direction: column;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  cursor: pointer;
}

.problem-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 12px 40px rgba(0, 255, 136, 0.1);
  border-color: rgba(0, 255, 136, 0.3);
}

.card-top {
  display: flex;
  justify-content: space-between;
  margin-bottom: 1rem;
}

.diff-badge {
  padding: 0.25rem 0.75rem;
  border-radius: 20px;
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: uppercase;
}

.diff-badge.easy   { background: rgba(39, 201, 63, 0.15);  color: #27c93f; }
.diff-badge.medium { background: rgba(255, 189, 46, 0.15); color: #ffbd2e; }
.diff-badge.hard   { background: rgba(255, 95, 86, 0.15);  color: #ff5f56; }

.id-badge {
  color: var(--text-muted);
  font-family: monospace;
  font-size: 0.9rem;
}

.problem-card h3 {
  font-size: 1.25rem;
  margin-bottom: 0.5rem;
  text-transform: capitalize;
}

.date {
  font-size: 0.85rem;
  color: var(--text-muted);
  margin-bottom: 2rem;
}

.card-action {
  margin-top: auto;
  display: flex;
  justify-content: space-between;
  color: var(--accent);
  font-weight: 600;
  font-size: 0.9rem;
  padding-top: 1rem;
  border-top: 1px solid var(--glass-border);
}

.arrow { transition: transform 0.2s; }
.problem-card:hover .arrow { transform: translateX(4px); }

.loading, .empty-state, .error-state {
  text-align: center;
  padding: 4rem;
  color: var(--text-muted);
}

.error-state { color: #ff5f56; }
</style>
