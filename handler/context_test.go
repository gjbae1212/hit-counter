package handler

import (
	"context"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/stretchr/testify/assert"
)

func TestHitCounterContext_ExtraLog(t *testing.T) {
	assert := assert.New(t)
	e := echo.New()
	r := httptest.NewRequest("GET", "http://localhost", nil)
	nctx := &HitCounterContext{e.NewContext(r, nil)}
	extraLog := nctx.ExtraLog()
	log.Println(extraLog)
	assert.Equal("GET", extraLog["method"])
	assert.Equal("localhost", extraLog["host"])
	assert.Len(extraLog, 6)
}

func TestHitCounterContext_ValueContext(t *testing.T) {
	assert := assert.New(t)

	e := echo.New()
	r := httptest.NewRequest("GET", "http://localhost", nil)
	nctx := &HitCounterContext{e.NewContext(r, nil)}

	nctx.WithContext("allan", "hi")
	value := nctx.ValueContext("allan")
	assert.Equal("hi", value.(string))
}

func TestHitCounterContext_WithContext(t *testing.T) {
	assert := assert.New(t)

	e := echo.New()
	r := httptest.NewRequest("GET", "http://localhost", nil)
	nctx := &HitCounterContext{e.NewContext(r, nil)}
	nctx.WithContext("allan", "hi")
	nctx.WithContext("test", "testhi")

	value := nctx.ValueContext("allan")
	assert.Equal("hi", value.(string))
}

func TestHitCounterContext_SetContext(t *testing.T) {
	assert := assert.New(t)

	e := echo.New()
	r := httptest.NewRequest("GET", "http://localhost", nil)
	nctx := &HitCounterContext{e.NewContext(r, nil)}

	ctx := context.WithValue(context.Background(), "test", "allan")
	nctx.SetContext(ctx)

	value := nctx.ValueContext("test")
	assert.Equal("allan", value.(string))
}

func TestHitCounterContext_GetContext(t *testing.T) {
	assert := assert.New(t)

	e := echo.New()
	r := httptest.NewRequest("GET", "http://localhost", nil)
	nctx := &HitCounterContext{e.NewContext(r, nil)}

	ctx := context.WithValue(context.Background(), "test", "allan")
	nctx.SetContext(ctx)

	vctx := nctx.GetContext()
	value := vctx.Value("test")
	assert.Equal("allan", value)
}
