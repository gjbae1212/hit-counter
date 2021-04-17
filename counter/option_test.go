package counter

import (
	"testing"

	"github.com/go-redis/redis/v8"

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
		c := &db{}
		opt(c)
		assert.Equal(c.redisClient, t.client)
	}
}
