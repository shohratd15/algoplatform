<template>
  <div class="detail-wrapper animate-fade-up" v-if="!loading && problem">
    <div class="split-pane">
      <!-- Left: Statement -->
      <div class="pane left-pane glass-panel">
        <router-link to="/problems" class="back-link">← Back to problems</router-link>
        
        <div class="problem-header">
          <h2>{{ problem.problem.slug.replace(/-/g, ' ') }}</h2>
          <span class="diff-badge" :class="problem.problem.difficulty.toLowerCase()">
            {{ problem.problem.difficulty }}
          </span>
        </div>

        <div class="content" v-if="problem.statements && problem.statements.length > 0">
          <!-- We just show the first statement translation for simplicity -->
          <h3>{{ problem.statements[0].title }}</h3>
          <p class="statement-text">{{ problem.statements[0].statement }}</p>
        </div>

        <div class="examples" v-if="problem.tests && problem.tests.length > 0">
          <h3>Examples</h3>
          <div v-for="(t, i) in sampleTests" :key="i" class="example-box">
            <div class="io-block">
              <strong>Input:</strong>
              <pre>{{ t.input_data }}</pre>
            </div>
            <div class="io-block">
              <strong>Expected Output:</strong>
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
            {{ submitting ? 'Submitting...' : 'Run Code' }}
          </button>
        </div>

        <textarea 
          class="code-editor glass-panel" 
          v-model="sourceCode" 
          spellcheck="false"
          placeholder="Write your code here..."
        ></textarea>
        
        <div class="result-box glass-panel" v-if="submissionResult">
          <h3>Submission Result</h3>
          <div class="status-badge" :class="submissionResult.status">
            {{ submissionResult.status.replace('_', ' ').toUpperCase() }}
          </div>
          <button class="btn btn-outline btn-sm" @click="pollResult" v-if="submissionResult.status === 'queued' || submissionResult.status === 'running'">
            Refresh Status
          </button>
        </div>
      </div>
    </div>
  </div>
  <div v-else class="loading-state">
    Loading challenge...
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'

const route = useRoute()
const router = useRouter()

const loading = ref(true)
const problem = ref(null)
const languageId = ref(71) // Default to Python
const sourceCode = ref('')
const submitting = ref(false)
const submissionResult = ref(null)

const sampleTests = computed(() => {
  if (!problem.value?.tests) return []
  return problem.value.tests.filter(t => t.is_sample)
})

const fetchProblem = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await axios.get(`/api/problems/detail?id=${route.params.id}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    problem.value = res.data
  } catch (err) {
    console.error(err)
    router.push('/problems')
  } finally {
    loading.value = false
  }
}

const submitCode = async () => {
  submitting.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await axios.post('/api/submissions', {
      problem_id: parseInt(route.params.id),
      language_id: languageId.value,
      source_code: sourceCode.value
    }, {
      headers: { Authorization: `Bearer ${token}` }
    })
    
    submissionResult.value = { id: res.data.id, status: 'queued' }
    
    // Start auto-poll
    setTimeout(pollResult, 1000)
    
  } catch (err) {
    alert("Submission failed: " + (err.response?.data || err.message))
  } finally {
    submitting.value = false
  }
}

const pollResult = async () => {
  if (!submissionResult.value?.id) return
  
  try {
    const token = localStorage.getItem('token')
    const res = await axios.get(`/api/submissions?id=${submissionResult.value.id}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    
    submissionResult.value = res.data
    
    if (res.data.status === 'queued' || res.data.status === 'running') {
      setTimeout(pollResult, 2000) // Poll again
    }
  } catch (err) {
    console.error("Poll failed", err)
  }
}

onMounted(() => {
  if (!localStorage.getItem('token')) {
    router.push('/login')
    return
  }
  fetchProblem()
})
</script>

<style scoped>
.detail-wrapper {
  height: calc(100vh - 70px); /* Fill screen minus navbar */
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
.diff-badge.easy { background: rgba(39, 201, 63, 0.15); color: #27c93f; }
.diff-badge.medium { background: rgba(255, 189, 46, 0.15); color: #ffbd2e; }
.diff-badge.hard { background: rgba(255, 95, 86, 0.15); color: #ff5f56; }

.content {
  margin-bottom: 3rem;
  color: #ccc;
}

.content h3 {
  margin-bottom: 1rem;
}

.statement-text {
  white-space: pre-wrap; /* render formatting */
}

.example-box {
  background: rgba(0,0,0,0.3);
  padding: 1.5rem;
  border-radius: 8px;
  margin-bottom: 1.5rem;
  border: 1px solid rgba(255,255,255,0.05);
}

.io-block {
  margin-bottom: 1rem;
}

.io-block:last-child {
  margin-bottom: 0;
}

.io-block pre {
  margin-top: 0.5rem;
  color: #a5b4fc;
  font-family: monospace;
}

/* Right pane */
.right-pane {
  gap: 1rem;
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 1.5rem;
}

.lang-select {
  background: rgba(0,0,0,0.3);
  color: #fff;
  border: 1px solid var(--glass-border);
  padding: 0.5rem 1rem;
  border-radius: 6px;
  outline: none;
}

.code-editor {
  flex: 1;
  resize: none;
  font-family: 'Fira Code', 'Courier New', monospace;
  font-size: 1rem;
  padding: 1.5rem;
  white-space: pre;
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

.status-badge.accepted { background: rgba(39, 201, 63, 0.2); color: #27c93f; }
.status-badge.wrong_answer { background: rgba(255, 95, 86, 0.2); color: #ff5f56; }
.status-badge.queued, .status-badge.running { background: rgba(255, 189, 46, 0.2); color: #ffbd2e; }

.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  height: calc(100vh - 70px);
  color: var(--text-muted);
}
</style>
