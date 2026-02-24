<template>
  <div class="picker-slide">
    <h2 class="animate-fade-up">Сколько тебе лет? 🎂</h2>
    <div class="picker-value animate-fade-up delay-1">
      <span class="big-number">{{ modelValue }}</span>
      <span class="unit">лет</span>
    </div>
    <WheelPicker
      :items="items"
      :modelValue="modelValue"
      @update:modelValue="$emit('update:modelValue', $event)"
    />
  </div>
</template>

<script setup lang="ts">
import WheelPicker from '../../components/WheelPicker.vue'

defineProps<{ modelValue: number }>()
defineEmits<{ 'update:modelValue': [value: number | string] }>()

const items = Array.from({ length: 91 }, (_, i) => ({
  value: i + 10,
  label: String(i + 10),
}))
</script>

<style scoped>
.picker-slide { text-align: center; padding: 20px 0; }
h2 { font-size: 22px; margin-bottom: 20px; }
.picker-value { margin-bottom: 20px; }
.big-number { font-size: 48px; font-weight: 700; color: var(--button-color); }
.unit { font-size: 20px; color: var(--hint-color); margin-left: 6px; }
.animate-fade-up { opacity: 0; animation: fadeUp 0.4s ease forwards; }
.delay-1 { animation-delay: 0.1s; }
@keyframes fadeUp { from { opacity:0; transform:translateY(12px); } to { opacity:1; transform:translateY(0); } }
</style>
