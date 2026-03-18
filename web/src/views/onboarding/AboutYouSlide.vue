<template>
  <div class="about-slide">
    <h2 class="animate-fade-up">Расскажи о себе 👤</h2>

    <div class="cards animate-fade-up delay-1">
      <div
        class="gender-card"
        :class="{ selected: gender === 'male' }"
        @click="selectGender('male')"
      >
        <span class="card-emoji">🙋‍♂️</span>
        <span class="card-label">Мужской</span>
      </div>
      <div
        class="gender-card"
        :class="{ selected: gender === 'female' }"
        @click="selectGender('female')"
      >
        <span class="card-emoji">🙋‍♀️</span>
        <span class="card-label">Женский</span>
      </div>
    </div>

    <div class="age-section animate-fade-up delay-2">
      <div class="section-label">Сколько тебе лет?</div>
      <div class="picker-value">
        <span class="big-number">{{ age }}</span>
        <span class="unit">лет</span>
      </div>
      <WheelPicker
        :items="ageItems"
        :modelValue="age"
        @update:modelValue="$emit('update:age', $event)"
        :visibleItems="3"
        :itemHeight="40"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import WheelPicker from '../../components/WheelPicker.vue'
import { useTelegram } from '../../composables/useTelegram'

defineProps<{
  gender: string
  age: number
}>()

const emit = defineEmits<{
  'update:gender': [value: string]
  'update:age': [value: number | string]
}>()

const { hapticImpact } = useTelegram()

function selectGender(val: string) {
  hapticImpact('medium')
  emit('update:gender', val)
}

const ageItems = Array.from({ length: 91 }, (_, i) => ({
  value: i + 10,
  label: String(i + 10),
}))
</script>

<style scoped>
.about-slide { text-align: center; padding: 20px 0; }
h2 { font-size: 22px; margin-bottom: 24px; }

.cards {
  display: flex;
  gap: 12px;
  justify-content: center;
  margin-bottom: 28px;
}

.gender-card {
  flex: 1;
  max-width: 150px;
  padding: 20px 16px;
  background: var(--secondary-bg);
  border-radius: 16px;
  border: 3px solid transparent;
  cursor: pointer;
  transition: transform 0.2s, border-color 0.2s, background 0.2s;
}

.gender-card:active { transform: scale(0.95); }

.gender-card.selected {
  border-color: var(--button-color);
  background: color-mix(in srgb, var(--button-color) 10%, var(--secondary-bg));
}

.card-emoji { font-size: 40px; display: block; margin-bottom: 8px; }
.card-label { font-size: 15px; font-weight: 600; }

.age-section { margin-top: 4px; }
.section-label { font-size: 14px; font-weight: 600; color: var(--hint-color); margin-bottom: 8px; }
.picker-value { margin-bottom: 8px; }
.big-number { font-size: 36px; font-weight: 700; color: var(--button-color); }
.unit { font-size: 18px; color: var(--hint-color); margin-left: 4px; }

.animate-fade-up { opacity: 0; animation: fadeUp 0.4s ease forwards; }
.delay-1 { animation-delay: 0.15s; }
.delay-2 { animation-delay: 0.25s; }
@keyframes fadeUp { from { opacity:0; transform:translateY(12px); } to { opacity:1; transform:translateY(0); } }
</style>
