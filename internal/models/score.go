package models

import "time"

type UserScore struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	ScoreType     string    `json:"score_type"`
	ReferenceType string    `json:"reference_type"`
	ReferenceID   int       `json:"reference_id"`
	Score         int       `json:"score"`
	Comment       string    `json:"comment"`
	Tags          []string  `json:"tags"`
	CreatedAt     time.Time `json:"created_at"`
}

// ReviewSummary holds aggregated review stats for a reference.
type ReviewSummary struct {
	ReferenceType string  `json:"reference_type"`
	ReferenceID   int     `json:"reference_id"`
	AverageScore  float64 `json:"average_score"`
	TotalReviews  int     `json:"total_reviews"`
}

// Preset tag sets by reference type.
var ReviewTags = map[string][]string{
	"workout": {
		"Эффективно", "Понятные инструкции", "Хорошая нагрузка",
		"Слишком сложно", "Слишком легко", "Отличное видео",
	},
	"rehab_session": {
		"Помогает", "Понятные упражнения", "Хорошее видео",
		"Слишком сложно", "Уменьшает боль", "Нужно больше объяснений",
	},
	"lesson": {
		"Полезно", "Понятно", "Интересно",
		"Мало информации", "Хочу ещё", "Скучно",
	},
	"meal_plan": {
		"Вкусные рецепты", "Легко готовить", "Разнообразно",
		"Сложные ингредиенты", "Мало вариантов", "Полезно",
	},
	"bot": {
		"Удобный", "Понятный", "Полезный контент",
		"Хороший тренер", "Быстрый", "Много функций",
	},
}
