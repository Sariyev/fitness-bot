<template>
  <div class="admin-page">
    <button class="back-btn" @click="router.push('/admin')">← Назад</button>
    <h1 class="page-title">Отзывы</h1>

    <div v-if="loading" class="loading">Загрузка...</div>

    <div v-else-if="summary" class="summary-card">
      <div class="stars">
        <span v-for="i in 5" :key="i" class="star" :class="{ filled: i <= Math.round(summary.average_score) }">★</span>
      </div>
      <div class="score-text">{{ summary.average_score.toFixed(1) }} / 5</div>
      <div class="review-count">{{ summary.total_reviews }} {{ reviewWord(summary.total_reviews) }}</div>
    </div>

    <div v-else class="empty">Отзывов пока нет</div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../api'
import type { ReviewSummary } from '../../types'

const router = useRouter()
const summary = ref<ReviewSummary | null>(null)
const loading = ref(true)

function reviewWord(n: number): string {
  if (n % 10 === 1 && n % 100 !== 11) return 'отзыв'
  if ([2, 3, 4].includes(n % 10) && ![12, 13, 14].includes(n % 100)) return 'отзыва'
  return 'отзывов'
}

onMounted(async () => {
  try {
    summary.value = await api.getAdminReviewsSummary()
  } catch {
    summary.value = null
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.admin-page {
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

.page-title {
  font-size: 20px;
  font-weight: 700;
  margin-bottom: 16px;
}

.loading, .empty {
  text-align: center;
  color: var(--hint-color);
  padding: 40px;
}

.summary-card {
  background: var(--secondary-bg);
  border-radius: 12px;
  padding: 32px 20px;
  text-align: center;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease forwards;
}

.stars {
  font-size: 36px;
  margin-bottom: 8px;
}

.star {
  color: #d0d0d0;
}

.star.filled {
  color: #ffb800;
}

.score-text {
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 4px;
}

.review-count {
  color: var(--hint-color);
  font-size: 15px;
}

@keyframes fadeSlideUp {
  from { opacity: 0; transform: translateY(12px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
