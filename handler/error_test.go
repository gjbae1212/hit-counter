package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Error(t *testing.T) {
	assert := assert.New(t)
	defer mockRedis.FlushAll()

	e := echo.New()
	h, err := NewHandler(mockRedis.Addr())
	assert.NoError(err)

	request := httptest.NewRequest("GET", "http://localhost", nil)
	w := httptest.NewRecorder()
	ectx := e.NewContext(request, w)

	h.Error(fmt.Errorf("[err] test"), ectx)

	resp := w.Result()
	assert.Equal(http.StatusInternalServerError, resp.StatusCode)
}
