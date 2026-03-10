<template>
  <div class="diary-page">
    <!-- Header -->
    <div class="header">
      <button class="back-btn" @click="router.back()">← Назад</button>
      <h1 class="page-title">Дневник питания</h1>
    </div>

    <!-- Date selector -->
    <div class="date-picker">
      <button class="date-arrow" @click="prevDay">&#8249;</button>
      <div class="date-display">
        <span class="date-weekday">{{ weekdayLabel }}</span>
        <span class="date-full">{{ formattedDate }}</span>
      </div>
      <button class="date-arrow" @click="nextDay">&#8250;</button>
    </div>

    <!-- Loading state -->
    <div v-if="loading" class="skeleton-list">
      <SkeletonCard v-for="i in 4" :key="i" />
    </div>

    <template v-else>
      <!-- Daily summary card -->
      <div class="summary-card">
        <div class="summary-calories">
          <span class="calories-number">{{ summary.calories }}</span>
          <span class="calories-unit">ккал</span>
        </div>
        <div class="macro-bars">
          <div class="macro-bar-item">
            <div class="macro-bar-track">
              <div
                class="macro-bar-fill protein-fill"
                :style="{ width: macroBarWidth(summary.protein, 'protein') }"
              ></div>
            </div>
            <span class="macro-bar-label">Белки {{ summary.protein }}г</span>
          </div>
          <div class="macro-bar-item">
            <div class="macro-bar-track">
              <div
                class="macro-bar-fill fat-fill"
                :style="{ width: macroBarWidth(summary.fat, 'fat') }"
              ></div>
            </div>
            <span class="macro-bar-label">Жиры {{ summary.fat }}г</span>
          </div>
          <div class="macro-bar-item">
            <div class="macro-bar-track">
              <div
                class="macro-bar-fill carbs-fill"
                :style="{ width: macroBarWidth(summary.carbs, 'carbs') }"
              ></div>
            </div>
            <span class="macro-bar-label">Углеводы {{ summary.carbs }}г</span>
          </div>
        </div>
      </div>

      <!-- Entries grouped by meal type -->
      <div v-for="group in groupedEntries" :key="group.type" class="meal-group">
        <div class="meal-group-header">
          <span class="meal-emoji">{{ mealEmoji(group.type) }}</span>
          <h3 class="meal-type-title">{{ mealTypeLabel(group.type) }}</h3>
          <span class="meal-group-total">{{ groupCalories(group.entries) }} ккал</span>
        </div>
        <div class="entries-list">
          <div
            v-for="entry in group.entries"
            :key="entry.id"
            class="entry-card"
          >
            <div class="entry-info">
              <span class="entry-name">{{ entry.food_name }}</span>
              <span class="entry-macros">
                {{ entry.calories }} ккал
                · Б {{ entry.protein }}г
                · Ж {{ entry.fat }}г
                · У {{ entry.carbs }}г
              </span>
            </div>
            <button
              class="delete-btn"
              @click="deleteEntry(entry.id)"
              :disabled="deletingId === entry.id"
            >
              <span v-if="deletingId === entry.id" class="delete-spinner"></span>
              <span v-else>🗑</span>
            </button>
          </div>
        </div>
      </div>

      <!-- Empty state -->
      <div v-if="entries.length === 0 && !showForm" class="empty-state">
        <div class="empty-icon">📋</div>
        <p class="empty-title">Записей пока нет</p>
        <p class="empty-hint">Добавьте первый прием пищи за этот день</p>
      </div>

      <!-- Inline add form -->
      <Transition name="form-slide">
        <div v-if="showForm" class="add-form">
          <h3 class="form-title">Новая запись</h3>

          <!-- Meal type chips -->
          <div class="form-field">
            <label class="field-label">Прием пищи</label>
            <div class="meal-chips">
              <button
                v-for="mt in mealTypes"
                :key="mt.value"
                class="meal-chip"
                :class="{ active: form.meal_type === mt.value }"
                @click="form.meal_type = mt.value"
              >
                <span class="chip-emoji">{{ mt.emoji }}</span>
                <span class="chip-text">{{ mt.label }}</span>
              </button>
            </div>
          </div>

          <!-- Food name -->
          <div class="form-field">
            <label class="field-label">Название продукта</label>
            <input
              v-model="form.food_name"
              type="text"
              placeholder="Например: Овсянка с бананом"
              class="form-input"
              ref="foodNameInput"
            />
          </div>

          <!-- Calories -->
          <div class="form-field">
            <label class="field-label">Калории (ккал)</label>
            <input
              v-model.number="form.calories"
              type="number"
              inputmode="numeric"
              placeholder="0"
              min="0"
              class="form-input"
            />
          </div>

          <!-- Macros row -->
          <div class="form-row-3">
            <div class="form-field">
              <label class="field-label protein-label">Белки (г)</label>
              <input
                v-model.number="form.protein"
                type="number"
                inputmode="numeric"
                placeholder="0"
                min="0"
                class="form-input"
              />
            </div>
            <div class="form-field">
              <label class="field-label fat-label">Жиры (г)</label>
              <input
                v-model.number="form.fat"
                type="number"
                inputmode="numeric"
                placeholder="0"
                min="0"
                class="form-input"
              />
            </div>
            <div class="form-field">
              <label class="field-label carbs-label">Углеводы (г)</label>
              <input
                v-model.number="form.carbs"
                type="number"
                inputmode="numeric"
                placeholder="0"
                min="0"
                class="form-input"
              />
            </div>
          </div>

          <!-- Form actions -->
          <div class="form-actions">
            <button
              class="btn btn-primary"
              @click="submitEntry"
              :disabled="submitting || !form.food_name.trim()"
            >
              {{ submitting ? 'Сохранение...' : 'Сохранить' }}
            </button>
            <button class="btn btn-secondary" @click="closeForm">Отмена</button>
          </div>

          <div v-if="formError" class="error-msg">{{ formError }}</div>
        </div>
      </Transition>

      <!-- Floating add button -->
      <button v-if="!showForm" class="fab" @click="openForm">
        <span class="fab-icon">+</span>
        <span class="fab-text">Добавить</span>
      </button>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../api'
import type { FoodLogEntry, DailySummary } from '../types'
import SkeletonCard from '../components/SkeletonCard.vue'

const router = useRouter()

const loading = ref(true)
const entries = ref<FoodLogEntry[]>([])
const summary = ref<DailySummary>({ calories: 0, protein: 0, fat: 0, carbs: 0 })
const currentDate = ref(new Date())
const showForm = ref(false)
const submitting = ref(false)
const deletingId = ref<number | null>(null)
const formError = ref('')
const foodNameInput = ref<HTMLInputElement | null>(null)

const form = reactive({
  food_name: '',
  meal_type: 'breakfast',
  calories: 0,
  protein: 0,
  fat: 0,
  carbs: 0,
})

interface MealTypeOption {
  value: string
  label: string
  emoji: string
}

const mealTypes: MealTypeOption[] = [
  { value: 'breakfast', label: 'Завтрак', emoji: '🌅' },
  { value: 'lunch', label: 'Обед', emoji: '☀️' },
  { value: 'dinner', label: 'Ужин', emoji: '🌙' },
  { value: 'snack', label: 'Перекус', emoji: '🍎' },
]

const mealTypeLabels: Record<string, string> = {
  breakfast: 'Завтрак',
  lunch: 'Обед',
  dinner: 'Ужин',
  snack: 'Перекус',
}

const mealEmojis: Record<string, string> = {
  breakfast: '🌅',
  lunch: '☀️',
  dinner: '🌙',
  snack: '🍎',
}

function mealTypeLabel(type: string): string {
  return mealTypeLabels[type] || type
}

function mealEmoji(type: string): string {
  return mealEmojis[type] || '🍽'
}

function groupCalories(groupEntries: FoodLogEntry[]): number {
  return groupEntries.reduce((sum, e) => sum + e.calories, 0)
}

// Date utilities
const weekdayLabel = computed(() => {
  const today = new Date()
  const d = currentDate.value
  if (dateString(d) === dateString(today)) return 'Сегодня'
  const yesterday = new Date(today)
  yesterday.setDate(today.getDate() - 1)
  if (dateString(d) === dateString(yesterday)) return 'Вчера'
  const tomorrow = new Date(today)
  tomorrow.setDate(today.getDate() + 1)
  if (dateString(d) === dateString(tomorrow)) return 'Завтра'
  return d.toLocaleDateString('ru-RU', { weekday: 'long' })
})

const formattedDate = computed(() => {
  const d = currentDate.value
  return d.toLocaleDateString('ru-RU', { day: 'numeric', month: 'long', year: 'numeric' })
})

function dateString(d: Date): string {
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

function prevDay(): void {
  const d = new Date(currentDate.value)
  d.setDate(d.getDate() - 1)
  currentDate.value = d
}

function nextDay(): void {
  const d = new Date(currentDate.value)
  d.setDate(d.getDate() + 1)
  currentDate.value = d
}

// Macro bar width calculation (proportional to a reasonable daily max)
function macroBarWidth(grams: number, type: string): string {
  const maxValues: Record<string, number> = {
    protein: 200,
    fat: 150,
    carbs: 400,
  }
  const max = maxValues[type] || 200
  const pct = Math.min(100, (grams / max) * 100)
  return `${pct}%`
}

// Grouped entries
interface MealGroup {
  type: string
  entries: FoodLogEntry[]
}

const mealOrder = ['breakfast', 'lunch', 'dinner', 'snack']

const groupedEntries = computed<MealGroup[]>(() => {
  const map = new Map<string, FoodLogEntry[]>()
  for (const e of entries.value) {
    const t = e.meal_type || 'snack'
    if (!map.has(t)) map.set(t, [])
    map.get(t)!.push(e)
  }
  return mealOrder
    .filter(t => map.has(t))
    .map(t => ({ type: t, entries: map.get(t)! }))
})

// Data loading
async function loadData(): Promise<void> {
  loading.value = true
  const ds = dateString(currentDate.value)
  try {
    const [logData, summaryData] = await Promise.all([
      api.getFoodLog(ds),
      api.getFoodLogSummary(ds),
    ])
    entries.value = logData
    summary.value = summaryData
  } catch (e) {
    console.error('Failed to load food log:', e)
    entries.value = []
    summary.value = { calories: 0, protein: 0, fat: 0, carbs: 0 }
  } finally {
    loading.value = false
  }
}

// Form management
function openForm(): void {
  form.food_name = ''
  form.meal_type = 'breakfast'
  form.calories = 0
  form.protein = 0
  form.fat = 0
  form.carbs = 0
  formError.value = ''
  showForm.value = true
  nextTick(() => {
    foodNameInput.value?.focus()
  })
}

function closeForm(): void {
  showForm.value = false
  formError.value = ''
}

async function submitEntry(): Promise<void> {
  if (submitting.value) return
  if (!form.food_name.trim()) {
    formError.value = 'Введите название продукта'
    return
  }
  submitting.value = true
  formError.value = ''
  try {
    await api.addFoodLog({
      date: dateString(currentDate.value),
      meal_type: form.meal_type,
      food_name: form.food_name.trim(),
      calories: form.calories || 0,
      protein: form.protein || 0,
      fat: form.fat || 0,
      carbs: form.carbs || 0,
    })
    showForm.value = false
    await loadData()
  } catch (e: any) {
    formError.value = e.message || 'Ошибка при сохранении'
  } finally {
    submitting.value = false
  }
}

async function deleteEntry(id: number): Promise<void> {
  if (deletingId.value !== null) return
  deletingId.value = id
  try {
    await api.deleteFoodLog(id)
    await loadData()
  } catch (e) {
    console.error('Failed to delete entry:', e)
  } finally {
    deletingId.value = null
  }
}

// Reload when date changes
watch(currentDate, () => {
  showForm.value = false
  formError.value = ''
  loadData()
})

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.diary-page {
  max-width: 400px;
  margin: 0 auto;
  padding: 0 0 100px;
}

/* Header */
.header {
  margin-bottom: 16px;
}

.back-btn {
  background: none;
  border: none;
  color: var(--button-color);
  font-size: 16px;
  cursor: pointer;
  padding: 4px 0;
  margin-bottom: 8px;
}

.page-title {
  font-size: 22px;
  font-weight: 700;
  color: var(--text-color);
}

/* Date picker */
.date-picker {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 8px;
  background: var(--secondary-bg);
  border-radius: 14px;
  margin-bottom: 16px;
}

.date-arrow {
  background: none;
  border: none;
  font-size: 28px;
  line-height: 1;
  color: var(--button-color);
  cursor: pointer;
  padding: 4px 14px;
  border-radius: 10px;
  transition: background 0.15s;
  user-select: none;
  -webkit-user-select: none;
}

.date-arrow:active {
  background: var(--bg-color);
}

.date-display {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
}

.date-weekday {
  font-size: 13px;
  color: var(--button-color);
  font-weight: 600;
  text-transform: capitalize;
}

.date-full {
  font-size: 15px;
  font-weight: 500;
  color: var(--text-color);
}

/* Loading skeleton */
.skeleton-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* Summary card */
.summary-card {
  background: var(--secondary-bg);
  border-radius: 14px;
  padding: 20px 16px;
  margin-bottom: 20px;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease forwards;
}

.summary-calories {
  text-align: center;
  margin-bottom: 16px;
}

.calories-number {
  font-size: 40px;
  font-weight: 800;
  color: var(--text-color);
  line-height: 1;
}

.calories-unit {
  display: block;
  font-size: 14px;
  color: var(--hint-color);
  margin-top: 4px;
}

/* Macro bars */
.macro-bars {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.macro-bar-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.macro-bar-track {
  height: 8px;
  background: var(--bg-color);
  border-radius: 4px;
  overflow: hidden;
}

.macro-bar-fill {
  height: 100%;
  border-radius: 4px;
  transition: width 0.4s ease;
  min-width: 2px;
}

.protein-fill { background: #5ac8fa; }
.fat-fill { background: #ff9500; }
.carbs-fill { background: #34c759; }

.macro-bar-label {
  font-size: 13px;
  color: var(--hint-color);
  font-weight: 500;
}

/* Meal groups */
.meal-group {
  margin-bottom: 16px;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease forwards;
}

.meal-group:nth-child(2) { animation-delay: 60ms; }
.meal-group:nth-child(3) { animation-delay: 120ms; }
.meal-group:nth-child(4) { animation-delay: 180ms; }

.meal-group-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
  padding: 0 4px;
}

.meal-emoji {
  font-size: 20px;
  line-height: 1;
}

.meal-type-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-color);
  flex: 1;
}

.meal-group-total {
  font-size: 13px;
  color: var(--hint-color);
  font-weight: 500;
}

.entries-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.entry-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: var(--secondary-bg);
  border-radius: 12px;
  transition: transform 0.15s;
}

.entry-card:active {
  transform: scale(0.98);
}

.entry-info {
  display: flex;
  flex-direction: column;
  gap: 3px;
  flex: 1;
  min-width: 0;
}

.entry-name {
  font-size: 15px;
  font-weight: 500;
  color: var(--text-color);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.entry-macros {
  font-size: 12px;
  color: var(--hint-color);
}

.delete-btn {
  background: none;
  border: none;
  font-size: 18px;
  cursor: pointer;
  padding: 6px 8px;
  border-radius: 8px;
  flex-shrink: 0;
  margin-left: 8px;
  transition: background 0.15s;
  line-height: 1;
}

.delete-btn:active {
  background: rgba(255, 59, 48, 0.1);
}

.delete-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.delete-spinner {
  display: inline-block;
  width: 16px;
  height: 16px;
  border: 2px solid var(--hint-color);
  border-top-color: transparent;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Empty state */
.empty-state {
  text-align: center;
  padding: 48px 24px;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease forwards;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 12px;
}

.empty-title {
  font-size: 17px;
  font-weight: 600;
  color: var(--text-color);
  margin-bottom: 4px;
}

.empty-hint {
  font-size: 14px;
  color: var(--hint-color);
}

/* Add form */
.add-form {
  background: var(--secondary-bg);
  border-radius: 14px;
  padding: 20px 16px;
  margin-bottom: 16px;
}

.form-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 16px;
  color: var(--text-color);
}

.form-field {
  margin-bottom: 14px;
  flex: 1;
}

.field-label {
  display: block;
  font-size: 13px;
  color: var(--hint-color);
  margin-bottom: 6px;
  font-weight: 500;
}

.protein-label { color: #5ac8fa; }
.fat-label { color: #ff9500; }
.carbs-label { color: #34c759; }

.form-input {
  width: 100%;
  padding: 11px 12px;
  border: 1.5px solid rgba(128, 128, 128, 0.2);
  border-radius: 10px;
  font-size: 15px;
  background: var(--bg-color);
  color: var(--text-color);
  font-family: inherit;
  transition: border-color 0.2s;
  box-sizing: border-box;
}

.form-input:focus {
  outline: none;
  border-color: var(--button-color);
}

.form-input::placeholder {
  color: var(--hint-color);
}

/* Meal chips */
.meal-chips {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
}

.meal-chip {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 10px 4px;
  border: 2px solid rgba(128, 128, 128, 0.2);
  border-radius: 12px;
  background: var(--bg-color);
  cursor: pointer;
  transition: all 0.2s ease;
  font-family: inherit;
}

.meal-chip.active {
  border-color: var(--button-color);
  background: var(--button-color);
}

.meal-chip.active .chip-text {
  color: var(--button-text-color);
}

.chip-emoji {
  font-size: 20px;
  line-height: 1;
}

.chip-text {
  font-size: 11px;
  font-weight: 500;
  color: var(--text-color);
}

/* 3-column macro row */
.form-row-3 {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
}

/* Form actions */
.form-actions {
  display: flex;
  gap: 10px;
  margin-top: 4px;
}

.btn {
  display: block;
  width: 100%;
  padding: 13px;
  border: none;
  border-radius: 12px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  text-align: center;
  font-family: inherit;
  transition: opacity 0.15s;
}

.btn:active {
  opacity: 0.85;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-primary {
  background: var(--button-color);
  color: var(--button-text-color);
}

.btn-secondary {
  background: var(--bg-color);
  color: var(--text-color);
  border: 1.5px solid rgba(128, 128, 128, 0.2);
}

.error-msg {
  text-align: center;
  color: #ff3b30;
  font-size: 14px;
  margin-top: 10px;
}

/* Floating action button */
.fab {
  position: fixed;
  bottom: 24px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 14px 28px;
  background: var(--button-color);
  color: var(--button-text-color);
  border: none;
  border-radius: 50px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
  transition: transform 0.2s, box-shadow 0.2s;
  font-family: inherit;
  z-index: 10;
}

.fab:active {
  transform: translateX(-50%) scale(0.95);
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.15);
}

.fab-icon {
  font-size: 20px;
  font-weight: 700;
  line-height: 1;
}

.fab-text {
  line-height: 1;
}

/* Form transition */
.form-slide-enter-active {
  transition: all 0.3s ease;
}

.form-slide-leave-active {
  transition: all 0.2s ease;
}

.form-slide-enter-from {
  opacity: 0;
  transform: translateY(16px);
}

.form-slide-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>
