<template>
  <div class="modules-view">
    <h1 class="page-title">Модули 📚</h1>
    <div v-if="loading" class="skeleton-list">
      <SkeletonCard v-for="i in 3" :key="i" />
    </div>
    <div v-else-if="modules.length === 0" class="empty">
      <p>Модули пока не доступны</p>
    </div>
    <div v-else class="module-grid">
      <div
        v-for="(mod, index) in modules"
        :key="mod.id"
        class="module-card"
        :style="{ animationDelay: index * 60 + 'ms' }"
        @click="openModule(mod.id)"
      >
        <div class="module-accent"></div>
        <div class="module-icon">{{ mod.icon || '📚' }}</div>
        <div class="module-info">
          <h3>{{ mod.name }}</h3>
          <p>{{ mod.description }}</p>
        </div>
        <div class="module-arrow">›</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../api'
import { useTelegram } from '../composables/useTelegram'
import type { Module } from '../types'
import SkeletonCard from '../components/SkeletonCard.vue'

const router = useRouter()
const { hapticSelection } = useTelegram()
const modules = ref<Module[]>([])
const loading = ref(true)

function openModule(id: number) {
  hapticSelection()
  router.push({ name: 'categories', params: { id } })
}

onMounted(async () => {
  window.Telegram?.WebApp?.BackButton?.hide()

  try {
    modules.value = await api.getModules()
  } catch (e) {
    console.error('Failed to load modules:', e)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.page-title {
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 16px;
}

.skeleton-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.module-grid {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.module-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: var(--secondary-bg);
  border-radius: 12px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease forwards;
  transition: transform 0.15s ease;
}

.module-card:active {
  transform: scale(0.98);
}

.module-accent {
  position: absolute;
  top: 0;
  left: 0;
  width: 4px;
  height: 100%;
  background: linear-gradient(180deg, var(--button-color), var(--link-color));
  border-radius: 4px 0 0 4px;
}

.module-icon {
  font-size: 32px;
  flex-shrink: 0;
}

.module-info {
  flex: 1;
  min-width: 0;
}

.module-info h3 {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 4px;
}

.module-info p {
  font-size: 13px;
  color: var(--hint-color);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.module-arrow {
  font-size: 20px;
  color: var(--hint-color);
  flex-shrink: 0;
}

.empty {
  text-align: center;
  padding: 40px 0;
  color: var(--hint-color);
}
</style>
