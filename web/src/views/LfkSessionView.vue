<template>
  <div class="lfk-session-page">
    <!-- Back button -->
    <button class="back-btn" @click="router.back()">← Назад</button>

    <!-- Loading state -->
    <div v-if="loading" class="skeleton-list">
      <SkeletonCard v-for="i in 4" :key="i" />
    </div>

    <!-- Session content -->
    <template v-else-if="session">
      <!-- Success state after completion -->
      <div v-if="completed" class="success-state">
        <div class="success-animation">
          <span class="success-check">✅</span>
          <div class="confetti">
            <span v-for="i in 12" :key="i" class="confetti-piece" :style="confettiStyle(i)"></span>
          </div>
        </div>
        <h2 class="success-title">Отлично!</h2>
        <p class="success-text">Занятие завершено. Так держать!</p>
        <button class="btn btn-primary" @click="router.back()">
          Вернуться к курсу
        </button>
      </div>

      <!-- Normal session view -->
      <template v-else>
        <!-- Header: Day + Stage -->
        <div class="session-header">
          <h1 class="session-title">День {{ session.day_number }}</h1>
          <span class="stage-badge">Этап {{ session.stage }}/3</span>
        </div>

        <!-- Stage progress bar -->
        <div class="stage-progress">
          <div
            v-for="s in 3"
            :key="s"
            class="stage-segment"
            :class="{ filled: s <= session.stage }"
          ></div>
        </div>

        <!-- Video placeholder -->
        <div class="video-section">
          <a
            v-if="session.video_url"
            :href="session.video_url"
            target="_blank"
            rel="noopener noreferrer"
            class="video-placeholder"
          >
            <span class="video-play-icon">▶️</span>
            <span class="video-label">Смотреть видео</span>
          </a>
          <div v-else class="video-placeholder video-unavailable">
            <span class="video-play-icon">🎬</span>
            <span class="video-label">Видео будет доступно позже</span>
          </div>
        </div>

        <!-- Duration badge -->
        <div class="meta-row" v-if="session.duration_minutes">
          <span class="duration-badge">⏱ {{ session.duration_minutes }} минут</span>
        </div>

        <!-- Description -->
        <div class="description-section" v-if="session.description">
          <h2 class="section-title">Описание</h2>
          <p class="description-text">{{ session.description }}</p>
        </div>

        <!-- Completion form -->
        <div class="completion-section">
          <h2 class="section-title">Завершить занятие</h2>

          <!-- Pain level slider -->
          <div class="pain-block">
            <div class="pain-header">
              <label class="field-label">Уровень боли</label>
              <span class="pain-value" :style="{ color: painColor }">{{ painLevel }}</span>
            </div>
            <div class="pain-slider-container">
              <input
                type="range"
                min="0"
                max="10"
                step="1"
                v-model.number="painLevel"
                class="pain-slider"
              />
              <div class="pain-track" :style="painTrackStyle"></div>
            </div>
            <div class="pain-labels">
              <span>0 — нет боли</span>
              <span>10 — сильная</span>
            </div>
            <div class="pain-emoji">{{ painEmoji }}</div>
          </div>

          <!-- Comment -->
          <div class="comment-block">
            <label class="field-label">Комментарий (необязательно)</label>
            <textarea
              v-model="comment"
              class="comment-input"
              rows="3"
              placeholder="Как прошло занятие? Что ощущаете?"
            ></textarea>
          </div>

          <!-- Submit -->
          <button
            class="btn btn-primary btn-complete"
            @click="completeSession"
            :disabled="completing"
          >
            {{ completing ? 'Сохранение...' : 'Завершить занятие' }}
          </button>

          <div v-if="error" class="error-msg">{{ error }}</div>
        </div>
      </template>
    </template>

    <!-- Not found state -->
    <div v-else class="empty">
      <p>Сессия не найдена</p>
      <button class="btn btn-secondary" @click="router.back()">Назад</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api } from '../api'
import type { RehabSession } from '../types'
import SkeletonCard from '../components/SkeletonCard.vue'

const props = defineProps<{ id: string }>()
const route = useRoute()
const router = useRouter()

// ---------- State ----------

const loading = ref(true)
const session = ref<RehabSession | null>(null)
const painLevel = ref(0)
const comment = ref('')
const completing = ref(false)
const completed = ref(false)
const error = ref('')

// ---------- Computed ----------

const painColor = computed(() => {
  const level = painLevel.value
  if (level <= 3) return '#34c759'
  if (level <= 6) return '#ffcc00'
  return '#ff3b30'
})

const painTrackStyle = computed(() => {
  const pct = (painLevel.value / 10) * 100
  return {
    width: `${pct}%`,
    background: `linear-gradient(to right, #34c759, ${painLevel.value > 5 ? '#ff3b30' : '#ffcc00'})`,
  }
})

const painEmoji = computed(() => {
  const level = painLevel.value
  if (level === 0) return '😊'
  if (level <= 2) return '🙂'
  if (level <= 4) return '😐'
  if (level <= 6) return '😕'
  if (level <= 8) return '😣'
  return '😫'
})

// ---------- Actions ----------

async function completeSession() {
  if (completing.value || !session.value) return
  completing.value = true
  error.value = ''

  try {
    await api.completeRehabSession(session.value.id, {
      pain_level: painLevel.value,
      comment: comment.value,
      day_number: session.value.day_number,
      course_id: session.value.course_id,
    })
    completed.value = true
  } catch (e: any) {
    error.value = e.message || 'Ошибка при сохранении'
  } finally {
    completing.value = false
  }
}

function confettiStyle(i: number): Record<string, string> {
  const colors = ['#ff3b30', '#ff9500', '#ffcc00', '#34c759', '#007aff', '#af52de']
  const angle = (i / 12) * 360
  const distance = 40 + Math.random() * 50
  const x = Math.cos((angle * Math.PI) / 180) * distance
  const y = Math.sin((angle * Math.PI) / 180) * distance
  return {
    '--x': `${x}px`,
    '--y': `${y}px`,
    '--color': colors[i % colors.length],
    '--delay': `${i * 50}ms`,
    '--rotation': `${Math.random() * 360}deg`,
  }
}

// ---------- Lifecycle ----------

onMounted(async () => {
  try {
    session.value = await api.getRehabSession(Number(props.id))
  } catch (e) {
    console.error('Failed to load session:', e)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.lfk-session-page {
  max-width: 400px;
  margin: 0 auto;
  padding-bottom: 40px;
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

.skeleton-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* ===== Session header ===== */
.session-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.session-title {
  font-size: 24px;
  font-weight: 700;
}

.stage-badge {
  background: var(--button-color);
  color: var(--button-text-color);
  padding: 4px 12px;
  border-radius: 16px;
  font-size: 13px;
  font-weight: 600;
  white-space: nowrap;
}

/* ===== Stage progress bar ===== */
.stage-progress {
  display: flex;
  gap: 4px;
  margin-bottom: 20px;
}

.stage-segment {
  flex: 1;
  height: 4px;
  border-radius: 2px;
  background: var(--secondary-bg);
  transition: background 0.3s ease;
}

.stage-segment.filled {
  background: var(--button-color);
}

/* ===== Video section ===== */
.video-section {
  margin-bottom: 16px;
}

.video-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 10px;
  width: 100%;
  aspect-ratio: 16 / 9;
  background: linear-gradient(135deg, var(--secondary-bg), color-mix(in srgb, var(--button-color) 8%, var(--secondary-bg)));
  border-radius: 14px;
  text-decoration: none;
  transition: transform 0.15s ease, box-shadow 0.15s ease;
  cursor: pointer;
}

.video-placeholder:active {
  transform: scale(0.98);
}

.video-unavailable {
  cursor: default;
  opacity: 0.7;
}

.video-unavailable:active {
  transform: none;
}

.video-play-icon {
  font-size: 48px;
  line-height: 1;
}

.video-label {
  font-size: 14px;
  color: var(--button-color);
  font-weight: 500;
}

.video-unavailable .video-label {
  color: var(--hint-color);
}

/* ===== Duration ===== */
.meta-row {
  margin-bottom: 20px;
}

.duration-badge {
  display: inline-block;
  background: var(--secondary-bg);
  padding: 6px 14px;
  border-radius: 16px;
  font-size: 14px;
  color: var(--hint-color);
  font-weight: 500;
}

/* ===== Description ===== */
.description-section {
  margin-bottom: 24px;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 10px;
}

.description-text {
  font-size: 15px;
  line-height: 1.5;
  color: var(--text-color);
  white-space: pre-line;
}

/* ===== Completion section ===== */
.completion-section {
  background: var(--secondary-bg);
  border-radius: 14px;
  padding: 20px 16px;
  margin-top: 8px;
}

/* Pain slider */
.pain-block {
  margin-bottom: 20px;
}

.pain-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
}

.field-label {
  font-size: 14px;
  color: var(--text-color);
  font-weight: 500;
}

.pain-value {
  font-size: 22px;
  font-weight: 700;
  min-width: 32px;
  text-align: center;
  transition: color 0.2s ease;
}

.pain-slider-container {
  position: relative;
  height: 28px;
  display: flex;
  align-items: center;
}

.pain-slider {
  -webkit-appearance: none;
  appearance: none;
  width: 100%;
  height: 8px;
  border-radius: 4px;
  background: linear-gradient(to right, #34c759 0%, #ffcc00 50%, #ff3b30 100%);
  outline: none;
  position: relative;
  z-index: 2;
}

.pain-slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  appearance: none;
  width: 26px;
  height: 26px;
  border-radius: 50%;
  background: var(--bg-color);
  border: 3px solid var(--button-color);
  cursor: pointer;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.15);
  transition: transform 0.1s ease;
}

.pain-slider::-webkit-slider-thumb:active {
  transform: scale(1.15);
}

.pain-slider::-moz-range-thumb {
  width: 26px;
  height: 26px;
  border-radius: 50%;
  background: var(--bg-color);
  border: 3px solid var(--button-color);
  cursor: pointer;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.15);
}

.pain-labels {
  display: flex;
  justify-content: space-between;
  font-size: 11px;
  color: var(--hint-color);
  margin-top: 6px;
}

.pain-emoji {
  text-align: center;
  font-size: 32px;
  margin-top: 8px;
  transition: all 0.2s ease;
}

/* Comment */
.comment-block {
  margin-bottom: 20px;
}

.comment-input {
  width: 100%;
  padding: 12px;
  border: 1px solid var(--hint-color);
  border-radius: 10px;
  font-size: 14px;
  background: var(--bg-color);
  color: var(--text-color);
  resize: vertical;
  font-family: inherit;
  box-sizing: border-box;
  margin-top: 6px;
  transition: border-color 0.15s ease;
}

.comment-input:focus {
  outline: none;
  border-color: var(--button-color);
}

.comment-input::placeholder {
  color: var(--hint-color);
}

/* ===== Success state ===== */
.success-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 60px 20px;
  min-height: 60vh;
}

.success-animation {
  position: relative;
  margin-bottom: 24px;
}

.success-check {
  font-size: 72px;
  display: block;
  animation: successPop 0.5s cubic-bezier(0.175, 0.885, 0.32, 1.275) forwards;
}

@keyframes successPop {
  0% {
    transform: scale(0);
    opacity: 0;
  }
  60% {
    transform: scale(1.2);
    opacity: 1;
  }
  100% {
    transform: scale(1);
    opacity: 1;
  }
}

.confetti {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 0;
  height: 0;
}

.confetti-piece {
  position: absolute;
  width: 8px;
  height: 8px;
  border-radius: 2px;
  background: var(--color);
  opacity: 0;
  animation: confettiBurst 0.7s var(--delay) ease-out forwards;
}

@keyframes confettiBurst {
  0% {
    transform: translate(0, 0) rotate(0deg) scale(1);
    opacity: 1;
  }
  100% {
    transform: translate(var(--x), var(--y)) rotate(var(--rotation)) scale(0);
    opacity: 0;
  }
}

.success-title {
  font-size: 26px;
  font-weight: 700;
  margin-bottom: 8px;
  animation: fadeSlideUp 0.4s 0.3s ease forwards;
  opacity: 0;
}

.success-text {
  font-size: 16px;
  color: var(--hint-color);
  margin-bottom: 32px;
  animation: fadeSlideUp 0.4s 0.45s ease forwards;
  opacity: 0;
}

.success-state .btn {
  animation: fadeSlideUp 0.4s 0.6s ease forwards;
  opacity: 0;
  max-width: 280px;
}

/* ===== Buttons ===== */
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
  transition: opacity 0.15s ease;
}

.btn:active {
  opacity: 0.85;
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
  background: var(--secondary-bg);
  color: var(--text-color);
}

.btn-complete {
  font-weight: 600;
}

.error-msg {
  text-align: center;
  color: #ff3b30;
  font-size: 14px;
  margin-top: 8px;
}

.empty {
  text-align: center;
  padding: 60px 20px;
  color: var(--hint-color);
}

.empty p {
  margin-bottom: 16px;
  font-size: 16px;
}

.empty .btn {
  max-width: 200px;
  margin: 0 auto;
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
</style>
