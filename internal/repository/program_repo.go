package repository

import (
	"context"
	"fitness-bot/internal/models"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type programRepo struct {
	pool *pgxpool.Pool
}

func NewProgramRepo(pool *pgxpool.Pool) ProgramRepository {
	return &programRepo{pool: pool}
}

func (r *programRepo) ListPrograms(ctx context.Context, format, goal, level string) ([]models.Program, error) {
	query := `SELECT id, slug, name, COALESCE(description,''), goal, format, level, duration_weeks,
			  is_active, sort_order, created_at, updated_at
			  FROM programs WHERE is_active = TRUE`
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

	programs := []models.Program{}
	for rows.Next() {
		var p models.Program
		if err := rows.Scan(&p.ID, &p.Slug, &p.Name, &p.Description, &p.Goal, &p.Format,
			&p.Level, &p.DurationWeeks, &p.IsActive, &p.SortOrder, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		programs = append(programs, p)
	}
	return programs, nil
}

func (r *programRepo) GetProgramByID(ctx context.Context, id int) (*models.Program, error) {
	p := &models.Program{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, slug, name, description, goal, format, level, duration_weeks,
			is_active, sort_order, created_at, updated_at
		 FROM programs WHERE id = $1`, id,
	).Scan(&p.ID, &p.Slug, &p.Name, &p.Description, &p.Goal, &p.Format,
		&p.Level, &p.DurationWeeks, &p.IsActive, &p.SortOrder, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *programRepo) CreateProgram(ctx context.Context, p *models.Program) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO programs (slug, name, description, goal, format, level, duration_weeks, is_active, sort_order)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 RETURNING id, created_at, updated_at`,
		p.Slug, p.Name, p.Description, p.Goal, p.Format, p.Level, p.DurationWeeks, p.IsActive, p.SortOrder,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *programRepo) UpdateProgram(ctx context.Context, p *models.Program) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE programs SET slug=$2, name=$3, description=$4, goal=$5, format=$6, level=$7,
			duration_weeks=$8, is_active=$9, sort_order=$10, updated_at=NOW()
		 WHERE id=$1`,
		p.ID, p.Slug, p.Name, p.Description, p.Goal, p.Format, p.Level,
		p.DurationWeeks, p.IsActive, p.SortOrder)
	return err
}

func (r *programRepo) EnrollUser(ctx context.Context, userID int64, programID int) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO user_program_enrollments (user_id, program_id, started_at, current_week, is_active)
		 VALUES ($1, $2, NOW(), 1, TRUE)`,
		userID, programID)
	return err
}

func (r *programRepo) GetActiveEnrollment(ctx context.Context, userID int64) (*models.UserProgramEnrollment, error) {
	e := &models.UserProgramEnrollment{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, program_id, started_at, current_week, is_active, created_at, updated_at
		 FROM user_program_enrollments
		 WHERE user_id = $1 AND is_active = TRUE
		 ORDER BY created_at DESC LIMIT 1`, userID,
	).Scan(&e.ID, &e.UserID, &e.ProgramID, &e.StartedAt, &e.CurrentWeek,
		&e.IsActive, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *programRepo) ListUserEnrollments(ctx context.Context, userID int64) ([]models.UserProgramEnrollment, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, program_id, started_at, current_week, is_active, created_at, updated_at
		 FROM user_program_enrollments
		 WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	enrollments := []models.UserProgramEnrollment{}
	for rows.Next() {
		var e models.UserProgramEnrollment
		if err := rows.Scan(&e.ID, &e.UserID, &e.ProgramID, &e.StartedAt, &e.CurrentWeek,
			&e.IsActive, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		enrollments = append(enrollments, e)
	}
	return enrollments, nil
}

