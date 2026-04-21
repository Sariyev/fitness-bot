<template>
  <div class="admin-page">
    <button class="back-btn" @click="router.push('/admin/exercises')">← Назад</button>
    <h1 class="page-title">{{ isEdit ? 'Редактировать упражнение' : 'Новое упражнение' }}</h1>

    <div v-if="loading" class="loading">Загрузка...</div>

    <form v-else class="form" @submit.prevent="save">
      <div class="field">
        <label>Название *</label>
        <input v-model="form.name" required />
      </div>
      <div class="field">
        <label>Техника выполнения</label>
        <textarea v-model="form.technique" rows="3"></textarea>
      </div>
      <div class="field">
        <label>Частые ошибки</label>
        <textarea v-model="form.common_mistakes" rows="3"></textarea>
      </div>
      <div class="field">
        <label>Упрощение</label>
        <textarea v-model="form.easier_modification" rows="2"></textarea>
      </div>
      <div class="field">
        <label>Усложнение</label>
        <textarea v-model="form.harder_modification" rows="2"></textarea>
      </div>
      <div class="field">
        <label>Отдых (сек)</label>
        <input v-model.number="form.rest_seconds" type="number" min="0" />
      </div>

      <div v-if="error" class="error-msg">{{ error }}</div>
      <button type="submit" class="btn btn-primary" :disabled="saving">
        {{ saving ? 'Сохранение...' : 'Сохранить' }}
      </button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
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
  technique: '',
  common_mistakes: '',
  easier_modification: '',
  harder_modification: '',
  rest_seconds: 60,
})

onMounted(async () => {
  if (isEdit) {
    try {
      const e = await api.getAdminExercise(Number(props.id))
      Object.assign(form, e)
    } catch {
      error.value = 'Не удалось загрузить'
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
      await api.updateAdminExercise(Number(props.id), { ...form })
    } else {
      await api.createAdminExercise({ ...form })
    }
    router.push('/admin/exercises')
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
.field input, .field textarea {
  padding: 10px 12px; border-radius: 10px; border: 1px solid var(--hint-color);
  background: var(--secondary-bg); color: var(--text-color); font-size: 15px;
}
.field textarea { resize: vertical; }
.error-msg { color: #ff3b30; font-size: 14px; text-align: center; }
.btn-primary {
  padding: 14px; border-radius: 12px; border: none;
  background: var(--button-color); color: var(--button-text-color);
  font-size: 16px; font-weight: 600; cursor: pointer;
}
.btn-primary:disabled { opacity: 0.5; }
</style>
