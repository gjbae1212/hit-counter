package counter

import (
	"github.com/go-redis/redis/v8"
)

type (
	Option     interface{ apply(d *db) }
	OptionFunc func(d *db)
)

func (f OptionFunc) apply(d *db) { f(d) }

// WithRedisClient returns a function which sets redis client.
func WithRedisClient(client *redis.Client) OptionFunc {
	return func(d *db) {
		d.redisClient = client
	}
}
