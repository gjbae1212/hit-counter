package api_handler

import (
	"os"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/gjbae1212/hit-counter/handler"
	"github.com/stretchr/testify/assert"
)

var (
	mockRedis *miniredis.Miniredis
)

func TestNewHandler(t *testing.T) {
	assert := assert.New(t)
	_, err := NewHandler(nil)
	assert.Error(err)

	h := &handler.Handler{}
	_, err = NewHandler(h)
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
