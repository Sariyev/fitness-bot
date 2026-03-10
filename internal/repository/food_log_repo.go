package repository

import (
	"context"
	"fitness-bot/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type foodLogRepo struct {
	pool *pgxpool.Pool
}

func NewFoodLogRepo(pool *pgxpool.Pool) FoodLogRepository {
	return &foodLogRepo{pool: pool}
}

func (r *foodLogRepo) Create(ctx context.Context, entry *models.FoodLogEntry) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO food_log_entries (user_id, date, meal_type, food_name, calories, protein, fat, carbs, photo_url)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 RETURNING id, created_at`,
		entry.UserID, entry.Date, entry.MealType, entry.FoodName, entry.Calories,
		entry.Protein, entry.Fat, entry.Carbs, entry.PhotoURL,
	).Scan(&entry.ID, &entry.CreatedAt)
}

func (r *foodLogRepo) Delete(ctx context.Context, userID int64, id int64) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM food_log_entries WHERE id = $1 AND user_id = $2`,
		id, userID)
	return err
}

func (r *foodLogRepo) ListByDate(ctx context.Context, userID int64, date string) ([]models.FoodLogEntry, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, date, meal_type, food_name, calories, protein, fat, carbs, photo_url, created_at
		 FROM food_log_entries
		 WHERE user_id = $1 AND date = $2
		 ORDER BY created_at`, userID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := []models.FoodLogEntry{}
	for rows.Next() {
		var e models.FoodLogEntry
		if err := rows.Scan(&e.ID, &e.UserID, &e.Date, &e.MealType, &e.FoodName,
			&e.Calories, &e.Protein, &e.Fat, &e.Carbs, &e.PhotoURL, &e.CreatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func (r *foodLogRepo) GetDailySummary(ctx context.Context, userID int64, date string) (int, float64, float64, float64, error) {
	var calories int
	var protein, fat, carbs float64
	err := r.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(calories), 0), COALESCE(SUM(protein), 0),
			COALESCE(SUM(fat), 0), COALESCE(SUM(carbs), 0)
		 FROM food_log_entries
		 WHERE user_id = $1 AND date = $2`, userID, date,
	).Scan(&calories, &protein, &fat, &carbs)
	return calories, protein, fat, carbs, err
}
