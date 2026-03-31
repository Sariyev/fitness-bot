package repository

import (
	"context"
	"fitness-bot/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetByTelegramID(ctx context.Context, telegramID int64) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	CreateProfile(ctx context.Context, profile *models.UserProfile) error
	GetProfileByUserID(ctx context.Context, userID int64) (*models.UserProfile, error)
	UpdateProfile(ctx context.Context, profile *models.UserProfile) error
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

type PaymentRepository interface {
	CreatePayment(ctx context.Context, p *models.Payment) error
	UpdatePayment(ctx context.Context, p *models.Payment) error
	GetPaymentByID(ctx context.Context, id int64) (*models.Payment, error)
	CreatePendingPayment(ctx context.Context, userID int64, amountKZT int, provider string) (int64, error)
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

type ProgramRepository interface {
	ListPrograms(ctx context.Context, format, goal, level string) ([]models.Program, error)
	GetProgramByID(ctx context.Context, id int) (*models.Program, error)
	CreateProgram(ctx context.Context, p *models.Program) error
	UpdateProgram(ctx context.Context, p *models.Program) error
	EnrollUser(ctx context.Context, userID int64, programID int) error
	GetActiveEnrollment(ctx context.Context, userID int64) (*models.UserProgramEnrollment, error)
	ListUserEnrollments(ctx context.Context, userID int64) ([]models.UserProgramEnrollment, error)
}

type WorkoutRepository interface {
	ListWorkouts(ctx context.Context, format, goal, level string) ([]models.Workout, error)
	GetWorkoutByID(ctx context.Context, id int) (*models.Workout, error)
	ListByProgram(ctx context.Context, programID int) ([]models.Workout, error)
	CreateWorkout(ctx context.Context, w *models.Workout) error
	UpdateWorkout(ctx context.Context, w *models.Workout) error
	ListExercises(ctx context.Context, workoutID int) ([]models.WorkoutExercise, error)
	AddExercise(ctx context.Context, we *models.WorkoutExercise) error
}

type ExerciseRepository interface {
	List(ctx context.Context) ([]models.Exercise, error)
	GetByID(ctx context.Context, id int) (*models.Exercise, error)
	Create(ctx context.Context, e *models.Exercise) error
	Update(ctx context.Context, e *models.Exercise) error
}

type RehabRepository interface {
	ListCourses(ctx context.Context, category string) ([]models.RehabCourse, error)
	GetCourseByID(ctx context.Context, id int) (*models.RehabCourse, error)
	CreateCourse(ctx context.Context, c *models.RehabCourse) error
	UpdateCourse(ctx context.Context, c *models.RehabCourse) error
	ListSessions(ctx context.Context, courseID int) ([]models.RehabSession, error)
	GetSessionByID(ctx context.Context, id int) (*models.RehabSession, error)
	CreateSession(ctx context.Context, s *models.RehabSession) error
	UpdateSession(ctx context.Context, s *models.RehabSession) error
	CreateProgress(ctx context.Context, p *models.UserRehabProgress) error
	ListUserProgress(ctx context.Context, userID int64, courseID int) ([]models.UserRehabProgress, error)
}

type NutritionRepository interface {
	ListPlans(ctx context.Context, goal string) ([]models.MealPlan, error)
	GetPlanByID(ctx context.Context, id int) (*models.MealPlan, error)
	CreatePlan(ctx context.Context, p *models.MealPlan) error
	UpdatePlan(ctx context.Context, p *models.MealPlan) error
	ListMeals(ctx context.Context, planID int) ([]models.Meal, error)
	CreateMeal(ctx context.Context, m *models.Meal) error
	UpdateMeal(ctx context.Context, m *models.Meal) error
}

type FoodLogRepository interface {
	Create(ctx context.Context, entry *models.FoodLogEntry) error
	Delete(ctx context.Context, userID int64, id int64) error
	ListByDate(ctx context.Context, userID int64, date string) ([]models.FoodLogEntry, error)
	GetDailySummary(ctx context.Context, userID int64, date string) (calories int, protein, fat, carbs float64, err error)
}

type ProgressRepository interface {
	Create(ctx context.Context, entry *models.ProgressEntry) error
	ListByUser(ctx context.Context, userID int64) ([]models.ProgressEntry, error)
	GetWeightHistory(ctx context.Context, userID int64) ([]models.WeightPoint, error)
}

type AchievementRepository interface {
	ListAll(ctx context.Context) ([]models.Achievement, error)
	ListByUser(ctx context.Context, userID int64) ([]models.UserAchievement, error)
	Unlock(ctx context.Context, userID int64, achievementID int) error
}

type DailyCompletionRepository interface {
	Create(ctx context.Context, c *models.DailyCompletion) error
	ListByDate(ctx context.Context, userID int64, date string) ([]models.DailyCompletion, error)
	GetStreak(ctx context.Context, userID int64) (current int, longest int, err error)
	GetCalendar(ctx context.Context, userID int64, year, month int) ([]string, error)
}
