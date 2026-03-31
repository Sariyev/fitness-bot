package payment

import "context"

type PaymentResult struct {
	ProviderTxID string
	Status       string // "completed", "pending", "failed"
	RedirectURL  string // URL to redirect user to for payment (empty for instant providers)
}

type Provider interface {
	CreatePayment(ctx context.Context, userID int64, amountKZT int) (*PaymentResult, error)
}

type CallbackVerifier interface {
	VerifyCallback(params map[string]string) (paymentID int64, err error)
}
