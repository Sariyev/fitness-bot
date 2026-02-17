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
