import type {
  Module,
  CategoryWithProgress,
  LessonWithProgress,
  LessonDetail,
  SubscriptionStatus,
  RegisterRequest,
  RegistrationStatus,
  PaymentStatus,
  PaymentResult,
  UserProfile,
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

  getPaymentStatus(): Promise<PaymentStatus> {
    return request('/app/api/payment/status')
  },

  processPayment(): Promise<PaymentResult> {
    return request('/app/api/payment/pay', { method: 'POST' })
  },

  getProfile(): Promise<UserProfile> {
    return request('/app/api/profile')
  },

  register(data: RegisterRequest): Promise<{ success: boolean }> {
    return request('/app/api/register', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  },

  getRegistrationStatus(): Promise<RegistrationStatus> {
    return request('/app/api/registration/status')
  },
}
