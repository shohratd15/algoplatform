<template>
  <div class="admin-wrapper">
    <div class="header-section animate-fade-up">
      <h1>Admin Problem Management</h1>
      <p class="subtitle">Create, update, and delete algorithm problems.</p>
    </div>

    <div class="status" v-if="error">{{ error }}</div>
    <div class="status success" v-if="success">{{ success }}</div>

    <div class="layout">
      <section class="glass-panel panel">
        <h2>Create Problem</h2>
        <form @submit.prevent="createProblem" class="form">
          <input v-model="createForm.slug" placeholder="slug (e.g. two-sum)" required />
          <select v-model="createForm.difficulty" required>
            <option value="easy">easy</option>
            <option value="medium">medium</option>
            <option value="hard">hard</option>
          </select>
          <textarea
            v-model="createForm.statementsJson"
            rows="7"
            placeholder='statements JSON: [{"language":"en","title":"Two Sum","statement":"..."}]'
            required
          />
          <textarea
            v-model="createForm.testsJson"
            rows="7"
            placeholder='tests JSON: [{"input_data":"1 2","expected_output":"3","is_sample":true}]'
            required
          />
          <button class="btn btn-primary" :disabled="loading">{{ loading ? 'Saving...' : 'Create' }}</button>
        </form>
      </section>

      <section class="glass-panel panel">
        <h2>Update Problem</h2>
        <form @submit.prevent="updateProblem" class="form">
          <input v-model.number="updateForm.id" type="number" min="1" placeholder="problem id" required />
          <input v-model="updateForm.slug" placeholder="new slug" required />
          <select v-model="updateForm.difficulty" required>
            <option value="easy">easy</option>
            <option value="medium">medium</option>
            <option value="hard">hard</option>
          </select>
          <textarea v-model="updateForm.statementsJson" rows="7" required />
          <textarea v-model="updateForm.testsJson" rows="7" required />
          <button class="btn btn-primary" :disabled="loading">{{ loading ? 'Updating...' : 'Update' }}</button>
        </form>
      </section>
    </div>

    <section class="glass-panel panel">
      <h2>Delete Problem</h2>
      <form @submit.prevent="deleteProblem" class="inline-form">
        <input v-model.number="deleteId" type="number" min="1" placeholder="problem id" required />
        <button class="btn btn-outline" :disabled="loading">{{ loading ? 'Deleting...' : 'Delete' }}</button>
      </form>
    </section>

    <section class="glass-panel panel">
      <div class="list-head">
        <h2>Existing Problems</h2>
        <button class="btn btn-outline" @click="fetchProblems" :disabled="loading">Refresh</button>
      </div>
      <div v-if="listLoading">Loading...</div>
      <div v-else-if="problems.length === 0">No problems yet.</div>
      <ul v-else class="problems">
        <li v-for="p in problems" :key="p.id">
          <strong>#{{ p.id }}</strong> - {{ p.slug }} ({{ p.difficulty }})
        </li>
      </ul>
    </section>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import client from '../api/client'

const problems = ref([])
const loading = ref(false)
const listLoading = ref(false)
const error = ref('')
const success = ref('')

const createForm = ref({
  slug: '',
  difficulty: 'easy',
  statementsJson: '[{"language":"en","title":"Title","statement":"Statement"}]',
  testsJson: '[{"input_data":"1 2","expected_output":"3","is_sample":true}]',
})

const updateForm = ref({
  id: null,
  slug: '',
  difficulty: 'easy',
  statementsJson: '[{"language":"en","title":"Title","statement":"Statement"}]',
  testsJson: '[{"input_data":"1 2","expected_output":"3","is_sample":true}]',
})

const deleteId = ref(null)

const resetMessages = () => {
  error.value = ''
  success.value = ''
}

const parseJsonField = (raw, fieldName) => {
  try {
    const parsed = JSON.parse(raw)
    if (!Array.isArray(parsed)) throw new Error('must be array')
    return parsed
  } catch {
    throw new Error(`${fieldName} must be valid JSON array`)
  }
}

const fetchProblems = async () => {
  listLoading.value = true
  resetMessages()
  try {
    const res = await client.get('/problems')
    problems.value = res.data || []
  } catch (err) {
    error.value = err.response?.data || 'Failed to load problems'
  } finally {
    listLoading.value = false
  }
}

const createProblem = async () => {
  loading.value = true
  resetMessages()
  try {
    const statements = parseJsonField(createForm.value.statementsJson, 'statements')
    const tests = parseJsonField(createForm.value.testsJson, 'tests')
    await client.post('/admin/problems', {
      slug: createForm.value.slug,
      difficulty: createForm.value.difficulty,
      statements,
      tests,
    })
    success.value = 'Problem created'
    await fetchProblems()
  } catch (err) {
    error.value = err.message || err.response?.data || 'Failed to create problem'
  } finally {
    loading.value = false
  }
}

const updateProblem = async () => {
  loading.value = true
  resetMessages()
  try {
    const statements = parseJsonField(updateForm.value.statementsJson, 'statements')
    const tests = parseJsonField(updateForm.value.testsJson, 'tests')
    await client.put(`/admin/problems?id=${updateForm.value.id}`, {
      slug: updateForm.value.slug,
      difficulty: updateForm.value.difficulty,
      statements,
      tests,
    })
    success.value = 'Problem updated'
    await fetchProblems()
  } catch (err) {
    error.value = err.message || err.response?.data || 'Failed to update problem'
  } finally {
    loading.value = false
  }
}

const deleteProblem = async () => {
  loading.value = true
  resetMessages()
  try {
    await client.delete(`/admin/problems?id=${deleteId.value}`)
    success.value = 'Problem deleted'
    await fetchProblems()
  } catch (err) {
    error.value = err.response?.data || 'Failed to delete problem'
  } finally {
    loading.value = false
  }
}

onMounted(fetchProblems)
</script>

<style scoped>
.admin-wrapper {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
  width: 100%;
}

.header-section {
  margin-bottom: 1.5rem;
}

.subtitle {
  color: var(--text-muted);
}

.layout {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
  margin-bottom: 1rem;
}

.panel {
  padding: 1rem;
  margin-bottom: 1rem;
}

.form,
.inline-form {
  display: flex;
  flex-direction: column;
  gap: 0.7rem;
}

.inline-form {
  flex-direction: row;
  align-items: center;
}

input,
select,
textarea {
  padding: 0.7rem;
  border: 1px solid var(--glass-border);
  border-radius: 6px;
  background: rgba(0, 0, 0, 0.3);
  color: #fff;
}

.status {
  color: #ff5f56;
  margin-bottom: 1rem;
}

.status.success {
  color: #27c93f;
}

.list-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.problems {
  margin-top: 0.7rem;
}

@media (max-width: 900px) {
  .layout {
    grid-template-columns: 1fr;
  }
}
</style>
