<template>
  <div class="admin-page">
    <h1 class="page-title">{{ isEdit ? 'Редактировать курс ЛФК' : 'Новый курс ЛФК' }}</h1>

    <div v-if="loading" class="loading">Загрузка...</div>

    <form v-else class="admin-form" @submit.prevent="save">
      <fieldset class="admin-section">
        <legend class="admin-section-title">Основные данные</legend>
        <div class="field">
          <label>Название *</label>
          <input v-model="form.name" required />
        </div>
        <div class="row">
          <div class="field">
            <label>Категория</label>
            <select v-model="form.category">
              <option value="">—</option>
              <option value="back">Спина</option>
              <option value="lower_back">Поясница</option>
              <option value="neck">Шея</option>
              <option value="shoulder">Плечо</option>
              <option value="knee">Колено</option>
              <option value="general">Общая</option>
            </select>
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
        <div class="field">
          <label>Предупреждения</label>
          <textarea v-model="form.warnings" rows="2" placeholder="Противопоказания, важные замечания"></textarea>
        </div>
      </fieldset>

      <fieldset class="admin-section">
        <legend class="admin-section-title">Доступ</legend>
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
        <legend class="admin-section-title">Статус</legend>
        <div class="field">
          <label>Активен</label>
          <div class="toggle-group">
            <button type="button" :class="{ active: form.is_active }" @click="form.is_active = true">Да</button>
            <button type="button" :class="{ active: !form.is_active }" @click="form.is_active = false">Нет</button>
          </div>
        </div>
      </fieldset>

      <div v-if="error" class="error-msg">{{ error }}</div>
      <div class="admin-save-bar">
        <button type="submit" class="admin-save-btn" :disabled="saving">
          {{ saving ? 'Сохранение...' : 'Сохранить' }}
        </button>
      </div>
    </form>

    <!-- Sessions list (only for existing course) -->
    <div v-if="isEdit && !loading" class="sessions-section">
      <div class="section-header">
        <h2 class="section-title">Занятия ({{ sessions.length }})</h2>
        <button class="add-btn" type="button" @click="addSession">+</button>
      </div>
      <div v-if="sessions.length === 0" class="empty-hint">
        Нет занятий. Нажмите «+» чтобы добавить первое.
      </div>
      <div
        v-for="s in sessions"
        :key="s.id"
        class="content-card"
        @click="router.push(`/admin/rehab/sessions/${s.id}`)"
      >
        <div class="content-main">
          <span class="content-name">День {{ s.day_number }} · Этап {{ s.stage }}</span>
          <span class="content-meta">{{ s.duration_minutes }} мин{{ s.video_url ? ' · 🎥' : '' }}</span>
        </div>
        <span class="arrow">→</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../api'
import type { RehabSession } from '../../types'

const props = defineProps<{ id?: string }>()
const router = useRouter()
const isEdit = !!props.id
const loading = ref(isEdit)
const saving = ref(false)
const error = ref('')

const form = reactive({
  name: '',
  slug: '',
  category: '',
  description: '',
  warnings: '',
  access_tier: 'paid' as 'free' | 'trial' | 'paid',
  sort_order: 0,
  is_active: true,
})

const sessions = ref<RehabSession[]>([])

onMounted(async () => {
  if (isEdit) {
    try {
      const data = await api.getAdminRehabCourse(Number(props.id))
      Object.assign(form, data.course)
      sessions.value = data.sessions || []
    } catch {
      error.value = 'Не удалось загрузить курс'
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
      await api.updateAdminRehabCourse(Number(props.id), { ...form })
    } else {
      const created = await api.createAdminRehabCourse({ ...form })
      router.replace(`/admin/rehab/courses/${created.id}`)
      return
    }
    router.replace('/admin/content')
  } catch (e: any) {
    error.value = e.message || 'Ошибка сохранения'
  } finally {
    saving.value = false
  }
}

function addSession() {
  router.push(`/admin/rehab/sessions/new?course_id=${props.id}`)
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
.toggle-group button.active {
  background: var(--button-color); color: var(--button-text-color); border-color: var(--button-color);
}
.error-msg { color: #ff3b30; font-size: 14px; text-align: center; }
.btn-primary {
  padding: 14px; border-radius: 12px; border: none;
  background: var(--button-color); color: var(--button-text-color);
  font-size: 16px; font-weight: 600; cursor: pointer;
}
.btn-primary:disabled { opacity: 0.5; }

.sessions-section { margin-top: 28px; }
.section-header {
  display: flex; justify-content: space-between; align-items: center;
  margin: 16px 0 10px;
}
.section-title { font-size: 16px; font-weight: 600; color: var(--hint-color); }
.add-btn {
  width: 32px; height: 32px; border-radius: 50%; border: none;
  background: var(--button-color); color: var(--button-text-color);
  font-size: 20px; font-weight: 700; cursor: pointer;
  display: flex; align-items: center; justify-content: center;
}
.empty-hint { color: var(--hint-color); font-size: 14px; padding: 12px; text-align: center; }
.content-card {
  background: var(--secondary-bg); border-radius: 12px;
  padding: 12px 16px; margin-bottom: 8px;
  display: flex; justify-content: space-between; align-items: center; cursor: pointer;
}
.content-main { display: flex; flex-direction: column; }
.content-name { font-weight: 500; font-size: 15px; }
.content-meta { color: var(--hint-color); font-size: 12px; margin-top: 2px; }
.arrow { font-size: 18px; color: var(--hint-color); }
.hint { font-size: 12px; color: var(--hint-color); margin-top: 4px; line-height: 1.4; }
</style>
