<template>
  <div class="welcome-slide">
    <div class="welcome-emoji animate-bounce-in">👋</div>
    <h1 class="animate-fade-up">Привет! Я Андрей</h1>
    <p class="subtitle animate-fade-up delay-1">Тренер-реабилитолог. Помогу восстановить здоровье и избавиться от боли</p>

    <div class="features">
      <div class="feature animate-fade-up delay-2">
        <span class="feature-icon">🏥</span>
        <div>
          <strong>ЛФК-программы</strong>
          <p>Персональные курсы при грыжах, протрузиях, сколиозе</p>
        </div>
      </div>
      <div class="feature animate-fade-up delay-3">
        <span class="feature-icon">📋</span>
        <div>
          <strong>Безопасные упражнения</strong>
          <p>Подобраны под твой диагноз и уровень боли</p>
        </div>
      </div>
      <div class="feature animate-fade-up delay-4">
        <span class="feature-icon">💪</span>
        <div>
          <strong>Тренировки</strong>
          <p>Фитнес-программы для поддержания формы и здоровья</p>
        </div>
      </div>
    </div>

    <p class="cta animate-fade-up delay-4">Жми НАЧАТЬ 👇</p>
  </div>
</template>

<style scoped>
.welcome-slide {
  text-align: center;
  padding: 20px 0;
}

.welcome-emoji {
  font-size: 72px;
  margin-bottom: 16px;
}

h1 {
  font-size: 26px;
  margin-bottom: 4px;
}

.subtitle {
  color: var(--hint-color);
  font-size: 16px;
  margin-bottom: 32px;
}

.features {
  text-align: left;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.feature {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px 16px;
  background: var(--secondary-bg);
  border-radius: 14px;
}

.feature-icon {
  font-size: 28px;
  flex-shrink: 0;
}

.feature strong {
  display: block;
  font-size: 15px;
  margin-bottom: 2px;
}

.feature p {
  font-size: 13px;
  color: var(--hint-color);
  margin: 0;
}

.cta {
  margin-top: 28px;
  font-size: 22px;
  font-weight: 700;
  color: var(--button-color);
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 0.6; }
  50% { opacity: 1; }
}

.animate-bounce-in {
  animation: bounceIn 0.6s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
}

.animate-fade-up {
  opacity: 0;
  animation: fadeUp 0.5s ease forwards;
}

.delay-1 { animation-delay: 0.1s; }
.delay-2 { animation-delay: 0.2s; }
.delay-3 { animation-delay: 0.3s; }
.delay-4 { animation-delay: 0.4s; }

@keyframes bounceIn {
  0% { transform: scale(0.3); opacity: 0; }
  50% { transform: scale(1.05); }
  100% { transform: scale(1); opacity: 1; }
}

@keyframes fadeUp {
  from { opacity: 0; transform: translateY(16px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
