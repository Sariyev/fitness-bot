<template>
  <div class="admin-page">
    <button class="back-btn" @click="router.push('/admin/content')">← Назад</button>
    <h1 class="page-title">{{ isEdit ? 'Редактировать план питания' : 'Новый план питания' }}</h1>

    <div v-if="loading" class="loading">Загрузка...</div>

    <form v-else class="form" @submit.prevent="save">
      <div class="field">
        <label>Наз��ание *</label>
        <input v-model="form.name" required />
      </div>
      <div class="field">
        <label>Slug</label>
        <input v-model="form.slug" placeholder="auto-from-name" />
      </div>
      <div class="row">
        <div class="field">
          <label>Цель</label>
          <select v-model="form.goal">
            <option value="">—</option>
            <option value="weight_loss">Похудение</option>
            <option value="muscle_gain">Набор массы</option>
            <option value="maintenance">Поддержание</option>
          </select>
        </div>
        <div class="field">
          <label>День</label>
          <input v-model.number="form.day_number" type="number" min="1" />
        </div>
      </div>
      <div class="row">
        <div class="field">
          <label>Калории</label>
          <input v-model.number="form.calories" type="number" />
        </div>
        <div class="field">
          <label>Белки (г)</label>
          <input v-model.number="form.protein" type="number" step="0.1" />
        </div>
      </div>
      <div class="row">
        <div class="field">
          <label>Жиры (г)</label>
          <input v-model.number="form.fat" type="number" step="0.1" />
        </div>
        <div class="field">
          <label>Углеводы (г)</label>
          <input v-model.number="form.carbs" type="number" step="0.1" />
        </div>
      </div>
      <div class="row">
        <div class="field">
          <label>Порядок</label>
          <input v-model.number="form.sort_order" type="number" />
        </div>
        <div class="field">
          <label>Активен</label>
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

    <!-- Meals section for edit mode -->
    <div v-if="isEdit && !loading" class="meals-section">
      <div class="section-header">
        <h2 class="section-title">Блюда ({{ meals.length }})</h2>
        <button class="add-btn" @click="router.push(`/admin/meals/new?plan_id=${props.id}`)">+</button>
      </div>
      <div v-for="m in meals" :key="m.id" class="content-card" @click="router.push(`/admin/meals/${m.id}`)">
        <div class="content-main">
          <span class="content-name">{{ m.name }}</span>
          <span class="content-meta">{{ mealTypeLabel(m.meal_type) }} | {{ m.calories }} ккал</span>
        </div>
        <span class="arrow">→</span>
      </div>
      <div v-if="meals.length === 0" class="empty">Нет блюд</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../api'
import type { Meal } from '../../types'

const props = defineProps<{ id?: string }>()
const router = useRouter()
const isEdit = !!props.id
const loading = ref(isEdit)
const saving = ref(false)
const error = ref('')
const meals = ref<Meal[]>([])

const form = reactive({
  name: '',
  slug: '',
  goal: '',
  day_number: 1,
  calories: 0,
  protein: 0,
  fat: 0,
  carbs: 0,
  sort_order: 0,
  is_active: true,
})

function mealTypeLabel(t: string) {
  const map: Record<string, string> = { breakfast: 'Завтрак', lunch: 'Обед', dinner: 'Ужин', snack: 'Перекус' }
  return map[t] || t
}

onMounted(async () => {
  if (isEdit) {
    try {
      const data = await api.getAdminMealPlan(Number(props.id))
      Object.assign(form, data.plan)
      meals.value = data.meals || []
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
      await api.updateAdminMealPlan(Number(props.id), { ...form })
    } else {
      await api.createAdminMealPlan({ ...form })
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
.empty { text-align: center; color: var(--hint-color); padding: 20px; }
.form { display: flex; flex-direction: column; gap: 14px; }
.field { display: flex; flex-direction: column; gap: 4px; }
.field label { font-size: 13px; color: var(--hint-color); font-weight: 500; }
.field input, .field select, .field textarea {
  padding: 10px 12px; border-radius: 10px; border: 1px solid var(--hint-color);
  background: var(--secondary-bg); color: var(--text-color); font-size: 15px;
}
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
.meals-section { margin-top: 24px; }
.section-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px; }
.section-title { font-size: 16px; font-weight: 600; color: var(--hint-color); }
.add-btn {
  width: 32px; height: 32px; border-radius: 50%; border: none;
  background: var(--button-color); color: var(--button-text-color);
  font-size: 20px; font-weight: 700; cursor: pointer;
  display: flex; align-items: center; justify-content: center;
}
.content-card {
  background: var(--secondary-bg); border-radius: 12px; padding: 12px 16px; margin-bottom: 8px;
  display: flex; justify-content: space-between; align-items: center; cursor: pointer;
}
.content-main { display: flex; flex-direction: column; }
.content-name { font-weight: 500; font-size: 14px; }
.content-meta { color: var(--hint-color); font-size: 12px; margin-top: 2px; }
.arrow { font-size: 18px; color: var(--hint-color); }
</style>
