package redis

import (
	"context"
	"time"
)

type RedisRepository interface {
	Get(ctx context.Context, key string) (any, error)
	Set(ctx context.Context, key string, value string) error
	SetTTL(ctx context.Context, key string, value string, expiration time.Duration) error
	Has(ctx context.Context, key string) (bool, error)
	Delete(ctx context.Context, key string) error
	Scan(ctx context.Context, cursor uint64, pattern string, count int64) (keys []string, nextCursor uint64, err error)
	MGet(ctx context.Context, keys ...string) (any, error)
}

type redisRepository struct {
	// redisClient *redis.Client
}

func New( /*redisClient *redis.Client*/ ) RedisRepository {
	return &redisRepository{
		// redisClient: redisClient,
	}
}
