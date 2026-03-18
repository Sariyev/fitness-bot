<template>
  <div class="body-slide">
    <h2 class="animate-fade-up">Физические параметры 💪</h2>

    <div class="pickers animate-fade-up delay-1">
      <div class="picker-col">
        <div class="picker-label">Рост</div>
        <div class="picker-value">
          <span class="num">{{ heightCm }}</span>
          <span class="unit">см</span>
        </div>
        <WheelPicker
          :items="heightItems"
          :modelValue="heightCm"
          @update:modelValue="$emit('update:heightCm', $event)"
          :visibleItems="3"
          :itemHeight="36"
        />
      </div>
      <div class="picker-col">
        <div class="picker-label">Вес</div>
        <div class="picker-value">
          <span class="num">{{ weightKg }}</span>
          <span class="unit">кг</span>
        </div>
        <WheelPicker
          :items="weightItems"
          :modelValue="weightKg"
          @update:modelValue="$emit('update:weightKg', $event)"
          :visibleItems="3"
          :itemHeight="36"
        />
      </div>
    </div>

    <div class="fitness-section animate-fade-up delay-2">
      <div class="section-label">Уровень подготовки</div>
      <div class="fitness-chips">
        <div
          class="fitness-chip"
          :class="{ selected: fitnessLevel === 'beginner' }"
          @click="selectFitness('beginner')"
        >🌱 Новичок</div>
        <div
          class="fitness-chip"
          :class="{ selected: fitnessLevel === 'intermediate' }"
          @click="selectFitness('intermediate')"
        >⚡ Средний</div>
        <div
          class="fitness-chip"
          :class="{ selected: fitnessLevel === 'advanced' }"
          @click="selectFitness('advanced')"
        >🔥 Продвинутый</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import WheelPicker from '../../components/WheelPicker.vue'
import { useTelegram } from '../../composables/useTelegram'

defineProps<{
  heightCm: number
  weightKg: number
  fitnessLevel: string
}>()

const emit = defineEmits<{
  'update:heightCm': [value: number | string]
  'update:weightKg': [value: number | string]
  'update:fitnessLevel': [value: string]
}>()

const { hapticImpact } = useTelegram()

function selectFitness(val: string) {
  hapticImpact('medium')
  emit('update:fitnessLevel', val)
}

const heightItems = Array.from({ length: 151 }, (_, i) => ({
  value: i + 100,
  label: String(i + 100),
}))

const weightItems = Array.from({ length: 341 }, (_, i) => {
  const val = 30 + i * 0.5
  return { value: val, label: val % 1 === 0 ? String(val) : val.toFixed(1) }
})
</script>

<style scoped>
.body-slide { text-align: center; padding: 20px 0; }
h2 { font-size: 22px; margin-bottom: 24px; }

.pickers {
  display: flex;
  gap: 16px;
  margin-bottom: 28px;
}

.picker-col {
  flex: 1;
  text-align: center;
}

.picker-label {
  font-size: 13px;
  font-weight: 600;
  color: var(--hint-color);
  margin-bottom: 4px;
}

.picker-value { margin-bottom: 4px; }
.num { font-size: 28px; font-weight: 700; color: var(--button-color); }
.unit { font-size: 14px; color: var(--hint-color); margin-left: 2px; }

.fitness-section { text-align: left; }
.section-label { font-size: 14px; font-weight: 600; color: var(--hint-color); margin-bottom: 10px; text-align: center; }

.fitness-chips {
  display: flex;
  gap: 8px;
}

.fitness-chip {
  flex: 1;
  padding: 12px 8px;
  border-radius: 12px;
  font-size: 13px;
  font-weight: 600;
  text-align: center;
  background: var(--secondary-bg);
  border: 2px solid transparent;
  cursor: pointer;
  transition: all 0.2s ease;
}

.fitness-chip:active { transform: scale(0.95); }

.fitness-chip.selected {
  background: color-mix(in srgb, var(--button-color) 15%, transparent);
  border-color: var(--button-color);
  color: var(--button-color);
}

.animate-fade-up { opacity: 0; animation: fadeUp 0.4s ease forwards; }
.delay-1 { animation-delay: 0.15s; }
.delay-2 { animation-delay: 0.25s; }
@keyframes fadeUp { from { opacity:0; transform:translateY(12px); } to { opacity:1; transform:translateY(0); } }
</style>
