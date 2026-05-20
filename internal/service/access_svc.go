package service

import (
	"context"
	"errors"
	"fitness-bot/internal/models"
	"fitness-bot/internal/repository"
	"time"
)

// TrialWindow defines how long after signup a user can access content
// tagged "trial" without paying. After this period, trial items behave
// like paid items.
//
// Constant for now; can be promoted to env config or per-category override
// if the business model requires it.
const TrialWindow = 7 * 24 * time.Hour

// ErrUnknownTier and ErrUnknownCategory protect against typos when callers
// pass enum values through configs / external input.
var (
	ErrUnknownTier     = errors.New("access: unknown tier")
	ErrUnknownCategory = errors.New("access: unknown category")
)

type AccessService struct {
	pricingRepo repository.PricingRepository
	accessRepo  repository.AccessRepository
}

func NewAccessService(p repository.PricingRepository, a repository.AccessRepository) *AccessService {
	return &AccessService{pricingRepo: p, accessRepo: a}
}

// CanAccess answers: is this user allowed to use an item with the given tier
// in the given category, right now?
//
//   - free  → always yes
//   - trial → yes if within TrialWindow of signup; otherwise treat as paid
//   - paid  → yes if user has a grant in user_category_access for this
//             category, OR the legacy users.is_paid flag is true (grandfathers
//             users who paid before the multi-category split)
func (s *AccessService) CanAccess(ctx context.Context, user *models.User, tier models.AccessTier, category models.Category) (bool, error) {
	if user == nil {
		return false, nil
	}
	if !category.IsValid() {
		return false, ErrUnknownCategory
	}
	switch tier {
	case models.AccessFree:
		return true, nil
	case models.AccessTrial:
		if time.Since(user.CreatedAt) < TrialWindow {
			return true, nil
		}
		// Trial expired — fall through to paid check.
		fallthrough
	case models.AccessPaid:
		if user.IsPaid {
			return true, nil
		}
		return s.accessRepo.HasAccess(ctx, user.ID, category)
	default:
		return false, ErrUnknownTier
	}
}

// Grant records that the user has bought permanent access to a category.
// Call after payment confirmation. Idempotent — re-grants are no-ops.
func (s *AccessService) Grant(ctx context.Context, userID int64, category models.Category, paymentID *int64) error {
	if !category.IsValid() {
		return ErrUnknownCategory
	}
	return s.accessRepo.Grant(ctx, userID, category, paymentID)
}

// ListGranted returns the categories the user has paid for.
func (s *AccessService) ListGranted(ctx context.Context, userID int64) ([]models.Category, error) {
	return s.accessRepo.ListGranted(ctx, userID)
}

// GetPrice returns the current price (in KZT) for a category.
func (s *AccessService) GetPrice(ctx context.Context, category models.Category) (int, error) {
	if !category.IsValid() {
		return 0, ErrUnknownCategory
	}
	return s.pricingRepo.GetPrice(ctx, category)
}

// ListPrices returns all category prices in one shot for the price-list UI.
func (s *AccessService) ListPrices(ctx context.Context) (map[models.Category]int, error) {
	return s.pricingRepo.ListPrices(ctx)
}

// SetPrice updates the price for a category (admin only).
func (s *AccessService) SetPrice(ctx context.Context, category models.Category, priceKZT int) error {
	if !category.IsValid() {
		return ErrUnknownCategory
	}
	if priceKZT <= 0 {
		return errors.New("access: price must be positive")
	}
	return s.pricingRepo.SetPrice(ctx, category, priceKZT)
}

// TrialRemaining returns how much of the trial window is left for the user.
// Zero or negative means the trial has expired.
func (s *AccessService) TrialRemaining(user *models.User) time.Duration {
	if user == nil {
		return 0
	}
	return TrialWindow - time.Since(user.CreatedAt)
}
