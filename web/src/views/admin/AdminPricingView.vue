<template>
  <div class="admin-page">
    <h1 class="page-title">Цены категорий</h1>

    <div v-if="loading" class="loading">Загрузка...</div>

    <div v-else class="admin-form">
      <fieldset class="admin-section">
        <legend class="admin-section-title">Цены за категории (₸)</legend>
        <p class="hint">Каждая цена — одноразовая оплата за пожизненный доступ к платному бакету категории.</p>

        <div class="price-row" v-for="cat in categories" :key="cat.key">
          <div class="price-info">
            <span class="price-icon">{{ cat.icon }}</span>
            <span class="price-name">{{ cat.label }}</span>
          </div>
          <div class="price-edit">
            <input
              v-model.number="prices[cat.key]"
              type="number"
              min="100"
              step="100"
              class="price-input"
              @change="markDirty(cat.key)"
            />
            <span class="price-currency">₸</span>
            <button
              v-if="dirty[cat.key]"
              class="price-save"
              :disabled="savingKey === cat.key"
              @click="save(cat.key)"
            >
              {{ savingKey === cat.key ? '...' : 'OK' }}
            </button>
            <span v-else-if="savedKey === cat.key" class="price-saved">✓ Сохранено</span>
          </div>
        </div>

        <div v-if="error" class="error-msg">{{ error }}</div>
      </fieldset>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { api } from '../../api'

type Key = 'workouts' | 'lfk' | 'nutrition'

const categories: { key: Key; label: string; icon: string }[] = [
  { key: 'workouts', label: 'Тренировки', icon: '🏋️' },
  { key: 'lfk', label: 'ЛФК', icon: '❤️‍🩹' },
  { key: 'nutrition', label: 'Питание', icon: '🍽️' },
]

const prices = reactive<Record<Key, number>>({ workouts: 0, lfk: 0, nutrition: 0 })
const dirty = reactive<Record<Key, boolean>>({ workouts: false, lfk: false, nutrition: false })
const loading = ref(true)
const savingKey = ref<Key | null>(null)
const savedKey = ref<Key | null>(null)
const error = ref('')

onMounted(async () => {
  try {
    const data = await api.getAdminPricing()
    for (const cat of categories) {
      prices[cat.key] = data[cat.key] ?? 0
    }
  } catch {
    error.value = 'Не удалось загрузить цены'
  } finally {
    loading.value = false
  }
})

function markDirty(key: Key) {
  dirty[key] = true
  savedKey.value = null
}

async function save(key: Key) {
  if (prices[key] <= 0) {
    error.value = 'Цена должна быть положительной'
    return
  }
  error.value = ''
  savingKey.value = key
  try {
    await api.setAdminPricing(key, prices[key])
    dirty[key] = false
    savedKey.value = key
    setTimeout(() => { if (savedKey.value === key) savedKey.value = null }, 2000)
  } catch (e: any) {
    error.value = e.message || 'Ошибка сохранения'
  } finally {
    savingKey.value = null
  }
}
</script>

<style scoped>
.admin-page { max-width: 400px; margin: 0 auto; padding-bottom: 24px; }
.page-title { font-size: 20px; font-weight: 700; margin-bottom: 16px; }
.loading { text-align: center; color: var(--hint-color); padding: 40px; }

.price-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
  gap: 12px;
}

.price-info {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
}

.price-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 10px;
  background: var(--bg-color);
  font-size: 18px;
}

.price-name {
  font-size: 16px;
  font-weight: 600;
}

.price-edit {
  display: flex;
  align-items: center;
  gap: 6px;
}

.price-input {
  width: 100px;
  padding: 8px 10px;
  border-radius: 10px;
  border: 1px solid var(--hint-color);
  background: var(--bg-color);
  color: var(--text-color);
  font-size: 15px;
  text-align: right;
}

.price-currency {
  color: var(--hint-color);
  font-size: 14px;
}

.price-save {
  padding: 6px 12px;
  border-radius: 10px;
  border: none;
  background: var(--button-color);
  color: var(--button-text-color);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
}

.price-save:disabled { opacity: 0.5; }

.price-saved {
  color: #34c759;
  font-size: 13px;
  font-weight: 600;
}

.hint { font-size: 12px; color: var(--hint-color); line-height: 1.4; }
.error-msg { color: #ff3b30; font-size: 14px; text-align: center; margin-top: 8px; }
</style>
