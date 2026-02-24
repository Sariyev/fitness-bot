<template>
  <div class="app">
    <div v-if="loading" class="loading-screen">
      <div class="loading-brand">💪</div>
      <div class="spinner"></div>
    </div>
    <template v-else>
      <nav class="bottom-nav" v-if="showNav">
        <router-link
          to="/"
          class="nav-item"
          :class="{ active: $route.name === 'modules' }"
          @click="hapticTap"
        >
          <span class="nav-icon">📚</span>
          <span class="nav-label">Модули</span>
        </router-link>
        <router-link
          to="/payment"
          class="nav-item"
          :class="{ active: $route.name === 'payment' }"
          @click="hapticTap"
        >
          <span class="nav-icon">💳</span>
          <span class="nav-label">Оплата</span>
        </router-link>
        <router-link
          to="/profile"
          class="nav-item"
          :class="{ active: $route.name === 'profile' }"
          @click="hapticTap"
        >
          <span class="nav-icon">👤</span>
          <span class="nav-label">Профиль</span>
        </router-link>
      </nav>
      <div class="content" :class="{ 'with-nav': showNav }">
        <router-view v-slot="{ Component }">
          <Transition :name="transitionName" mode="out-in">
            <component :is="Component" />
          </Transition>
        </router-view>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()
const loading = ref(true)
const transitionName = ref('fade')

const showNav = computed(() => {
  if (route.meta.hideNav) return false
  const navRoutes = ['modules', 'payment', 'profile']
  return navRoutes.includes(route.name as string)
})

function hapticTap() {
  window.Telegram?.WebApp?.HapticFeedback?.selectionChanged()
}

// Track route depth for transition direction
const routeDepth: Record<string, number> = {
  onboarding: 0,
  modules: 1,
  payment: 1,
  profile: 1,
  categories: 2,
  lessons: 3,
  lesson: 4,
}

watch(
  () => route.name,
  (to, from) => {
    const toDepth = routeDepth[to as string] ?? 1
    const fromDepth = routeDepth[from as string] ?? 1
    if (toDepth > fromDepth) {
      transitionName.value = 'slide-left'
    } else if (toDepth < fromDepth) {
      transitionName.value = 'slide-right'
    } else {
      transitionName.value = 'fade'
    }
  }
)

onMounted(async () => {
  const tg = window.Telegram?.WebApp
  if (tg) {
    tg.ready()
    tg.expand()
  }

  // Small delay for branded loader feel
  await new Promise(r => setTimeout(r, 300))
  loading.value = false
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

/* Branded loading screen */
.loading-screen {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  gap: 20px;
}

.loading-brand {
  font-size: 48px;
  animation: pulse-brand 1.2s ease-in-out infinite;
}

@keyframes pulse-brand {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.15); }
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

/* Bottom nav */
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
  transition: color 0.2s;
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

/* Page transitions */
.slide-left-enter-active,
.slide-left-leave-active,
.slide-right-enter-active,
.slide-right-leave-active {
  transition: all 0.25s ease;
}

.slide-left-enter-from { transform: translateX(30px); opacity: 0; }
.slide-left-leave-to { transform: translateX(-30px); opacity: 0; }
.slide-right-enter-from { transform: translateX(-30px); opacity: 0; }
.slide-right-leave-to { transform: translateX(30px); opacity: 0; }

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* Stagger animation utility */
@keyframes fadeSlideUp {
  from { opacity: 0; transform: translateY(12px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Skeleton shimmer */
@keyframes shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}
</style>
