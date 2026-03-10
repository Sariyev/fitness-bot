package models

import (
	"encoding/json"
	"time"
)

type Achievement struct {
	ID          int             `json:"id"`
	Slug        string          `json:"slug"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Icon        string          `json:"icon"`
	Criteria    json.RawMessage `json:"criteria"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type UserAchievement struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	AchievementID int       `json:"achievement_id"`
	EarnedAt      time.Time `json:"earned_at"`
}
