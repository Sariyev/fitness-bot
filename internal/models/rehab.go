package models

import "time"

type RehabCourse struct {
	ID          int       `json:"id"`
	Slug        string    `json:"slug"`
	Category    string    `json:"category"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Warnings    string    `json:"warnings"`
	IsActive    bool      `json:"is_active"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RehabSession struct {
	ID              int       `json:"id"`
	CourseID        int       `json:"course_id"`
	DayNumber       int       `json:"day_number"`
	Stage           int       `json:"stage"`
	VideoURL        string    `json:"video_url"`
	DurationMinutes int       `json:"duration_minutes"`
	Description     string    `json:"description"`
	SortOrder       int       `json:"sort_order"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type UserRehabProgress struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	CourseID    int       `json:"course_id"`
	SessionID   int       `json:"session_id"`
	DayNumber   int       `json:"day_number"`
	CompletedAt time.Time `json:"completed_at"`
	PainLevel   int       `json:"pain_level"`
	Comment     string    `json:"comment"`
}
