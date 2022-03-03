package limiter

import (
	"testing"

	redis "github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestWithRedisClient(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]struct {
		client *redis.Client
	}{
		"success": {client: mockClient},
	}

	for _, t := range tests {
		opt := WithRedisClient(t.client)
		c := &limiter{}
		opt(c)
		assert.Equal(c.Client, t.client)
	}
}
