package models

import "time"

// Workout is a single training session. Workouts used to be children of a
// `Program` row (with week/day numbering), but the curriculum layer was
// dropped — migration 028 flattened everything into the workouts table and
// copied the program metadata down into each workout.
type Workout struct {
	ID              int        `json:"id"`
	Slug            string     `json:"slug"`
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	Goal            string     `json:"goal"`
	Format          string     `json:"format"`
	Level           string     `json:"level"`
	DurationMinutes int        `json:"duration_minutes"`
	Equipment       []string   `json:"equipment"`
	ExpectedResult  string     `json:"expected_result"`
	// VideoURL is the public/external URL the client plays. For admin-
	// uploaded videos it's populated from media at handler-resolution time;
	// for external (YouTube etc.) it's the raw column value.
	VideoURL string `json:"video_url"`
	// VideoMediaID points at the media row for admin-uploaded videos. NULL
	// when no upload (external URL or no video).
	VideoMediaID *int64     `json:"video_media_id"`
	AccessTier   AccessTier `json:"access_tier"`
	// Locked is set per-request by the handler based on the viewing user's
	// access; not persisted, not scanned from DB.
	Locked    bool      `json:"locked"`
	SortOrder int       `json:"sort_order"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
