package counter

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	timeout = 10 * time.Second
)

type (
	Counter interface {
		IncreaseHitOfDaily(ctx context.Context, id string, t time.Time) (*Score, error)
		IncreaseHitOfTotal(ctx context.Context, id string) (*Score, error)
		GetHitOfDaily(ctx context.Context, id string, t time.Time) (*Score, error)
		GetHitOfTotal(ctx context.Context, id string) (*Score, error)
		GetHitOfDailyAndTotal(ctx context.Context, id string, t time.Time) (daily *Score, total *Score, err error)
		IncreaseRankOfDaily(ctx context.Context, group, id string, t time.Time) (*Score, error)
		IncreaseRankOfTotal(ctx context.Context, group, id string) (*Score, error)
		GetRankDailyByLimit(ctx context.Context, group string, limit int, t time.Time) ([]*Score, error)
		GetRankTotalByLimit(ctx context.Context, group string, limit int) ([]*Score, error)
		GetHitOfDailyByRange(ctx context.Context, id string, timeRange []time.Time) (scores []*Score, err error)
	}

	db struct {
		redisClient *redis.Client
	}
)

// Score presents result for response.
type Score struct {
	Name  string
	Value int64
}

// NewCounter returns an object implemented counter interface.
func NewCounter(opts ...Option) (Counter, error) {
	c := &db{}
	for _, opt := range opts {
		opt.apply(c)
	}

	// if redis client doesn't exist, a default redis will be set up to have `localhost:6379`.
	if c.redisClient == nil {
		c.redisClient = redis.NewClient(&redis.Options{
			Addr:       "localhost:6379",
			Password:   "",
			DB:         0,
			MaxRetries: 1,
		})
	}
	return Counter(c), nil
}
