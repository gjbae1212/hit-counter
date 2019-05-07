package counter

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	allan_util "github.com/gjbae1212/go-module/util"
	"github.com/stretchr/testify/assert"
)

func TestDb_IncreaseHitOfDaily(t *testing.T) {
	assert := assert.New(t)

	s, err := miniredis.Run()
	assert.NoError(err)
	defer s.Close()

	counter, err := NewCounter(WithRedisOption([]string{s.Addr()}))
	assert.NoError(err)

	_, err = counter.IncreaseHitOfDaily("")
	assert.Error(err)

	for i := 0; i < 2; i++ {
		count, err := counter.IncreaseHitOfDaily("test")
		assert.NoError(err)
		assert.Equal(i+1, int(count.Value))
	}

	daily := allan_util.TimeToDailyStringFormat(time.Now())
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

	_, err = counter.GetHitOfDaily("")
	assert.Error(err)

	v, err := counter.GetHitOfDaily("empty")
	assert.NoError(err)
	assert.Nil(v)

	for i := 0; i < 1000; i++ {
		count, err := counter.IncreaseHitOfDaily("test")
		assert.NoError(err)
		assert.Equal(i+1, int(count.Value))
	}

	v, err = counter.GetHitOfDaily("test")
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
