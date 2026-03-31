<template>
  <div class="payment-page">
    <div v-if="loading" class="skeleton-list">
      <SkeletonCard />
      <SkeletonCard />
    </div>

    <div v-else-if="paid" class="success-card">
      <div class="success-icon">✅</div>
      <h2>Доступ оплачен</h2>
      <p>Все модули доступны. Приятных тренировок! 💪</p>
      <button class="btn btn-primary" @click="goToModules">Перейти к модулям</button>
    </div>

    <div v-else class="payment-flow">
      <Transition name="fade" mode="out-in">
        <!-- Step 1: Product card -->
        <div v-if="step === 'product'" key="product" class="product-card">
          <div class="product-gradient"></div>
          <h2>🏋️ Полный доступ</h2>
          <div class="features">
            <div class="feature">🏥 ЛФК — упражнения при проблемах со здоровьем</div>
            <div class="feature">💪 Тренировки — программы по группам мышц</div>
            <div class="feature">🥗 Питание — рецепты и планы</div>
          </div>
          <div class="price">5 000 ₸ <span class="price-note">разовая оплата</span></div>
          <button class="btn btn-primary" @click="goToConfirm">Оплатить 💳</button>
        </div>

        <!-- Step 2: Confirmation -->
        <div v-else-if="step === 'confirm'" key="confirm" class="confirm-card">
          <h2>Подтверждение 🔐</h2>
          <p>Полный доступ к платформе</p>
          <div class="price">5 000 ₸</div>
          <div class="btn-row">
            <button class="btn btn-primary" @click="pay">Подтвердить ✅</button>
            <button class="btn btn-secondary" @click="goToProduct">Отмена</button>
          </div>
        </div>

        <!-- Step 3: Processing (instant providers) -->
        <div v-else-if="step === 'processing'" key="processing" class="processing-card">
          <div class="spinner"></div>
          <p>Обработка платежа... ⏳</p>
        </div>

        <!-- Step 3b: Waiting for Robokassa -->
        <div v-else-if="step === 'waiting'" key="waiting" class="processing-card">
          <div class="spinner"></div>
          <h2>Завершите оплату</h2>
          <p>Окно оплаты открыто. После завершения вернитесь сюда.</p>
          <button class="btn btn-secondary" style="margin-top: 16px" @click="cancelWaiting">Отмена</button>
        </div>

        <!-- Step 4: Success -->
        <div v-else-if="step === 'success'" key="success" class="success-card">
          <div class="success-icon">🎉</div>
          <h2>Оплата прошла успешно!</h2>
          <p>Полный доступ ко всем модулям открыт. 💪</p>
          <button class="btn btn-primary" @click="goToModules">Перейти к модулям</button>
        </div>

        <!-- Error -->
        <div v-else-if="step === 'error'" key="error" class="error-card">
          <div class="error-icon">❌</div>
          <h2>Ошибка</h2>
          <p>{{ errorMsg }}</p>
          <button class="btn btn-primary" @click="goToProduct">Попробовать снова</button>
        </div>
      </Transition>
    </div>

    <ConfettiCanvas :active="showConfetti" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { api } from '../api'
import { useTelegram } from '../composables/useTelegram'
import SkeletonCard from '../components/SkeletonCard.vue'
import ConfettiCanvas from '../components/ConfettiCanvas.vue'

const router = useRouter()
const route = useRoute()
const { hapticImpact, hapticNotification } = useTelegram()
const loading = ref(true)
const paid = ref(false)
const step = ref<'product' | 'confirm' | 'processing' | 'waiting' | 'success' | 'error'>('product')
const errorMsg = ref('')
const showConfetti = ref(false)
let pollTimer: ReturnType<typeof setInterval> | null = null

onMounted(async () => {
  try {
    const status = await api.getPaymentStatus()
    paid.value = status.is_paid
  } catch {
    paid.value = false
  } finally {
    loading.value = false
  }

  // Handle redirect back from Robokassa
  const queryStatus = route.query.status as string
  if (queryStatus === 'success' && !paid.value) {
    step.value = 'waiting'
    startPolling()
  } else if (queryStatus === 'fail') {
    errorMsg.value = 'Оплата отменена или не прошла'
    step.value = 'error'
    hapticNotification('error')
  }
})

onUnmounted(() => {
  stopPolling()
})

function goToConfirm() {
  hapticImpact('light')
  step.value = 'confirm'
}

function goToProduct() {
  hapticImpact('light')
  step.value = 'product'
}

async function pay() {
  hapticImpact('medium')
  step.value = 'processing'
  try {
    const result = await api.processPayment()
    if (result.redirect_url) {
      // Async provider (Robokassa) — open payment page and start polling
      step.value = 'waiting'
      window.Telegram?.WebApp?.openLink(result.redirect_url)
      startPolling()
    } else if (result.success) {
      // Sync provider (Dummy) — payment completed instantly
      step.value = 'success'
      paid.value = true
      showConfetti.value = true
      hapticNotification('success')
    } else {
      errorMsg.value = result.message || 'Неизвестная ошибка'
      step.value = 'error'
      hapticNotification('error')
    }
  } catch (e: any) {
    errorMsg.value = e.message || 'Ошибка соединения'
    step.value = 'error'
    hapticNotification('error')
  }
}

function startPolling() {
  stopPolling()
  pollTimer = setInterval(async () => {
    try {
      const status = await api.getPaymentStatus()
      if (status.is_paid) {
        stopPolling()
        paid.value = true
        step.value = 'success'
        showConfetti.value = true
        hapticNotification('success')
      }
    } catch {
      // Ignore polling errors, keep trying
    }
  }, 3000)
}

function stopPolling() {
  if (pollTimer) {
    clearInterval(pollTimer)
    pollTimer = null
  }
}

function cancelWaiting() {
  stopPolling()
  goToProduct()
}

function goToModules() {
  router.push('/')
}
</script>

<style scoped>
.payment-page {
  max-width: 400px;
  margin: 0 auto;
}

.skeleton-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-top: 20px;
}

.product-card, .confirm-card, .success-card, .error-card, .processing-card {
  text-align: center;
  padding: 24px;
  background: var(--secondary-bg);
  border-radius: 12px;
  margin-top: 20px;
  position: relative;
  overflow: hidden;
}

.product-gradient {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, var(--button-color), #34c759, var(--link-color));
}

.features {
  text-align: left;
  margin: 16px 0;
}

.feature {
  padding: 8px 0;
  border-bottom: 1px solid rgba(0,0,0,0.05);
  font-size: 14px;
}

.price {
  font-size: 28px;
  font-weight: bold;
  margin: 16px 0;
  color: var(--button-color);
}

.price-note {
  font-size: 14px;
  font-weight: normal;
  color: var(--hint-color);
}

.btn-row {
  display: flex;
  gap: 12px;
  justify-content: center;
  margin-top: 16px;
}

.btn-primary {
  background: var(--button-color);
  color: var(--button-text-color);
}

.btn-secondary {
  background: var(--secondary-bg);
  color: var(--text-color);
  border: 1px solid var(--hint-color);
}

.success-icon, .error-icon {
  font-size: 48px;
  margin-bottom: 12px;
}

.processing-card p {
  margin-top: 16px;
  color: var(--hint-color);
}

h2 {
  margin-bottom: 8px;
}

p {
  color: var(--hint-color);
  margin-bottom: 8px;
}
</style>
