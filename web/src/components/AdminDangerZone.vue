<template>
  <fieldset class="admin-section danger-zone">
    <legend class="admin-section-title danger-title">Опасная зона</legend>
    <div class="danger-row">
      <button
        v-if="canDeactivate"
        type="button"
        class="btn-deactivate"
        :disabled="working"
        @click="onDeactivate"
      >
        {{ isActive ? 'Снять с публикации' : 'Опубликовать' }}
      </button>
      <button
        type="button"
        class="btn-delete"
        :disabled="working"
        @click="onDelete"
      >
        {{ working ? 'Удаление...' : 'Удалить' }}
      </button>
    </div>
    <p class="hint">
      «Снять с публикации» делает контент невидимым для пользователей, но сохраняет данные.
      «Удалить» — необратимо. Удаление родителя удаляет и его дочерние элементы.
    </p>
  </fieldset>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const props = defineProps<{
  label: string                   // e.g. "программу", "тренировку" — used in confirm prompt
  isActive?: boolean              // undefined = entity has no is_active field; hides toggle
  canDeactivate?: boolean         // explicit override; defaults to (isActive !== undefined)
}>()

const emit = defineEmits<{
  delete: []
  toggleActive: [next: boolean]
}>()

const working = ref(false)

const canDeactivate = props.canDeactivate ?? (props.isActive !== undefined)

function onDelete() {
  if (!window.confirm(`Удалить ${props.label}? Это действие необратимо.`)) return
  working.value = true
  emit('delete')
}

function onDeactivate() {
  emit('toggleActive', !props.isActive)
}

// Parent should call .reset() after async finishes; expose for cases where save fails.
defineExpose({
  reset: () => { working.value = false },
})
</script>

<style scoped>
.danger-zone {
  border: 1px solid rgba(255, 59, 48, 0.3);
  border-radius: 12px;
  padding: 14px;
  margin-top: 18px;
}

.danger-title {
  color: #ff3b30;
  font-size: 14px;
}

.danger-row {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  margin-bottom: 8px;
}

.btn-deactivate,
.btn-delete {
  flex: 1;
  min-width: 140px;
  padding: 10px 14px;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  border: 1px solid;
  background: var(--secondary-bg);
}

.btn-deactivate {
  border-color: var(--hint-color);
  color: var(--text-color);
}

.btn-delete {
  border-color: #ff3b30;
  color: #ff3b30;
}

.btn-delete:active {
  background: rgba(255, 59, 48, 0.08);
}

.btn-deactivate:disabled,
.btn-delete:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.hint {
  font-size: 12px;
  color: var(--hint-color);
  margin: 0;
  line-height: 1.4;
}
</style>
