<template>
  <div class="skeleton-card">
    <div class="skeleton-line skeleton-icon"></div>
    <div class="skeleton-body">
      <div class="skeleton-line skeleton-title"></div>
      <div class="skeleton-line skeleton-text"></div>
    </div>
  </div>
</template>

<style scoped>
.skeleton-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: var(--secondary-bg);
  border-radius: 12px;
}

.skeleton-line {
  background: linear-gradient(
    90deg,
    var(--secondary-bg) 25%,
    color-mix(in srgb, var(--hint-color) 15%, transparent) 50%,
    var(--secondary-bg) 75%
  );
  background-size: 200% 100%;
  animation: shimmer 1.5s ease-in-out infinite;
  border-radius: 6px;
}

.skeleton-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  flex-shrink: 0;
}

.skeleton-body {
  flex: 1;
}

.skeleton-title {
  height: 16px;
  width: 60%;
  margin-bottom: 8px;
}

.skeleton-text {
  height: 12px;
  width: 80%;
}

@keyframes shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}
</style>
