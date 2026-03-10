package repository

import (
	"context"
	"fitness-bot/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type exerciseRepo struct {
	pool *pgxpool.Pool
}

func NewExerciseRepo(pool *pgxpool.Pool) ExerciseRepository {
	return &exerciseRepo{pool: pool}
}

func (r *exerciseRepo) List(ctx context.Context) ([]models.Exercise, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, COALESCE(technique,''), COALESCE(common_mistakes,''), COALESCE(easier_modification,''),
			COALESCE(harder_modification,''), COALESCE(rest_seconds,60), created_at, updated_at
		 FROM exercises ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	exercises := []models.Exercise{}
	for rows.Next() {
		var e models.Exercise
		if err := rows.Scan(&e.ID, &e.Name, &e.Technique, &e.CommonMistakes,
			&e.EasierModification, &e.HarderModification, &e.RestSeconds,
			&e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		exercises = append(exercises, e)
	}
	return exercises, nil
}

func (r *exerciseRepo) GetByID(ctx context.Context, id int) (*models.Exercise, error) {
	e := &models.Exercise{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, COALESCE(technique,''), COALESCE(common_mistakes,''), COALESCE(easier_modification,''),
			COALESCE(harder_modification,''), COALESCE(rest_seconds,60), created_at, updated_at
		 FROM exercises WHERE id = $1`, id,
	).Scan(&e.ID, &e.Name, &e.Technique, &e.CommonMistakes,
		&e.EasierModification, &e.HarderModification, &e.RestSeconds,
		&e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *exerciseRepo) Create(ctx context.Context, e *models.Exercise) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO exercises (name, technique, common_mistakes, easier_modification,
			harder_modification, rest_seconds)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`,
		e.Name, e.Technique, e.CommonMistakes, e.EasierModification,
		e.HarderModification, e.RestSeconds,
	).Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt)
}

func (r *exerciseRepo) Update(ctx context.Context, e *models.Exercise) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE exercises SET name=$2, technique=$3, common_mistakes=$4,
			easier_modification=$5, harder_modification=$6, rest_seconds=$7, updated_at=NOW()
		 WHERE id=$1`,
		e.ID, e.Name, e.Technique, e.CommonMistakes,
		e.EasierModification, e.HarderModification, e.RestSeconds)
	return err
}
