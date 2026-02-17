package models

import "time"

type Questionnaire struct {
	ID          int        `json:"id"`
	Slug        string     `json:"slug"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	IsActive    bool       `json:"is_active"`
	SortOrder   int        `json:"sort_order"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Questions   []Question `json:"questions,omitempty"`
}

type Question struct {
	ID              int              `json:"id"`
	QuestionnaireID int              `json:"questionnaire_id"`
	Text            string           `json:"text"`
	QuestionType    string           `json:"question_type"`
	SortOrder       int              `json:"sort_order"`
	IsRequired      bool             `json:"is_required"`
	Metadata        map[string]any   `json:"metadata"`
	Options         []QuestionOption `json:"options,omitempty"`
	CreatedAt       time.Time        `json:"created_at"`
}

type QuestionOption struct {
	ID         int       `json:"id"`
	QuestionID int       `json:"question_id"`
	Text       string    `json:"text"`
	Value      string    `json:"value"`
	SortOrder  int       `json:"sort_order"`
	CreatedAt  time.Time `json:"created_at"`
}

type QuestionnaireSubmission struct {
	ID              int64      `json:"id"`
	UserID          int64      `json:"user_id"`
	QuestionnaireID int        `json:"questionnaire_id"`
	CompletedAt     *time.Time `json:"completed_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

type QuestionnaireAnswer struct {
	ID           int64     `json:"id"`
	SubmissionID int64     `json:"submission_id"`
	QuestionID   int       `json:"question_id"`
	AnswerText   *string   `json:"answer_text"`
	AnswerValue  *string   `json:"answer_value"`
	AnswerValues []string  `json:"answer_values"`
	CreatedAt    time.Time `json:"created_at"`
}
