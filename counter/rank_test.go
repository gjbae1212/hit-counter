package counter

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestDb_IncreaseRankOfDaily(t *testing.T) {
	assert := assert.New(t)
	defer mockRedis.FlushAll()

	ctx := context.Background()
	counter, err := NewCounter(WithRedisClient(mockClient))
	assert.NoError(err)

	now := time.Now()
	_, err = counter.IncreaseRankOfDaily(ctx, "", "", time.Time{}, time.Minute)
	assert.Error(err)

	group := "github.com"
	ids := []string{"allan", "gjbae1212", "dong", "jung"}
	values := make([]int64, 4)

	for i := 0; i < 1000; i++ {
		rnd := rand.Int() % 4
		v, err := counter.IncreaseRankOfDaily(ctx, group, ids[rnd], now, time.Minute)
		assert.NoError(err)
		assert.Equal(ids[rnd], v.Name)
		assert.Equal(values[rnd]+1, v.Value)
		values[rnd] = v.Value
	}
}

func TestDb_IncreaseRankOfTotal(t *testing.T) {
	assert := assert.New(t)
	defer mockRedis.FlushAll()

	ctx := context.Background()
	counter, err := NewCounter(WithRedisClient(mockClient))
	assert.NoError(err)

	_, err = counter.IncreaseRankOfTotal(ctx, "", "")
	assert.Error(err)

	group := "github.com"
	ids := []string{"allan", "gjbae1212", "dong", "jung"}
	values := make([]int64, 4)

	for i := 0; i < 1000; i++ {
		rnd := rand.Int() % 4
		v, err := counter.IncreaseRankOfTotal(ctx, group, ids[rnd])
		assert.NoError(err)
		assert.Equal(ids[rnd], v.Name)
		assert.Equal(values[rnd]+1, v.Value)
		values[rnd] = v.Value
	}
}

func TestDb_GetRankDailyByLimit(t *testing.T) {
	assert := assert.New(t)
	defer mockRedis.FlushAll()

	ctx := context.Background()
	counter, err := NewCounter(WithRedisClient(mockClient))
	assert.NoError(err)

	group := "github.com"
	ids := []string{"allan", "gjbae1212", "dong", "jung"}
	values := make([]int64, 4)
	now := time.Now()

	for i := 0; i < 1000; i++ {
		rnd := rand.Int() % 4
		v, err := counter.IncreaseRankOfDaily(ctx, group, ids[rnd], now, time.Minute)
		assert.NoError(err)
		assert.Equal(ids[rnd], v.Name)
		assert.Equal(values[rnd]+1, v.Value)
		values[rnd] = v.Value
	}

	_, err = counter.GetRankDailyByLimit(ctx, "", 0, time.Time{})
	assert.Error(err)

	scores, err := counter.GetRankDailyByLimit(ctx, "empty", 10, now)
	assert.NoError(err)
	assert.Len(scores, 0)

	scores, err = counter.GetRankDailyByLimit(ctx, group, 2, now)
	assert.NoError(err)
	assert.Len(scores, 2)
	for _, s := range scores {
		for i, id := range ids {
			if id == s.Name {
				assert.Equal(values[i], s.Value)
			}
		}
	}

	scores, err = counter.GetRankDailyByLimit(ctx, group, 100, now)
	assert.NoError(err)
	assert.Len(scores, 4)
	for _, s := range scores {
		spew.Dump(s)
		for i, id := range ids {
			if id == s.Name {
				assert.Equal(values[i], s.Value)
			}
		}
	}
}

func TestDb_GetRankTotalByLimit(t *testing.T) {
	assert := assert.New(t)
	defer mockRedis.FlushAll()

	ctx := context.Background()
	counter, err := NewCounter(WithRedisClient(mockClient))
	assert.NoError(err)

	group := "github.com"
	ids := []string{"allan", "gjbae1212", "dong", "jung"}
	values := make([]int64, 4)

	for i := 0; i < 1000; i++ {
		rnd := rand.Int() % 4
		v, err := counter.IncreaseRankOfTotal(ctx, group, ids[rnd])
		assert.NoError(err)
		assert.Equal(ids[rnd], v.Name)
		assert.Equal(values[rnd]+1, v.Value)
		values[rnd] = v.Value
	}

	_, err = counter.GetRankTotalByLimit(ctx, "", 0)
	assert.Error(err)

	scores, err := counter.GetRankTotalByLimit(ctx, "empty", 10)
	assert.NoError(err)
	assert.Len(scores, 0)

	scores, err = counter.GetRankTotalByLimit(ctx, group, 2)
	assert.NoError(err)
	assert.Len(scores, 2)
	for _, s := range scores {
		for i, id := range ids {
			if id == s.Name {
				assert.Equal(values[i], s.Value)
			}
		}
	}

	scores, err = counter.GetRankTotalByLimit(ctx, group, 100)
	assert.NoError(err)
	assert.Len(scores, 4)
	for _, s := range scores {
		spew.Dump(s)
		for i, id := range ids {
			if id == s.Name {
				assert.Equal(values[i], s.Value)
			}
		}
	}
}
