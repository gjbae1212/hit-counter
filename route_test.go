package main

import (
	"testing"

	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAddRoute(t *testing.T) {
	assert := assert.New(t)
	err := AddRoute(nil)
	assert.Error(err)

	err = AddRoute(echo.New())
	assert.NoError(err)
}
