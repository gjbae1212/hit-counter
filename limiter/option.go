package limiter

import redis "github.com/go-redis/redis/v8"

type (
	Option     interface{ apply(limiter *limiter) }
	OptionFunc func(limiter *limiter)
)

func (f OptionFunc) apply(limiter *limiter) { f(limiter) }

// WithRedisClient returns a function which sets redis client.
func WithRedisClient(client *redis.Client) OptionFunc {
	return func(limiter *limiter) {
		limiter.Client = client
	}
}
