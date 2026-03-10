package repository

import (
	"context"
	"encoding/json"
	"fitness-bot/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type progressRepo struct {
	pool *pgxpool.Pool
}

func NewProgressRepo(pool *pgxpool.Pool) ProgressRepository {
	return &progressRepo{pool: pool}
}

func (r *progressRepo) Create(ctx context.Context, entry *models.ProgressEntry) error {
	measurementsJSON, err := json.Marshal(entry.Measurements)
	if err != nil {
		return err
	}

	return r.pool.QueryRow(ctx,
		`INSERT INTO progress_entries (user_id, date, weight_kg, measurements, photo_url, wellbeing, pain_level)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, created_at, updated_at`,
		entry.UserID, entry.Date, entry.WeightKg, measurementsJSON,
		entry.PhotoURL, entry.Wellbeing, entry.PainLevel,
	).Scan(&entry.ID, &entry.CreatedAt, &entry.UpdatedAt)
}

func (r *progressRepo) ListByUser(ctx context.Context, userID int64) ([]models.ProgressEntry, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, date, weight_kg, measurements, photo_url, wellbeing, pain_level,
			created_at, updated_at
		 FROM progress_entries
		 WHERE user_id = $1
		 ORDER BY date DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.ProgressEntry
	for rows.Next() {
		var e models.ProgressEntry
		var measurementsRaw []byte
		if err := rows.Scan(&e.ID, &e.UserID, &e.Date, &e.WeightKg, &measurementsRaw,
			&e.PhotoURL, &e.Wellbeing, &e.PainLevel, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		if measurementsRaw != nil {
			if err := json.Unmarshal(measurementsRaw, &e.Measurements); err != nil {
				return nil, err
			}
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func (r *progressRepo) GetWeightHistory(ctx context.Context, userID int64) ([]models.WeightPoint, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT date, weight_kg
		 FROM progress_entries
		 WHERE user_id = $1 AND weight_kg IS NOT NULL
		 ORDER BY date`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []models.WeightPoint
	for rows.Next() {
		var p models.WeightPoint
		if err := rows.Scan(&p.Date, &p.WeightKg); err != nil {
			return nil, err
		}
		points = append(points, p)
	}
	return points, nil
}
