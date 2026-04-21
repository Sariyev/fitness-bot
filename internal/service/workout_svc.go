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

func (s *WorkoutService) GetWorkoutExercisesWithDetails(ctx context.Context, workoutID int) ([]models.WorkoutExerciseWithDetails, error) {
	exercises, err := s.workoutRepo.ListExercises(ctx, workoutID)
	if err != nil {
		return nil, err
	}

	result := make([]models.WorkoutExerciseWithDetails, 0, len(exercises))
	for _, we := range exercises {
		detail := models.WorkoutExerciseWithDetails{WorkoutExercise: we}
		if ex, err := s.exerciseRepo.GetByID(ctx, we.ExerciseID); err == nil {
			detail.ExerciseName = ex.Name
			detail.Technique = ex.Technique
			detail.CommonMistakes = ex.CommonMistakes
			detail.EasierModification = ex.EasierModification
			detail.HarderModification = ex.HarderModification
			detail.RestSeconds = ex.RestSeconds
		}
		result = append(result, detail)
	}
	return result, nil
}

func (s *WorkoutService) ListAllPrograms(ctx context.Context) ([]models.Program, error) {
	return s.programRepo.ListAllPrograms(ctx)
}

func (s *WorkoutService) CreateProgram(ctx context.Context, p *models.Program) error {
	return s.programRepo.CreateProgram(ctx, p)
}

func (s *WorkoutService) UpdateProgram(ctx context.Context, p *models.Program) error {
	return s.programRepo.UpdateProgram(ctx, p)
}

func (s *WorkoutService) ListAllWorkouts(ctx context.Context) ([]models.Workout, error) {
	return s.workoutRepo.ListAllWorkouts(ctx)
}

func (s *WorkoutService) CreateWorkout(ctx context.Context, w *models.Workout) error {
	return s.workoutRepo.CreateWorkout(ctx, w)
}

func (s *WorkoutService) UpdateWorkout(ctx context.Context, w *models.Workout) error {
	return s.workoutRepo.UpdateWorkout(ctx, w)
}

func (s *WorkoutService) ListExercises(ctx context.Context) ([]models.Exercise, error) {
	return s.exerciseRepo.List(ctx)
}

func (s *WorkoutService) GetExercise(ctx context.Context, id int) (*models.Exercise, error) {
	return s.exerciseRepo.GetByID(ctx, id)
}

func (s *WorkoutService) CreateExercise(ctx context.Context, e *models.Exercise) error {
	return s.exerciseRepo.Create(ctx, e)
}

func (s *WorkoutService) UpdateExercise(ctx context.Context, e *models.Exercise) error {
	return s.exerciseRepo.Update(ctx, e)
}

func (s *WorkoutService) AddExerciseToWorkout(ctx context.Context, we *models.WorkoutExercise) error {
	return s.workoutRepo.AddExercise(ctx, we)
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
