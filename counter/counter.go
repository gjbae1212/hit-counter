package counter

import (
	"time"

	"github.com/gjbae1212/go-module/redis"
	"github.com/pkg/errors"
)

var (
	timeout = 10 * time.Second
)

// It is a option for the counter interface
type (
	Option interface {
		apply(d *db) error
	}

	OptionFunc func(d *db) error
)

func (f OptionFunc) apply(d *db) error {
	return f(d)
}

// It is the counter interface wrapped a db object
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
	}

	db struct {
		redis redis.Manager
	}
)

// It is a wrapper struct which returned result for response.
type Score struct {
	Name  string
	Value int64
}

// It is passed redis addresses.
func WithRedisOption(addrs []string) OptionFunc {
	return func(d *db) error {
		rs, err := redis.NewManager(addrs)
		if err != nil {
			return errors.Wrap(err, "[err] WithRedisAddrs")
		}
		d.redis = rs
		return nil
	}
}

func NewCounter(opts ...Option) (Counter, error) {
	c := &db{}
	for _, opt := range opts {
		if err := opt.apply(c); err != nil {
			return nil, errors.Wrap(err, "[err] NewCounter")
		}
	}
	// If a redis do not exist and a default redis will be set up to have `localhost:6379`.
	if c.redis == nil {
		rs, err := redis.NewManager([]string{"localhost:6379"})
		if err != nil {
			return nil, errors.Wrap(err, "[err]  NewCounter")
		}
		c.redis = rs
	}
	return Counter(c), nil
}
