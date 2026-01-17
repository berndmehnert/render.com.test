<script setup>
import { ref, onMounted } from 'vue'

const quote = ref(null)
const error = ref(null)

// Use an environment variable for the API URL
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

const fetchQuote = async () => {
  try {
    const res = await fetch(`${API_URL}/api/quote`)
    if (!res.ok) throw new Error('Failed to connect to Go')
    quote.value = await res.json()
  } catch (err) {
    error.value = err.message
  }
}

onMounted(() => {
  fetchQuote()
})
</script>

<template>
  <div class="container">
    <h1>Vue + Go on Render</h1>
    
    <div v-if="quote" class="card">
      <h2>"{{ quote.message }}"</h2>
      <p>- {{ quote.author }}</p>
    </div>

    <div v-else-if="error" class="error">
      Error: {{ error }}
    </div>

    <div v-else>
      Loading Gopher wisdom...
    </div>
    
    <button @click="fetchQuote">Get New Quote</button>
  </div>
</template>

<style scoped>
.container { text-align: center; font-family: sans-serif; margin-top: 50px; }
.card { background: #f4f4f4; padding: 20px; border-radius: 8px; display: inline-block; }
.error { color: red; }
button { margin-top: 20px; padding: 10px 20px; cursor: pointer; }
</style>