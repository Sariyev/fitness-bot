<template>
  <div class="schedule-slide">
    <h2 class="animate-fade-up">Режим тренировок ⏰</h2>
    <p class="subtitle animate-fade-up delay-1">Настроим под твой график</p>

    <div class="section animate-fade-up delay-2">
      <div class="section-label">Где тренируешься?</div>
      <div class="chips">
        <div
          v-for="opt in accessOptions"
          :key="opt.value"
          class="chip"
          :class="{ selected: trainingAccess === opt.value }"
          @click="$emit('update:trainingAccess', opt.value); hapticSelection()"
        >
          {{ opt.label }}
        </div>
      </div>
    </div>

    <div class="section animate-fade-up delay-2">
      <div class="section-label">Дней в неделю: {{ daysPerWeek }}</div>
      <input
        type="range"
        min="2"
        max="6"
        :value="daysPerWeek"
        @input="$emit('update:daysPerWeek', Number(($event.target as HTMLInputElement).value))"
        class="range-slider"
      />
      <div class="slider-labels">
        <span>2</span><span>3</span><span>4</span><span>5</span><span>6</span>
      </div>
    </div>

    <div class="section animate-fade-up delay-3">
      <div class="section-label">Длительность занятия</div>
      <div class="chips">
        <div
          v-for="dur in durationOptions"
          :key="dur"
          class="chip"
          :class="{ selected: sessionDuration === dur }"
          @click="$emit('update:sessionDuration', dur); hapticSelection()"
        >
          {{ dur }} мин
        </div>
      </div>
    </div>

    <div class="section animate-fade-up delay-3">
      <div class="section-label">Предпочтительное время</div>
      <div class="chips">
        <div
          v-for="t in timeOptions"
          :key="t.value"
          class="chip"
          :class="{ selected: preferredTime === t.value }"
          @click="$emit('update:preferredTime', t.value); hapticSelection()"
        >
          {{ t.label }}
        </div>
      </div>
    </div>

    <div class="section animate-fade-up delay-3">
      <div class="section-label">Инвентарь</div>
      <div class="chips">
        <div
          v-for="eq in equipmentOptions"
          :key="eq.value"
          class="chip"
          :class="{ selected: equipment.includes(eq.value) }"
          @click="toggleEquipment(eq.value)"
        >
          {{ eq.label }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useTelegram } from '../../composables/useTelegram'

const props = defineProps<{
  trainingAccess: string
  daysPerWeek: number
  sessionDuration: number
  preferredTime: string
  equipment: string[]
}>()

const emit = defineEmits<{
  'update:trainingAccess': [value: string]
  'update:daysPerWeek': [value: number]
  'update:sessionDuration': [value: number]
  'update:preferredTime': [value: string]
  'update:equipment': [value: string[]]
}>()

const { hapticSelection } = useTelegram()

const accessOptions = [
  { value: 'home', label: 'Дом' },
  { value: 'gym', label: 'Зал' },
  { value: 'both', label: 'Дом + Зал' },
]

const durationOptions = [15, 25, 35, 45, 60]

const timeOptions = [
  { value: 'morning', label: 'Утро' },
  { value: 'afternoon', label: 'День' },
  { value: 'evening', label: 'Вечер' },
]

const equipmentOptions = [
  { value: 'mat', label: 'Коврик' },
  { value: 'dumbbells', label: 'Гантели' },
  { value: 'bands', label: 'Резинки' },
  { value: 'barbell', label: 'Штанга' },
  { value: 'pullup_bar', label: 'Турник' },
]

function toggleEquipment(val: string) {
  hapticSelection()
  const arr = [...props.equipment]
  const idx = arr.indexOf(val)
  if (idx >= 0) arr.splice(idx, 1)
  else arr.push(val)
  emit('update:equipment', arr)
}
</script>

<style scoped>
.schedule-slide { padding: 20px 0; }
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
}

.range-slider {
  width: 100%;
  height: 6px;
  -webkit-appearance: none;
  appearance: none;
  border-radius: 3px;
  background: var(--secondary-bg);
  outline: none;
}
.range-slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: var(--button-color);
  cursor: pointer;
}
.slider-labels {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: var(--hint-color);
  margin-top: 4px;
  padding: 0 4px;
}

.animate-fade-up { opacity: 0; animation: fadeUp 0.4s ease forwards; }
.delay-1 { animation-delay: 0.1s; }
.delay-2 { animation-delay: 0.15s; }
.delay-3 { animation-delay: 0.25s; }
@keyframes fadeUp { from { opacity:0; transform:translateY(12px); } to { opacity:1; transform:translateY(0); } }
</style>
