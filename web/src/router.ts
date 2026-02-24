import { createRouter, createWebHistory } from 'vue-router'
import ModulesView from './views/ModulesView.vue'
import CategoriesView from './views/CategoriesView.vue'
import LessonsView from './views/LessonsView.vue'
import LessonView from './views/LessonView.vue'
import PaymentView from './views/PaymentView.vue'
import ProfileView from './views/ProfileView.vue'
import OnboardingView from './views/OnboardingView.vue'
import { api } from './api'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/onboarding',
      name: 'onboarding',
      component: OnboardingView,
      meta: { hideNav: true },
    },
    {
      path: '/',
      name: 'modules',
      component: ModulesView,
    },
    {
      path: '/module/:id',
      name: 'categories',
      component: CategoriesView,
      props: true,
    },
    {
      path: '/category/:id',
      name: 'lessons',
      component: LessonsView,
      props: true,
    },
    {
      path: '/lesson/:id',
      name: 'lesson',
      component: LessonView,
      props: true,
    },
    {
      path: '/payment',
      name: 'payment',
      component: PaymentView,
    },
    {
      path: '/profile',
      name: 'profile',
      component: ProfileView,
    },
  ],
})

let registrationChecked = false
let isRegistered = false

router.beforeEach(async (to) => {
  if (to.name === 'onboarding') return true

  if (!registrationChecked) {
    try {
      const status = await api.getRegistrationStatus()
      isRegistered = status.is_registered
    } catch {
      isRegistered = false
    }
    registrationChecked = true
  }

  if (!isRegistered) {
    return { name: 'onboarding' }
  }

  return true
})

export function markRegistered() {
  isRegistered = true
  registrationChecked = true
}

export default router
