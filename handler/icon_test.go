package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_IconAll(t *testing.T) {
	assert := assert.New(t)

	e := echo.New()

	s, err := miniredis.Run()
	assert.NoError(err)
	defer s.Close()

	h, err := NewHandler([]string{s.Addr()})
	assert.NoError(err)

	tests := map[string]struct {
	}{
		"success": {},
	}

	for _, _ = range tests {
		r := httptest.NewRequest("GET", "http://localhost:8080/icon/all.json", nil)
		w := httptest.NewRecorder()
		ctx := &HitCounterContext{Context: e.NewContext(r, w)}
		err := h.IconAll(ctx)
		assert.NoError(err)

		resp := w.Result()
		assert.Equal(http.StatusOK, resp.StatusCode)
		raw, err := ioutil.ReadAll(resp.Body)
		assert.NoError(err)
		var rm []interface{}
		err = json.Unmarshal(raw, &rm)
		assert.NoError(err)
		assert.Len(rm, len(h.IconsList))
	}

}

func TestHandler_Icon(t *testing.T) {
	assert := assert.New(t)
	e := echo.New()

	s, err := miniredis.Run()
	assert.NoError(err)
	defer s.Close()

	h, err := NewHandler([]string{s.Addr()})
	assert.NoError(err)

	tests := map[string]struct {
		status int
		path   string
	}{
		"empty": {
			status: 404,
			path:   "empty.svg",
		},
		"success": {
			status: 200,
			path:   "github.svg",
		},
	}

	for _, t := range tests {
		r := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:8080/icon/%s", t.path), nil)
		w := httptest.NewRecorder()
		ctx := &HitCounterContext{Context: e.NewContext(r, w)}
		ctx.SetParamNames("icon")
		ctx.SetParamValues(t.path)
		err := h.Icon(ctx)
		assert.NoError(err)

		resp := w.Result()

		assert.Equal(t.status, resp.StatusCode)
		if resp.StatusCode == 200 {
			raw, err := ioutil.ReadAll(resp.Body)
			assert.NoError(err)
			assert.Equal(string(h.Icons[t.path].Origin), string(raw))
		}
	}
}
