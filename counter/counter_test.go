package counter

import (
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
)

func TestWithRedisOption(t *testing.T) {
	assert := assert.New(t)

	s, err := miniredis.Run()
	assert.NoError(err)
	defer s.Close()

	d := &db{}
	opt := WithRedisOption([]string{s.Addr()})
	err = opt.apply(d)
	assert.NoError(err)
	assert.NotNil(d.redis)
}

func TestNewCounter(t *testing.T) {
	assert := assert.New(t)

	s, err := miniredis.Run()
	assert.NoError(err)
	defer s.Close()

	counter, err := NewCounter(WithRedisOption([]string{s.Addr()}))
	assert.NoError(err)
	assert.NotNil(counter.(*db).redis)
}
