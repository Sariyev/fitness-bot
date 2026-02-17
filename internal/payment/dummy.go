package payment

import (
	"context"
	"fmt"
	"time"
)

type DummyProvider struct{}

func NewDummyProvider() Provider {
	return &DummyProvider{}
}

func (d *DummyProvider) CreatePayment(_ context.Context, userID int64, amountKZT int) (*PaymentResult, error) {
	return &PaymentResult{
		ProviderTxID: fmt.Sprintf("dummy_%d_%d", userID, time.Now().UnixNano()),
		Status:       "completed",
	}, nil
}
