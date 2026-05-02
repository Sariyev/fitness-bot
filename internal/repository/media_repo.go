package repository

import (
	"context"
	"fitness-bot/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type mediaRepo struct {
	pool *pgxpool.Pool
}

func NewMediaRepo(pool *pgxpool.Pool) MediaRepository {
	return &mediaRepo{pool: pool}
}

func (r *mediaRepo) Create(ctx context.Context, m *models.Media) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO media
		 (storage_key, bucket, content_type, size_bytes, owner_user_id, reference_type, reference_id, is_public, confirmed)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 RETURNING id, created_at`,
		m.StorageKey, m.Bucket, m.ContentType, m.SizeBytes,
		m.OwnerUserID, m.ReferenceType, m.ReferenceID,
		m.IsPublic, m.Confirmed,
	).Scan(&m.ID, &m.CreatedAt)
}

func (r *mediaRepo) GetByID(ctx context.Context, id int64) (*models.Media, error) {
	m := &models.Media{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, storage_key, bucket, content_type, size_bytes,
		        owner_user_id, reference_type, reference_id, is_public, confirmed, created_at
		 FROM media WHERE id = $1`, id,
	).Scan(&m.ID, &m.StorageKey, &m.Bucket, &m.ContentType, &m.SizeBytes,
		&m.OwnerUserID, &m.ReferenceType, &m.ReferenceID,
		&m.IsPublic, &m.Confirmed, &m.CreatedAt)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *mediaRepo) MarkConfirmed(ctx context.Context, id int64, sizeBytes int64) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE media SET confirmed = TRUE, size_bytes = $2 WHERE id = $1`,
		id, sizeBytes,
	)
	return err
}

func (r *mediaRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM media WHERE id = $1`, id)
	return err
}

func (r *mediaRepo) TotalConfirmedBytes(ctx context.Context) (int64, error) {
	var total int64
	err := r.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(size_bytes), 0) FROM media WHERE confirmed = TRUE`,
	).Scan(&total)
	return total, err
}
