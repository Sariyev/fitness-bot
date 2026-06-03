<template>
  <div class="admin-page">
    <h1 class="page-title">{{ isEdit ? 'Редактировать программу' : 'Новая программа' }}</h1>

    <div v-if="loading" class="loading">Загрузка...</div>

    <form v-else class="admin-form" @submit.prevent="save">
      <fieldset class="admin-section">
        <legend class="admin-section-title"><span class="icon">📝</span>Основные данные</legend>
        <div class="field">
          <label>Название *</label>
          <input v-model="form.name" required />
        </div>
        <div class="field">
          <label>Описание</label>
          <textarea v-model="form.description" rows="3"></textarea>
        </div>
      </fieldset>

      <fieldset class="admin-section">
        <legend class="admin-section-title icon-goals"><span class="icon">🎯</span>Параметры программы</legend>
        <div class="row">
          <div class="field">
            <label>Цель</label>
            <select v-model="form.goal">
              <option value="">—</option>
              <option value="weight_loss">Похудение</option>
              <option value="muscle_gain">Набор массы</option>
              <option value="general_fitness">Общая форма</option>
            </select>
          </div>
          <div class="field">
            <label>Формат</label>
            <select v-model="form.format">
              <option value="">—</option>
              <option value="home">Дома</option>
              <option value="gym">Зал</option>
            </select>
          </div>
        </div>
        <div class="row">
          <div class="field">
            <label>Уровень</label>
            <select v-model="form.level">
              <option value="">—</option>
              <option value="beginner">Начинающий</option>
              <option value="intermediate">Средний</option>
              <option value="advanced">Продвинутый</option>
            </select>
          </div>
          <div class="field">
            <label>Недель</label>
            <input v-model.number="form.duration_weeks" type="number" min="1" max="52" />
          </div>
        </div>
      </fieldset>

      <fieldset class="admin-section">
        <legend class="admin-section-title icon-access"><span class="icon">🔓</span>Доступ</legend>
        <div class="field">
          <label>Бакет</label>
          <div class="toggle-group">
            <button type="button" :class="{ active: form.access_tier === 'free' }" @click="form.access_tier = 'free'">Бесплатно</button>
            <button type="button" :class="{ active: form.access_tier === 'trial' }" @click="form.access_tier = 'trial'">Триал</button>
            <button type="button" :class="{ active: form.access_tier === 'paid' }" @click="form.access_tier = 'paid'">Платно</button>
          </div>
          <span class="hint">Бесплатно — всем; Триал — первые 7 дней после регистрации; Платно — только купившим категорию.</span>
        </div>
      </fieldset>

      <fieldset class="admin-section">
        <legend class="admin-section-title icon-status"><span class="icon">⚡️</span>Статус</legend>
        <div class="row">
          <div class="field">
            <label>Порядок</label>
            <input v-model.number="form.sort_order" type="number" min="0" max="9999" />
          </div>
          <div class="field">
            <label>Активна</label>
            <div class="toggle-group">
              <button type="button" :class="{ active: form.is_active }" @click="form.is_active = true">Да</button>
              <button type="button" :class="{ active: !form.is_active }" @click="form.is_active = false">Нет</button>
            </div>
          </div>
        </div>
      </fieldset>

      <div v-if="error" class="error-msg">{{ error }}</div>
      <div class="admin-save-bar">
        <button type="submit" class="admin-save-btn" :disabled="saving">
          {{ saving ? 'Сохранение...' : 'Сохранить' }}
        </button>
      </div>

      <AdminDangerZone
        v-if="isEdit"
        ref="dangerZone"
        label="программу"
        :is-active="form.is_active"
        @delete="onDelete"
        @toggle-active="onToggleActive"
      />
    </form>

    <!-- Workouts attached to this program (edit mode only) -->
    <div v-if="isEdit && !loading" class="workouts-section">
      <div class="section-header">
        <h2 class="section-title">Тренировки ({{ programWorkouts.length }})</h2>
        <button class="add-btn" @click="router.push(`/admin/workouts/new?program_id=${props.id}`)">+</button>
      </div>
      <div v-if="programWorkouts.length === 0" class="empty">
        Нет тренировок в программе. Нажмите «+» чтобы добавить первую.
      </div>
      <div
        v-for="w in programWorkouts"
        :key="w.id"
        class="content-card"
        :class="{ inactive: !w.is_active }"
        @click="router.push(`/admin/workouts/${w.id}`)"
      >
        <div class="content-main">
          <span class="content-name">{{ w.name }}</span>
          <span class="content-meta">
            <template v-if="w.week_number">Нед. {{ w.week_number }}</template>
            <template v-if="w.week_number && w.day_number"> · </template>
            <template v-if="w.day_number">День {{ w.day_number }}</template>
            <template v-if="!w.week_number && !w.day_number">без расписания</template>
            <template v-if="w.duration_minutes"> · {{ w.duration_minutes }} мин</template>
          </span>
        </div>
        <span class="arrow">→</span>
      </div>
      <p class="hint">
        Упражнения добавляются на странице самой тренировки — откройте её и прокрутите до раздела «Упражнения».
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../api'
import AdminDangerZone from '../../components/AdminDangerZone.vue'
import type { Workout } from '../../types'

const props = defineProps<{ id?: string }>()
const router = useRouter()
const isEdit = !!props.id
const loading = ref(isEdit)
const saving = ref(false)
const error = ref('')

const form = reactive({
  name: '',
  slug: '',
  description: '',
  goal: '',
  format: '',
  level: '',
  duration_weeks: 4,
  access_tier: 'paid' as 'free' | 'trial' | 'paid',
  sort_order: 0,
  is_active: true,
})

const allWorkouts = ref<Workout[]>([])
const programWorkouts = computed(() => {
  const id = Number(props.id)
  return allWorkouts.value
    .filter((w) => w.program_id === id)
    .sort((a, b) => {
      const aw = a.week_number ?? 999, bw = b.week_number ?? 999
      if (aw !== bw) return aw - bw
      const ad = a.day_number ?? 999, bd = b.day_number ?? 999
      return ad - bd
    })
})

onMounted(async () => {
  if (isEdit) {
    try {
      const [p, ws] = await Promise.all([
        api.getAdminProgram(Number(props.id)),
        api.getAdminWorkouts(),
      ])
      Object.assign(form, p)
      allWorkouts.value = ws || []
    } catch {
      error.value = 'Не удалось загрузить программу'
    } finally {
      loading.value = false
    }
  }
})

async function save() {
  saving.value = true
  error.value = ''
  try {
    if (isEdit) {
      await api.updateAdminProgram(Number(props.id), { ...form })
    } else {
      await api.createAdminProgram({ ...form })
    }
    router.replace('/admin/content')
  } catch (e: any) {
    error.value = e.message || 'Ошибка сохранения'
  } finally {
    saving.value = false
  }
}

const dangerZone = ref<{ reset: () => void } | null>(null)

async function onDelete() {
  if (!isEdit) return
  try {
    await api.deleteAdminProgram(Number(props.id))
    router.replace('/admin/content')
  } catch (e: any) {
    error.value = e.message || 'Не удалось удалить'
    dangerZone.value?.reset()
  }
}

async function onToggleActive(next: boolean) {
  if (!isEdit) return
  try {
    await api.updateAdminProgram(Number(props.id), { ...form, is_active: next })
    form.is_active = next
  } catch (e: any) {
    error.value = e.message || 'Не удалось изменить статус'
  }
}
</script>

<style scoped>
.admin-page { max-width: 400px; margin: 0 auto; padding-bottom: 24px; }
.page-title { font-size: 20px; font-weight: 700; margin-bottom: 16px; }
.loading { text-align: center; color: var(--hint-color); padding: 40px; }
.form { display: flex; flex-direction: column; gap: 14px; }
.field { display: flex; flex-direction: column; gap: 4px; }
.field label { font-size: 13px; color: var(--hint-color); font-weight: 500; }
.field input, .field select, .field textarea {
  padding: 10px 12px; border-radius: 10px; border: 1px solid var(--hint-color);
  background: var(--secondary-bg); color: var(--text-color); font-size: 15px;
}
.field textarea { resize: vertical; }
.row { display: flex; gap: 12px; }
.row .field { flex: 1; }
.toggle-group { display: flex; gap: 4px; }
.toggle-group button {
  flex: 1; padding: 8px; border-radius: 8px; border: 1px solid var(--hint-color);
  background: var(--secondary-bg); color: var(--text-color); font-size: 14px; cursor: pointer;
}
.toggle-group button.active { background: var(--button-color); color: var(--button-text-color); border-color: var(--button-color); }
.error-msg { color: #ff3b30; font-size: 14px; text-align: center; }
.btn-primary {
  padding: 14px; border-radius: 12px; border: none;
  background: var(--button-color); color: var(--button-text-color);
  font-size: 16px; font-weight: 600; cursor: pointer;
}
.btn-primary:disabled { opacity: 0.5; }
.hint { font-size: 12px; color: var(--hint-color); margin-top: 4px; line-height: 1.4; }

/* ===== Workouts attached to this program ===== */
.workouts-section { margin-top: 24px; }
.section-header {
  display: flex; align-items: center; justify-content: space-between;
  margin: 16px 0 10px;
}
.section-title { font-size: 16px; font-weight: 600; color: var(--hint-color); }
.add-btn {
  width: 32px; height: 32px; border-radius: 50%; border: none;
  background: var(--button-color); color: var(--button-text-color);
  font-size: 20px; font-weight: 700; cursor: pointer;
  display: flex; align-items: center; justify-content: center;
}
.content-card {
  background: var(--secondary-bg); border-radius: 12px; padding: 12px 16px;
  margin-bottom: 8px; display: flex; justify-content: space-between; align-items: center;
  cursor: pointer; touch-action: manipulation;
}
.content-card.inactive .content-name { color: var(--hint-color); }
.content-main { display: flex; flex-direction: column; min-width: 0; flex: 1; }
.content-name { font-weight: 500; font-size: 15px; }
.content-meta { color: var(--hint-color); font-size: 12px; margin-top: 2px; }
.arrow { font-size: 18px; color: var(--hint-color); flex-shrink: 0; margin-left: 12px; }
.empty { text-align: center; color: var(--hint-color); padding: 18px; font-size: 13px; }
</style>
