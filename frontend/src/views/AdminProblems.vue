<template>
  <div class="admin-wrapper">
    <div class="header-section animate-fade-up">
      <h1>{{ ui.t('adminTitle') }}</h1>
      <p class="subtitle">{{ ui.t('adminSubtitle') }}</p>
    </div>

    <div class="status" v-if="error">{{ error }}</div>
    <div class="status success" v-if="success">{{ success }}</div>

    <div class="mode-switch">
      <button class="btn btn-outline btn-sm" :class="{ active: mode === 'create' }" @click="mode = 'create'">
        {{ ui.t('createSection') }}
      </button>
      <button class="btn btn-outline btn-sm" :class="{ active: mode === 'manage' }" @click="mode = 'manage'">
        {{ ui.t('manageSection') }}
      </button>
    </div>

    <div class="layout" v-if="mode === 'create'">
      <section class="glass-panel panel">
        <h2>{{ ui.t('createSection') }}</h2>
        <form @submit.prevent="createProblem" class="form">
          <input v-model="createForm.slug" :placeholder="ui.t('adminSlugPlaceholder')" required />
          <select v-model="createForm.difficulty" required>
            <option value="easy">easy</option>
            <option value="medium">medium</option>
            <option value="hard">hard</option>
          </select>

          <h3>{{ ui.t('adminStatements') }}</h3>
          <div v-for="(statement, idx) in createForm.statements" :key="idx" class="row-card">
            <select v-model="statement.language" required>
              <option value="en">EN</option>
              <option value="ru">RU</option>
              <option value="tm">TM</option>
            </select>
            <input v-model="statement.title" :placeholder="ui.t('adminTitlePlaceholder')" required />
            <textarea v-model="statement.statement" rows="4" :placeholder="ui.t('adminStatementPlaceholder')" required />
            <button class="btn btn-outline btn-sm" type="button" @click="removeStatement(createForm, idx)">{{ ui.t('adminRemove') }}</button>
          </div>
          <button class="btn btn-outline" type="button" @click="addStatement(createForm)">{{ ui.t('adminAddStatement') }}</button>

          <h3>{{ ui.t('adminTests') }}</h3>
          <div v-for="(test, idx) in createForm.tests" :key="idx" class="row-card">
            <textarea v-model="test.input_data" rows="2" :placeholder="ui.t('adminInput')" required />
            <textarea v-model="test.expected_output" rows="2" :placeholder="ui.t('adminExpected')" required />
            <label class="checkbox-line">
              <input type="checkbox" v-model="test.is_sample" />
              {{ ui.t('detailSampleTest') }}
            </label>
            <button class="btn btn-outline btn-sm" type="button" @click="removeTest(createForm, idx)">{{ ui.t('adminRemove') }}</button>
          </div>
          <button class="btn btn-outline" type="button" @click="addTest(createForm)">{{ ui.t('adminAddTest') }}</button>

          <button class="btn btn-primary" :disabled="loading">{{ loading ? ui.t('adminSaving') : ui.t('adminCreate') }}</button>
        </form>
      </section>
    </div>

    <div class="layout" v-else>
      <section class="glass-panel panel">
        <h2>{{ ui.t('manageSection') }}</h2>
        <form @submit.prevent="updateProblem" class="form">
          <input v-model.number="updateForm.id" type="number" min="1" :placeholder="ui.t('adminProblemIdPlaceholder')" required />
          <button class="btn btn-outline btn-sm" type="button" @click="loadProblemToUpdate" :disabled="loading">
            Load by ID
          </button>
          <input v-model="updateForm.slug" :placeholder="ui.t('adminNewSlugPlaceholder')" required />
          <select v-model="updateForm.difficulty" required>
            <option value="easy">easy</option>
            <option value="medium">medium</option>
            <option value="hard">hard</option>
          </select>

          <h3>{{ ui.t('adminStatements') }}</h3>
          <div v-for="(statement, idx) in updateForm.statements" :key="idx" class="row-card">
            <select v-model="statement.language" required>
              <option value="en">EN</option>
              <option value="ru">RU</option>
              <option value="tm">TM</option>
            </select>
            <input v-model="statement.title" :placeholder="ui.t('adminTitlePlaceholder')" required />
            <textarea v-model="statement.statement" rows="4" :placeholder="ui.t('adminStatementPlaceholder')" required />
            <button class="btn btn-outline btn-sm" type="button" @click="removeStatement(updateForm, idx)">{{ ui.t('adminRemove') }}</button>
          </div>
          <button class="btn btn-outline" type="button" @click="addStatement(updateForm)">{{ ui.t('adminAddStatement') }}</button>

          <h3>{{ ui.t('adminTests') }}</h3>
          <div v-for="(test, idx) in updateForm.tests" :key="idx" class="row-card">
            <textarea v-model="test.input_data" rows="2" :placeholder="ui.t('adminInput')" required />
            <textarea v-model="test.expected_output" rows="2" :placeholder="ui.t('adminExpected')" required />
            <label class="checkbox-line">
              <input type="checkbox" v-model="test.is_sample" />
              {{ ui.t('detailSampleTest') }}
            </label>
            <button class="btn btn-outline btn-sm" type="button" @click="removeTest(updateForm, idx)">{{ ui.t('adminRemove') }}</button>
          </div>
          <button class="btn btn-outline" type="button" @click="addTest(updateForm)">{{ ui.t('adminAddTest') }}</button>

          <button class="btn btn-primary" :disabled="loading">{{ loading ? ui.t('adminUpdating') : ui.t('adminUpdate') }}</button>
        </form>
      </section>
    </div>

    <section class="glass-panel panel" v-if="mode === 'manage'">
      <h2>{{ ui.t('adminDeleteProblem') }}</h2>
      <form @submit.prevent="deleteProblem" class="inline-form">
        <input v-model.number="deleteId" type="number" min="1" :placeholder="ui.t('adminProblemIdPlaceholder')" required />
        <button class="btn btn-outline" :disabled="loading">{{ loading ? ui.t('adminDeleting') : ui.t('adminDelete') }}</button>
      </form>
    </section>

    <section class="glass-panel panel">
      <div class="list-head">
        <h2>{{ ui.t('adminExisting') }}</h2>
        <button class="btn btn-outline" @click="fetchProblems" :disabled="loading">{{ ui.t('adminRefresh') }}</button>
      </div>
      <div v-if="listLoading">{{ ui.t('adminLoading') }}</div>
      <div v-else-if="problems.length === 0">{{ ui.t('adminNoProblems') }}</div>
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
import { useUIStore } from '../stores/ui'

const problems = ref([])
const loading = ref(false)
const listLoading = ref(false)
const error = ref('')
const success = ref('')
const mode = ref('create')
const ui = useUIStore()

const createForm = ref({
  slug: '',
  difficulty: 'easy',
  statements: [{ language: 'en', title: '', statement: '' }],
  tests: [{ input_data: '', expected_output: '', is_sample: true }],
})

const updateForm = ref({
  id: null,
  slug: '',
  difficulty: 'easy',
  statements: [{ language: 'en', title: '', statement: '' }],
  tests: [{ input_data: '', expected_output: '', is_sample: true }],
})

const deleteId = ref(null)

const resetMessages = () => {
  error.value = ''
  success.value = ''
}

const addStatement = (formRef) => {
  formRef.value.statements.push({ language: 'en', title: '', statement: '' })
}

const removeStatement = (formRef, idx) => {
  if (formRef.value.statements.length === 1) return
  formRef.value.statements.splice(idx, 1)
}

const addTest = (formRef) => {
  formRef.value.tests.push({ input_data: '', expected_output: '', is_sample: false })
}

const removeTest = (formRef, idx) => {
  if (formRef.value.tests.length === 1) return
  formRef.value.tests.splice(idx, 1)
}

const ensureFormValid = (form) => {
  if (!form.slug.trim()) throw new Error('Slug is required')
  if (!form.statements.length) throw new Error('At least one statement is required')
  if (!form.tests.length) throw new Error('At least one test is required')
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
    ensureFormValid(createForm.value)
    await client.post('/admin/problems', {
      slug: createForm.value.slug,
      difficulty: createForm.value.difficulty,
      statements: createForm.value.statements,
      tests: createForm.value.tests,
    })
    success.value = 'Problem created'
    createForm.value = {
      slug: '',
      difficulty: 'easy',
      statements: [{ language: 'en', title: '', statement: '' }],
      tests: [{ input_data: '', expected_output: '', is_sample: true }],
    }
    await fetchProblems()
  } catch (err) {
    error.value = err.message || err.response?.data || 'Failed to create problem'
  } finally {
    loading.value = false
  }
}

const loadProblemToUpdate = async () => {
  if (!updateForm.value.id) {
    error.value = 'Enter problem id first'
    return
  }
  loading.value = true
  resetMessages()
  try {
    const res = await client.get(`/problems/detail?id=${updateForm.value.id}`)
    updateForm.value.slug = res.data.problem.slug
    updateForm.value.difficulty = res.data.problem.difficulty
    updateForm.value.statements = (res.data.statements || []).map((s) => ({
      language: s.language?.slice(0, 2) || 'en',
      title: s.title || '',
      statement: s.statement || '',
    }))
    updateForm.value.tests = (res.data.tests || []).map((t) => ({
      input_data: t.input_data || '',
      expected_output: t.expected_output || '',
      is_sample: !!t.is_sample,
    }))
    if (!updateForm.value.statements.length) addStatement(updateForm)
    if (!updateForm.value.tests.length) addTest(updateForm)
    success.value = 'Problem loaded'
  } catch (err) {
    error.value = err.response?.data || 'Failed to load problem by id'
  } finally {
    loading.value = false
  }
}

const updateProblem = async () => {
  loading.value = true
  resetMessages()
  try {
    ensureFormValid(updateForm.value)
    await client.put(`/admin/problems?id=${updateForm.value.id}`, {
      slug: updateForm.value.slug,
      difficulty: updateForm.value.difficulty,
      statements: updateForm.value.statements,
      tests: updateForm.value.tests,
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
  grid-template-columns: 1fr;
  gap: 1rem;
  margin-bottom: 1rem;
}

.panel {
  padding: 1rem;
  margin-bottom: 1rem;
}

h3 {
  margin-top: 0.6rem;
  margin-bottom: 0.2rem;
  color: var(--text-muted);
  font-size: 0.95rem;
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

.row-card {
  border: 1px solid var(--glass-border);
  border-radius: 8px;
  padding: 0.8rem;
  display: flex;
  flex-direction: column;
  gap: 0.55rem;
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

.mode-switch {
  display: flex;
  gap: 0.6rem;
  margin-bottom: 1rem;
}

.mode-switch .active {
  border-color: rgba(59, 130, 246, 0.4);
  background: rgba(59, 130, 246, 0.1);
}

.checkbox-line {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: var(--text-muted);
  font-size: 0.9rem;
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
