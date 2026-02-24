<template>
  <div class="goals-slide">
    <h2 class="animate-fade-up">Выбери свои цели 🎯</h2>
    <p class="subtitle animate-fade-up delay-1">Можно несколько</p>

    <div class="section animate-fade-up delay-2">
      <div class="section-label">💪 Фитнес</div>
      <div class="chips">
        <div
          v-for="g in fitnessGoals"
          :key="g.value"
          class="chip"
          :class="{ selected: modelValue.includes(g.value) }"
          @click="toggle(g.value)"
        >
          {{ g.label }}
        </div>
      </div>
    </div>

    <div class="section animate-fade-up delay-3">
      <div class="section-label">🏥 Здоровье</div>
      <div class="chips">
        <div
          v-for="g in healthGoals"
          :key="g.value"
          class="chip"
          :class="{ selected: modelValue.includes(g.value) }"
          @click="toggle(g.value)"
        >
          {{ g.label }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useTelegram } from '../../composables/useTelegram'

const props = defineProps<{ modelValue: string[] }>()
const emit = defineEmits<{ 'update:modelValue': [value: string[]] }>()

const { hapticSelection } = useTelegram()

const fitnessGoals = [
  { value: 'weight_loss', label: 'Похудеть' },
  { value: 'muscle_gain', label: 'Набрать массу' },
  { value: 'strength', label: 'Больше силы' },
  { value: 'endurance', label: 'Выносливость' },
  { value: 'maintenance', label: 'Поддержание формы' },
]

const healthGoals = [
  { value: 'hernia', label: 'Грыжа' },
  { value: 'protrusion', label: 'Протрузии' },
  { value: 'scoliosis', label: 'Сколиоз' },
  { value: 'kyphosis', label: 'Кифоз' },
  { value: 'lordosis', label: 'Лордоз' },
]

function toggle(val: string) {
  hapticSelection()
  const current = [...props.modelValue]
  const idx = current.indexOf(val)
  if (idx >= 0) {
    current.splice(idx, 1)
  } else {
    current.push(val)
  }
  emit('update:modelValue', current)
}
</script>

<style scoped>
.goals-slide { padding: 20px 0; }
h2 { font-size: 22px; text-align: center; margin-bottom: 4px; }
.subtitle { text-align: center; color: var(--hint-color); font-size: 14px; margin-bottom: 24px; }

.section { margin-bottom: 20px; }
.section-label { font-size: 14px; font-weight: 600; color: var(--hint-color); margin-bottom: 10px; }

.chips { display: flex; flex-wrap: wrap; gap: 8px; }

.chip {
  padding: 10px 16px;
  border-radius: 20px;
  font-size: 14px;
  background: var(--secondary-bg);
  border: 2px solid transparent;
  cursor: pointer;
  transition: all 0.2s ease;
}

.chip:active { transform: scale(0.95); }

.chip.selected {
  background: color-mix(in srgb, var(--button-color) 15%, transparent);
  border-color: var(--button-color);
  color: var(--button-color);
  font-weight: 600;
  transform: scale(1.02);
}

.animate-fade-up { opacity: 0; animation: fadeUp 0.4s ease forwards; }
.delay-1 { animation-delay: 0.1s; }
.delay-2 { animation-delay: 0.15s; }
.delay-3 { animation-delay: 0.25s; }
@keyframes fadeUp { from { opacity:0; transform:translateY(12px); } to { opacity:1; transform:translateY(0); } }
</style>
