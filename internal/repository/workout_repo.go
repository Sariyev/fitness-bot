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

// SELECT list shared by ListWorkouts / ListAllWorkouts / GetWorkoutByID.
// Keep in lockstep with the Scan in scanWorkout below.
const workoutSelectCols = `id, slug, name, COALESCE(description,''), COALESCE(goal,''),
	COALESCE(format,''), COALESCE(level,''), COALESCE(duration_minutes,0), equipment,
	COALESCE(expected_result,''), COALESCE(video_url,''), video_media_id, access_tier,
	sort_order, is_active, created_at, updated_at`

func scanWorkout(rows interface {
	Scan(dest ...interface{}) error
}, w *models.Workout) error {
	return rows.Scan(&w.ID, &w.Slug, &w.Name, &w.Description, &w.Goal, &w.Format, &w.Level,
		&w.DurationMinutes, &w.Equipment, &w.ExpectedResult, &w.VideoURL, &w.VideoMediaID,
		&w.AccessTier, &w.SortOrder, &w.IsActive, &w.CreatedAt, &w.UpdatedAt)
}

func (r *workoutRepo) ListWorkouts(ctx context.Context, format, goal, level string) ([]models.Workout, error) {
	query := `SELECT ` + workoutSelectCols + ` FROM workouts WHERE is_active = TRUE`
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
		if err := scanWorkout(rows, &w); err != nil {
			return nil, err
		}
		workouts = append(workouts, w)
	}
	return workouts, nil
}

func (r *workoutRepo) ListAllWorkouts(ctx context.Context) ([]models.Workout, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+workoutSelectCols+` FROM workouts ORDER BY sort_order`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workouts := []models.Workout{}
	for rows.Next() {
		var w models.Workout
		if err := scanWorkout(rows, &w); err != nil {
			return nil, err
		}
		workouts = append(workouts, w)
	}
	return workouts, nil
}

func (r *workoutRepo) GetWorkoutByID(ctx context.Context, id int) (*models.Workout, error) {
	w := &models.Workout{}
	err := scanWorkout(r.pool.QueryRow(ctx,
		`SELECT `+workoutSelectCols+` FROM workouts WHERE id = $1`, id), w)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (r *workoutRepo) CreateWorkout(ctx context.Context, w *models.Workout) error {
	tier := w.AccessTier
	if tier == "" {
		tier = models.AccessPaid
	}
	return r.pool.QueryRow(ctx,
		`INSERT INTO workouts (slug, name, description, goal, format, level,
			duration_minutes, equipment, expected_result, video_url, video_media_id,
			access_tier, sort_order, is_active)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		 RETURNING id, created_at, updated_at`,
		w.Slug, w.Name, w.Description, w.Goal, w.Format, w.Level,
		w.DurationMinutes, w.Equipment, w.ExpectedResult, w.VideoURL, w.VideoMediaID,
		tier, w.SortOrder, w.IsActive,
	).Scan(&w.ID, &w.CreatedAt, &w.UpdatedAt)
}

func (r *workoutRepo) UpdateWorkout(ctx context.Context, w *models.Workout) error {
	tier := w.AccessTier
	if tier == "" {
		tier = models.AccessPaid
	}
	_, err := r.pool.Exec(ctx,
		`UPDATE workouts SET slug=$2, name=$3, description=$4, goal=$5, format=$6,
			level=$7, duration_minutes=$8, equipment=$9, expected_result=$10, video_url=$11,
			video_media_id=$12, access_tier=$13, sort_order=$14, is_active=$15,
			updated_at=NOW()
		 WHERE id=$1`,
		w.ID, w.Slug, w.Name, w.Description, w.Goal, w.Format,
		w.Level, w.DurationMinutes, w.Equipment, w.ExpectedResult, w.VideoURL,
		w.VideoMediaID, tier, w.SortOrder, w.IsActive)
	return err
}

func (r *workoutRepo) DeleteWorkout(ctx context.Context, id int) error {
	// workout_exercises is a pure join — clean it up so the workout delete
	// doesn't fail on FK violation.
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `DELETE FROM workout_exercises WHERE workout_id=$1`, id); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `DELETE FROM workouts WHERE id=$1`, id); err != nil {
		return err
	}
	return tx.Commit(ctx)
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
