<template>
  <div class="admin-page">
    <button class="back-btn" @click="router.push('/admin')">← Назад</button>
    <h1 class="page-title">Пользователи <span class="total-badge" v-if="total">({{ total }})</span></h1>

    <div v-if="loading" class="loading">Загрузка...</div>

    <div v-else>
      <div
        v-for="user in users"
        :key="user.id"
        class="user-card"
        @click="router.push(`/admin/users/${user.id}`)"
      >
        <div class="user-info">
          <span class="user-name">{{ user.first_name }} {{ user.last_name }}</span>
          <span class="user-username" v-if="user.username">@{{ user.username }}</span>
        </div>
        <div class="user-badges">
          <span class="badge badge-role" :class="user.role">{{ user.role }}</span>
          <span class="badge" :class="user.is_paid ? 'badge-paid' : 'badge-unpaid'">
            {{ user.is_paid ? 'Paid' : 'Free' }}
          </span>
        </div>
      </div>

      <div class="pagination" v-if="total > pageSize">
        <button class="page-btn" :disabled="offset === 0" @click="prevPage">← Назад</button>
        <span class="page-info">{{ Math.floor(offset / pageSize) + 1 }} / {{ Math.ceil(total / pageSize) }}</span>
        <button class="page-btn" :disabled="offset + pageSize >= total" @click="nextPage">Вперёд →</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../api'
import type { AdminUser } from '../../types'

const router = useRouter()
const users = ref<AdminUser[]>([])
const total = ref(0)
const loading = ref(true)
const pageSize = 20
const offset = ref(0)

async function loadUsers() {
  loading.value = true
  try {
    const res = await api.getAdminUsers(pageSize, offset.value)
    users.value = res.users || []
    total.value = res.total
  } catch {
    users.value = []
  } finally {
    loading.value = false
  }
}

function prevPage() {
  offset.value = Math.max(0, offset.value - pageSize)
  loadUsers()
}

function nextPage() {
  offset.value += pageSize
  loadUsers()
}

onMounted(loadUsers)
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

.total-badge {
  color: var(--hint-color);
  font-weight: 400;
  font-size: 16px;
}

.loading {
  text-align: center;
  color: var(--hint-color);
  padding: 40px;
}

.user-card {
  background: var(--secondary-bg);
  border-radius: 12px;
  padding: 14px 16px;
  margin-bottom: 8px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  cursor: pointer;
  animation: fadeSlideUp 0.3s ease both;
}

.user-info {
  display: flex;
  flex-direction: column;
}

.user-name {
  font-weight: 600;
  font-size: 15px;
}

.user-username {
  color: var(--hint-color);
  font-size: 13px;
}

.user-badges {
  display: flex;
  gap: 6px;
}

.badge {
  padding: 3px 10px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.badge-role.admin {
  background: #ff9500;
  color: #fff;
}

.badge-role.client {
  background: rgba(0,0,0,0.08);
  color: var(--hint-color);
}

.badge-paid {
  background: #34c759;
  color: #fff;
}

.badge-unpaid {
  background: rgba(0,0,0,0.08);
  color: var(--hint-color);
}

.pagination {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 16px;
  padding: 8px 0;
}

.page-btn {
  background: var(--secondary-bg);
  border: none;
  padding: 8px 16px;
  border-radius: 8px;
  color: var(--button-color);
  cursor: pointer;
  font-size: 14px;
}

.page-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.page-info {
  color: var(--hint-color);
  font-size: 14px;
}

@keyframes fadeSlideUp {
  from { opacity: 0; transform: translateY(12px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
