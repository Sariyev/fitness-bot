<template>
  <div class="profile-page">
    <div v-if="loading" class="loading">
      <div class="spinner"></div>
    </div>

    <div v-else-if="profile" class="profile-content">
      <div class="profile-header">
        <div class="avatar">{{ initials }}</div>
        <h2>{{ profile.first_name }} {{ profile.last_name }}</h2>
        <p v-if="profile.username" class="username">@{{ profile.username }}</p>
      </div>

      <div class="info-card">
        <div class="info-row">
          <span class="label">Возраст</span>
          <span class="value">{{ profile.age }}</span>
        </div>
        <div class="info-row">
          <span class="label">Рост</span>
          <span class="value">{{ profile.height_cm }} см</span>
        </div>
        <div class="info-row">
          <span class="label">Вес</span>
          <span class="value">{{ profile.weight_kg }} кг</span>
        </div>
        <div class="info-row">
          <span class="label">Пол</span>
          <span class="value">{{ genderLabel }}</span>
        </div>
        <div class="info-row">
          <span class="label">Уровень</span>
          <span class="value">{{ fitnessLabel }}</span>
        </div>
      </div>

      <div class="info-card" v-if="profile.goals && profile.goals.length">
        <h3>Цели</h3>
        <div class="goals-list">
          <span v-for="goal in profile.goals" :key="goal" class="goal-tag">{{ goalLabel(goal) }}</span>
        </div>
      </div>

      <div class="info-card">
        <div class="info-row">
          <span class="label">Доступ</span>
          <span class="value" :class="profile.is_paid ? 'paid' : 'unpaid'">
            {{ profile.is_paid ? '✅ Оплачено' : 'Не оплачено' }}
          </span>
        </div>
      </div>

      <button v-if="!profile.is_paid" class="btn btn-primary" @click="$router.push('/payment')">
        Оплатить доступ
      </button>
      <button class="btn btn-secondary" @click="$router.push('/')">
        К модулям
      </button>
    </div>

    <div v-else class="error">
      <p>Не удалось загрузить профиль</p>
      <button class="btn" @click="$router.push('/')">Назад</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { api } from '../api'
import type { UserProfile } from '../types'

const loading = ref(true)
const profile = ref<UserProfile | null>(null)

const goalLabels: Record<string, string> = {
  weight_loss: 'Похудеть',
  muscle_gain: 'Набрать массу',
  strength: 'Больше силы',
  endurance: 'Выносливость',
  maintenance: 'Поддержание формы',
  hernia: 'Грыжа',
  protrusion: 'Протрузии',
  scoliosis: 'Сколиоз',
  kyphosis: 'Кифоз',
  lordosis: 'Лордоз',
}

function goalLabel(key: string): string {
  return goalLabels[key] || key
}

const initials = computed(() => {
  if (!profile.value) return '?'
  const f = profile.value.first_name?.[0] || ''
  const l = profile.value.last_name?.[0] || ''
  return (f + l).toUpperCase() || '?'
})

const genderLabel = computed(() => {
  return profile.value?.gender === 'male' ? 'Мужской' : 'Женский'
})

const fitnessLabel = computed(() => {
  const labels: Record<string, string> = {
    beginner: 'Новичок',
    intermediate: 'Средний',
    advanced: 'Продвинутый',
  }
  return labels[profile.value?.fitness_level || ''] || profile.value?.fitness_level || ''
})

onMounted(async () => {
  try {
    profile.value = await api.getProfile()
  } catch {
    profile.value = null
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.profile-page {
  max-width: 400px;
  margin: 0 auto;
}

.profile-header {
  text-align: center;
  margin-bottom: 20px;
}

.avatar {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: var(--button-color);
  color: var(--button-text-color);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: bold;
  margin: 0 auto 12px;
}

.username {
  color: var(--hint-color);
  font-size: 14px;
}

.info-card {
  background: var(--secondary-bg);
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 12px;
}

.info-card h3 {
  margin-bottom: 8px;
  font-size: 14px;
  color: var(--hint-color);
  text-transform: uppercase;
}

.info-row {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  border-bottom: 1px solid rgba(0,0,0,0.05);
}

.info-row:last-child {
  border-bottom: none;
}

.label {
  color: var(--hint-color);
}

.value.paid {
  color: #34c759;
}

.value.unpaid {
  color: #ff3b30;
}

.goals-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.goal-tag {
  background: var(--button-color);
  color: var(--button-text-color);
  padding: 4px 12px;
  border-radius: 16px;
  font-size: 13px;
}

.btn {
  display: block;
  width: 100%;
  padding: 14px;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  cursor: pointer;
  margin-bottom: 8px;
  text-align: center;
}

.btn-primary {
  background: var(--button-color);
  color: var(--button-text-color);
}

.btn-secondary {
  background: var(--secondary-bg);
  color: var(--text-color);
}

.error {
  text-align: center;
  padding: 40px;
  color: var(--hint-color);
}
</style>
