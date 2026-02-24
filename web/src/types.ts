export interface Module {
  id: number
  slug: string
  name: string
  description: string
  icon: string
  requires_subscription: boolean
  is_active: boolean
  sort_order: number
}

export interface CategoryWithProgress {
  id: number
  module_id: number
  slug: string
  name: string
  description: string
  icon: string
  sort_order: number
  completed: number
  total: number
}

export interface LessonWithProgress {
  id: number
  category_id: number
  slug: string
  title: string
  description: string
  sort_order: number
  status: 'not_started' | 'in_progress' | 'completed'
}

export interface LessonContent {
  id: number
  lesson_id: number
  content_type: string
  title: string
  body: string
  video_url: string
  telegram_file_id: string
  file_url: string
  sort_order: number
}

export interface LessonDetail {
  id: number
  category_id: number
  title: string
  description: string
  status: string
  content: LessonContent[]
}

export interface SubscriptionStatus {
  active: boolean
}

export interface RegisterRequest {
  age: number
  height_cm: number
  weight_kg: number
  gender: 'male' | 'female'
  fitness_level: 'beginner' | 'intermediate' | 'advanced'
  goals: string[]
}

export interface RegistrationStatus {
  is_registered: boolean
}

export interface PaymentStatus {
  is_paid: boolean
  price_kzt: number
}

export interface PaymentResult {
  success: boolean
  message: string
}

export interface UserProfile {
  first_name: string
  last_name: string
  username: string
  age: number
  height_cm: number
  weight_kg: number
  gender: string
  fitness_level: string
  goals: string[]
  is_paid: boolean
}
