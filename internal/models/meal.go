package models

import "time"

type MealPlan struct {
	ID            int        `json:"id"`
	Slug          string     `json:"slug"`
	Name          string     `json:"name"`
	Goal          string     `json:"goal"`
	DayNumber     int        `json:"day_number"`
	Calories      int        `json:"calories"`
	Protein       float64    `json:"protein"`
	Fat           float64    `json:"fat"`
	Carbs         float64    `json:"carbs"`
	ImageMediaID  *int64     `json:"image_media_id,omitempty"`
	ImageURL      string     `json:"image_url,omitempty"` // populated by handler; not persisted.
	AccessTier    AccessTier `json:"access_tier"`
	// Locked: per-request, set by handler; not persisted.
	Locked    bool      `json:"locked"`
	IsActive  bool      `json:"is_active"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Meal struct {
	ID           int       `json:"id"`
	MealPlanID   int       `json:"meal_plan_id"`
	MealType     string    `json:"meal_type"`
	Name         string    `json:"name"`
	Recipe       string    `json:"recipe"`
	Calories     int       `json:"calories"`
	Protein      float64   `json:"protein"`
	Fat          float64   `json:"fat"`
	Carbs        float64   `json:"carbs"`
	Alternatives string    `json:"alternatives"`
	ImageMediaID *int64    `json:"image_media_id,omitempty"`
	ImageURL     string    `json:"image_url,omitempty"` // populated by handler; not persisted.
	SortOrder    int       `json:"sort_order"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type FoodLogEntry struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Date      string    `json:"date"`
	MealType  string    `json:"meal_type"`
	FoodName  string    `json:"food_name"`
	Calories  int       `json:"calories"`
	Protein   float64   `json:"protein"`
	Fat       float64   `json:"fat"`
	Carbs     float64   `json:"carbs"`
	PhotoURL  string    `json:"photo_url"`
	CreatedAt time.Time `json:"created_at"`
}
