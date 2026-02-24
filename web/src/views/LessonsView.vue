<template>
  <div class="lessons-view">
    <h1 class="page-title">Занятия 📝</h1>
    <div v-if="loading" class="skeleton-list">
      <SkeletonCard v-for="i in 5" :key="i" />
    </div>
    <div v-else-if="lessons.length === 0" class="empty">
      <p>Занятия пока не доступны</p>
    </div>
    <div v-else class="lesson-list">
      <div
        v-for="(lesson, index) in lessons"
        :key="lesson.id"
        class="lesson-card"
        :style="{ animationDelay: index * 50 + 'ms' }"
        @click="openLesson(lesson.id)"
      >
        <div class="lesson-number" :class="lesson.status">
          <span v-if="lesson.status === 'completed'">✓</span>
          <span v-else>{{ index + 1 }}</span>
        </div>
        <div class="lesson-info">
          <h3>{{ lesson.title }}</h3>
          <p v-if="lesson.description">{{ lesson.description }}</p>
        </div>
        <div class="lesson-status">
          <span v-if="lesson.status === 'completed'" class="status-badge completed">Пройдено ✅</span>
          <span v-else-if="lesson.status === 'in_progress'" class="status-badge in-progress">В процессе ⏳</span>
          <span class="lesson-arrow">›</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../api'
import { useTelegram } from '../composables/useTelegram'
import type { LessonWithProgress } from '../types'
import SkeletonCard from '../components/SkeletonCard.vue'

const props = defineProps<{ id: string }>()
const router = useRouter()
const { hapticSelection } = useTelegram()
const lessons = ref<LessonWithProgress[]>([])
const loading = ref(true)

function openLesson(id: number) {
  hapticSelection()
  router.push({ name: 'lesson', params: { id } })
}

function goBack() {
  router.back()
}

onMounted(async () => {
  const backBtn = window.Telegram?.WebApp?.BackButton
  if (backBtn) {
    backBtn.show()
    backBtn.onClick(goBack)
  }

  try {
    lessons.value = await api.getLessons(Number(props.id))
  } catch (e) {
    console.error('Failed to load lessons:', e)
  } finally {
    loading.value = false
  }
})

onUnmounted(() => {
  window.Telegram?.WebApp?.BackButton?.offClick(goBack)
})
</script>

<style scoped>
.page-title {
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 16px;
}

.skeleton-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.lesson-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.lesson-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  background: var(--secondary-bg);
  border-radius: 12px;
  cursor: pointer;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease forwards;
  transition: transform 0.15s ease;
}

.lesson-card:active {
  transform: scale(0.98);
}

.lesson-number {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
  background: var(--button-color);
  color: var(--button-text-color);
  flex-shrink: 0;
}

.lesson-number.completed {
  background: #34c759;
}

.lesson-number.in_progress {
  background: #ff9500;
  animation: pulse-status 1.5s ease-in-out infinite;
}

@keyframes pulse-status {
  0%, 100% { transform: scale(1); opacity: 1; }
  50% { transform: scale(1.08); opacity: 0.85; }
}

.lesson-info {
  flex: 1;
  min-width: 0;
}

.lesson-info h3 {
  font-size: 15px;
  font-weight: 600;
}

.lesson-info p {
  font-size: 13px;
  color: var(--hint-color);
  margin-top: 2px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.lesson-status {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.status-badge {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 10px;
  font-weight: 500;
}

.status-badge.completed {
  background: #34c75920;
  color: #34c759;
}

.status-badge.in-progress {
  background: #ff950020;
  color: #ff9500;
}

.lesson-arrow {
  font-size: 20px;
  color: var(--hint-color);
}

.empty {
  text-align: center;
  padding: 40px 0;
  color: var(--hint-color);
}
</style>
