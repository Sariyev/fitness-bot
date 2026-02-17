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
	CreatedAt     time.Time `json:"created_at"`
}
