package storage

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// R2Config carries the credentials and bucket names. Construct from
// internal/config so secrets stay out of this package.
type R2Config struct {
	AccountID       string
	AccessKeyID     string
	SecretAccessKey string
	BucketPrivate   string
	BucketPublic    string
	// PublicURL is the stable base for objects in BucketPublic, e.g.
	// "https://pub-xxxx.r2.dev" or a custom subdomain. If empty, falls back
	// to the R2 path-style URL (slower, no CDN, only for dev).
	PublicURL string
}

type R2Provider struct {
	cfg     R2Config
	client  *s3.Client
	presign *s3.PresignClient
}

func NewR2Provider(cfg R2Config) (*R2Provider, error) {
	if cfg.AccountID == "" || cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" {
		return nil, errors.New("r2: missing required credentials")
	}
	if cfg.BucketPrivate == "" || cfg.BucketPublic == "" {
		return nil, errors.New("r2: missing bucket names")
	}

	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountID)

	awsCfg := aws.Config{
		Region:      "auto",
		Credentials: credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		// R2 supports path-style; virtual-host style works too once a custom
		// domain is bound. Path-style is the safe default.
		o.UsePathStyle = true
	})

	return &R2Provider{
		cfg:     cfg,
		client:  client,
		presign: s3.NewPresignClient(client),
	}, nil
}

func (p *R2Provider) bucketName(kind BucketKind) (string, error) {
	switch kind {
	case BucketPrivate:
		return p.cfg.BucketPrivate, nil
	case BucketPublic:
		return p.cfg.BucketPublic, nil
	default:
		return "", fmt.Errorf("r2: unknown bucket kind %q", kind)
	}
}

func (p *R2Provider) PresignPut(ctx context.Context, kind BucketKind, key, contentType string, ttl time.Duration) (string, error) {
	bucket, err := p.bucketName(kind)
	if err != nil {
		return "", err
	}
	req, err := p.presign.PresignPutObject(ctx,
		&s3.PutObjectInput{
			Bucket:      aws.String(bucket),
			Key:         aws.String(key),
			ContentType: aws.String(contentType),
		},
		s3.WithPresignExpires(ttl),
	)
	if err != nil {
		return "", fmt.Errorf("r2: presign put: %w", err)
	}
	return req.URL, nil
}

func (p *R2Provider) PresignGet(ctx context.Context, kind BucketKind, key string, ttl time.Duration) (string, error) {
	bucket, err := p.bucketName(kind)
	if err != nil {
		return "", err
	}
	req, err := p.presign.PresignGetObject(ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		},
		s3.WithPresignExpires(ttl),
	)
	if err != nil {
		return "", fmt.Errorf("r2: presign get: %w", err)
	}
	return req.URL, nil
}

func (p *R2Provider) PublicURL(kind BucketKind, key string) string {
	if kind != BucketPublic {
		return ""
	}
	if p.cfg.PublicURL != "" {
		return strings.TrimRight(p.cfg.PublicURL, "/") + "/" + key
	}
	return fmt.Sprintf("https://%s.r2.cloudflarestorage.com/%s/%s", p.cfg.AccountID, p.cfg.BucketPublic, key)
}

func (p *R2Provider) HeadObject(ctx context.Context, kind BucketKind, key string) (*ObjectInfo, error) {
	bucket, err := p.bucketName(kind)
	if err != nil {
		return nil, err
	}
	out, err := p.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("r2: head: %w", err)
	}
	info := &ObjectInfo{}
	if out.ContentType != nil {
		info.ContentType = *out.ContentType
	}
	if out.ContentLength != nil {
		info.SizeBytes = *out.ContentLength
	}
	return info, nil
}

func (p *R2Provider) Delete(ctx context.Context, kind BucketKind, key string) error {
	bucket, err := p.bucketName(kind)
	if err != nil {
		return err
	}
	_, err = p.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("r2: delete: %w", err)
	}
	return nil
}
