import { ref } from 'vue'
import { api } from '../api'

export interface UploadOptions {
  reference_type?: string
  reference_id?: number | null
  is_public?: boolean
  // onProgress receives a value in [0, 1].
  onProgress?: (fraction: number) => void
}

export interface UploadResult {
  media_id: number
  key: string
}

/**
 * Upload a single File to R2 via the backend's presigned-PUT flow.
 *
 * 1. POST /api/media/presign       -> get a presigned PUT URL + media_id
 * 2. PUT direct to R2 (no Go server in the byte path)
 * 3. POST /api/media/{id}/confirm  -> backend HEADs the object and flips
 *                                     confirmed=true
 *
 * Reactive `progress` (0..1) and `uploading` are exposed so views can
 * render their own spinners or progress bars.
 */
export function useMediaUpload() {
  const uploading = ref(false)
  const progress = ref(0)
  const error = ref<string | null>(null)

  async function upload(file: File, opts: UploadOptions = {}): Promise<UploadResult> {
    uploading.value = true
    progress.value = 0
    error.value = null

    try {
      // 1) presign
      const presigned = await api.presignMediaUpload({
        content_type: file.type,
        size_bytes: file.size,
        reference_type: opts.reference_type,
        reference_id: opts.reference_id ?? null,
        is_public: opts.is_public ?? false,
      })

      // 2) PUT direct to R2 with progress
      await new Promise<void>((resolve, reject) => {
        const xhr = new XMLHttpRequest()
        xhr.open('PUT', presigned.upload_url)
        xhr.setRequestHeader('Content-Type', file.type)
        xhr.upload.onprogress = (e) => {
          if (e.lengthComputable) {
            const frac = e.loaded / e.total
            progress.value = frac
            opts.onProgress?.(frac)
          }
        }
        xhr.onload = () => {
          if (xhr.status >= 200 && xhr.status < 300) {
            progress.value = 1
            opts.onProgress?.(1)
            resolve()
          } else {
            reject(new Error(`R2 upload failed: ${xhr.status} ${xhr.statusText}`))
          }
        }
        xhr.onerror = () => reject(new Error('R2 upload network error'))
        xhr.send(file)
      })

      // 3) confirm
      await api.confirmMediaUpload(presigned.media_id)

      return { media_id: presigned.media_id, key: presigned.key }
    } catch (e: any) {
      error.value = e.message || 'Upload failed'
      throw e
    } finally {
      uploading.value = false
    }
  }

  return { uploading, progress, error, upload }
}
