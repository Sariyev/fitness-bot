<template>
  <div class="gender-slide">
    <h2 class="animate-fade-up">Укажи свой пол</h2>
    <div class="cards">
      <div
        class="gender-card animate-fade-up delay-1"
        :class="{ selected: modelValue === 'male' }"
        @click="select('male')"
      >
        <span class="card-emoji">🙋‍♂️</span>
        <span class="card-label">Мужской</span>
      </div>
      <div
        class="gender-card animate-fade-up delay-2"
        :class="{ selected: modelValue === 'female' }"
        @click="select('female')"
      >
        <span class="card-emoji">🙋‍♀️</span>
        <span class="card-label">Женский</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useTelegram } from '../../composables/useTelegram'

defineProps<{ modelValue: string }>()
const emit = defineEmits<{ 'update:modelValue': [value: string] }>()

const { hapticImpact } = useTelegram()

function select(val: string) {
  hapticImpact('medium')
  emit('update:modelValue', val)
}
</script>

<style scoped>
.gender-slide { text-align: center; padding: 20px 0; }
h2 { font-size: 22px; margin-bottom: 32px; }

.cards {
  display: flex;
  gap: 12px;
  justify-content: center;
}

.gender-card {
  flex: 1;
  max-width: 160px;
  padding: 28px 16px;
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

.card-emoji { font-size: 48px; display: block; margin-bottom: 12px; }
.card-label { font-size: 16px; font-weight: 600; }

.animate-fade-up { opacity: 0; animation: fadeUp 0.4s ease forwards; }
.delay-1 { animation-delay: 0.15s; }
.delay-2 { animation-delay: 0.25s; }
@keyframes fadeUp { from { opacity:0; transform:translateY(12px); } to { opacity:1; transform:translateY(0); } }
</style>
