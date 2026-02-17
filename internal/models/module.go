package models

import "time"

type Module struct {
	ID                   int       `json:"id"`
	Slug                 string    `json:"slug"`
	Name                 string    `json:"name"`
	Description          string    `json:"description"`
	Icon                 string    `json:"icon"`
	RequiresSubscription bool      `json:"requires_subscription"`
	IsActive             bool      `json:"is_active"`
	SortOrder            int       `json:"sort_order"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type ModuleCategory struct {
	ID          int       `json:"id"`
	ModuleID    int       `json:"module_id"`
	Slug        string    `json:"slug"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	SortOrder   int       `json:"sort_order"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Lesson struct {
	ID          int       `json:"id"`
	CategoryID  int       `json:"category_id"`
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	SortOrder   int       `json:"sort_order"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type LessonContent struct {
	ID             int       `json:"id"`
	LessonID       int       `json:"lesson_id"`
	ContentType    string    `json:"content_type"`
	Title          string    `json:"title"`
	Body           string    `json:"body"`
	VideoURL       string    `json:"video_url"`
	TelegramFileID string    `json:"telegram_file_id"`
	FileURL        string    `json:"file_url"`
	SortOrder      int       `json:"sort_order"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type UserLessonProgress struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"user_id"`
	LessonID    int        `json:"lesson_id"`
	Status      string     `json:"status"`
	StartedAt   time.Time  `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

type UserModuleSelection struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	CategoryID int       `json:"category_id"`
	SelectedAt time.Time `json:"selected_at"`
}
