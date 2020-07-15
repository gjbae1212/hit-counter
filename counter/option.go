package counter

import (
	"fmt"

	redis "github.com/gjbae1212/go-redis"
)

type (
	Option interface {
		apply(d *db) error
	}

	OptionFunc func(d *db) error
)

func (f OptionFunc) apply(d *db) error {
	return f(d)
}

// WithRedisOption returns a function which sets a redis addr to db object.
func WithRedisOption(addrs []string) OptionFunc {
	return func(d *db) error {
		rs, err := redis.NewManager(addrs)
		if err != nil {
			return fmt.Errorf("[err] WithRedisOption %w", err)
		}
		d.redis = rs
		return nil
	}
}
