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
    // { path: '/nutrition/calculator', name: 'macro-calculator', component: () => import('./views/MacroCalculatorView.vue'), meta: { hideNav: true } },
    { path: '/progress', name: 'progress', component: () => import('./views/ProgressView.vue') },
    { path: '/profile', name: 'profile', component: () => import('./views/ProfileView.vue'), meta: { hideNav: true } },
    { path: '/payment', name: 'payment', component: () => import('./views/PaymentView.vue'), meta: { hideNav: true } },
    // Legacy module routes
    { path: '/modules', name: 'modules', component: () => import('./views/ModulesView.vue'), meta: { hideNav: true } },
    { path: '/module/:id', name: 'categories', component: () => import('./views/CategoriesView.vue'), props: true, meta: { hideNav: true } },
    { path: '/category/:id', name: 'lessons', component: () => import('./views/LessonsView.vue'), props: true, meta: { hideNav: true } },
    { path: '/lesson/:id', name: 'lesson', component: () => import('./views/LessonView.vue'), props: true, meta: { hideNav: true } },
    { path: '/review', name: 'review', component: () => import('./views/ReviewView.vue'), meta: { hideNav: true } },
    // Admin routes
    { path: '/admin', name: 'admin', component: () => import('./views/admin/AdminDashboardView.vue'), meta: { hideNav: true, requiresAdmin: true } },
    { path: '/admin/users', name: 'admin-users', component: () => import('./views/admin/AdminUsersView.vue'), meta: { hideNav: true, requiresAdmin: true } },
    { path: '/admin/users/:id', name: 'admin-user-detail', component: () => import('./views/admin/AdminUserDetailView.vue'), props: true, meta: { hideNav: true, requiresAdmin: true } },
    { path: '/admin/content', name: 'admin-content', component: () => import('./views/admin/AdminContentView.vue'), meta: { hideNav: true, requiresAdmin: true } },
    { path: '/admin/programs/new', name: 'admin-program-new', component: () => import('./views/admin/AdminProgramFormView.vue'), meta: { hideNav: true, requiresAdmin: true } },
    { path: '/admin/programs/:id', name: 'admin-program-edit', component: () => import('./views/admin/AdminProgramFormView.vue'), props: true, meta: { hideNav: true, requiresAdmin: true } },
    { path: '/admin/workouts/new', name: 'admin-workout-new', component: () => import('./views/admin/AdminWorkoutFormView.vue'), meta: { hideNav: true, requiresAdmin: true } },
    { path: '/admin/workouts/:id', name: 'admin-workout-edit', component: () => import('./views/admin/AdminWorkoutFormView.vue'), props: true, meta: { hideNav: true, requiresAdmin: true } },
    { path: '/admin/exercises', name: 'admin-exercises', component: () => import('./views/admin/AdminExercisesView.vue'), meta: { hideNav: true, requiresAdmin: true } },
    { path: '/admin/exercises/new', name: 'admin-exercise-new', component: () => import('./views/admin/AdminExerciseFormView.vue'), meta: { hideNav: true, requiresAdmin: true } },
    { path: '/admin/exercises/:id', name: 'admin-exercise-edit', component: () => import('./views/admin/AdminExerciseFormView.vue'), props: true, meta: { hideNav: true, requiresAdmin: true } },
    { path: '/admin/meal-plans/new', name: 'admin-mealplan-new', component: () => import('./views/admin/AdminMealPlanFormView.vue'), meta: { hideNav: true, requiresAdmin: true } },
    { path: '/admin/meal-plans/:id', name: 'admin-mealplan-edit', component: () => import('./views/admin/AdminMealPlanFormView.vue'), props: true, meta: { hideNav: true, requiresAdmin: true } },
    { path: '/admin/meals/new', name: 'admin-meal-new', component: () => import('./views/admin/AdminMealFormView.vue'), meta: { hideNav: true, requiresAdmin: true } },
    { path: '/admin/meals/:id', name: 'admin-meal-edit', component: () => import('./views/admin/AdminMealFormView.vue'), props: true, meta: { hideNav: true, requiresAdmin: true } },
    { path: '/admin/reviews', name: 'admin-reviews', component: () => import('./views/admin/AdminReviewsView.vue'), meta: { hideNav: true, requiresAdmin: true } },
  ],
})

let onboardingDone = false
let authDone = false

router.beforeEach(async (to) => {
  // Authenticate once on app start — get session token
  if (!authDone) {
    try { await api.authenticate() } catch { /* continue without token */ }
    authDone = true
  }

  // Admin route guard
  if (to.meta.requiresAdmin && !api.isAdmin()) {
    return { name: 'today' }
  }

  // Always show onboarding first when app opens
  if (!onboardingDone && to.name !== 'onboarding') {
    return { name: 'onboarding' }
  }
  if (onboardingDone && to.name === 'onboarding') {
    return { name: 'today' }
  }
  return true
})

export function markRegistered() {
  onboardingDone = true
}

export default router
