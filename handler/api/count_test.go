package api_handler

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gjbae1212/go-badge"

	"github.com/labstack/echo/v4"

	"context"

	"fmt"

	"github.com/alicebob/miniredis"
	"github.com/davecgh/go-spew/spew"
	"github.com/gjbae1212/hit-counter/counter"
	"github.com/gjbae1212/hit-counter/handler"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestWebSocketMessage_GetMessage(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]struct {
		input string
		want  string
	}{"step1": {input: "hi", want: "hi"}}

	for _, v := range tests {
		wsm := &WebSocketMessage{Payload: []byte(v.input)}
		assert.Equal(wsm.Payload, wsm.GetMessage())
	}
}

func TestRankTask_Process(t *testing.T) {
	assert := assert.New(t)

	s, err := miniredis.Run()
	assert.NoError(err)
	defer s.Close()

	h, err := handler.NewHandler([]string{s.Addr()})
	assert.NoError(err)

	api, err := NewHandler(h)
	assert.NoError(err)

	tests := map[string]struct {
		input *RankTask
		wants []*counter.Score
	}{
		"not_github": {input: &RankTask{
			Counter:   h.Counter,
			Domain:    "allan.com",
			Path:      "aa/bb",
			CreatedAt: time.Now(),
		}, wants: []*counter.Score{
			nil,
			nil,
			&counter.Score{
				Name:  "allan.com",
				Value: 1,
			},
			&counter.Score{
				Name:  "allan.com",
				Value: 1,
			},
		}},
		"github": {input: &RankTask{
			Counter:   h.Counter,
			Domain:    "github.com",
			Path:      "gjbae1212/test",
			CreatedAt: time.Now(),
		}, wants: []*counter.Score{
			&counter.Score{
				Name:  "gjbae1212/test",
				Value: 1,
			},
			&counter.Score{
				Name:  "gjbae1212/test",
				Value: 1,
			},
			&counter.Score{
				Name:  "github.com",
				Value: 1,
			},
			&counter.Score{
				Name:  "github.com",
				Value: 1,
			},
		}},
	}

	ctx := context.Background()

	test := tests["not_github"]
	err = api.AsyncTask.AddTask(ctx, test.input)
	assert.NoError(err)
	time.Sleep(1 * time.Second)
	scores, err := api.Counter.GetRankDailyByLimit(test.input.Domain, 10, time.Now())
	assert.NoError(err)
	assert.Len(scores, 0)
	scores, err = api.Counter.GetRankTotalByLimit(test.input.Domain, 10)
	assert.NoError(err)
	assert.Len(scores, 0)
	scores, err = api.Counter.GetRankDailyByLimit("domain", 10, time.Now())
	assert.NoError(err)
	assert.Len(scores, 1)
	assert.True(cmp.Equal(test.wants[2], scores[0]))

	scores, err = api.Counter.GetRankTotalByLimit("domain", 10)
	assert.NoError(err)
	assert.Len(scores, 1)
	assert.True(cmp.Equal(test.wants[3], scores[0]))

	test = tests["github"]
	err = api.AsyncTask.AddTask(ctx, test.input)
	assert.NoError(err)
	time.Sleep(1 * time.Second)
	scores, err = api.Counter.GetRankDailyByLimit(test.input.Domain, 10, time.Now())
	assert.NoError(err)
	assert.Len(scores, 1)
	assert.True(cmp.Equal(test.wants[0], scores[0]))
	spew.Dump(scores)

	scores, err = api.Counter.GetRankTotalByLimit(test.input.Domain, 10)
	assert.NoError(err)
	assert.Len(scores, 1)
	assert.True(cmp.Equal(test.wants[1], scores[0]))
	spew.Dump(scores)

	scores, err = api.Counter.GetRankDailyByLimit("domain", 10, time.Now())
	assert.NoError(err)
	assert.Len(scores, 2)
	if cmp.Equal(test.wants[2], scores[0]) {
	} else if cmp.Equal(test.wants[2], scores[1]) {
	} else {
		assert.NoError(fmt.Errorf("error"))
	}

	scores, err = api.Counter.GetRankTotalByLimit("domain", 10)
	assert.NoError(err)
	assert.Len(scores, 2)
	if cmp.Equal(test.wants[3], scores[0]) {
	} else if cmp.Equal(test.wants[3], scores[1]) {
	} else {
		assert.NoError(fmt.Errorf("error"))
	}
}

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
