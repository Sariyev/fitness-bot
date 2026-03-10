package models

import "time"

type Program struct {
	ID            int       `json:"id"`
	Slug          string    `json:"slug"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Goal          string    `json:"goal"`
	Format        string    `json:"format"`
	Level         string    `json:"level"`
	DurationWeeks int       `json:"duration_weeks"`
	IsActive      bool      `json:"is_active"`
	SortOrder     int       `json:"sort_order"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Workout struct {
	ID              int       `json:"id"`
	ProgramID       *int      `json:"program_id"`
	Slug            string    `json:"slug"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Goal            string    `json:"goal"`
	Format          string    `json:"format"`
	Level           string    `json:"level"`
	DurationMinutes int       `json:"duration_minutes"`
	Equipment       []string  `json:"equipment"`
	ExpectedResult  string    `json:"expected_result"`
	VideoURL        string    `json:"video_url"`
	SortOrder       int       `json:"sort_order"`
	WeekNumber      *int      `json:"week_number"`
	DayNumber       *int      `json:"day_number"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type UserProgramEnrollment struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	ProgramID   int       `json:"program_id"`
	StartedAt   time.Time `json:"started_at"`
	CurrentWeek int       `json:"current_week"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
