package main

import (
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gjbae1212/hit-counter/handler"

	"github.com/alicebob/miniredis"
	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAddRoute(t *testing.T) {
	assert := assert.New(t)
	err := AddRoute(nil, "")
	assert.Error(err)

	s, err := miniredis.Run()
	assert.NoError(err)
	defer s.Close()

	err = AddRoute(echo.New(), s.Addr())
	assert.NoError(err)
}

func TestGroup(t *testing.T) {
	assert := assert.New(t)
	e := echo.New()

	mockHandler := func(c echo.Context) error {
		log.Println("call????")
		return nil
	}

	r := httptest.NewRequest("GET", "http://localhost?url=github.com", nil)
	r.Header.Set(echo.HeaderXForwardedFor, "127.0.0.1")
	w := httptest.NewRecorder()
	hctx := &handler.HitCounterContext{Context: e.NewContext(r, w)}

	// group api
	funcs, err := groupApiCount()
	assert.NoError(err)

	f := funcs[0]
	err = f(mockHandler)(hctx)
	assert.NoError(err)
	assert.NotNil(hctx.Get("host"))
	assert.NotNil(hctx.Get("path"))
	assert.NotNil(hctx.Get("title"))
}
