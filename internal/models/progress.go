package models

import "time"

type ProgressEntry struct {
	ID           int64          `json:"id"`
	UserID       int64          `json:"user_id"`
	Date         string         `json:"date"`
	WeightKg     *float64       `json:"weight_kg"`
	Measurements map[string]any `json:"measurements"`
	PhotoURL     string         `json:"photo_url"`
	Wellbeing    string         `json:"wellbeing"`
	PainLevel    int            `json:"pain_level"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type DailyCompletion struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	EntityType string    `json:"entity_type"`
	EntityID   int       `json:"entity_id"`
	Date       string    `json:"date"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type WeightPoint struct {
	Date     string  `json:"date"`
	WeightKg float64 `json:"weight_kg"`
}
