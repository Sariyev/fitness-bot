<template>
  <div class="admin-page">
    <button class="back-btn" @click="router.push('/profile')">← Назад</button>
    <h1 class="page-title">Админ панель</h1>

    <div class="stats-grid">
      <div class="stat-card">
        <span class="stat-value">{{ stats?.total_users ?? '-' }}</span>
        <span class="stat-label">Пользователей</span>
      </div>
      <div class="stat-card">
        <span class="stat-value">{{ reviewSummary?.average_score?.toFixed(1) ?? '-' }}</span>
        <span class="stat-label">Рейтинг ({{ reviewSummary?.total_reviews ?? 0 }})</span>
      </div>
    </div>

    <div class="nav-cards">
      <router-link to="/admin/users" class="nav-card">
        <span class="nav-icon">👥</span>
        <span>Пользователи</span>
      </router-link>
      <router-link to="/admin/content" class="nav-card">
        <span class="nav-icon">🏋️</span>
        <span>Контент</span>
      </router-link>
      <router-link to="/admin/reviews" class="nav-card">
        <span class="nav-icon">⭐</span>
        <span>Отзывы</span>
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../api'
import type { ReviewSummary } from '../../types'

const router = useRouter()
const stats = ref<{ total_users: number } | null>(null)
const reviewSummary = ref<ReviewSummary | null>(null)

onMounted(async () => {
  const [s, r] = await Promise.allSettled([
    api.getAdminStats(),
    api.getAdminReviewsSummary(),
  ])
  if (s.status === 'fulfilled') stats.value = s.value
  if (r.status === 'fulfilled') reviewSummary.value = r.value
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
  font-size: 22px;
  font-weight: 700;
  margin-bottom: 16px;
}

.stats-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  margin-bottom: 20px;
}

.stat-card {
  background: var(--secondary-bg);
  border-radius: 12px;
  padding: 20px 16px;
  text-align: center;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease forwards;
}

.stat-value {
  display: block;
  font-size: 28px;
  font-weight: 700;
  color: var(--button-color);
}

.stat-label {
  display: block;
  font-size: 13px;
  color: var(--hint-color);
  margin-top: 4px;
}

.nav-cards {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.nav-card {
  display: flex;
  align-items: center;
  gap: 12px;
  background: var(--secondary-bg);
  border-radius: 12px;
  padding: 16px;
  text-decoration: none;
  color: var(--text-color);
  font-size: 16px;
  font-weight: 500;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease forwards;
}

.nav-card:nth-child(1) { animation-delay: 80ms; }
.nav-card:nth-child(2) { animation-delay: 160ms; }
.nav-card:nth-child(3) { animation-delay: 240ms; }

.nav-icon {
  font-size: 24px;
}

@keyframes fadeSlideUp {
  from { opacity: 0; transform: translateY(12px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
