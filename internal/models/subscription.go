package models

import "time"

type SubscriptionPlan struct {
	ID           int       `json:"id"`
	Slug         string    `json:"slug"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	PriceKZT     int       `json:"price_kzt"`
	DurationDays int       `json:"duration_days"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}

type Subscription struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	PlanID    int       `json:"plan_id"`
	Status    string    `json:"status"`
	StartsAt  time.Time `json:"starts_at"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Payment struct {
	ID             int64          `json:"id"`
	UserID         int64          `json:"user_id"`
	SubscriptionID *int64         `json:"subscription_id"`
	PlanID         int            `json:"plan_id"`
	AmountKZT      int            `json:"amount_kzt"`
	Status         string         `json:"status"`
	Provider       string         `json:"provider"`
	ProviderTxID   *string        `json:"provider_tx_id"`
	Metadata       map[string]any `json:"metadata"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}
