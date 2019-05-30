package counter

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	allan_util "github.com/gjbae1212/go-module/util"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestDb_IncreaseHitOfDaily(t *testing.T) {
	assert := assert.New(t)

	s, err := miniredis.Run()
	assert.NoError(err)
	defer s.Close()

	counter, err := NewCounter(WithRedisOption([]string{s.Addr()}))
	assert.NoError(err)

	_, err = counter.IncreaseHitOfDaily("", time.Time{})
	assert.Error(err)

	now := time.Now()
	for i := 0; i < 2; i++ {
		count, err := counter.IncreaseHitOfDaily("test", now)
		assert.NoError(err)
		assert.Equal(i+1, int(count.Value))
	}

	daily := allan_util.TimeToDailyStringFormat(now)
	key := fmt.Sprintf(hitDailyFormat, daily, "test")
	log.Println(key)
	v, err := counter.(*db).redis.Do("GET", key)
	assert.NoError(err)
	assert.Equal("2", string(v.([]byte)))
}

func TestDb_IncreaseHitOfTotal(t *testing.T) {
	assert := assert.New(t)
	s, err := miniredis.Run()
	assert.NoError(err)
	defer s.Close()

	counter, err := NewCounter(WithRedisOption([]string{s.Addr()}))
	assert.NoError(err)

	_, err = counter.IncreaseHitOfTotal("")
	assert.Error(err)

	for i := 0; i < 2; i++ {
		count, err := counter.IncreaseHitOfTotal("test")
		assert.NoError(err)
		assert.Equal(i+1, int(count.Value))
	}

	key := fmt.Sprintf(hitTotalFormat, "test")
	log.Println(key)
	v, err := counter.(*db).redis.Do("GET", key)
	assert.NoError(err)
	assert.Equal("2", string(v.([]byte)))
}

func TestDb_GetHitOfDaily(t *testing.T) {
	assert := assert.New(t)
	s, err := miniredis.Run()
	assert.NoError(err)
	defer s.Close()

	counter, err := NewCounter(WithRedisOption([]string{s.Addr()}))
	assert.NoError(err)

	now := time.Now()

	_, err = counter.GetHitOfDaily("", time.Time{})
	assert.Error(err)

	v, err := counter.GetHitOfDaily("empty", now)
	assert.NoError(err)
	assert.Nil(v)

	for i := 0; i < 1000; i++ {
		count, err := counter.IncreaseHitOfDaily("test", now)
		assert.NoError(err)
		assert.Equal(i+1, int(count.Value))
	}

	v, err = counter.GetHitOfDaily("test", now)
	assert.NoError(err)
	assert.Equal(1000, int(v.Value))
}

func TestDb_GetHitOfTotal(t *testing.T) {
	assert := assert.New(t)

	s, err := miniredis.Run()
	assert.NoError(err)
	defer s.Close()

	counter, err := NewCounter(WithRedisOption([]string{s.Addr()}))
	assert.NoError(err)

	_, err = counter.GetHitOfTotal("")
	assert.Error(err)

	v, err := counter.GetHitOfTotal("empty")
	assert.NoError(err)
	assert.Nil(v)

	for i := 0; i < 1000; i++ {
		count, err := counter.IncreaseHitOfTotal("test")
		assert.NoError(err)
		assert.Equal(i+1, int(count.Value))
	}

	v, err = counter.GetHitOfTotal("test")
	assert.NoError(err)
	assert.Equal(1000, int(v.Value))
}

func TestDb_GetHitOfDailyAndTotal(t *testing.T) {
	assert := assert.New(t)

	s, err := miniredis.Run()
	assert.NoError(err)
	defer s.Close()

	counter, err := NewCounter(WithRedisOption([]string{s.Addr()}))
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
		"onlytotal": {inputs: []interface{}{"onlytotal", now}, wants: []*Score{nil, &Score{Name: "onlytotal", Value: 10}}, err: false},
		"onlydaily": {inputs: []interface{}{"onlydaily", now}, wants: []*Score{&Score{Name: "onlydaily", Value: 10}, nil}, err: false},
		"both":      {inputs: []interface{}{"both", now}, wants: []*Score{&Score{Name: "both", Value: 10}, &Score{Name: "both", Value: 10}}, err: false},
	}

	test := tests["error1"]
	_, _, err = counter.GetHitOfDailyAndTotal(test.inputs[0].(string), test.inputs[1].(time.Time))
	assert.Error(err)
	test = tests["error2"]
	_, _, err = counter.GetHitOfDailyAndTotal(test.inputs[0].(string), test.inputs[1].(time.Time))
	assert.Error(err)

	test = tests["empty"]
	daily, total, err := counter.GetHitOfDailyAndTotal(test.inputs[0].(string), test.inputs[1].(time.Time))
	assert.NoError(err)
	assert.Equal(test.wants[0], daily)
	assert.Equal(test.wants[1], total)

	for i := 0; i < 10; i++ {
		_, err := counter.IncreaseHitOfTotal("onlytotal")
		assert.NoError(err)
	}
	test = tests["onlytotal"]
	daily, total, err = counter.GetHitOfDailyAndTotal(test.inputs[0].(string), test.inputs[1].(time.Time))
	assert.NoError(err)
	assert.True(cmp.Equal(test.wants[0], daily))
	assert.True(cmp.Equal(test.wants[1], total))

	for i := 0; i < 10; i++ {
		_, err := counter.IncreaseHitOfDaily("onlydaily", now)
		assert.NoError(err)
	}

	test = tests["onlydaily"]
	daily, total, err = counter.GetHitOfDailyAndTotal(test.inputs[0].(string), test.inputs[1].(time.Time))
	assert.NoError(err)
	assert.True(cmp.Equal(test.wants[0], daily))
	assert.True(cmp.Equal(test.wants[1], total))

	for i := 0; i < 10; i++ {
		_, err := counter.IncreaseHitOfDaily("both", now)
		assert.NoError(err)
		_, err = counter.IncreaseHitOfTotal("both")
		assert.NoError(err)
	}

	test = tests["both"]
	daily, total, err = counter.GetHitOfDailyAndTotal(test.inputs[0].(string), test.inputs[1].(time.Time))
	assert.NoError(err)
	assert.True(cmp.Equal(test.wants[0], daily))
	assert.True(cmp.Equal(test.wants[1], total))
}
