package models

import (
	"encoding/json"
	"time"
)

type ConversationState struct {
	ID         int64           `json:"id"`
	TelegramID int64           `json:"telegram_id"`
	State      string          `json:"state"`
	Data       json.RawMessage `json:"data"`
	ExpiresAt  time.Time       `json:"expires_at"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

const (
	StateIdle = "idle"

	// Registration
	StateRegAge          = "reg:age"
	StateRegHeight       = "reg:height"
	StateRegWeight       = "reg:weight"
	StateRegGender       = "reg:gender"
	StateRegFitnessLevel = "reg:fitness_level"
	StateRegGoal         = "reg:goal"

	// Questionnaire
	StateQuizInProgress = "quiz:in_progress"

	// Payment
	StatePayShowProduct  = "sub:show_product"
	StatePayAwaitConfirm = "sub:await_confirm"
	StatePayProcessing   = "sub:processing"

	// Module browsing
	StateModBrowse         = "mod:browse"
	StateModChooseCategory = "mod:choose_category"
	StateModChooseLesson   = "mod:choose_lesson"
	StateModViewLesson     = "mod:view_lesson"

	// Score collection
	StateScoreAwaitRating  = "score:await_rating"
	StateScoreAwaitComment = "score:await_comment"
)

type RegistrationData struct {
	Age          int      `json:"age,omitempty"`
	HeightCm     int      `json:"height_cm,omitempty"`
	WeightKg     float64  `json:"weight_kg,omitempty"`
	Gender       string   `json:"gender,omitempty"`
	FitnessLevel string   `json:"fitness_level,omitempty"`
	Goals        []string `json:"goals,omitempty"`
	GoalMsgID    int      `json:"goal_msg_id,omitempty"`
}

type QuizFlowData struct {
	SubmissionID    int64 `json:"submission_id"`
	QuestionnaireID int   `json:"questionnaire_id"`
	CurrentQuestionIdx int `json:"current_question_idx"`
	TotalQuestions  int   `json:"total_questions"`
}

type PaymentFlowData struct {
	MessageID int `json:"message_id,omitempty"`
}

type ModuleBrowseData struct {
	ModuleID   int `json:"module_id,omitempty"`
	CategoryID int `json:"category_id,omitempty"`
	LessonID   int `json:"lesson_id,omitempty"`
	ContentIdx int `json:"content_idx,omitempty"`
}

type ScoreFlowData struct {
	ScoreType     string `json:"score_type"`
	ReferenceType string `json:"reference_type"`
	ReferenceID   int    `json:"reference_id"`
	Score         int    `json:"score,omitempty"`
}
