package api_handler

import (
	"testing"

	"github.com/gjbae1212/hit-counter/handler"
	"github.com/stretchr/testify/assert"
)

func TestNewHandler(t *testing.T) {
	assert := assert.New(t)
	_, err := NewHandler(nil)
	assert.Error(err)

	h := &handler.Handler{}
	_, err = NewHandler(h)
	assert.NoError(err)
}
