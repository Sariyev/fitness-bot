<template>
  <div class="admin-page">
    <h1 class="page-title">Контент</h1>

    <div v-if="loading" class="loading">Загрузка...</div>

    <div v-else>
      <div v-if="loadError" class="load-error">
        ⚠️ Не удалось загрузить часть данных. Подробности в консоли.
      </div>
      <div class="content-card" @click="router.push('/admin/pricing')">
        <div class="content-main">
          <span class="content-name">💰 Цены категорий</span>
          <span class="content-meta">Тренировки · ЛФК · Питание</span>
        </div>
        <span class="arrow">→</span>
      </div>

      <div class="section-header">
        <h2 class="section-title">Тренировки ({{ workouts.length }})</h2>
        <button class="add-btn" @click="router.push('/admin/workouts/new')">+</button>
      </div>
      <div v-for="w in workouts" :key="w.id" class="content-card" :class="{ inactive: !w.is_active }" @click="router.push(`/admin/workouts/${w.id}`)">
        <div class="content-main">
          <span class="content-name">{{ w.name }}</span>
          <span class="content-meta">{{ w.level }} | {{ w.duration_minutes }} мин</span>
        </div>
        <span class="badge" :class="w.is_active ? 'badge-active' : 'badge-inactive'">
          {{ w.is_active ? 'Active' : 'Off' }}
        </span>
      </div>

      <div class="section-header">
        <h2 class="section-title">Планы питания ({{ mealPlans.length }})</h2>
        <button class="add-btn" @click="router.push('/admin/meal-plans/new')">+</button>
      </div>
      <div v-for="mp in mealPlans" :key="mp.id" class="content-card" :class="{ inactive: !mp.is_active }" @click="router.push(`/admin/meal-plans/${mp.id}`)">
        <div class="content-main">
          <span class="content-name">{{ mp.name }}</span>
          <span class="content-meta">{{ mp.goal }} | День {{ mp.day_number }} | {{ mp.calories }} ккал</span>
        </div>
        <span class="badge" :class="mp.is_active ? 'badge-active' : 'badge-inactive'">
          {{ mp.is_active ? 'Active' : 'Off' }}
        </span>
      </div>

      <div class="section-header">
        <h2 class="section-title">ЛФК курсы ({{ rehabCourses.length }})</h2>
        <button class="add-btn" @click="router.push('/admin/rehab/courses/new')">+</button>
      </div>
      <div v-for="rc in rehabCourses" :key="rc.id" class="content-card" :class="{ inactive: !rc.is_active }" @click="router.push(`/admin/rehab/courses/${rc.id}`)">
        <div class="content-main">
          <span class="content-name">{{ rc.name }}</span>
          <span class="content-meta">{{ rc.category }}</span>
        </div>
        <span class="badge" :class="rc.is_active ? 'badge-active' : 'badge-inactive'">
          {{ rc.is_active ? 'Active' : 'Off' }}
        </span>
      </div>

      <div class="section-header">
        <h2 class="section-title">Упражнения</h2>
      </div>
      <div class="content-card" @click="router.push('/admin/exercises')">
        <div class="content-main">
          <span class="content-name">Управление упражнениями</span>
          <span class="content-meta">Просмотр и редактирование базы упражнений</span>
        </div>
        <span class="arrow">→</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../api'
import type { Workout, MealPlan, RehabCourse } from '../../types'

const router = useRouter()
const workouts = ref<Workout[]>([])
const mealPlans = ref<MealPlan[]>([])
const rehabCourses = ref<RehabCourse[]>([])
const loading = ref(true)
const loadError = ref(false)

onMounted(async () => {
  const results = await Promise.allSettled([
    api.getAdminWorkouts(),
    api.getAdminMealPlans(),
    api.getAdminRehabCourses(),
  ])
  const labels = ['workouts', 'meal-plans', 'rehab-courses']
  const refs = [workouts, mealPlans, rehabCourses] as const
  results.forEach((r, i) => {
    if (r.status === 'fulfilled') {
      ;(refs[i].value as unknown[]) = (r.value as unknown[]) || []
    } else {
      loadError.value = true
      console.error(`AdminContent: failed to load ${labels[i]}`, r.reason)
    }
  })
  loading.value = false
})
</script>

<style scoped>
.admin-page {
  max-width: 400px;
  margin: 0 auto;
  padding-bottom: 24px;
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

.load-error {
  background: rgba(255, 59, 48, 0.1);
  color: #ff3b30;
  border-radius: 12px;
  padding: 10px 14px;
  margin-bottom: 12px;
  font-size: 13px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin: 16px 0 10px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--hint-color);
}

.add-btn {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: none;
  background: var(--button-color);
  color: var(--button-text-color);
  font-size: 20px;
  font-weight: 700;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.content-card {
  background: var(--secondary-bg);
  border-radius: 12px;
  padding: 12px 16px;
  margin-bottom: 8px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  cursor: pointer;
  animation: fadeSlideUp 0.3s ease both;
  touch-action: manipulation;
}

.content-card.inactive {
  animation: fadeSlideUp 0.3s ease both;
}
.content-card.inactive .content-name {
  color: var(--hint-color);
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
  flex-shrink: 0;
}

.badge-active {
  background: #34c759;
  color: #fff;
}

.badge-inactive {
  background: rgba(0,0,0,0.08);
  color: var(--hint-color);
}

.arrow {
  font-size: 18px;
  color: var(--hint-color);
}

@keyframes fadeSlideUp {
  from { opacity: 0; transform: translateY(12px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
