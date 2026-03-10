<template>
  <div class="health-slide">
    <h2 class="animate-fade-up">Здоровье и ограничения 🏥</h2>
    <p class="subtitle animate-fade-up delay-1">Это поможет подобрать безопасный план</p>

    <div class="section animate-fade-up delay-2">
      <div class="section-label">Есть ли боли?</div>
      <div class="chips">
        <div
          v-for="loc in painOptions"
          :key="loc.value"
          class="chip"
          :class="{ selected: painLocations.includes(loc.value) }"
          @click="togglePain(loc.value)"
        >
          {{ loc.label }}
        </div>
      </div>
    </div>

    <div v-if="painLocations.length > 0" class="section animate-fade-up">
      <div class="section-label">Уровень боли: {{ painLevel }}</div>
      <input
        type="range"
        min="0"
        max="10"
        :value="painLevel"
        @input="$emit('update:painLevel', Number(($event.target as HTMLInputElement).value))"
        class="pain-slider"
      />
      <div class="slider-labels">
        <span>Нет</span>
        <span>Сильная</span>
      </div>
    </div>

    <div class="section animate-fade-up delay-3">
      <div class="section-label">Диагнозы</div>
      <div class="chips">
        <div
          v-for="d in diagnosisOptions"
          :key="d.value"
          class="chip"
          :class="{ selected: diagnoses.includes(d.value) }"
          @click="toggleDiagnosis(d.value)"
        >
          {{ d.label }}
        </div>
      </div>
    </div>

    <div class="disclaimer animate-fade-up delay-3">
      Приложение не заменяет врача. При острой боли обратитесь к специалисту.
    </div>
  </div>
</template>

<script setup lang="ts">
import { useTelegram } from '../../composables/useTelegram'

const props = defineProps<{
  painLocations: string[]
  painLevel: number
  diagnoses: string[]
}>()

const emit = defineEmits<{
  'update:painLocations': [value: string[]]
  'update:painLevel': [value: number]
  'update:diagnoses': [value: string[]]
}>()

const { hapticSelection } = useTelegram()

const painOptions = [
  { value: 'shoulder', label: 'Плечо' },
  { value: 'knee', label: 'Колено' },
  { value: 'back', label: 'Спина' },
  { value: 'neck', label: 'Шея' },
  { value: 'other', label: 'Другое' },
]

const diagnosisOptions = [
  { value: 'hernia', label: 'Грыжа' },
  { value: 'protrusion', label: 'Протрузия' },
  { value: 'scoliosis', label: 'Сколиоз' },
  { value: 'kyphosis', label: 'Кифоз' },
  { value: 'lordosis', label: 'Лордоз' },
  { value: 'other', label: 'Другое' },
]

function togglePain(val: string) {
  hapticSelection()
  const arr = [...props.painLocations]
  const idx = arr.indexOf(val)
  if (idx >= 0) arr.splice(idx, 1)
  else arr.push(val)
  emit('update:painLocations', arr)
}

function toggleDiagnosis(val: string) {
  hapticSelection()
  const arr = [...props.diagnoses]
  const idx = arr.indexOf(val)
  if (idx >= 0) arr.splice(idx, 1)
  else arr.push(val)
  emit('update:diagnoses', arr)
}
</script>

<style scoped>
.health-slide { padding: 20px 0; }
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

.pain-slider {
  width: 100%;
  height: 6px;
  -webkit-appearance: none;
  appearance: none;
  border-radius: 3px;
  background: linear-gradient(to right, #34c759, #ffcc00, #ff3b30);
  outline: none;
}
.pain-slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: white;
  box-shadow: 0 2px 6px rgba(0,0,0,0.2);
  cursor: pointer;
}
.slider-labels {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: var(--hint-color);
  margin-top: 4px;
}

.disclaimer {
  margin-top: 16px;
  padding: 12px;
  border-radius: 8px;
  background: color-mix(in srgb, #ff3b30 10%, transparent);
  color: #ff3b30;
  font-size: 13px;
  line-height: 1.4;
}

.animate-fade-up { opacity: 0; animation: fadeUp 0.4s ease forwards; }
.delay-1 { animation-delay: 0.1s; }
.delay-2 { animation-delay: 0.15s; }
.delay-3 { animation-delay: 0.25s; }
@keyframes fadeUp { from { opacity:0; transform:translateY(12px); } to { opacity:1; transform:translateY(0); } }
</style>
