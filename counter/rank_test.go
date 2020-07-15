package counter

import (
	"testing"

	"math/rand"

	"time"

	"github.com/alicebob/miniredis"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestDb_IncreaseRankOfDaily(t *testing.T) {
	assert := assert.New(t)

	s, err := miniredis.Run()
	assert.NoError(err)
	defer s.Close()

	counter, err := NewCounter(WithRedisOption([]string{s.Addr()}))
	assert.NoError(err)

	now := time.Now()
	_, err = counter.IncreaseRankOfDaily("", "", time.Time{})
	assert.Error(err)

	group := "github.com"
	ids := []string{"allan", "gjbae1212", "dong", "jung"}
	values := make([]int64, 4)

	for i := 0; i < 1000; i++ {
		rnd := rand.Int() % 4
		v, err := counter.IncreaseRankOfDaily(group, ids[rnd], now)
		assert.NoError(err)
		assert.Equal(ids[rnd], v.Name)
		assert.Equal(values[rnd]+1, v.Value)
		values[rnd] = v.Value
	}
}

func TestDb_IncreaseRankOfTotal(t *testing.T) {
	assert := assert.New(t)

	s, err := miniredis.Run()
	assert.NoError(err)
	defer s.Close()

	counter, err := NewCounter(WithRedisOption([]string{s.Addr()}))
	assert.NoError(err)

	_, err = counter.IncreaseRankOfTotal("", "")
	assert.Error(err)

	group := "github.com"
	ids := []string{"allan", "gjbae1212", "dong", "jung"}
	values := make([]int64, 4)

	for i := 0; i < 1000; i++ {
		rnd := rand.Int() % 4
		v, err := counter.IncreaseRankOfTotal(group, ids[rnd])
		assert.NoError(err)
		assert.Equal(ids[rnd], v.Name)
		assert.Equal(values[rnd]+1, v.Value)
		values[rnd] = v.Value
	}
}

func TestDb_GetRankDailyByLimit(t *testing.T) {
	assert := assert.New(t)

	s, err := miniredis.Run()
	assert.NoError(err)
	defer s.Close()

	counter, err := NewCounter(WithRedisOption([]string{s.Addr()}))
	assert.NoError(err)

	group := "github.com"
	ids := []string{"allan", "gjbae1212", "dong", "jung"}
	values := make([]int64, 4)
	now := time.Now()

	for i := 0; i < 1000; i++ {
		rnd := rand.Int() % 4
		v, err := counter.IncreaseRankOfDaily(group, ids[rnd], now)
		assert.NoError(err)
		assert.Equal(ids[rnd], v.Name)
		assert.Equal(values[rnd]+1, v.Value)
		values[rnd] = v.Value
	}

	_, err = counter.GetRankDailyByLimit("", 0, time.Time{})
	assert.Error(err)

	scores, err := counter.GetRankDailyByLimit("empty", 10, now)
	assert.NoError(err)
	assert.Len(scores, 0)

	scores, err = counter.GetRankDailyByLimit(group, 2, now)
	assert.NoError(err)
	assert.Len(scores, 2)
	for _, s := range scores {
		for i, id := range ids {
			if id == s.Name {
				assert.Equal(values[i], s.Value)
			}
		}
	}

	scores, err = counter.GetRankDailyByLimit(group, 100, now)
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

	s, err := miniredis.Run()
	assert.NoError(err)
	defer s.Close()

	counter, err := NewCounter(WithRedisOption([]string{s.Addr()}))
	assert.NoError(err)

	group := "github.com"
	ids := []string{"allan", "gjbae1212", "dong", "jung"}
	values := make([]int64, 4)

	for i := 0; i < 1000; i++ {
		rnd := rand.Int() % 4
		v, err := counter.IncreaseRankOfTotal(group, ids[rnd])
		assert.NoError(err)
		assert.Equal(ids[rnd], v.Name)
		assert.Equal(values[rnd]+1, v.Value)
		values[rnd] = v.Value
	}

	_, err = counter.GetRankTotalByLimit("", 0)
	assert.Error(err)

	scores, err := counter.GetRankTotalByLimit("empty", 10)
	assert.NoError(err)
	assert.Len(scores, 0)

	scores, err = counter.GetRankTotalByLimit(group, 2)
	assert.NoError(err)
	assert.Len(scores, 2)
	for _, s := range scores {
		for i, id := range ids {
			if id == s.Name {
				assert.Equal(values[i], s.Value)
			}
		}
	}

	scores, err = counter.GetRankTotalByLimit(group, 100)
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
