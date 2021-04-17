package handler

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_HealthCheck(t *testing.T) {
	assert := assert.New(t)
	defer mockRedis.FlushAll()

	e := echo.New()
	h, err := NewHandler(mockRedis.Addr())
	assert.NoError(err)

	r := httptest.NewRequest("GET", "http://localhost:8080", nil)
	w := httptest.NewRecorder()

	ectx := e.NewContext(r, w)
	err = h.HealthCheck(ectx)
	assert.NoError(err)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal("health check!", string(body))
}
