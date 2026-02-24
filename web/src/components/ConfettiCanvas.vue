<template>
  <canvas ref="canvas" class="confetti-canvas" v-show="active"></canvas>
</template>

<script setup lang="ts">
import { ref, watch, onUnmounted } from 'vue'

const props = withDefaults(defineProps<{
  active: boolean
  duration?: number
}>(), {
  duration: 3000,
})

const canvas = ref<HTMLCanvasElement>()

interface Particle {
  x: number; y: number
  vx: number; vy: number
  size: number; color: string
  rotation: number; rotationSpeed: number
  opacity: number
}

let particles: Particle[] = []
let animId = 0
let generating = false

const colors = ['#FF6B6B', '#4ECDC4', '#45B7D1', '#96CEB4', '#FFEAA7', '#DDA0DD', '#98D8C8', '#F7DC6F']

function createParticle(): Particle {
  return {
    x: Math.random() * (canvas.value?.width || 300),
    y: -10,
    vx: (Math.random() - 0.5) * 6,
    vy: Math.random() * 3 + 2,
    size: Math.random() * 8 + 4,
    color: colors[Math.floor(Math.random() * colors.length)],
    rotation: Math.random() * 360,
    rotationSpeed: (Math.random() - 0.5) * 10,
    opacity: 1,
  }
}

function animate() {
  const c = canvas.value
  if (!c) return
  const ctx = c.getContext('2d')
  if (!ctx) return

  c.width = window.innerWidth
  c.height = window.innerHeight

  ctx.clearRect(0, 0, c.width, c.height)

  if (generating && particles.length < 150) {
    for (let i = 0; i < 3; i++) particles.push(createParticle())
  }

  particles = particles.filter(p => {
    p.x += p.vx
    p.y += p.vy
    p.vy += 0.1 // gravity
    p.rotation += p.rotationSpeed
    p.opacity -= 0.003

    if (p.opacity <= 0 || p.y > c.height + 20) return false

    ctx.save()
    ctx.translate(p.x, p.y)
    ctx.rotate((p.rotation * Math.PI) / 180)
    ctx.globalAlpha = p.opacity
    ctx.fillStyle = p.color
    ctx.fillRect(-p.size / 2, -p.size / 2, p.size, p.size * 0.6)
    ctx.restore()

    return true
  })

  if (particles.length > 0 || generating) {
    animId = requestAnimationFrame(animate)
  }
}

function start() {
  particles = []
  generating = true
  animate()
  setTimeout(() => {
    generating = false
  }, props.duration)
}

function stop() {
  generating = false
  cancelAnimationFrame(animId)
  particles = []
}

watch(() => props.active, (val) => {
  if (val) start()
  else stop()
})

onUnmounted(() => {
  cancelAnimationFrame(animId)
})
</script>

<style scoped>
.confetti-canvas {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  pointer-events: none;
  z-index: 9999;
}
</style>
