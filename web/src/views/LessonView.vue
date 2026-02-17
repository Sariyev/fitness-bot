<template>
  <div class="lesson-view">
    <div v-if="loading" class="loading">
      <div class="spinner"></div>
    </div>
    <template v-else-if="lesson">
      <h1 class="page-title">{{ lesson.title }}</h1>
      <p v-if="lesson.description" class="lesson-desc">{{ lesson.description }}</p>

      <div class="content-list">
        <ContentBlock
          v-for="item in lesson.content"
          :key="item.id"
          :content="item"
        />
      </div>

      <div v-if="lesson.status !== 'completed'" class="complete-section">
        <button
          class="btn complete-btn"
          :disabled="completing"
          @click="markComplete"
        >
          {{ completing ? 'Сохранение...' : 'Завершить занятие' }}
        </button>
      </div>
      <div v-else class="completed-badge">
        <span>✓ Занятие пройдено</span>
      </div>
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
import type { LessonDetail } from '../types'
import ContentBlock from '../components/ContentBlock.vue'

const props = defineProps<{ id: string }>()
const router = useRouter()
const lesson = ref<LessonDetail | null>(null)
const loading = ref(true)
const completing = ref(false)

async function markComplete() {
  if (!lesson.value || completing.value) return
  completing.value = true

  try {
    await api.completeLesson(lesson.value.id)
    lesson.value.status = 'completed'
    window.Telegram?.WebApp?.HapticFeedback?.notificationOccurred('success')
  } catch (e) {
    console.error('Failed to complete lesson:', e)
    window.Telegram?.WebApp?.HapticFeedback?.notificationOccurred('error')
  } finally {
    completing.value = false
  }
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
    const data = await api.getLesson(Number(props.id))
    lesson.value = data

    // Auto-mark as started
    if (data.status === 'not_started') {
      api.startLesson(data.id).catch(() => {})
    }
  } catch (e) {
    console.error('Failed to load lesson:', e)
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
  font-size: 22px;
  font-weight: 700;
  margin-bottom: 8px;
}

.lesson-desc {
  color: var(--hint-color);
  font-size: 14px;
  margin-bottom: 20px;
}

.content-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.complete-section {
  margin-top: 24px;
  padding-bottom: 24px;
}

.complete-btn {
  width: 100%;
  padding: 14px;
  font-size: 16px;
  font-weight: 600;
}

.complete-btn:disabled {
  opacity: 0.6;
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
