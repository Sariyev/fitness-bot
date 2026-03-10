package service

import (
	"context"
	"fitness-bot/internal/models"
	"fitness-bot/internal/repository"
	"time"
)

type WorkoutService struct {
	programRepo    repository.ProgramRepository
	workoutRepo    repository.WorkoutRepository
	exerciseRepo   repository.ExerciseRepository
	completionRepo repository.DailyCompletionRepository
}

func NewWorkoutService(
	programRepo repository.ProgramRepository,
	workoutRepo repository.WorkoutRepository,
	exerciseRepo repository.ExerciseRepository,
	completionRepo repository.DailyCompletionRepository,
) *WorkoutService {
	return &WorkoutService{
		programRepo:    programRepo,
		workoutRepo:    workoutRepo,
		exerciseRepo:   exerciseRepo,
		completionRepo: completionRepo,
	}
}

func (s *WorkoutService) ListPrograms(ctx context.Context, format, goal, level string) ([]models.Program, error) {
	return s.programRepo.ListPrograms(ctx, format, goal, level)
}

func (s *WorkoutService) GetProgram(ctx context.Context, id int) (*models.Program, error) {
	return s.programRepo.GetProgramByID(ctx, id)
}

func (s *WorkoutService) EnrollInProgram(ctx context.Context, userID int64, programID int) error {
	return s.programRepo.EnrollUser(ctx, userID, programID)
}

func (s *WorkoutService) GetActiveEnrollment(ctx context.Context, userID int64) (*models.UserProgramEnrollment, error) {
	return s.programRepo.GetActiveEnrollment(ctx, userID)
}

func (s *WorkoutService) ListWorkouts(ctx context.Context, format, goal, level string) ([]models.Workout, error) {
	return s.workoutRepo.ListWorkouts(ctx, format, goal, level)
}

func (s *WorkoutService) GetWorkout(ctx context.Context, id int) (*models.Workout, error) {
	return s.workoutRepo.GetWorkoutByID(ctx, id)
}

func (s *WorkoutService) GetWorkoutExercises(ctx context.Context, workoutID int) ([]models.WorkoutExercise, error) {
	return s.workoutRepo.ListExercises(ctx, workoutID)
}

func (s *WorkoutService) CompleteWorkout(ctx context.Context, userID int64, workoutID int) error {
	completion := &models.DailyCompletion{
		UserID:     userID,
		EntityType: "workout",
		EntityID:   workoutID,
		Date:       time.Now().Format("2006-01-02"),
		Status:     "completed",
	}
	return s.completionRepo.Create(ctx, completion)
}
