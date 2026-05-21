package service

import (
	"context"
	"fitness-bot/internal/models"
	"fitness-bot/internal/payment"
	"fitness-bot/internal/repository"
	"fmt"
	"strconv"
)

type PaymentService struct {
	paymentRepo repository.PaymentRepository
	userRepo    repository.UserRepository
	accessSvc   *AccessService
	provider    payment.Provider
	// fallbackPriceKZT is used only as a last resort if the access service
	// can't return a price for the requested category (e.g. unconfigured row).
	fallbackPriceKZT int
}

func NewPaymentService(
	paymentRepo repository.PaymentRepository,
	userRepo repository.UserRepository,
	accessSvc *AccessService,
	provider payment.Provider,
	fallbackPriceKZT int,
) *PaymentService {
	return &PaymentService{
		paymentRepo:      paymentRepo,
		userRepo:         userRepo,
		accessSvc:        accessSvc,
		provider:         provider,
		fallbackPriceKZT: fallbackPriceKZT,
	}
}

// priceFor returns the current price for a category, falling back if the
// access service can't answer.
func (s *PaymentService) priceFor(ctx context.Context, category models.Category) int {
	if s.accessSvc != nil {
		if p, err := s.accessSvc.GetPrice(ctx, category); err == nil && p > 0 {
			return p
		}
	}
	return s.fallbackPriceKZT
}

// ProcessPayment is the synchronous flow used by the Telegram bot
// (DummyProvider). Kept on the workouts category so old bot flows still work.
func (s *PaymentService) ProcessPayment(ctx context.Context, user *models.User) error {
	return s.processSync(ctx, user, models.CategoryWorkouts)
}

func (s *PaymentService) processSync(ctx context.Context, user *models.User, category models.Category) error {
	price := s.priceFor(ctx, category)
	result, err := s.provider.CreatePayment(ctx, user.ID, price)
	if err != nil {
		return err
	}

	txID := result.ProviderTxID
	cat := string(category)
	pmt := &models.Payment{
		UserID:       user.ID,
		AmountKZT:    price,
		Status:       "completed",
		Provider:     "dummy",
		ProviderTxID: &txID,
		Category:     &cat,
	}
	if err := s.paymentRepo.CreatePayment(ctx, pmt); err != nil {
		return err
	}

	if s.accessSvc != nil {
		if err := s.accessSvc.Grant(ctx, user.ID, category, &pmt.ID); err != nil {
			return fmt.Errorf("grant access: %w", err)
		}
	}

	user.IsPaid = true // legacy denormalised flag — keep for backwards compat
	return s.userRepo.Update(ctx, user)
}

// InitiatePayment starts the payment flow for a specific category. For async
// providers (Robokassa) returns a redirect URL; for sync (Dummy) completes
// inline and grants access immediately.
func (s *PaymentService) InitiatePayment(ctx context.Context, user *models.User, category models.Category) (*payment.PaymentResult, error) {
	if !category.IsValid() {
		return nil, fmt.Errorf("invalid category: %s", category)
	}
	price := s.priceFor(ctx, category)

	result, err := s.provider.CreatePayment(ctx, user.ID, price)
	if err != nil {
		return nil, err
	}

	// Sync provider (Dummy): complete inline + grant access.
	if result.Status == "completed" {
		txID := result.ProviderTxID
		cat := string(category)
		pmt := &models.Payment{
			UserID:       user.ID,
			AmountKZT:    price,
			Status:       "completed",
			Provider:     "dummy",
			ProviderTxID: &txID,
			Category:     &cat,
		}
		if err := s.paymentRepo.CreatePayment(ctx, pmt); err != nil {
			return nil, err
		}
		if s.accessSvc != nil {
			if err := s.accessSvc.Grant(ctx, user.ID, category, &pmt.ID); err != nil {
				return nil, fmt.Errorf("grant access: %w", err)
			}
		}
		user.IsPaid = true
		if err := s.userRepo.Update(ctx, user); err != nil {
			return nil, err
		}
		return result, nil
	}

	// Async provider (Robokassa): the provider just created a pending row.
	// Tag it with the category so ConfirmPayment can grant the right bucket.
	if result.ProviderTxID != "" {
		if id, parseErr := strconv.ParseInt(result.ProviderTxID, 10, 64); parseErr == nil {
			_ = s.paymentRepo.SetCategory(ctx, id, string(category))
		}
	}
	return result, nil
}

// ConfirmPayment is called by the Robokassa callback to mark a pending payment
// as completed and grant the corresponding category access.
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

	// Grant per-category access if we know which category this paid for.
	// Legacy payments (NULL category) just keep the global users.is_paid.
	if pmt.Category != nil && *pmt.Category != "" && s.accessSvc != nil {
		category := models.Category(*pmt.Category)
		if category.IsValid() {
			if err := s.accessSvc.Grant(ctx, user.ID, category, &pmt.ID); err != nil {
				return fmt.Errorf("grant access: %w", err)
			}
		}
	}

	user.IsPaid = true
	return s.userRepo.Update(ctx, user)
}
