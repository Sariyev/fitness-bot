<template>
  <div class="content-block">
    <h3 v-if="content.title" class="content-title">{{ content.title }}</h3>

    <!-- Text content -->
    <div v-if="content.content_type === 'text'" class="text-content" v-html="content.body"></div>

    <!-- Video content -->
    <div v-else-if="content.content_type === 'video'" class="video-content">
      <div v-if="youtubeId" class="video-wrapper">
        <iframe
          :src="`https://www.youtube.com/embed/${youtubeId}`"
          frameborder="0"
          allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
          allowfullscreen
        ></iframe>
      </div>
      <video
        v-else-if="content.file_url || content.video_url"
        controls
        playsinline
        class="video-player"
      >
        <source :src="content.file_url || content.video_url" />
      </video>
      <p v-else class="video-unavailable">Видео доступно только в боте</p>
    </div>

    <!-- Image content -->
    <div v-else-if="content.content_type === 'image'" class="image-content">
      <img
        v-if="content.file_url"
        :src="content.file_url"
        :alt="content.title"
        loading="lazy"
      />
    </div>

    <!-- Document content -->
    <div v-else-if="content.content_type === 'document'" class="doc-content">
      <a v-if="content.file_url" :href="content.file_url" target="_blank" class="doc-link">
        {{ content.title || 'Скачать документ' }}
      </a>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { LessonContent } from '../types'

const props = defineProps<{
  content: LessonContent
}>()

const youtubeId = computed(() => {
  const url = props.content.video_url
  if (!url) return null

  // Match youtube.com/watch?v=ID or youtu.be/ID
  const match = url.match(
    /(?:youtube\.com\/(?:watch\?v=|embed\/)|youtu\.be\/)([a-zA-Z0-9_-]{11})/
  )
  return match ? match[1] : null
})
</script>

<style scoped>
.content-block {
  background: var(--secondary-bg);
  border-radius: 12px;
  padding: 16px;
  overflow: hidden;
}

.content-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 8px;
}

.text-content {
  font-size: 14px;
  line-height: 1.6;
  word-wrap: break-word;
}

.video-wrapper {
  position: relative;
  width: 100%;
  padding-bottom: 56.25%; /* 16:9 */
  border-radius: 8px;
  overflow: hidden;
}

.video-wrapper iframe {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}

.video-player {
  width: 100%;
  border-radius: 8px;
}

.video-unavailable {
  color: var(--hint-color);
  font-size: 13px;
  font-style: italic;
}

.image-content img {
  width: 100%;
  border-radius: 8px;
  display: block;
}

.doc-link {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: var(--link-color);
  text-decoration: none;
  font-size: 14px;
}

.doc-link::before {
  content: '📄';
}
</style>
