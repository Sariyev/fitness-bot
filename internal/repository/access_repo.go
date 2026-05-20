package repository

import (
	"context"
	"fitness-bot/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type accessRepo struct {
	pool *pgxpool.Pool
}

func NewAccessRepo(pool *pgxpool.Pool) AccessRepository {
	return &accessRepo{pool: pool}
}

func (r *accessRepo) HasAccess(ctx context.Context, userID int64, category models.Category) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(
		     SELECT 1 FROM user_category_access
		     WHERE user_id = $1 AND category = $2
		 )`,
		userID, category,
	).Scan(&exists)
	return exists, err
}

func (r *accessRepo) Grant(ctx context.Context, userID int64, category models.Category, paymentID *int64) error {
	// Idempotent: if the row already exists, keep the original grant date
	// and don't overwrite payment_id (first payment wins).
	_, err := r.pool.Exec(ctx,
		`INSERT INTO user_category_access (user_id, category, payment_id)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (user_id, category) DO NOTHING`,
		userID, category, paymentID,
	)
	return err
}

func (r *accessRepo) ListGranted(ctx context.Context, userID int64) ([]models.Category, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT category FROM user_category_access WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, rows.Err()
}
