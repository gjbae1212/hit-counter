package limiter

import (
	"context"

	"github.com/gjbae1212/hit-counter/internal"
	redis "github.com/go-redis/redis/v8"
)

type Limiter interface {
	IsBlackList(ctx context.Context, id string) (bool, error)
	AddBlackList(ctx context.Context, id string) error
	IsCacheLimit(ctx context.Context, id string) (bool, error)
	AddCacheLimit(ctx context.Context, id string) error
}

const (
	blacklistKey  = "hits:blacklist"
	cacheLimitKey = "hits:cache-limit"
)

type limiter struct {
	*redis.Client
}

// IsBlackList checks that this id is blacklist or not.
func (lt *limiter) IsBlackList(ctx context.Context, id string) (bool, error) {
	if ctx == nil || id == "" {
		return false, internal.ErrorEmptyParams
	}
	result, err := lt.SIsMember(ctx, blacklistKey, id).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return result, nil
}

// AddBlackList adds this argument to blacklist.
func (lt *limiter) AddBlackList(ctx context.Context, id string) error {
	if ctx == nil || id == "" {
		return internal.ErrorEmptyParams
	}
	return lt.SAdd(ctx, blacklistKey, id).Err()
}

// IsCacheLimit checks that this id is limited to save to cache.
func (lt *limiter) IsCacheLimit(ctx context.Context, id string) (bool, error) {
	if ctx == nil || id == "" {
		return false, internal.ErrorEmptyParams
	}

	result, err := lt.SIsMember(ctx, cacheLimitKey, id).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return result, nil
}

// AddCacheLimit adds this argument to cache limit.
func (lt *limiter) AddCacheLimit(ctx context.Context, id string) error {
	if ctx == nil || id == "" {
		return internal.ErrorEmptyParams
	}
	return lt.SAdd(ctx, cacheLimitKey, id).Err()
}

// NewLimiter creates limiter.
func NewLimiter(opts ...Option) (Limiter, error) {
	lt := &limiter{}
	for _, opt := range opts {
		opt.apply(lt)
	}

	// if redis client doesn't exist, a default redis will be set up to have `localhost:6379`.
	if lt.Client == nil {
		lt.Client = redis.NewClient(&redis.Options{
			Addr:       "localhost:6379",
			Password:   "",
			DB:         0,
			MaxRetries: 1,
		})
	}
	return lt, nil
}
