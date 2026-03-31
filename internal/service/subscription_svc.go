package service

import (
	"context"
	"fmt"
	"fitness-bot/internal/models"
	"fitness-bot/internal/payment"
	"fitness-bot/internal/repository"
)

type PaymentService struct {
	paymentRepo repository.PaymentRepository
	userRepo    repository.UserRepository
	provider    payment.Provider
	priceKZT    int
}

func NewPaymentService(paymentRepo repository.PaymentRepository, userRepo repository.UserRepository, provider payment.Provider, priceKZT int) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		userRepo:    userRepo,
		provider:    provider,
		priceKZT:    priceKZT,
	}
}

// ProcessPayment is the synchronous flow used by the Telegram bot (DummyProvider).
func (s *PaymentService) ProcessPayment(ctx context.Context, user *models.User) error {
	result, err := s.provider.CreatePayment(ctx, user.ID, s.priceKZT)
	if err != nil {
		return err
	}

	txID := result.ProviderTxID
	pmt := &models.Payment{
		UserID:       user.ID,
		AmountKZT:    s.priceKZT,
		Status:       "completed",
		Provider:     "dummy",
		ProviderTxID: &txID,
	}
	if err := s.paymentRepo.CreatePayment(ctx, pmt); err != nil {
		return err
	}

	user.IsPaid = true
	return s.userRepo.Update(ctx, user)
}

// InitiatePayment starts the payment flow. For async providers (Robokassa), returns
// a redirect URL. For sync providers (Dummy), completes the payment inline.
func (s *PaymentService) InitiatePayment(ctx context.Context, user *models.User) (*payment.PaymentResult, error) {
	result, err := s.provider.CreatePayment(ctx, user.ID, s.priceKZT)
	if err != nil {
		return nil, err
	}

	// Sync provider (Dummy) — complete payment immediately
	if result.Status == "completed" {
		txID := result.ProviderTxID
		pmt := &models.Payment{
			UserID:       user.ID,
			AmountKZT:    s.priceKZT,
			Status:       "completed",
			Provider:     "dummy",
			ProviderTxID: &txID,
		}
		if err := s.paymentRepo.CreatePayment(ctx, pmt); err != nil {
			return nil, err
		}
		user.IsPaid = true
		if err := s.userRepo.Update(ctx, user); err != nil {
			return nil, err
		}
	}
	// Async provider (Robokassa) — pending record already created by provider

	return result, nil
}

// ConfirmPayment is called by the Robokassa callback to mark a pending payment as completed.
func (s *PaymentService) ConfirmPayment(ctx context.Context, paymentID int64) error {
	pmt, err := s.paymentRepo.GetPaymentByID(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("payment not found: %w", err)
	}

	if pmt.Status == "completed" {
		return nil // idempotent
	}

	pmt.Status = "completed"
	if err := s.paymentRepo.UpdatePayment(ctx, pmt); err != nil {
		return fmt.Errorf("update payment: %w", err)
	}

	user, err := s.userRepo.GetByID(ctx, pmt.UserID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	user.IsPaid = true
	return s.userRepo.Update(ctx, user)
}
