<template>
  <div class="onboarding">
    <!-- Progress dots (slides 1-3, hidden on welcome=0 and success=4) -->
    <div class="progress-dots" v-if="currentSlide > 0 && currentSlide < 4">
      <span
        v-for="i in 3"
        :key="i"
        class="dot"
        :class="{ active: i <= currentSlide, current: i === currentSlide }"
      ></span>
    </div>

    <!-- Slides -->
    <Transition :name="slideDirection" mode="out-in">
      <WelcomeSlide v-if="currentSlide === 0" :key="0" />
      <AboutYouSlide
        v-else-if="currentSlide === 1"
        :key="1"
        :gender="formData.gender"
        :age="formData.age"
        @update:gender="formData.gender = $event"
        @update:age="formData.age = Number($event)"
      />
      <HealthSlide
        v-else-if="currentSlide === 2"
        :key="2"
        :painLocations="formData.painLocations"
        :painLevel="formData.painLevel"
        :diagnoses="formData.diagnoses"
        @update:painLocations="formData.painLocations = $event"
        @update:painLevel="formData.painLevel = $event"
        @update:diagnoses="formData.diagnoses = $event"
      />
      <BodyFitnessSlide
        v-else-if="currentSlide === 3"
        :key="3"
        :heightCm="formData.heightCm"
        :weightKg="formData.weightKg"
        :fitnessLevel="formData.fitnessLevel"
        @update:heightCm="formData.heightCm = Number($event)"
        @update:weightKg="formData.weightKg = Number($event)"
        @update:fitnessLevel="formData.fitnessLevel = $event"
      />
      <SuccessSlide
        v-else-if="currentSlide === 4"
        :key="4"
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
import AboutYouSlide from './onboarding/AboutYouSlide.vue'
import HealthSlide from './onboarding/HealthSlide.vue'
import BodyFitnessSlide from './onboarding/BodyFitnessSlide.vue'
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
  fitnessLevel: 'beginner',
  // Health
  painLocations: [] as string[],
  painLevel: 0,
  diagnoses: [] as string[],
  // Defaults for removed fields
  trainingAccess: 'home',
  daysPerWeek: 3,
  sessionDuration: 35,
  preferredTime: 'evening',
  equipment: [] as string[],
})

const mainButtonLabels = [
  'Начать',            // 0: welcome
  'Далее',             // 1: about you
  'Далее',             // 2: health
  'Готово',            // 3: body & fitness
  'К ЛФК-курсам 💪',  // 4: success
]

function canProceed(): boolean {
  switch (currentSlide.value) {
    case 1: return formData.gender !== ''
    case 3: return formData.fitnessLevel !== ''
    default: return true
  }
}

function next() {
  if (!canProceed()) {
    hapticNotification('error')
    if (currentSlide.value === 1) error.value = 'Выбери пол'
    else if (currentSlide.value === 3) error.value = 'Выбери уровень'
    return
  }
  error.value = ''
  hapticImpact('light')

  if (currentSlide.value === 3) {
    submitRegistration()
    return
  }

  if (currentSlide.value === 4) {
    router.replace('/')
    return
  }

  slideDirection.value = 'slide-left'
  currentSlide.value++
}

function back() {
  if (currentSlide.value > 0 && currentSlide.value < 4) {
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

  // Auto-derive goals from diagnoses
  const validGoals = ['hernia', 'protrusion', 'scoliosis', 'kyphosis', 'lordosis']
  const autoGoals = formData.diagnoses.filter(d => validGoals.includes(d))

  try {
    await api.register({
      age: formData.age,
      height_cm: formData.heightCm,
      weight_kg: formData.weightKg,
      gender: formData.gender as 'male' | 'female',
      fitness_level: formData.fitnessLevel as 'beginner' | 'intermediate' | 'advanced',
      goals: autoGoals,
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
    currentSlide.value = 4
  } catch (e: any) {
    error.value = e.message
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
  if (currentSlide.value > 0 && currentSlide.value < 4) {
    showBackButton(backRef)
  } else {
    hideBackButton()
  }
}

watch(currentSlide, updateButtons)

onMounted(async () => {
  // Pre-fill form with existing profile data for returning users
  try {
    const profile = await api.getProfile()
    if (profile) {
      formData.age = profile.age || 25
      formData.heightCm = profile.height_cm || 170
      formData.weightKg = profile.weight_kg || 70
      formData.gender = profile.gender || ''
      formData.fitnessLevel = profile.fitness_level || 'beginner'
      formData.painLocations = profile.pain_locations || []
      formData.painLevel = profile.pain_level || 0
      formData.diagnoses = profile.diagnoses || []
      formData.trainingAccess = profile.training_access || 'home'
      formData.daysPerWeek = profile.days_per_week || 3
      formData.sessionDuration = profile.session_duration || 35
      formData.preferredTime = profile.preferred_time || 'evening'
    }
  } catch {
    // No profile yet — use defaults
  }
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
