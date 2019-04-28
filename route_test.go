package main

import (
	"testing"

	"github.com/alicebob/miniredis"
	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAddRoute(t *testing.T) {
	assert := assert.New(t)
	err := AddRoute(nil, nil, 0)
	assert.Error(err)

	s, err := miniredis.Run()
	assert.NoError(err)
	defer s.Close()

	err = AddRoute(echo.New(), []string{s.Addr()}, 10)
	assert.NoError(err)
}
