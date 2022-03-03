package counter

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gjbae1212/hit-counter/internal"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestDb_IncreaseHitOfDaily(t *testing.T) {
	assert := assert.New(t)
	defer mockRedis.FlushAll()

	ctx := context.Background()
	counter, err := NewCounter(WithRedisClient(mockClient))
	assert.NoError(err)

	_, err = counter.IncreaseHitOfDaily(ctx, "", time.Time{}, time.Minute)
	assert.Error(err)

	now := time.Now()
	for i := 0; i < 2; i++ {
		count, err := counter.IncreaseHitOfDaily(ctx, "test", now, time.Minute)
		assert.NoError(err)
		assert.Equal(i+1, int(count.Value))
	}

	daily := internal.TimeToDailyStringFormat(now)
	key := fmt.Sprintf(hitDailyFormat, daily, "test")
	log.Println(key)

	v, err := counter.(*db).redisClient.Get(ctx, key).Result()
	assert.NoError(err)
	assert.Equal("2", v)
}

func TestDb_IncreaseHitOfTotal(t *testing.T) {
	assert := assert.New(t)
	defer mockRedis.FlushAll()

	ctx := context.Background()
	counter, err := NewCounter(WithRedisClient(mockClient))
	assert.NoError(err)

	_, err = counter.IncreaseHitOfTotal(ctx, "")
	assert.Error(err)

	for i := 0; i < 2; i++ {
		count, err := counter.IncreaseHitOfTotal(ctx, "test")
		assert.NoError(err)
		assert.Equal(i+1, int(count.Value))
	}

	key := fmt.Sprintf(hitTotalFormat, "test")
	log.Println(key)
	v, err := counter.(*db).redisClient.Get(ctx, key).Result()
	assert.NoError(err)
	assert.Equal("2", v)
}

func TestDb_GetHitOfDaily(t *testing.T) {
	assert := assert.New(t)
	defer mockRedis.FlushAll()

	ctx := context.Background()
	counter, err := NewCounter(WithRedisClient(mockClient))
	assert.NoError(err)

	now := time.Now()
	_, err = counter.GetHitOfDaily(ctx, "", time.Time{})
	assert.Error(err)

	v, err := counter.GetHitOfDaily(ctx, "empty", now)
	assert.NoError(err)
	assert.Nil(v)

	for i := 0; i < 1000; i++ {
		count, err := counter.IncreaseHitOfDaily(ctx, "test", now, time.Minute)
		assert.NoError(err)
		assert.Equal(i+1, int(count.Value))
	}

	v, err = counter.GetHitOfDaily(ctx, "test", now)
	assert.NoError(err)
	assert.Equal(1000, int(v.Value))
	assert.Equal(fmt.Sprintf(hitDailyFormat, internal.TimeToDailyStringFormat(now), "test"), v.Name)
}

func TestDb_GetHitOfTotal(t *testing.T) {
	assert := assert.New(t)
	defer mockRedis.FlushAll()

	ctx := context.Background()
	counter, err := NewCounter(WithRedisClient(mockClient))
	assert.NoError(err)

	_, err = counter.GetHitOfTotal(ctx, "")
	assert.Error(err)

	v, err := counter.GetHitOfTotal(ctx, "empty")
	assert.NoError(err)
	assert.Nil(v)

	for i := 0; i < 1000; i++ {
		count, err := counter.IncreaseHitOfTotal(ctx, "test")
		assert.NoError(err)
		assert.Equal(i+1, int(count.Value))
	}

	v, err = counter.GetHitOfTotal(ctx, "test")
	assert.NoError(err)
	assert.Equal(1000, int(v.Value))
	assert.Equal(fmt.Sprintf(hitTotalFormat, "test"), v.Name)
}

func TestDb_GetHitOfDailyAndTotal(t *testing.T) {
	assert := assert.New(t)
	defer mockRedis.FlushAll()

	ctx := context.Background()
	counter, err := NewCounter(WithRedisClient(mockClient))
	assert.NoError(err)

	id := "allan"
	now := time.Now()
	tests := map[string]struct {
		inputs []interface{}
		wants  []*Score
		err    bool
	}{
		"error1":    {inputs: []interface{}{"", time.Time{}}, wants: []*Score{nil, nil}, err: true},
		"error2":    {inputs: []interface{}{id, time.Time{}}, wants: []*Score{nil, nil}, err: true},
		"empty":     {inputs: []interface{}{id, now}, wants: []*Score{nil, nil}, err: false},
		"onlytotal": {inputs: []interface{}{"onlytotal", now}, wants: []*Score{nil, &Score{Name: fmt.Sprintf(hitTotalFormat, "onlytotal"), Value: 10}}, err: false},
		"onlydaily": {inputs: []interface{}{"onlydaily", now}, wants: []*Score{&Score{Name: fmt.Sprintf(hitDailyFormat, internal.TimeToDailyStringFormat(now), "onlydaily"), Value: 10}, nil}, err: false},
		"both":      {inputs: []interface{}{"both", now}, wants: []*Score{&Score{Name: fmt.Sprintf(hitDailyFormat, internal.TimeToDailyStringFormat(now), "both"), Value: 10}, &Score{Name: fmt.Sprintf(hitTotalFormat, "both"), Value: 10}}, err: false},
	}

	test := tests["error1"]
	_, _, err = counter.GetHitOfDailyAndTotal(ctx, test.inputs[0].(string), test.inputs[1].(time.Time))
	assert.Error(err)
	test = tests["error2"]
	_, _, err = counter.GetHitOfDailyAndTotal(ctx, test.inputs[0].(string), test.inputs[1].(time.Time))
	assert.Error(err)

	test = tests["empty"]
	daily, total, err := counter.GetHitOfDailyAndTotal(ctx, test.inputs[0].(string), test.inputs[1].(time.Time))
	assert.NoError(err)
	assert.Equal(test.wants[0], daily)
	assert.Equal(test.wants[1], total)

	for i := 0; i < 10; i++ {
		_, err := counter.IncreaseHitOfTotal(ctx, "onlytotal")
		assert.NoError(err)
	}
	test = tests["onlytotal"]
	daily, total, err = counter.GetHitOfDailyAndTotal(ctx, test.inputs[0].(string), test.inputs[1].(time.Time))
	assert.NoError(err)
	assert.True(cmp.Equal(test.wants[0], daily))
	assert.True(cmp.Equal(test.wants[1], total))

	for i := 0; i < 10; i++ {
		_, err := counter.IncreaseHitOfDaily(ctx, "onlydaily", now, time.Minute)
		assert.NoError(err)
	}

	test = tests["onlydaily"]
	daily, total, err = counter.GetHitOfDailyAndTotal(ctx, test.inputs[0].(string), test.inputs[1].(time.Time))
	assert.NoError(err)
	assert.True(cmp.Equal(test.wants[0], daily))
	assert.True(cmp.Equal(test.wants[1], total))

	for i := 0; i < 10; i++ {
		_, err := counter.IncreaseHitOfDaily(ctx, "both", now, time.Minute)
		assert.NoError(err)
		_, err = counter.IncreaseHitOfTotal(ctx, "both")
		assert.NoError(err)
	}

	test = tests["both"]
	daily, total, err = counter.GetHitOfDailyAndTotal(ctx, test.inputs[0].(string), test.inputs[1].(time.Time))
	assert.NoError(err)
	assert.True(cmp.Equal(test.wants[0], daily))
	assert.True(cmp.Equal(test.wants[1], total))
}

func TestDb_GetHitOfDailyByRange(t *testing.T) {
	assert := assert.New(t)
	defer mockRedis.FlushAll()

	ctx := context.Background()
	counter, err := NewCounter(WithRedisClient(mockClient))
	assert.NoError(err)

	_, err = counter.GetHitOfDailyByRange(ctx, "", []time.Time{})
	assert.Error(err)

	scores, err := counter.GetHitOfDailyByRange(ctx, "test.com", []time.Time{time.Now(), time.Now().Add(-1 * 24 * time.Hour)})
	assert.NoError(err)
	assert.Len(scores, 2)
	for _, s := range scores {
		assert.Nil(s)
	}

	var timeRange []time.Time
	prev := time.Now().Add(-30 * 24 * time.Hour)
	now := time.Now()
	for now.Unix() > prev.Unix() {
		timeRange = append(timeRange, prev)
		_, err := counter.IncreaseHitOfDaily(ctx, "test.com", prev, time.Minute)
		assert.NoError(err)
		prev = prev.Add(24 * time.Hour)
	}

	scores, err = counter.GetHitOfDailyByRange(ctx, "test.com", timeRange)
	assert.NoError(err)
	assert.Len(scores, 30)
	spew.Dump(scores)
	for _, score := range scores {
		assert.Equal(1, int(score.Value))
	}
}
