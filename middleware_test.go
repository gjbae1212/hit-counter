package main

import (
	"testing"

	"github.com/gjbae1212/hit-counter/internal"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAddMiddleware(t *testing.T) {
	assert := assert.New(t)

	err := AddMiddleware(nil)
	assert.Error(err)

	e := echo.New()
	err = AddMiddleware(e, WithDebugOption(false))
	assert.NoError(err)
	assert.False(e.Debug)

	clogger, err := internal.NewLogger("", "")
	assert.NoError(err)

	err = AddMiddleware(e, WithLogger(clogger))
	assert.NoError(err)
	assert.True(e.Debug)
	assert.Equal(clogger, e.Logger)
}
