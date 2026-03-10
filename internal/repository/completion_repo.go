package repository

import (
	"context"
	"fitness-bot/internal/models"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type completionRepo struct {
	pool *pgxpool.Pool
}

func NewCompletionRepo(pool *pgxpool.Pool) DailyCompletionRepository {
	return &completionRepo{pool: pool}
}

func (r *completionRepo) Create(ctx context.Context, c *models.DailyCompletion) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO daily_completions (user_id, entity_type, entity_id, date, status)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, created_at`,
		c.UserID, c.EntityType, c.EntityID, c.Date, c.Status,
	).Scan(&c.ID, &c.CreatedAt)
}

func (r *completionRepo) ListByDate(ctx context.Context, userID int64, date string) ([]models.DailyCompletion, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, entity_type, entity_id, date, status, created_at
		 FROM daily_completions
		 WHERE user_id = $1 AND date = $2
		 ORDER BY created_at`, userID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	completions := []models.DailyCompletion{}
	for rows.Next() {
		var c models.DailyCompletion
		if err := rows.Scan(&c.ID, &c.UserID, &c.EntityType, &c.EntityID,
			&c.Date, &c.Status, &c.CreatedAt); err != nil {
			return nil, err
		}
		completions = append(completions, c)
	}
	return completions, nil
}

func (r *completionRepo) GetStreak(ctx context.Context, userID int64) (int, int, error) {
	var current, longest int
	err := r.pool.QueryRow(ctx,
		`WITH distinct_dates AS (
			SELECT DISTINCT date::date AS d
			FROM daily_completions
			WHERE user_id = $1 AND status = 'done'
		),
		grouped AS (
			SELECT d, d - (ROW_NUMBER() OVER (ORDER BY d))::int * INTERVAL '1 day' AS grp
			FROM distinct_dates
		),
		streaks AS (
			SELECT COUNT(*) AS streak_len, MAX(d) AS streak_end
			FROM grouped
			GROUP BY grp
		)
		SELECT
			COALESCE((SELECT streak_len FROM streaks WHERE streak_end >= CURRENT_DATE - INTERVAL '1 day' ORDER BY streak_end DESC LIMIT 1), 0),
			COALESCE((SELECT MAX(streak_len) FROM streaks), 0)`,
		userID,
	).Scan(&current, &longest)
	return current, longest, err
}

func (r *completionRepo) GetCalendar(ctx context.Context, userID int64, year, month int) ([]string, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT DISTINCT date::text
		 FROM daily_completions
		 WHERE user_id = $1
		   AND date >= $2::date
		   AND date < ($2::date + INTERVAL '1 month')
		   AND status = 'done'
		 ORDER BY date`,
		userID, fmt.Sprintf("%d-%02d-01", year, month))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dates := []string{}
	for rows.Next() {
		var d string
		if err := rows.Scan(&d); err != nil {
			return nil, err
		}
		dates = append(dates, d)
	}
	return dates, nil
}
