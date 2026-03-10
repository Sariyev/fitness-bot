<template>
  <div class="calculator-page">
    <!-- Header -->
    <div class="header">
      <button class="back-btn" @click="router.back()">← Назад</button>
      <h1 class="page-title">Калькулятор КБЖУ</h1>
    </div>

    <!-- Input form -->
    <div class="calc-form">
      <!-- Gender toggle -->
      <div class="form-field">
        <label class="field-label">Пол</label>
        <div class="gender-toggle">
          <button
            class="toggle-btn"
            :class="{ active: gender === 'male' }"
            @click="gender = 'male'"
          >
            Мужской
          </button>
          <button
            class="toggle-btn"
            :class="{ active: gender === 'female' }"
            @click="gender = 'female'"
          >
            Женский
          </button>
        </div>
      </div>

      <!-- Weight -->
      <div class="form-field">
        <label class="field-label">Вес (кг)</label>
        <input
          v-model.number="weight"
          type="number"
          inputmode="numeric"
          placeholder="70"
          min="30"
          max="300"
          class="form-input"
        />
      </div>

      <!-- Height -->
      <div class="form-field">
        <label class="field-label">Рост (см)</label>
        <input
          v-model.number="height"
          type="number"
          inputmode="numeric"
          placeholder="170"
          min="100"
          max="250"
          class="form-input"
        />
      </div>

      <!-- Age -->
      <div class="form-field">
        <label class="field-label">Возраст</label>
        <input
          v-model.number="age"
          type="number"
          inputmode="numeric"
          placeholder="25"
          min="14"
          max="100"
          class="form-input"
        />
      </div>

      <!-- Goal chips -->
      <div class="form-field">
        <label class="field-label">Цель</label>
        <div class="goal-chips">
          <button
            v-for="g in goals"
            :key="g.value"
            class="goal-chip"
            :class="{ active: goal === g.value }"
            @click="goal = g.value"
          >
            <span class="chip-emoji">{{ g.emoji }}</span>
            <span class="chip-label">{{ g.label }}</span>
          </button>
        </div>
      </div>

      <!-- Calculate button -->
      <button
        class="btn btn-primary"
        @click="calculate"
        :disabled="calculating || !isFormValid"
      >
        <span v-if="calculating" class="btn-spinner"></span>
        {{ calculating ? 'Считаем...' : 'Рассчитать' }}
      </button>

      <div v-if="error" class="error-msg">{{ error }}</div>
    </div>

    <!-- Results section -->
    <Transition name="results-fade">
      <div v-if="results" class="results-section">
        <h2 class="results-heading">Ваша норма</h2>

        <!-- Large calorie display -->
        <div class="calorie-card">
          <span class="calorie-number">{{ results.calories }}</span>
          <span class="calorie-label">ккал/день</span>
        </div>

        <!-- 3 macro cards -->
        <div class="macro-cards">
          <div class="macro-card protein-card">
            <div class="macro-accent protein-accent"></div>
            <span class="macro-value">{{ results.protein }}г</span>
            <span class="macro-name">Белки</span>
            <span class="macro-pct">{{ proteinPct }}%</span>
          </div>
          <div class="macro-card fat-card">
            <div class="macro-accent fat-accent"></div>
            <span class="macro-value">{{ results.fat }}г</span>
            <span class="macro-name">Жиры</span>
            <span class="macro-pct">{{ fatPct }}%</span>
          </div>
          <div class="macro-card carbs-card">
            <div class="macro-accent carbs-accent"></div>
            <span class="macro-value">{{ results.carbs }}г</span>
            <span class="macro-name">Углеводы</span>
            <span class="macro-pct">{{ carbsPct }}%</span>
          </div>
        </div>

        <!-- Macro distribution bar -->
        <div class="distribution-bar">
          <div class="dist-segment protein-seg" :style="{ width: proteinPct + '%' }"></div>
          <div class="dist-segment fat-seg" :style="{ width: fatPct + '%' }"></div>
          <div class="dist-segment carbs-seg" :style="{ width: carbsPct + '%' }"></div>
        </div>
        <div class="distribution-legend">
          <span class="legend-item"><span class="legend-dot protein-dot"></span>Белки</span>
          <span class="legend-item"><span class="legend-dot fat-dot"></span>Жиры</span>
          <span class="legend-item"><span class="legend-dot carbs-dot"></span>Углеводы</span>
        </div>

        <!-- Explanation -->
        <p class="explanation">
          Рассчитано по формуле Миффлина-Сан Жеора с учетом вашей цели и уровня активности
        </p>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../api'
import type { MacroTargets } from '../types'

const router = useRouter()

const gender = ref('male')
const weight = ref<number | null>(null)
const height = ref<number | null>(null)
const age = ref<number | null>(null)
const goal = ref('maintenance')
const calculating = ref(false)
const error = ref('')
const results = ref<MacroTargets | null>(null)

interface GoalOption {
  value: string
  label: string
  emoji: string
}

const goals: GoalOption[] = [
  { value: 'weight_loss', label: 'Похудение', emoji: '🔥' },
  { value: 'maintenance', label: 'Поддержка', emoji: '⚖️' },
  { value: 'muscle_gain', label: 'Набор массы', emoji: '💪' },
]

const isFormValid = computed(() => {
  return (
    weight.value !== null && weight.value > 0 &&
    height.value !== null && height.value > 0 &&
    age.value !== null && age.value > 0 &&
    gender.value !== '' &&
    goal.value !== ''
  )
})

// Macro percentage calculations
const totalMacroCals = computed(() => {
  if (!results.value) return 0
  return results.value.protein * 4 + results.value.fat * 9 + results.value.carbs * 4
})

const proteinPct = computed(() => {
  if (!totalMacroCals.value) return 0
  return Math.round((results.value!.protein * 4 / totalMacroCals.value) * 100)
})

const fatPct = computed(() => {
  if (!totalMacroCals.value) return 0
  return Math.round((results.value!.fat * 9 / totalMacroCals.value) * 100)
})

const carbsPct = computed(() => {
  if (!totalMacroCals.value) return 0
  return 100 - proteinPct.value - fatPct.value
})

async function calculate(): Promise<void> {
  if (calculating.value) return
  if (!isFormValid.value) {
    error.value = 'Заполните все поля'
    return
  }

  calculating.value = true
  error.value = ''
  results.value = null

  try {
    results.value = await api.calculateMacros({
      gender: gender.value,
      weight_kg: weight.value!,
      height_cm: height.value!,
      age: age.value!,
      goal: goal.value,
    })
  } catch (e: any) {
    error.value = e.message || 'Ошибка при расчете'
  } finally {
    calculating.value = false
  }
}
</script>

<style scoped>
.calculator-page {
  max-width: 400px;
  margin: 0 auto;
  padding-bottom: 40px;
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

/* Form */
.calc-form {
  background: var(--secondary-bg);
  border-radius: 14px;
  padding: 20px 16px;
  margin-bottom: 24px;
}

.form-field {
  margin-bottom: 16px;
}

.field-label {
  display: block;
  font-size: 13px;
  color: var(--hint-color);
  margin-bottom: 6px;
  font-weight: 500;
}

.form-input {
  width: 100%;
  padding: 12px 14px;
  border: 1.5px solid rgba(128, 128, 128, 0.2);
  border-radius: 10px;
  font-size: 16px;
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

/* Gender toggle */
.gender-toggle {
  display: flex;
  gap: 8px;
}

.toggle-btn {
  flex: 1;
  padding: 12px;
  border: 2px solid rgba(128, 128, 128, 0.2);
  border-radius: 12px;
  background: var(--bg-color);
  color: var(--text-color);
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  font-family: inherit;
}

.toggle-btn.active {
  border-color: var(--button-color);
  background: var(--button-color);
  color: var(--button-text-color);
}

.toggle-btn:active {
  transform: scale(0.97);
}

/* Goal chips */
.goal-chips {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
}

.goal-chip {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 14px 8px;
  border: 2px solid rgba(128, 128, 128, 0.2);
  border-radius: 12px;
  background: var(--bg-color);
  cursor: pointer;
  transition: all 0.2s ease;
  font-family: inherit;
}

.goal-chip.active {
  border-color: var(--button-color);
  background: var(--button-color);
}

.goal-chip.active .chip-label {
  color: var(--button-text-color);
}

.goal-chip:active {
  transform: scale(0.95);
}

.chip-emoji {
  font-size: 24px;
  line-height: 1;
}

.chip-label {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-color);
}

/* Button */
.btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  padding: 14px;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  text-align: center;
  font-family: inherit;
  transition: opacity 0.15s;
  margin-top: 4px;
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

.btn-spinner {
  display: inline-block;
  width: 18px;
  height: 18px;
  border: 2.5px solid var(--button-text-color);
  border-top-color: transparent;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.error-msg {
  text-align: center;
  color: #ff3b30;
  font-size: 14px;
  margin-top: 10px;
}

/* Results */
.results-section {
  margin-bottom: 24px;
}

.results-heading {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 14px;
  color: var(--text-color);
}

/* Calorie card */
.calorie-card {
  background: var(--button-color);
  border-radius: 16px;
  padding: 28px 20px;
  text-align: center;
  margin-bottom: 14px;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease forwards;
}

.calorie-number {
  display: block;
  font-size: 48px;
  font-weight: 800;
  color: var(--button-text-color);
  line-height: 1;
}

.calorie-label {
  display: block;
  font-size: 15px;
  color: var(--button-text-color);
  opacity: 0.8;
  margin-top: 6px;
  font-weight: 500;
}

/* Macro cards */
.macro-cards {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
  margin-bottom: 16px;
}

.macro-card {
  background: var(--secondary-bg);
  border-radius: 14px;
  padding: 16px 10px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  position: relative;
  overflow: hidden;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease forwards;
}

.macro-card:nth-child(1) { animation-delay: 80ms; }
.macro-card:nth-child(2) { animation-delay: 160ms; }
.macro-card:nth-child(3) { animation-delay: 240ms; }

.macro-accent {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  border-radius: 4px 4px 0 0;
}

.protein-accent { background: #5ac8fa; }
.fat-accent { background: #ff9500; }
.carbs-accent { background: #34c759; }

.macro-value {
  font-size: 22px;
  font-weight: 700;
  margin-top: 4px;
}

.protein-card .macro-value { color: #5ac8fa; }
.fat-card .macro-value { color: #ff9500; }
.carbs-card .macro-value { color: #34c759; }

.macro-name {
  font-size: 13px;
  color: var(--hint-color);
  font-weight: 500;
}

.macro-pct {
  font-size: 12px;
  color: var(--hint-color);
  opacity: 0.7;
}

/* Distribution bar */
.distribution-bar {
  display: flex;
  height: 8px;
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 8px;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease 0.3s forwards;
}

.dist-segment {
  transition: width 0.5s ease;
}

.protein-seg { background: #5ac8fa; }
.fat-seg { background: #ff9500; }
.carbs-seg { background: #34c759; }

.distribution-legend {
  display: flex;
  justify-content: center;
  gap: 16px;
  margin-bottom: 20px;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease 0.35s forwards;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 12px;
  color: var(--hint-color);
}

.legend-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.protein-dot { background: #5ac8fa; }
.fat-dot { background: #ff9500; }
.carbs-dot { background: #34c759; }

/* Explanation */
.explanation {
  font-size: 13px;
  color: var(--hint-color);
  text-align: center;
  line-height: 1.5;
  padding: 12px 16px;
  background: var(--secondary-bg);
  border-radius: 12px;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease 0.4s forwards;
}

/* Results transition */
.results-fade-enter-active {
  transition: all 0.4s ease;
}

.results-fade-leave-active {
  transition: all 0.2s ease;
}

.results-fade-enter-from {
  opacity: 0;
  transform: translateY(20px);
}

.results-fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
