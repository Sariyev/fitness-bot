<template>
  <div class="nutrition-page">
    <!-- Loading -->
    <div v-if="loading" class="skeleton-list">
      <div class="skeleton-card" style="height: 120px"></div>
      <div class="skeleton-card" style="height: 90px" v-for="i in 3" :key="i"></div>
    </div>

    <template v-else>
      <!-- Daily macro summary -->
      <section v-if="macroTargets" class="macro-section">
        <h3 class="section-title">Дневные макросы</h3>
        <div class="macro-bars">
          <div class="macro-bar-item">
            <div class="macro-bar-header">
              <span class="macro-label">Калории</span>
              <span class="macro-value">{{ dailySummary.calories }} / {{ macroTargets.calories }} ккал</span>
            </div>
            <div class="bar-track">
              <div
                class="bar-fill bar-calories"
                :style="{ width: macroPercent(dailySummary.calories, macroTargets.calories) + '%' }"
              ></div>
            </div>
          </div>
          <div class="macro-bar-item">
            <div class="macro-bar-header">
              <span class="macro-label">Белки</span>
              <span class="macro-value">{{ dailySummary.protein }} / {{ macroTargets.protein }} г</span>
            </div>
            <div class="bar-track">
              <div
                class="bar-fill bar-protein"
                :style="{ width: macroPercent(dailySummary.protein, macroTargets.protein) + '%' }"
              ></div>
            </div>
          </div>
          <div class="macro-bar-item">
            <div class="macro-bar-header">
              <span class="macro-label">Жиры</span>
              <span class="macro-value">{{ dailySummary.fat }} / {{ macroTargets.fat }} г</span>
            </div>
            <div class="bar-track">
              <div
                class="bar-fill bar-fat"
                :style="{ width: macroPercent(dailySummary.fat, macroTargets.fat) + '%' }"
              ></div>
            </div>
          </div>
          <div class="macro-bar-item">
            <div class="macro-bar-header">
              <span class="macro-label">Углеводы</span>
              <span class="macro-value">{{ dailySummary.carbs }} / {{ macroTargets.carbs }} г</span>
            </div>
            <div class="bar-track">
              <div
                class="bar-fill bar-carbs"
                :style="{ width: macroPercent(dailySummary.carbs, macroTargets.carbs) + '%' }"
              ></div>
            </div>
          </div>
        </div>
      </section>

      <!-- Meal plans section -->
      <section class="section">
        <h3 class="section-title">Планы питания</h3>

        <div v-if="mealPlans.length === 0" class="empty-state">
          <p class="empty-text">Планы питания пока не добавлены</p>
        </div>

        <div v-else class="plans-list">
          <div
            v-for="(plan, index) in mealPlans"
            :key="plan.id"
            class="plan-card"
            :style="{ animationDelay: (index * 60) + 'ms' }"
          >
            <div class="plan-card-top" @click="togglePlan(plan.id)">
              <div class="plan-card-info">
                <span class="plan-name">{{ plan.name }}</span>
                <div class="plan-meta">
                  <span v-if="plan.goal" class="plan-goal-badge">{{ goalLabel(plan.goal) }}</span>
                  <span class="plan-kcal">{{ plan.calories }} ккал</span>
                </div>
              </div>
              <span class="plan-toggle" :class="{ open: expandedPlan === plan.id }">&rsaquo;</span>
            </div>

            <!-- Expanded meals -->
            <Transition name="meals-slide">
              <div v-if="expandedPlan === plan.id" class="meals-list">
                <div v-if="loadingMeals" class="meals-loading">Загрузка...</div>
                <template v-else-if="planMeals.length > 0">
                  <div
                    v-for="meal in planMeals"
                    :key="meal.id"
                    class="meal-item"
                  >
                    <span class="meal-type-badge">{{ mealTypeLabel(meal.meal_type) }}</span>
                    <div class="meal-info">
                      <span class="meal-name">{{ meal.name }}</span>
                      <span class="meal-kcal">{{ meal.calories }} ккал</span>
                    </div>
                  </div>
                </template>
                <p v-else class="meals-empty">Блюда не найдены</p>
              </div>
            </Transition>
          </div>
        </div>
      </section>

      <!-- Bottom action buttons -->
      <div class="action-buttons">
        <router-link to="/nutrition/diary" class="action-btn">
          <span class="action-icon">&#x1F4D3;</span>
          <span class="action-label">Дневник питания</span>
        </router-link>
        <router-link to="/nutrition/calculator" class="action-btn">
          <span class="action-icon">&#x1F9EE;</span>
          <span class="action-label">Калькулятор БЖУ</span>
        </router-link>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '../api'
import type { MealPlan, Meal, DailySummary, MacroTargets } from '../types'

const loading = ref(true)
const mealPlans = ref<MealPlan[]>([])
const macroTargets = ref<MacroTargets | null>(null)
const dailySummary = ref<DailySummary>({ calories: 0, protein: 0, fat: 0, carbs: 0 })

// Expanded plan
const expandedPlan = ref<number | null>(null)
const loadingMeals = ref(false)
const planMeals = ref<Meal[]>([])

const goalLabels: Record<string, string> = {
  weight_loss: 'Похудение',
  muscle_gain: 'Набор массы',
  maintenance: 'Поддержание',
  strength: 'Сила',
}

const mealTypeLabels: Record<string, string> = {
  breakfast: 'Завтрак',
  lunch: 'Обед',
  dinner: 'Ужин',
  snack: 'Перекус',
}

function goalLabel(key: string): string {
  return goalLabels[key] || key
}

function mealTypeLabel(key: string): string {
  return mealTypeLabels[key] || key
}

function macroPercent(current: number, target: number): number {
  if (!target || target <= 0) return 0
  return Math.min(100, Math.round((current / target) * 100))
}

function todayString(): string {
  const d = new Date()
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

async function togglePlan(planId: number) {
  if (expandedPlan.value === planId) {
    expandedPlan.value = null
    planMeals.value = []
    return
  }

  expandedPlan.value = planId
  loadingMeals.value = true
  try {
    const data = await api.getMealPlan(planId)
    planMeals.value = data.meals || []
  } catch {
    planMeals.value = []
  } finally {
    loadingMeals.value = false
  }
}

async function loadInitialData() {
  loading.value = true
  try {
    const [plans, summary] = await Promise.all([
      api.getMealPlans(),
      api.getFoodLogSummary(todayString()).catch(() => ({ calories: 0, protein: 0, fat: 0, carbs: 0 })),
    ])
    mealPlans.value = plans

    dailySummary.value = summary

    // Try to load user profile for macro targets
    try {
      const profile = await api.getProfile()
      if (profile.goals && profile.goals.length > 0) {
        const targets = await api.calculateMacros({
          gender: profile.gender,
          weight_kg: profile.weight_kg,
          height_cm: profile.height_cm,
          age: profile.age,
          goal: profile.goals[0],
        })
        macroTargets.value = targets
      }
    } catch {
      // Macro targets are optional -- continue without them
      macroTargets.value = null
    }
  } catch {
    mealPlans.value = []
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadInitialData()
})
</script>

<style scoped>
.nutrition-page {
  max-width: 480px;
  margin: 0 auto;
  padding-bottom: 24px;
}

/* ===== Skeleton ===== */
.skeleton-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
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

/* ===== Section ===== */
.section {
  margin-bottom: 24px;
}

.section-title {
  font-size: 18px;
  font-weight: 700;
  margin-bottom: 12px;
}

/* ===== Macro summary ===== */
.macro-section {
  margin-bottom: 24px;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease forwards;
}

.macro-bars {
  background: var(--secondary-bg);
  border-radius: 12px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.macro-bar-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.macro-bar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.macro-label {
  font-size: 13px;
  font-weight: 500;
}

.macro-value {
  font-size: 12px;
  color: var(--hint-color);
}

.bar-track {
  height: 8px;
  background: var(--bg-color);
  border-radius: 4px;
  overflow: hidden;
}

.bar-fill {
  height: 100%;
  border-radius: 4px;
  transition: width 0.6s ease;
  min-width: 2px;
}

.bar-calories { background: #ff9500; }
.bar-protein { background: #5856d6; }
.bar-fat { background: #ff3b30; }
.bar-carbs { background: #34c759; }

/* ===== Plans list ===== */
.plans-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.plan-card {
  background: var(--secondary-bg);
  border-radius: 12px;
  overflow: hidden;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease forwards;
}

.plan-card-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  cursor: pointer;
  transition: background 0.15s ease;
}

.plan-card-top:active {
  background: var(--bg-color);
}

.plan-card-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.plan-name {
  font-size: 15px;
  font-weight: 600;
}

.plan-meta {
  display: flex;
  gap: 8px;
  align-items: center;
}

.plan-goal-badge {
  background: var(--button-color);
  color: var(--button-text-color);
  padding: 2px 10px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 600;
}

.plan-kcal {
  font-size: 13px;
  color: var(--hint-color);
}

.plan-toggle {
  flex-shrink: 0;
  font-size: 22px;
  color: var(--hint-color);
  transition: transform 0.2s ease;
}

.plan-toggle.open {
  transform: rotate(90deg);
}

/* ===== Meals (expanded) ===== */
.meals-list {
  padding: 0 16px 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.meals-loading {
  font-size: 13px;
  color: var(--hint-color);
  text-align: center;
  padding: 12px 0;
}

.meals-empty {
  font-size: 13px;
  color: var(--hint-color);
  text-align: center;
  padding: 8px 0;
}

.meal-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  background: var(--bg-color);
  border-radius: 10px;
}

.meal-type-badge {
  flex-shrink: 0;
  background: var(--button-color);
  color: var(--button-text-color);
  padding: 2px 8px;
  border-radius: 8px;
  font-size: 11px;
  font-weight: 600;
}

.meal-info {
  flex: 1;
  min-width: 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
}

.meal-name {
  font-size: 14px;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.meal-kcal {
  flex-shrink: 0;
  font-size: 12px;
  color: var(--hint-color);
}

.meals-slide-enter-active,
.meals-slide-leave-active {
  transition: all 0.25s ease;
  overflow: hidden;
}

.meals-slide-enter-from,
.meals-slide-leave-to {
  opacity: 0;
  max-height: 0;
  padding-top: 0;
  padding-bottom: 0;
}

/* ===== Empty ===== */
.empty-state {
  text-align: center;
  padding: 32px 16px;
  color: var(--hint-color);
}

.empty-text {
  font-size: 14px;
}

/* ===== Action buttons ===== */
.action-buttons {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-top: 8px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: var(--secondary-bg);
  border-radius: 12px;
  text-decoration: none;
  color: var(--text-color);
  font-size: 15px;
  font-weight: 600;
  transition: transform 0.15s ease;
}

.action-btn:active {
  transform: scale(0.98);
}

.action-icon {
  font-size: 22px;
}

.action-label {
  flex: 1;
}

@keyframes fadeSlideUp {
  from { opacity: 0; transform: translateY(12px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
