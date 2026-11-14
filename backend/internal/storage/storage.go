package storage

import (
	"context"
	"io"
)

// Storage is blob storage (MinIO now; swappable later).
type Storage interface {
	Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) error
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
}
