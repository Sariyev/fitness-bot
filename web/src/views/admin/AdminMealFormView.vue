<template>
  <div class="admin-page">
    <button class="back-btn" @click="goBack">← Назад</button>
    <h1 class="page-title">{{ isEdit ? 'Редактировать блюдо' : 'Новое блюдо' }}</h1>

    <div v-if="loading" class="loading">Загрузка...</div>

    <form v-else class="form" @submit.prevent="save">
      <div class="field">
        <label>План питания</label>
        <select v-model.number="form.meal_plan_id" required>
          <option :value="0">— Выбрать —</option>
          <option v-for="p in plans" :key="p.id" :value="p.id">{{ p.name }}</option>
        </select>
      </div>
      <div class="field">
        <label>Тип приёма</label>
        <select v-model="form.meal_type">
          <option value="breakfast">Завтрак</option>
          <option value="lunch">Обед</option>
          <option value="dinner">Ужин</option>
          <option value="snack">Перекус</option>
        </select>
      </div>
      <div class="field">
        <label>Название *</label>
        <input v-model="form.name" required />
      </div>
      <div class="field">
        <label>Рецепт</label>
        <textarea v-model="form.recipe" rows="4"></textarea>
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
      <div class="field">
        <label>Альтернативы</label>
        <textarea v-model="form.alternatives" rows="2"></textarea>
      </div>
      <div class="field">
        <label>Порядок</label>
        <input v-model.number="form.sort_order" type="number" />
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
import { useRouter, useRoute } from 'vue-router'
import { api } from '../../api'
import type { MealPlan } from '../../types'

const props = defineProps<{ id?: string }>()
const router = useRouter()
const route = useRoute()
const isEdit = !!props.id
const loading = ref(true)
const saving = ref(false)
const error = ref('')
const plans = ref<MealPlan[]>([])

const form = reactive({
  meal_plan_id: 0,
  meal_type: 'breakfast',
  name: '',
  recipe: '',
  calories: 0,
  protein: 0,
  fat: 0,
  carbs: 0,
  alternatives: '',
  sort_order: 0,
})

function goBack() {
  if (form.meal_plan_id) {
    router.push(`/admin/meal-plans/${form.meal_plan_id}`)
  } else {
    router.push('/admin/content')
  }
}

onMounted(async () => {
  try {
    plans.value = (await api.getAdminMealPlans()) || []

    if (isEdit) {
      const m = await api.getAdminMeal(Number(props.id))
      Object.assign(form, m)
    } else {
      const planId = Number(route.query.plan_id)
      if (planId) form.meal_plan_id = planId
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
      await api.updateAdminMeal(Number(props.id), { ...form })
    } else {
      await api.createAdminMeal({ ...form })
    }
    if (form.meal_plan_id) {
      router.push(`/admin/meal-plans/${form.meal_plan_id}`)
    } else {
      router.push('/admin/content')
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
.error-msg { color: #ff3b30; font-size: 14px; text-align: center; }
.btn-primary {
  padding: 14px; border-radius: 12px; border: none;
  background: var(--button-color); color: var(--button-text-color);
  font-size: 16px; font-weight: 600; cursor: pointer;
}
.btn-primary:disabled { opacity: 0.5; }
</style>
