<template>
  <div class="admin-page">
    <button class="back-btn" @click="router.push('/admin')">← Назад</button>
    <h1 class="page-title">Контент</h1>

    <div v-if="loading" class="loading">Загрузка...</div>

    <div v-else>
      <h2 class="section-title">Программы ({{ programs.length }})</h2>
      <div v-for="p in programs" :key="p.id" class="content-card">
        <div class="content-main">
          <span class="content-name">{{ p.name }}</span>
          <span class="content-meta">{{ p.level }} | {{ p.format }} | {{ p.duration_weeks }} нед.</span>
        </div>
        <span class="badge" :class="p.is_active ? 'badge-active' : 'badge-inactive'">
          {{ p.is_active ? 'Active' : 'Off' }}
        </span>
      </div>

      <h2 class="section-title">Тренировки ({{ workouts.length }})</h2>
      <div v-for="w in workouts" :key="w.id" class="content-card">
        <div class="content-main">
          <span class="content-name">{{ w.name }}</span>
          <span class="content-meta">{{ w.level }} | {{ w.duration_minutes }} мин</span>
        </div>
        <span class="badge" :class="w.is_active ? 'badge-active' : 'badge-inactive'">
          {{ w.is_active ? 'Active' : 'Off' }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../api'
import type { Program, Workout } from '../../types'

const router = useRouter()
const programs = ref<Program[]>([])
const workouts = ref<Workout[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    const [p, w] = await Promise.all([
      api.getAdminPrograms(),
      api.getAdminWorkouts(),
    ])
    programs.value = p || []
    workouts.value = w || []
  } catch {
    // ignore
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.admin-page {
  max-width: 400px;
  margin: 0 auto;
  padding-bottom: 24px;
}

.back-btn {
  background: none;
  border: none;
  color: var(--button-color);
  font-size: 16px;
  cursor: pointer;
  padding: 4px 0;
  margin-bottom: 12px;
}

.page-title {
  font-size: 20px;
  font-weight: 700;
  margin-bottom: 16px;
}

.loading {
  text-align: center;
  color: var(--hint-color);
  padding: 40px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  margin: 16px 0 10px;
  color: var(--hint-color);
}

.content-card {
  background: var(--secondary-bg);
  border-radius: 12px;
  padding: 12px 16px;
  margin-bottom: 8px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  opacity: 0;
  animation: fadeSlideUp 0.3s ease forwards;
}

.content-main {
  display: flex;
  flex-direction: column;
}

.content-name {
  font-weight: 500;
  font-size: 15px;
}

.content-meta {
  color: var(--hint-color);
  font-size: 12px;
  margin-top: 2px;
}

.badge {
  padding: 3px 10px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.badge-active {
  background: #34c759;
  color: #fff;
}

.badge-inactive {
  background: rgba(0,0,0,0.08);
  color: var(--hint-color);
}

@keyframes fadeSlideUp {
  from { opacity: 0; transform: translateY(12px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
