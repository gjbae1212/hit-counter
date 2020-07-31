package api_handler

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gjbae1212/go-badge"

	"github.com/labstack/echo/v4"

	"fmt"

	"github.com/alicebob/miniredis"
	"github.com/gjbae1212/hit-counter/handler"
	"github.com/stretchr/testify/assert"
)



func TestHandler_KeepCount(t *testing.T) {
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
	err = api.KeepCount(hctx)
	assert.Error(err)

	hctx.Set("ckid", "test")
	hctx.Set("host", "github.com")
	hctx.Set("path", "gjbae1212/hit-counter")
	hctx.Set("title", "hello")
	err = api.KeepCount(hctx)
	assert.NoError(err)
	assert.Equal(200, w.Code)
	text := fmt.Sprintf(badgeFormat, 0, 0)
	badge, err := badge.RenderBytes("hello", text, "#79c83d")
	assert.NoError(err)
	assert.Equal(string(badge), w.Body.String())
}

func TestHandler_IncrCount(t *testing.T) {
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
	err = api.IncrCount(hctx)
	assert.Error(err)

	for i := 1; i < 20; i++ {
		r := httptest.NewRequest("GET", "http://localhost:8080", nil)
		r.Header.Set("User-Agent", fmt.Sprintf("%d", i))
		w := httptest.NewRecorder()
		hctx := &handler.HitCounterContext{Context: e.NewContext(r, w)}
		hctx.Set("ckid", "test")
		hctx.Set("host", "github.com")
		hctx.Set("path", "gjbae1212/hit-counter")
		hctx.Set("title", "")
		err = api.IncrCount(hctx)
		assert.NoError(err)
		assert.Equal(200, w.Code)
		text := fmt.Sprintf(badgeFormat, i, i)
		badge, err := badge.RenderBytes("hits", text, "#79c83d")
		assert.NoError(err)
		assert.Equal(string(badge), w.Body.String())
	}
	time.Sleep(1 * time.Second)
	scores, err := api.Counter.GetRankDailyByLimit("github.com", 10, time.Now())
	assert.NoError(err)
	assert.Len(scores, 1)
	assert.Equal(int64(19), scores[0].Value)
	scores, err = api.Counter.GetRankTotalByLimit("github.com", 10)
	assert.NoError(err)
	assert.Len(scores, 1)
	assert.Equal(int64(19), scores[0].Value)
	scores, err = api.Counter.GetRankTotalByLimit("domain", 10)
	assert.NoError(err)
	assert.Len(scores, 1)
	assert.Equal(int64(19), scores[0].Value)

}
