<template>
  <div class="admin-page">
    <button class="back-btn" @click="router.push('/admin/content')">← Назад</button>
    <h1 class="page-title">{{ isEdit ? 'Редактировать программу' : 'Новая программа' }}</h1>

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
          <label>Ф��рмат</label>
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
          <input v-model.number="form.duration_weeks" type="number" min="1" />
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
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../api'

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
  sort_order: 0,
  is_active: true,
})

onMounted(async () => {
  if (isEdit) {
    try {
      const p = await api.getAdminProgram(Number(props.id))
      Object.assign(form, p)
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
    router.push('/admin/content')
  } catch (e: any) {
    error.value = e.message || 'Ошибка сохранения'
  } finally {
    saving.value = false
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
</style>
