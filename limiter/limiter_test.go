package limiter

import (
	"context"
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

func TestNewLimiter(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]struct {
		opts  []Option
		isErr bool
	}{
		"success": {isErr: false},
	}

	for _, t := range tests {
		_, err := NewLimiter(t.opts...)
		assert.Equal(t.isErr, err != nil)
	}
}

func TestLimiter_AddBlackList(t *testing.T) {
	assert := assert.New(t)

	lt, err := NewLimiter(WithRedisClient(mockClient))
	assert.NoError(err)

	tests := map[string]struct {
		id    string
		isErr bool
	}{
		"fail":    {isErr: true},
		"success": {isErr: false, id: "allan-1"},
	}

	for _, t := range tests {
		err := lt.AddBlackList(context.Background(), t.id)
		assert.Equal(t.isErr, err != nil)
		if err == nil {
			result, err := lt.IsBlackList(context.Background(), t.id)
			assert.NoError(err)
			assert.True(result)
		}
	}
}

func TestLimiter_IsBlackList(t *testing.T) {
	assert := assert.New(t)

	lt, err := NewLimiter(WithRedisClient(mockClient))
	assert.NoError(err)
	err = lt.AddBlackList(context.Background(), "allan-2")
	assert.NoError(err)

	tests := map[string]struct {
		id    string
		ok    bool
		isErr bool
	}{
		"fail":      {isErr: true},
		"found":     {isErr: false, id: "allan-2", ok: true},
		"not found": {isErr: false, id: "empty", ok: false},
	}

	for _, t := range tests {
		ok, err := lt.IsBlackList(context.Background(), t.id)
		assert.Equal(t.isErr, err != nil)
		if err == nil {
			assert.Equal(t.ok, ok)
		}
	}
}

func TestLimiter_AddCacheLimit(t *testing.T) {
	assert := assert.New(t)

	lt, err := NewLimiter(WithRedisClient(mockClient))
	assert.NoError(err)

	tests := map[string]struct {
		id    string
		isErr bool
	}{
		"fail":    {isErr: true},
		"success": {isErr: false, id: "allan-3"},
	}

	for _, t := range tests {
		err := lt.AddCacheLimit(context.Background(), t.id)
		assert.Equal(t.isErr, err != nil)
		if err == nil {
			result, err := lt.IsCacheLimit(context.Background(), t.id)
			assert.NoError(err)
			assert.True(result)
		}
	}
}

func TestLimiter_IsCacheLimit(t *testing.T) {
	assert := assert.New(t)

	lt, err := NewLimiter(WithRedisClient(mockClient))
	assert.NoError(err)
	err = lt.AddCacheLimit(context.Background(), "allan-4")
	assert.NoError(err)

	tests := map[string]struct {
		id    string
		ok    bool
		isErr bool
	}{
		"fail":      {isErr: true},
		"found":     {isErr: false, id: "allan-4", ok: true},
		"not found": {isErr: false, id: "empty", ok: false},
	}

	for _, t := range tests {
		ok, err := lt.IsCacheLimit(context.Background(), t.id)
		assert.Equal(t.isErr, err != nil)
		if err == nil {
			assert.Equal(t.ok, ok)
		}
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
