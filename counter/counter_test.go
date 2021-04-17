package counter

import (
	"os"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

var (
	mockRedis  *miniredis.Miniredis
	mockClient *redis.Client
)

func TestNewCounter(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]struct {
		opts  []Option
		isErr bool
	}{
		"success": {
			opts:  []Option{WithRedisClient(mockClient)},
			isErr: false,
		},
	}

	for _, t := range tests {
		_, err := NewCounter(t.opts...)
		assert.Equal(t.isErr, err != nil)
	}
}

func TestMain(m *testing.M) {
	var err error
	mockRedis, err = miniredis.Run()
	if err != nil {
		panic(err)
	}
	mockClient = redis.NewClient(&redis.Options{Addr: mockRedis.Addr()})
	code := m.Run()
	mockRedis.Close()
	mockClient.Close()
	os.Exit(code)
}
