<template>
  <div class="level-slide">
    <h2 class="animate-fade-up">Твой уровень подготовки 🏋️</h2>
    <p class="subtitle animate-fade-up delay-1">Это поможет подобрать нагрузку</p>

    <div class="level-cards animate-fade-up delay-2">
      <div
        class="level-card"
        :class="{ selected: modelValue === 'beginner' }"
        @click="select('beginner')"
      >
        <div class="level-icon">🌱</div>
        <div class="level-title">Новичок</div>
        <div class="level-desc">Мало или нет опыта тренировок</div>
      </div>
      <div
        class="level-card"
        :class="{ selected: modelValue === 'intermediate' }"
        @click="select('intermediate')"
      >
        <div class="level-icon">⚡</div>
        <div class="level-title">Средний</div>
        <div class="level-desc">Есть опыт, тренируюсь периодически</div>
      </div>
      <div
        class="level-card"
        :class="{ selected: modelValue === 'advanced' }"
        @click="select('advanced')"
      >
        <div class="level-icon">🔥</div>
        <div class="level-title">Продвинутый</div>
        <div class="level-desc">Регулярные тренировки, хороший опыт</div>
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
.level-slide { text-align: center; padding: 20px 0; }
h2 { font-size: 24px; margin-bottom: 8px; }
.subtitle { font-size: 15px; color: var(--hint-color); margin-bottom: 32px; }

.level-cards {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.level-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 18px 20px;
  border-radius: 16px;
  background: var(--secondary-bg);
  border: 2px solid transparent;
  cursor: pointer;
  transition: all 0.2s ease;
  text-align: left;
}

.level-card:active { transform: scale(0.97); }

.level-card.selected {
  background: color-mix(in srgb, var(--button-color) 12%, transparent);
  border-color: var(--button-color);
}

.level-icon { font-size: 32px; flex-shrink: 0; }
.level-title { font-size: 17px; font-weight: 700; margin-bottom: 2px; }
.level-desc { font-size: 13px; color: var(--hint-color); }

.level-card.selected .level-title { color: var(--button-color); }

.animate-fade-up { opacity: 0; animation: fadeUp 0.4s ease forwards; }
.delay-1 { animation-delay: 0.1s; }
.delay-2 { animation-delay: 0.2s; }
@keyframes fadeUp { from { opacity:0; transform:translateY(12px); } to { opacity:1; transform:translateY(0); } }
</style>
