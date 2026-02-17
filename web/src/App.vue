<template>
  <div class="app">
    <div v-if="loading" class="loading">
      <div class="spinner"></div>
    </div>
    <div v-else-if="!subscribed" class="no-subscription">
      <div class="no-sub-content">
        <h2>Требуется подписка</h2>
        <p>Для доступа к модулям необходима активная подписка.</p>
        <p>Используйте команду <strong>/subscribe</strong> в боте.</p>
        <button class="btn" @click="closeApp">Закрыть</button>
      </div>
    </div>
    <router-view v-else />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from './api'

const loading = ref(true)
const subscribed = ref(false)

function closeApp() {
  window.Telegram?.WebApp?.close()
}

onMounted(async () => {
  const tg = window.Telegram?.WebApp
  if (tg) {
    tg.ready()
    tg.expand()
  }

  try {
    const status = await api.getSubscriptionStatus()
    subscribed.value = status.active
  } catch {
    subscribed.value = false
  } finally {
    loading.value = false
  }
})
</script>

<style>
:root {
  --bg-color: var(--tg-theme-bg-color, #ffffff);
  --text-color: var(--tg-theme-text-color, #000000);
  --hint-color: var(--tg-theme-hint-color, #999999);
  --link-color: var(--tg-theme-link-color, #2481cc);
  --button-color: var(--tg-theme-button-color, #2481cc);
  --button-text-color: var(--tg-theme-button-text-color, #ffffff);
  --secondary-bg: var(--tg-theme-secondary-bg-color, #f0f0f0);
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  background-color: var(--bg-color);
  color: var(--text-color);
  -webkit-font-smoothing: antialiased;
}

.app {
  min-height: 100vh;
  padding: 16px;
}

.loading {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 60vh;
}

.spinner {
  width: 32px;
  height: 32px;
  border: 3px solid var(--secondary-bg);
  border-top-color: var(--button-color);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.no-subscription {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 60vh;
}

.no-sub-content {
  text-align: center;
  padding: 24px;
}

.no-sub-content h2 {
  margin-bottom: 12px;
}

.no-sub-content p {
  color: var(--hint-color);
  margin-bottom: 8px;
}

.btn {
  display: inline-block;
  padding: 12px 24px;
  margin-top: 16px;
  background: var(--button-color);
  color: var(--button-text-color);
  border: none;
  border-radius: 8px;
  font-size: 16px;
  cursor: pointer;
}

.btn:active {
  opacity: 0.8;
}
</style>
