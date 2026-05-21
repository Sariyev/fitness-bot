<template>
  <div class="locked-card">
    <div class="lock-icon">🔒</div>
    <h2 class="locked-title">{{ titleByCategory }}</h2>
    <p class="locked-text">
      Эта {{ noun }} в платном бакете категории <strong>{{ categoryLabel }}</strong>.
      Одна оплата открывает все платные {{ noun }}-материалы навсегда.
    </p>
    <p v-if="tier === 'trial'" class="locked-trial-hint">
      Триал-доступ доступен первые 7 дней после регистрации.
    </p>
    <button
      class="unlock-btn"
      :disabled="paying"
      @click="unlock"
    >
      {{ paying ? 'Открываем оплату…' : `Открыть за ${priceKzt} ₸` }}
    </button>
    <p v-if="error" class="error-msg">{{ error }}</p>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { api } from '../api'
import type { ContentCategory } from '../types'

const props = defineProps<{
  category: ContentCategory
  tier: 'free' | 'trial' | 'paid'
  priceKzt: number
}>()

const paying = ref(false)
const error = ref('')

const categoryLabel = computed(() => {
  const map: Record<ContentCategory, string> = {
    workouts: 'Тренировки',
    lfk: 'ЛФК',
    nutrition: 'Питание',
  }
  return map[props.category] || props.category
})

const noun = computed(() => {
  const map: Record<ContentCategory, string> = {
    workouts: 'программа',
    lfk: 'курс',
    nutrition: 'план',
  }
  return map[props.category] || 'материал'
})

const titleByCategory = computed(() => {
  const map: Record<ContentCategory, string> = {
    workouts: 'Программа закрыта',
    lfk: 'Курс закрыт',
    nutrition: 'План закрыт',
  }
  return map[props.category] || 'Контент закрыт'
})

async function unlock() {
  if (paying.value) return
  paying.value = true
  error.value = ''
  try {
    const res = await api.processPayment(props.category)
    if (res.redirect_url) {
      // Robokassa flow — full-page redirect into the provider's checkout.
      window.location.href = res.redirect_url
      return
    }
    // Sync provider (dummy in dev) — payment already confirmed server-side.
    // Reload to refetch the content with unlocked status.
    window.location.reload()
  } catch (e: any) {
    error.value = e.message || 'Не удалось начать оплату'
    paying.value = false
  }
}
</script>

<style scoped>
.locked-card {
  background: var(--secondary-bg);
  border-radius: 16px;
  padding: 32px 20px;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  margin: 24px 0;
  animation: fadeSlideUp 0.4s ease both;
}

.lock-icon {
  font-size: 48px;
  line-height: 1;
  opacity: 0.8;
}

.locked-title {
  font-size: 20px;
  font-weight: 700;
  margin: 0;
}

.locked-text {
  font-size: 14px;
  color: var(--hint-color);
  line-height: 1.5;
  margin: 0;
  max-width: 320px;
}

.locked-trial-hint {
  font-size: 13px;
  color: #b35700;
  background: rgba(255, 149, 0, 0.14);
  padding: 8px 12px;
  border-radius: 10px;
  margin: 0;
}

.unlock-btn {
  padding: 14px 28px;
  border-radius: 14px;
  border: none;
  background: var(--button-color);
  color: var(--button-text-color);
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.12);
  transition: transform 0.1s, opacity 0.15s;
  -webkit-tap-highlight-color: transparent;
  margin-top: 8px;
}

.unlock-btn:active { transform: scale(0.97); }
.unlock-btn:disabled { opacity: 0.5; cursor: wait; }

.error-msg {
  color: #ff3b30;
  font-size: 13px;
  margin: 4px 0 0;
}

@keyframes fadeSlideUp {
  from { opacity: 0; transform: translateY(12px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
