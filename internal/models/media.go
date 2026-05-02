package models

import "time"

// Media is one row in the media table — a tracked object stored in R2.
type Media struct {
	ID            int64     `json:"id"`
	StorageKey    string    `json:"storage_key"`
	Bucket        string    `json:"bucket"` // "private" | "public"
	ContentType   string    `json:"content_type"`
	SizeBytes     int64     `json:"size_bytes"`
	OwnerUserID   *int64    `json:"owner_user_id,omitempty"`
	ReferenceType *string   `json:"reference_type,omitempty"`
	ReferenceID   *int64    `json:"reference_id,omitempty"`
	IsPublic      bool      `json:"is_public"`
	Confirmed     bool      `json:"confirmed"`
	CreatedAt     time.Time `json:"created_at"`
}
