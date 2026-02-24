<template>
  <div class="wheel-picker" ref="container"
    @touchstart="onTouchStart"
    @touchmove="onTouchMove"
    @touchend="onTouchEnd"
    @mousedown="onMouseDown"
  >
    <div class="wheel-picker__highlight"></div>
    <div class="wheel-picker__items" ref="itemsEl" :style="itemsStyle">
      <div
        v-for="(item, i) in items"
        :key="i"
        class="wheel-picker__item"
        :style="{ height: itemHeight + 'px', lineHeight: itemHeight + 'px' }"
        :class="{ active: item.value === modelValue }"
      >
        {{ item.label }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useTelegram } from '../composables/useTelegram'

interface PickerItem {
  value: number | string
  label: string
}

const props = withDefaults(defineProps<{
  items: PickerItem[]
  modelValue: number | string
  visibleItems?: number
  itemHeight?: number
}>(), {
  visibleItems: 5,
  itemHeight: 44,
})

const emit = defineEmits<{
  'update:modelValue': [value: number | string]
}>()

const { hapticSelection } = useTelegram()

const container = ref<HTMLElement>()
const itemsEl = ref<HTMLElement>()
const offset = ref(0)
const isDragging = ref(false)
const startY = ref(0)
const startOffset = ref(0)
const velocity = ref(0)
const lastY = ref(0)
const lastTime = ref(0)
const animating = ref(false)
const lastEmittedIndex = ref(-1)

const containerHeight = computed(() => props.itemHeight * props.visibleItems)
const paddingTop = computed(() => props.itemHeight * Math.floor(props.visibleItems / 2))

const itemsStyle = computed(() => ({
  transform: `translateY(${offset.value + paddingTop.value}px)`,
  transition: isDragging.value ? 'none' : 'transform 0.3s cubic-bezier(0.2, 0, 0.2, 1)',
}))

function getSelectedIndex(): number {
  const idx = Math.round(-offset.value / props.itemHeight)
  return Math.max(0, Math.min(props.items.length - 1, idx))
}

function snapToIndex(index: number) {
  offset.value = -index * props.itemHeight
  const item = props.items[index]
  if (item && item.value !== props.modelValue) {
    emit('update:modelValue', item.value)
  }
}

function onTouchStart(e: TouchEvent) {
  isDragging.value = true
  startY.value = e.touches[0].clientY
  startOffset.value = offset.value
  velocity.value = 0
  lastY.value = e.touches[0].clientY
  lastTime.value = Date.now()
}

function onTouchMove(e: TouchEvent) {
  if (!isDragging.value) return
  e.preventDefault()
  const y = e.touches[0].clientY
  const delta = y - startY.value
  const now = Date.now()
  const dt = now - lastTime.value
  if (dt > 0) {
    velocity.value = (y - lastY.value) / dt
  }
  lastY.value = y
  lastTime.value = now
  offset.value = startOffset.value + delta

  // Haptic on crossing item boundary
  const idx = getSelectedIndex()
  if (idx !== lastEmittedIndex.value) {
    lastEmittedIndex.value = idx
    hapticSelection()
  }
}

function onTouchEnd() {
  isDragging.value = false
  const momentum = velocity.value * 150
  offset.value += momentum
  const idx = getSelectedIndex()
  snapToIndex(idx)
}

// Mouse support for desktop testing
function onMouseDown(e: MouseEvent) {
  isDragging.value = true
  startY.value = e.clientY
  startOffset.value = offset.value
  velocity.value = 0
  lastY.value = e.clientY
  lastTime.value = Date.now()

  const onMouseMove = (e: MouseEvent) => {
    if (!isDragging.value) return
    const y = e.clientY
    const delta = y - startY.value
    const now = Date.now()
    const dt = now - lastTime.value
    if (dt > 0) velocity.value = (y - lastY.value) / dt
    lastY.value = y
    lastTime.value = now
    offset.value = startOffset.value + delta
    const idx = getSelectedIndex()
    if (idx !== lastEmittedIndex.value) {
      lastEmittedIndex.value = idx
      hapticSelection()
    }
  }

  const onMouseUp = () => {
    isDragging.value = false
    const momentum = velocity.value * 150
    offset.value += momentum
    snapToIndex(getSelectedIndex())
    window.removeEventListener('mousemove', onMouseMove)
    window.removeEventListener('mouseup', onMouseUp)
  }

  window.addEventListener('mousemove', onMouseMove)
  window.addEventListener('mouseup', onMouseUp)
}

// Initialize position
onMounted(() => {
  const idx = props.items.findIndex(item => item.value === props.modelValue)
  if (idx >= 0) {
    offset.value = -idx * props.itemHeight
    lastEmittedIndex.value = idx
  }
})

watch(() => props.modelValue, (val) => {
  const idx = props.items.findIndex(item => item.value === val)
  if (idx >= 0 && !isDragging.value) {
    offset.value = -idx * props.itemHeight
    lastEmittedIndex.value = idx
  }
})
</script>

<style scoped>
.wheel-picker {
  position: relative;
  overflow: hidden;
  height: v-bind("containerHeight + 'px'");
  touch-action: none;
  user-select: none;
  -webkit-mask-image: linear-gradient(
    to bottom,
    transparent 0%,
    black 25%,
    black 75%,
    transparent 100%
  );
  mask-image: linear-gradient(
    to bottom,
    transparent 0%,
    black 25%,
    black 75%,
    transparent 100%
  );
}

.wheel-picker__highlight {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  height: v-bind("itemHeight + 'px'");
  left: 0;
  right: 0;
  border-top: 2px solid var(--button-color);
  border-bottom: 2px solid var(--button-color);
  border-radius: 8px;
  background: color-mix(in srgb, var(--button-color) 8%, transparent);
  pointer-events: none;
  z-index: 1;
}

.wheel-picker__items {
  will-change: transform;
}

.wheel-picker__item {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  color: var(--hint-color);
  transition: color 0.15s, font-weight 0.15s;
}

.wheel-picker__item.active {
  color: var(--text-color);
  font-weight: 600;
  font-size: 22px;
}
</style>
