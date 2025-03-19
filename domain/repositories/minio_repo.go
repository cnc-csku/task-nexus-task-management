package repositories

import (
	"context"
	"os"
)

type MinioRepository interface {
	Upload(ctx context.Context, key string, object *os.File, contentType string) error
	GeneratePutPresignedURL(ctx context.Context, key string) (string, error)
	GetFullURL(key string) string
}
