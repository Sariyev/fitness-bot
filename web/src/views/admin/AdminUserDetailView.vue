<template>
  <div class="admin-page">
    <button class="back-btn" @click="router.push('/admin/users')">← Назад</button>

    <div v-if="loading" class="loading">Загрузка...</div>

    <div v-else-if="user">
      <div class="user-header">
        <div class="avatar">{{ initials }}</div>
        <h2>{{ user.first_name }} {{ user.last_name }}</h2>
        <p class="username" v-if="user.username">@{{ user.username }}</p>
        <p class="meta">ID: {{ user.id }} | TG: {{ user.telegram_id }}</p>
      </div>

      <div class="info-card">
        <div class="info-row">
          <span class="label">Роль</span>
          <div class="toggle-group">
            <button
              class="toggle-btn"
              :class="{ active: user.role === 'client' }"
              @click="setRole('client')"
              :disabled="saving"
            >Client</button>
            <button
              class="toggle-btn"
              :class="{ active: user.role === 'admin' }"
              @click="setRole('admin')"
              :disabled="saving"
            >Admin</button>
          </div>
        </div>
        <div class="info-row">
          <span class="label">Оплата</span>
          <div class="toggle-group">
            <button
              class="toggle-btn"
              :class="{ active: !user.is_paid }"
              @click="setPaid(false)"
              :disabled="saving"
            >Free</button>
            <button
              class="toggle-btn paid-btn"
              :class="{ active: user.is_paid }"
              @click="setPaid(true)"
              :disabled="saving"
            >Paid</button>
          </div>
        </div>
        <div class="info-row">
          <span class="label">Регистрация</span>
          <span class="value">{{ user.is_registered ? 'Да' : 'Нет' }}</span>
        </div>
        <div class="info-row">
          <span class="label">Создан</span>
          <span class="value">{{ formatDate(user.created_at) }}</span>
        </div>
      </div>

      <div class="info-card" v-if="profile">
        <h3>Профиль</h3>
        <div class="info-row">
          <span class="label">Возраст</span>
          <span class="value">{{ profile.age }}</span>
        </div>
        <div class="info-row">
          <span class="label">Пол</span>
          <span class="value">{{ profile.gender === 'male' ? 'М' : 'Ж' }}</span>
        </div>
        <div class="info-row">
          <span class="label">Рост / Вес</span>
          <span class="value">{{ profile.height_cm }} см / {{ profile.weight_kg }} кг</span>
        </div>
        <div class="info-row">
          <span class="label">Уровень</span>
          <span class="value">{{ profile.fitness_level }}</span>
        </div>
        <div class="info-row" v-if="profile.goals && profile.goals.length">
          <span class="label">Цели</span>
          <span class="value">{{ profile.goals.join(', ') }}</span>
        </div>
      </div>

      <p v-if="error" class="error-msg">{{ error }}</p>
    </div>

    <div v-else class="loading">Пользователь не найден</div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../api'
import type { AdminUser, UserProfile } from '../../types'

const props = defineProps<{ id: string | number }>()
const router = useRouter()

const user = ref<AdminUser | null>(null)
const profile = ref<UserProfile | null>(null)
const loading = ref(true)
const saving = ref(false)
const error = ref('')

const initials = computed(() => {
  if (!user.value) return '?'
  const f = user.value.first_name?.[0] || ''
  const l = user.value.last_name?.[0] || ''
  return (f + l).toUpperCase() || '?'
})

function formatDate(d: string): string {
  return new Date(d).toLocaleDateString('ru-RU', { day: 'numeric', month: 'short', year: 'numeric' })
}

async function setRole(role: string) {
  if (!user.value || user.value.role === role || saving.value) return
  saving.value = true
  error.value = ''
  try {
    await api.updateAdminUser(user.value.id, { role })
    user.value.role = role
  } catch (e: any) {
    error.value = e.message || 'Ошибка'
  } finally {
    saving.value = false
  }
}

async function setPaid(isPaid: boolean) {
  if (!user.value || user.value.is_paid === isPaid || saving.value) return
  saving.value = true
  error.value = ''
  try {
    await api.updateAdminUser(user.value.id, { is_paid: isPaid })
    user.value.is_paid = isPaid
  } catch (e: any) {
    error.value = e.message || 'Ошибка'
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  try {
    const res = await api.getAdminUser(Number(props.id))
    user.value = res.user
    profile.value = res.profile
  } catch {
    user.value = null
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

.loading {
  text-align: center;
  color: var(--hint-color);
  padding: 40px;
}

.user-header {
  text-align: center;
  background: var(--secondary-bg);
  border-radius: 12px;
  padding: 24px 16px;
  margin-bottom: 12px;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease forwards;
}

.avatar {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: var(--button-color);
  color: var(--button-text-color);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22px;
  font-weight: bold;
  margin: 0 auto 10px;
}

.username {
  color: var(--hint-color);
  font-size: 14px;
}

.meta {
  color: var(--hint-color);
  font-size: 12px;
  margin-top: 4px;
}

.info-card {
  background: var(--secondary-bg);
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 12px;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease forwards;
  animation-delay: 80ms;
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
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid rgba(0,0,0,0.05);
}

.info-row:last-child {
  border-bottom: none;
}

.label {
  color: var(--hint-color);
  font-size: 14px;
}

.value {
  font-size: 14px;
  text-align: right;
  max-width: 55%;
}

.toggle-group {
  display: flex;
  gap: 4px;
  background: var(--bg-color);
  border-radius: 8px;
  padding: 2px;
}

.toggle-btn {
  padding: 5px 14px;
  border: none;
  border-radius: 6px;
  font-size: 13px;
  cursor: pointer;
  background: transparent;
  color: var(--hint-color);
  transition: all 0.15s;
}

.toggle-btn.active {
  background: var(--button-color);
  color: var(--button-text-color);
}

.toggle-btn.paid-btn.active {
  background: #34c759;
}

.toggle-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.error-msg {
  text-align: center;
  color: #ff3b30;
  font-size: 14px;
  margin-top: 8px;
}

@keyframes fadeSlideUp {
  from { opacity: 0; transform: translateY(12px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
