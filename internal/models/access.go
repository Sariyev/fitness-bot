package models

import "time"

// Category identifies one of the three paid content domains.
// The values are persisted in the DB; do not rename without a migration.
type Category string

const (
	CategoryWorkouts  Category = "workouts"
	CategoryLFK       Category = "lfk"
	CategoryNutrition Category = "nutrition"
)

// AccessTier is which bucket a content item lives in within its category.
// Values are persisted in the DB.
type AccessTier string

const (
	AccessFree  AccessTier = "free"
	AccessTrial AccessTier = "trial"
	AccessPaid  AccessTier = "paid"
)

// IsValid reports whether the value is one of the known tiers.
func (t AccessTier) IsValid() bool {
	switch t {
	case AccessFree, AccessTrial, AccessPaid:
		return true
	}
	return false
}

func (c Category) IsValid() bool {
	switch c {
	case CategoryWorkouts, CategoryLFK, CategoryNutrition:
		return true
	}
	return false
}

// CategoryPricing is the admin-editable price for a category.
type CategoryPricing struct {
	Category  Category  `json:"category"`
	PriceKZT  int       `json:"price_kzt"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserCategoryAccess records that a user permanently has the paid tier
// of one category (set on successful payment).
type UserCategoryAccess struct {
	UserID    int64     `json:"user_id"`
	Category  Category  `json:"category"`
	PaymentID *int64    `json:"payment_id"`
	GrantedAt time.Time `json:"granted_at"`
}
