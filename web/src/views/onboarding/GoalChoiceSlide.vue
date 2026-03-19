<template>
  <div class="goal-slide">
    <h2 class="animate-fade-up">Что тебя интересует? 🎯</h2>
    <p class="subtitle animate-fade-up delay-1">Выбери направление</p>

    <div class="goal-cards animate-fade-up delay-2">
      <div
        class="goal-card"
        :class="{ selected: modelValue === 'lfk' }"
        @click="select('lfk')"
      >
        <div class="goal-icon">🏥</div>
        <div class="goal-title">ЛФК</div>
        <div class="goal-desc">Лечебная физкультура, реабилитация, восстановление</div>
      </div>
      <div
        class="goal-card"
        :class="{ selected: modelValue === 'training' }"
        @click="select('training')"
      >
        <div class="goal-icon">💪</div>
        <div class="goal-title">Тренировки</div>
        <div class="goal-desc">Фитнес, набор массы, похудение, общая форма</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useTelegram } from '../../composables/useTelegram'

defineProps<{
  modelValue: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const { hapticImpact } = useTelegram()

function select(val: string) {
  hapticImpact('medium')
  emit('update:modelValue', val)
}
</script>

<style scoped>
.goal-slide { text-align: center; padding: 20px 0; }
h2 { font-size: 24px; margin-bottom: 8px; }
.subtitle { font-size: 15px; color: var(--hint-color); margin-bottom: 32px; }

.goal-cards {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.goal-card {
  padding: 24px 20px;
  border-radius: 16px;
  background: var(--secondary-bg);
  border: 2px solid transparent;
  cursor: pointer;
  transition: all 0.2s ease;
  text-align: left;
}

.goal-card:active { transform: scale(0.97); }

.goal-card.selected {
  background: color-mix(in srgb, var(--button-color) 12%, transparent);
  border-color: var(--button-color);
}

.goal-icon { font-size: 36px; margin-bottom: 12px; }
.goal-title { font-size: 20px; font-weight: 700; margin-bottom: 6px; }
.goal-desc { font-size: 14px; color: var(--hint-color); line-height: 1.4; }

.goal-card.selected .goal-title { color: var(--button-color); }

.animate-fade-up { opacity: 0; animation: fadeUp 0.4s ease forwards; }
.delay-1 { animation-delay: 0.1s; }
.delay-2 { animation-delay: 0.2s; }
@keyframes fadeUp { from { opacity:0; transform:translateY(12px); } to { opacity:1; transform:translateY(0); } }
</style>
