package payment

import "context"

type PaymentResult struct {
	ProviderTxID string
	Status       string // "completed", "pending", "failed"
}

type Provider interface {
	CreatePayment(ctx context.Context, userID int64, amountKZT int) (*PaymentResult, error)
}
