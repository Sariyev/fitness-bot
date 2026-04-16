<template>
  <div class="review-page">
    <button class="back-btn" @click="router.back()">← Назад</button>

    <!-- Bot review summary -->
    <div class="summary-card" v-if="summary">
      <div class="summary-stars">
        <span v-for="i in 5" :key="i" class="star" :class="{ filled: i <= Math.round(summary.average_score) }">★</span>
      </div>
      <div class="summary-text">
        {{ summary.average_score.toFixed(1) }} / 5
        <span class="review-count">({{ summary.total_reviews }} {{ reviewWord(summary.total_reviews) }})</span>
      </div>
    </div>

    <!-- Write review form -->
    <div class="review-form" v-if="!submitted">
      <h2>Оставить отзыв</h2>
      <p class="subtitle">{{ isBot ? 'Как вам наш бот?' : 'Оцените контент' }}</p>

      <!-- Star rating -->
      <div class="star-selector">
        <button
          v-for="i in 5"
          :key="i"
          class="star-btn"
          :class="{ active: i <= selectedScore }"
          @click="selectedScore = i"
        >
          ★
        </button>
      </div>
      <p v-if="selectedScore" class="score-label">{{ scoreLabels[selectedScore] }}</p>

      <!-- Tags -->
      <div class="tags-section" v-if="availableTags.length">
        <p class="tags-label">Выберите теги:</p>
        <div class="tags-list">
          <button
            v-for="tag in availableTags"
            :key="tag"
            class="tag-btn"
            :class="{ selected: selectedTags.includes(tag) }"
            @click="toggleTag(tag)"
          >
            {{ tag }}
          </button>
        </div>
      </div>

      <!-- Comment -->
      <div class="comment-section">
        <textarea
          v-model="comment"
          placeholder="Комментарий (необязательно)"
          rows="3"
          class="comment-input"
        ></textarea>
      </div>

      <!-- Submit -->
      <button
        class="btn btn-primary"
        :disabled="!selectedScore || submitting"
        @click="submitReview"
      >
        {{ submitting ? 'Отправка...' : 'Отправить отзыв' }}
      </button>
      <p v-if="error" class="error-msg">{{ error }}</p>
    </div>

    <!-- Success state -->
    <div class="success-card" v-else>
      <div class="success-icon">✓</div>
      <h3>Спасибо за отзыв!</h3>
      <div class="submitted-stars">
        <span v-for="i in 5" :key="i" class="star" :class="{ filled: i <= selectedScore }">★</span>
      </div>
      <div v-if="selectedTags.length" class="submitted-tags">
        <span v-for="tag in selectedTags" :key="tag" class="submitted-tag">{{ tag }}</span>
      </div>
      <p v-if="comment" class="submitted-comment">"{{ comment }}"</p>
      <button class="btn btn-secondary" @click="router.push('/')">На главную</button>
    </div>

    <!-- Recent reviews -->
    <div class="reviews-list" v-if="reviews.length">
      <h3>Последние отзывы</h3>
      <div
        v-for="review in reviews"
        :key="review.id"
        class="review-card"
      >
        <div class="review-header">
          <div class="review-stars">
            <span v-for="i in 5" :key="i" class="star-sm" :class="{ filled: i <= review.score }">★</span>
          </div>
          <span class="review-date">{{ formatDate(review.created_at) }}</span>
        </div>
        <div v-if="review.tags && review.tags.length" class="review-tags">
          <span v-for="tag in review.tags" :key="tag" class="review-tag">{{ tag }}</span>
        </div>
        <p v-if="review.comment" class="review-comment">{{ review.comment }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api } from '../api'
import type { Review, ReviewSummary } from '../types'

const route = useRoute()
const router = useRouter()

const referenceType = computed(() => (route.query.type as string) || 'bot')
const referenceId = computed(() => Number(route.query.id) || 0)
const isBot = computed(() => referenceType.value === 'bot')

const summary = ref<ReviewSummary | null>(null)
const reviews = ref<Review[]>([])
const availableTags = ref<string[]>([])
const selectedScore = ref(0)
const selectedTags = ref<string[]>([])
const comment = ref('')
const submitting = ref(false)
const submitted = ref(false)
const error = ref('')

const scoreLabels: Record<number, string> = {
  1: 'Плохо',
  2: 'Ниже среднего',
  3: 'Нормально',
  4: 'Хорошо',
  5: 'Отлично!',
}

function toggleTag(tag: string) {
  const idx = selectedTags.value.indexOf(tag)
  if (idx >= 0) {
    selectedTags.value.splice(idx, 1)
  } else {
    selectedTags.value.push(tag)
  }
}

async function submitReview() {
  if (!selectedScore.value || submitting.value) return
  submitting.value = true
  error.value = ''

  try {
    await api.createReview({
      reference_type: referenceType.value,
      reference_id: referenceId.value,
      score: selectedScore.value,
      comment: comment.value || undefined,
      tags: selectedTags.value.length ? selectedTags.value : undefined,
    })
    submitted.value = true
    // Refresh summary
    await loadSummary()
    await loadReviews()
  } catch (e: any) {
    error.value = e.message || 'Ошибка при отправке'
  } finally {
    submitting.value = false
  }
}

function reviewWord(n: number): string {
  if (n % 10 === 1 && n % 100 !== 11) return 'отзыв'
  if ([2, 3, 4].includes(n % 10) && ![12, 13, 14].includes(n % 100)) return 'отзыва'
  return 'отзывов'
}

function formatDate(dateStr: string): string {
  const d = new Date(dateStr)
  return d.toLocaleDateString('ru-RU', { day: 'numeric', month: 'short' })
}

async function loadSummary() {
  try {
    if (isBot.value) {
      summary.value = await api.getBotReviewSummary()
    } else {
      summary.value = await api.getReviewSummary(referenceType.value, referenceId.value)
    }
  } catch {
    // ignore
  }
}

async function loadReviews() {
  try {
    if (isBot.value) {
      reviews.value = await api.getReviews('bot', 0)
    } else {
      reviews.value = await api.getReviews(referenceType.value, referenceId.value)
    }
  } catch {
    // ignore
  }
}

onMounted(async () => {
  await Promise.all([
    loadSummary(),
    loadReviews(),
    api.getReviewTags(referenceType.value).then(res => {
      availableTags.value = res.tags
    }).catch(() => {}),
  ])
})
</script>

<style scoped>
.review-page {
  max-width: 400px;
  margin: 0 auto;
  padding-bottom: 24px;
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

.summary-card {
  background: var(--secondary-bg);
  border-radius: 12px;
  padding: 20px;
  text-align: center;
  margin-bottom: 16px;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease forwards;
}

.summary-stars {
  font-size: 28px;
  margin-bottom: 4px;
}

.star {
  color: #d0d0d0;
}

.star.filled {
  color: #ffb800;
}

.summary-text {
  font-size: 16px;
  font-weight: 600;
}

.review-count {
  color: var(--hint-color);
  font-weight: 400;
  font-size: 14px;
}

.review-form {
  background: var(--secondary-bg);
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 16px;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease forwards;
  animation-delay: 80ms;
}

.review-form h2 {
  margin: 0 0 4px;
  font-size: 18px;
}

.subtitle {
  color: var(--hint-color);
  font-size: 14px;
  margin: 0 0 16px;
}

.star-selector {
  display: flex;
  justify-content: center;
  gap: 8px;
  margin-bottom: 8px;
}

.star-btn {
  background: none;
  border: none;
  font-size: 36px;
  color: #d0d0d0;
  cursor: pointer;
  padding: 4px;
  transition: color 0.15s, transform 0.15s;
}

.star-btn.active {
  color: #ffb800;
  transform: scale(1.15);
}

.score-label {
  text-align: center;
  color: var(--hint-color);
  font-size: 14px;
  margin: 0 0 16px;
}

.tags-section {
  margin-bottom: 16px;
}

.tags-label {
  font-size: 14px;
  color: var(--hint-color);
  margin: 0 0 8px;
}

.tags-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tag-btn {
  background: var(--bg-color);
  border: 1px solid rgba(0, 0, 0, 0.1);
  border-radius: 16px;
  padding: 6px 14px;
  font-size: 13px;
  cursor: pointer;
  color: var(--text-color);
  transition: all 0.15s;
}

.tag-btn.selected {
  background: var(--button-color);
  color: var(--button-text-color);
  border-color: var(--button-color);
}

.comment-section {
  margin-bottom: 16px;
}

.comment-input {
  width: 100%;
  background: var(--bg-color);
  border: 1px solid rgba(0, 0, 0, 0.1);
  border-radius: 10px;
  padding: 12px;
  font-size: 14px;
  color: var(--text-color);
  resize: none;
  box-sizing: border-box;
  font-family: inherit;
}

.comment-input::placeholder {
  color: var(--hint-color);
}

.btn {
  display: block;
  width: 100%;
  padding: 14px;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  cursor: pointer;
  text-align: center;
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
  background: var(--secondary-bg);
  color: var(--text-color);
  margin-top: 12px;
}

.error-msg {
  color: #ff3b30;
  text-align: center;
  font-size: 14px;
  margin-top: 8px;
}

.success-card {
  background: var(--secondary-bg);
  border-radius: 12px;
  padding: 28px 20px;
  text-align: center;
  margin-bottom: 16px;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease forwards;
}

.success-icon {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: #34c759;
  color: #fff;
  font-size: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 12px;
}

.success-card h3 {
  margin: 0 0 8px;
}

.submitted-stars {
  font-size: 24px;
  margin-bottom: 12px;
}

.submitted-tags {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 6px;
  margin-bottom: 12px;
}

.submitted-tag {
  background: var(--button-color);
  color: var(--button-text-color);
  padding: 4px 12px;
  border-radius: 16px;
  font-size: 13px;
}

.submitted-comment {
  color: var(--hint-color);
  font-style: italic;
  font-size: 14px;
}

.reviews-list {
  margin-top: 16px;
}

.reviews-list h3 {
  font-size: 16px;
  margin: 0 0 12px;
}

.review-card {
  background: var(--secondary-bg);
  border-radius: 12px;
  padding: 14px;
  margin-bottom: 10px;
  opacity: 0;
  animation: fadeSlideUp 0.3s ease forwards;
}

.review-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.review-stars {
  font-size: 14px;
}

.star-sm {
  color: #d0d0d0;
}

.star-sm.filled {
  color: #ffb800;
}

.review-date {
  color: var(--hint-color);
  font-size: 12px;
}

.review-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-bottom: 6px;
}

.review-tag {
  background: rgba(0, 0, 0, 0.05);
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 12px;
  color: var(--hint-color);
}

.review-comment {
  font-size: 14px;
  margin: 0;
  color: var(--text-color);
}

@keyframes fadeSlideUp {
  from {
    opacity: 0;
    transform: translateY(12px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
