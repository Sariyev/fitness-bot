<template>
  <div class="admin-page">
    <button class="back-btn" @click="router.push('/admin/content')">← Назад</button>
    <h1 class="page-title">{{ isEdit ? 'Редактировать тренировку' : 'Новая тренировка' }}</h1>

    <div v-if="loading" class="loading">Загрузка...</div>

    <form v-else class="form" @submit.prevent="save">
      <div class="field">
        <label>Название *</label>
        <input v-model="form.name" required />
      </div>
      <div class="field">
        <label>Slug</label>
        <input v-model="form.slug" placeholder="auto-from-name" />
      </div>
      <div class="field">
        <label>Описание</label>
        <textarea v-model="form.description" rows="3"></textarea>
      </div>
      <div class="field">
        <label>Программа</label>
        <select v-model="form.program_id">
          <option :value="undefined">— Без программы —</option>
          <option v-for="p in programs" :key="p.id" :value="p.id">{{ p.name }}</option>
        </select>
      </div>
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
          <label>Минут</label>
          <input v-model.number="form.duration_minutes" type="number" min="1" />
        </div>
      </div>
      <div class="field">
        <label>Оборудование (через запятую)</label>
        <input v-model="equipmentStr" />
      </div>
      <div class="field">
        <label>Ожидаемый результат</label>
        <textarea v-model="form.expected_result" rows="2"></textarea>
      </div>
      <div class="field">
        <label>Видео URL</label>
        <input v-model="form.video_url" />
      </div>
      <div v-if="form.program_id" class="row">
        <div class="field">
          <label>Неделя</label>
          <input v-model.number="form.week_number" type="number" min="1" />
        </div>
        <div class="field">
          <label>День</label>
          <input v-model.number="form.day_number" type="number" min="1" />
        </div>
      </div>
      <div class="row">
        <div class="field">
          <label>Порядок</label>
          <input v-model.number="form.sort_order" type="number" />
        </div>
        <div class="field">
          <label>Активна</label>
          <div class="toggle-group">
            <button type="button" :class="{ active: form.is_active }" @click="form.is_active = true">Да</button>
            <button type="button" :class="{ active: !form.is_active }" @click="form.is_active = false">Нет</button>
          </div>
        </div>
      </div>

      <div v-if="error" class="error-msg">{{ error }}</div>
      <button type="submit" class="btn btn-primary" :disabled="saving">
        {{ saving ? 'Сохранение...' : 'Сохранить' }}
      </button>
    </form>

    <!-- Exercises section for edit mode -->
    <div v-if="isEdit && !loading" class="exercises-section">
      <div class="section-header">
        <h2 class="section-title">Упражнения ({{ exercises.length }})</h2>
      </div>
      <div v-for="ex in exercises" :key="ex.id" class="exercise-card">
        <span class="ex-name">{{ ex.exercise_name }}</span>
        <span class="ex-meta">{{ ex.sets }}×{{ ex.reps }} | отдых {{ ex.rest_seconds }}с</span>
      </div>

      <details class="add-exercise-form">
        <summary>Добавить упражнение</summary>
        <div class="form" style="margin-top: 12px;">
          <div class="field">
            <label>Упражнение</label>
            <select v-model.number="newEx.exercise_id">
              <option :value="0">— Выбрать —</option>
              <option v-for="e in allExercises" :key="e.id" :value="e.id">{{ e.name }}</option>
            </select>
          </div>
          <div class="row">
            <div class="field">
              <label>Подходы</label>
              <input v-model.number="newEx.sets" type="number" min="1" />
            </div>
            <div class="field">
              <label>Повторы</label>
              <input v-model="newEx.reps" />
            </div>
          </div>
          <div class="row">
            <div class="field">
              <label>Сек (длит.)</label>
              <input v-model.number="newEx.duration_seconds" type="number" />
            </div>
            <div class="field">
              <label>Порядок</label>
              <input v-model.number="newEx.sort_order" type="number" />
            </div>
          </div>
          <button type="button" class="btn btn-primary" :disabled="addingEx" @click="addExercise">
            {{ addingEx ? 'Добавление...' : 'Добавить' }}
          </button>
        </div>
      </details>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../api'
import type { Program, WorkoutExercise, Exercise } from '../../types'

const props = defineProps<{ id?: string }>()
const router = useRouter()
const isEdit = !!props.id
const loading = ref(true)
const saving = ref(false)
const error = ref('')

const programs = ref<Program[]>([])
const exercises = ref<WorkoutExercise[]>([])
const allExercises = ref<Exercise[]>([])
const addingEx = ref(false)

const form = reactive({
  name: '',
  slug: '',
  description: '',
  program_id: undefined as number | undefined,
  goal: '',
  format: '',
  level: '',
  duration_minutes: 30,
  equipment: [] as string[],
  expected_result: '',
  video_url: '',
  sort_order: 0,
  week_number: undefined as number | undefined,
  day_number: undefined as number | undefined,
  is_active: true,
})

const equipmentStr = computed({
  get: () => form.equipment?.join(', ') || '',
  set: (v: string) => { form.equipment = v.split(',').map(s => s.trim()).filter(Boolean) },
})

const newEx = reactive({ exercise_id: 0, sets: 3, reps: '10-12', duration_seconds: 0, sort_order: exercises.value.length })

onMounted(async () => {
  try {
    const [p, exList] = await Promise.all([
      api.getAdminPrograms(),
      api.getAdminExercises(),
    ])
    programs.value = p || []
    allExercises.value = exList || []

    if (isEdit) {
      const data = await api.getAdminWorkout(Number(props.id))
      Object.assign(form, data.workout)
      exercises.value = data.exercises || []
    }
  } catch {
    error.value = 'Ошибка загрузки'
  } finally {
    loading.value = false
  }
})

async function save() {
  saving.value = true
  error.value = ''
  try {
    if (isEdit) {
      await api.updateAdminWorkout(Number(props.id), { ...form })
    } else {
      await api.createAdminWorkout({ ...form })
    }
    router.push('/admin/content')
  } catch (e: any) {
    error.value = e.message || 'Ошибка сохранения'
  } finally {
    saving.value = false
  }
}

async function addExercise() {
  if (!newEx.exercise_id) return
  addingEx.value = true
  try {
    await api.addWorkoutExercise({
      workout_id: Number(props.id),
      exercise_id: newEx.exercise_id,
      sets: newEx.sets,
      reps: newEx.reps,
      duration_seconds: newEx.duration_seconds,
      sort_order: newEx.sort_order,
    })
    const data = await api.getAdminWorkout(Number(props.id))
    exercises.value = data.exercises || []
    newEx.exercise_id = 0
    newEx.sort_order = exercises.value.length
  } catch {
    // ignore
  } finally {
    addingEx.value = false
  }
}
</script>

<style scoped>
.admin-page { max-width: 400px; margin: 0 auto; padding-bottom: 24px; }
.back-btn { background: none; border: none; color: var(--button-color); font-size: 16px; cursor: pointer; padding: 4px 0; margin-bottom: 12px; }
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
.exercises-section { margin-top: 24px; }
.section-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px; }
.section-title { font-size: 16px; font-weight: 600; color: var(--hint-color); }
.exercise-card {
  background: var(--secondary-bg); border-radius: 10px; padding: 10px 14px; margin-bottom: 6px;
  display: flex; flex-direction: column;
}
.ex-name { font-weight: 500; font-size: 14px; }
.ex-meta { color: var(--hint-color); font-size: 12px; margin-top: 2px; }
.add-exercise-form { margin-top: 12px; }
.add-exercise-form summary { color: var(--button-color); font-weight: 500; cursor: pointer; font-size: 14px; }
</style>
