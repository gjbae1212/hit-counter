package main

import (
	"log"
	"reflect"
	"testing"

	"github.com/gjbae1212/go-module/logger"
	"github.com/labstack/echo/v4"

	"github.com/stretchr/testify/assert"
)

func TestWithDebugOption(t *testing.T) {
	assert := assert.New(t)
	e := echo.New()
	opt := WithDebugOption(false)
	opt.apply(e)
	assert.False(e.Debug)

	opt = WithDebugOption(true)
	opt.apply(e)
	assert.True(e.Debug)
}

func TestWithLoggerOption(t *testing.T) {
	assert := assert.New(t)
	e := echo.New()
	opt := WithLoggerOption(nil)
	opt.apply(e)
	assert.IsType(reflect.TypeOf(&log.Logger{}), reflect.TypeOf(e.Logger))

	clogger, err := logger.NewLogger("", "")
	assert.NoError(err)
	opt = WithLoggerOption(clogger)
	opt.apply(e)
	assert.Equal(clogger, e.Logger)
}

func TestAddMiddleware(t *testing.T) {
	assert := assert.New(t)

	err := AddMiddleware(nil)
	assert.Error(err)

	e := echo.New()
	err = AddMiddleware(e, WithDebugOption(false))
	assert.NoError(err)
	assert.False(e.Debug)

	clogger, err := logger.NewLogger("", "")
	assert.NoError(err)

	err = AddMiddleware(e, WithLoggerOption(clogger))
	assert.NoError(err)
	assert.True(e.Debug)
	assert.Equal(clogger, e.Logger)
}


