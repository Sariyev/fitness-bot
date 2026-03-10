import { createRouter, createWebHistory } from 'vue-router'
import { api } from './api'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/onboarding',
      name: 'onboarding',
      component: () => import('./views/OnboardingView.vue'),
      meta: { hideNav: true },
    },
    { path: '/', name: 'today', component: () => import('./views/TodayView.vue') },
    { path: '/workouts', name: 'workouts', component: () => import('./views/WorkoutsView.vue') },
    { path: '/workouts/program/:id', name: 'workout-program', component: () => import('./views/WorkoutProgramView.vue'), props: true, meta: { hideNav: true } },
    { path: '/workouts/session/:id', name: 'workout-session', component: () => import('./views/WorkoutSessionView.vue'), props: true, meta: { hideNav: true } },
    { path: '/lfk', name: 'lfk', component: () => import('./views/LfkView.vue') },
    { path: '/lfk/course/:id', name: 'lfk-course', component: () => import('./views/LfkCourseView.vue'), props: true, meta: { hideNav: true } },
    { path: '/lfk/session/:id', name: 'lfk-session', component: () => import('./views/LfkSessionView.vue'), props: true, meta: { hideNav: true } },
    { path: '/nutrition', name: 'nutrition', component: () => import('./views/NutritionView.vue') },
    { path: '/nutrition/diary', name: 'food-diary', component: () => import('./views/FoodDiaryView.vue'), meta: { hideNav: true } },
    { path: '/nutrition/calculator', name: 'macro-calculator', component: () => import('./views/MacroCalculatorView.vue'), meta: { hideNav: true } },
    { path: '/progress', name: 'progress', component: () => import('./views/ProgressView.vue') },
    { path: '/profile', name: 'profile', component: () => import('./views/ProfileView.vue'), meta: { hideNav: true } },
    { path: '/payment', name: 'payment', component: () => import('./views/PaymentView.vue'), meta: { hideNav: true } },
    // Legacy module routes
    { path: '/modules', name: 'modules', component: () => import('./views/ModulesView.vue'), meta: { hideNav: true } },
    { path: '/module/:id', name: 'categories', component: () => import('./views/CategoriesView.vue'), props: true, meta: { hideNav: true } },
    { path: '/category/:id', name: 'lessons', component: () => import('./views/LessonsView.vue'), props: true, meta: { hideNav: true } },
    { path: '/lesson/:id', name: 'lesson', component: () => import('./views/LessonView.vue'), props: true, meta: { hideNav: true } },
  ],
})

let registrationChecked = false
let isRegistered = false

router.beforeEach(async (to) => {
  if (!registrationChecked) {
    try {
      const status = await api.getRegistrationStatus()
      isRegistered = status.is_registered
    } catch {
      isRegistered = false
    }
    registrationChecked = true
  }
  if (isRegistered && to.name === 'onboarding') return { name: 'today' }
  if (!isRegistered && to.name !== 'onboarding') return { name: 'onboarding' }
  return true
})

export function markRegistered() {
  isRegistered = true
  registrationChecked = true
}

export default router
