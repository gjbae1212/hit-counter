package handler

import (
	"os"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
)

var (
	mockRedis *miniredis.Miniredis
)

func TestNewHandler(t *testing.T) {
	assert := assert.New(t)
	defer mockRedis.FlushAll()

	_, err := NewHandler("")
	assert.Error(err)

	_, err = NewHandler(mockRedis.Addr())
	assert.NoError(err)
}

func TestMain(m *testing.M) {
	var err error
	mockRedis, err = miniredis.Run()
	if err != nil {
		panic(err)
	}
	code := m.Run()
	mockRedis.Close()
	os.Exit(code)
}
