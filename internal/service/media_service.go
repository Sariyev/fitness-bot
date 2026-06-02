package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fitness-bot/internal/models"
	"fitness-bot/internal/repository"
	"fitness-bot/internal/storage"
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	presignPutTTL = 15 * time.Minute
	presignGetTTL = 1 * time.Hour
)

// allowedContentTypes maps a content-type to its max size in bytes.
var allowedContentTypes = map[string]int64{
	"image/jpeg": 10 * 1024 * 1024,  // 10 MB
	"image/png":  10 * 1024 * 1024,  // 10 MB
	"image/webp": 10 * 1024 * 1024,  // 10 MB
	"video/mp4":  500 * 1024 * 1024, // 500 MB
}

// MediaService coordinates presigned uploads against R2 and tracks the
// resulting objects in the media table. A nil-safe constructor (NewMediaService)
// returns nil when storage isn't configured; callers should treat that as
// "media feature disabled" and skip wiring the routes.
type MediaService struct {
	repo       repository.MediaRepository
	provider   storage.Provider
	quotaBytes int64
}

func NewMediaService(repo repository.MediaRepository, provider storage.Provider, quotaBytes int64) *MediaService {
	return &MediaService{repo: repo, provider: provider, quotaBytes: quotaBytes}
}

// UploadRequest is what RequestUpload accepts from a handler.
type UploadRequest struct {
	ContentType   string
	SizeBytes     int64
	ReferenceType string // optional, e.g. "user_avatar"
	ReferenceID   *int64 // optional
	IsPublic      bool   // admin-only when true; service enforces this
}

// UploadResult is the presigned-PUT bundle returned to the client.
type UploadResult struct {
	UploadURL string
	MediaID   int64
	Key       string
}

var (
	ErrInvalidContentType  = errors.New("media: content-type not allowed")
	ErrSizeExceedsTypeCap  = errors.New("media: file too large for content-type")
	ErrSizeNotPositive     = errors.New("media: size must be > 0")
	ErrQuotaExceeded       = errors.New("media: storage quota exceeded")
	ErrPublicRequiresAdmin = errors.New("media: only admins may upload public media")
	ErrNotOwner            = errors.New("media: not owner")
	ErrSizeMismatch        = errors.New("media: uploaded size differs from declared size")
)

func (s *MediaService) RequestUpload(ctx context.Context, user *models.User, req UploadRequest) (*UploadResult, error) {
	if user == nil {
		return nil, errors.New("media: missing user")
	}

	// Validate content-type + size against the allowlist.
	maxBytes, ok := allowedContentTypes[req.ContentType]
	if !ok {
		return nil, ErrInvalidContentType
	}
	if req.SizeBytes <= 0 {
		return nil, ErrSizeNotPositive
	}
	if req.SizeBytes > maxBytes {
		return nil, ErrSizeExceedsTypeCap
	}

	// is_public is admin-only.
	if req.IsPublic && !user.IsAdmin() {
		return nil, ErrPublicRequiresAdmin
	}

	// Quota check: confirmed + this upload must stay under the cap.
	used, err := s.repo.TotalConfirmedBytes(ctx)
	if err != nil {
		return nil, fmt.Errorf("media: quota lookup: %w", err)
	}
	if used+req.SizeBytes > s.quotaBytes {
		return nil, ErrQuotaExceeded
	}
	if pct := float64(used+req.SizeBytes) / float64(s.quotaBytes); pct > 0.8 {
		log.Printf("[MEDIA] quota at %.0f%% (%d/%d bytes) — consider raising MEDIA_QUOTA_BYTES or pruning", pct*100, used+req.SizeBytes, s.quotaBytes)
	}

	// Build a key. Random hex prefix avoids collisions and makes user
	// content effectively unguessable even before signing.
	rnd, err := randomHex(8)
	if err != nil {
		return nil, fmt.Errorf("media: random: %w", err)
	}
	ext := extForContentType(req.ContentType)
	var key string
	if req.IsPublic {
		// public: organize by reference_type
		key = fmt.Sprintf("%s/%s%s", coalesce(req.ReferenceType, "misc"), rnd, ext)
	} else {
		// private: scoped per user
		key = fmt.Sprintf("users/%d/%s/%s%s", user.ID, coalesce(req.ReferenceType, "misc"), rnd, ext)
	}

	bucketKind := storage.BucketPrivate
	bucketName := "private"
	if req.IsPublic {
		bucketKind = storage.BucketPublic
		bucketName = "public"
	}

	// Insert unconfirmed row first; we have a valid id to return to the
	// client which it'll pass back during /confirm.
	media := &models.Media{
		StorageKey:    key,
		Bucket:        bucketName,
		ContentType:   req.ContentType,
		SizeBytes:     req.SizeBytes,
		OwnerUserID:   &user.ID,
		ReferenceType: nilIfEmpty(req.ReferenceType),
		ReferenceID:   req.ReferenceID,
		IsPublic:      req.IsPublic,
		Confirmed:     false,
	}
	if err := s.repo.Create(ctx, media); err != nil {
		return nil, fmt.Errorf("media: insert: %w", err)
	}

	// Now ask R2 for the presigned PUT.
	uploadURL, err := s.provider.PresignPut(ctx, bucketKind, key, req.ContentType, presignPutTTL)
	if err != nil {
		// Best-effort cleanup of the orphan row.
		_ = s.repo.Delete(ctx, media.ID)
		return nil, fmt.Errorf("media: presign: %w", err)
	}

	return &UploadResult{
		UploadURL: uploadURL,
		MediaID:   media.ID,
		Key:       key,
	}, nil
}

// ConfirmUpload is called by the client after a successful PUT. We verify
// the object exists with the expected size and flip confirmed=true.
func (s *MediaService) ConfirmUpload(ctx context.Context, user *models.User, mediaID int64) (*models.Media, error) {
	m, err := s.repo.GetByID(ctx, mediaID)
	if err != nil {
		return nil, err
	}
	if m.OwnerUserID == nil || *m.OwnerUserID != user.ID {
		return nil, ErrNotOwner
	}
	if m.Confirmed {
		return m, nil
	}

	info, err := s.provider.HeadObject(ctx, bucketKindFor(m.Bucket), m.StorageKey)
	if err != nil {
		return nil, fmt.Errorf("media: head: %w", err)
	}
	if info.SizeBytes != m.SizeBytes {
		// Trust the actual object size over the declared one. We still
		// store it so quota math is accurate.
		log.Printf("[MEDIA] size differs at confirm: media_id=%d declared=%d actual=%d", m.ID, m.SizeBytes, info.SizeBytes)
	}
	if err := s.repo.MarkConfirmed(ctx, m.ID, info.SizeBytes); err != nil {
		return nil, fmt.Errorf("media: mark confirmed: %w", err)
	}
	m.Confirmed = true
	m.SizeBytes = info.SizeBytes
	return m, nil
}

// GetURL returns a fetchable URL for the media. Public media gets the stable
// public URL; private media gets a short-lived presigned GET that only the
// owner can request.
func (s *MediaService) GetURL(ctx context.Context, user *models.User, mediaID int64) (string, error) {
	m, err := s.repo.GetByID(ctx, mediaID)
	if err != nil {
		return "", err
	}
	if m.IsPublic {
		return s.provider.PublicURL(storage.BucketPublic, m.StorageKey), nil
	}
	if user == nil || m.OwnerUserID == nil || *m.OwnerUserID != user.ID {
		return "", ErrNotOwner
	}
	return s.provider.PresignGet(ctx, storage.BucketPrivate, m.StorageKey, presignGetTTL)
}

// PresignReadURL returns a fetchable URL for any media without checking
// ownership. Callers must verify the user can access the parent content
// (workout / rehab session / etc.) first — used from content handlers that
// gate at the program / course / plan level via accessSvc.CanAccess.
//
// This decouples "who uploaded the video" (admin) from "who can play it"
// (any user with access to the parent content).
func (s *MediaService) PresignReadURL(ctx context.Context, mediaID int64) (string, error) {
	m, err := s.repo.GetByID(ctx, mediaID)
	if err != nil {
		return "", err
	}
	if m.IsPublic {
		return s.provider.PublicURL(storage.BucketPublic, m.StorageKey), nil
	}
	return s.provider.PresignGet(ctx, storage.BucketPrivate, m.StorageKey, presignGetTTL)
}

// Delete removes the media row and the underlying object. Owner or admin only.
func (s *MediaService) Delete(ctx context.Context, user *models.User, mediaID int64) error {
	m, err := s.repo.GetByID(ctx, mediaID)
	if err != nil {
		return err
	}
	if !user.IsAdmin() && (m.OwnerUserID == nil || *m.OwnerUserID != user.ID) {
		return ErrNotOwner
	}
	// Try to delete the object first; if R2 fails, leave the row alone so we
	// don't lose track of the orphan.
	if err := s.provider.Delete(ctx, bucketKindFor(m.Bucket), m.StorageKey); err != nil {
		return fmt.Errorf("media: r2 delete: %w", err)
	}
	return s.repo.Delete(ctx, mediaID)
}

// helpers

func bucketKindFor(bucket string) storage.BucketKind {
	if bucket == "public" {
		return storage.BucketPublic
	}
	return storage.BucketPrivate
}

func randomHex(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func extForContentType(ct string) string {
	switch ct {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "video/mp4":
		return ".mp4"
	}
	return ""
}

func coalesce(s, fallback string) string {
	if strings.TrimSpace(s) == "" {
		return fallback
	}
	return s
}

func nilIfEmpty(s string) *string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return &s
}
