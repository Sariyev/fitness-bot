<template>
  <div class="session-page">
    <button class="back-btn" @click="router.back()">← Назад</button>

    <div v-if="loading" class="skeleton-list">
      <SkeletonCard v-for="i in 3" :key="i" />
    </div>

    <template v-else-if="workout">
      <h1 class="session-title">{{ workout.name }}</h1>
      <p v-if="workout.description" class="session-desc">{{ workout.description }}</p>

      <!-- Video player area -->
      <div class="video-area">
        <iframe
          v-if="workout.video_url"
          :src="workout.video_url"
          class="video-player"
          frameborder="0"
          allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
          allowfullscreen
        ></iframe>
        <div v-else class="video-placeholder">
          <span class="placeholder-icon">🎬</span>
          <span class="placeholder-text">Видео будет доступно позже</span>
        </div>
      </div>

      <!-- Workout meta -->
      <div class="workout-meta-row">
        <div class="meta-chip" v-if="workout.duration_minutes">
          {{ workout.duration_minutes }} мин
        </div>
        <div class="meta-chip" v-if="workout.format">
          {{ workout.format }}
        </div>
        <div class="meta-chip" v-if="workout.equipment && workout.equipment.length">
          {{ workout.equipment.join(', ') }}
        </div>
      </div>

      <p v-if="workout.expected_result" class="expected-result">
        {{ workout.expected_result }}
      </p>

      <!-- Exercises list -->
      <div class="exercises-section" v-if="workout.exercises && workout.exercises.length">
        <h2 class="section-title">Упражнения</h2>
        <div class="exercise-list">
          <div
            v-for="(ex, index) in workout.exercises"
            :key="ex.id"
            class="exercise-card"
            :style="{ animationDelay: index * 60 + 'ms' }"
          >
            <div class="exercise-number">{{ index + 1 }}</div>
            <div class="exercise-info">
              <span class="exercise-name">Упражнение #{{ ex.exercise_id }}</span>
              <span class="exercise-detail">
                <template v-if="ex.sets && ex.reps">{{ ex.sets }} x {{ ex.reps }}</template>
                <template v-else-if="ex.duration_seconds">{{ ex.duration_seconds }} сек</template>
              </span>
            </div>
          </div>
        </div>
      </div>

      <button
        class="btn btn-primary btn-complete"
        @click="complete"
        :disabled="completing"
      >
        {{ completing ? 'Сохранение...' : 'Завершить тренировку' }}
      </button>

      <div v-if="error" class="error-msg">{{ error }}</div>
    </template>

    <div v-else class="empty">
      <p>Тренировка не найдена</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api } from '../api'
import type { WorkoutWithExercises } from '../types'
import SkeletonCard from '../components/SkeletonCard.vue'

const props = defineProps<{ id: string }>()
const route = useRoute()
const router = useRouter()

const loading = ref(true)
const workout = ref<WorkoutWithExercises | null>(null)
const completing = ref(false)
const error = ref('')

async function complete() {
  if (completing.value || !workout.value) return
  completing.value = true
  error.value = ''
  try {
    await api.completeWorkout(workout.value.id)
    router.back()
  } catch (e: any) {
    error.value = e.message || 'Ошибка при завершении'
  } finally {
    completing.value = false
  }
}

onMounted(async () => {
  try {
    workout.value = await api.getWorkout(Number(props.id))
  } catch (e) {
    console.error('Failed to load workout:', e)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.session-page {
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

.skeleton-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.session-title {
  font-size: 22px;
  font-weight: 700;
  margin-bottom: 6px;
}

.session-desc {
  color: var(--hint-color);
  font-size: 14px;
  margin-bottom: 16px;
  line-height: 1.4;
}

.video-area {
  margin-bottom: 16px;
  border-radius: 12px;
  overflow: hidden;
  background: var(--secondary-bg);
}

.video-player {
  width: 100%;
  aspect-ratio: 16 / 9;
  display: block;
}

.video-placeholder {
  width: 100%;
  aspect-ratio: 16 / 9;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.placeholder-icon {
  font-size: 40px;
}

.placeholder-text {
  font-size: 14px;
  color: var(--hint-color);
}

.workout-meta-row {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 12px;
}

.meta-chip {
  background: var(--secondary-bg);
  padding: 4px 12px;
  border-radius: 16px;
  font-size: 13px;
  color: var(--hint-color);
}

.expected-result {
  font-size: 14px;
  color: var(--hint-color);
  font-style: italic;
  margin-bottom: 16px;
  padding: 12px;
  background: var(--secondary-bg);
  border-radius: 12px;
  line-height: 1.4;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 12px;
}

.exercise-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 24px;
}

.exercise-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  background: var(--secondary-bg);
  border-radius: 12px;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease forwards;
}

.exercise-number {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: var(--button-color);
  color: var(--button-text-color);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
  flex-shrink: 0;
}

.exercise-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.exercise-name {
  font-size: 15px;
  font-weight: 500;
}

.exercise-detail {
  font-size: 13px;
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
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-primary {
  background: var(--button-color);
  color: var(--button-text-color);
}

.btn-complete {
  margin-top: 8px;
  margin-bottom: 24px;
}

.error-msg {
  text-align: center;
  color: #ff3b30;
  font-size: 14px;
  margin-top: 8px;
}

.empty {
  text-align: center;
  padding: 40px 0;
  color: var(--hint-color);
}
</style>
