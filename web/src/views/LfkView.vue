<template>
  <div class="lfk-page">
    <!-- Category cards grid -->
    <div class="categories-grid">
      <div
        v-for="cat in categories"
        :key="cat.value"
        class="category-card"
        :class="{ selected: selectedCategory === cat.value }"
        @click="selectCategory(cat.value)"
      >
        <span class="category-icon">{{ cat.icon }}</span>
        <span class="category-name">{{ cat.label }}</span>
        <span class="category-desc">{{ cat.description }}</span>
      </div>
    </div>

    <!-- Courses section (visible when a category is selected) -->
    <template v-if="selectedCategory">
      <h3 class="section-title">Курсы: {{ selectedCategoryLabel }}</h3>

      <!-- Loading courses -->
      <div v-if="loadingCourses" class="skeleton-list">
        <div class="skeleton-card" v-for="i in 3" :key="i">
          <div class="skeleton-line" style="height: 14px; width: 40%"></div>
          <div class="skeleton-line" style="height: 18px; width: 70%; margin-top: 8px"></div>
          <div class="skeleton-line" style="height: 12px; width: 90%; margin-top: 8px"></div>
        </div>
      </div>

      <!-- Error -->
      <div v-else-if="error" class="error-state">
        <p class="error-text">{{ error }}</p>
        <button class="btn-secondary" @click="loadCourses">Попробовать снова</button>
      </div>

      <!-- Empty -->
      <div v-else-if="courses.length === 0" class="empty-state">
        <p class="empty-text">Курсы в этой категории пока не добавлены</p>
      </div>

      <!-- Course cards -->
      <div v-else class="courses-list">
        <router-link
          v-for="(course, index) in courses"
          :key="course.id"
          :to="`/lfk/course/${course.id}`"
          class="course-card"
          :style="{ animationDelay: (index * 60) + 'ms' }"
        >
          <div class="course-card-top">
            <span class="course-name">{{ course.name }}</span>
            <span class="course-badge">14 дней</span>
          </div>
          <p class="course-desc">{{ truncate(course.description, 100) }}</p>
        </router-link>
      </div>
    </template>

    <!-- Warning -->
    <div class="warning-card">
      <span class="warning-icon">&#x26A0;&#xFE0F;</span>
      <p>При острой боли обратитесь к врачу</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { api } from '../api'
import type { RehabCourse } from '../types'

interface CategoryOption {
  value: string
  label: string
  icon: string
  description: string
}

const categories: CategoryOption[] = [
  { value: 'shoulder', label: 'Плечо', icon: '\uD83D\uDCAA', description: 'Восстановление плечевого сустава' },
  { value: 'knee', label: 'Колено', icon: '\uD83E\uDDB5', description: 'Укрепление коленного сустава' },
  { value: 'back', label: 'Грыжа/Спина', icon: '\uD83D\uDD19', description: 'Реабилитация при грыжах и протрузиях' },
  { value: 'scoliosis', label: 'Сколиоз', icon: '\uD83E\uDDB4', description: 'Коррекция искривления позвоночника' },
  { value: 'posture', label: 'Холка/Осанка', icon: '\uD83E\uDDCD', description: 'Формирование правильной осанки' },
]

const selectedCategory = ref<string | null>(null)
const loadingCourses = ref(false)
const error = ref('')
const courses = ref<RehabCourse[]>([])

const selectedCategoryLabel = computed(() => {
  const cat = categories.find(c => c.value === selectedCategory.value)
  return cat?.label || ''
})

function truncate(text: string, maxLen: number): string {
  if (!text) return ''
  if (text.length <= maxLen) return text
  return text.slice(0, maxLen).trimEnd() + '...'
}

async function selectCategory(value: string) {
  if (selectedCategory.value === value) {
    selectedCategory.value = null
    courses.value = []
    return
  }
  selectedCategory.value = value
  await loadCourses()
}

async function loadCourses() {
  if (!selectedCategory.value) return
  loadingCourses.value = true
  error.value = ''
  try {
    courses.value = await api.getRehabCourses(selectedCategory.value)
  } catch (e: any) {
    error.value = e.message || 'Не удалось загрузить курсы'
    courses.value = []
  } finally {
    loadingCourses.value = false
  }
}

onMounted(() => {
  // No initial load -- user picks a category first
})
</script>

<style scoped>
.lfk-page {
  max-width: 480px;
  margin: 0 auto;
  padding-bottom: 24px;
}

/* ===== Category grid ===== */
.categories-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
  margin-bottom: 20px;
}

.category-card {
  background: var(--secondary-bg);
  border-radius: 12px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  cursor: pointer;
  transition: transform 0.15s ease, box-shadow 0.15s ease;
  border: 2px solid transparent;
}

.category-card:active {
  transform: scale(0.97);
}

.category-card.selected {
  border-color: var(--button-color);
  box-shadow: 0 0 0 1px var(--button-color);
}

.category-icon {
  font-size: 28px;
  margin-bottom: 4px;
}

.category-name {
  font-size: 15px;
  font-weight: 600;
}

.category-desc {
  font-size: 12px;
  color: var(--hint-color);
  line-height: 1.3;
}

/* ===== Section title ===== */
.section-title {
  font-size: 18px;
  font-weight: 700;
  margin-bottom: 12px;
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

@keyframes shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

/* ===== Error ===== */
.error-state {
  text-align: center;
  padding: 32px 0;
}

.error-text {
  color: #ff3b30;
  font-size: 14px;
  margin-bottom: 16px;
}

/* ===== Empty ===== */
.empty-state {
  text-align: center;
  padding: 40px 0;
  color: var(--hint-color);
}

.empty-text {
  font-size: 14px;
}

/* ===== Course cards ===== */
.courses-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 20px;
}

.course-card {
  background: var(--secondary-bg);
  border-radius: 12px;
  padding: 16px;
  text-decoration: none;
  color: var(--text-color);
  transition: transform 0.15s ease;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease forwards;
}

.course-card:active {
  transform: scale(0.98);
}

.course-card-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 8px;
}

.course-name {
  font-size: 15px;
  font-weight: 600;
  flex: 1;
  min-width: 0;
}

.course-badge {
  flex-shrink: 0;
  background: var(--button-color);
  color: var(--button-text-color);
  padding: 3px 10px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 600;
}

.course-desc {
  font-size: 13px;
  color: var(--hint-color);
  line-height: 1.4;
}

/* ===== Warning ===== */
.warning-card {
  background: #ff3b3015;
  border: 1px solid #ff3b3040;
  border-radius: 12px;
  padding: 14px 16px;
  display: flex;
  gap: 10px;
  align-items: center;
  margin-top: 8px;
}

.warning-icon {
  font-size: 20px;
  flex-shrink: 0;
}

.warning-card p {
  font-size: 13px;
  color: #ff3b30;
  line-height: 1.4;
  font-weight: 500;
}

/* ===== Buttons ===== */
.btn-secondary {
  display: block;
  width: 100%;
  padding: 14px;
  background: var(--secondary-bg);
  color: var(--text-color);
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
