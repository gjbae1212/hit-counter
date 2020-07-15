package counter

import (
	"fmt"
	"time"

	redis "github.com/gjbae1212/go-redis"
)

var (
	timeout = 10 * time.Second
)

type (
	Counter interface {
		IncreaseHitOfDaily(id string, t time.Time) (*Score, error)
		IncreaseHitOfTotal(id string) (*Score, error)
		GetHitOfDaily(id string, t time.Time) (*Score, error)
		GetHitOfTotal(id string) (*Score, error)
		GetHitOfDailyAndTotal(id string, t time.Time) (daily *Score, total *Score, err error)
		IncreaseRankOfDaily(group, id string, t time.Time) (*Score, error)
		IncreaseRankOfTotal(group, id string) (*Score, error)
		GetRankDailyByLimit(group string, limit int, t time.Time) ([]*Score, error)
		GetRankTotalByLimit(group string, limit int) ([]*Score, error)
		GetHitOfDailyByRange(id string, timeRange []time.Time) (scores []*Score, err error)
	}

	db struct {
		redis redis.Manager
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
	// inject config
	for _, opt := range opts {
		if err := opt.apply(c); err != nil {
			return nil, fmt.Errorf("[err] NewCounter %w", err)
		}
	}

	// if a redis doesn't exist, a default redis will be set up to have `localhost:6379`.
	if c.redis == nil {
		rs, err := redis.NewManager([]string{"localhost:6379"})
		if err != nil {
			return nil, fmt.Errorf("[err] NewCounter %w", err)
		}
		c.redis = rs
	}
	return Counter(c), nil
}
