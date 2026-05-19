<template>
  <div class="categories-view">
    <h1 class="page-title">Категории 📂</h1>
    <div v-if="loading" class="skeleton-list">
      <SkeletonCard v-for="i in 4" :key="i" />
    </div>
    <div v-else-if="categories.length === 0" class="empty">
      <p>Категории пока не доступны</p>
    </div>
    <div v-else class="category-list">
      <div
        v-for="(cat, index) in categories"
        :key="cat.id"
        class="category-card"
        :style="{ animationDelay: index * 60 + 'ms' }"
        @click="openCategory(cat.id)"
      >
        <div class="category-header">
          <span class="category-icon">{{ cat.icon || '📂' }}</span>
          <div class="category-info">
            <h3>{{ cat.name }}</h3>
            <p v-if="cat.description">{{ cat.description }}</p>
          </div>
          <span class="category-arrow">›</span>
        </div>
        <ProgressBar
          v-if="cat.total > 0"
          :completed="cat.completed"
          :total="cat.total"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../api'
import { useTelegram } from '../composables/useTelegram'
import type { CategoryWithProgress } from '../types'
import ProgressBar from '../components/ProgressBar.vue'
import SkeletonCard from '../components/SkeletonCard.vue'

const props = defineProps<{ id: string }>()
const router = useRouter()
const { hapticSelection } = useTelegram()
const categories = ref<CategoryWithProgress[]>([])
const loading = ref(true)

function openCategory(id: number) {
  hapticSelection()
  router.push({ name: 'lessons', params: { id } })
}

// BackButton is wired centrally in router.ts (top-left back arrow).

onMounted(async () => {
  try {
    categories.value = await api.getCategories(Number(props.id))
  } catch (e) {
    console.error('Failed to load categories:', e)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.page-title {
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 16px;
}

.skeleton-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.category-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.category-card {
  padding: 16px;
  background: var(--secondary-bg);
  border-radius: 12px;
  cursor: pointer;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease forwards;
  transition: transform 0.15s ease;
}

.category-card:active {
  transform: scale(0.98);
}

.category-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.category-icon {
  font-size: 24px;
  flex-shrink: 0;
}

.category-info {
  flex: 1;
  min-width: 0;
}

.category-info h3 {
  font-size: 16px;
  font-weight: 600;
}

.category-info p {
  font-size: 13px;
  color: var(--hint-color);
  margin-top: 2px;
}

.category-arrow {
  font-size: 20px;
  color: var(--hint-color);
  flex-shrink: 0;
}

.empty {
  text-align: center;
  padding: 40px 0;
  color: var(--hint-color);
}
</style>
