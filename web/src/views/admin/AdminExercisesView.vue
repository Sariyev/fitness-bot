<template>
  <div class="admin-page">
    <button class="back-btn" @click="router.push('/admin/content')">← Назад</button>
    <div class="section-header">
      <h1 class="page-title">Упражнения</h1>
      <button class="add-btn" @click="router.push('/admin/exercises/new')">+</button>
    </div>

    <div v-if="loading" class="loading">Загрузка...</div>

    <div v-else>
      <div v-for="e in exercises" :key="e.id" class="content-card" @click="router.push(`/admin/exercises/${e.id}`)">
        <div class="content-main">
          <span class="content-name">{{ e.name }}</span>
          <span class="content-meta">Отдых: {{ e.rest_seconds }}с</span>
        </div>
        <span class="arrow">→</span>
      </div>
      <div v-if="exercises.length === 0" class="empty">Нет упражнений</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../api'
import type { Exercise } from '../../types'

const router = useRouter()
const exercises = ref<Exercise[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    exercises.value = (await api.getAdminExercises()) || []
  } catch {
    // ignore
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.admin-page { max-width: 400px; margin: 0 auto; padding-bottom: 24px; }
.back-btn { background: none; border: none; color: var(--button-color); font-size: 16px; cursor: pointer; padding: 4px 0; margin-bottom: 12px; }
.page-title { font-size: 20px; font-weight: 700; }
.section-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.loading { text-align: center; color: var(--hint-color); padding: 40px; }
.empty { text-align: center; color: var(--hint-color); padding: 40px; }
.add-btn {
  width: 32px; height: 32px; border-radius: 50%; border: none;
  background: var(--button-color); color: var(--button-text-color);
  font-size: 20px; font-weight: 700; cursor: pointer;
  display: flex; align-items: center; justify-content: center;
}
.content-card {
  background: var(--secondary-bg); border-radius: 12px; padding: 12px 16px; margin-bottom: 8px;
  display: flex; justify-content: space-between; align-items: center; cursor: pointer;
  opacity: 0; animation: fadeSlideUp 0.3s ease forwards;
}
.content-main { display: flex; flex-direction: column; }
.content-name { font-weight: 500; font-size: 15px; }
.content-meta { color: var(--hint-color); font-size: 12px; margin-top: 2px; }
.arrow { font-size: 18px; color: var(--hint-color); }
@keyframes fadeSlideUp {
  from { opacity: 0; transform: translateY(12px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
