<template>
  <div class="payment-page">
    <div v-if="loading" class="loading">
      <div class="spinner"></div>
    </div>

    <div v-else-if="paid" class="success-card">
      <div class="success-icon">✅</div>
      <h2>Доступ оплачен</h2>
      <p>Все модули доступны. Приятных тренировок!</p>
      <button class="btn" @click="goToModules">Перейти к модулям</button>
    </div>

    <div v-else class="payment-flow">
      <!-- Step 1: Product card -->
      <div v-if="step === 'product'" class="product-card">
        <h2>🏋️ Полный доступ</h2>
        <div class="features">
          <div class="feature">🏥 ЛФК — упражнения при проблемах со здоровьем</div>
          <div class="feature">💪 Тренировки — программы по группам мышц</div>
          <div class="feature">🥗 Питание — рецепты и планы</div>
        </div>
        <div class="price">5 000 ₸ <span class="price-note">разовая оплата</span></div>
        <button class="btn btn-primary" @click="step = 'confirm'">Оплатить</button>
      </div>

      <!-- Step 2: Confirmation -->
      <div v-else-if="step === 'confirm'" class="confirm-card">
        <h2>Подтверждение</h2>
        <p>Полный доступ к платформе</p>
        <div class="price">5 000 ₸</div>
        <div class="btn-row">
          <button class="btn btn-primary" @click="pay">Подтвердить</button>
          <button class="btn btn-secondary" @click="step = 'product'">Отмена</button>
        </div>
      </div>

      <!-- Step 3: Processing -->
      <div v-else-if="step === 'processing'" class="processing-card">
        <div class="spinner"></div>
        <p>Обработка платежа...</p>
      </div>

      <!-- Step 4: Success -->
      <div v-else-if="step === 'success'" class="success-card">
        <div class="success-icon">✅</div>
        <h2>Оплата прошла успешно!</h2>
        <p>Полный доступ ко всем модулям открыт.</p>
        <button class="btn btn-primary" @click="goToModules">Перейти к модулям</button>
      </div>

      <!-- Error -->
      <div v-else-if="step === 'error'" class="error-card">
        <div class="error-icon">❌</div>
        <h2>Ошибка</h2>
        <p>{{ errorMsg }}</p>
        <button class="btn" @click="step = 'product'">Попробовать снова</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../api'

const router = useRouter()
const loading = ref(true)
const paid = ref(false)
const step = ref<'product' | 'confirm' | 'processing' | 'success' | 'error'>('product')
const errorMsg = ref('')

onMounted(async () => {
  try {
    const status = await api.getPaymentStatus()
    paid.value = status.is_paid
  } catch {
    paid.value = false
  } finally {
    loading.value = false
  }
})

async function pay() {
  step.value = 'processing'
  try {
    const result = await api.processPayment()
    if (result.success) {
      step.value = 'success'
      paid.value = true
    } else {
      errorMsg.value = result.message || 'Неизвестная ошибка'
      step.value = 'error'
    }
  } catch (e: any) {
    errorMsg.value = e.message || 'Ошибка соединения'
    step.value = 'error'
  }
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

.product-card, .confirm-card, .success-card, .error-card, .processing-card {
  text-align: center;
  padding: 24px;
  background: var(--secondary-bg);
  border-radius: 12px;
  margin-top: 20px;
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
