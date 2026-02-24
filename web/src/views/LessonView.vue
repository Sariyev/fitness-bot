<template>
  <div class="lesson-view">
    <div v-if="loading" class="skeleton-list">
      <SkeletonCard v-for="i in 3" :key="i" />
    </div>
    <template v-else-if="lesson">
      <h1 class="page-title">{{ lesson.title }}</h1>
      <p v-if="lesson.description" class="lesson-desc">{{ lesson.description }}</p>

      <div class="content-list">
        <div
          v-for="(item, index) in lesson.content"
          :key="item.id"
          class="content-item"
          :style="{ animationDelay: index * 80 + 'ms' }"
        >
          <ContentBlock :content="item" />
        </div>
      </div>

      <div v-if="lesson.status === 'completed'" class="completed-badge">
        <span>✅ Занятие пройдено</span>
      </div>

      <ConfettiCanvas :active="showConfetti" />
    </template>
    <div v-else class="empty">
      <p>Занятие не найдено</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../api'
import { useTelegram } from '../composables/useTelegram'
import type { LessonDetail } from '../types'
import ContentBlock from '../components/ContentBlock.vue'
import SkeletonCard from '../components/SkeletonCard.vue'
import ConfettiCanvas from '../components/ConfettiCanvas.vue'

const props = defineProps<{ id: string }>()
const router = useRouter()
const { hapticNotification, showMainButton, hideMainButton } = useTelegram()
const lesson = ref<LessonDetail | null>(null)
const loading = ref(true)
const completing = ref(false)
const showConfetti = ref(false)

async function markComplete() {
  if (!lesson.value || completing.value) return
  completing.value = true

  try {
    await api.completeLesson(lesson.value.id)
    lesson.value.status = 'completed'
    hapticNotification('success')
    showConfetti.value = true
    hideMainButton()
  } catch (e) {
    console.error('Failed to complete lesson:', e)
    hapticNotification('error')
  } finally {
    completing.value = false
  }
}

function goBack() {
  router.back()
}

const markCompleteRef = markComplete

onMounted(async () => {
  const backBtn = window.Telegram?.WebApp?.BackButton
  if (backBtn) {
    backBtn.show()
    backBtn.onClick(goBack)
  }

  try {
    const data = await api.getLesson(Number(props.id))
    lesson.value = data

    if (data.status === 'not_started') {
      api.startLesson(data.id).catch(() => {})
    }

    if (data.status !== 'completed') {
      showMainButton('Завершить занятие ✅', markCompleteRef)
    }
  } catch (e) {
    console.error('Failed to load lesson:', e)
  } finally {
    loading.value = false
  }
})

onUnmounted(() => {
  window.Telegram?.WebApp?.BackButton?.offClick(goBack)
  hideMainButton()
})
</script>

<style scoped>
.page-title {
  font-size: 22px;
  font-weight: 700;
  margin-bottom: 8px;
}

.lesson-desc {
  color: var(--hint-color);
  font-size: 14px;
  margin-bottom: 20px;
}

.skeleton-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.content-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.content-item {
  opacity: 0;
  animation: fadeSlideUp 0.4s ease forwards;
}

.completed-badge {
  margin-top: 24px;
  padding: 14px;
  text-align: center;
  background: #34c75920;
  color: #34c759;
  border-radius: 12px;
  font-weight: 600;
  font-size: 16px;
}

.empty {
  text-align: center;
  padding: 40px 0;
  color: var(--hint-color);
}
</style>
