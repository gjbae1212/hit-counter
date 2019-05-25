package api_handler

import (
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"

	"net/http/httptest"

	"encoding/json"

	"github.com/gjbae1212/hit-counter/handler"
	"github.com/labstack/echo/v4"
)

func TestHandler_GithubRankOfTotal(t *testing.T) {
	assert := assert.New(t)

	s, err := miniredis.Run()
	assert.NoError(err)
	defer s.Close()

	h, err := handler.NewHandler([]string{s.Addr()})
	assert.NoError(err)

	api, err := NewHandler(h)
	assert.NoError(err)

	e := echo.New()
	r := httptest.NewRequest("GET", "http://localhost:8080", nil)
	w := httptest.NewRecorder()
	hctx := &handler.HitCounterContext{Context: e.NewContext(r, w)}
	err = api.GithubRankOfTotal(hctx)
	assert.NoError(err)
	var ranks []string
	err = json.Unmarshal(w.Body.Bytes(), &ranks)
	assert.NoError(err)
	assert.Len(ranks, 0)

	for i := 0; i < 10; i++ {
		_, err := api.Counter.IncreaseRankOfTotal("github.com", "/allan/hi")
		assert.NoError(err)
	}

	for i := 0; i < 5; i++ {
		_, err := api.Counter.IncreaseRankOfTotal("github.com", "/allan/hello")
		assert.NoError(err)
	}


	r = httptest.NewRequest("GET", "http://localhost:8080", nil)
	w = httptest.NewRecorder()
	hctx = &handler.HitCounterContext{Context: e.NewContext(r, w)}
	err = api.GithubRankOfTotal(hctx)
	assert.NoError(err)
	ranks = []string{}
	err = json.Unmarshal(w.Body.Bytes(), &ranks)
	assert.NoError(err)
	assert.Len(ranks, 2)
}
