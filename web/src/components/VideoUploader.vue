<template>
  <div class="video-uploader">
    <input
      ref="fileInput"
      type="file"
      accept="video/mp4"
      style="display: none"
      @change="handleFileChange"
    />

    <!-- Preview when we already have a media_id -->
    <div v-if="modelValue && previewUrl && !uploading" class="preview-block">
      <video :src="previewUrl" controls preload="metadata" class="preview-video"></video>
      <div class="preview-actions">
        <button type="button" class="btn-link" @click="pickFile">Заменить</button>
        <button type="button" class="btn-link danger" @click="remove">Удалить</button>
      </div>
    </div>

    <!-- Progress while uploading -->
    <div v-else-if="uploading" class="upload-progress">
      <div class="progress-track">
        <div class="progress-bar" :style="{ width: percent + '%' }"></div>
      </div>
      <span class="progress-label">{{ percent }}% · {{ uploadingFileName }}</span>
      <button type="button" class="btn-link" @click="cancel">Отменить</button>
    </div>

    <!-- Empty state -->
    <button v-else type="button" class="dropzone" @click="pickFile">
      <span class="dropzone-icon">🎥</span>
      <span class="dropzone-title">Загрузить видео</span>
      <span class="dropzone-hint">MP4, до 100 МБ. Хранится в R2.</span>
    </button>

    <div v-if="errorMsg" class="error-msg">{{ errorMsg }}</div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useMediaUpload } from '../composables/useMediaUpload'
import { api } from '../api'

const props = defineProps<{
  modelValue: number | null
  referenceType?: string
}>()
const emit = defineEmits<{
  'update:modelValue': [value: number | null]
}>()

const fileInput = ref<HTMLInputElement | null>(null)
const { uploading, progress, error: uploadErr, upload } = useMediaUpload()
const uploadingFileName = ref('')
const previewUrl = ref<string | null>(null)
const errorMsg = ref('')

const percent = ref(0)
watch(progress, (p) => { percent.value = Math.round(p * 100) })

const MAX_SIZE = 100 * 1024 * 1024 // 100 MB — matches media_service.go allowlist

async function loadPreview(mediaID: number) {
  try {
    const r = await api.getMediaURL(mediaID)
    previewUrl.value = r.url
  } catch {
    previewUrl.value = null
  }
}

watch(() => props.modelValue, (id) => {
  if (id) {
    loadPreview(id)
  } else {
    previewUrl.value = null
  }
})

onMounted(() => {
  if (props.modelValue) loadPreview(props.modelValue)
})

function pickFile() {
  errorMsg.value = ''
  fileInput.value?.click()
}

async function handleFileChange(e: Event) {
  const inputEl = e.target as HTMLInputElement
  const file = inputEl.files?.[0]
  inputEl.value = '' // allow re-picking same file
  if (!file) return

  if (file.type !== 'video/mp4') {
    errorMsg.value = 'Только MP4-видео'
    return
  }
  if (file.size > MAX_SIZE) {
    errorMsg.value = `Файл слишком большой (${Math.round(file.size / 1024 / 1024)} МБ). Максимум 100 МБ.`
    return
  }

  uploadingFileName.value = file.name
  errorMsg.value = ''
  try {
    const result = await upload(file, {
      reference_type: props.referenceType ?? 'content_video',
      // Private bucket (same as avatars) — server-side GetWorkout / GetSession
      // resolves video_media_id to a presigned URL at view time.
      is_public: false,
    })
    emit('update:modelValue', result.media_id)
  } catch (e: any) {
    errorMsg.value = e.message || uploadErr.value || 'Ошибка загрузки'
  } finally {
    uploadingFileName.value = ''
  }
}

function remove() {
  emit('update:modelValue', null)
  previewUrl.value = null
  errorMsg.value = ''
}

function cancel() {
  // We can't actually abort the in-flight XHR (useMediaUpload doesn't expose
  // an AbortController). For now this only clears the local progress state on
  // the next upload-complete. Best effort — the upload will finish in the
  // background. A future refactor can wire AbortController through.
  errorMsg.value = 'Отмена не поддерживается. Дождитесь окончания загрузки.'
}
</script>

<style scoped>
.video-uploader {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.dropzone {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 24px 16px;
  border: 1.5px dashed var(--hint-color);
  border-radius: 12px;
  background: var(--secondary-bg);
  color: var(--text-color);
  cursor: pointer;
  font-family: inherit;
  -webkit-tap-highlight-color: transparent;
  transition: border-color 0.15s, transform 0.1s;
}

.dropzone:hover {
  border-color: var(--button-color);
}

.dropzone:active {
  transform: scale(0.99);
}

.dropzone-icon {
  font-size: 32px;
  line-height: 1;
}

.dropzone-title {
  font-size: 15px;
  font-weight: 600;
}

.dropzone-hint {
  font-size: 12px;
  color: var(--hint-color);
}

.preview-block {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.preview-video {
  width: 100%;
  max-height: 280px;
  border-radius: 12px;
  background: #000;
}

.preview-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

.btn-link {
  background: none;
  border: none;
  color: var(--button-color);
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  padding: 4px 8px;
}

.btn-link.danger {
  color: #ff3b30;
}

.upload-progress {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 14px;
  border-radius: 12px;
  background: var(--secondary-bg);
}

.progress-track {
  height: 6px;
  border-radius: 3px;
  background: var(--bg-color);
  overflow: hidden;
}

.progress-bar {
  height: 100%;
  background: var(--button-color);
  transition: width 0.15s linear;
}

.progress-label {
  font-size: 13px;
  color: var(--hint-color);
}

.error-msg {
  color: #ff3b30;
  font-size: 13px;
  padding: 2px 4px;
}
</style>
