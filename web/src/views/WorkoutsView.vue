<template>
  <div class="workouts-page">
    <!-- Format toggle: Дом / Зал -->
    <div class="format-toggle">
      <button
        class="toggle-btn"
        :class="{ active: filters.format === 'home' }"
        @click="setFormat('home')"
      >Дом</button>
      <button
        class="toggle-btn"
        :class="{ active: filters.format === 'gym' }"
        @click="setFormat('gym')"
      >Зал</button>
    </div>

    <!-- Goal filter chips -->
    <div class="filter-row">
      <button
        class="filter-chip"
        :class="{ active: filters.goal === 'weight_loss' }"
        @click="toggleGoal('weight_loss')"
      >Похудение</button>
      <button
        class="filter-chip"
        :class="{ active: filters.goal === 'muscle_gain' }"
        @click="toggleGoal('muscle_gain')"
      >Набор мышц</button>
    </div>

    <!-- Level filter chips -->
    <div class="filter-row">
      <button
        class="filter-chip"
        :class="{ active: filters.level === 'beginner' }"
        @click="toggleLevel('beginner')"
      >Новичок</button>
      <button
        class="filter-chip"
        :class="{ active: filters.level === 'intermediate' }"
        @click="toggleLevel('intermediate')"
      >Средний</button>
      <button
        class="filter-chip"
        :class="{ active: filters.level === 'advanced' }"
        @click="toggleLevel('advanced')"
      >Продвинутый</button>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="skeleton-list">
      <div class="skeleton-card" style="height: 100px" v-for="i in 2" :key="'p-' + i"></div>
      <div class="skeleton-card" style="height: 80px" v-for="i in 4" :key="'w-' + i"></div>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="error-state">
      <p class="error-text">{{ error }}</p>
      <button class="btn-primary" @click="loadData">Попробовать снова</button>
    </div>

    <template v-else>
      <!-- Programs section -->
      <section v-if="programs.length > 0" class="section">
        <h2 class="section-title">Программы</h2>
        <div class="programs-list">
          <router-link
            v-for="(program, index) in programs"
            :key="program.id"
            :to="`/workouts/program/${program.id}`"
            class="program-card"
            :style="{ animationDelay: (index * 80) + 'ms' }"
          >
            <div class="program-card-top">
              <span class="program-name">{{ program.name }}</span>
              <span class="level-badge">{{ levelLabel(program.level) }}</span>
            </div>
            <div class="program-meta">
              <span v-if="program.duration_weeks" class="meta-item">
                {{ program.duration_weeks }} {{ weeksLabel(program.duration_weeks) }}
              </span>
              <span v-if="program.goal" class="meta-item">
                {{ goalLabel(program.goal) }}
              </span>
            </div>
          </router-link>
        </div>
      </section>

      <!-- Workouts library section -->
      <section class="section">
        <h2 class="section-title">Библиотека</h2>

        <div v-if="workouts.length > 0" class="workouts-list">
          <router-link
            v-for="(workout, index) in workouts"
            :key="workout.id"
            :to="`/workouts/session/${workout.id}`"
            class="workout-card"
            :style="{ animationDelay: (index * 60) + 'ms' }"
          >
            <div class="workout-card-body">
              <span class="workout-name">{{ workout.name }}</span>
              <div class="workout-meta">
                <span v-if="workout.duration_minutes" class="meta-tag">{{ workout.duration_minutes }} мин</span>
                <span
                  v-for="eq in (workout.equipment || []).slice(0, 3)"
                  :key="eq"
                  class="meta-tag equipment-tag"
                >{{ eq }}</span>
              </div>
            </div>
            <span class="workout-arrow">&rsaquo;</span>
          </router-link>
        </div>

        <div v-else class="empty-state">
          <p class="empty-text">Тренировок не найдено</p>
          <p class="empty-hint">Попробуйте изменить фильтры</p>
        </div>
      </section>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch, onMounted } from 'vue'
import { api } from '../api'
import type { Program, Workout } from '../types'

const loading = ref(true)
const error = ref('')
const programs = ref<Program[]>([])
const workouts = ref<Workout[]>([])

const filters = reactive({
  format: 'home' as string,
  goal: '' as string,
  level: '' as string,
})

const goalLabels: Record<string, string> = {
  weight_loss: 'Похудение',
  muscle_gain: 'Набор мышц',
  strength: 'Сила',
  endurance: 'Выносливость',
  maintenance: 'Поддержание',
}

const levelLabels: Record<string, string> = {
  beginner: 'Новичок',
  intermediate: 'Средний',
  advanced: 'Продвинутый',
}

function goalLabel(key: string): string {
  return goalLabels[key] || key
}

function levelLabel(key: string): string {
  return levelLabels[key] || key
}

function weeksLabel(n: number): string {
  if (n % 10 === 1 && n % 100 !== 11) return 'неделя'
  if (n % 10 >= 2 && n % 10 <= 4 && (n % 100 < 10 || n % 100 >= 20)) return 'недели'
  return 'недель'
}

function setFormat(format: string) {
  filters.format = format
}

function toggleGoal(goal: string) {
  filters.goal = filters.goal === goal ? '' : goal
}

function toggleLevel(level: string) {
  filters.level = filters.level === level ? '' : level
}

function buildFilters(): { format?: string; goal?: string; level?: string } {
  const result: { format?: string; goal?: string; level?: string } = {}
  if (filters.format) result.format = filters.format
  if (filters.goal) result.goal = filters.goal
  if (filters.level) result.level = filters.level
  return result
}

async function loadData() {
  loading.value = true
  error.value = ''
  try {
    const filterParams = buildFilters()
    const [programsData, workoutsData] = await Promise.all([
      api.getPrograms(filterParams),
      api.getWorkouts(filterParams),
    ])
    programs.value = programsData
    workouts.value = workoutsData
  } catch (e: any) {
    error.value = e.message || 'Ошибка загрузки'
    programs.value = []
    workouts.value = []
  } finally {
    loading.value = false
  }
}

watch(filters, () => {
  loadData()
})

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.workouts-page {
  max-width: 480px;
  margin: 0 auto;
  padding-bottom: 24px;
}

/* ===== Format toggle ===== */
.format-toggle {
  display: flex;
  background: var(--secondary-bg);
  border-radius: 12px;
  padding: 4px;
  margin-bottom: 12px;
}

.toggle-btn {
  flex: 1;
  padding: 10px 16px;
  border: none;
  border-radius: 10px;
  background: transparent;
  color: var(--hint-color);
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  font-family: inherit;
}

.toggle-btn.active {
  background: var(--button-color);
  color: var(--button-text-color);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

/* ===== Filter chips ===== */
.filter-row {
  display: flex;
  gap: 6px;
  overflow-x: auto;
  padding-bottom: 8px;
  scrollbar-width: none;
  -ms-overflow-style: none;
}

.filter-row::-webkit-scrollbar {
  display: none;
}

.filter-chip {
  flex-shrink: 0;
  padding: 6px 14px;
  border: 1.5px solid var(--secondary-bg);
  border-radius: 20px;
  background: var(--bg-color);
  color: var(--hint-color);
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;
  font-family: inherit;
}

.filter-chip.active {
  background: var(--button-color);
  color: var(--button-text-color);
  border-color: var(--button-color);
}

/* ===== Skeleton ===== */
.skeleton-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-top: 8px;
}

.skeleton-card {
  background: var(--secondary-bg);
  border-radius: 12px;
  background: linear-gradient(
    90deg,
    var(--secondary-bg) 25%,
    color-mix(in srgb, var(--hint-color) 15%, transparent) 50%,
    var(--secondary-bg) 75%
  );
  background-size: 200% 100%;
  animation: shimmer 1.5s ease-in-out infinite;
}

@keyframes shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

/* ===== Error ===== */
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 40px 16px;
  text-align: center;
}

.error-text {
  font-size: 15px;
  color: var(--hint-color);
  margin-bottom: 16px;
}

/* ===== Section ===== */
.section {
  margin-top: 20px;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease forwards;
}

.section-title {
  font-size: 18px;
  font-weight: 700;
  margin-bottom: 12px;
}

/* ===== Programs ===== */
.programs-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.program-card {
  background: var(--secondary-bg);
  border-radius: 12px;
  padding: 16px;
  text-decoration: none;
  color: var(--text-color);
  transition: transform 0.15s ease;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease forwards;
}

.program-card:active {
  transform: scale(0.98);
}

.program-card-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 8px;
}

.program-name {
  font-size: 15px;
  font-weight: 600;
  flex: 1;
  min-width: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.level-badge {
  flex-shrink: 0;
  background: var(--button-color);
  color: var(--button-text-color);
  padding: 3px 10px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 600;
}

.program-meta {
  display: flex;
  gap: 12px;
}

.meta-item {
  font-size: 13px;
  color: var(--hint-color);
}

/* ===== Workouts library ===== */
.workouts-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.workout-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  background: var(--secondary-bg);
  border-radius: 12px;
  text-decoration: none;
  color: var(--text-color);
  transition: transform 0.15s ease;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease forwards;
}

.workout-card:active {
  transform: scale(0.98);
}

.workout-card-body {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.workout-name {
  font-size: 15px;
  font-weight: 600;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.workout-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.meta-tag {
  padding: 2px 8px;
  border-radius: 8px;
  font-size: 12px;
  color: var(--hint-color);
  background: var(--bg-color);
}

.equipment-tag {
  color: var(--text-color);
}

.workout-arrow {
  flex-shrink: 0;
  font-size: 22px;
  color: var(--hint-color);
}

/* ===== Empty ===== */
.empty-state {
  text-align: center;
  padding: 40px 16px;
}

.empty-text {
  font-size: 16px;
  font-weight: 600;
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
