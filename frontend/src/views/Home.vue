<template>
  <div class="home-container">
    <div class="hero">
      <div class="hero-content animate-fade-up">
        <h1>Master Algorithms.<br>Level up your code.</h1>
        <p class="subtitle">Platform designed to help you prepare for technical interviews and competitive programming with a premium environment.</p>
        <div class="cta-group">
          <router-link to="/problems" class="btn btn-primary btn-large">Start Coding</router-link>
          <router-link to="/register" class="btn btn-outline btn-large">Create Account</router-link>
        </div>
      </div>
      
      <!-- Visual Element -->
      <div class="hero-visual animate-fade-up" style="animation-delay: 0.2s">
        <div class="code-window glass-panel">
          <div class="window-header">
            <span class="dot red"></span>
            <span class="dot yellow"></span>
            <span class="dot green"></span>
          </div>
          <pre><code><span class="keyword">func</span> <span class="func-name">BinarySearch</span>(arr []<span class="type">int</span>, target <span class="type">int</span>) <span class="type">int</span> {
    l, r := <span class="number">0</span>, <span class="keyword">len</span>(arr)-<span class="number">1</span>
    <span class="keyword">for</span> l &lt;= r {
        m := l + (r-l)/<span class="number">2</span>
        <span class="keyword">if</span> arr[m] == target {
            <span class="keyword">return</span> m
        }
        <span class="keyword">if</span> arr[m] &lt; target {
            l = m + <span class="number">1</span>
        } <span class="keyword">else</span> {
            r = m - <span class="number">1</span>
        }
    }
    <span class="keyword">return</span> -<span class="number">1</span>
}</code></pre>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.home-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: calc(100vh - 70px);
  padding: 2rem;
}

.hero {
  display: flex;
  align-items: center;
  justify-content: space-between;
  max-width: 1200px;
  width: 100%;
  gap: 4rem;
}

.hero-content {
  flex: 1;
  max-width: 600px;
}

.subtitle {
  font-size: 1.2rem;
  color: var(--text-muted);
  margin-bottom: 2.5rem;
}

.cta-group {
  display: flex;
  gap: 1rem;
}

.btn-large {
  padding: 1rem 2rem;
  font-size: 1.1rem;
}

.hero-visual {
  flex: 1;
  display: flex;
  justify-content: flex-end;
}

.code-window {
  width: 100%;
  max-width: 500px;
  background: rgba(10, 12, 16, 0.7);
  padding: 0;
  overflow: hidden;
}

.window-header {
  padding: 1rem;
  background: rgba(255, 255, 255, 0.05);
  display: flex;
  gap: 0.5rem;
  border-bottom: 1px solid var(--glass-border);
}

.dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}
.red { background: #ff5f56; }
.yellow { background: #ffbd2e; }
.green { background: #27c93f; }

pre {
  padding: 1.5rem;
  font-family: 'Fira Code', 'Courier New', monospace;
  font-size: 0.95rem;
  line-height: 1.5;
  color: #e5e7eb;
}

/* Syntax highlighting fake classes */
.keyword { color: #f472b6; }
.func-name { color: #60a5fa; }
.type { color: #34d399; }
.number { color: #fbbf24; }

@media (max-width: 900px) {
  .hero {
    flex-direction: column;
    text-align: center;
  }
  .cta-group {
    justify-content: center;
  }
}
</style>
