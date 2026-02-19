package models

import "time"

type User struct {
	ID           int64     `json:"id"`
	TelegramID   int64     `json:"telegram_id"`
	Username     string    `json:"username"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	LanguageCode string    `json:"language_code"`
	IsRegistered bool      `json:"is_registered"`
	IsPaid       bool      `json:"is_paid"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserProfile struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Gender       string    `json:"gender"`
	Age          int       `json:"age"`
	WeightKg     float64   `json:"weight_kg"`
	HeightCm     int       `json:"height_cm"`
	FitnessLevel string    `json:"fitness_level"`
	Goal         string    `json:"goal"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
