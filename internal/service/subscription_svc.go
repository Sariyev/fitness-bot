package service

import (
	"context"
	"fitness-bot/internal/models"
	"fitness-bot/internal/payment"
	"fitness-bot/internal/repository"
	"time"

	"github.com/jackc/pgx/v4"
)

type SubscriptionService struct {
	repo     repository.SubscriptionRepository
	provider payment.Provider
}

func NewSubscriptionService(repo repository.SubscriptionRepository, provider payment.Provider) *SubscriptionService {
	return &SubscriptionService{repo: repo, provider: provider}
}

func (s *SubscriptionService) HasActiveSubscription(ctx context.Context, userID int64) (bool, error) {
	sub, err := s.repo.GetActiveSubscription(ctx, userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return sub != nil, nil
}

func (s *SubscriptionService) GetActiveSubscription(ctx context.Context, userID int64) (*models.Subscription, error) {
	sub, err := s.repo.GetActiveSubscription(ctx, userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return sub, nil
}

func (s *SubscriptionService) ListPlans(ctx context.Context) ([]models.SubscriptionPlan, error) {
	return s.repo.ListActivePlans(ctx)
}

func (s *SubscriptionService) GetPlanByID(ctx context.Context, id int) (*models.SubscriptionPlan, error) {
	return s.repo.GetPlanByID(ctx, id)
}

func (s *SubscriptionService) Subscribe(ctx context.Context, userID int64, planID int) (*models.Subscription, error) {
	plan, err := s.repo.GetPlanByID(ctx, planID)
	if err != nil {
		return nil, err
	}

	// Process payment via provider
	result, err := s.provider.CreatePayment(ctx, userID, plan.PriceKZT)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	sub := &models.Subscription{
		UserID:    userID,
		PlanID:    planID,
		Status:    "active",
		StartsAt:  now,
		ExpiresAt: now.AddDate(0, 0, plan.DurationDays),
	}
	if err := s.repo.CreateSubscription(ctx, sub); err != nil {
		return nil, err
	}

	txID := result.ProviderTxID
	pmt := &models.Payment{
		UserID:         userID,
		SubscriptionID: &sub.ID,
		PlanID:         planID,
		AmountKZT:      plan.PriceKZT,
		Status:         "completed",
		Provider:       "dummy",
		ProviderTxID:   &txID,
	}
	s.repo.CreatePayment(ctx, pmt)

	return sub, nil
}
