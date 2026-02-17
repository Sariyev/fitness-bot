package repository

import (
	"context"
	"fitness-bot/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByTelegramID(ctx context.Context, telegramID int64) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	CreateProfile(ctx context.Context, profile *models.UserProfile) error
	GetProfileByUserID(ctx context.Context, userID int64) (*models.UserProfile, error)
}

type ConversationRepository interface {
	GetState(ctx context.Context, telegramID int64) (*models.ConversationState, error)
	UpsertState(ctx context.Context, state *models.ConversationState) error
	ClearState(ctx context.Context, telegramID int64) error
	CleanupExpired(ctx context.Context) error
}

type QuestionnaireRepository interface {
	GetBySlug(ctx context.Context, slug string) (*models.Questionnaire, error)
	GetByID(ctx context.Context, id int) (*models.Questionnaire, error)
	ListActive(ctx context.Context) ([]models.Questionnaire, error)
	Create(ctx context.Context, q *models.Questionnaire) error
	Update(ctx context.Context, q *models.Questionnaire) error

	GetQuestionsByQuestionnaireID(ctx context.Context, qID int) ([]models.Question, error)
	GetQuestionByID(ctx context.Context, id int) (*models.Question, error)
	CreateQuestion(ctx context.Context, q *models.Question) error
	UpdateQuestion(ctx context.Context, q *models.Question) error
	DeleteQuestion(ctx context.Context, id int) error

	GetOptionsByQuestionID(ctx context.Context, questionID int) ([]models.QuestionOption, error)
	CreateOption(ctx context.Context, o *models.QuestionOption) error
	DeleteOption(ctx context.Context, id int) error

	CreateSubmission(ctx context.Context, s *models.QuestionnaireSubmission) error
	CompleteSubmission(ctx context.Context, submissionID int64) error
	SaveAnswer(ctx context.Context, a *models.QuestionnaireAnswer) error
	GetSubmissionAnswers(ctx context.Context, submissionID int64) ([]models.QuestionnaireAnswer, error)
	GetUserSubmissions(ctx context.Context, userID int64, questionnaireID int) ([]models.QuestionnaireSubmission, error)
}

type SubscriptionRepository interface {
	ListActivePlans(ctx context.Context) ([]models.SubscriptionPlan, error)
	GetPlanByID(ctx context.Context, id int) (*models.SubscriptionPlan, error)
	CreatePlan(ctx context.Context, plan *models.SubscriptionPlan) error
	UpdatePlan(ctx context.Context, plan *models.SubscriptionPlan) error

	GetActiveSubscription(ctx context.Context, userID int64) (*models.Subscription, error)
	CreateSubscription(ctx context.Context, sub *models.Subscription) error
	ExpireSubscription(ctx context.Context, subID int64) error

	CreatePayment(ctx context.Context, p *models.Payment) error
	UpdatePayment(ctx context.Context, p *models.Payment) error
	GetPaymentByID(ctx context.Context, id int64) (*models.Payment, error)
}

type ModuleRepository interface {
	ListActiveModules(ctx context.Context) ([]models.Module, error)
	GetModuleByID(ctx context.Context, id int) (*models.Module, error)
	GetModuleBySlug(ctx context.Context, slug string) (*models.Module, error)
	CreateModule(ctx context.Context, m *models.Module) error
	UpdateModule(ctx context.Context, m *models.Module) error

	ListCategoriesByModule(ctx context.Context, moduleID int) ([]models.ModuleCategory, error)
	GetCategoryByID(ctx context.Context, id int) (*models.ModuleCategory, error)
	CreateCategory(ctx context.Context, c *models.ModuleCategory) error
	UpdateCategory(ctx context.Context, c *models.ModuleCategory) error

	ListLessonsByCategory(ctx context.Context, categoryID int) ([]models.Lesson, error)
	GetLessonByID(ctx context.Context, id int) (*models.Lesson, error)
	CreateLesson(ctx context.Context, l *models.Lesson) error
	UpdateLesson(ctx context.Context, l *models.Lesson) error

	ListContentByLesson(ctx context.Context, lessonID int) ([]models.LessonContent, error)
	GetContentByID(ctx context.Context, id int) (*models.LessonContent, error)
	CreateContent(ctx context.Context, c *models.LessonContent) error
	UpdateContent(ctx context.Context, c *models.LessonContent) error
	DeleteContent(ctx context.Context, id int) error
	UpdateTelegramFileID(ctx context.Context, contentID int, fileID string) error

	UpsertProgress(ctx context.Context, userID int64, lessonID int) error
	CompleteLesson(ctx context.Context, userID int64, lessonID int) error
	GetUserProgress(ctx context.Context, userID int64, lessonID int) (*models.UserLessonProgress, error)
	GetUserCategoryProgress(ctx context.Context, userID int64, categoryID int) (completed int, total int, err error)

	SelectCategory(ctx context.Context, userID int64, categoryID int) error
	GetUserSelections(ctx context.Context, userID int64) ([]models.UserModuleSelection, error)
}

type ScoreRepository interface {
	Create(ctx context.Context, s *models.UserScore) error
	GetByReference(ctx context.Context, userID int64, refType string, refID int) (*models.UserScore, error)
	ListByUser(ctx context.Context, userID int64) ([]models.UserScore, error)
}
