package main

import (
	"log"
	"reflect"
	"testing"

	"github.com/gjbae1212/go-module/logger"
	echo "github.com/labstack/echo/v4"

	"net/http/httptest"

	"github.com/gjbae1212/hit-counter/handler"
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

func TestMiddleWare(t *testing.T) {
	assert := assert.New(t)
	e := echo.New()

	mockHandler := func(c echo.Context) error {
		log.Println("call????")
		return nil
	}

	r := httptest.NewRequest("GET", "http://localhost", nil)
	r.Header.Set(echo.HeaderXForwardedFor, "127.0.0.1")
	w := httptest.NewRecorder()
	hctx := &handler.HitCounterContext{Context: e.NewContext(r, w)}

	// group api
	funcs, err := middlewareGroupApi()
	assert.NoError(err)
	f := funcs[0]
	err = f(mockHandler)(hctx)
	assert.NoError(err)
	assert.NotNil(hctx.Get("uid"))
	assert.NoError(err)
	assert.NotEmpty(w.Header().Get("Set-Cookie"))
	log.Println(w.Header().Get("Set-Cookie"))
}
