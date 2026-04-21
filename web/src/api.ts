import type {
  Module,
  CategoryWithProgress,
  LessonWithProgress,
  LessonDetail,
  SubscriptionStatus,
  RegisterRequest,
  RegisterResponse,
  RegistrationStatus,
  PaymentStatus,
  PaymentResult,
  UserProfile,
  DashboardData,
  Program,
  Workout,
  WorkoutWithExercises,
  WorkoutExercise,
  Exercise,
  WorkoutExerciseLink,
  RehabCourse,
  RehabCourseWithSessions,
  RehabSession,
  RehabProgress,
  CompleteRehabRequest,
  MealPlan,
  Meal,
  MealPlanWithMeals,
  MacroTargets,
  FoodLogEntry,
  CreateFoodLogRequest,
  DailySummary,
  ProgressEntry,
  CreateProgressRequest,
  ProgressStats,
  WeightPoint,
  Achievement,
  UserAchievement,
  Review,
  CreateReviewRequest,
  ReviewSummary,
  ReviewTagsResponse,
  AdminUser,
} from './types'

function getInitData(): string {
  return window.Telegram?.WebApp?.initData || ''
}

let authToken: string | null = null
let userRole: string | null = null

function getAuthHeaders(): Record<string, string> {
  if (authToken) {
    return { Authorization: `Bearer ${authToken}` }
  }
  return { 'X-Telegram-Init-Data': getInitData() }
}

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(path, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...getAuthHeaders(),
      ...options?.headers,
    },
  })

  // On 401 with a token, clear it and retry with initData
  if (res.status === 401 && authToken) {
    authToken = null
    return request(path, options)
  }

  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: 'Request failed' }))
    throw new Error(err.error || `HTTP ${res.status}`)
  }

  return res.json()
}

async function authenticate(): Promise<void> {
  const initData = getInitData()
  if (!initData) return

  try {
    const res = await fetch('/app/api/auth', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-Telegram-Init-Data': initData,
      },
    })
    if (res.ok) {
      const data = await res.json()
      authToken = data.token
      userRole = data.role || 'client'
    }
  } catch {
    // Auth failed silently — requests will fall back to initData
  }
}

export const api = {
  authenticate,
  // ====== LEGACY MODULES ======
  getModules(): Promise<Module[]> {
    return request('/app/api/modules')
  },
  getCategories(moduleId: number): Promise<CategoryWithProgress[]> {
    return request(`/app/api/modules/${moduleId}/categories`)
  },
  getLessons(categoryId: number): Promise<LessonWithProgress[]> {
    return request(`/app/api/categories/${categoryId}/lessons`)
  },
  getLesson(lessonId: number): Promise<LessonDetail> {
    return request(`/app/api/lessons/${lessonId}`)
  },
  startLesson(lessonId: number): Promise<{ status: string }> {
    return request(`/app/api/lessons/${lessonId}/start`, { method: 'POST' })
  },
  completeLesson(lessonId: number): Promise<{ status: string }> {
    return request(`/app/api/lessons/${lessonId}/complete`, { method: 'POST' })
  },
  getSubscriptionStatus(): Promise<SubscriptionStatus> {
    return request('/app/api/subscription/status')
  },

  // ====== AUTH & REGISTRATION ======
  getRegistrationStatus(): Promise<RegistrationStatus> {
    return request('/app/api/registration/status')
  },
  register(data: RegisterRequest): Promise<RegisterResponse> {
    return request('/app/api/register', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  },
  getProfile(): Promise<UserProfile> {
    return request('/app/api/profile')
  },
  updateProfile(data: Partial<RegisterRequest>): Promise<{ success: boolean }> {
    return request('/app/api/profile', {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  },

  // ====== PAYMENT ======
  getPaymentStatus(): Promise<PaymentStatus> {
    return request('/app/api/payment/status')
  },
  processPayment(): Promise<PaymentResult> {
    return request('/app/api/payment/pay', { method: 'POST' })
  },

  // ====== DASHBOARD ======
  getDashboard(): Promise<DashboardData> {
    return request('/app/api/dashboard')
  },

  // ====== PROGRAMS ======
  getPrograms(filters?: { format?: string; goal?: string; level?: string }): Promise<Program[]> {
    const params = new URLSearchParams()
    if (filters?.format) params.set('format', filters.format)
    if (filters?.goal) params.set('goal', filters.goal)
    if (filters?.level) params.set('level', filters.level)
    const qs = params.toString()
    return request(`/app/api/programs${qs ? '?' + qs : ''}`)
  },
  getProgram(id: number): Promise<Program> {
    return request(`/app/api/programs/${id}`)
  },
  enrollProgram(id: number): Promise<{ success: boolean }> {
    return request(`/app/api/programs/${id}/enroll`, { method: 'POST' })
  },

  // ====== WORKOUTS ======
  getWorkouts(filters?: { format?: string; goal?: string; level?: string }): Promise<Workout[]> {
    const params = new URLSearchParams()
    if (filters?.format) params.set('format', filters.format)
    if (filters?.goal) params.set('goal', filters.goal)
    if (filters?.level) params.set('level', filters.level)
    const qs = params.toString()
    return request(`/app/api/workouts${qs ? '?' + qs : ''}`)
  },
  getWorkout(id: number): Promise<WorkoutWithExercises> {
    return request(`/app/api/workouts/${id}`)
  },
  completeWorkout(id: number): Promise<{ success: boolean }> {
    return request(`/app/api/workouts/${id}/complete`, { method: 'POST' })
  },

  // ====== REHAB / LFK ======
  getRehabCourses(category?: string): Promise<RehabCourse[]> {
    const qs = category ? `?category=${category}` : ''
    return request(`/app/api/rehab/courses${qs}`)
  },
  getRehabCourse(id: number): Promise<RehabCourseWithSessions> {
    return request(`/app/api/rehab/courses/${id}`)
  },
  getRehabSession(id: number): Promise<RehabSession> {
    return request(`/app/api/rehab/sessions/${id}`)
  },
  completeRehabSession(sessionId: number, data: CompleteRehabRequest): Promise<{ success: boolean }> {
    return request(`/app/api/rehab/sessions/${sessionId}/complete`, {
      method: 'POST',
      body: JSON.stringify(data),
    })
  },
  getRehabProgress(courseId: number): Promise<RehabProgress[]> {
    return request(`/app/api/rehab/progress/${courseId}`)
  },

  // ====== NUTRITION ======
  getMealPlans(goal?: string): Promise<MealPlan[]> {
    const qs = goal ? `?goal=${goal}` : ''
    return request(`/app/api/nutrition/plans${qs}`)
  },
  getMealPlan(id: number): Promise<MealPlanWithMeals> {
    return request(`/app/api/nutrition/plans/${id}`)
  },
  calculateMacros(params: {
    gender: string
    weight_kg: number
    height_cm: number
    age: number
    goal: string
  }): Promise<MacroTargets> {
    const qs = new URLSearchParams(
      Object.entries(params).reduce((acc, [k, v]) => ({ ...acc, [k]: String(v) }), {} as Record<string, string>)
    ).toString()
    return request(`/app/api/nutrition/calculator?${qs}`)
  },

  // ====== FOOD LOG ======
  getFoodLog(date: string): Promise<FoodLogEntry[]> {
    return request(`/app/api/food-log?date=${date}`)
  },
  addFoodLog(data: CreateFoodLogRequest): Promise<{ success: boolean; id: number }> {
    return request('/app/api/food-log', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  },
  deleteFoodLog(id: number): Promise<{ success: boolean }> {
    return request(`/app/api/food-log/${id}`, { method: 'DELETE' })
  },
  getFoodLogSummary(date: string): Promise<DailySummary> {
    return request(`/app/api/food-log/summary?date=${date}`)
  },

  // ====== PROGRESS ======
  getProgressEntries(): Promise<ProgressEntry[]> {
    return request('/app/api/progress')
  },
  addProgressEntry(data: CreateProgressRequest): Promise<{ success: boolean; id: number }> {
    return request('/app/api/progress', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  },
  getProgressStats(): Promise<ProgressStats> {
    return request('/app/api/progress/stats')
  },
  getWeightHistory(): Promise<WeightPoint[]> {
    return request('/app/api/progress')
  },

  // ====== ACHIEVEMENTS ======
  getAchievements(): Promise<{ achievements: Achievement[]; user_achievements: UserAchievement[] }> {
    return request('/app/api/progress/achievements')
  },

  // ====== REVIEWS ======
  createReview(data: CreateReviewRequest): Promise<Review> {
    return request('/app/api/reviews', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  },
  getReviews(referenceType: string, referenceId: number): Promise<Review[]> {
    return request(`/app/api/reviews?reference_type=${referenceType}&reference_id=${referenceId}`)
  },
  getMyReviews(): Promise<Review[]> {
    return request('/app/api/reviews/my')
  },
  getReviewSummary(referenceType: string, referenceId: number): Promise<ReviewSummary> {
    return request(`/app/api/reviews/summary?reference_type=${referenceType}&reference_id=${referenceId}`)
  },
  getBotReviewSummary(): Promise<ReviewSummary> {
    return request('/app/api/reviews/bot/summary')
  },
  getReviewTags(referenceType: string): Promise<ReviewTagsResponse> {
    return request(`/app/api/reviews/tags?reference_type=${referenceType}`)
  },

  // ====== ADMIN ======
  isAdmin(): boolean {
    return userRole === 'admin'
  },
  getAdminUsers(limit = 20, offset = 0): Promise<{ users: AdminUser[]; total: number }> {
    return request(`/app/api/admin/users?limit=${limit}&offset=${offset}`)
  },
  getAdminUser(id: number): Promise<{ user: AdminUser; profile: UserProfile | null }> {
    return request(`/app/api/admin/users/${id}`)
  },
  updateAdminUser(id: number, data: { role?: string; is_paid?: boolean }): Promise<{ success: boolean }> {
    return request(`/app/api/admin/users/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  },
  // Admin Programs CRUD
  getAdminPrograms(): Promise<Program[]> {
    return request('/app/api/admin/programs')
  },
  getAdminProgram(id: number): Promise<Program> {
    return request(`/app/api/admin/programs/${id}`)
  },
  createAdminProgram(data: Partial<Program>): Promise<Program> {
    return request('/app/api/admin/programs', { method: 'POST', body: JSON.stringify(data) })
  },
  updateAdminProgram(id: number, data: Partial<Program>): Promise<Program> {
    return request(`/app/api/admin/programs/${id}`, { method: 'PUT', body: JSON.stringify(data) })
  },

  // Admin Workouts CRUD
  getAdminWorkouts(): Promise<Workout[]> {
    return request('/app/api/admin/workouts')
  },
  getAdminWorkout(id: number): Promise<{ workout: Workout; exercises: WorkoutExercise[] }> {
    return request(`/app/api/admin/workouts/${id}`)
  },
  createAdminWorkout(data: Partial<Workout>): Promise<Workout> {
    return request('/app/api/admin/workouts', { method: 'POST', body: JSON.stringify(data) })
  },
  updateAdminWorkout(id: number, data: Partial<Workout>): Promise<Workout> {
    return request(`/app/api/admin/workouts/${id}`, { method: 'PUT', body: JSON.stringify(data) })
  },

  // Admin Exercises CRUD
  getAdminExercises(): Promise<Exercise[]> {
    return request('/app/api/admin/exercises')
  },
  getAdminExercise(id: number): Promise<Exercise> {
    return request(`/app/api/admin/exercises/${id}`)
  },
  createAdminExercise(data: Partial<Exercise>): Promise<Exercise> {
    return request('/app/api/admin/exercises', { method: 'POST', body: JSON.stringify(data) })
  },
  updateAdminExercise(id: number, data: Partial<Exercise>): Promise<Exercise> {
    return request(`/app/api/admin/exercises/${id}`, { method: 'PUT', body: JSON.stringify(data) })
  },

  // Admin Workout Exercises
  addWorkoutExercise(data: WorkoutExerciseLink): Promise<WorkoutExerciseLink> {
    return request('/app/api/admin/workout-exercises', { method: 'POST', body: JSON.stringify(data) })
  },

  // Admin Meal Plans CRUD
  getAdminMealPlans(): Promise<MealPlan[]> {
    return request('/app/api/admin/meal-plans')
  },
  getAdminMealPlan(id: number): Promise<{ plan: MealPlan; meals: Meal[] }> {
    return request(`/app/api/admin/meal-plans/${id}`)
  },
  createAdminMealPlan(data: Partial<MealPlan>): Promise<MealPlan> {
    return request('/app/api/admin/meal-plans', { method: 'POST', body: JSON.stringify(data) })
  },
  updateAdminMealPlan(id: number, data: Partial<MealPlan>): Promise<MealPlan> {
    return request(`/app/api/admin/meal-plans/${id}`, { method: 'PUT', body: JSON.stringify(data) })
  },

  // Admin Meals CRUD
  getAdminMeals(planId: number): Promise<Meal[]> {
    return request(`/app/api/admin/meals?plan_id=${planId}`)
  },
  getAdminMeal(id: number): Promise<Meal> {
    return request(`/app/api/admin/meals/${id}`)
  },
  createAdminMeal(data: Partial<Meal>): Promise<Meal> {
    return request('/app/api/admin/meals', { method: 'POST', body: JSON.stringify(data) })
  },
  updateAdminMeal(id: number, data: Partial<Meal>): Promise<Meal> {
    return request(`/app/api/admin/meals/${id}`, { method: 'PUT', body: JSON.stringify(data) })
  },

  // Admin Reviews & Stats
  getAdminReviewsSummary(): Promise<ReviewSummary> {
    return request('/app/api/admin/reviews')
  },
  getAdminStats(): Promise<{ total_users: number }> {
    return request('/app/api/admin/stats')
  },
}
