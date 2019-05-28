package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Wasm(t *testing.T) {
	assert := assert.New(t)

	e := echo.New()

	s, err := miniredis.Run()
	assert.NoError(err)
	defer s.Close()

	h, err := NewHandler([]string{s.Addr()})
	assert.NoError(err)

	r := httptest.NewRequest("GET", "http://localhost:8080", nil)
	w := httptest.NewRecorder()

	hctx := &HitCounterContext{Context: e.NewContext(r, w)}
	err = h.Wasm(hctx)
	assert.NoError(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal("application/wasm", resp.Header.Get("Content-Type"))
	assert.Equal("gzip", resp.Header.Get("Content-Encoding"))
}
