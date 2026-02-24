<template>
  <div class="onboarding">
    <!-- Progress dots -->
    <div class="progress-dots" v-if="currentSlide > 0 && currentSlide < 7">
      <span
        v-for="i in 7"
        :key="i"
        class="dot"
        :class="{ active: i - 1 <= currentSlide, current: i - 1 === currentSlide }"
      ></span>
    </div>

    <!-- Slides -->
    <Transition :name="slideDirection" mode="out-in">
      <WelcomeSlide v-if="currentSlide === 0" :key="0" />
      <AgeSlide v-else-if="currentSlide === 1" :key="1" v-model="formData.age" />
      <HeightSlide v-else-if="currentSlide === 2" :key="2" v-model="formData.heightCm" />
      <WeightSlide v-else-if="currentSlide === 3" :key="3" v-model="formData.weightKg" />
      <GenderSlide v-else-if="currentSlide === 4" :key="4" v-model="formData.gender" />
      <FitnessSlide v-else-if="currentSlide === 5" :key="5" v-model="formData.fitnessLevel" />
      <GoalsSlide v-else-if="currentSlide === 6" :key="6" v-model="formData.goals" />
      <SuccessSlide
        v-else-if="currentSlide === 7"
        :key="7"
        :age="formData.age"
        :heightCm="formData.heightCm"
        :weightKg="formData.weightKg"
      />
    </Transition>

    <div v-if="error" class="error-msg">{{ error }}</div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../api'
import { markRegistered } from '../router'
import { useTelegram } from '../composables/useTelegram'

import WelcomeSlide from './onboarding/WelcomeSlide.vue'
import AgeSlide from './onboarding/AgeSlide.vue'
import HeightSlide from './onboarding/HeightSlide.vue'
import WeightSlide from './onboarding/WeightSlide.vue'
import GenderSlide from './onboarding/GenderSlide.vue'
import FitnessSlide from './onboarding/FitnessSlide.vue'
import GoalsSlide from './onboarding/GoalsSlide.vue'
import SuccessSlide from './onboarding/SuccessSlide.vue'

const router = useRouter()
const { hapticImpact, hapticNotification, showMainButton, hideMainButton, showBackButton, hideBackButton } = useTelegram()

const currentSlide = ref(0)
const slideDirection = ref('slide-left')
const error = ref('')
const submitting = ref(false)

const formData = reactive({
  age: 25,
  heightCm: 170,
  weightKg: 70,
  gender: '',
  fitnessLevel: '',
  goals: [] as string[],
})

const mainButtonLabels = [
  'Начать',     // 0: welcome
  'Далее',      // 1: age
  'Далее',      // 2: height
  'Далее',      // 3: weight
  'Далее',      // 4: gender
  'Далее',      // 5: fitness
  'Готово ✨',  // 6: goals
  'К модулям 💪', // 7: success
]

function canProceed(): boolean {
  switch (currentSlide.value) {
    case 4: return formData.gender !== ''
    case 5: return formData.fitnessLevel !== ''
    case 6: return true
    default: return true
  }
}

function next() {
  if (!canProceed()) {
    hapticNotification('error')
    if (currentSlide.value === 4) error.value = 'Выбери пол'
    else if (currentSlide.value === 5) error.value = 'Выбери уровень'
    return
  }
  error.value = ''
  hapticImpact('light')

  if (currentSlide.value === 6) {
    submitRegistration()
    return
  }

  if (currentSlide.value === 7) {
    router.replace('/')
    return
  }

  slideDirection.value = 'slide-left'
  currentSlide.value++
}

function back() {
  if (currentSlide.value > 0 && currentSlide.value < 7) {
    error.value = ''
    slideDirection.value = 'slide-right'
    currentSlide.value--
    hapticImpact('light')
  }
}

async function submitRegistration() {
  if (submitting.value) return
  submitting.value = true
  error.value = ''

  try {
    await api.register({
      age: formData.age,
      height_cm: formData.heightCm,
      weight_kg: formData.weightKg,
      gender: formData.gender as 'male' | 'female',
      fitness_level: formData.fitnessLevel as 'beginner' | 'intermediate' | 'advanced',
      goals: formData.goals,
    })
    markRegistered()
    slideDirection.value = 'slide-left'
    currentSlide.value = 7
  } catch (e: any) {
    error.value = e.message || 'Ошибка регистрации'
    hapticNotification('error')
  } finally {
    submitting.value = false
  }
}

// MainButton and BackButton management
const nextRef = next
const backRef = back

function updateButtons() {
  showMainButton(mainButtonLabels[currentSlide.value], nextRef)
  if (currentSlide.value > 0 && currentSlide.value < 7) {
    showBackButton(backRef)
  } else {
    hideBackButton()
  }
}

watch(currentSlide, updateButtons)

onMounted(() => {
  updateButtons()
})

onUnmounted(() => {
  hideMainButton()
  hideBackButton()
})
</script>

<style scoped>
.onboarding {
  min-height: 80vh;
  display: flex;
  flex-direction: column;
  padding: 16px;
}

.progress-dots {
  display: flex;
  justify-content: center;
  gap: 6px;
  margin-bottom: 20px;
}

.dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--secondary-bg);
  transition: all 0.3s ease;
}

.dot.active {
  background: var(--button-color);
}

.dot.current {
  width: 24px;
  border-radius: 4px;
}

.error-msg {
  text-align: center;
  color: #ff3b30;
  font-size: 14px;
  margin-top: 12px;
}

/* Slide transitions */
.slide-left-enter-active,
.slide-left-leave-active,
.slide-right-enter-active,
.slide-right-leave-active {
  transition: all 0.3s ease;
}

.slide-left-enter-from { transform: translateX(40px); opacity: 0; }
.slide-left-leave-to { transform: translateX(-40px); opacity: 0; }
.slide-right-enter-from { transform: translateX(-40px); opacity: 0; }
.slide-right-leave-to { transform: translateX(40px); opacity: 0; }
</style>
