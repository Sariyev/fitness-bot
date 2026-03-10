<template>
  <div class="today-page">
    <!-- Loading skeleton -->
    <div v-if="loading" class="skeleton-wrapper">
      <div class="skeleton-card skeleton-greeting-card">
        <div class="skeleton-line" style="height: 22px; width: 55%"></div>
        <div class="skeleton-line" style="height: 14px; width: 85%; margin-top: 10px"></div>
        <div class="skeleton-line" style="height: 14px; width: 65%; margin-top: 6px"></div>
      </div>
      <div class="skeleton-chips-row">
        <div class="skeleton-line skeleton-chip-ph" v-for="i in 3" :key="i"></div>
      </div>
      <div class="skeleton-card" v-for="i in 3" :key="'plan-' + i">
        <div class="skeleton-line" style="height: 12px; width: 30%"></div>
        <div class="skeleton-line" style="height: 16px; width: 60%; margin-top: 8px"></div>
        <div class="skeleton-line" style="height: 44px; width: 100%; margin-top: 12px; border-radius: 10px"></div>
      </div>
    </div>

    <!-- Error state -->
    <div v-else-if="error" class="error-state">
      <p class="error-text">Не удалось загрузить данные</p>
      <p class="error-detail">{{ error }}</p>
      <button class="btn-primary" @click="loadDashboard">Попробовать снова</button>
    </div>

    <!-- Dashboard content -->
    <template v-else-if="dashboard">
      <!-- Greeting + trainer message card -->
      <div class="greeting-card">
        <h2 class="greeting-text">{{ dashboard.greeting }}</h2>
        <p v-if="dashboard.trainer_message" class="trainer-message">{{ dashboard.trainer_message }}</p>
      </div>

      <!-- Streak counter -->
      <div v-if="dashboard.current_streak > 0" class="streak-row">
        <span class="streak-fire">&#x1F525;</span>
        <span class="streak-count">{{ dashboard.current_streak }}</span>
        <span class="streak-label">{{ streakDaysLabel }} подряд</span>
      </div>

      <!-- Goals chips -->
      <div v-if="dashboard.goals && dashboard.goals.length" class="goals-row">
        <span
          v-for="goal in dashboard.goals"
          :key="goal"
          class="goal-chip"
        >{{ goalLabel(goal) }}</span>
      </div>

      <!-- Section: today plan -->
      <h3 class="section-title">План на сегодня</h3>

      <div class="plan-cards">
        <!-- Workout card -->
        <div v-if="dashboard.today_workout" class="plan-card" :class="{ done: dashboard.today_workout.done }">
          <div class="plan-card-top">
            <span class="plan-badge badge-workout">Тренировка</span>
            <span v-if="dashboard.today_workout.done" class="done-check">&#x2705;</span>
          </div>
          <h4 class="plan-card-title">{{ dashboard.today_workout.title }}</h4>
          <router-link
            :to="`/workouts/session/${dashboard.today_workout.id}`"
            class="plan-btn"
            :class="{ 'plan-btn-done': dashboard.today_workout.done }"
          >
            {{ dashboard.today_workout.done ? 'Выполнено' : 'Начать' }}
          </router-link>
        </div>

        <!-- Meal card -->
        <div v-if="dashboard.today_meal" class="plan-card" :class="{ done: dashboard.today_meal.done }">
          <div class="plan-card-top">
            <span class="plan-badge badge-meal">Питание</span>
            <span v-if="dashboard.today_meal.done" class="done-check">&#x2705;</span>
          </div>
          <h4 class="plan-card-title">{{ dashboard.today_meal.title }}</h4>
          <router-link
            to="/nutrition"
            class="plan-btn"
            :class="{ 'plan-btn-done': dashboard.today_meal.done }"
          >
            {{ dashboard.today_meal.done ? 'Выполнено' : 'Начать' }}
          </router-link>
        </div>

        <!-- Rehab session card -->
        <div v-if="dashboard.today_rehab" class="plan-card" :class="{ done: dashboard.today_rehab.done }">
          <div class="plan-card-top">
            <span class="plan-badge badge-rehab">ЛФК</span>
            <span v-if="dashboard.today_rehab.done" class="done-check">&#x2705;</span>
          </div>
          <h4 class="plan-card-title">{{ dashboard.today_rehab.title }}</h4>
          <router-link
            :to="`/lfk/session/${dashboard.today_rehab.id}`"
            class="plan-btn"
            :class="{ 'plan-btn-done': dashboard.today_rehab.done }"
          >
            {{ dashboard.today_rehab.done ? 'Выполнено' : 'Начать' }}
          </router-link>
        </div>

        <!-- Empty state -->
        <div v-if="!dashboard.today_workout && !dashboard.today_meal && !dashboard.today_rehab" class="empty-plan">
          <p class="empty-text">На сегодня нет запланированных занятий.</p>
          <p class="empty-hint">Выбери программу или курс ЛФК, чтобы начать!</p>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { api } from '../api'
import type { DashboardData } from '../types'

const loading = ref(true)
const error = ref('')
const dashboard = ref<DashboardData | null>(null)

const goalLabels: Record<string, string> = {
  weight_loss: 'Похудение',
  muscle_gain: 'Набор массы',
  strength: 'Сила',
  endurance: 'Выносливость',
  maintenance: 'Поддержание',
  hernia: 'Грыжа',
  protrusion: 'Протрузии',
  scoliosis: 'Сколиоз',
  kyphosis: 'Кифоз',
  lordosis: 'Лордоз',
}

function goalLabel(key: string): string {
  return goalLabels[key] || key
}

const streakDaysLabel = computed(() => {
  if (!dashboard.value) return ''
  const n = dashboard.value.current_streak
  if (n % 10 === 1 && n % 100 !== 11) return 'день'
  if (n % 10 >= 2 && n % 10 <= 4 && (n % 100 < 10 || n % 100 >= 20)) return 'дня'
  return 'дней'
})

async function loadDashboard() {
  loading.value = true
  error.value = ''
  try {
    dashboard.value = await api.getDashboard()
  } catch (e: any) {
    error.value = e.message || 'Ошибка загрузки'
    dashboard.value = null
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadDashboard()
})
</script>

<style scoped>
.today-page {
  max-width: 480px;
  margin: 0 auto;
  padding-bottom: 24px;
}

/* ===== Skeleton loading ===== */
.skeleton-wrapper {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.skeleton-card {
  background: var(--secondary-bg);
  border-radius: 12px;
  padding: 16px;
}

.skeleton-line {
  background: linear-gradient(
    90deg,
    var(--secondary-bg) 25%,
    color-mix(in srgb, var(--hint-color) 15%, transparent) 50%,
    var(--secondary-bg) 75%
  );
  background-size: 200% 100%;
  animation: shimmer 1.5s ease-in-out infinite;
  border-radius: 6px;
}

.skeleton-chips-row {
  display: flex;
  gap: 8px;
}

.skeleton-chip-ph {
  height: 30px;
  width: 80px;
  border-radius: 16px;
}

@keyframes shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

/* ===== Error state ===== */
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 48px 16px;
  text-align: center;
}

.error-text {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 4px;
}

.error-detail {
  font-size: 14px;
  color: var(--hint-color);
  margin-bottom: 20px;
}

/* ===== Greeting card ===== */
.greeting-card {
  background: var(--secondary-bg);
  border-radius: 12px;
  padding: 20px 16px;
  margin-bottom: 16px;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease forwards;
}

.greeting-text {
  font-size: 22px;
  font-weight: 700;
  margin-bottom: 8px;
}

.trainer-message {
  font-size: 14px;
  color: var(--hint-color);
  line-height: 1.5;
}

/* ===== Streak ===== */
.streak-row {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 16px;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease forwards;
  animation-delay: 80ms;
}

.streak-fire {
  font-size: 22px;
}

.streak-count {
  font-size: 26px;
  font-weight: 800;
  color: var(--button-color);
}

.streak-label {
  font-size: 14px;
  color: var(--hint-color);
}

/* ===== Goals ===== */
.goals-row {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 20px;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease forwards;
  animation-delay: 120ms;
}

.goal-chip {
  background: var(--button-color);
  color: var(--button-text-color);
  padding: 5px 14px;
  border-radius: 16px;
  font-size: 13px;
  font-weight: 500;
}

/* ===== Section title ===== */
.section-title {
  font-size: 18px;
  font-weight: 700;
  margin-bottom: 12px;
}

/* ===== Plan cards ===== */
.plan-cards {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.plan-card {
  background: var(--secondary-bg);
  border-radius: 12px;
  padding: 16px;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease forwards;
  animation-delay: 160ms;
}

.plan-card:nth-child(2) { animation-delay: 220ms; }
.plan-card:nth-child(3) { animation-delay: 280ms; }

.plan-card.done {
  opacity: 0.65;
}

.plan-card-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.plan-badge {
  display: inline-block;
  padding: 3px 10px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
  color: #fff;
}

.badge-workout { background: #5856d6; }
.badge-meal { background: #34c759; }
.badge-rehab { background: #ff9500; }

.done-check {
  font-size: 16px;
}

.plan-card-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 12px;
}

.plan-btn {
  display: block;
  width: 100%;
  padding: 14px;
  background: var(--button-color);
  color: var(--button-text-color);
  border: none;
  border-radius: 10px;
  font-size: 16px;
  font-weight: 600;
  text-align: center;
  text-decoration: none;
  cursor: pointer;
}

.plan-btn:active {
  opacity: 0.85;
}

.plan-btn-done {
  background: var(--secondary-bg);
  color: var(--hint-color);
  border: 1px solid var(--hint-color);
}

/* ===== Empty ===== */
.empty-plan {
  background: var(--secondary-bg);
  border-radius: 12px;
  padding: 28px 16px;
  text-align: center;
}

.empty-text {
  font-size: 15px;
  margin-bottom: 4px;
}

.empty-hint {
  font-size: 13px;
  color: var(--hint-color);
}

/* ===== Buttons ===== */
.btn-primary {
  display: block;
  width: 100%;
  padding: 14px;
  background: var(--button-color);
  color: var(--button-text-color);
  border: none;
  border-radius: 12px;
  font-size: 16px;
  cursor: pointer;
  text-align: center;
  font-family: inherit;
}

@keyframes fadeSlideUp {
  from { opacity: 0; transform: translateY(12px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
