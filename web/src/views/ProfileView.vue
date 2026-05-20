<template>
  <div class="profile-page">

    <div v-if="loading" class="skeleton-list">
      <SkeletonCard v-for="i in 4" :key="i" />
    </div>

    <div v-else-if="profile" class="profile-content">
      <div class="profile-header">
        <div class="header-gradient"></div>
        <div class="avatar avatar-button" @click="pickAvatar" :class="{ uploading: avatarUploading }">
          <div class="avatar-circle">
            <img v-if="avatarUrl" :src="avatarUrl" alt="avatar" class="avatar-img" />
            <span v-else>{{ initials }}</span>
            <div v-if="avatarUploading" class="avatar-overlay">
              <span>{{ Math.round(avatarProgress * 100) }}%</span>
            </div>
          </div>
          <div v-if="!avatarUploading" class="avatar-edit-badge" aria-label="Изменить фото">📷</div>
        </div>
        <input
          ref="avatarInput"
          type="file"
          accept="image/jpeg,image/png,image/webp"
          @change="handleAvatarChange"
          style="display: none"
        />
        <h2 class="profile-name">{{ profile.first_name }} {{ profile.last_name }}</h2>
        <p v-if="profile.username" class="username">@{{ profile.username }}</p>
        <p v-if="avatarError" class="avatar-error">{{ avatarError }}</p>

        <!-- Vitals row — quick-glance stats, only in read mode -->
        <div v-if="!editing && (profile.age || profile.height_cm || profile.weight_kg)" class="vitals-row">
          <div class="vital-tile" v-if="profile.age">
            <span class="vital-value">{{ profile.age }}</span>
            <span class="vital-label">Возраст</span>
          </div>
          <div class="vital-tile" v-if="profile.height_cm">
            <span class="vital-value">{{ profile.height_cm }}</span>
            <span class="vital-label">Рост, см</span>
          </div>
          <div class="vital-tile" v-if="profile.weight_kg">
            <span class="vital-value">{{ profile.weight_kg }}</span>
            <span class="vital-label">Вес, кг</span>
          </div>
        </div>

        <button class="edit-toggle" @click="toggleEdit">
          {{ editing ? 'Отмена' : 'Редактировать' }}
        </button>
      </div>

      <!-- Basic info -->
      <div class="info-card" style="animation-delay: 80ms">
        <div class="info-row">
          <span class="label">Возраст</span>
          <input v-if="editing" v-model.number="editData.age" type="number" class="edit-input" />
          <span v-else class="value">{{ profile.age }}</span>
        </div>
        <div class="info-row">
          <span class="label">Рост</span>
          <div v-if="editing" class="edit-with-unit">
            <input v-model.number="editData.height_cm" type="number" class="edit-input" />
            <span class="unit">см</span>
          </div>
          <span v-else class="value">{{ profile.height_cm }} см</span>
        </div>
        <div class="info-row">
          <span class="label">Вес</span>
          <div v-if="editing" class="edit-with-unit">
            <input v-model.number="editData.weight_kg" type="number" step="0.1" class="edit-input" />
            <span class="unit">кг</span>
          </div>
          <span v-else class="value">{{ profile.weight_kg }} кг</span>
        </div>
        <div class="info-row">
          <span class="label">Пол</span>
          <select v-if="editing" v-model="editData.gender" class="edit-select">
            <option value="male">Мужской</option>
            <option value="female">Женский</option>
          </select>
          <span v-else class="value">{{ genderLabel }}</span>
        </div>
        <div class="info-row">
          <span class="label">Уровень</span>
          <select v-if="editing" v-model="editData.fitness_level" class="edit-select">
            <option value="beginner">Новичок</option>
            <option value="intermediate">Средний</option>
            <option value="advanced">Продвинутый</option>
          </select>
          <span v-else class="value">{{ fitnessLabel }}</span>
        </div>
      </div>

      <!-- Goals -->
      <div class="info-card section-goals" v-if="profile.goals && profile.goals.length" style="animation-delay: 160ms">
        <h3><span class="section-icon">🎯</span> Цели</h3>
        <div class="goals-list">
          <span
            v-for="(goal, index) in profile.goals"
            :key="goal"
            class="goal-tag"
            :style="{ animationDelay: index * 60 + 'ms' }"
          >{{ goalLabel(goal) }}</span>
        </div>
      </div>

      <!-- Training info -->
      <div class="info-card section-training" v-if="profile.training_access || profile.training_experience" style="animation-delay: 240ms">
        <h3><span class="section-icon">💪</span> Тренировки</h3>
        <div class="info-row" v-if="profile.training_access">
          <span class="label">Доступ</span>
          <span class="value">{{ accessLabel }}</span>
        </div>
        <div class="info-row" v-if="profile.training_experience">
          <span class="label">Опыт</span>
          <span class="value">{{ experienceLabel }}</span>
        </div>
      </div>

      <!-- Schedule -->
      <div class="info-card section-schedule" v-if="editing || profile.days_per_week || profile.session_duration || profile.preferred_time" style="animation-delay: 320ms">
        <h3><span class="section-icon">📅</span> Расписание</h3>
        <div class="info-row">
          <span class="label">Дней в неделю</span>
          <input v-if="editing" v-model.number="editData.days_per_week" type="number" min="1" max="7" class="edit-input" />
          <span v-else class="value">{{ profile.days_per_week }}</span>
        </div>
        <div class="info-row">
          <span class="label">Длительность</span>
          <div v-if="editing" class="edit-with-unit">
            <input v-model.number="editData.session_duration" type="number" class="edit-input" />
            <span class="unit">мин</span>
          </div>
          <span v-else class="value">{{ profile.session_duration }} мин</span>
        </div>
        <div class="info-row">
          <span class="label">Время</span>
          <select v-if="editing" v-model="editData.preferred_time" class="edit-select">
            <option value="morning">Утро</option>
            <option value="afternoon">День</option>
            <option value="evening">Вечер</option>
            <option value="any">Любое</option>
          </select>
          <span v-else class="value">{{ timeLabel }}</span>
        </div>
        <div class="info-row" v-if="!editing && profile.equipment && profile.equipment.length">
          <span class="label">Оборудование</span>
          <span class="value">{{ equipmentText }}</span>
        </div>
      </div>

      <!-- Health info -->
      <div class="info-card health-card" v-if="profile.has_pain" style="animation-delay: 400ms">
        <h3><span class="section-icon">❤️</span> Здоровье</h3>
        <div class="info-row" v-if="profile.pain_locations && profile.pain_locations.length">
          <span class="label">Зоны боли</span>
          <span class="value">{{ painLocationsText }}</span>
        </div>
        <div class="info-row" v-if="profile.pain_level != null">
          <span class="label">Уровень боли</span>
          <span class="value" :style="{ color: painColor(profile.pain_level) }">{{ profile.pain_level }} / 10</span>
        </div>
        <div class="info-row" v-if="profile.diagnoses && profile.diagnoses.length">
          <span class="label">Диагнозы</span>
          <span class="value">{{ diagnosesText }}</span>
        </div>
        <div class="info-row" v-if="profile.contraindications">
          <span class="label">Противопоказания</span>
          <span class="value">{{ profile.contraindications }}</span>
        </div>
      </div>

      <!-- Payment -->
      <div class="info-card" style="animation-delay: 480ms">
        <div class="info-row">
          <span class="label">Доступ</span>
          <span class="value" :class="profile.is_paid ? 'paid' : 'unpaid'">
            {{ profile.is_paid ? 'Оплачено' : 'Не оплачено' }}
          </span>
        </div>
      </div>

      <div v-if="editing" class="edit-actions">
        <button class="btn btn-primary" @click="saveProfile" :disabled="saving">
          {{ saving ? 'Сохранение...' : 'Сохранить' }}
        </button>
        <div v-if="saveError" class="save-msg error-msg">{{ saveError }}</div>
        <div v-if="saveSuccess" class="save-msg success-msg">Профиль обновлён</div>
      </div>

      <button v-if="!editing && !profile.is_paid" class="btn btn-primary" @click="router.push('/payment')">
        Оплатить доступ
      </button>
      <button v-if="!editing && profile.role === 'admin'" class="btn btn-primary" @click="router.push('/admin')">
        Админ панель
      </button>
      <button v-if="!editing" class="btn btn-secondary" @click="router.push('/')">
        На главную
      </button>
    </div>

    <div v-else class="error">
      <p>Не удалось загрузить профиль</p>
      <button class="btn btn-secondary" @click="router.push('/')">Назад</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api } from '../api'
import type { UserProfile } from '../types'
import SkeletonCard from '../components/SkeletonCard.vue'
import { useMediaUpload } from '../composables/useMediaUpload'
import { useTelegram } from '../composables/useTelegram'

const { setClosingGuard } = useTelegram()

const route = useRoute()
const router = useRouter()

const loading = ref(true)
const profile = ref<UserProfile | null>(null)
const editing = ref(false)
const saving = ref(false)
const saveError = ref('')
const saveSuccess = ref(false)
const editData = ref<{
  age: number
  height_cm: number
  weight_kg: number
  gender: 'male' | 'female'
  fitness_level: 'beginner' | 'intermediate' | 'advanced'
  days_per_week: number
  session_duration: number
  preferred_time: string
}>({
  age: 0,
  height_cm: 0,
  weight_kg: 0,
  gender: 'male',
  fitness_level: 'beginner',
  days_per_week: 3,
  session_duration: 60,
  preferred_time: 'any',
})

function toggleEdit() {
  if (editing.value) {
    editing.value = false
    return
  }
  if (profile.value) {
    editData.value = {
      age: profile.value.age,
      height_cm: profile.value.height_cm,
      weight_kg: profile.value.weight_kg,
      gender: (profile.value.gender as 'male' | 'female') || 'male',
      fitness_level: (profile.value.fitness_level as 'beginner' | 'intermediate' | 'advanced') || 'beginner',
      days_per_week: profile.value.days_per_week || 3,
      session_duration: profile.value.session_duration || 60,
      preferred_time: profile.value.preferred_time || 'any',
    }
  }
  saveError.value = ''
  saveSuccess.value = false
  editing.value = true
}

async function saveProfile() {
  if (saving.value) return
  saving.value = true
  saveError.value = ''
  saveSuccess.value = false
  try {
    await api.updateProfile(editData.value)
    profile.value = await api.getProfile()
    saveSuccess.value = true
    setTimeout(() => {
      editing.value = false
      saveSuccess.value = false
    }, 1000)
  } catch (e: any) {
    saveError.value = e.message || 'Ошибка при сохранении'
  } finally {
    saving.value = false
  }
}

const goalLabels: Record<string, string> = {
  weight_loss: 'Похудеть',
  muscle_gain: 'Набрать массу',
  strength: 'Больше силы',
  endurance: 'Выносливость',
  maintenance: 'Поддержание формы',
  hernia: 'Грыжа',
  protrusion: 'Протрузии',
  scoliosis: 'Сколиоз',
  kyphosis: 'Кифоз',
  lordosis: 'Лордоз',
  lfk: 'ЛФК',
  training: 'Тренировки',
}

// Stored values are English keys (see web/src/views/onboarding/HealthSlide.vue);
// these dictionaries translate them for display in the read-only profile.
const painLocationLabels: Record<string, string> = {
  back: 'Спина',
  lower_back: 'Поясница',
  neck: 'Шея',
  shoulder: 'Плечо',
  knee: 'Колено',
  other: 'Другое',
}

const diagnosisLabels: Record<string, string> = {
  hernia: 'Грыжа',
  protrusion: 'Протрузия',
  scoliosis: 'Сколиоз',
  osteochondrosis: 'Остеохондроз',
  kyphosis: 'Кифоз',
  lordosis: 'Лордоз',
  other: 'Другое',
}

const equipmentLabels: Record<string, string> = {
  dumbbells: 'Гантели',
  barbell: 'Штанга',
  kettlebell: 'Гиря',
  bench: 'Скамья',
  pullup_bar: 'Турник',
  resistance_band: 'Резинки',
  mat: 'Коврик',
  none: 'Без инвентаря',
}

const accessLabels: Record<string, string> = {
  home: 'Дома',
  gym: 'Зал',
  outdoor: 'На улице',
  mixed: 'Смешанный',
}

const experienceLabels: Record<string, string> = {
  none: 'Нет опыта',
  less_1y: 'Менее 1 года',
  '1_3y': '1-3 года',
  more_3y: 'Более 3 лет',
}

const timeLabels: Record<string, string> = {
  morning: 'Утро',
  afternoon: 'День',
  evening: 'Вечер',
  any: 'Любое',
}

function goalLabel(key: string): string {
  return goalLabels[key] || key
}

function painLocationLabel(key: string): string {
  return painLocationLabels[key] || key
}

function diagnosisLabel(key: string): string {
  return diagnosisLabels[key] || key
}

function equipmentLabel(key: string): string {
  return equipmentLabels[key] || key
}

const painLocationsText = computed(() =>
  (profile.value?.pain_locations || []).map(painLocationLabel).join(', '),
)

const diagnosesText = computed(() =>
  (profile.value?.diagnoses || []).map(diagnosisLabel).join(', '),
)

const equipmentText = computed(() =>
  (profile.value?.equipment || []).map(equipmentLabel).join(', '),
)

function painColor(level: number): string {
  if (level <= 3) return '#34c759'
  if (level <= 6) return '#ffcc00'
  return '#ff3b30'
}

const initials = computed(() => {
  if (!profile.value) return '?'
  const f = profile.value.first_name?.[0] || ''
  const l = profile.value.last_name?.[0] || ''
  return (f + l).toUpperCase() || '?'
})

const genderLabel = computed(() => {
  return profile.value?.gender === 'male' ? 'Мужской' : 'Женский'
})

const fitnessLabel = computed(() => {
  const labels: Record<string, string> = {
    beginner: 'Новичок',
    intermediate: 'Средний',
    advanced: 'Продвинутый',
  }
  return labels[profile.value?.fitness_level || ''] || profile.value?.fitness_level || ''
})

const accessLabel = computed(() => {
  return accessLabels[profile.value?.training_access || ''] || profile.value?.training_access || ''
})

const experienceLabel = computed(() => {
  return experienceLabels[profile.value?.training_experience || ''] || profile.value?.training_experience || ''
})

const timeLabel = computed(() => {
  return timeLabels[profile.value?.preferred_time || ''] || profile.value?.preferred_time || ''
})

// Avatar upload
const avatarInput = ref<HTMLInputElement | null>(null)
const avatarUrl = ref<string | null>(null)
const avatarError = ref('')
const { uploading: avatarUploading, progress: avatarProgress, upload: uploadAvatar } = useMediaUpload()

async function loadAvatar(mediaId: number | null | undefined) {
  if (!mediaId) {
    avatarUrl.value = null
    return
  }
  try {
    const res = await api.getMediaURL(mediaId)
    avatarUrl.value = res.url
  } catch {
    avatarUrl.value = null
  }
}

function pickAvatar() {
  if (avatarUploading.value) return
  avatarInput.value?.click()
}

async function handleAvatarChange(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  // Reset the input so picking the same file again still triggers change
  input.value = ''
  avatarError.value = ''
  try {
    const result = await uploadAvatar(file, { reference_type: 'user_avatar' })
    await api.updateProfile({ avatar_media_id: result.media_id } as any)
    if (profile.value) profile.value.avatar_media_id = result.media_id
    await loadAvatar(result.media_id)
  } catch (err: any) {
    avatarError.value = err.message || 'Ошибка загрузки фото'
  }
}

// Guard the close button only while the form is in edit mode (dirty surface).
// Browsing the read-only profile is safe to leave at any time.
watch(editing, (val) => setClosingGuard(val))

onMounted(async () => {
  try {
    profile.value = await api.getProfile()
    if (profile.value?.avatar_media_id) {
      await loadAvatar(profile.value.avatar_media_id)
    }
  } catch {
    profile.value = null
  } finally {
    loading.value = false
  }
})

onUnmounted(() => setClosingGuard(false))
</script>

<style scoped>
.profile-page {
  max-width: 400px;
  margin: 0 auto;
}


.skeleton-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.profile-header {
  text-align: center;
  margin-bottom: 16px;
  padding: 32px 16px 20px;
  background: var(--secondary-bg);
  border-radius: 16px;
  position: relative;
  overflow: hidden;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease forwards;
}

.header-gradient {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 96px;
  background: linear-gradient(135deg, var(--button-color), var(--link-color));
  opacity: 0.22;
}

/* Outer .avatar is the positioning context for the edit badge.
   Inner .avatar-circle does the circular masking. This split lets the
   badge sit OUTSIDE the circle's boundary while still inside the card. */
.avatar {
  width: 104px;
  height: 104px;
  margin: 0 auto 14px;
  position: relative;
  z-index: 1;
}

.avatar-circle {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--button-color), var(--link-color));
  color: var(--button-text-color);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 38px;
  font-weight: 700;
  overflow: hidden;
  box-shadow: 0 4px 14px rgba(0, 0, 0, 0.12);
  transition: transform 0.12s;
}

.avatar-button {
  cursor: pointer;
  -webkit-tap-highlight-color: transparent;
}

.avatar-button:active .avatar-circle {
  transform: scale(0.96);
}

.avatar-button.uploading {
  cursor: wait;
}

.avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
  border-radius: 50%;
}

/* Floating edit badge — anchored at the bottom-right of the avatar
   bounding box (outside the inscribed circle). Thick page-bg border +
   shadow gives the iOS-style "floats above" feel. */
.avatar-edit-badge {
  position: absolute;
  bottom: -2px;
  right: -2px;
  width: 34px;
  height: 34px;
  background: var(--button-color);
  color: var(--button-text-color);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 15px;
  border: 3px solid var(--bg-color);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.18);
  z-index: 2;
}

.avatar-error {
  color: #ff3b30;
  font-size: 13px;
  margin-top: 6px;
}

.profile-name {
  font-size: 22px;
  font-weight: 700;
  margin: 4px 0 2px;
  letter-spacing: -0.2px;
}

.username {
  color: var(--hint-color);
  font-size: 14px;
}

/* Vitals tile row — Apple-Health-style quick-glance stats */
.vitals-row {
  display: flex;
  gap: 8px;
  margin: 16px 0 4px;
  justify-content: center;
}

.vital-tile {
  flex: 1;
  max-width: 110px;
  background: var(--bg-color);
  border-radius: 12px;
  padding: 10px 8px 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}

.vital-value {
  font-size: 20px;
  font-weight: 700;
  color: var(--text-color);
  line-height: 1.1;
}

.vital-label {
  font-size: 11px;
  color: var(--hint-color);
  margin-top: 4px;
  letter-spacing: 0.1px;
}

.info-card {
  background: var(--secondary-bg);
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 12px;
  opacity: 0;
  animation: fadeSlideUp 0.35s ease forwards;
}

.info-card h3 {
  margin-bottom: 12px;
  font-size: 18px;
  font-weight: 700;
  display: flex;
  align-items: center;
  gap: 8px;
}

.section-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 8px;
  font-size: 16px;
  flex-shrink: 0;
}

.section-goals .section-icon { background: rgba(52, 199, 89, 0.18); }
.section-training .section-icon { background: rgba(36, 129, 204, 0.18); }
.section-schedule .section-icon { background: rgba(255, 149, 0, 0.18); }
.health-card .section-icon { background: rgba(255, 59, 48, 0.18); }

.section-goals { border-left: 3px solid #34c759; }
.section-training { border-left: 3px solid var(--button-color); }
.section-schedule { border-left: 3px solid #ff9500; }
.health-card { border-left: 3px solid #ff3b30; }

.info-row {
  display: flex;
  justify-content: space-between;
  padding: 10px 0;
}

.label {
  color: var(--hint-color);
}

.value {
  text-align: right;
  max-width: 55%;
}

.value.paid {
  color: #34c759;
}

.value.unpaid {
  color: #ff3b30;
}

.goals-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.goal-tag {
  background: var(--button-color);
  color: var(--button-text-color);
  padding: 4px 12px;
  border-radius: 16px;
  font-size: 13px;
  opacity: 0;
  animation: bounceIn 0.4s ease forwards;
}

@keyframes bounceIn {
  0% { opacity: 0; transform: scale(0.6); }
  60% { opacity: 1; transform: scale(1.05); }
  100% { opacity: 1; transform: scale(1); }
}

.btn {
  display: block;
  width: 100%;
  padding: 14px;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  cursor: pointer;
  margin-bottom: 8px;
  text-align: center;
}

.btn-primary {
  background: var(--button-color);
  color: var(--button-text-color);
}

.btn-secondary {
  background: var(--secondary-bg);
  color: var(--text-color);
}

.error {
  text-align: center;
  padding: 40px;
  color: var(--hint-color);
}

.edit-toggle {
  background: var(--bg-color);
  border: none;
  color: var(--button-color);
  padding: 8px 20px;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  margin-top: 16px;
  position: relative;
  z-index: 1;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  transition: transform 0.1s;
  -webkit-tap-highlight-color: transparent;
}

.edit-toggle:active {
  transform: scale(0.97);
}

.edit-input {
  background: var(--bg-color);
  border: 1px solid var(--hint-color);
  border-radius: 8px;
  padding: 6px 10px;
  font-size: 14px;
  color: var(--text-color);
  text-align: right;
  width: 80px;
}

.edit-select {
  background: var(--bg-color);
  border: 1px solid var(--hint-color);
  border-radius: 8px;
  padding: 6px 10px;
  font-size: 14px;
  color: var(--text-color);
  text-align: right;
}

.edit-with-unit {
  display: flex;
  align-items: center;
  gap: 4px;
}

.unit {
  font-size: 13px;
  color: var(--hint-color);
}

.edit-actions {
  margin-bottom: 12px;
}

.save-msg {
  text-align: center;
  font-size: 14px;
  margin-top: 8px;
}

.error-msg {
  color: #ff3b30;
}

.success-msg {
  color: #34c759;
}
</style>
