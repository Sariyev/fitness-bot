<template>
  <div class="admin-page">
    <h1 class="page-title">{{ isEdit ? 'Редактировать упражнение' : 'Новое упражнение' }}</h1>

    <div v-if="loading" class="loading">Загрузка...</div>

    <form v-else class="admin-form" @submit.prevent="save">
      <fieldset class="admin-section">
        <legend class="admin-section-title">Основные данные</legend>
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
      </fieldset>

      <fieldset class="admin-section">
        <legend class="admin-section-title">Варианты выполнения</legend>
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
          <input v-model.number="form.rest_seconds" type="number" min="0" max="600" />
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
        label="упражнение"
        @delete="onDelete"
      />
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../api'
import AdminDangerZone from '../../components/AdminDangerZone.vue'

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
    router.replace('/admin/exercises')
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
    await api.deleteAdminExercise(Number(props.id))
    router.replace('/admin/exercises')
  } catch (e: any) {
    error.value = e.message || 'Не удалось удалить'
    dangerZone.value?.reset()
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
