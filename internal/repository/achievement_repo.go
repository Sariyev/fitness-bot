package repository

import (
	"context"
	"fitness-bot/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type achievementRepo struct {
	pool *pgxpool.Pool
}

func NewAchievementRepo(pool *pgxpool.Pool) AchievementRepository {
	return &achievementRepo{pool: pool}
}

func (r *achievementRepo) ListAll(ctx context.Context) ([]models.Achievement, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, slug, name, description, icon, criteria, created_at, updated_at
		 FROM achievements ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var achievements []models.Achievement
	for rows.Next() {
		var a models.Achievement
		if err := rows.Scan(&a.ID, &a.Slug, &a.Name, &a.Description,
			&a.Icon, &a.Criteria, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		achievements = append(achievements, a)
	}
	return achievements, nil
}

func (r *achievementRepo) ListByUser(ctx context.Context, userID int64) ([]models.UserAchievement, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, achievement_id, earned_at
		 FROM user_achievements
		 WHERE user_id = $1
		 ORDER BY earned_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	achievements := []models.UserAchievement{}
	for rows.Next() {
		var a models.UserAchievement
		if err := rows.Scan(&a.ID, &a.UserID, &a.AchievementID, &a.EarnedAt); err != nil {
			return nil, err
		}
		achievements = append(achievements, a)
	}
	return achievements, nil
}

func (r *achievementRepo) Unlock(ctx context.Context, userID int64, achievementID int) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO user_achievements (user_id, achievement_id, earned_at)
		 VALUES ($1, $2, NOW())
		 ON CONFLICT (user_id, achievement_id) DO NOTHING`,
		userID, achievementID)
	return err
}
