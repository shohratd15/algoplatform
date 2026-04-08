<template>
  <div class="detail-wrapper animate-fade-up" v-if="!loading && problem">
    <div class="split-pane">
      <!-- Left: Statement -->
      <div class="pane left-pane glass-panel">
        <router-link to="/problems" class="back-link">{{ ui.t('detailBack') }}</router-link>

        <div class="problem-header">
          <h2>{{ problem.problem.slug.replace(/-/g, ' ') }}</h2>
          <span class="diff-badge" :class="problem.problem.difficulty.toLowerCase()">
            {{ problem.problem.difficulty }}
          </span>
        </div>

        <div class="content" v-if="problem.statements && problem.statements.length > 0">
          <div class="lang-switch">
            <button
              v-for="lang in availableStatementLanguages"
              :key="lang"
              class="lang-btn"
              :class="{ active: selectedStatementLanguage === lang }"
              @click="selectedStatementLanguage = lang"
            >
              {{ languageLabel(lang) }}
            </button>
          </div>
          <h3>{{ selectedStatement?.title }}</h3>
          <p class="statement-text">{{ selectedStatement?.statement }}</p>
        </div>

        <div class="examples" v-if="sampleTests.length > 0">
          <h3>{{ ui.t('detailExamples') }}</h3>
          <div v-for="(t, i) in sampleTests" :key="i" class="example-box">
            <div class="io-block">
              <strong>{{ ui.t('detailInput') }}</strong>
              <pre>{{ t.input_data }}</pre>
            </div>
            <div class="io-block">
              <strong>{{ ui.t('detailOutput') }}</strong>
              <pre>{{ t.expected_output }}</pre>
            </div>
          </div>
        </div>
      </div>

      <!-- Right: Editor -->
      <div class="pane right-pane">
        <div class="editor-header glass-panel">
          <select v-model="languageId" class="lang-select">
            <option :value="71">Python (3.8.1)</option>
            <option :value="60">Go (1.13.5)</option>
            <option :value="54">C++ (GCC 9.2.0)</option>
            <option :value="62">Java (OpenJDK 13.0.1)</option>
            <option :value="93">JavaScript (Node.js)</option>
          </select>
          <button class="btn btn-primary" @click="submitCode" :disabled="submitting">
            {{ submitting ? ui.t('detailSubmitting') : ui.t('detailRunCode') }}
          </button>
        </div>

        <div class="editor-shell glass-panel">
          <div class="line-numbers" ref="lineNumbersRef">
            <div v-for="line in lineCount" :key="line">{{ line }}</div>
          </div>
          <textarea
            ref="editorRef"
            class="code-editor"
            v-model="sourceCode"
            spellcheck="false"
            :placeholder="ui.t('detailWriteCode')"
            @keydown.tab.prevent="handleTab"
            @scroll="syncScroll"
          ></textarea>
        </div>

        <!-- Ошибка сабмита — в UI, не через alert() -->
        <div class="error-msg" v-if="submitError">{{ submitError }}</div>

        <div class="result-box glass-panel" v-if="submissionResult">
          <h3>{{ ui.t('detailResult') }}</h3>
          <div class="status-badge" :class="submissionResult.status">
            {{ submissionResult.status.replace(/_/g, ' ').toUpperCase() }}
          </div>
          <button
            class="btn btn-outline btn-sm"
            @click="pollResult"
            v-if="submissionResult.status === 'queued' || submissionResult.status === 'running'"
          >
            Refresh Status
          </button>
        </div>
      </div>
    </div>
  </div>
  <div v-else class="loading-state">
    {{ ui.t('detailLoading') }}
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useUIStore } from '../stores/ui'
import { useRoute, useRouter } from 'vue-router'
import client from '../api/client'

const route = useRoute()
const ui = useUIStore()
const router = useRouter()

const loading = ref(true)
const problem = ref(null)
const languageId = ref(71)
const sourceCode = ref('')
const submitting = ref(false)
const submissionResult = ref(null)
const submitError = ref('')
const selectedStatementLanguage = ref('eng')
const editorRef = ref(null)
const lineNumbersRef = ref(null)

// Таймер храним в ref, чтобы очистить при unmount
const pollTimer = ref(null)
const MAX_POLL_ATTEMPTS = 15
let pollAttempts = 0

const sampleTests = computed(() => {
  if (!problem.value?.tests) return []
  return problem.value.tests.filter((t) => t.is_sample)
})
const lineCount = computed(() => Math.max(1, sourceCode.value.split('\n').length))
const availableStatementLanguages = computed(() => {
  const stmts = problem.value?.statements || []
  return [...new Set(stmts.map((s) => normalizeLang(s.language)).filter(Boolean))]
})
const selectedStatement = computed(() => {
  const stmts = problem.value?.statements || []
  if (!stmts.length) return null
  const found = stmts.find((s) => normalizeLang(s.language) === selectedStatementLanguage.value)
  return found || stmts[0]
})

const fetchProblem = async () => {
  try {
    const res = await client.get(`/problems/detail?id=${route.params.id}`)
    problem.value = res.data
    const langs = [...new Set((res.data?.statements || []).map((s) => normalizeLang(s.language)).filter(Boolean))]
    if (langs.length) selectedStatementLanguage.value = langs[0]
  } catch (err) {
    // 401 обрабатывает interceptor; остальное — редирект на список
    if (err.response?.status !== 401) {
      router.push('/problems')
    }
  } finally {
    loading.value = false
  }
}

const normalizeLang = (lang) => {
  const l = (lang || '').toLowerCase()
  if (l === 'en' || l === 'eng') return 'eng'
  if (l === 'ru' || l === 'rus') return 'rus'
  if (l === 'tm' || l === 'tkm' || l === 'tk') return 'tkm'
  return l
}

const languageLabel = (lang) => {
  if (lang === 'eng') return 'EN'
  if (lang === 'rus') return 'RU'
  if (lang === 'tkm') return 'TM'
  return (lang || '').toUpperCase()
}

const submitCode = async () => {
  submitting.value = true
  submitError.value = ''
  submissionResult.value = null
  pollAttempts = 0

  try {
    const res = await client.post('/submissions', {
      problem_id: parseInt(route.params.id),
      language_id: languageId.value,
      source_code: sourceCode.value,
    })

    submissionResult.value = { id: res.data.id, status: 'queued' }
    schedulePoll(1000)
  } catch (err) {
    submitError.value = 'Submission failed: ' + (err.response?.data || err.message)
  } finally {
    submitting.value = false
  }
}

const schedulePoll = (delay) => {
  pollTimer.value = setTimeout(pollResult, delay)
}

const pollResult = async () => {
  if (!submissionResult.value?.id) return
  if (pollAttempts >= MAX_POLL_ATTEMPTS) {
    submissionResult.value.status = 'error'
    return
  }

  pollAttempts++

  try {
    const res = await client.get(`/submissions?id=${submissionResult.value.id}`)
    submissionResult.value = res.data

    if (res.data.status === 'queued' || res.data.status === 'running') {
      schedulePoll(2000)
    }
  } catch (err) {
    console.error('Poll failed', err)
  }
}

const syncScroll = () => {
  if (!editorRef.value || !lineNumbersRef.value) return
  lineNumbersRef.value.scrollTop = editorRef.value.scrollTop
}

const handleTab = (e) => {
  const el = e.target
  const tab = '  '
  const start = el.selectionStart
  const end = el.selectionEnd
  const current = sourceCode.value

  sourceCode.value = `${current.slice(0, start)}${tab}${current.slice(end)}`
  requestAnimationFrame(() => {
    el.selectionStart = el.selectionEnd = start + tab.length
    syncScroll()
  })
}

onMounted(fetchProblem)

// Очищаем таймер при уходе со страницы — не допускаем утечки
onUnmounted(() => {
  if (pollTimer.value) clearTimeout(pollTimer.value)
})
</script>

<style scoped>
.detail-wrapper {
  height: calc(100vh - 70px);
  padding: 1.5rem;
  overflow: hidden;
}

.split-pane {
  display: flex;
  gap: 1.5rem;
  height: 100%;
}

.pane {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.left-pane {
  overflow-y: auto;
  padding: 2rem;
}

.back-link {
  font-size: 0.9rem;
  color: var(--text-muted);
  margin-bottom: 2rem;
  display: inline-block;
}

.problem-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 2rem;
}

.problem-header h2 {
  font-size: 2rem;
  text-transform: capitalize;
  margin: 0;
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

.content {
  margin-bottom: 3rem;
  color: #ccc;
}

.content h3 { margin-bottom: 1rem; }
.statement-text { white-space: pre-wrap; }

.lang-switch {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.lang-btn {
  border: 1px solid var(--glass-border);
  background: rgba(0, 0, 0, 0.25);
  color: var(--text-muted);
  border-radius: 6px;
  padding: 0.35rem 0.7rem;
  cursor: pointer;
}

.lang-btn.active {
  color: #fff;
  border-color: rgba(59, 130, 246, 0.4);
  background: rgba(59, 130, 246, 0.1);
}

.example-box {
  background: rgba(0, 0, 0, 0.3);
  padding: 1.5rem;
  border-radius: 8px;
  margin-bottom: 1.5rem;
  border: 1px solid rgba(255, 255, 255, 0.05);
}

.io-block { margin-bottom: 1rem; }
.io-block:last-child { margin-bottom: 0; }
.io-block pre {
  margin-top: 0.5rem;
  color: #a5b4fc;
  font-family: monospace;
}

.right-pane { gap: 1rem; }

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 1.5rem;
}

.lang-select {
  background: rgba(0, 0, 0, 0.3);
  color: #fff;
  border: 1px solid var(--glass-border);
  padding: 0.5rem 1rem;
  border-radius: 6px;
  outline: none;
}

.editor-shell {
  flex: 1;
  display: grid;
  grid-template-columns: 56px 1fr;
  overflow: hidden;
}

.line-numbers {
  padding: 1rem 0.5rem;
  overflow: hidden;
  background: rgba(0, 0, 0, 0.35);
  border-right: 1px solid var(--glass-border);
  color: var(--text-muted);
  font-family: 'Fira Code', 'Courier New', monospace;
  font-size: 0.9rem;
  line-height: 1.6;
  text-align: right;
}

.code-editor {
  width: 100%;
  height: 100%;
  border: none;
  outline: none;
  background: transparent;
  color: #fff;
  resize: none;
  font-family: 'Fira Code', 'Courier New', monospace;
  font-size: 1rem;
  padding: 1rem;
  line-height: 1.6;
  white-space: pre;
}

.error-msg {
  color: #ff5f56;
  font-size: 0.9rem;
  background: rgba(255, 95, 86, 0.1);
  padding: 0.75rem 1rem;
  border-radius: 8px;
}

.result-box {
  padding: 1.5rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.status-badge {
  padding: 0.5rem 1rem;
  border-radius: 6px;
  font-weight: bold;
}

.status-badge.accepted     { background: rgba(39, 201, 63, 0.2);  color: #27c93f; }
.status-badge.wrong_answer { background: rgba(255, 95, 86, 0.2);  color: #ff5f56; }
.status-badge.error        { background: rgba(255, 95, 86, 0.2);  color: #ff5f56; }
.status-badge.time_limit   { background: rgba(255, 95, 86, 0.2);  color: #ff5f56; }
.status-badge.queued,
.status-badge.running      { background: rgba(255, 189, 46, 0.2); color: #ffbd2e; }

.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  height: calc(100vh - 70px);
  color: var(--text-muted);
}
</style>
