package storage

import (
	"context"
	"time"
)

// BucketKind selects between the two R2 buckets we maintain.
type BucketKind string

const (
	BucketPrivate BucketKind = "private"
	BucketPublic  BucketKind = "public"
)

// ObjectInfo is the subset of metadata returned by HeadObject.
type ObjectInfo struct {
	ContentType string
	SizeBytes   int64
}

// Provider abstracts an S3-compatible object store. R2Provider is the only
// implementation today; the interface exists so tests can stub it out and so
// a future MinIO local-dev provider would be a drop-in.
type Provider interface {
	// PresignPut returns a URL the client can PUT a file directly to. The URL
	// is bound to the given content-type and expires after ttl.
	PresignPut(ctx context.Context, kind BucketKind, key, contentType string, ttl time.Duration) (string, error)

	// PresignGet returns a URL the client can GET the object from. Used for
	// private content; for public content, callers should use PublicURL.
	PresignGet(ctx context.Context, kind BucketKind, key string, ttl time.Duration) (string, error)

	// PublicURL returns the stable, CDN-cacheable URL for an object in the
	// public bucket. Returns an empty string for the private bucket.
	PublicURL(kind BucketKind, key string) string

	// HeadObject fetches metadata. Used after a presigned upload to confirm
	// the client actually wrote what it claimed.
	HeadObject(ctx context.Context, kind BucketKind, key string) (*ObjectInfo, error)

	// Delete removes the object.
	Delete(ctx context.Context, kind BucketKind, key string) error
}
