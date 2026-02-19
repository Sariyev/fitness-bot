package repository

import (
	"context"
	"fitness-bot/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type paymentRepo struct {
	pool *pgxpool.Pool
}

func NewPaymentRepo(pool *pgxpool.Pool) PaymentRepository {
	return &paymentRepo{pool: pool}
}

func (r *paymentRepo) CreatePayment(ctx context.Context, p *models.Payment) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO payments (user_id, amount_kzt, status, provider, provider_tx_id, metadata)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`,
		p.UserID, p.AmountKZT, p.Status, p.Provider, p.ProviderTxID, p.Metadata,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *paymentRepo) UpdatePayment(ctx context.Context, p *models.Payment) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE payments SET status=$2, provider_tx_id=$3, metadata=$4, updated_at=NOW()
		 WHERE id=$1`,
		p.ID, p.Status, p.ProviderTxID, p.Metadata)
	return err
}

func (r *paymentRepo) GetPaymentByID(ctx context.Context, id int64) (*models.Payment, error) {
	p := &models.Payment{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, amount_kzt, status, provider, provider_tx_id, metadata, created_at, updated_at
		 FROM payments WHERE id = $1`, id,
	).Scan(&p.ID, &p.UserID, &p.AmountKZT, &p.Status,
		&p.Provider, &p.ProviderTxID, &p.Metadata, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}
