package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis"
	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_WebSocket(t *testing.T) {
	assert := assert.New(t)

	e := echo.New()

	s, err := miniredis.Run()
	assert.NoError(err)
	defer s.Close()

	h, err := NewHandler([]string{s.Addr()})
	assert.NoError(err)

	r := httptest.NewRequest("GET", "http://localhost:8080", nil)
	r.Header.Set("Connection", "upgrade")
	r.Header.Set("Upgrade", "websocket")
	r.Header.Set("Sec-Websocket-Version", "13")
	r.Header.Set("Sec-WebSocket-Key", "allan")
	w := httptest.NewRecorder()
	hctx := &HitCounterContext{Context: e.NewContext(r, w)}
	//err = h.WebSocket(hctx)
	_ = hctx
	_ = h
}
