<template>
  <div class="fitness-slide">
    <h2 class="animate-fade-up">Твой уровень подготовки 💪</h2>
    <div class="options">
      <div
        v-for="(opt, i) in options"
        :key="opt.value"
        class="fitness-card animate-fade-up"
        :style="{ animationDelay: (0.1 + i * 0.1) + 's' }"
        :class="{ selected: modelValue === opt.value }"
        @click="select(opt.value)"
      >
        <span class="card-emoji">{{ opt.emoji }}</span>
        <div>
          <strong>{{ opt.label }}</strong>
          <p>{{ opt.desc }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useTelegram } from '../../composables/useTelegram'

defineProps<{ modelValue: string }>()
const emit = defineEmits<{ 'update:modelValue': [value: string] }>()

const { hapticImpact } = useTelegram()

const options = [
  { value: 'beginner', emoji: '🌱', label: 'Новичок', desc: 'Только начинаю заниматься' },
  { value: 'intermediate', emoji: '⚡', label: 'Средний', desc: 'Занимаюсь регулярно' },
  { value: 'advanced', emoji: '🔥', label: 'Продвинутый', desc: 'Тренируюсь давно и серьёзно' },
]

function select(val: string) {
  hapticImpact('medium')
  emit('update:modelValue', val)
}
</script>

<style scoped>
.fitness-slide { text-align: center; padding: 20px 0; }
h2 { font-size: 22px; margin-bottom: 24px; }

.options { display: flex; flex-direction: column; gap: 12px; }

.fitness-card {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 16px;
  background: var(--secondary-bg);
  border-radius: 14px;
  text-align: left;
  border: 3px solid transparent;
  cursor: pointer;
  transition: transform 0.2s, border-color 0.2s, background 0.2s;
}

.fitness-card:active { transform: scale(0.97); }

.fitness-card.selected {
  border-color: var(--button-color);
  background: color-mix(in srgb, var(--button-color) 10%, var(--secondary-bg));
}

.card-emoji { font-size: 32px; flex-shrink: 0; }
.fitness-card strong { font-size: 16px; display: block; margin-bottom: 2px; }
.fitness-card p { font-size: 13px; color: var(--hint-color); margin: 0; }

.animate-fade-up { opacity: 0; animation: fadeUp 0.4s ease forwards; }
@keyframes fadeUp { from { opacity:0; transform:translateY(12px); } to { opacity:1; transform:translateY(0); } }
</style>
