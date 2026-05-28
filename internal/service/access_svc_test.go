package service

import (
	"context"
	"fitness-bot/internal/models"
	"testing"
	"time"
)

// fakeAccessRepo is an in-memory AccessRepository for unit tests.
type fakeAccessRepo struct {
	granted map[int64]map[models.Category]bool
}

func newFakeAccessRepo() *fakeAccessRepo {
	return &fakeAccessRepo{granted: map[int64]map[models.Category]bool{}}
}

func (r *fakeAccessRepo) HasAccess(_ context.Context, userID int64, cat models.Category) (bool, error) {
	return r.granted[userID][cat], nil
}

func (r *fakeAccessRepo) Grant(_ context.Context, userID int64, cat models.Category, _ *int64) error {
	if r.granted[userID] == nil {
		r.granted[userID] = map[models.Category]bool{}
	}
	r.granted[userID][cat] = true
	return nil
}

func (r *fakeAccessRepo) ListGranted(_ context.Context, userID int64) ([]models.Category, error) {
	cats := []models.Category{}
	for c, ok := range r.granted[userID] {
		if ok {
			cats = append(cats, c)
		}
	}
	return cats, nil
}

// fakePricingRepo is a static-price stub for tests.
type fakePricingRepo struct{ price int }

func (r *fakePricingRepo) GetPrice(_ context.Context, _ models.Category) (int, error) {
	return r.price, nil
}
func (r *fakePricingRepo) ListPrices(_ context.Context) (map[models.Category]int, error) {
	return map[models.Category]int{
		models.CategoryWorkouts:  r.price,
		models.CategoryLFK:       r.price,
		models.CategoryNutrition: r.price,
	}, nil
}
func (r *fakePricingRepo) SetPrice(_ context.Context, _ models.Category, p int) error {
	r.price = p
	return nil
}

func newTestUser(daysAgo int, paid bool) *models.User {
	return &models.User{
		ID:        1,
		CreatedAt: time.Now().Add(-time.Duration(daysAgo) * 24 * time.Hour),
		IsPaid:    paid,
	}
}

func TestCanAccess_Free_AlwaysYes(t *testing.T) {
	svc := NewAccessService(&fakePricingRepo{price: 5000}, newFakeAccessRepo())

	for _, days := range []int{0, 3, 7, 30, 365} {
		u := newTestUser(days, false)
		ok, err := svc.CanAccess(context.Background(), u, models.AccessFree, models.CategoryWorkouts)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if !ok {
			t.Errorf("free should always be accessible (user age %dd)", days)
		}
	}
}

func TestCanAccess_Trial_WithinWindow(t *testing.T) {
	svc := NewAccessService(&fakePricingRepo{price: 5000}, newFakeAccessRepo())
	u := newTestUser(3, false) // 3 days old, within 7-day window

	ok, err := svc.CanAccess(context.Background(), u, models.AccessTrial, models.CategoryWorkouts)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !ok {
		t.Error("trial within window should be accessible")
	}
}

func TestCanAccess_Trial_ExpiredFallsThroughToPaid(t *testing.T) {
	repo := newFakeAccessRepo()
	svc := NewAccessService(&fakePricingRepo{price: 5000}, repo)
	u := newTestUser(10, false) // older than 7-day trial

	// No grant + no legacy paid → trial-expired user is locked out.
	ok, err := svc.CanAccess(context.Background(), u, models.AccessTrial, models.CategoryWorkouts)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if ok {
		t.Error("expired trial without payment should be locked")
	}

	// Now grant access — should unlock.
	_ = repo.Grant(context.Background(), u.ID, models.CategoryWorkouts, nil)
	ok, _ = svc.CanAccess(context.Background(), u, models.AccessTrial, models.CategoryWorkouts)
	if !ok {
		t.Error("after grant, trial-tier content should unlock")
	}
}

func TestCanAccess_Paid_RequiresGrant(t *testing.T) {
	repo := newFakeAccessRepo()
	svc := NewAccessService(&fakePricingRepo{price: 5000}, repo)
	u := newTestUser(3, false)

	ok, _ := svc.CanAccess(context.Background(), u, models.AccessPaid, models.CategoryWorkouts)
	if ok {
		t.Error("paid content should be locked without grant or legacy is_paid")
	}

	_ = repo.Grant(context.Background(), u.ID, models.CategoryWorkouts, nil)
	ok, _ = svc.CanAccess(context.Background(), u, models.AccessPaid, models.CategoryWorkouts)
	if !ok {
		t.Error("paid content should unlock after grant")
	}
}

func TestCanAccess_LegacyIsPaid_GrandfathersAllCategories(t *testing.T) {
	repo := newFakeAccessRepo() // no grants
	svc := NewAccessService(&fakePricingRepo{price: 5000}, repo)
	u := newTestUser(30, true) // legacy paid, trial long expired

	for _, cat := range []models.Category{models.CategoryWorkouts, models.CategoryLFK, models.CategoryNutrition} {
		ok, _ := svc.CanAccess(context.Background(), u, models.AccessPaid, cat)
		if !ok {
			t.Errorf("legacy is_paid user should have %s access without per-category grant", cat)
		}
	}
}

func TestCanAccess_PerCategoryIsolation(t *testing.T) {
	repo := newFakeAccessRepo()
	svc := NewAccessService(&fakePricingRepo{price: 5000}, repo)
	u := newTestUser(30, false)

	_ = repo.Grant(context.Background(), u.ID, models.CategoryWorkouts, nil)

	ok, _ := svc.CanAccess(context.Background(), u, models.AccessPaid, models.CategoryWorkouts)
	if !ok {
		t.Error("granted category should unlock")
	}
	ok, _ = svc.CanAccess(context.Background(), u, models.AccessPaid, models.CategoryLFK)
	if ok {
		t.Error("non-granted category should stay locked (no cross-bucket leak)")
	}
}

func TestCanAccess_InvalidCategoryRejected(t *testing.T) {
	svc := NewAccessService(&fakePricingRepo{price: 5000}, newFakeAccessRepo())
	u := newTestUser(0, true)

	if _, err := svc.CanAccess(context.Background(), u, models.AccessPaid, "bogus"); err != ErrUnknownCategory {
		t.Errorf("expected ErrUnknownCategory, got %v", err)
	}
}

func TestCanAccess_NilUser_LocksOut(t *testing.T) {
	svc := NewAccessService(&fakePricingRepo{price: 5000}, newFakeAccessRepo())
	ok, _ := svc.CanAccess(context.Background(), nil, models.AccessFree, models.CategoryWorkouts)
	if ok {
		t.Error("nil user should never have access")
	}
}

func TestSetPrice_RejectsNonPositive(t *testing.T) {
	svc := NewAccessService(&fakePricingRepo{price: 5000}, newFakeAccessRepo())
	if err := svc.SetPrice(context.Background(), models.CategoryWorkouts, 0); err == nil {
		t.Error("expected error for price <= 0")
	}
	if err := svc.SetPrice(context.Background(), models.CategoryWorkouts, -100); err == nil {
		t.Error("expected error for negative price")
	}
}

func TestTrialRemaining(t *testing.T) {
	svc := NewAccessService(&fakePricingRepo{price: 5000}, newFakeAccessRepo())

	fresh := newTestUser(0, false)
	if rem := svc.TrialRemaining(fresh); rem < 6*24*time.Hour {
		t.Errorf("fresh user should have ~7d trial, got %v", rem)
	}

	expired := newTestUser(10, false)
	if rem := svc.TrialRemaining(expired); rem > 0 {
		t.Errorf("10-day-old user should have 0 trial remaining, got %v", rem)
	}
}
