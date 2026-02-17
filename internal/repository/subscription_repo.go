package repository

import (
	"context"
	"fitness-bot/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type subscriptionRepo struct {
	pool *pgxpool.Pool
}

func NewSubscriptionRepo(pool *pgxpool.Pool) SubscriptionRepository {
	return &subscriptionRepo{pool: pool}
}

func (r *subscriptionRepo) ListActivePlans(ctx context.Context) ([]models.SubscriptionPlan, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, slug, name, description, price_kzt, duration_days, is_active, created_at
		 FROM subscription_plans WHERE is_active = TRUE ORDER BY price_kzt`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []models.SubscriptionPlan
	for rows.Next() {
		var p models.SubscriptionPlan
		if err := rows.Scan(&p.ID, &p.Slug, &p.Name, &p.Description, &p.PriceKZT,
			&p.DurationDays, &p.IsActive, &p.CreatedAt); err != nil {
			return nil, err
		}
		plans = append(plans, p)
	}
	return plans, nil
}

func (r *subscriptionRepo) GetPlanByID(ctx context.Context, id int) (*models.SubscriptionPlan, error) {
	p := &models.SubscriptionPlan{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, slug, name, description, price_kzt, duration_days, is_active, created_at
		 FROM subscription_plans WHERE id = $1`, id,
	).Scan(&p.ID, &p.Slug, &p.Name, &p.Description, &p.PriceKZT,
		&p.DurationDays, &p.IsActive, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *subscriptionRepo) CreatePlan(ctx context.Context, plan *models.SubscriptionPlan) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO subscription_plans (slug, name, description, price_kzt, duration_days, is_active)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`,
		plan.Slug, plan.Name, plan.Description, plan.PriceKZT, plan.DurationDays, plan.IsActive,
	).Scan(&plan.ID, &plan.CreatedAt)
}

func (r *subscriptionRepo) UpdatePlan(ctx context.Context, plan *models.SubscriptionPlan) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE subscription_plans SET slug=$2, name=$3, description=$4, price_kzt=$5, duration_days=$6, is_active=$7
		 WHERE id=$1`,
		plan.ID, plan.Slug, plan.Name, plan.Description, plan.PriceKZT, plan.DurationDays, plan.IsActive)
	return err
}

func (r *subscriptionRepo) GetActiveSubscription(ctx context.Context, userID int64) (*models.Subscription, error) {
	s := &models.Subscription{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, plan_id, status, starts_at, expires_at, created_at
		 FROM subscriptions
		 WHERE user_id = $1 AND status = 'active' AND expires_at > NOW()
		 ORDER BY expires_at DESC LIMIT 1`, userID,
	).Scan(&s.ID, &s.UserID, &s.PlanID, &s.Status, &s.StartsAt, &s.ExpiresAt, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *subscriptionRepo) CreateSubscription(ctx context.Context, sub *models.Subscription) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO subscriptions (user_id, plan_id, status, starts_at, expires_at)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`,
		sub.UserID, sub.PlanID, sub.Status, sub.StartsAt, sub.ExpiresAt,
	).Scan(&sub.ID, &sub.CreatedAt)
}

func (r *subscriptionRepo) ExpireSubscription(ctx context.Context, subID int64) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE subscriptions SET status = 'expired' WHERE id = $1`, subID)
	return err
}

func (r *subscriptionRepo) CreatePayment(ctx context.Context, p *models.Payment) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO payments (user_id, subscription_id, plan_id, amount_kzt, status, provider, provider_tx_id, metadata)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at, updated_at`,
		p.UserID, p.SubscriptionID, p.PlanID, p.AmountKZT, p.Status, p.Provider, p.ProviderTxID, p.Metadata,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *subscriptionRepo) UpdatePayment(ctx context.Context, p *models.Payment) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE payments SET status=$2, provider_tx_id=$3, subscription_id=$4, metadata=$5, updated_at=NOW()
		 WHERE id=$1`,
		p.ID, p.Status, p.ProviderTxID, p.SubscriptionID, p.Metadata)
	return err
}

func (r *subscriptionRepo) GetPaymentByID(ctx context.Context, id int64) (*models.Payment, error) {
	p := &models.Payment{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, subscription_id, plan_id, amount_kzt, status, provider, provider_tx_id, metadata, created_at, updated_at
		 FROM payments WHERE id = $1`, id,
	).Scan(&p.ID, &p.UserID, &p.SubscriptionID, &p.PlanID, &p.AmountKZT, &p.Status,
		&p.Provider, &p.ProviderTxID, &p.Metadata, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}
