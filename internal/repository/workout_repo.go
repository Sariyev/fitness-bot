package repository

import (
	"context"
	"fitness-bot/internal/models"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type workoutRepo struct {
	pool *pgxpool.Pool
}

func NewWorkoutRepo(pool *pgxpool.Pool) WorkoutRepository {
	return &workoutRepo{pool: pool}
}

func (r *workoutRepo) ListWorkouts(ctx context.Context, format, goal, level string) ([]models.Workout, error) {
	query := `SELECT id, program_id, slug, name, COALESCE(description,''), COALESCE(goal,''), COALESCE(format,''), COALESCE(level,''),
			  COALESCE(duration_minutes,0), equipment, COALESCE(expected_result,''), COALESCE(video_url,''), sort_order,
			  week_number, day_number, is_active, created_at, updated_at
			  FROM workouts WHERE is_active = TRUE`
	args := []interface{}{}
	idx := 1

	if format != "" {
		query += fmt.Sprintf(" AND format = $%d", idx)
		args = append(args, format)
		idx++
	}
	if goal != "" {
		query += fmt.Sprintf(" AND goal = $%d", idx)
		args = append(args, goal)
		idx++
	}
	if level != "" {
		query += fmt.Sprintf(" AND level = $%d", idx)
		args = append(args, level)
		idx++
	}
	query += " ORDER BY sort_order"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workouts := []models.Workout{}
	for rows.Next() {
		var w models.Workout
		if err := rows.Scan(&w.ID, &w.ProgramID, &w.Slug, &w.Name, &w.Description,
			&w.Goal, &w.Format, &w.Level, &w.DurationMinutes, &w.Equipment,
			&w.ExpectedResult, &w.VideoURL, &w.SortOrder, &w.WeekNumber, &w.DayNumber,
			&w.IsActive, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		workouts = append(workouts, w)
	}
	return workouts, nil
}

func (r *workoutRepo) ListAllWorkouts(ctx context.Context) ([]models.Workout, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, program_id, slug, name, COALESCE(description,''), COALESCE(goal,''), COALESCE(format,''), COALESCE(level,''),
			COALESCE(duration_minutes,0), equipment, COALESCE(expected_result,''), COALESCE(video_url,''), sort_order,
			week_number, day_number, is_active, created_at, updated_at
		 FROM workouts ORDER BY sort_order`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workouts := []models.Workout{}
	for rows.Next() {
		var w models.Workout
		if err := rows.Scan(&w.ID, &w.ProgramID, &w.Slug, &w.Name, &w.Description,
			&w.Goal, &w.Format, &w.Level, &w.DurationMinutes, &w.Equipment,
			&w.ExpectedResult, &w.VideoURL, &w.SortOrder, &w.WeekNumber, &w.DayNumber,
			&w.IsActive, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		workouts = append(workouts, w)
	}
	return workouts, nil
}

func (r *workoutRepo) GetWorkoutByID(ctx context.Context, id int) (*models.Workout, error) {
	w := &models.Workout{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, program_id, slug, name, COALESCE(description,''), COALESCE(goal,''), COALESCE(format,''), COALESCE(level,''),
			COALESCE(duration_minutes,0), equipment, COALESCE(expected_result,''), COALESCE(video_url,''), sort_order,
			week_number, day_number, is_active, created_at, updated_at
		 FROM workouts WHERE id = $1`, id,
	).Scan(&w.ID, &w.ProgramID, &w.Slug, &w.Name, &w.Description,
		&w.Goal, &w.Format, &w.Level, &w.DurationMinutes, &w.Equipment,
		&w.ExpectedResult, &w.VideoURL, &w.SortOrder, &w.WeekNumber, &w.DayNumber,
		&w.IsActive, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (r *workoutRepo) ListByProgram(ctx context.Context, programID int) ([]models.Workout, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, program_id, slug, name, COALESCE(description,''), COALESCE(goal,''), COALESCE(format,''), COALESCE(level,''),
			COALESCE(duration_minutes,0), equipment, COALESCE(expected_result,''), COALESCE(video_url,''), sort_order,
			week_number, day_number, is_active, created_at, updated_at
		 FROM workouts WHERE program_id = $1 AND is_active = TRUE
		 ORDER BY week_number, day_number, sort_order`, programID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workouts := []models.Workout{}
	for rows.Next() {
		var w models.Workout
		if err := rows.Scan(&w.ID, &w.ProgramID, &w.Slug, &w.Name, &w.Description,
			&w.Goal, &w.Format, &w.Level, &w.DurationMinutes, &w.Equipment,
			&w.ExpectedResult, &w.VideoURL, &w.SortOrder, &w.WeekNumber, &w.DayNumber,
			&w.IsActive, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		workouts = append(workouts, w)
	}
	return workouts, nil
}

func (r *workoutRepo) CreateWorkout(ctx context.Context, w *models.Workout) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO workouts (program_id, slug, name, description, goal, format, level,
			duration_minutes, equipment, expected_result, video_url, sort_order,
			week_number, day_number, is_active)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		 RETURNING id, created_at, updated_at`,
		w.ProgramID, w.Slug, w.Name, w.Description, w.Goal, w.Format, w.Level,
		w.DurationMinutes, w.Equipment, w.ExpectedResult, w.VideoURL, w.SortOrder,
		w.WeekNumber, w.DayNumber, w.IsActive,
	).Scan(&w.ID, &w.CreatedAt, &w.UpdatedAt)
}

func (r *workoutRepo) UpdateWorkout(ctx context.Context, w *models.Workout) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE workouts SET program_id=$2, slug=$3, name=$4, description=$5, goal=$6, format=$7,
			level=$8, duration_minutes=$9, equipment=$10, expected_result=$11, video_url=$12,
			sort_order=$13, week_number=$14, day_number=$15, is_active=$16, updated_at=NOW()
		 WHERE id=$1`,
		w.ID, w.ProgramID, w.Slug, w.Name, w.Description, w.Goal, w.Format,
		w.Level, w.DurationMinutes, w.Equipment, w.ExpectedResult, w.VideoURL,
		w.SortOrder, w.WeekNumber, w.DayNumber, w.IsActive)
	return err
}

func (r *workoutRepo) ListExercises(ctx context.Context, workoutID int) ([]models.WorkoutExercise, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, workout_id, exercise_id, COALESCE(sets,0), COALESCE(reps,''), COALESCE(duration_seconds,0), sort_order
		 FROM workout_exercises WHERE workout_id = $1 ORDER BY sort_order`, workoutID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	exercises := []models.WorkoutExercise{}
	for rows.Next() {
		var we models.WorkoutExercise
		if err := rows.Scan(&we.ID, &we.WorkoutID, &we.ExerciseID, &we.Sets,
			&we.Reps, &we.DurationSeconds, &we.SortOrder); err != nil {
			return nil, err
		}
		exercises = append(exercises, we)
	}
	return exercises, nil
}

func (r *workoutRepo) AddExercise(ctx context.Context, we *models.WorkoutExercise) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, duration_seconds, sort_order)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		we.WorkoutID, we.ExerciseID, we.Sets, we.Reps, we.DurationSeconds, we.SortOrder,
	).Scan(&we.ID)
}
