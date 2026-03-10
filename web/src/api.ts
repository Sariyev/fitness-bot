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
  RehabCourse,
  RehabCourseWithSessions,
  RehabSession,
  RehabProgress,
  CompleteRehabRequest,
  MealPlan,
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
} from './types'

function getInitData(): string {
  return window.Telegram?.WebApp?.initData || ''
}

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(path, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      'X-Telegram-Init-Data': getInitData(),
      ...options?.headers,
    },
  })

  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: 'Request failed' }))
    throw new Error(err.error || `HTTP ${res.status}`)
  }

  return res.json()
}

export const api = {
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
}
