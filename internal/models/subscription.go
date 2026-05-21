package models

import "time"

type Payment struct {
	ID           int64          `json:"id"`
	UserID       int64          `json:"user_id"`
	AmountKZT    int            `json:"amount_kzt"`
	Status       string         `json:"status"`
	Provider     string         `json:"provider"`
	ProviderTxID *string        `json:"provider_tx_id"`
	Category     *string        `json:"category"` // NULL for legacy pre-split payments
	Metadata     map[string]any `json:"metadata"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}
