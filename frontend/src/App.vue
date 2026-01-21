<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { marked } from 'marked'
import DOMPurify from 'dompurify'

const error = ref(null)

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

const resp = ref('')

// Computed property for markdown rendering
const renderedMarkdown = computed(() => {
  if (!resp.value) return ''
  const raw = marked.parse(resp.value)
  return DOMPurify.sanitize(raw)
})

let sse = null

function fetchQuote() {
  if (sse) return

  // Clear previous quote when fetching new one
  resp.value = ''

  sse = new EventSource(`${API_URL}/stream`)

  sse.onopen = () => {
    console.log('SSE connection opened')
  }

  sse.onmessage = (event) => {
    const data = event.data
    if (!data) return

    if (data === 'stream complete') {
      closeSse()
      return
    }

    try {
      const parsed = JSON.parse(data)
      resp.value += parsed.chunk ?? data
    } catch {
      resp.value += data
    }
  }

  sse.onerror = (err) => {
    error.value = 'Connection failed' // TODO: improve error handling
    closeSse()
  }
}

function closeSse() {
  if (sse) {
    sse.close()
    sse = null
    console.log('SSE closed')
  }
}

onBeforeUnmount(() => {
  closeSse()
})

onMounted(() => {
  fetchQuote()
})
</script>

<template>
  <div class="container">
    <h1>Vue + Go on Render</h1>

    <div v-if="resp" class="card">
      <!-- Use v-html for rendered markdown -->
      <div class="markdown-content" v-html="renderedMarkdown"></div>
      <p class="author">- Gopher</p>
    </div>

    <!--div v-else-if="error" class="error">
      Error: {{ error }}
    </div-->

    <div v-else>
      Loading Gopher wisdom...
    </div>

    <button @click="fetchQuote">Get New Quote</button>
  </div>
</template>

<style scoped>
.container {
  text-align: center;
  font-family: sans-serif;
  margin-top: 50px;
  max-width: 800px;
  margin-left: auto;
  margin-right: auto;
  padding: 0 20px;
}

.card {
  background: #f4f4f4;
  padding: 20px 30px;
  border-radius: 8px;
  display: block;
  text-align: left;
}

.markdown-content {
  line-height: 1.6;
}

/* Basic markdown styling */
.markdown-content :deep(h1),
.markdown-content :deep(h2),
.markdown-content :deep(h3) {
  margin-top: 0.5em;
  margin-bottom: 0.5em;
}

.markdown-content :deep(code) {
  background: #e0e0e0;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: monospace;
}

.markdown-content :deep(pre) {
  background: #2d2d2d;
  color: #f8f8f2;
  padding: 15px;
  border-radius: 6px;
  overflow-x: auto;
}

.markdown-content :deep(pre code) {
  background: none;
  padding: 0;
  color: inherit;
}

.markdown-content :deep(blockquote) {
  border-left: 4px solid #00ADD8;
  margin-left: 0;
  padding-left: 15px;
  color: #555;
}

.author {
  text-align: right;
  font-style: italic;
  margin-top: 15px;
  color: #666;
}

.error {
  color: red;
}

button {
  margin-top: 20px;
  padding: 10px 20px;
  cursor: pointer;
  background: #00ADD8;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 16px;
}

button:hover {
  background: #0090b8;
}
</style>