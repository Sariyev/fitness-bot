<template>
  <div class="app">
    <div v-if="loading" class="loading-screen">
      <div class="loading-brand">💪</div>
      <div class="spinner"></div>
    </div>
    <template v-else>
      <!-- Top header with profile icon -->
      <header class="app-header" v-if="showNav">
        <span class="header-title">{{ pageTitle }}</span>
        <router-link to="/profile" class="header-profile" @click="hapticTap">👤</router-link>
      </header>

      <div class="content" :class="{ 'with-nav': showNav, 'with-header': showNav }">
        <router-view v-slot="{ Component }">
          <Transition :name="transitionName" mode="out-in">
            <component :is="Component" />
          </Transition>
        </router-view>
      </div>

      <!-- 5-tab bottom nav -->
      <nav class="bottom-nav" v-if="showNav">
        <router-link to="/" class="nav-item" :class="{ active: $route.name === 'today' }" @click="hapticTap">
          <span class="nav-icon">🏠</span>
          <span class="nav-label">Сегодня</span>
        </router-link>
        <router-link to="/workouts" class="nav-item" :class="{ active: $route.name === 'workouts' }" @click="hapticTap">
          <span class="nav-icon">🏋️</span>
          <span class="nav-label">Тренировки</span>
        </router-link>
        <router-link to="/lfk" class="nav-item" :class="{ active: $route.name === 'lfk' }" @click="hapticTap">
          <span class="nav-icon">❤️‍🩹</span>
          <span class="nav-label">ЛФК</span>
        </router-link>
        <router-link to="/nutrition" class="nav-item" :class="{ active: $route.name === 'nutrition' }" @click="hapticTap">
          <span class="nav-icon">🍽️</span>
          <span class="nav-label">Питание</span>
        </router-link>
        <router-link to="/progress" class="nav-item" :class="{ active: $route.name === 'progress' }" @click="hapticTap">
          <span class="nav-icon">📊</span>
          <span class="nav-label">Прогресс</span>
        </router-link>
      </nav>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const loading = ref(true)
const transitionName = ref('fade')

const mainTabs = ['today', 'workouts', 'lfk', 'nutrition', 'progress']

const showNav = computed(() => {
  if (route.meta.hideNav) return false
  return mainTabs.includes(route.name as string)
})

const pageTitle = computed(() => {
  const titles: Record<string, string> = {
    today: 'Сегодня',
    workouts: 'Тренировки',
    lfk: 'ЛФК 14 дней',
    nutrition: 'Питание',
    progress: 'Прогресс',
  }
  return titles[route.name as string] || ''
})

function hapticTap() {
  window.Telegram?.WebApp?.HapticFeedback?.selectionChanged()
}

const routeDepth: Record<string, number> = {
  onboarding: 0,
  today: 1,
  workouts: 1,
  lfk: 1,
  nutrition: 1,
  progress: 1,
  profile: 2,
  payment: 2,
  'workout-program': 2,
  'workout-session': 3,
  'lfk-course': 2,
  'lfk-session': 3,
  'food-diary': 2,
  'macro-calculator': 2,
  modules: 1,
  categories: 2,
  lessons: 3,
  lesson: 4,
  admin: 2,
  'admin-users': 3,
  'admin-user-detail': 4,
  'admin-content': 3,
  'admin-reviews': 3,
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

* { margin: 0; padding: 0; box-sizing: border-box; }

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  background-color: var(--bg-color);
  color: var(--text-color);
  -webkit-font-smoothing: antialiased;
}

.app { min-height: 100vh; }
.content { padding: 16px; }
.content.with-nav { padding-bottom: calc(80px + env(safe-area-inset-bottom, 0px)); }
.content.with-header { padding-top: 56px; }

/* App header */
.app-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  background: var(--bg-color);
  border-bottom: 1px solid var(--secondary-bg);
  z-index: 100;
}

.header-title {
  font-size: 18px;
  font-weight: 700;
}

.header-profile {
  font-size: 22px;
  text-decoration: none;
}

/* Loading */
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

@keyframes spin { to { transform: rotate(360deg); } }

/* Bottom nav - 5 tabs */
.bottom-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  background: var(--bg-color);
  border-top: 1px solid var(--secondary-bg);
  padding: 8px 0;
  padding-bottom: calc(8px + env(safe-area-inset-bottom, 0px));
  z-index: 100;
}

.nav-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-decoration: none;
  color: var(--hint-color);
  font-size: 10px;
  gap: 2px;
  transition: color 0.2s;
}

.nav-item.active { color: var(--button-color); }
.nav-icon { font-size: 18px; }
.nav-label { font-size: 10px; }

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
  width: 100%;
  text-align: center;
}

.btn:active { opacity: 0.8; }

/* Transitions */
.slide-left-enter-active, .slide-left-leave-active,
.slide-right-enter-active, .slide-right-leave-active {
  transition: all 0.25s ease;
}
.slide-left-enter-from { transform: translateX(30px); opacity: 0; }
.slide-left-leave-to { transform: translateX(-30px); opacity: 0; }
.slide-right-enter-from { transform: translateX(-30px); opacity: 0; }
.slide-right-leave-to { transform: translateX(30px); opacity: 0; }

.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }

@keyframes fadeSlideUp {
  from { opacity: 0; transform: translateY(12px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}
</style>
