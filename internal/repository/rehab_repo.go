package repository

import (
	"context"
	"fitness-bot/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type rehabRepo struct {
	pool *pgxpool.Pool
}

func NewRehabRepo(pool *pgxpool.Pool) RehabRepository {
	return &rehabRepo{pool: pool}
}

func (r *rehabRepo) ListCourses(ctx context.Context, category string) ([]models.RehabCourse, error) {
	var query string
	var args []interface{}

	if category != "" {
		query = `SELECT id, slug, category, name, description, warnings, is_active, sort_order, created_at, updated_at
				 FROM rehab_courses WHERE is_active = TRUE AND category = $1
				 ORDER BY sort_order`
		args = append(args, category)
	} else {
		query = `SELECT id, slug, category, name, description, warnings, is_active, sort_order, created_at, updated_at
				 FROM rehab_courses WHERE is_active = TRUE
				 ORDER BY sort_order`
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	courses := []models.RehabCourse{}
	for rows.Next() {
		var c models.RehabCourse
		if err := rows.Scan(&c.ID, &c.Slug, &c.Category, &c.Name, &c.Description,
			&c.Warnings, &c.IsActive, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}
	return courses, nil
}

func (r *rehabRepo) GetCourseByID(ctx context.Context, id int) (*models.RehabCourse, error) {
	c := &models.RehabCourse{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, slug, category, name, description, warnings, is_active, sort_order, created_at, updated_at
		 FROM rehab_courses WHERE id = $1`, id,
	).Scan(&c.ID, &c.Slug, &c.Category, &c.Name, &c.Description,
		&c.Warnings, &c.IsActive, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *rehabRepo) CreateCourse(ctx context.Context, c *models.RehabCourse) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO rehab_courses (slug, category, name, description, warnings, is_active, sort_order)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, created_at, updated_at`,
		c.Slug, c.Category, c.Name, c.Description, c.Warnings, c.IsActive, c.SortOrder,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *rehabRepo) UpdateCourse(ctx context.Context, c *models.RehabCourse) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE rehab_courses SET slug=$2, category=$3, name=$4, description=$5, warnings=$6,
			is_active=$7, sort_order=$8, updated_at=NOW()
		 WHERE id=$1`,
		c.ID, c.Slug, c.Category, c.Name, c.Description, c.Warnings, c.IsActive, c.SortOrder)
	return err
}

func (r *rehabRepo) ListSessions(ctx context.Context, courseID int) ([]models.RehabSession, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, course_id, day_number, stage, video_url, duration_minutes,
			description, sort_order, created_at, updated_at
		 FROM rehab_sessions WHERE course_id = $1
		 ORDER BY sort_order`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions := []models.RehabSession{}
	for rows.Next() {
		var s models.RehabSession
		if err := rows.Scan(&s.ID, &s.CourseID, &s.DayNumber, &s.Stage, &s.VideoURL,
			&s.DurationMinutes, &s.Description, &s.SortOrder, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}
	return sessions, nil
}

func (r *rehabRepo) GetSessionByID(ctx context.Context, id int) (*models.RehabSession, error) {
	s := &models.RehabSession{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, course_id, day_number, stage, video_url, duration_minutes,
			description, sort_order, created_at, updated_at
		 FROM rehab_sessions WHERE id = $1`, id,
	).Scan(&s.ID, &s.CourseID, &s.DayNumber, &s.Stage, &s.VideoURL,
		&s.DurationMinutes, &s.Description, &s.SortOrder, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *rehabRepo) CreateSession(ctx context.Context, s *models.RehabSession) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO rehab_sessions (course_id, day_number, stage, video_url, duration_minutes,
			description, sort_order)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, created_at, updated_at`,
		s.CourseID, s.DayNumber, s.Stage, s.VideoURL, s.DurationMinutes,
		s.Description, s.SortOrder,
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
}

func (r *rehabRepo) UpdateSession(ctx context.Context, s *models.RehabSession) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE rehab_sessions SET course_id=$2, day_number=$3, stage=$4, video_url=$5,
			duration_minutes=$6, description=$7, sort_order=$8, updated_at=NOW()
		 WHERE id=$1`,
		s.ID, s.CourseID, s.DayNumber, s.Stage, s.VideoURL,
		s.DurationMinutes, s.Description, s.SortOrder)
	return err
}

func (r *rehabRepo) CreateProgress(ctx context.Context, p *models.UserRehabProgress) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO user_rehab_progress (user_id, course_id, session_id, day_number, completed_at, pain_level, comment)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id`,
		p.UserID, p.CourseID, p.SessionID, p.DayNumber, p.CompletedAt, p.PainLevel, p.Comment,
	).Scan(&p.ID)
}

func (r *rehabRepo) ListUserProgress(ctx context.Context, userID int64, courseID int) ([]models.UserRehabProgress, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, course_id, session_id, day_number, completed_at, pain_level, comment
		 FROM user_rehab_progress
		 WHERE user_id = $1 AND course_id = $2
		 ORDER BY completed_at`, userID, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	progress := []models.UserRehabProgress{}
	for rows.Next() {
		var p models.UserRehabProgress
		if err := rows.Scan(&p.ID, &p.UserID, &p.CourseID, &p.SessionID,
			&p.DayNumber, &p.CompletedAt, &p.PainLevel, &p.Comment); err != nil {
			return nil, err
		}
		progress = append(progress, p)
	}
	return progress, nil
}
