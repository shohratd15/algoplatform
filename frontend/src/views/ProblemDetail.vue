<template>
  <div class="detail-wrapper animate-fade-up" v-if="!loading && problem">
    <!-- Mobile tabs -->
    <div class="mobile-tabs">
      <button 
        class="mobile-tab" 
        :class="{ active: activeMobileTab === 'statement' }"
        @click="activeMobileTab = 'statement'"
      >
        {{ ui.t('detailTask') || 'Задача' }}
      </button>
      <button 
        class="mobile-tab" 
        :class="{ active: activeMobileTab === 'editor' }"
        @click="activeMobileTab = 'editor'"
      >
        {{ ui.t('detailEditor') || 'Редактор' }}
      </button>
    </div>

    <div class="split-pane">
      <!-- Left: Statement -->
      <div class="pane left-pane glass-panel" :class="{ 'mobile-hidden': activeMobileTab !== 'statement' }">
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
      <div class="pane right-pane" :class="{ 'mobile-hidden': activeMobileTab !== 'editor' }">
        <div class="editor-header glass-panel">
          <select v-model="languageId" class="lang-select">
            <option :value="71">Python (3.8.1)</option>
            <option :value="60">Go (1.13.5)</option>
            <option :value="54">C++ (GCC 9.2.0)</option>
            <option :value="62">Java (OpenJDK 13.0.1)</option>
            <option :value="63">JavaScript (Node.js)</option>
          </select>
          <button class="btn btn-primary" @click="submitCode" :disabled="submitting">
            {{ submitting ? ui.t('detailSubmitting') : ui.t('detailRunCode') }}
          </button>
        </div>

        <div class="editor-shell glass-panel" style="padding: 1rem 0;">
          <div ref="editorContainer" style="width: 100%; height: 400px;"></div>
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

          <div class="result-details" v-if="submissionResult.message || submissionResult.compile_output || submissionResult.stderr || submissionResult.stdout || submissionResult.expected_output">
            <div v-if="submissionResult.message">
              <strong>Message:</strong>
              <pre>{{ submissionResult.message }}</pre>
            </div>
            <div v-if="submissionResult.stdout">
              <strong>{{ ui.t('detailYourOutput') }}:</strong>
              <pre>{{ submissionResult.stdout }}</pre>
            </div>
            <div v-if="submissionResult.expected_output">
              <strong>{{ ui.t('detailOutput') }}:</strong>
              <pre>{{ submissionResult.expected_output }}</pre>
            </div>
            <div v-if="submissionResult.compile_output">
              <strong>Compile output:</strong>
              <pre>{{ submissionResult.compile_output }}</pre>
            </div>
            <div v-if="submissionResult.stderr">
              <strong>Runtime stderr:</strong>
              <pre>{{ submissionResult.stderr }}</pre>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <div v-else class="loading-state">
    {{ ui.t('detailLoading') }}
  </div>
</template>

<script setup>
import { ref, shallowRef, markRaw, onMounted, onUnmounted, computed, watch, nextTick } from 'vue'
import { useUIStore } from '../stores/ui'
import { useRoute, useRouter } from 'vue-router'
import * as monaco from 'monaco-editor'
import editorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker'
import jsonWorker from 'monaco-editor/esm/vs/language/json/json.worker?worker'
import cssWorker from 'monaco-editor/esm/vs/language/css/css.worker?worker'
import htmlWorker from 'monaco-editor/esm/vs/language/html/html.worker?worker'
import tsWorker from 'monaco-editor/esm/vs/language/typescript/ts.worker?worker'
import client from '../api/client'

self.MonacoEnvironment = {
  getWorker(_, label) {
    if (label === 'json') return new jsonWorker()
    if (label === 'css' || label === 'scss' || label === 'less') return new cssWorker()
    if (label === 'html' || label === 'handlebars' || label === 'razor') return new htmlWorker()
    if (label === 'typescript' || label === 'javascript') return new tsWorker()
    return new editorWorker()
  }
}

const route = useRoute()
const ui = useUIStore()
const router = useRouter()

const loading = ref(true)
const problem = ref(null)
const languageId = ref(60)
const sourceCode = ref('')
const submitting = ref(false)
const submissionResult = ref(null)
const submitError = ref('')
const selectedStatementLanguage = ref('')
const editorInstance = shallowRef(null)
const activeMobileTab = ref('statement')

// Таймер храним в ref, чтобы очистить при unmount
const pollTimer = ref(null)
const MAX_POLL_ATTEMPTS = 15
let pollAttempts = 0

const sampleTests = computed(() => {
  if (!problem.value?.tests) return []
  return problem.value.tests.filter((t) => t.is_sample)
})

// Monaco Editor configuration
const currentLanguage = computed(() => {
  const langMap = {
    71: 'python',
    60: 'go',
    54: 'cpp',
    62: 'java',
    63: 'javascript'
  }
  return langMap[languageId.value] || 'plaintext'
})

const editorTheme = 'vs-dark'

const editorOptions = {
  fontSize: 14,
  fontFamily: "'JetBrains Mono', 'Fira Code', 'Consolas', monospace",
  minimap: { enabled: false },
  scrollBeyondLastLine: false,
  automaticLayout: true,
  tabSize: 2,
  insertSpaces: true,
  wordWrap: 'on',
  lineNumbers: 'on',
  glyphMargin: false,
  folding: true,
  lineDecorationsWidth: 10,
  lineNumbersMinChars: 3,
  renderLineHighlight: 'line',
  scrollbar: {
    verticalScrollbarSize: 10,
    horizontalScrollbarSize: 10
  },
  padding: { top: 16, bottom: 16 }
}

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
    if (langs.length && !langs.includes(selectedStatementLanguage.value)) {
      selectedStatementLanguage.value = langs[0]
    }
  } catch (err) {
    // 401 обрабатывает interceptor; остальное — редирект на список
    if (err.response?.status !== 401) {
      router.push('/problems')
    }
  } finally {
    loading.value = false
    nextTick(() => {
      initEditor()
    })
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

const editorContainer = ref(null)

const handleEditorMount = () => {
  // Editor is now initialized in onMounted
}

const initEditor = () => {
  if (editorContainer.value && !editorInstance.value) {
    monaco.editor.defineTheme('transparent-dark', {
      base: 'vs-dark',
      inherit: true,
      rules: [],
      colors: {
        'editor.background': '#00000000',
        'editor.marginBackground': '#00000000',
        'editorGutter.background': '#00000000',
        'editorLineNumber.background': '#00000000'
      }
    })

    const editor = monaco.editor.create(editorContainer.value, {
      value: sourceCode.value,
      language: currentLanguage.value,
      theme: 'transparent-dark',
      fontSize: 14,
      fontFamily: "'JetBrains Mono', 'Fira Code', 'Consolas', monospace",
      minimap: { enabled: false },
      scrollBeyondLastLine: false,
      automaticLayout: true,
      lineNumbers: 'on',
      roundedSelection: false,
      readOnly: false,
      cursorStyle: 'line',
    })
    editorInstance.value = markRaw(editor)
    
    editorInstance.value.onDidChangeModelContent(() => {
      sourceCode.value = editorInstance.value.getValue()
    })
  }
}

onMounted(fetchProblem)

watch(languageId, () => {
  if (editorInstance.value) {
    monaco.editor.setModelLanguage(
      editorInstance.value.getModel(),
      currentLanguage.value
    )
  }
})

// Очищаем таймер при уходе со страницы — не допускаем утечки
onUnmounted(() => {
  if (pollTimer.value) clearTimeout(pollTimer.value)
  if (editorInstance.value) {
    editorInstance.value.dispose()
  }
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
  display: block;
  overflow: hidden;
  border-radius: 12px;
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

/* Mobile responsive styles */
.mobile-tabs {
  display: none;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.mobile-tab {
  flex: 1;
  padding: 0.75rem 1rem;
  background: var(--glass-bg);
  border: 1px solid var(--glass-border);
  border-radius: 8px;
  color: var(--text-muted);
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.mobile-tab.active {
  background: rgba(59, 130, 246, 0.15);
  border-color: rgba(59, 130, 246, 0.4);
  color: #fff;
}

@media (max-width: 768px) {
  .detail-wrapper {
    padding: 1rem;
    height: auto;
    min-height: calc(100vh - 70px);
    overflow-y: auto;
  }

  .split-pane {
    flex-direction: column;
    gap: 1rem;
  }

  .pane.mobile-hidden {
    display: none;
  }

  .pane {
    width: 100%;
  }

  .left-pane {
    max-height: none;
    padding: 1.5rem;
  }

  .right-pane {
    min-height: 60vh;
  }

  .editor-header {
    flex-wrap: wrap;
    gap: 0.75rem;
    padding: 0.75rem 1rem;
  }

  .lang-select {
    flex: 1;
    min-width: 120px;
  }

  .editor-shell {
    min-height: 300px;
  }

  .result-box {
    margin-top: 1rem;
  }

  .problem-header h2 {
    font-size: 1.5rem;
  }

  .content {
    font-size: 0.95rem;
  }

  .example-box {
    padding: 1rem;
  }
}
</style>
