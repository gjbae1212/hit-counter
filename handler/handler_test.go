package handler

import (
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
)

func TestNewHandler(t *testing.T) {
	assert := assert.New(t)

	_, err := NewHandler(nil)
	assert.Error(err)

	s, err := miniredis.Run()
	assert.NoError(err)
	defer s.Close()

	_, err = NewHandler([]string{s.Addr()})
	assert.NoError(err)
}
