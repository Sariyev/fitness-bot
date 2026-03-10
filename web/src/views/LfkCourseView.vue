<template>
  <div class="course-page">
    <button class="back-btn" @click="router.back()">← Назад</button>

    <div v-if="loading" class="skeleton-list">
      <SkeletonCard v-for="i in 5" :key="i" />
    </div>

    <template v-else-if="course">
      <div class="course-header">
        <h1 class="course-name">{{ course.name }}</h1>
        <p class="course-desc">{{ course.description }}</p>
      </div>

      <!-- Warnings card -->
      <div v-if="course.warnings" class="warning-card">
        <span class="warning-icon">⚠️</span>
        <p>{{ course.warnings }}</p>
      </div>

      <!-- 14-day timeline -->
      <div class="timeline">
        <!-- Stage 1: days 1-4 -->
        <div class="stage-block">
          <div class="stage-label">Этап 1 <span class="stage-days">дни 1-4</span></div>
          <div class="days-grid">
            <div
              v-for="day in getDays(1, 4)"
              :key="day"
              class="day-cell"
              :class="{
                completed: isDayCompleted(day),
                current: !isDayCompleted(day) && isNextDay(day),
              }"
              @click="openDay(day)"
            >
              <span class="day-number">{{ day }}</span>
              <span v-if="isDayCompleted(day)" class="day-check">✓</span>
              <span v-if="isDayCompleted(day) && getDayPain(day) !== null" class="day-pain"
                :style="{ color: painColor(getDayPain(day)!) }">
                {{ getDayPain(day) }}
              </span>
            </div>
          </div>
        </div>

        <!-- Stage 2: days 5-9 -->
        <div class="stage-block">
          <div class="stage-label">Этап 2 <span class="stage-days">дни 5-9</span></div>
          <div class="days-grid">
            <div
              v-for="day in getDays(5, 9)"
              :key="day"
              class="day-cell"
              :class="{
                completed: isDayCompleted(day),
                current: !isDayCompleted(day) && isNextDay(day),
              }"
              @click="openDay(day)"
            >
              <span class="day-number">{{ day }}</span>
              <span v-if="isDayCompleted(day)" class="day-check">✓</span>
              <span v-if="isDayCompleted(day) && getDayPain(day) !== null" class="day-pain"
                :style="{ color: painColor(getDayPain(day)!) }">
                {{ getDayPain(day) }}
              </span>
            </div>
          </div>
        </div>

        <!-- Stage 3: days 10-14 -->
        <div class="stage-block">
          <div class="stage-label">Этап 3 <span class="stage-days">дни 10-14</span></div>
          <div class="days-grid">
            <div
              v-for="day in getDays(10, 14)"
              :key="day"
              class="day-cell"
              :class="{
                completed: isDayCompleted(day),
                current: !isDayCompleted(day) && isNextDay(day),
              }"
              @click="openDay(day)"
            >
              <span class="day-number">{{ day }}</span>
              <span v-if="isDayCompleted(day)" class="day-check">✓</span>
              <span v-if="isDayCompleted(day) && getDayPain(day) !== null" class="day-pain"
                :style="{ color: painColor(getDayPain(day)!) }">
                {{ getDayPain(day) }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </template>

    <div v-else class="empty">
      <p>Курс не найден</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api } from '../api'
import type { RehabCourseWithSessions, RehabProgress } from '../types'
import SkeletonCard from '../components/SkeletonCard.vue'

const props = defineProps<{ id: string }>()
const route = useRoute()
const router = useRouter()

const loading = ref(true)
const course = ref<RehabCourseWithSessions | null>(null)
const progress = ref<RehabProgress[]>([])

function getDays(from: number, to: number): number[] {
  const days: number[] = []
  for (let i = from; i <= to; i++) days.push(i)
  return days
}

const completedDays = computed(() => {
  const map = new Map<number, RehabProgress>()
  for (const p of progress.value) {
    map.set(p.day_number, p)
  }
  return map
})

function isDayCompleted(day: number): boolean {
  return completedDays.value.has(day)
}

function getDayPain(day: number): number | null {
  const p = completedDays.value.get(day)
  return p ? p.pain_level : null
}

function isNextDay(day: number): boolean {
  // Find the highest completed day
  let maxCompleted = 0
  for (const d of completedDays.value.keys()) {
    if (d > maxCompleted) maxCompleted = d
  }
  return day === maxCompleted + 1
}

function painColor(level: number): string {
  if (level <= 3) return '#34c759'
  if (level <= 6) return '#ffcc00'
  return '#ff3b30'
}

function openDay(day: number) {
  if (!course.value) return
  // Find the session for this day
  const session = course.value.sessions?.find(s => s.day_number === day)
  if (session) {
    router.push({ name: 'lfk-session', params: { id: session.id } })
  }
}

onMounted(async () => {
  try {
    const courseId = Number(props.id)
    const [courseData, progressData] = await Promise.all([
      api.getRehabCourse(courseId),
      api.getRehabProgress(courseId),
    ])
    course.value = courseData
    progress.value = progressData
  } catch (e) {
    console.error('Failed to load course:', e)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.course-page {
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

.course-header {
  background: var(--secondary-bg);
  border-radius: 12px;
  padding: 20px 16px;
  margin-bottom: 12px;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease forwards;
}

.course-name {
  font-size: 22px;
  font-weight: 700;
  margin-bottom: 8px;
}

.course-desc {
  color: var(--hint-color);
  font-size: 14px;
  line-height: 1.4;
}

.warning-card {
  background: #ff3b3015;
  border: 1px solid #ff3b3040;
  border-radius: 12px;
  padding: 14px 16px;
  margin-bottom: 16px;
  display: flex;
  gap: 10px;
  align-items: flex-start;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease 0.1s forwards;
}

.warning-icon {
  font-size: 20px;
  flex-shrink: 0;
}

.warning-card p {
  font-size: 14px;
  color: #ff3b30;
  line-height: 1.4;
}

.timeline {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.stage-block {
  opacity: 0;
  animation: fadeSlideUp 0.35s ease forwards;
}

.stage-block:nth-child(1) { animation-delay: 0.15s; }
.stage-block:nth-child(2) { animation-delay: 0.25s; }
.stage-block:nth-child(3) { animation-delay: 0.35s; }

.stage-label {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 8px;
}

.stage-days {
  font-size: 13px;
  font-weight: 400;
  color: var(--hint-color);
}

.days-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 8px;
}

.day-cell {
  background: var(--secondary-bg);
  border-radius: 12px;
  padding: 12px 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  cursor: pointer;
  transition: transform 0.15s ease, box-shadow 0.15s ease;
  position: relative;
}

.day-cell:active {
  transform: scale(0.95);
}

.day-cell.completed {
  background: #34c75920;
}

.day-cell.current {
  box-shadow: 0 0 0 2px var(--button-color);
}

.day-number {
  font-size: 16px;
  font-weight: 600;
}

.day-check {
  font-size: 14px;
  color: #34c759;
  font-weight: 700;
}

.day-pain {
  font-size: 11px;
  font-weight: 600;
}

.empty {
  text-align: center;
  padding: 40px 0;
  color: var(--hint-color);
}
</style>
