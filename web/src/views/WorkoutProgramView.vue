<template>
  <div class="program-page">
    <button class="back-btn" @click="router.back()">← Назад</button>

    <div v-if="loading" class="skeleton-list">
      <SkeletonCard v-for="i in 4" :key="i" />
    </div>

    <template v-else-if="program">
      <div class="program-header">
        <h1 class="program-name">{{ program.name }}</h1>
        <p class="program-desc">{{ program.description }}</p>
        <div class="badges">
          <span class="badge" v-if="program.goal">{{ goalLabel(program.goal) }}</span>
          <span class="badge" v-if="program.format">{{ formatLabel(program.format) }}</span>
          <span class="badge" v-if="program.level">{{ levelLabel(program.level) }}</span>
        </div>
        <div class="program-meta">
          <span v-if="program.duration_weeks">{{ program.duration_weeks }} нед.</span>
        </div>
      </div>

      <!-- Week timeline -->
      <div class="weeks-timeline" v-if="weekMap.length > 0">
        <div
          v-for="(week, wIndex) in weekMap"
          :key="wIndex"
          class="week-block"
          :style="{ animationDelay: wIndex * 80 + 'ms' }"
        >
          <h3 class="week-title">Неделя {{ week.week }}</h3>
          <div class="workout-cards">
            <div
              v-for="workout in week.workouts"
              :key="workout.id"
              class="workout-card"
              @click="openWorkout(workout.id)"
            >
              <div class="workout-card-info">
                <span class="workout-day" v-if="workout.day_number">День {{ workout.day_number }}</span>
                <span class="workout-name">{{ workout.name }}</span>
              </div>
              <div class="workout-card-meta">
                <span v-if="workout.duration_minutes">{{ workout.duration_minutes }} мин</span>
                <span class="arrow">›</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-else class="empty">
        <p>Тренировки для этой программы пока не добавлены</p>
      </div>

      <button class="btn btn-primary" @click="enroll" :disabled="enrolling">
        {{ enrolling ? 'Подключение...' : 'Начать программу' }}
      </button>

      <div v-if="enrollError" class="error-msg">{{ enrollError }}</div>
      <div v-if="enrollSuccess" class="success-msg">Вы записаны на программу!</div>
    </template>

    <div v-else class="empty">
      <p>Программа не найдена</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api } from '../api'
import type { Program, Workout } from '../types'
import SkeletonCard from '../components/SkeletonCard.vue'

const props = defineProps<{ id: string }>()
const route = useRoute()
const router = useRouter()

const loading = ref(true)
const program = ref<Program | null>(null)
const workouts = ref<Workout[]>([])
const enrolling = ref(false)
const enrollError = ref('')
const enrollSuccess = ref(false)

const goalLabels: Record<string, string> = {
  weight_loss: 'Похудение',
  muscle_gain: 'Масса',
  strength: 'Сила',
  endurance: 'Выносливость',
  maintenance: 'Поддержание',
  rehabilitation: 'Реабилитация',
}

const formatLabels: Record<string, string> = {
  home: 'Дома',
  gym: 'Зал',
  outdoor: 'На улице',
  mixed: 'Смешанный',
}

const levelLabels: Record<string, string> = {
  beginner: 'Новичок',
  intermediate: 'Средний',
  advanced: 'Продвинутый',
}

function goalLabel(key: string): string { return goalLabels[key] || key }
function formatLabel(key: string): string { return formatLabels[key] || key }
function levelLabel(key: string): string { return levelLabels[key] || key }

interface WeekGroup {
  week: number
  workouts: Workout[]
}

const weekMap = computed<WeekGroup[]>(() => {
  if (!workouts.value.length) return []
  const map = new Map<number, Workout[]>()
  for (const w of workouts.value) {
    const week = w.week_number || 1
    if (!map.has(week)) map.set(week, [])
    map.get(week)!.push(w)
  }
  return Array.from(map.entries())
    .sort(([a], [b]) => a - b)
    .map(([week, wks]) => ({
      week,
      workouts: wks.sort((a, b) => (a.day_number || 0) - (b.day_number || 0)),
    }))
})

function openWorkout(id: number) {
  router.push({ name: 'workout-session', params: { id } })
}

async function enroll() {
  if (enrolling.value) return
  enrolling.value = true
  enrollError.value = ''
  enrollSuccess.value = false
  try {
    await api.enrollProgram(Number(props.id))
    enrollSuccess.value = true
  } catch (e: any) {
    enrollError.value = e.message || 'Ошибка при записи'
  } finally {
    enrolling.value = false
  }
}

onMounted(async () => {
  try {
    const [programData, allWorkouts] = await Promise.all([
      api.getProgram(Number(props.id)),
      api.getWorkouts({}),
    ])
    program.value = programData
    workouts.value = allWorkouts.filter(w => w.program_id === programData.id)
  } catch (e) {
    console.error('Failed to load program:', e)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.program-page {
  max-width: 400px;
  margin: 0 auto;
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

.skeleton-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.program-header {
  background: var(--secondary-bg);
  border-radius: 12px;
  padding: 20px 16px;
  margin-bottom: 16px;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease forwards;
}

.program-name {
  font-size: 22px;
  font-weight: 700;
  margin-bottom: 8px;
}

.program-desc {
  color: var(--hint-color);
  font-size: 14px;
  margin-bottom: 12px;
  line-height: 1.4;
}

.badges {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 8px;
}

.badge {
  background: var(--button-color);
  color: var(--button-text-color);
  padding: 4px 12px;
  border-radius: 16px;
  font-size: 12px;
  font-weight: 500;
}

.program-meta {
  color: var(--hint-color);
  font-size: 13px;
}

.weeks-timeline {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 20px;
}

.week-block {
  opacity: 0;
  animation: fadeSlideUp 0.35s ease forwards;
}

.week-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 8px;
  color: var(--hint-color);
}

.workout-cards {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.workout-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 16px;
  background: var(--secondary-bg);
  border-radius: 12px;
  cursor: pointer;
  transition: transform 0.15s ease;
}

.workout-card:active {
  transform: scale(0.98);
}

.workout-card-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.workout-day {
  font-size: 12px;
  color: var(--hint-color);
}

.workout-name {
  font-size: 15px;
  font-weight: 500;
}

.workout-card-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--hint-color);
  font-size: 13px;
}

.arrow {
  font-size: 20px;
}

.btn {
  display: block;
  width: 100%;
  padding: 14px;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  cursor: pointer;
  margin-bottom: 8px;
  text-align: center;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-primary {
  background: var(--button-color);
  color: var(--button-text-color);
}

.error-msg {
  text-align: center;
  color: #ff3b30;
  font-size: 14px;
  margin-top: 8px;
}

.success-msg {
  text-align: center;
  color: #34c759;
  font-size: 14px;
  margin-top: 8px;
}

.empty {
  text-align: center;
  padding: 40px 0;
  color: var(--hint-color);
}
</style>
