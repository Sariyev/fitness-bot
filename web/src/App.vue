<template>
  <div class="app">
    <div v-if="loading" class="loading">
      <div class="spinner"></div>
    </div>
    <template v-else>
      <nav class="bottom-nav" v-if="showNav">
        <router-link to="/" class="nav-item" :class="{ active: $route.name === 'modules' }">
          <span class="nav-icon">📚</span>
          <span class="nav-label">Модули</span>
        </router-link>
        <router-link to="/payment" class="nav-item" :class="{ active: $route.name === 'payment' }">
          <span class="nav-icon">💳</span>
          <span class="nav-label">Оплата</span>
        </router-link>
        <router-link to="/profile" class="nav-item" :class="{ active: $route.name === 'profile' }">
          <span class="nav-icon">👤</span>
          <span class="nav-label">Профиль</span>
        </router-link>
      </nav>
      <div class="content" :class="{ 'with-nav': showNav }">
        <router-view />
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { api } from './api'

const router = useRouter()
const route = useRoute()
const loading = ref(true)
const paid = ref(false)

const showNav = computed(() => {
  const navRoutes = ['modules', 'payment', 'profile']
  return navRoutes.includes(route.name as string)
})

onMounted(async () => {
  const tg = window.Telegram?.WebApp
  if (tg) {
    tg.ready()
    tg.expand()
  }

  try {
    const status = await api.getSubscriptionStatus()
    paid.value = status.active
  } catch {
    paid.value = false
  } finally {
    loading.value = false
  }

  // Redirect unpaid users to payment page (unless already on payment/profile)
  if (!paid.value && route.name !== 'payment' && route.name !== 'profile') {
    router.replace('/payment')
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
}

.content {
  padding: 16px;
}

.content.with-nav {
  padding-bottom: 72px;
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

.bottom-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  background: var(--bg-color);
  border-top: 1px solid var(--secondary-bg);
  padding: 8px 0;
  padding-bottom: max(8px, env(safe-area-inset-bottom));
  z-index: 100;
}

.nav-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-decoration: none;
  color: var(--hint-color);
  font-size: 11px;
  gap: 2px;
}

.nav-item.active {
  color: var(--button-color);
}

.nav-icon {
  font-size: 20px;
}

.nav-label {
  font-size: 11px;
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
