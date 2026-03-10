<template>
  <div class="progress-page">
    <h1 class="page-title">Прогресс</h1>

    <!-- Loading state -->
    <div v-if="loading" class="skeleton-list">
      <SkeletonCard v-for="i in 5" :key="i" />
    </div>

    <template v-else>
      <!-- Stats summary row -->
      <div class="stats-row">
        <div class="stat-card">
          <span class="stat-icon">🔥</span>
          <span class="stat-value">{{ stats.current_streak }}</span>
          <span class="stat-label">Серия, дн.</span>
        </div>
        <div class="stat-card">
          <span class="stat-icon">🏆</span>
          <span class="stat-value">{{ stats.longest_streak }}</span>
          <span class="stat-label">Рекорд, дн.</span>
        </div>
        <div class="stat-card">
          <span class="stat-icon">📏</span>
          <span class="stat-value">{{ entries.length }}</span>
          <span class="stat-label">Замеров</span>
        </div>
      </div>

      <!-- Add entry toggle -->
      <button
        v-if="!showForm"
        class="btn btn-primary add-btn"
        @click="openForm"
      >
        + Добавить запись
      </button>

      <!-- Inline add form -->
      <Transition name="form-slide">
        <div v-if="showForm" class="add-form">
          <h3 class="form-title">Новая запись</h3>

          <div class="form-field">
            <label>Дата</label>
            <input
              v-model="form.date"
              type="date"
              class="form-input"
            />
          </div>

          <div class="form-field">
            <label>Вес (кг)</label>
            <input
              v-model.number="form.weight_kg"
              type="number"
              step="0.1"
              min="30"
              max="300"
              placeholder="0.0"
              class="form-input"
            />
          </div>

          <div class="form-field">
            <label>Самочувствие</label>
            <select v-model="form.wellbeing" class="form-input">
              <option value="">-- Выберите --</option>
              <option value="excellent">Отлично</option>
              <option value="good">Хорошо</option>
              <option value="normal">Нормально</option>
              <option value="bad">Плохо</option>
            </select>
          </div>

          <div class="form-actions">
            <button
              class="btn btn-primary"
              @click="submitEntry"
              :disabled="submitting"
            >
              {{ submitting ? 'Сохранение...' : 'Сохранить' }}
            </button>
            <button class="btn btn-secondary" @click="showForm = false">
              Отмена
            </button>
          </div>

          <div v-if="formError" class="error-msg">{{ formError }}</div>
        </div>
      </Transition>

      <!-- Weight history -->
      <section class="section" v-if="weightEntries.length > 0">
        <h2 class="section-title">История веса</h2>

        <!-- Mini bar chart -->
        <div class="weight-chart" v-if="weightEntries.length >= 2">
          <div class="chart-area">
            <div
              v-for="(point, idx) in chartPoints"
              :key="point.date"
              class="chart-bar-wrapper"
            >
              <span class="chart-value">{{ point.weight.toFixed(1) }}</span>
              <div
                class="chart-bar"
                :style="{
                  height: point.heightPct + '%',
                  animationDelay: idx * 60 + 'ms',
                }"
              ></div>
              <span class="chart-date">{{ point.label }}</span>
            </div>
          </div>
          <div class="chart-summary">
            <span v-if="weightDelta !== null" class="weight-delta" :class="weightDeltaClass">
              {{ weightDeltaSign }}{{ Math.abs(weightDelta).toFixed(1) }} кг
            </span>
            <span class="weight-delta-hint">за последние записи</span>
          </div>
        </div>

        <!-- Weight entries list -->
        <div class="entries-list">
          <div
            v-for="(entry, idx) in weightEntries.slice(0, 10)"
            :key="entry.id"
            class="entry-row"
            :style="{ animationDelay: idx * 40 + 'ms' }"
          >
            <div class="entry-date">{{ formatDate(entry.date) }}</div>
            <div class="entry-weight">
              {{ entry.weight_kg!.toFixed(1) }} кг
            </div>
            <div class="entry-trend">
              {{ getTrend(idx) }}
            </div>
          </div>
        </div>
      </section>

      <!-- Calendar streak -->
      <section class="section">
        <h2 class="section-title">Календарь активности</h2>
        <div class="calendar">
          <div class="calendar-header">
            <button class="cal-nav-btn" @click="prevMonth">‹</button>
            <span class="cal-month-label">{{ calendarMonthLabel }}</span>
            <button class="cal-nav-btn" @click="nextMonth">›</button>
          </div>
          <div class="calendar-weekdays">
            <span v-for="d in weekdays" :key="d" class="weekday">{{ d }}</span>
          </div>
          <div class="calendar-grid">
            <div
              v-for="(cell, idx) in calendarCells"
              :key="idx"
              class="calendar-cell"
              :class="{
                empty: cell.day === 0,
                active: cell.isActive,
                today: cell.isToday,
              }"
            >
              <span v-if="cell.day > 0" class="cell-day">{{ cell.day }}</span>
              <span v-if="cell.isActive" class="cell-dot"></span>
            </div>
          </div>
        </div>
      </section>

      <!-- Achievements -->
      <section class="section" v-if="allAchievements.length > 0">
        <h2 class="section-title">Достижения</h2>
        <div class="achievements-grid">
          <div
            v-for="(ach, idx) in allAchievements"
            :key="ach.id"
            class="achievement-badge"
            :class="{ earned: ach.earned }"
            :style="{ animationDelay: idx * 50 + 'ms' }"
          >
            <span class="ach-icon">{{ ach.icon || '🎖️' }}</span>
            <span class="ach-name">{{ ach.name }}</span>
            <span v-if="ach.earned" class="ach-earned-label">Получено</span>
          </div>
        </div>
      </section>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api } from '../api'
import type {
  ProgressEntry,
  ProgressStats,
  Achievement,
  UserAchievement,
  CreateProgressRequest,
} from '../types'
import SkeletonCard from '../components/SkeletonCard.vue'

const route = useRoute()
const router = useRouter()

// ---------- State ----------

const loading = ref(true)
const entries = ref<ProgressEntry[]>([])
const stats = ref<ProgressStats>({ current_streak: 0, longest_streak: 0, calendar: [] })
const achievements = ref<Achievement[]>([])
const userAchievements = ref<UserAchievement[]>([])

const showForm = ref(false)
const submitting = ref(false)
const formError = ref('')

const calendarOffset = ref(0) // 0 = current month, -1 = prev, etc.

function todayString(): string {
  const d = new Date()
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

const form = reactive<CreateProgressRequest>({
  date: todayString(),
  weight_kg: undefined,
  wellbeing: '',
})

// ---------- Computed ----------

const weightEntries = computed<ProgressEntry[]>(() => {
  return entries.value
    .filter((e) => e.weight_kg != null)
    .sort((a, b) => b.date.localeCompare(a.date))
})

/** Up to 8 most recent weight points for the mini bar chart. */
const chartPoints = computed(() => {
  const recent = weightEntries.value.slice(0, 8).reverse()
  if (recent.length === 0) return []

  const weights = recent.map((e) => e.weight_kg!)
  const min = Math.min(...weights)
  const max = Math.max(...weights)
  const range = max - min || 1

  return recent.map((e) => ({
    date: e.date,
    weight: e.weight_kg!,
    heightPct: 20 + ((e.weight_kg! - min) / range) * 70,
    label: formatShortDate(e.date),
  }))
})

const weightDelta = computed<number | null>(() => {
  if (weightEntries.value.length < 2) return null
  const newest = weightEntries.value[0].weight_kg!
  const oldest = weightEntries.value[weightEntries.value.length - 1].weight_kg!
  return newest - oldest
})

const weightDeltaSign = computed(() => {
  if (weightDelta.value === null) return ''
  return weightDelta.value > 0 ? '+' : weightDelta.value < 0 ? '' : ''
})

const weightDeltaClass = computed(() => {
  if (weightDelta.value === null) return ''
  if (weightDelta.value < 0) return 'delta-down'
  if (weightDelta.value > 0) return 'delta-up'
  return 'delta-same'
})

const weekdays = ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс']

const calendarDate = computed(() => {
  const d = new Date()
  d.setDate(1)
  d.setMonth(d.getMonth() + calendarOffset.value)
  return d
})

const calendarMonthLabel = computed(() => {
  return calendarDate.value.toLocaleDateString('ru-RU', {
    month: 'long',
    year: 'numeric',
  })
})

interface CalendarCell {
  day: number
  dateStr: string
  isActive: boolean
  isToday: boolean
}

const calendarCells = computed<CalendarCell[]>(() => {
  const year = calendarDate.value.getFullYear()
  const month = calendarDate.value.getMonth()
  const firstDay = new Date(year, month, 1)
  const daysInMonth = new Date(year, month + 1, 0).getDate()

  // Monday = 0, Sunday = 6 (ISO style)
  let startDay = firstDay.getDay() - 1
  if (startDay < 0) startDay = 6

  const today = todayString()
  const calendarSet = new Set(stats.value.calendar)

  const cells: CalendarCell[] = []

  // Leading empty cells
  for (let i = 0; i < startDay; i++) {
    cells.push({ day: 0, dateStr: '', isActive: false, isToday: false })
  }

  for (let d = 1; d <= daysInMonth; d++) {
    const dateStr = `${year}-${String(month + 1).padStart(2, '0')}-${String(d).padStart(2, '0')}`
    cells.push({
      day: d,
      dateStr,
      isActive: calendarSet.has(dateStr),
      isToday: dateStr === today,
    })
  }

  return cells
})

interface AchievementDisplay {
  id: number
  icon: string
  name: string
  description: string
  earned: boolean
  earned_at?: string
}

const allAchievements = computed<AchievementDisplay[]>(() => {
  const earnedIds = new Set(userAchievements.value.map((ua) => ua.achievement_id))
  return achievements.value.map((a) => ({
    id: a.id,
    icon: a.icon,
    name: a.name,
    description: a.description,
    earned: earnedIds.has(a.id),
    earned_at: userAchievements.value.find((ua) => ua.achievement_id === a.id)?.earned_at,
  }))
})

// ---------- Helpers ----------

function formatDate(dateStr: string): string {
  const d = new Date(dateStr + 'T00:00:00')
  return d.toLocaleDateString('ru-RU', {
    day: 'numeric',
    month: 'short',
  })
}

function formatShortDate(dateStr: string): string {
  const d = new Date(dateStr + 'T00:00:00')
  return d.toLocaleDateString('ru-RU', {
    day: 'numeric',
    month: 'short',
  }).replace('.', '')
}

function getTrend(idx: number): string {
  const list = weightEntries.value
  if (idx >= list.length - 1) return ''
  const current = list[idx].weight_kg!
  const prev = list[idx + 1].weight_kg!
  if (current > prev) return '↑'
  if (current < prev) return '↓'
  return '→'
}

function prevMonth() {
  calendarOffset.value--
}

function nextMonth() {
  calendarOffset.value++
}

// ---------- Form ----------

function openForm() {
  form.date = todayString()
  form.weight_kg = undefined
  form.wellbeing = ''
  formError.value = ''
  showForm.value = true
}

async function submitEntry() {
  if (submitting.value) return

  if (!form.weight_kg && !form.wellbeing) {
    formError.value = 'Укажите вес или самочувствие'
    return
  }

  submitting.value = true
  formError.value = ''

  try {
    const payload: CreateProgressRequest = {
      date: form.date,
    }
    if (form.weight_kg) payload.weight_kg = form.weight_kg
    if (form.wellbeing) payload.wellbeing = form.wellbeing

    await api.addProgressEntry(payload)
    showForm.value = false
    await loadData()
  } catch (e: any) {
    formError.value = e.message || 'Ошибка при сохранении'
  } finally {
    submitting.value = false
  }
}

// ---------- Data loading ----------

async function loadData() {
  try {
    const [entriesData, statsData, achievementsData] = await Promise.all([
      api.getProgressEntries(),
      api.getProgressStats(),
      api.getAchievements(),
    ])
    entries.value = entriesData
    stats.value = statsData
    achievements.value = achievementsData.achievements || []
    userAchievements.value = achievementsData.user_achievements || []
  } catch (e) {
    console.error('Failed to load progress data:', e)
  }
}

onMounted(async () => {
  try {
    await loadData()
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.progress-page {
  max-width: 400px;
  margin: 0 auto;
  padding-bottom: 80px;
}

.page-title {
  font-size: 22px;
  font-weight: 700;
  margin-bottom: 16px;
}

.skeleton-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* ===== Stats row ===== */
.stats-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
  margin-bottom: 16px;
}

.stat-card {
  background: var(--secondary-bg);
  border-radius: 12px;
  padding: 14px 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease forwards;
}

.stat-card:nth-child(1) { animation-delay: 0ms; }
.stat-card:nth-child(2) { animation-delay: 80ms; }
.stat-card:nth-child(3) { animation-delay: 160ms; }

.stat-icon {
  font-size: 22px;
  line-height: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-color);
}

.stat-label {
  font-size: 11px;
  color: var(--hint-color);
  text-align: center;
}

/* ===== Add button ===== */
.add-btn {
  margin-bottom: 20px;
}

/* ===== Inline form ===== */
.add-form {
  background: var(--secondary-bg);
  border-radius: 12px;
  padding: 20px 16px;
  margin-bottom: 20px;
}

.form-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 16px;
}

.form-field {
  margin-bottom: 12px;
}

.form-field label {
  display: block;
  font-size: 13px;
  color: var(--hint-color);
  margin-bottom: 4px;
}

.form-input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid var(--hint-color);
  border-radius: 10px;
  font-size: 14px;
  background: var(--bg-color);
  color: var(--text-color);
  font-family: inherit;
  box-sizing: border-box;
}

.form-input::placeholder {
  color: var(--hint-color);
}

select.form-input {
  -webkit-appearance: none;
  appearance: none;
}

.form-actions {
  display: flex;
  gap: 12px;
  margin-top: 4px;
}

.form-slide-enter-active,
.form-slide-leave-active {
  transition: all 0.25s ease;
}

.form-slide-enter-from,
.form-slide-leave-to {
  opacity: 0;
  transform: translateY(-12px);
}

/* ===== Weight chart ===== */
.weight-chart {
  background: var(--secondary-bg);
  border-radius: 12px;
  padding: 16px 12px;
  margin-bottom: 12px;
}

.chart-area {
  display: flex;
  align-items: flex-end;
  justify-content: space-around;
  height: 120px;
  gap: 4px;
  padding-bottom: 4px;
}

.chart-bar-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  flex: 1;
  min-width: 0;
  height: 100%;
  justify-content: flex-end;
}

.chart-value {
  font-size: 9px;
  color: var(--hint-color);
  white-space: nowrap;
}

.chart-bar {
  width: 100%;
  max-width: 28px;
  min-height: 4px;
  background: var(--button-color);
  border-radius: 4px 4px 2px 2px;
  opacity: 0;
  animation: barGrow 0.4s ease forwards;
}

.chart-date {
  font-size: 9px;
  color: var(--hint-color);
  white-space: nowrap;
  text-align: center;
}

.chart-summary {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  margin-top: 10px;
  padding-top: 10px;
  border-top: 1px solid rgba(128, 128, 128, 0.15);
}

.weight-delta {
  font-size: 15px;
  font-weight: 600;
}

.delta-down {
  color: #34c759;
}

.delta-up {
  color: #ff9500;
}

.delta-same {
  color: var(--hint-color);
}

.weight-delta-hint {
  font-size: 12px;
  color: var(--hint-color);
}

@keyframes barGrow {
  0% {
    opacity: 0;
    transform: scaleY(0);
    transform-origin: bottom;
  }
  100% {
    opacity: 1;
    transform: scaleY(1);
    transform-origin: bottom;
  }
}

/* ===== Weight entries list ===== */
.entries-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.entry-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 14px;
  background: var(--secondary-bg);
  border-radius: 10px;
  opacity: 0;
  animation: fadeSlideUp 0.3s ease forwards;
}

.entry-date {
  font-size: 14px;
  color: var(--hint-color);
  min-width: 70px;
}

.entry-weight {
  font-size: 15px;
  font-weight: 600;
  flex: 1;
  text-align: center;
}

.entry-trend {
  font-size: 18px;
  width: 24px;
  text-align: center;
}

/* ===== Section ===== */
.section {
  margin-top: 24px;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 12px;
}

/* ===== Calendar ===== */
.calendar {
  background: var(--secondary-bg);
  border-radius: 12px;
  padding: 16px;
}

.calendar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.cal-nav-btn {
  background: none;
  border: none;
  font-size: 22px;
  color: var(--button-color);
  cursor: pointer;
  padding: 4px 12px;
  border-radius: 8px;
  transition: background 0.15s;
}

.cal-nav-btn:active {
  background: var(--bg-color);
}

.cal-month-label {
  font-size: 15px;
  font-weight: 600;
  text-transform: capitalize;
}

.calendar-weekdays {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 2px;
  margin-bottom: 6px;
}

.weekday {
  text-align: center;
  font-size: 11px;
  color: var(--hint-color);
  font-weight: 500;
  padding: 4px 0;
}

.calendar-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 2px;
}

.calendar-cell {
  position: relative;
  aspect-ratio: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  gap: 2px;
}

.calendar-cell.empty {
  visibility: hidden;
}

.cell-day {
  font-size: 13px;
  color: var(--text-color);
  line-height: 1;
}

.calendar-cell.today .cell-day {
  font-weight: 700;
  color: var(--button-color);
}

.calendar-cell.active {
  background: color-mix(in srgb, var(--button-color) 12%, transparent);
}

.cell-dot {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: var(--button-color);
}

/* ===== Achievements ===== */
.achievements-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
}

.achievement-badge {
  background: var(--secondary-bg);
  border-radius: 12px;
  padding: 14px 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  text-align: center;
  opacity: 0.4;
  filter: grayscale(1);
  transition: all 0.3s ease;
  animation: fadeIn 0.3s ease forwards;
}

.achievement-badge.earned {
  opacity: 1;
  filter: grayscale(0);
  border: 1.5px solid var(--button-color);
}

.ach-icon {
  font-size: 28px;
  line-height: 1;
}

.ach-name {
  font-size: 11px;
  font-weight: 500;
  color: var(--text-color);
  line-height: 1.2;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.ach-earned-label {
  font-size: 9px;
  color: var(--button-color);
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.3px;
}

/* ===== Shared ===== */
.btn {
  display: block;
  width: 100%;
  padding: 14px;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  cursor: pointer;
  text-align: center;
  font-family: inherit;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-primary {
  background: var(--button-color);
  color: var(--button-text-color);
}

.btn-secondary {
  background: var(--bg-color);
  color: var(--text-color);
  border: 1px solid var(--hint-color);
}

.error-msg {
  text-align: center;
  color: #ff3b30;
  font-size: 14px;
  margin-top: 8px;
}

/* ===== Animations ===== */
@keyframes fadeSlideUp {
  0% {
    opacity: 0;
    transform: translateY(12px);
  }
  100% {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fadeIn {
  0% { opacity: 0; }
  100% { opacity: 0.4; }
}

.achievement-badge.earned {
  animation-name: fadeInFull;
}

@keyframes fadeInFull {
  0% { opacity: 0; }
  100% { opacity: 1; }
}
</style>
