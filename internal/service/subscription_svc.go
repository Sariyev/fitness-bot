package service

import (
	"context"
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
