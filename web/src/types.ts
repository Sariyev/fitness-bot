// ====== EXISTING TYPES (LEGACY MODULES) ======

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

// ====== USER & REGISTRATION ======

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
  training_access?: string
  training_experience?: string
  has_pain?: boolean
  pain_locations?: string[]
  pain_level?: number
  diagnoses?: string[]
  contraindications?: string
  days_per_week?: number
  session_duration?: number
  preferred_time?: string
  equipment?: string[]
}

export interface RegisterRequest {
  age: number
  height_cm: number
  weight_kg: number
  gender: 'male' | 'female'
  fitness_level: 'beginner' | 'intermediate' | 'advanced'
  goals: string[]
  training_access?: string
  training_experience?: string
  has_pain?: boolean
  pain_locations?: string[]
  pain_level?: number
  diagnoses?: string[]
  contraindications?: string
  days_per_week?: number
  session_duration?: number
  preferred_time?: string
  equipment?: string[]
}

export interface RegisterResponse {
  success: boolean
  recommendations?: Recommendation
}

export interface Recommendation {
  program_id?: number
  program_name?: string
  rehab_course_id?: number
  rehab_course_name?: string
  calories: number
  protein: number
  fat: number
  carbs: number
  message: string
}

// ====== DASHBOARD ======

export interface DashboardData {
  greeting: string
  trainer_message: string
  today_workout?: DashboardItem
  today_meal?: DashboardItem
  today_rehab?: DashboardItem
  goals: string[]
  current_streak: number
}

export interface DashboardItem {
  id: number
  title: string
  type: string
  done: boolean
}

// ====== WORKOUTS & PROGRAMS ======

export interface Program {
  id: number
  slug: string
  name: string
  description: string
  goal: string
  format: string
  level: string
  duration_weeks: number
  is_active: boolean
  sort_order: number
}

export interface Workout {
  id: number
  program_id?: number
  slug: string
  name: string
  description: string
  goal: string
  format: string
  level: string
  duration_minutes: number
  equipment: string[]
  expected_result: string
  video_url: string
  sort_order: number
  week_number?: number
  day_number?: number
  is_active: boolean
}

export interface WorkoutExercise {
  id: number
  workout_id: number
  exercise_id: number
  sets: number
  reps: string
  duration_seconds: number
  sort_order: number
}

export interface WorkoutWithExercises extends Workout {
  exercises: WorkoutExercise[]
}

// ====== REHAB / LFK ======

export interface RehabCourse {
  id: number
  slug: string
  category: string
  name: string
  description: string
  warnings: string
  is_active: boolean
  sort_order: number
}

export interface RehabSession {
  id: number
  course_id: number
  day_number: number
  stage: number
  video_url: string
  duration_minutes: number
  description: string
  sort_order: number
}

export interface RehabCourseWithSessions extends RehabCourse {
  sessions: RehabSession[]
}

export interface RehabProgress {
  id: number
  user_id: number
  course_id: number
  session_id: number
  day_number: number
  completed_at: string
  pain_level: number
  comment: string
}

export interface CompleteRehabRequest {
  pain_level: number
  comment: string
  day_number: number
  course_id: number
}

// ====== NUTRITION ======

export interface MealPlan {
  id: number
  slug: string
  name: string
  goal: string
  day_number: number
  calories: number
  protein: number
  fat: number
  carbs: number
  is_active: boolean
  sort_order: number
}

export interface Meal {
  id: number
  meal_plan_id: number
  meal_type: string
  name: string
  recipe: string
  calories: number
  protein: number
  fat: number
  carbs: number
  alternatives: string
  sort_order: number
}

export interface MealPlanWithMeals extends MealPlan {
  meals: Meal[]
}

export interface MacroTargets {
  calories: number
  protein: number
  fat: number
  carbs: number
}

export interface FoodLogEntry {
  id: number
  user_id: number
  date: string
  meal_type: string
  food_name: string
  calories: number
  protein: number
  fat: number
  carbs: number
  photo_url: string
  created_at: string
}

export interface CreateFoodLogRequest {
  date: string
  meal_type: string
  food_name: string
  calories: number
  protein: number
  fat: number
  carbs: number
}

export interface DailySummary {
  calories: number
  protein: number
  fat: number
  carbs: number
}

// ====== PROGRESS ======

export interface ProgressEntry {
  id: number
  user_id: number
  date: string
  weight_kg?: number
  measurements?: Record<string, number>
  photo_url: string
  wellbeing: string
  pain_level: number
  created_at: string
}

export interface CreateProgressRequest {
  date: string
  weight_kg?: number
  measurements?: Record<string, number>
  photo_url?: string
  wellbeing?: string
  pain_level?: number
}

export interface ProgressStats {
  current_streak: number
  longest_streak: number
  calendar: string[]
}

export interface WeightPoint {
  date: string
  weight_kg: number
}

// ====== ACHIEVEMENTS ======

export interface Achievement {
  id: number
  slug: string
  name: string
  description: string
  icon: string
  criteria: Record<string, unknown>
}

export interface UserAchievement {
  id: number
  user_id: number
  achievement_id: number
  earned_at: string
}
