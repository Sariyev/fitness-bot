package repository

import (
	"context"
	"fitness-bot/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type scoreRepo struct {
	pool *pgxpool.Pool
}

func NewScoreRepo(pool *pgxpool.Pool) ScoreRepository {
	return &scoreRepo{pool: pool}
}

func (r *scoreRepo) Create(ctx context.Context, s *models.UserScore) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO user_scores (user_id, score_type, reference_type, reference_id, score, comment)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`,
		s.UserID, s.ScoreType, s.ReferenceType, s.ReferenceID, s.Score, s.Comment,
	).Scan(&s.ID, &s.CreatedAt)
}

func (r *scoreRepo) GetByReference(ctx context.Context, userID int64, refType string, refID int) (*models.UserScore, error) {
	s := &models.UserScore{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, score_type, reference_type, reference_id, score, comment, created_at
		 FROM user_scores
		 WHERE user_id = $1 AND reference_type = $2 AND reference_id = $3
		 ORDER BY created_at DESC LIMIT 1`,
		userID, refType, refID,
	).Scan(&s.ID, &s.UserID, &s.ScoreType, &s.ReferenceType, &s.ReferenceID,
		&s.Score, &s.Comment, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *scoreRepo) ListByUser(ctx context.Context, userID int64) ([]models.UserScore, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, score_type, reference_type, reference_id, score, comment, created_at
		 FROM user_scores WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scores []models.UserScore
	for rows.Next() {
		var s models.UserScore
		if err := rows.Scan(&s.ID, &s.UserID, &s.ScoreType, &s.ReferenceType, &s.ReferenceID,
			&s.Score, &s.Comment, &s.CreatedAt); err != nil {
			return nil, err
		}
		scores = append(scores, s)
	}
	return scores, nil
}
