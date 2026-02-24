<template>
  <div class="progress-bar">
    <div class="progress-track">
      <div class="progress-fill" :style="{ width: animatedWidth + '%' }"></div>
    </div>
    <span class="progress-text">{{ completed }}/{{ total }}</span>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'

const props = defineProps<{
  completed: number
  total: number
}>()

const animated = ref(false)

const percentage = computed(() => {
  if (props.total === 0) return 0
  return Math.round((props.completed / props.total) * 100)
})

const animatedWidth = computed(() => animated.value ? percentage.value : 0)

onMounted(() => {
  requestAnimationFrame(() => {
    animated.value = true
  })
})
</script>

<style scoped>
.progress-bar {
  display: flex;
  align-items: center;
  gap: 8px;
}

.progress-track {
  flex: 1;
  height: 6px;
  background: var(--bg-color);
  border-radius: 3px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: var(--button-color);
  border-radius: 3px;
  transition: width 0.6s ease;
}

.progress-text {
  font-size: 12px;
  color: var(--hint-color);
  flex-shrink: 0;
  min-width: 32px;
  text-align: right;
}
</style>
