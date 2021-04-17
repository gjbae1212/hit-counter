package handler

import (
	"net/http/httptest"
	"testing"

	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_WebSocket(t *testing.T) {
	assert := assert.New(t)
	defer mockRedis.FlushAll()

	e := echo.New()
	h, err := NewHandler(mockRedis.Addr())
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
