<template>
  <div class="admin-page">
    <h1 class="page-title">{{ isEdit ? 'Редактировать занятие' : 'Новое занятие' }}</h1>

    <div v-if="loading" class="loading">Загрузка...</div>

    <form v-else class="admin-form" @submit.prevent="save">
      <fieldset class="admin-section">
        <legend class="admin-section-title">Основные данные</legend>
        <div class="field">
          <label>Курс ID *</label>
          <input v-model.number="form.course_id" type="number" required :disabled="isEdit" />
        </div>
        <div class="row">
          <div class="field">
            <label>День</label>
            <input v-model.number="form.day_number" type="number" min="1" />
          </div>
          <div class="field">
            <label>Этап</label>
            <input v-model.number="form.stage" type="number" min="1" max="3" />
          </div>
        </div>
        <div class="row">
          <div class="field">
            <label>Длительность (мин)</label>
            <input v-model.number="form.duration_minutes" type="number" min="0" />
          </div>
          <div class="field">
            <label>Порядок</label>
            <input v-model.number="form.sort_order" type="number" />
          </div>
        </div>
        <div class="field">
          <label>Описание</label>
          <textarea v-model="form.description" rows="3"></textarea>
        </div>
      </fieldset>

      <fieldset class="admin-section">
        <legend class="admin-section-title">Видео</legend>
        <div class="field">
          <label>Загрузить видео (MP4 в R2)</label>
          <VideoUploader v-model="form.video_media_id" reference-type="rehab_session_video" />
        </div>
        <div class="field">
          <label>Или внешняя ссылка (YouTube и т.п.)</label>
          <input v-model="form.video_url" placeholder="Использовать вместо загрузки" />
        </div>
      </fieldset>

      <div v-if="error" class="error-msg">{{ error }}</div>
      <div class="admin-save-bar">
        <button type="submit" class="admin-save-btn" :disabled="saving">
          {{ saving ? 'Сохранение...' : 'Сохранить' }}
        </button>
        <button v-if="form.course_id" type="button" class="btn-secondary" @click="router.push(`/admin/rehab/courses/${form.course_id}`)">
          Назад к курсу
        </button>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api } from '../../api'
import VideoUploader from '../../components/VideoUploader.vue'

const props = defineProps<{ id?: string }>()
const route = useRoute()
const router = useRouter()
const isEdit = !!props.id
const loading = ref(isEdit)
const saving = ref(false)
const error = ref('')

const form = reactive({
  course_id: 0,
  day_number: 1,
  stage: 1,
  video_url: '',
  video_media_id: null as number | null,
  duration_minutes: 15,
  description: '',
  sort_order: 0,
})

onMounted(async () => {
  if (isEdit) {
    try {
      const s = await api.getAdminRehabSession(Number(props.id))
      Object.assign(form, s)
    } catch {
      error.value = 'Не удалось загрузить занятие'
    } finally {
      loading.value = false
    }
  } else {
    // Pre-fill course_id from ?course_id= query param when navigating from a course
    const qp = route.query.course_id
    if (typeof qp === 'string') {
      const id = Number(qp)
      if (id > 0) form.course_id = id
    }
  }
})

async function save() {
  saving.value = true
  error.value = ''
  try {
    if (isEdit) {
      await api.updateAdminRehabSession(Number(props.id), { ...form })
      router.replace(`/admin/rehab/courses/${form.course_id}`)
    } else {
      const created = await api.createAdminRehabSession({ ...form })
      router.replace(`/admin/rehab/courses/${created.course_id}`)
    }
  } catch (e: any) {
    error.value = e.message || 'Ошибка сохранения'
  } finally {
    saving.value = false
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
.field input:disabled { opacity: 0.6; }
.field textarea { resize: vertical; }
.row { display: flex; gap: 12px; }
.row .field { flex: 1; }
.hint { font-size: 12px; color: var(--hint-color); margin-top: 4px; }
.error-msg { color: #ff3b30; font-size: 14px; text-align: center; }
.btn-primary {
  padding: 14px; border-radius: 12px; border: none;
  background: var(--button-color); color: var(--button-text-color);
  font-size: 16px; font-weight: 600; cursor: pointer;
}
.btn-primary:disabled { opacity: 0.5; }
.btn-secondary {
  padding: 12px; border-radius: 12px; border: 1px solid var(--hint-color);
  background: transparent; color: var(--text-color); font-size: 14px; cursor: pointer;
}
</style>
