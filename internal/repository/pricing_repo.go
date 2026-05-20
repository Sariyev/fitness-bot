package repository

import (
	"context"
	"fitness-bot/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type pricingRepo struct {
	pool *pgxpool.Pool
}

func NewPricingRepo(pool *pgxpool.Pool) PricingRepository {
	return &pricingRepo{pool: pool}
}

func (r *pricingRepo) GetPrice(ctx context.Context, category models.Category) (int, error) {
	var price int
	err := r.pool.QueryRow(ctx,
		`SELECT price_kzt FROM category_pricing WHERE category = $1`,
		category,
	).Scan(&price)
	return price, err
}

func (r *pricingRepo) ListPrices(ctx context.Context) (map[models.Category]int, error) {
	rows, err := r.pool.Query(ctx, `SELECT category, price_kzt FROM category_pricing`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	prices := make(map[models.Category]int)
	for rows.Next() {
		var cat models.Category
		var price int
		if err := rows.Scan(&cat, &price); err != nil {
			return nil, err
		}
		prices[cat] = price
	}
	return prices, rows.Err()
}

func (r *pricingRepo) SetPrice(ctx context.Context, category models.Category, priceKZT int) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE category_pricing
		 SET price_kzt = $1, updated_at = NOW()
		 WHERE category = $2`,
		priceKZT, category,
	)
	return err
}
