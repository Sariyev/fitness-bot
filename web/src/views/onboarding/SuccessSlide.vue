<template>
  <div class="success-slide">
    <ConfettiCanvas :active="showConfetti" />
    <div class="success-emoji animate-bounce-in">🎉</div>
    <h1 class="animate-fade-up delay-1">Отлично, ты готов!</h1>
    <p class="subtitle animate-fade-up delay-2">
      Андрей подберёт тебе персональную программу тренировок и питания
    </p>
    <div class="summary animate-fade-up delay-3">
      <div class="summary-item">🎂 {{ age }} лет</div>
      <div class="summary-item">📏 {{ heightCm }} см</div>
      <div class="summary-item">⚖️ {{ weightKg }} кг</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import ConfettiCanvas from '../../components/ConfettiCanvas.vue'
import { useTelegram } from '../../composables/useTelegram'

defineProps<{
  age: number
  heightCm: number
  weightKg: number
}>()

const { hapticNotification } = useTelegram()
const showConfetti = ref(false)

onMounted(() => {
  showConfetti.value = true
  hapticNotification('success')
})
</script>

<style scoped>
.success-slide { text-align: center; padding: 40px 0 20px; }

.success-emoji { font-size: 80px; margin-bottom: 16px; }

h1 { font-size: 28px; margin-bottom: 8px; }

.subtitle { color: var(--hint-color); font-size: 15px; margin-bottom: 32px; padding: 0 16px; }

.summary {
  display: flex;
  justify-content: center;
  gap: 16px;
}

.summary-item {
  padding: 10px 16px;
  background: var(--secondary-bg);
  border-radius: 12px;
  font-size: 14px;
  font-weight: 600;
}

.animate-bounce-in {
  animation: bounceIn 0.6s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
}
.animate-fade-up { opacity: 0; animation: fadeUp 0.5s ease forwards; }
.delay-1 { animation-delay: 0.2s; }
.delay-2 { animation-delay: 0.4s; }
.delay-3 { animation-delay: 0.6s; }

@keyframes bounceIn {
  0% { transform: scale(0.3); opacity: 0; }
  50% { transform: scale(1.05); }
  100% { transform: scale(1); opacity: 1; }
}
@keyframes fadeUp { from { opacity:0; transform:translateY(12px); } to { opacity:1; transform:translateY(0); } }
</style>
