<template>
  <div class="onboarding">
    <!-- Progress dots -->
    <div class="progress-dots" v-if="currentSlide > 0 && currentSlide < totalSlides">
      <span
        v-for="i in totalSlides"
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
      <HealthSlide
        v-else-if="currentSlide === 7"
        :key="7"
        :painLocations="formData.painLocations"
        :painLevel="formData.painLevel"
        :diagnoses="formData.diagnoses"
        @update:painLocations="formData.painLocations = $event"
        @update:painLevel="formData.painLevel = $event"
        @update:diagnoses="formData.diagnoses = $event"
      />
      <ScheduleSlide
        v-else-if="currentSlide === 8"
        :key="8"
        :trainingAccess="formData.trainingAccess"
        :daysPerWeek="formData.daysPerWeek"
        :sessionDuration="formData.sessionDuration"
        :preferredTime="formData.preferredTime"
        :equipment="formData.equipment"
        @update:trainingAccess="formData.trainingAccess = $event"
        @update:daysPerWeek="formData.daysPerWeek = $event"
        @update:sessionDuration="formData.sessionDuration = $event"
        @update:preferredTime="formData.preferredTime = $event"
        @update:equipment="formData.equipment = $event"
      />
      <SuccessSlide
        v-else-if="currentSlide === 9"
        :key="9"
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
import HealthSlide from './onboarding/HealthSlide.vue'
import ScheduleSlide from './onboarding/ScheduleSlide.vue'
import SuccessSlide from './onboarding/SuccessSlide.vue'

const router = useRouter()
const { hapticImpact, hapticNotification, showMainButton, hideMainButton, showBackButton, hideBackButton } = useTelegram()

const totalSlides = 9 // slides 1-8 (excluding welcome=0 and success=9)
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
  // Health
  painLocations: [] as string[],
  painLevel: 0,
  diagnoses: [] as string[],
  // Schedule
  trainingAccess: 'home',
  daysPerWeek: 3,
  sessionDuration: 35,
  preferredTime: 'evening',
  equipment: [] as string[],
})

const mainButtonLabels = [
  'Начать',        // 0: welcome
  'Далее',         // 1: age
  'Далее',         // 2: height
  'Далее',         // 3: weight
  'Далее',         // 4: gender
  'Далее',         // 5: fitness
  'Далее',         // 6: goals
  'Далее',         // 7: health
  'Готово ✨',     // 8: schedule
  'К тренировкам 💪', // 9: success
]

function canProceed(): boolean {
  switch (currentSlide.value) {
    case 4: return formData.gender !== ''
    case 5: return formData.fitnessLevel !== ''
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

  if (currentSlide.value === 8) {
    submitRegistration()
    return
  }

  if (currentSlide.value === 9) {
    router.replace('/')
    return
  }

  slideDirection.value = 'slide-left'
  currentSlide.value++
}

function back() {
  if (currentSlide.value > 0 && currentSlide.value < 9) {
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
      training_access: formData.trainingAccess,
      has_pain: formData.painLocations.length > 0,
      pain_locations: formData.painLocations,
      pain_level: formData.painLevel,
      diagnoses: formData.diagnoses,
      days_per_week: formData.daysPerWeek,
      session_duration: formData.sessionDuration,
      preferred_time: formData.preferredTime,
      equipment: formData.equipment,
    })
    markRegistered()
    slideDirection.value = 'slide-left'
    currentSlide.value = 9
  } catch (e: any) {
    const tg = window.Telegram?.WebApp
    const initLen = tg?.initData?.length ?? -1
    const platform = tg?.platform ?? 'none'
    const version = tg?.version ?? 'none'
    const hash = window.location.hash?.length ?? 0
    const href = window.location.href
    const uid = (tg as any)?.initDataUnsafe?.user?.id ?? 'no-user'
    error.value = `${e.message} [init:${initLen}, p:${platform}, v:${version}, hash:${hash}, uid:${uid}, url:${href}]`
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
  if (currentSlide.value > 0 && currentSlide.value < 9) {
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
