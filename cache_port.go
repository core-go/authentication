package auth

import (
	"context"
	"time"
)

type CachePort interface {
	Put(ctx context.Context, key string, obj interface{}, timeToLive time.Duration) error
	GetMany(ctx context.Context, key []string) (map[string]string, []string, error)
	Get(ctx context.Context, key string) (string, error)
}
