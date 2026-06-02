<template>
  <div class="admin-page">
    <h1 class="page-title">{{ isEdit ? 'Редактировать план питания' : 'Новый план питания' }}</h1>

    <div v-if="loading" class="loading">Загрузка...</div>

    <form v-else class="admin-form" @submit.prevent="save">
      <fieldset class="admin-section">
        <legend class="admin-section-title">Основные данные</legend>
        <div class="field">
          <label>Изображение</label>
          <ImageUploader v-model="form.image_media_id" reference-type="meal_plan_image" />
        </div>
        <div class="field">
          <label>Название *</label>
          <input v-model="form.name" required />
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
            <input v-model.number="form.day_number" type="number" min="1" max="365" />
          </div>
        </div>
      </fieldset>

      <fieldset class="admin-section">
        <legend class="admin-section-title">Питательная ценность</legend>
        <div class="row">
          <div class="field">
            <label>Калории</label>
            <input v-model.number="form.calories" type="number" min="0" max="10000" />
          </div>
          <div class="field">
            <label>Белки (г)</label>
            <input v-model.number="form.protein" type="number" step="0.1" min="0" max="1000" />
          </div>
        </div>
        <div class="row">
          <div class="field">
            <label>Жиры (г)</label>
            <input v-model.number="form.fat" type="number" step="0.1" min="0" max="1000" />
          </div>
          <div class="field">
            <label>Углеводы (г)</label>
            <input v-model.number="form.carbs" type="number" step="0.1" min="0" max="1000" />
          </div>
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
        <div class="row">
          <div class="field">
            <label>Порядок</label>
            <input v-model.number="form.sort_order" type="number" min="0" max="9999" />
          </div>
          <div class="field">
            <label>Активен</label>
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
import ImageUploader from '../../components/ImageUploader.vue'

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
  image_media_id: null as number | null,
  access_tier: 'paid' as 'free' | 'trial' | 'paid',
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
    router.replace('/admin/content')
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
.hint { font-size: 12px; color: var(--hint-color); margin-top: 4px; line-height: 1.4; }
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
